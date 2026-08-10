package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestNavigationDrawerPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, contract := range []string{
		`.ui-navigation-drawer {`,
		`width: var(--ui-navigation-drawer-container-width);`,
		`height: 100%;`,
		`background: var(--ui-navigation-drawer-container-color);`,
		`box-shadow: var(--ui-navigation-drawer-container-elevation);`,
		`border-start-end-radius: var(--ui-radius-lg);`,
		`border-end-end-radius: var(--ui-radius-lg);`,
		`.ui-navigation-drawer-list {`,
		`list-style: none;`,
		`.ui-navigation-drawer-item {`,
		`.ui-navigation-drawer-destination {`,
		`min-height: var(--ui-navigation-drawer-item-height);`,
		`gap: 12px;`,
		`font: var(--ui-type-label-lg);`,
		`.ui-navigation-drawer-glyph {`,
		`width: var(--ui-navigation-drawer-icon-size);`,
		`.ui-navigation-drawer-label {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-drawer CSS is missing contract %q", contract)
		}
	}
}

func TestNavigationDrawerActiveIndicatorIsFullWidthPill(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, contract := range []string{
		`.ui-navigation-drawer-indicator {`,
		`position: absolute;`,
		`inset: 0;`,
		`border-radius: var(--ui-radius-full);`,
		`background: var(--ui-navigation-drawer-indicator-color);`,
		`opacity: 0;`,
		`pointer-events: none;`,
		`.ui-navigation-drawer-destination--active .ui-navigation-drawer-indicator {`,
		`opacity: 1;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-drawer CSS is missing indicator contract %q", contract)
		}
	}
}

