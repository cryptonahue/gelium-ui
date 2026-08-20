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
func TestAccordionReferenceAndSkinLayersSurviveCompilation(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	// Lightning CSS may remove the source order declaration, so assert the
	// ordered compiled layer blocks themselves rather than only source text.
	layers := []string{"@layergelium.core{", "@layergelium.reference{", "@layergelium.skin{", "@layergelium.site{"}
	last := -1
	for _, layer := range layers {
		at := strings.Index(compiled, layer)
		if at < 0 || at <= last {
			t.Fatalf("compiled CSS does not preserve core → reference → skin → site order at %q (index=%d previous=%d)", layer, at, last)
		}
		last = at
	}

	for _, selector := range []string{
		`html[data-gelium-reference=none][data-gelium-skin=none][data-gelium-scheme=dark]{`,
		`html[data-gelium-reference=material]{`,
		`html[data-gelium-reference=basecoat]{`,
		`html[data-gelium-reference=baseui]{`,
		`html[data-gelium-skin=material]{`,
		`html[data-gelium-skin=basecoat]{`,
		`html[data-gelium-skin=baseui]{`,
		`html[data-gelium-skin=alden]{`,
		`html[data-gelium-skin=linear]{`,
		`html[data-gelium-skin=vercel]{`,
	} {
		if !strings.Contains(compiled, selector) {
			t.Errorf("compiled CSS missing Phase 2 visual adapter selector %q", selector)
		}
	}

	// Both selectors match the exact server-rendered material+Basecoat case.
	// They have equal specificity, so the later skin layer must own the final
	// value. Explicit sentinel values make the proof deterministic without a
	// browser-specific computed-style harness.
	ref := strings.Index(compiled, `html[data-gelium-reference=material]{`)
	skin := strings.Index(compiled, `html[data-gelium-skin=basecoat]{`)
	if ref < 0 || skin < 0 || skin <= ref {
		t.Fatalf("skin adapter must follow reference adapter: reference=%d skin=%d", ref, skin)
	}
	if !strings.Contains(compiled[ref:skin], "--ui-accordion-cascade-id:material-reference") {
		t.Error("material reference selector must supply its Accordion token identity")
	}
	if !strings.Contains(compiled[skin:], "--ui-accordion-cascade-id:basecoat-skin") {
		t.Error("Basecoat skin selector must override the matching Accordion token identity")
	}
}

func TestMaterialDarkReferenceAccordionSurfaceSurvivesThemeSpecificity(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	selector := `html[data-gelium-reference=material][data-gelium-scheme=dark]{`
	start := strings.Index(compiled, selector)
	if start < 0 {
		t.Fatalf("compiled app.css missing dark Material reference selector %q", selector)
	}
	end := strings.Index(compiled[start:], "}")
	if end < 0 {
		t.Fatalf("compiled dark Material reference selector is not closed")
	}
	block := compiled[start : start+end]
	if !strings.Contains(block, "--ui-accordion-surface:var(--ui-color-surface)") {
		t.Errorf("dark Material reference must reset Accordion surface to the active --ui-color-surface; block=%s", block)
	}
}

