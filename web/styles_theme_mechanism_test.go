package web

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// compiledAppCSS reads the compiled, embedded bundle — the single asset served
// to every theme. It is the artifact produced by `npm run build` from
// web/styles/app.css, so these assertions document the real bundle, not the
// source imports.
func compiledAppCSS(t *testing.T) string {
	t.Helper()
	css, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read compiled app.css: %v", err)
	}
	return string(css)
}

// TestEveryThemeShipsRootSelectorInSourceAndBundle proves Phase H contract
// points (a) and (b): every theme that exists on disk defines its own root
// selector .theme-<name> in themes/<name>/theme.css AND that selector reaches
// the compiled bundle served to the browser. Class-driven selection can only
// work if the selector is inside the single served asset — a theme whose root
// selector is missing from the bundle is a theme that cannot be selected.
func TestEveryThemeShipsRootSelectorInSourceAndBundle(t *testing.T) {
	themes := availableThemes(t)
	if len(themes) == 0 {
		t.Fatal("no themes found on disk; the bundle contract cannot be verified")
	}
	compiled := compactCSS(t, compiledAppCSS(t))
	for _, name := range themes {
		source := compactCSS(t, themeCSS(t, name))
		// The directory is themes/<name>/ and the root selector is .<name>
		// (e.g. themes/theme-material/ → .theme-material).
		root := "." + name + "{"
		if !strings.Contains(source, root) {
			t.Errorf("%s theme.css must open its root block with %q", name, root)
		}
		if !strings.Contains(compiled, root) {
			t.Errorf("compiled bundle is missing the %q root selector (theme %s is not selectable at runtime)", root, name)
		}
	}
}

// TestAppCSSImportsEveryThemeExplicitly proves contract point (b) at the entry
// level: the bundle lists each theme on disk as an explicit @import. CSS does
// not glob, so a theme is in the bundle if and only if its import line exists
// here — adding theme-basecoat is literally one import line plus its class.
func TestAppCSSImportsEveryThemeExplicitly(t *testing.T) {
	entry, err := sourceStyles.ReadFile("styles/app.css")
	if err != nil {
		t.Fatalf("read styles/app.css: %v", err)
	}
	imports := map[string]bool{}
	for _, m := range regexp.MustCompile(`@import\s+"\.\./\.\./themes/([a-z0-9-]+)/theme\.css"`).FindAllStringSubmatch(string(entry), -1) {
		imports[m[1]] = true
	}
	if len(imports) == 0 {
		t.Fatal("app.css must import at least one theme explicitly")
	}
	for _, name := range availableThemes(t) {
		if !imports[name] {
			t.Errorf("theme %s exists on disk but app.css does not import it (add the explicit @import)", name)
		}
	}
	for name := range imports {
		if _, err := os.Stat(filepath.Join(repositoryRoot(t), "themes", name, "theme.css")); err != nil {
			t.Errorf("app.css imports theme %s but themes/%s/theme.css does not exist", name, name)
		}
	}
}

// TestThemeSelectionIsClassDrivenWithoutJS proves contract point (c): the
// served bundle is a single inlined asset (zero residual @imports) and the
// shipped app.js carries no theme-selection logic. Switching themes is purely a
// class swap on the document root driven by the server ({{.ThemeClass}}) — no
// runtime JavaScript picks or toggles a theme, no rebuild is needed.
func TestThemeSelectionIsClassDrivenWithoutJS(t *testing.T) {
	compiled := compiledAppCSS(t)
	if strings.Contains(compiled, "@import") {
		t.Error("compiled bundle must inline every theme (no residual @import served to the browser)")
	}

	appJS, err := Assets.ReadFile("static/app.js")
	if err != nil {
		t.Fatalf("read static/app.js: %v", err)
	}
	js := string(appJS)
	for _, forbidden := range []string{".theme-", "data-theme"} {
		if strings.Contains(js, forbidden) {
			t.Errorf("app.js must not contain %q (theme selection is class-driven server-side, no JS)", forbidden)
		}
	}
}

// compactCSS removes every whitespace run so minified-selector assertions
// (.theme-<name>{) hold for both the pretty source and the lightningcss output.
func compactCSS(t *testing.T, css string) string {
	t.Helper()
	return regexp.MustCompile(`\s+`).ReplaceAllString(css, "")
}
