package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestBannerPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "banner.css"), " ")

	for _, contract := range []string{
		`.ui-banner {`,
		`display: flex;`,
		`align-items: center;`,
		`width: 100%;`,
		`gap: var(--ui-banner-gap);`,
		`padding: var(--ui-banner-padding);`,
		`border-radius: var(--ui-banner-radius);`,
		`background: var(--ui-banner-bg);`,
		`color: var(--ui-banner-fg);`,
		`.ui-banner--error {`,
		`.ui-banner--success {`,
		`.ui-banner--info {`,
		`.ui-banner--warning {`,
		`.ui-banner-icon {`,
		`width: var(--ui-banner-icon-size);`,
		`height: var(--ui-banner-icon-size);`,
		`.ui-banner-content {`,
		`flex: 1;`,
		`.ui-banner-title {`,
		`font: var(--ui-type-label-lg);`,
		`.ui-banner-body {`,
		`font: var(--ui-type-body-sm);`,
		`.ui-banner-dismiss {`,
		`margin: 0 0 0 auto;`,
		`--ui-banner-padding: var(--ui-space-3) var(--ui-space-4);`,
		`--ui-banner-gap: var(--ui-space-3);`,
		`--ui-banner-radius: var(--ui-radius-sm);`,
		`--ui-banner-icon-size: var(--ui-size-icon-sm);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source banner CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("banner.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("banner.css must not declare transitions or animations")
	}
}

func TestBannerContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-banner`,
		`.ui-banner--error`,
		`.ui-banner--success`,
		`.ui-banner--info`,
		`.ui-banner--warning`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled banner CSS is missing %q", contract)
		}
	}
}

func TestBannerClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "banner.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "banner.css"), " ")

	for _, cls := range []string{
		"ui-banner",
		"ui-banner-icon",
		"ui-banner-content",
		"ui-banner-title",
		"ui-banner-body",
		"ui-banner-dismiss",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("banner.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("banner.css is missing selector .%s", cls)
		}
	}

	// The tone modifier is interpolated (ui-banner--{{.Tone}}), so the
	// template carries the prefix and the CSS owns the closed tone vocabulary.
	if !strings.Contains(tmpl, "ui-banner--") {
		t.Error("banner.html must interpolate the tone modifier ui-banner--")
	}
	for _, tone := range []string{"error", "success", "info", "warning"} {
		if !strings.Contains(css, ".ui-banner--"+tone) {
			t.Errorf("banner.css is missing tone selector .ui-banner--%s", tone)
		}
	}

	if strings.Contains(tmpl, "ui-banner-demo") {
		t.Error("banner.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-banner-demo") {
		t.Error("banner.css must not define .ui-banner-demo selectors")
	}
}

func TestBannerTonesUseCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "banner.css"), " ")

	// Every tone resolves bg/fg through core tokens (or a color-mix over them),
	// never a raw color literal.
	for _, contract := range []string{
		`.ui-banner--error { --ui-banner-bg: var(--ui-color-danger-container); --ui-banner-fg: var(--ui-color-danger); }`,
		`.ui-banner--warning { --ui-banner-bg: var(--ui-color-warning-container); --ui-banner-fg: var(--ui-color-warning-fg); }`,
		`.ui-banner--success { --ui-banner-bg: color-mix(in srgb, var(--ui-color-success) 12%, var(--ui-color-surface)); --ui-banner-fg: var(--ui-color-success); }`,
		`.ui-banner--info { --ui-banner-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface)); --ui-banner-fg: var(--ui-color-info); }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("banner.css tone is missing contract %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("banner.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
