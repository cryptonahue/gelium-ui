package web

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

//go:embed styles/*.css
var sourceStyles embed.FS

func sourceAppCSS(t *testing.T) string {
	t.Helper()
	// Canonical import order of web/styles/app.css. The entry file only
	// imports the split files, so raw-source assertions concatenate them in
	// that exact order plus the entry's own top-level tail (keyframes,
	// reduced-motion, forced-colors). Keep this list in sync with app.css.
	paths := []string{
		"styles/base.css",
		"styles/button.css",
		"styles/text-field.css",
		"styles/dialog.css",
		"styles/toast.css",
		"styles/focus-ring.css",
		"styles/elevation.css",
		"styles/icon.css",
		"styles/divider.css",
		"styles/card.css",
		"styles/badge.css",
		"styles/checkbox.css",
		"styles/radio.css",
		"styles/switch.css",
		"styles/select.css",
		"styles/select-menu.css",
		"styles/slider.css",
		"styles/progress.css",
		"styles/list.css",
		"styles/app.css",
	}
	var sb strings.Builder
	for _, path := range paths {
		css, err := sourceStyles.ReadFile(path)
		if err != nil {
			t.Fatalf("read source app CSS %s: %v", path, err)
		}
		sb.Write(css)
	}
	return sb.String()
}

func repositoryFile(t *testing.T, path ...string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate styles contract test")
	}
	parts := append([]string{filepath.Dir(filename), ".."}, path...)
	content, err := os.ReadFile(filepath.Join(parts...))
	if err != nil {
		t.Fatalf("read repository file: %v", err)
	}
	return string(content)
}
