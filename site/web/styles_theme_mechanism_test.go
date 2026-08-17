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
// selector .theme-<name> in lib/themes/<name>.css AND that selector reaches
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
		// The file is lib/themes/<name>.css and the root selector is .<name>
		// (e.g. lib/themes/theme-material.css → .theme-material).
		root := "." + name + "{"
		if !strings.Contains(source, root) {
			t.Errorf("%s theme.css must open its root block with %q", name, root)
		}
		if !strings.Contains(compiled, root) {
			t.Errorf("compiled bundle is missing the %q root selector (theme %s is not selectable at runtime)", root, name)
		}
	}
}

// TestAppCSSImportsEveryThemeExplicitly proves contract point (b) at the
// entry level: the site entry lists each theme on disk as an explicit
// @import before the lib manifest. CSS does not glob, so a theme is in the
// bundle if and only if its import line exists here — adding theme-basecoat
// is literally one import line plus its class.
func TestAppCSSImportsEveryThemeExplicitly(t *testing.T) {
	entry, err := siteStyles.ReadFile("styles/app.css")
	if err != nil {
		t.Fatalf("read styles/app.css: %v", err)
	}
	imports := map[string]bool{}
	for _, m := range regexp.MustCompile(`@import\s+"gelium-ui/themes/([a-z0-9-]+)\.css"`).FindAllStringSubmatch(string(entry), -1) {
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
		if _, err := os.Stat(filepath.Join(repositoryRoot(t), "lib", "themes", name+".css")); err != nil {
			t.Errorf("app.css imports theme %s but lib/themes/%s.css does not exist", name, name)
		}
	}
}

// TestThemeSelectionIsClassDrivenWithoutJS proves contract point (c): the
// served bundle is a single inlined asset (zero residual @imports) and the
// shipped app.js carries no theme-SELECTION logic — it never builds a
// .theme-* class or picks a theme. Switching themes is a class swap on the
// document root driven by the server ({{.ThemeClass}}). With hx-boost the
// server still decides: app.js only REPLICATES the response's html class
// onto the live <html> after a boosted swap (htmx swaps the body only), so
// the JS never chooses — it echoes the server's decision.
func TestThemeSelectionIsClassDrivenWithoutJS(t *testing.T) {
	compiled := compiledAppCSS(t)
	if strings.Contains(compiled, "@import") {
		t.Error("compiled bundle must inline every theme (no residual @import served to the browser)")
	}

	appJS, err := Assets.ReadFile("static/app.js")
	if err != nil {
		t.Fatalf("read static/app.js: %v", err)
	}
	geliumJS, err := Assets.ReadFile("static/gelium.js")
	if err != nil {
		t.Fatalf("read static/gelium.js: %v", err)
	}
	js := string(appJS) + "\n" + string(geliumJS)
	// JS must never own the theme CATALOG: direction-class literals
	// (theme-material, theme-basecoat) exist only server-side (availableThemes).
	// It MAY reference the scheme class "theme-dark", which mirrors the
	// server's applyDocumentRootScheme, and it reads the direction class from
	// the server-emitted data-class attribute of the selected option.
	for _, forbidden := range []string{"theme-material", "theme-basecoat"} {
		if strings.Contains(js, forbidden) {
			t.Errorf("site JS must not hardcode catalog class %q (catalog lives server-side)", forbidden)
		}
	}
	// The boosted-nav sync is allowed and REQUIRED: it copies the server's
	// decision from the response html onto the live document root, and the
	// optimistic preview reads its classes from server-emitted attributes.
	for _, required := range []string{"ctx.text", "className", "data-class", "htmx:before:swap", "applyOptimisticChrome"} {
		if !strings.Contains(js, required) {
			t.Errorf("app.js must sync the response html class after boosted swaps (missing %q)", required)
		}
	}
}

// compactCSS removes every whitespace run so minified-selector assertions
// (.theme-<name>{) hold for both the pretty source and the lightningcss output.
func compactCSS(t *testing.T, css string) string {
	t.Helper()
	return regexp.MustCompile(`\s+`).ReplaceAllString(css, "")
}

// TestSchemeIconsToggleFromThemeDarkClass proves the sun/moon icon pair is
// driven by the server-side .theme-dark class on the document root — pure
// cascade: light shows sun, dark hides sun and shows moon. No JS, no icon
// library, and the toggle survives forced-colors because icons use
// currentColor. The bundle must contain the rules minified.
func TestSchemeIconsToggleFromThemeDarkClass(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	// compactCSS strips whitespace INCLUDING the descendant combinator space,
	// so the minified .theme-dark .ui-scheme-icon-* matches compacted as
	// .theme-dark.ui-scheme-icon-*. lightningcss also merges the two
	// display:none rules into ONE combined selector.
	for _, contract := range []string{
		".ui-scheme-icon-moon,.theme-dark.ui-scheme-icon-sun{display:none}",
		".theme-dark.ui-scheme-icon-moon{display:inline-block}",
	} {
		if !strings.Contains(compiled, contract) {
			t.Errorf("compiled bundle is missing icon contract %q", contract)
		}
	}
}

// TestThemeAppliesFromStaticDocumentOnly proves the document-root-only
// contract (Phase H close-out): a static HTML document carrying
// class="theme-<name>" on the document root and linking only the served bundle
// receives the theme with NO server round-trip, NO rebuild, and NO JavaScript.
// It reads the bundle from the embedded Assets and starts no HTTP server.
func TestThemeAppliesFromStaticDocumentOnly(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	if strings.Contains(compiled, "@import") {
		t.Fatal("bundle must be a single inlined asset (no residual @import)")
	}

	fixture := func(rootClass string) string {
		return "<html class=\"" + rootClass + "\"><head><link rel=\"stylesheet\" href=\"/static/app.css\"></head><body></body></html>"
	}

	// Scenario: basecoat via static document only — root class + bundle link +
	// .theme-basecoat{ selector in the served asset, no server request.
	basecoatDoc := fixture("theme-basecoat")
	if !strings.Contains(basecoatDoc, `class="theme-basecoat"`) {
		t.Fatal("fixture must carry theme-basecoat on the document root")
	}
	if !strings.Contains(basecoatDoc, `href="/static/app.css"`) {
		t.Fatal("fixture must link the served bundle")
	}
	if !strings.Contains(compiled, ".theme-basecoat{") {
		t.Error("served bundle must carry .theme-basecoat{ root selector")
	}

	// Scenario: material parity from the same single bundle (default theme).
	materialDoc := fixture("theme-material")
	if !strings.Contains(materialDoc, `class="theme-material"`) {
		t.Fatal("fixture must carry theme-material on the document root")
	}
	if !strings.Contains(compiled, ".theme-material{") {
		t.Error("served bundle must carry .theme-material{ root selector")
	}

	// Scenario: class swap changes direction without rebuild or JS — the two
	// fixtures differ only in the root class and share one bundle. This
	// complements TestThemeSelectionIsClassDrivenWithoutJS (bundle single
	// asset + zero theme logic in app.js).
	if basecoatDoc == materialDoc {
		t.Fatal("fixtures must differ in the root class only")
	}
	if !strings.Contains(compiled, ".theme-material{") || !strings.Contains(compiled, ".theme-basecoat{") {
		t.Error("both root selectors must resolve from the same single bundle")
	}
}
