package build

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	operaDriverBinary = "operadriver_linux64/operadriver"

	// Opera N is built on Chromium N+16.
	operaChromiumOffset = 16
)

type Opera struct {
	Requirements
}

func (o *Opera) Build() error {

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "opera/apt"
	pkgSrcPath, pkgVersion, err := o.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		srcDir = "opera/local"
		pkgDestDir := filepath.Join(devDestDir, srcDir)
		err := os.MkdirAll(pkgDestDir, 0755)
		if err != nil {
			return fmt.Errorf("create %v temporary dir: %v", pkgDestDir, err)
		}
		pkgDestPath := filepath.Join(pkgDestDir, "opera.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	pkgTagVersion := extractVersion(pkgVersion)
	devImageTag := fmt.Sprintf("websummoner/dev_opera:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: o.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(srcDir, devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devBuildArgs := []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
	devBuildArgs = append(devBuildArgs, o.channelToBuildArgs()...)
	devImage.BuildArgs = devBuildArgs
	if pkgSrcPath != "" {
		devImage.FileServer = true
	}

	err = devImage.Build()
	if err != nil {
		return fmt.Errorf("build dev image: %v", err)
	}

	// Build main image
	destDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create temporary dir: %v", err)
	}

	image, err := NewImage("opera", destDir, o.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	driverVersion, driverName, err := o.downloadOperaDriver(image.Dir, pkgTagVersion)
	if err != nil {
		return fmt.Errorf("failed to download operadriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=%s:%s", driverName, driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(o.TestsDir, "opera", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (o *Opera) channelToBuildArgs() []string {
	switch o.BrowserChannel {
	case "beta":
		return []string{"PACKAGE=opera-beta"}
	case "dev":
		return []string{"PACKAGE=opera-developer"}
	default:
		return []string{}
	}
}

// downloadOperaDriver puts a driver for this Opera build into dir.
//
// operachromiumdriver tags follow the Chromium version Opera is built on, not
// Opera's own line: Opera N ships Chromium N+16. Matching on Opera's line gets
// a driver many majors too old, which refuses every session.
//
// Opera publishes late, so fall back to the newest operadriver rather than a
// chromedriver — the version check warns, it does not refuse.
//
// Returns the driver version and its name, for the image label.
func (o *Opera) downloadOperaDriver(dir string, browserVersion string) (string, string, error) {
	if o.DriverVersion != LatestVersion && o.DriverVersion != "" {
		u := fmt.Sprintf("https://github.com/operasoftware/operachromiumdriver/releases/download/v.%s/operadriver_linux64.zip", o.DriverVersion)
		if _, err := downloadDriver(u, operaDriverBinary, dir); err != nil {
			return "", "", fmt.Errorf("download Operadriver %s: %v", o.DriverVersion, err)
		}
		return o.DriverVersion, "operadriver", nil
	}

	chromiumMajor, err := operaChromiumMajor(browserVersion)
	if err != nil {
		return "", "", err
	}

	// Prefer Opera's own driver when one exists for this Chromium line.
	if tag, err := latestGithubReleaseWithPrefix("operasoftware/operachromiumdriver", fmt.Sprintf("v.%d.", chromiumMajor)); err == nil && tag != "" {
		version := strings.TrimPrefix(tag, "v.")
		u := fmt.Sprintf("https://github.com/operasoftware/operachromiumdriver/releases/download/%s/operadriver_linux64.zip", tag)
		if _, err := downloadDriver(u, operaDriverBinary, dir); err == nil {
			return version, "operadriver", nil
		}
	}

	// No driver on this line yet. The newest operadriver still drives it — the
	// version check only warns. A chromedriver starts but crashes on window.open.
	tag, err := latestGithubReleaseWithPrefix("operasoftware/operachromiumdriver", "v.")
	if err != nil {
		return "", "", fmt.Errorf("look up latest operadriver: %v", err)
	}
	if tag == "" {
		return "", "", fmt.Errorf("no operadriver release found for Chromium %d", chromiumMajor)
	}
	version := strings.TrimPrefix(tag, "v.")
	u := fmt.Sprintf("https://github.com/operasoftware/operachromiumdriver/releases/download/%s/operadriver_linux64.zip", tag)
	if _, err := downloadDriver(u, operaDriverBinary, dir); err != nil {
		return "", "", fmt.Errorf("download operadriver %s: %v", version, err)
	}
	return version, "operadriver", nil
}

// operaChromiumMajor maps an Opera release to the Chromium major it ships.
func operaChromiumMajor(operaVersion string) (int, error) {
	major, err := strconv.Atoi(majorVersion(operaVersion))
	if err != nil {
		return 0, fmt.Errorf("parse Opera version %q: %v", operaVersion, err)
	}
	return major + operaChromiumOffset, nil
}

// latestGithubReleaseWithPrefix returns the newest release tag of repo starting
// with prefix, or an empty string when the repository has none.
func latestGithubReleaseWithPrefix(repo string, prefix string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=100", repo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}
	for _, r := range releases {
		if strings.HasPrefix(r.TagName, prefix) {
			return r.TagName, nil
		}
	}
	return "", nil
}
