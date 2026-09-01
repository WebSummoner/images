package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	geckoDriverBinary = "geckodriver"
)

type Firefox struct {
	Requirements
}

func (c *Firefox) Build() error {
	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	devSrcDir := "firefox/apt"
	pkgSrcPath, pkgVersion, err := c.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		devSrcDir = "firefox/local"
		pkgDestDir := filepath.Join(devDestDir, devSrcDir)
		err := os.MkdirAll(pkgDestDir, 0755)
		if err != nil {
			return fmt.Errorf("create %v temporary dir: %v", pkgDestDir, err)
		}
		pkgDestPath := filepath.Join(pkgDestDir, "firefox.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	pkgTagVersion := extractVersion(pkgVersion)
	devImageTag := fmt.Sprintf("websummoner/dev_firefox:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: c.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(devSrcDir, devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devBuildArgs := []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
	devBuildArgs = append(devBuildArgs, c.channelToBuildArgs()...)
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

	image, err := NewImage("firefox/geckodriver", destDir, c.Requirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	firefoxMajorMinorVersion := majorMinorVersion(pkgTagVersion)
	driverVersion, err := c.downloadGeckoDriver(image.Dir)
	if err != nil {
		return fmt.Errorf("failed to download geckodriver: %v", err)
	}
	// geckodriver is served directly: the hub handles file upload, session
	// timeouts and cleanup from outside the container, so nothing needs to
	// run alongside the driver here.
	image.Labels = []string{fmt.Sprintf("driver=geckodriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(c.TestsDir, "firefox", firefoxMajorMinorVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (c *Firefox) channelToBuildArgs() []string {
	switch c.BrowserChannel {
	case "beta":
		return []string{"PACKAGE=firefox-beta"}
	case "dev":
		return []string{"PACKAGE=firefox-nightly"}
	case "esr":
		return []string{"PACKAGE=firefox-esr"}
	default:
		return []string{}
	}
}

func (c *Firefox) downloadGeckoDriver(dir string) (string, error) {
	version := c.DriverVersion
	if version == LatestVersion {
		v, err := latestGithubRelease("mozilla/geckodriver")
		if err != nil {
			return "", fmt.Errorf("latest geckodriver version: %v", err)
		}
		version = strings.TrimPrefix(v, "v")
	}

	u := fmt.Sprintf("https://github.com/mozilla/geckodriver/releases/download/v%s/geckodriver-v%s-linux64.tar.gz", version, version)
	_, err := downloadDriver(u, geckoDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download geckodriver: %v", err)
	}
	return version, nil
}
