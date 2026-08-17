package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestInlineAlertPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "inline-alert.css"), " ")

	for _, contract := range []string{
		`.ui-inline-alert {`,
		`display: flex;`,
		`align-items: flex-start;`,
		`gap: var(--ui-inline-alert-gap);`,
		`padding: var(--ui-inline-alert-padding);`,
		`border-radius: var(--ui-inline-alert-radius);`,
		`background: var(--ui-inline-alert-bg);`,
		`color: var(--ui-inline-alert-fg);`,
		`.ui-inline-alert--error {`,
		`.ui-inline-alert--success {`,
		`.ui-inline-alert--info {`,
		`.ui-inline-alert--warning {`,
		`.ui-inline-alert-icon {`,
		`width: var(--ui-inline-alert-icon-size);`,
		`height: var(--ui-inline-alert-icon-size);`,
		`.ui-inline-alert-title {`,
		`font: var(--ui-type-label-lg);`,
		`.ui-inline-alert-body {`,
		`font: var(--ui-type-body-sm);`,
		`--ui-inline-alert-padding: var(--ui-space-3) var(--ui-space-4);`,
		`--ui-inline-alert-gap: var(--ui-space-2);`,
		`--ui-inline-alert-radius: var(--ui-radius-sm);`,
		`--ui-inline-alert-icon-size: var(--ui-size-icon-sm);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source inline-alert CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("inline-alert.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("inline-alert.css must not declare transitions or animations")
	}
}

func TestInlineAlertContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-inline-alert`,
		`.ui-inline-alert--error`,
		`.ui-inline-alert--success`,
		`.ui-inline-alert--info`,
		`.ui-inline-alert--warning`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled inline-alert CSS is missing %q", contract)
		}
	}
}

func TestInlineAlertClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "inline-alert.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "inline-alert.css"), " ")

	for _, cls := range []string{
		"ui-inline-alert",
		"ui-inline-alert-icon",
		"ui-inline-alert-title",
		"ui-inline-alert-body",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("inline-alert.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("inline-alert.css is missing selector .%s", cls)
		}
	}

	// The tone modifier is interpolated (ui-inline-alert--{{.Tone}}), so the
	// template carries the prefix and the CSS owns the closed tone vocabulary.
	if !strings.Contains(tmpl, "ui-inline-alert--") {
		t.Error("inline-alert.html must interpolate the tone modifier ui-inline-alert--")
	}
	for _, tone := range []string{"error", "success", "info", "warning"} {
		if !strings.Contains(css, ".ui-inline-alert--"+tone) {
			t.Errorf("inline-alert.css is missing tone selector .ui-inline-alert--%s", tone)
		}
	}

	if strings.Contains(tmpl, "ui-inline-alert-demo") {
		t.Error("inline-alert.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-inline-alert-demo") {
		t.Error("inline-alert.css must not define .ui-inline-alert-demo selectors")
	}
}

func TestInlineAlertTonesUseCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "inline-alert.css"), " ")

	// Every tone resolves bg/fg through core tokens (or a color-mix over them),
	// never a raw color literal.
	for _, contract := range []string{
		`.ui-inline-alert--error { --ui-inline-alert-bg: var(--ui-color-danger-container); --ui-inline-alert-fg: var(--ui-color-danger); }`,
		`.ui-inline-alert--warning { --ui-inline-alert-bg: var(--ui-color-warning-container); --ui-inline-alert-fg: var(--ui-color-warning-fg); }`,
		`.ui-inline-alert--success { --ui-inline-alert-bg: color-mix(in srgb, var(--ui-color-success) 12%, var(--ui-color-surface)); --ui-inline-alert-fg: var(--ui-color-success); }`,
		`.ui-inline-alert--info { --ui-inline-alert-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface)); --ui-inline-alert-fg: var(--ui-color-info); }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("inline-alert.css tone is missing contract %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("inline-alert.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
