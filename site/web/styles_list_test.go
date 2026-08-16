package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestListPrimitiveCSSMapsAnatomyAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-list {`,
		`display: flex;`,
		`flex-direction: column;`,
		`list-style: none;`,
		`background: var(--ui-list-container-color);`,
		`.ui-list-item {`,
		`min-height: var(--ui-list-item-one-line-height);`,
		`padding: var(--ui-space-3) var(--ui-list-item-trailing-space) var(--ui-space-3) var(--ui-list-item-leading-space);`,
		`.ui-list-item--two-line { min-height: var(--ui-list-item-two-line-height);`,
		`.ui-list-item--three-line { min-height: var(--ui-list-item-three-line-height);`,
		`.ui-list-item-headline {`,
		`.ui-list-item-supporting {`,
		`gap: var(--ui-space-4);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing list contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-list-item::before {`,
		`.ui-list-item:hover::before {`,
		`.ui-list-item:focus-within::before {`,
		`.ui-list-item:active::before {`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing list state selector %q", sel)
		}
	}

	// Interactive semantics: real anchor nav and native checkbox selection.
	for _, sel := range []string{
		`.ui-list-item-link {`,
		`.ui-list-item-link:focus-visible {`,
		`.ui-list-item-label {`,
		`.ui-list-item-label input[type="checkbox"] {`,
	} {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing list interactive selector %q", sel)
		}
	}

	// Disabled state on individual items.
	if !strings.Contains(css, `input:disabled`) {
		t.Error("source CSS is missing the disabled-item rule")
	}

	if forcedIndex := strings.Index(css, "@media (forced-colors: active)"); forcedIndex >= 0 {
		forced := css[forcedIndex:]
		if !strings.Contains(forced, `.ui-list { background: Canvas;`) {
			t.Errorf("list must stay discernible in forced colors; missing container rule")
		}
		if !strings.Contains(forced, `.ui-list-item-link { color: LinkText;`) {
			t.Errorf("list links must stay in LinkText in forced colors")
		}
	} else {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
}

func TestListDemoScaffoldingSelectorsMatchTemplateClasses(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	// Demo scaffolding classes must be applied without the ui- prefix and must
	// match the template exactly (see prior QA bug: grid selectors mismatched
	// the template classes and never applied display:flex).
	for _, sel := range []string{
		`.list-demo-grid {`,
		`.list-demo-group {`,
		`.list-demo-caption {`,
	} {
		if !strings.Contains(css, sel) {
			t.Errorf("list demo scaffold is missing CSS selector %q", sel)
		}
	}
}

func TestListComponentCSSDeclaresScopedUIFamily(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, token := range []string{
		`--ui-list-container-color:`,
		`--ui-list-item-leading-space:`,
		`--ui-list-item-trailing-space:`,
		`--ui-list-item-icon-size:`,
		`--ui-list-item-one-line-height:`,
		`--ui-list-item-two-line-height:`,
		`--ui-list-item-three-line-height:`,
		`--ui-list-item-label-color:`,
		`--ui-list-item-supporting-color:`,
		`--ui-list-item-icon-color:`,
		`--ui-list-item-hover-opacity:`,
		`--ui-list-item-pressed-opacity:`,
		`--ui-list-item-focus-opacity:`,
		`--ui-list-item-disabled-opacity:`,
	} {
		if !strings.Contains(css, token) {
			t.Errorf("list component CSS is missing token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesListContract(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-list`,
		`.ui-list-item`,
		`--ui-list-item-one-line-height`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled list CSS is missing %q", contract)
		}
	}
}
