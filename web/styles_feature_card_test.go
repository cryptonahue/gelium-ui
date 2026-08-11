package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestFeatureCardCompositionCSSIsMinimal(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "feature-card.css"), " ")

	for _, contract := range []string{
		`.ui-feature-card {`,
		`padding: 0;`,
		`overflow: hidden;`,
		`.ui-feature-card-media {`,
		`aspect-ratio: 16 / 9;`,
		`background: var(--ui-color-surface-container);`,
		`.ui-feature-card-media img,`,
		`.ui-feature-card-media video {`,
		`object-fit: cover;`,
		`.ui-feature-card-body {`,
		`padding: var(--ui-space-4);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source feature-card CSS is missing contract %q", contract)
		}
	}

	// The composition must not re-declare any surface, shadow, focus or state
	// signal: those come from .ui-card itself. Aspect-ratio is literal
	// (structural geometry, same rule as Video) and never tokenized.
	if strings.Contains(css, "ui-card-elevated") || strings.Contains(css, "ui-card-outlined") || strings.Contains(css, "ui-card-filled") {
		t.Error("feature-card.css must not re-declare card surfaces (compose .ui-card instead)")
	}
	if strings.Contains(css, "box-shadow") || strings.Contains(css, "focus-visible") {
		t.Error("feature-card.css must not re-declare card shadow or focus (Card owns those)")
	}
	if strings.Contains(css, `aspect-ratio: var(--ui-`) {
		t.Error("feature-card.css must keep aspect-ratio literal (structural geometry, not a token)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("feature-card.css must not declare transitions or animations")
	}
}

func TestFeatureCardContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-feature-card`,
		`.ui-feature-card-media`,
		`.ui-feature-card-body`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled feature-card CSS is missing %q", contract)
		}
	}
}

func TestFeatureCardComposesCardAndCTA(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "feature-card.html")

	// The feature card is a composition, not a primitive: it must reuse the
	// real .ui-card surface, title/body anatomy and the Button partial for the
	// CTA, adding only the feature media slot and layout wrapper.
	for _, contract := range []string{
		`<article class="ui-card ui-card-elevated ui-feature-card">`,
		`<div class="ui-feature-card-media">`,
		`<h3 class="ui-card-title">`,
		`<p class="ui-card-body">`,
		`<div class="ui-card-action">`,
		`{{template "button" .CTA}}`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("feature-card.html is missing composition contract %q", contract)
		}
	}
	if strings.Contains(tmpl, "ui-feature-card-demo") {
		t.Error("feature-card.html must not ui-prefix demo scaffolding")
	}
}

func TestFeatureCardUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "feature-card.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-surface-container)`,
		`var(--ui-space-4)`,
		`var(--ui-space-2)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("feature-card.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("feature-card.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
