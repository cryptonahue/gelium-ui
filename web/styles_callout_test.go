package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestCalloutPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "callout.css"), " ")

	for _, contract := range []string{
		`.ui-callout {`,
		`display: flex;`,
		`align-items: flex-start;`,
		`gap: var(--ui-callout-gap);`,
		`padding: var(--ui-callout-padding);`,
		`border-radius: var(--ui-callout-radius);`,
		`border-left: 4px solid var(--ui-callout-accent);`,
		`background: var(--ui-callout-bg);`,
		`color: var(--ui-callout-fg);`,
		`.ui-callout--info {`,
		`.ui-callout--tip {`,
		`.ui-callout-icon {`,
		`width: var(--ui-callout-icon-size);`,
		`height: var(--ui-callout-icon-size);`,
		`flex: none;`,
		`.ui-callout-heading {`,
		`font: var(--ui-type-label-lg);`,
		`color: var(--ui-callout-heading-color);`,
		`.ui-callout-body {`,
		`font: var(--ui-type-body-sm);`,
		`color: var(--ui-callout-body-color);`,
		`--ui-callout-padding: var(--ui-space-3) var(--ui-space-4);`,
		`--ui-callout-gap: var(--ui-space-2);`,
		`--ui-callout-radius: var(--ui-radius-sm);`,
		`--ui-callout-icon-size: var(--ui-size-icon-sm);`,
		`--ui-callout-bg: var(--ui-color-surface-container);`,
		`--ui-callout-fg: var(--ui-color-fg);`,
		`--ui-callout-heading-color: var(--ui-color-fg);`,
		`--ui-callout-body-color: var(--ui-color-fg-muted);`,
		`--ui-callout-accent: var(--ui-color-fg-muted);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source callout CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("callout.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("callout.css must not declare transitions or animations")
	}
}

func TestCalloutContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-callout`,
		`.ui-callout--info`,
		`.ui-callout--tip`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled callout CSS is missing %q", contract)
		}
	}
}

func TestCalloutClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "callout.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "callout.css"), " ")

	for _, cls := range []string{
		"ui-callout",
		"ui-callout-icon",
		"ui-callout-heading",
		"ui-callout-body",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("callout.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("callout.css is missing selector .%s", cls)
		}
	}

	// The variant modifier is interpolated (ui-callout--{{.Variant}}), so the
	// template carries the prefix and the CSS owns the closed variant
	// vocabulary. Default (no suffix) is the neutral note; only info and tip
	// are defined — error/warning/success belong to inline-alert and banner.
	if !strings.Contains(tmpl, "ui-callout--") {
		t.Error("callout.html must interpolate the variant modifier ui-callout--")
	}
	for _, variant := range []string{"info", "tip"} {
		if !strings.Contains(css, ".ui-callout--"+variant) {
			t.Errorf("callout.css is missing variant selector .ui-callout--%s", variant)
		}
	}

	if strings.Contains(tmpl, "ui-callout-demo") {
		t.Error("callout.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-callout-demo") {
		t.Error("callout.css must not define .ui-callout-demo selectors")
	}
}

func TestCalloutTonesUseCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "callout.css"), " ")

	// Every variant resolves its surface/accent through core tokens (or a
	// color-mix over them), never a raw color literal.
	for _, contract := range []string{
		`.ui-callout--info { --ui-callout-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface)); --ui-callout-accent: var(--ui-color-info); }`,
		`.ui-callout--tip { --ui-callout-bg: color-mix(in srgb, var(--ui-color-secondary) 20%, var(--ui-color-surface)); --ui-callout-accent: var(--ui-color-secondary-fg); }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("callout.css variant is missing contract %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("callout.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
