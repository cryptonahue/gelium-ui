package web

import (
	"regexp"
	"strings"
	"testing"
)

// sourceTooltipCSS reads web/styles/tooltip.css through the shared embed used by
// styles_contract_test.go, mirroring how newer components (menu, chips, tabs)
// read their own single file instead of mutating the shared sourceAppCSS list.
func sourceTooltipCSS(t *testing.T) string {
	t.Helper()
	css, err := sourceStyles.ReadFile("styles/tooltip.css")
	if err != nil {
		t.Fatalf("read source tooltip CSS: %v", err)
	}
	return string(css)
}

func TestTooltipPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceTooltipCSS(t), " ")

	for _, contract := range []string{
		`.ui-tooltip-host {`,
		`position: relative;`,
		`.ui-tooltip {`,
		`position: absolute;`,
		`visibility: hidden;`,
		`opacity: 0;`,
		`var(--ui-tooltip-container)`,
		`var(--ui-tooltip-fg)`,
		`var(--ui-radius-xs)`,
		`var(--ui-tooltip-padding)`,
		`var(--ui-tooltip-supporting-text)`,
		`.ui-tooltip--rich {`,
		`.ui-tooltip-subhead {`,
		`.ui-tooltip-supporting-text {`,
		`.ui-tooltip-action {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source tooltip CSS is missing contract %q", contract)
		}
	}
}

// TestTooltipRevealIsCSSOnly guards the roadmap's "accessible visible fallback":
// the tooltip appears on :hover and :focus-within of the host with no component
// JavaScript. Interest Invokers is not Baseline (audited 2026-08), so the reveal
// must not depend on interesttarget/interestaction or a script.
func TestTooltipRevealIsCSSOnly(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceTooltipCSS(t), " ")

	for _, sel := range []string{
		`.ui-tooltip-host:hover .ui-tooltip, .ui-tooltip-host:focus-within .ui-tooltip {`,
		`visibility: visible;`,
		`opacity: 1;`,
		`transition: opacity`,
	} {
		if !strings.Contains(css, sel) {
			t.Errorf("source tooltip CSS is missing reveal contract %q", sel)
		}
	}

	if strings.Contains(css, `interestaction`) {
		t.Error("tooltip must not depend on the not-yet-Baseline Interest Invokers API")
	}
}

func TestTooltipDemoClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "tooltip.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceTooltipCSS(t), " ")

	for _, cls := range []string{
		"tooltip-demo-grid",
		"tooltip-demo-group",
		"tooltip-demo-heading",
		"tooltip-demo-note",
		"tooltip-demo-row",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("tooltip.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("tooltip.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-tooltip-host",
		"ui-tooltip",
		"ui-tooltip--rich",
		"ui-tooltip--top",
		"ui-tooltip-subhead",
		"ui-tooltip-supporting-text",
		"ui-tooltip-action",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("tooltip.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("tooltip.css is missing primitive selector .%s", cls)
		}
	}

	if strings.Contains(tmpl, "ui-tooltip-demo") {
		t.Error("tooltip.html must not ui-prefix the demo scaffolding classes")
	}
	if strings.Contains(css, ".ui-tooltip-demo") {
		t.Error("tooltip.css must not define .ui-tooltip-demo selectors")
	}
}

func TestTooltipReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceTooltipCSS(t), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("tooltip.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `transition: none;`) {
		t.Error("tooltip reduced-motion must disable transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("tooltip.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-tooltip { border: 1px solid CanvasText;`) {
		t.Error("tooltip surface must stay discernible in forced colors")
	}
	if !strings.Contains(css, `.ui-tooltip-action:focus-visible { outline-color: Highlight;`) {
		t.Error("tooltip action focus ring must switch to Highlight in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesTooltipContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-tooltip-host`,
		`.ui-tooltip--rich`,
		`.ui-tooltip--top`,
		`.ui-tooltip-subhead`,
		`.tooltip-demo-grid`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled tooltip CSS is missing %q", contract)
		}
	}
}
