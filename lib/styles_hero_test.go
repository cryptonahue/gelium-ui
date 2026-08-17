package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestHeroPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "hero.css"), " ")

	for _, contract := range []string{
		`.ui-hero {`,
		`display: flex;`,
		`flex-direction: column;`,
		`align-items: center;`,
		`text-align: center;`,
		`.ui-hero-media {`,
		`.ui-hero--has-media::after {`,
		`.ui-hero-content {`,
		`.ui-hero-eyebrow {`,
		`.ui-hero-title {`,
		`.ui-hero-subtitle {`,
		`.ui-hero-actions {`,
		`--ui-hero-title-type: var(--ui-type-display-lg);`,
		`--ui-hero-surface: var(--ui-color-surface);`,
		`--ui-hero-scrim: var(--ui-color-scrim);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source hero CSS is missing contract %q", contract)
		}
	}

	// The overlay keeps the copy legible over the optional background media and
	// is a no-op without media (only .ui-hero--has-media carries it).
	if !strings.Contains(css, `.ui-hero--has-media::after`) {
		t.Error("hero.css must declare the media scrim overlay")
	}
	if !strings.Contains(css, `@media (min-width: 48rem)`) {
		t.Error("hero.css must declare the wide responsive breakpoint")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("hero.css must not declare transitions or animations")
	}
}

func TestHeroContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-hero`,
		`.ui-hero-media`,
		`.ui-hero-content`,
		`.ui-hero-title`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled hero CSS is missing %q", contract)
		}
	}
}

func TestHeroClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "hero.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "hero.css"), " ")

	for _, cls := range []string{
		"ui-hero",
		"ui-hero-media",
		"ui-hero-content",
		"ui-hero-eyebrow",
		"ui-hero-title",
		"ui-hero-subtitle",
		"ui-hero-actions",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("hero.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("hero.css is missing selector .%s", cls)
		}
	}

	// The hero owns the page's single h1 (display), subtitle and CTA links
	// (Button partial); media is an optional background layer.
	for _, contract := range []string{
		`<section class="ui-hero`,
		`<h1 class="ui-hero-title">`,
		`{{template "button" .}}`,
		`{{if .Media}}<div class="ui-hero-media">`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("hero.html is missing contract %q", contract)
		}
	}
	if strings.Contains(tmpl, "<h2") {
		t.Error("hero.html must not render a section h2 (the hero owns the page h1)")
	}
	if strings.Contains(tmpl, "ui-hero-demo") {
		t.Error("hero.html must not ui-prefix demo scaffolding")
	}
}

func TestHeroUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "hero.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-fg)`,
		`var(--ui-color-scrim)`,
		`var(--ui-type-display-lg)`,
		`var(--ui-type-body-lg)`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("hero.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("hero.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
