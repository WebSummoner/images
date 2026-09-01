package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	hv "github.com/hashicorp/go-version"
)

// Brave builds a Brave browser image. Brave is Chromium-based, so the
// WebDriver binary is ChromeDriver for the Chromium version Brave embeds —
// detected from the dev image itself, not guessed from the Brave version.
type Brave struct {
	Requirements
}

func (b *Brave) Build() error {
	pkgSrcPath, pkgVersion, err := b.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	pkgTagVersion := extractVersion(pkgVersion)

	// Build dev image (Brave + deps)
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "brave/apt"
	if pkgSrcPath != "" {
		srcDir = "brave/local"
		pkgDestDir := filepath.Join(devDestDir, srcDir)
		if err := os.MkdirAll(pkgDestDir, 0755); err != nil {
			return fmt.Errorf("create %v temporary dir: %v", pkgDestDir, err)
		}
		pkgDestPath := filepath.Join(pkgDestDir, "brave-browser.deb")
		if err := os.Rename(pkgSrcPath, pkgDestPath); err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	devImageTag := fmt.Sprintf("websummoner/dev_brave:%s", pkgTagVersion)
	devImage, err := NewImage(srcDir, devDestDir, Requirements{NoCache: b.NoCache, Tags: []string{devImageTag}})
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devImage.BuildArgs = []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
	if pkgSrcPath != "" {
		devImage.FileServer = true
	}
	if err := devImage.Build(); err != nil {
		return fmt.Errorf("build dev image: %v", err)
	}

	// Detect the Chromium major version Brave is built on
	chromiumMajor, err := chromiumMajorOf(devImageTag)
	if err != nil {
		return err
	}

	chromeDriverVersions, err := fetchChromeDriverVersions()
	if err != nil {
		return fmt.Errorf("fetch chromedriver versions: %v", err)
	}

	driverVersion, err := latestChromeDriverForMajor(chromiumMajor, chromeDriverVersions)
	if err != nil {
		return err
	}

	// Build main image
	destDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create temporary dir: %v", err)
	}

	image, err := NewImage("brave", destDir, b.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	if err := downloadChromeDriver(image.Dir, driverVersion, chromeDriverVersions); err != nil {
		return fmt.Errorf("failed to download chromedriver: %v", err)
	}
	image.Labels = []string{
		fmt.Sprintf("driver=chromedriver:%s", driverVersion),
		fmt.Sprintf("chromium=%s", chromiumMajor),
	}

	if err := image.Build(); err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	if err := image.Test(b.TestsDir, "brave", pkgTagVersion); err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	if err := image.Push(); err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

// chromiumMajorOf runs `brave-browser-stable --version` inside the dev image
// (the brave-browser wrapper script swallows the output) and extracts the
// Chromium major Brave embeds in its own version, e.g.
// "Brave Browser 152.1.94.117" -> "152".
func chromiumMajorOf(devImageTag string) (string, error) {
	out, _ := exec.Command("docker", "run", "--rm", "--entrypoint",
		"/usr/bin/brave-browser-stable", devImageTag, "--version").CombinedOutput()
	if len(strings.TrimSpace(string(out))) == 0 {
		return "", fmt.Errorf("empty brave version output for %s", devImageTag)
	}
	m := regexp.MustCompile(`Brave Browser\s*(\d+)\.`).FindStringSubmatch(string(out))
	if m == nil {
		return "", fmt.Errorf("cannot detect Chromium version from: %s", strings.TrimSpace(string(out)))
	}
	return m[1], nil
}

// latestChromeDriverForMajor picks the newest ChromeDriver version whose
// major matches Brave's Chromium major.
func latestChromeDriverForMajor(major string, versions map[string]string) (string, error) {
	var matching []string
	for v := range versions {
		if strings.HasPrefix(v, major+".") {
			matching = append(matching, v)
		}
	}
	if len(matching) == 0 {
		return "", fmt.Errorf("no chromedriver found for Chromium %s", major)
	}
	sort.SliceStable(matching, func(i, j int) bool {
		lv, lerr := hv.NewVersion(matching[i])
		rv, rerr := hv.NewVersion(matching[j])
		if lerr != nil || rerr != nil {
			return matching[i] > matching[j]
		}
		return lv.LessThan(rv)
	})
	return matching[len(matching)-1], nil
}
