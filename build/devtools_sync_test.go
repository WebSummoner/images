package build

import (
	"os"
	"path/filepath"
	"testing"
)

// Each Chromium image needs its own copy: a build context cannot reach outside
// itself. A copy that falls behind loses CDP and HAR only in the published image.
const canonicalDevtools = "../static/chrome/devtools"

var devtoolsCopies = []string{
	"../static/brave/devtools",
	"../static/edge/devtools",
	"../static/opera/devtools",
	"../static/yandex/devtools",
}

func TestDevtoolsCopiesMatchChrome(t *testing.T) {
	want := readTree(t, canonicalDevtools)
	for _, dir := range devtoolsCopies {
		got := readTree(t, dir)
		for name, content := range want {
			other, ok := got[name]
			if !ok {
				t.Errorf("%s is missing %s; copy it from %s", dir, name, canonicalDevtools)
				continue
			}
			if content != other {
				t.Errorf("%s/%s differs from %s/%s; the copies must stay identical", dir, name, canonicalDevtools, name)
			}
		}
		for name := range got {
			if _, ok := want[name]; !ok {
				t.Errorf("%s has an extra file %s", dir, name)
			}
		}
	}
}

func readTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		files[rel] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	if len(files) == 0 {
		t.Fatalf("%s has no files", dir)
	}
	return files
}
