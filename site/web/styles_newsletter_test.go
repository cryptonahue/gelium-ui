package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestNewsletterPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "newsletter.css"), " ")

	for _, contract := range []string{
		`.ui-newsletter {`,
		`display: flex;`,
		`flex-direction: column;`,
		`.ui-newsletter-title {`,
		`.ui-newsletter-description {`,
		`.ui-newsletter-success {`,
		`.ui-newsletter-form {`,
		`.ui-newsletter-field {`,
		`.ui-newsletter-row {`,
		`.ui-newsletter-label {`,
		`.ui-newsletter-input {`,
		`.ui-newsletter-input[aria-invalid="true"] {`,
		`--ui-newsletter-surface: var(--ui-color-surface-container);`,
		`--ui-newsletter-title-type: var(--ui-type-headline-sm);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source newsletter CSS is missing contract %q", contract)
		}
	}

	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("newsletter.css must not declare transitions or animations")
	}
}

func TestNewsletterContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-newsletter`,
		`.ui-newsletter-title`,
		`.ui-newsletter-input`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled newsletter CSS is missing %q", contract)
		}
	}
}

func TestNewsletterClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "newsletter.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "newsletter.css"), " ")

	for _, cls := range []string{
		"ui-newsletter",
		"ui-newsletter-title",
		"ui-newsletter-description",
		"ui-newsletter-success",
		"ui-newsletter-form",
		"ui-newsletter-field",
		"ui-newsletter-row",
		"ui-newsletter-label",
		"ui-newsletter-input",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("newsletter.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("newsletter.css is missing selector .%s", cls)
		}
	}

	// The aside is a zero-JS subscription form: native POST form with a real
	// email input (required) and the Button partial as the submit. Success
	// replaces the form with a status paragraph; the error reuses the
	// inline-alert primitive.
	for _, contract := range []string{
		`<aside class="ui-newsletter" aria-labelledby="{{.ID}}-title">`,
		`<form class="ui-newsletter-form" method="post" action="{{.Action}}"`,
		`type="email"`,
		`required`,
		`{{template "button" .Submit}}`,
		`{{if .Success}}<p class="ui-newsletter-success" role="status">`,
		`{{template "inline-alert" .Error}}`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("newsletter.html is missing contract %q", contract)
		}
	}
	if strings.Contains(tmpl, "ui-newsletter-demo") {
		t.Error("newsletter.html must not ui-prefix demo scaffolding")
	}
}

func TestNewsletterUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "newsletter.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-surface-container)`,
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-fg)`,
		`var(--ui-color-error)`,
		`var(--ui-type-headline-sm)`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("newsletter.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("newsletter.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
