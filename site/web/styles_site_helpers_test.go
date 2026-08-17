package web

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Site-side contract-test helpers. The lib package carries the component CSS
// suites (lib/styles_*_test.go) with their own sourceStyles embed; these
// helpers cover the docs-shell / demo / recipe CSS that stays in site/web.
// repositoryRoot arithmetic here is site/web → repo root (two levels up).

//go:embed styles/*.css
var siteStyles embed.FS

// repositoryRoot returns the repository root as seen from this test file.
func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate site contract test")
	}
	return filepath.Dir(filename) + "/../.."
}

// repositoryFile reads a repository-relative file (any depth below root).
func repositoryFile(t *testing.T, path ...string) string {
	t.Helper()
	p := filepath.Join(append([]string{repositoryRoot(t)}, path...)...)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read repository file %v: %v", path, err)
	}
	return string(b)
}

// availableThemes discovers the themes that actually exist on disk.
func availableThemes(t *testing.T) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repositoryRoot(t), "themes", "*", "theme.css"))
	if err != nil {
		t.Fatalf("discover themes: %v", err)
	}
	var themes []string
	for _, m := range matches {
		themes = append(themes, filepath.Base(filepath.Dir(m)))
	}
	return themes
}

// themeCSS reads one theme's CSS from the repository themes directory.
func themeCSS(t *testing.T, name string) string {
	t.Helper()
	found := false
	for _, n := range availableThemes(t) {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("theme %q not present on disk", name)
	}
	return repositoryFile(t, "themes", name, "theme.css")
}

// sourceComponentCSS reads a source CSS file. Site styles come from the local
// embed; component styles that moved into lib/ are read from the repository
// tree (base.css, focus-ring.css are lib-owned but consumed by shell tests).
func sourceComponentCSS(t *testing.T, name string) string {
	t.Helper()
	if b, err := siteStyles.ReadFile("styles/" + name); err == nil {
		return string(b)
	}
	libPath := filepath.Join(repositoryRoot(t), "lib", "styles", name)
	if b, err := os.ReadFile(libPath); err == nil {
		return string(b)
	}
	t.Fatalf("read site component CSS %s from embed or lib", name)
	return ""
}

// splitThemeSchemes splits one theme's CSS into its light block and its
// explicit dark-class block (site-side copy of the lib contract helper;
// the theme matrix suite moved to lib/, this is the docs-shell consumer).
func splitThemeSchemes(t *testing.T, theme string) (light, darkClass, darkMedia string) {
	t.Helper()
	css := themeCSS(t, theme)

	darkClassStart := strings.Index(css, ".theme-dark")
	if darkClassStart < 0 {
		t.Errorf("theme %q must declare an explicit dark class route (.theme-dark / .dark / [data-theme=\"dark\"])", theme)
		return "", "", ""
	}

	light = css[:darkClassStart]
	darkClass = css[darkClassStart:]
	if mediaStart := strings.Index(css, "@media (prefers-color-scheme: dark)"); mediaStart >= 0 {
		darkClass = css[darkClassStart:mediaStart]
		darkMedia = css[mediaStart:]
	}
	return light, darkClass, darkMedia
}
