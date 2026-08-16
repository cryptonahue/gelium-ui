package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestValidationSummaryPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "validation-summary.css"), " ")

	for _, contract := range []string{
		`.ui-validation-summary {`,
		`display: flex;`,
		`flex-direction: column;`,
		`gap: var(--ui-validation-summary-gap);`,
		`padding: var(--ui-validation-summary-padding);`,
		`border-radius: var(--ui-validation-summary-radius);`,
		`background: var(--ui-validation-summary-bg);`,
		`color: var(--ui-validation-summary-fg);`,
		`.ui-validation-summary-title {`,
		`font: var(--ui-type-title-md);`,
		`margin: 0;`,
		`.ui-validation-summary-list {`,
		`list-style: none;`,
		`.ui-validation-summary-item {`,
		`font: var(--ui-type-body-sm);`,
		`.ui-validation-summary-item a {`,
		`text-decoration: underline;`,
		`.ui-validation-summary-item a:focus-visible`,
		`--ui-validation-summary-padding: var(--ui-space-3) var(--ui-space-4);`,
		`--ui-validation-summary-gap: var(--ui-space-2);`,
		`--ui-validation-summary-radius: var(--ui-radius-sm);`,
		`--ui-validation-summary-bg: var(--ui-color-danger-container);`,
		`--ui-validation-summary-fg: var(--ui-color-danger);`,
		`--ui-validation-summary-title-color: var(--ui-color-danger);`,
		`--ui-validation-summary-item-color: var(--ui-color-danger);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source validation-summary CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("validation-summary.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("validation-summary.css must not declare transitions or animations")
	}
}

func TestValidationSummaryContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-validation-summary`,
		`.ui-validation-summary-title`,
		`.ui-validation-summary-list`,
		`.ui-validation-summary-item`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled validation-summary CSS is missing %q", contract)
		}
	}
}

func TestValidationSummaryClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "validation-summary.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "validation-summary.css"), " ")

	for _, cls := range []string{
		"ui-validation-summary",
		"ui-validation-summary-title",
		"ui-validation-summary-list",
		"ui-validation-summary-item",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("validation-summary.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("validation-summary.css is missing selector .%s", cls)
		}
	}

	// The anchor role: each item links to a real field error id, so the
	// template must emit the href on the item anchor and the CSS must style it.
	if !strings.Contains(tmpl, `<a href="{{.Href}}">`) {
		t.Error("validation-summary.html must emit a real anchor per item")
	}

	if strings.Contains(tmpl, "ui-validation-summary-demo") {
		t.Error("validation-summary.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-validation-summary-demo") {
		t.Error("validation-summary.css must not define .ui-validation-summary-demo selectors")
	}
}

func TestValidationSummaryUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "validation-summary.css"), " ")

	// The whole surface resolves through core tokens (danger + container,
	// space scale, radius, typescale, focus) and the scoped aliases over them,
	// never a raw color literal.
	for _, contract := range []string{
		`var(--ui-color-danger-container)`,
		`var(--ui-color-danger)`,
		`var(--ui-space-1)`,
		`var(--ui-space-2)`,
		`var(--ui-space-3)`,
		`var(--ui-space-4)`,
		`var(--ui-radius-sm)`,
		`var(--ui-type-title-md)`,
		`var(--ui-type-body-sm)`,
		`var(--ui-focus-thickness)`,
		`var(--ui-color-focus-ring)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("validation-summary.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("validation-summary.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
	if strings.Contains(css, "rgb(") {
		t.Error("validation-summary.css must not contain a raw rgb() literal")
	}
}