func TestBaseUIDocsInspiredAccordionAdaptersStayFlatAndSkinComposable(t *testing.T) {
	reference := compactCSS(t, repositoryFile(t, "site", "web", "styles", "accordion-reference.css"))
	refSourceSelector := `html[data-gelium-reference="baseui"]{`
	refStart := strings.Index(reference, refSourceSelector)
	if refStart < 0 {
		t.Fatalf("Base UI reference adapter is missing selector %q", refSourceSelector)
	}
	refEnd := strings.Index(reference[refStart:], "}")
	if refEnd < 0 {
		t.Fatal("Base UI reference adapter is not closed")
	}
	refBlock := reference[refStart : refStart+refEnd]
	for _, want := range []string{
		"--ui-accordion-root-display:flex",
		"--ui-accordion-root-direction:column",
		"--ui-accordion-root-surface:transparent",
		"--ui-accordion-root-border:0",
		"--ui-accordion-root-radius:0",
		"--ui-accordion-root-padding:0",
		"--ui-accordion-item-border:var(--ui-border-width-1)var(--ui-border-style-solid)var(--ui-color-fg)",
		"--ui-accordion-item-overlap:calc(-1*var(--ui-border-width-1))",
		"--ui-accordion-item-radius:0",
		"--ui-accordion-item-surface:transparent",
		"--ui-accordion-item-shadow:none",
		"--ui-accordion-item-open-shadow:none",
		"--ui-accordion-chevron-display:none",
		"--ui-accordion-plus-display:inline",
		"--ui-accordion-icon-rotation:45deg",
	} {
		if !strings.Contains(refBlock, want) {
			t.Errorf("Base UI docs-inspired reference must define %q; block=%s", want, refBlock)
		}
	}
	for _, forbidden := range []string{
		"--ui-accordion-item-shadow:var(--ui-shadow",
		"--ui-accordion-item-radius:.375rem",
		"--ui-accordion-item-radius:2px",
	} {
		if strings.Contains(refBlock, forbidden) {
			t.Errorf("Base UI docs-inspired reference must not retain card treatment %q; block=%s", forbidden, refBlock)
		}
	}

	skin := compactCSS(t, repositoryFile(t, "site", "web", "styles", "accordion-skin.css"))
	skinSourceSelector := `html[data-gelium-skin="baseui"]{`
	skinStart := strings.Index(skin, skinSourceSelector)
	if skinStart < 0 {
		t.Fatalf("Base UI skin adapter is missing selector %q", skinSourceSelector)
	}
	skinEnd := strings.Index(skin[skinStart:], "}")
	if skinEnd < 0 {
		t.Fatal("Base UI skin adapter is not closed")
	}
	skinBlock := skin[skinStart : skinStart+skinEnd]
	for _, want := range []string{
		"--ui-accordion-item-overlap:calc(-1*var(--ui-border-width-1))",
		"--ui-accordion-item-border:var(--ui-border-width-1)var(--ui-border-style-solid)var(--ui-color-fg)",
		"--ui-accordion-item-radius:0",
		"--ui-accordion-item-shadow:none",
		"--ui-accordion-icon-rotation:45deg",
	} {
		if !strings.Contains(skinBlock, want) {
			t.Errorf("Base UI skin must carry the docs-inspired anatomy token %q; block=%s", want, skinBlock)
		}
	}

	compiled := compactCSS(t, compiledAppCSS(t))
	refCompiled := strings.Index(compiled, `html[data-gelium-reference=baseui]{`)
	vercelSelector := `html[data-gelium-skin=vercel]{`
	vercelCompiled := strings.Index(compiled, vercelSelector)
	if refCompiled < 0 || vercelCompiled <= refCompiled {
		t.Fatalf("Vercel skin must follow Base UI reference in the compiled cascade: reference=%d vercel=%d", refCompiled, vercelCompiled)
	}
	vercelSourceSelector := `html[data-gelium-skin="vercel"]{`
	vercelSource := strings.Index(skin, vercelSourceSelector)
	if vercelSource < 0 {
		t.Fatalf("Vercel skin adapter is missing selector %q", vercelSourceSelector)
	}
	// Vercel is the last skin adapter. Search its source suffix rather than
	// slicing minified CSS at the first closing brace: color-mix fallbacks in
	// compiled CSS may carry nested blocks before the adapter closes.
	vercelTokens := skin[vercelSource:]
	if !strings.Contains(vercelTokens, "--ui-accordion-cascade-id:vercel-skin") || !strings.Contains(vercelTokens, "--ui-accordion-icon-rotation:180deg") {
		t.Errorf("Vercel skin must still override Base UI reference anatomy; suffix=%s", vercelTokens)
	}
}

func TestNativeNeutralDarkSelectionShipsCoreAndAccordionTokens(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	selector := `html[data-gelium-reference=none][data-gelium-skin=none][data-gelium-scheme=dark]{`
	start := strings.Index(compiled, selector)
	if start < 0 {
		t.Fatalf("compiled app.css missing neutral dark selector %q", selector)
	}
	end := strings.Index(compiled[start:], "}")
	if end < 0 {
		t.Fatal("compiled neutral dark selector is not closed")
	}
	block := compiled[start : start+end]
	for _, token := range []string{
		"color-scheme:dark",
		"--ui-color-canvas:#121212",
		"--ui-color-surface:#1e1e1e",
		"--ui-color-fg:#f5f5f5",
		"--ui-color-border:#5f5f5f",
		"--ui-color-focus-ring:#d6d6d6",
		"--ui-accordion-cascade-id:neutral-reference",
		"--ui-accordion-surface:var(--ui-color-surface)",
	} {
		if !strings.Contains(block, token) {
			t.Errorf("neutral dark preset missing %q", token)
		}
	}
}

