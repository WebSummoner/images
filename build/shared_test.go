package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The build helpers resolve static/ from the repository root.
func fromRepoRoot(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

// Every image that runs the devtools helper builds it from the shared copy, so
// a Dockerfile COPY of it must resolve inside the generated context.
func TestSharedFilesReachEveryContext(t *testing.T) {
	fromRepoRoot(t)
	for _, browser := range []string{"chrome", "brave", "edge", "opera", "yandex"} {
		dest := t.TempDir()
		dir, err := copySourceFiles(browser, dest)
		if err != nil {
			t.Fatalf("%s: %v", browser, err)
		}
		if err := copySharedFiles(dir); err != nil {
			t.Fatalf("%s: %v", browser, err)
		}
		for _, name := range []string{"devtools/devtools.go", "devtools/go.mod"} {
			if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
				t.Errorf("%s: %s missing from the build context", browser, name)
			}
		}
	}
}

func TestSharedCopyIsIdenticalToSource(t *testing.T) {
	dest := t.TempDir()
	fromRepoRoot(t)
	dir, err := copySourceFiles("chrome", dest)
	if err != nil {
		t.Fatal(err)
	}
	if err := copySharedFiles(dir); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("static", sharedSourceDir, "devtools", "devtools.go"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "devtools", "devtools.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Error("the context copy differs from static/_shared")
	}
}

// Compiling the helper is not enough: the final stage must ship it too.
func TestEveryDevtoolsImageShipsTheBinary(t *testing.T) {
	fromRepoRoot(t)
	for _, browser := range []string{"chrome", "brave", "edge", "opera", "yandex"} {
		entrypoint, err := os.ReadFile(filepath.Join("static", browser, "entrypoint.sh"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(entrypoint), "/usr/bin/devtools") {
			continue
		}
		dockerfile, err := os.ReadFile(filepath.Join("static", browser, "Dockerfile"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(dockerfile), "COPY --from=go /devtools/devtools /usr/bin/") {
			t.Errorf("%s: entrypoint starts devtools but the image never copies it", browser)
		}
	}
}