func TestNavigationDrawerStatesCoverHoverFocusActiveSelected(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, contract := range []string{
		`.ui-navigation-drawer-destination::before {`,
		`position: absolute;`,
		`inset: 0;`,
		`border-radius: var(--ui-radius-full);`,
		`background: var(--ui-navigation-drawer-state-layer-color);`,
		`opacity: 0;`,
		`pointer-events: none;`,
		`.ui-navigation-drawer-destination:hover::before {`,
		`opacity: var(--ui-navigation-drawer-hover-opacity);`,
		`.ui-navigation-drawer-destination:focus-visible::before {`,
		`opacity: var(--ui-navigation-drawer-focus-opacity);`,
		`.ui-navigation-drawer-destination:active::before {`,
		`opacity: var(--ui-navigation-drawer-pressed-opacity);`,
		`.ui-navigation-drawer-destination:focus-visible {`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-navigation-drawer-destination--active::before {`,
		`background: var(--ui-navigation-drawer-active-state-layer-color);`,
		`.ui-navigation-drawer-destination:hover .ui-navigation-drawer-label,`,
		`.ui-navigation-drawer-destination--active {`,
		`font-weight: var(--ui-navigation-drawer-active-label-weight);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-drawer CSS is missing state contract %q", contract)
		}
	}
}

func TestNavigationDrawerModalVariantOverNativeDialog(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, contract := range []string{
		`.ui-navigation-drawer--modal {`,
		`--ui-navigation-drawer-container-color: var(--ui-dialog-container);`,
		`--ui-navigation-drawer-container-elevation: var(--ui-shadow-1);`,
		`position: fixed;`,
		`inset-block: 0;`,
		`inset-inline-start: 0;`,
		`margin: 0;`,
		`border: 0;`,
		`.ui-navigation-drawer--modal::backdrop {`,
		`background: var(--ui-dialog-scrim);`,
		`opacity: 0;`,
		`.ui-navigation-drawer--modal[open] {`,
		`translate: 0;`,
		`.ui-navigation-drawer--modal[open]::backdrop {`,
		`opacity: 1;`,
		`@starting-style {`,
		`transition-behavior: allow-discrete;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-drawer CSS is missing modal contract %q", contract)
		}
	}
}

func TestNavigationDrawerStandardVariantEmbeddedInLayout(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, contract := range []string{
		`.ui-navigation-drawer--standard {`,
		`display: flex;`,
		`flex-direction: column;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-drawer CSS is missing standard variant contract %q", contract)
		}
	}
}

func TestNavigationDrawerBadgeReuseComposesUiBadge(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")
	tmpl := repositoryFile(t, "web", "templates", "navigation-drawer.html")

	for _, contract := range []string{
		`.ui-navigation-drawer .ui-badge {`,
		`position: absolute;`,
		`inset-block-start: -3px;`,
		`inset-inline-start: calc(100% - 3px);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("navigation-drawer CSS is missing badge composition contract %q", contract)
		}
	}
	// The template must compose the existing .ui-badge primitive rather than a
	// new drawer badge class.
	for _, contract := range []string{
		`class="ui-badge"`,
		`class="ui-badge ui-badge-large"`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("navigation-drawer template is missing reused .ui-badge %q", contract)
		}
	}
}

// TestNavigationDrawerDemoGridClassVocabularyIsClosed is the TDD regression for
// the closed vocabulary rule: demo/preview scaffolding classes in
// navigation-drawer.html MUST match the CSS selectors exactly, and the ui-
// prefix is reserved for the component primitives (.ui-navigation-drawer*),
// never for the demo scaffolding.
func TestNavigationDrawerDemoGridClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "navigation-drawer.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	for _, cls := range []string{
		"navigation-drawer-demo-grid",
		"navigation-drawer-demo-group",
		"navigation-drawer-demo-heading",
		"navigation-drawer-demo-frame",
		"navigation-drawer-demo-content",
		"navigation-drawer-demo-content-title",
		"navigation-drawer-demo-content-body",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-drawer.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-drawer.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-navigation-drawer",
		"ui-navigation-drawer--standard",
		"ui-navigation-drawer--modal",
		"ui-navigation-drawer-list",
		"ui-navigation-drawer-item",
		"ui-navigation-drawer-destination",
		"ui-navigation-drawer-destination--active",
		"ui-navigation-drawer-indicator",
		"ui-navigation-drawer-glyph",
		"ui-navigation-drawer-label",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-drawer.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-drawer.css is missing primitive selector .%s", cls)
		}
	}

	// The glyph-svg class is carried by the trusted SVG constants in
	// navigation_drawer.go (rendered inside the aria-hidden glyph slot), so the
	// CSS selector must exist here and the rendered page is asserted in the Go
	// route tests; it cannot appear as template literal text.
	if !strings.Contains(css, ".ui-navigation-drawer-glyph-svg") {
		t.Error("navigation-drawer.css is missing primitive selector .ui-navigation-drawer-glyph-svg")
	}
}

// TestNavigationDrawerNoFirstOrLastChildGeometrySelectors is the TDD regression
// for the segmented-button lesson: item geometry must never depend on
// :first-child / :last-child of a container whose first/last element may not be
// the item itself. The drawer destinations are full pills, so no leading/trailing
// radius adjustment exists at all.
func TestNavigationDrawerNoFirstOrLastChildGeometrySelectors(t *testing.T) {
	css := sourceComponentCSS(t, "navigation-drawer.css")

	if strings.Contains(css, ":first-child") {
		t.Error("navigation-drawer.css must not use :first-child for item geometry")
	}
	if strings.Contains(css, ":last-child") {
		t.Error("navigation-drawer.css must not use :last-child for item geometry")
	}
}

func TestNavigationDrawerReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-drawer.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("navigation-drawer.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-navigation-drawer-destination::before,`) {
		t.Error("navigation-drawer reduced-motion must disable the state layer transition")
	}
	if !strings.Contains(css, `.ui-navigation-drawer--modal,`) {
		t.Error("navigation-drawer reduced-motion must disable the modal transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("navigation-drawer.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-navigation-drawer { box-shadow: none; border: 1px solid CanvasText;`) {
		t.Error("the drawer must keep a visible boundary in forced colors")
	}
	if !strings.Contains(css, `.ui-navigation-drawer-indicator { border: 1px solid Highlight;`) {
		t.Error("the active indicator must stay visible in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesNavigationDrawerContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-navigation-drawer{`,
		`.ui-navigation-drawer--standard`,
		`.ui-navigation-drawer--modal`,
		`.ui-navigation-drawer-list`,
		`.ui-navigation-drawer-item`,
		`.ui-navigation-drawer-destination`,
		`.ui-navigation-drawer-destination--active`,
		`.ui-navigation-drawer-indicator`,
		`.ui-navigation-drawer-glyph`,
		`.ui-navigation-drawer-glyph-svg`,
		`.ui-navigation-drawer-label`,
		`.ui-navigation-drawer .ui-badge`,
		`.navigation-drawer-demo-grid`,
		`.navigation-drawer-demo-group`,
		`.navigation-drawer-demo-heading`,
		`.navigation-drawer-demo-frame`,
		`.navigation-drawer-demo-content`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled navigation-drawer CSS is missing %q", contract)
		}
	}
}