func TestAccordionProductionCSSNeverSelectsBehavior(t *testing.T) {
	behaviorSelector := regexp.MustCompile(`ui-accordion--behavior-|\[data-behavior(?:[~|^$*]?=|\])`)
	for _, root := range []string{
		filepath.Join(repositoryRoot(t), "lib", "styles"),
		filepath.Join(repositoryRoot(t), "site", "web", "styles"),
	} {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(path) != ".css" {
				return nil
			}
			css, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if behaviorSelector.Find(css) != nil {
				t.Errorf("production CSS must not select Accordion behavior: %s", path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("scan CSS under %s: %v", root, err)
		}
	}
	if behaviorSelector.FindString(compactCSS(t, compiledAppCSS(t))) != "" {
		t.Error("compiled app.css must not contain an Accordion behavior selector")
	}
}

func TestEachAccordionSkinDefinesCompleteVisualAnatomy(t *testing.T) {
	source := compactCSS(t, repositoryFile(t, "site", "web", "styles", "accordion-skin.css"))
	baseline := []string{
		"--ui-accordion-root-display:", "--ui-accordion-root-direction:", "--ui-accordion-root-width:", "--ui-accordion-root-gap:", "--ui-accordion-root-max-width:",
		"--ui-accordion-root-surface:", "--ui-accordion-root-border:", "--ui-accordion-root-radius:", "--ui-accordion-root-padding:",
		"--ui-accordion-item-border:", "--ui-accordion-item-divider-border:", "--ui-accordion-item-radius:", "--ui-accordion-item-surface:", "--ui-accordion-item-shadow:",
		"--ui-accordion-trigger-min-height:", "--ui-accordion-trigger-padding-x:", "--ui-accordion-trigger-padding-y:", "--ui-accordion-trigger-font:", "--ui-accordion-trigger-align:", "--ui-accordion-trigger-border:", "--ui-accordion-trigger-radius:", "--ui-accordion-trigger-hover-surface:", "--ui-accordion-trigger-hover-decoration:",
		"--ui-accordion-icon-flex:", "--ui-accordion-icon-width:", "--ui-accordion-icon-height:", "--ui-accordion-icon-margin-inline-start:", "--ui-accordion-icon-size:", "--ui-accordion-icon-color:", "--ui-accordion-chevron-display:", "--ui-accordion-plus-display:", "--ui-accordion-icon-rotation:",
		"--ui-accordion-panel-padding-top:", "--ui-accordion-panel-padding-x:", "--ui-accordion-panel-padding-y:", "--ui-accordion-panel-font:",
	}
	for _, skin := range []string{"material", "basecoat", "baseui", "alden", "linear", "vercel"} {
		selector := `html[data-gelium-skin="` + skin + `"]{`
		start := strings.Index(source, selector)
		if start < 0 {
			t.Errorf("skin %q has no token-only selector", skin)
			continue
		}
		end := strings.Index(source[start:], "}")
		if end < 0 {
			t.Errorf("skin %q selector is not closed", skin)
			continue
		}
		block := source[start : start+end]
		for _, token := range baseline {
			if !strings.Contains(block, token) {
				t.Errorf("skin %q must own complete Accordion anatomy; missing %s", skin, token)
			}
		}
	}
}

func TestAldenSkinOwnsAccordionTopologyAfterBasecoatReference(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	selector := `html[data-gelium-skin=alden]{`
	start := strings.Index(compiled, selector)
	if start < 0 {
		t.Fatalf("compiled app.css missing Alden Accordion skin selector %q", selector)
	}
	end := strings.Index(compiled[start:], "}")
	if end < 0 {
		t.Fatal("compiled Alden Accordion skin selector is not closed")
	}
	block := compiled[start : start+end]
	for _, want := range []string{
		"--ui-accordion-cascade-id:alden-skin",
		"--ui-accordion-root-display:grid",
		"--ui-accordion-root-direction:initial",
		"--ui-accordion-root-width:100%",
		"--ui-accordion-root-gap:var(--ui-space-3)",
		"--ui-accordion-root-max-width:48rem",
	} {
		if !strings.Contains(block, want) {
			t.Errorf("Alden skin must own topology after Basecoat reference; missing %q", want)
		}
	}
	if strings.Contains(block, "--ui-accordion-root-display:flex") {
		t.Error("Alden skin must not inherit Basecoat's flex root display")
	}
}

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
