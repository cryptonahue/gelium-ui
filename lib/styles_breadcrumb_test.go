package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestBreadcrumbPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "breadcrumb.css"), " ")

	for _, contract := range []string{
		`.ui-breadcrumb {`,
		`display: flex;`,
		`flex-wrap: wrap;`,
		`list-style: none;`,
		`.ui-breadcrumb-item {`,
		`font: var(--ui-breadcrumb-type);`,
		`.ui-breadcrumb-item + .ui-breadcrumb-item::before {`,
		`content: var(--ui-breadcrumb-separator);`,
		`.ui-breadcrumb-item a {`,
		`.ui-breadcrumb-item span {`,
		`--ui-breadcrumb-gap: var(--ui-space-1);`,
		`--ui-breadcrumb-type: var(--ui-type-label-sm);`,
		`--ui-breadcrumb-color: var(--ui-color-fg-muted);`,
		`--ui-breadcrumb-current-color: var(--ui-color-fg);`,
		`--ui-breadcrumb-separator: "›";`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source breadcrumb CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("breadcrumb.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("breadcrumb.css must not declare transitions or animations")
	}
}

func TestBreadcrumbContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-breadcrumb`,
		`.ui-breadcrumb-item`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled breadcrumb CSS is missing %q", contract)
		}
	}
}

func TestBreadcrumbClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "breadcrumb.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "breadcrumb.css"), " ")

	for _, cls := range []string{
		"ui-breadcrumb",
		"ui-breadcrumb-item",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("breadcrumb.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("breadcrumb.css is missing selector .%s", cls)
		}
	}

	// The P1 markup contract (seo-patterns.md:50-64): nav → ol → li, the
	// current crumb is a span with aria-current="page" and never an <a>.
	for _, contract := range []string{
		`<nav aria-label="Breadcrumb">`,
		`<ol class="ui-breadcrumb">`,
		`<li class="ui-breadcrumb-item">`,
		`<span aria-current="page">`,
		`<a href="{{.Href}}">`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("breadcrumb.html is missing contract %q", contract)
		}
	}
	if strings.Contains(tmpl, `aria-current="page"`+`<a`) {
		t.Error("breadcrumb.html must never render the current crumb as a link")
	}
	if strings.Contains(tmpl, "ui-breadcrumb-demo") {
		t.Error("breadcrumb.html must not ui-prefix demo scaffolding")
	}
}

func TestBreadcrumbUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "breadcrumb.css"), " ")

	// Every color/type/space reference resolves through core (or theme) tokens,
	// never a raw literal.
	for _, contract := range []string{
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-fg)`,
		`var(--ui-type-label-sm)`,
		`var(--ui-space-1)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("breadcrumb.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("breadcrumb.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
