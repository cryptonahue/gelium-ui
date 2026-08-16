package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestSectionHeadingUtilityCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "section-heading.css"), " ")

	for _, contract := range []string{
		`.ui-section-heading {`,
		`margin: var(--ui-section-heading-margin);`,
		`color: var(--ui-section-heading-color);`,
		`font: var(--ui-section-heading-type);`,
		`.ui-section-heading--centered {`,
		`.ui-section-heading-eyebrow {`,
		`--ui-section-heading-type: var(--ui-type-headline-sm);`,
		`--ui-section-heading-color: var(--ui-color-fg);`,
		`--ui-section-heading-margin: var(--ui-space-6) 0 var(--ui-space-3);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source section-heading CSS is missing contract %q", contract)
		}
	}

	// The utility is a static text rule: no motion and no raw colors.
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("section-heading.css must not declare transitions or animations")
	}
	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("section-heading.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}

func TestSectionHeadingContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-section-heading`,
		`.ui-section-heading--centered`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled section-heading CSS is missing %q", contract)
		}
	}
}

func TestSectionHeadingClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "section-heading.html")

	// The partial must always render h2 (the page owns a single h1, P2) and the
	// closed utility vocabulary, with the eyebrow as the only optional piece.
	if !strings.Contains(tmpl, `<h2 class="ui-section-heading`) {
		t.Error("section-heading.html must render an <h2 class=\"ui-section-heading\"> element")
	}
	if strings.Contains(tmpl, "<h1") {
		t.Error("section-heading.html must never render h1 (P2: single h1 per page)")
	}
	if strings.Contains(tmpl, "ui-section-heading-demo") {
		t.Error("section-heading.html must not ui-prefix demo scaffolding")
	}
}
