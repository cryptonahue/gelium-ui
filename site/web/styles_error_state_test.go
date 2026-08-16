package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestErrorStatePrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "error-state.css"), " ")

	for _, contract := range []string{
		`.ui-error-state {`,
		`display: flex;`,
		`flex-direction: column;`,
		`align-items: center;`,
		`place-content: center;`,
		`gap: var(--ui-error-state-gap);`,
		`padding: var(--ui-error-state-padding);`,
		`text-align: center;`,
		`.ui-error-state-code {`,
		`margin: 0;`,
		`font: var(--ui-type-display-lg);`,
		`color: var(--ui-error-state-code-color);`,
		`.ui-error-state-title {`,
		`font: var(--ui-type-title-lg);`,
		`color: var(--ui-error-state-title-color);`,
		`.ui-error-state-body {`,
		`font: var(--ui-type-body-sm);`,
		`color: var(--ui-error-state-body-color);`,
		`.ui-error-state .ui-button {`,
		`margin-top: var(--ui-space-2);`,
		`--ui-error-state-padding: var(--ui-space-8);`,
		`--ui-error-state-gap: var(--ui-space-2);`,
		`--ui-error-state-code-color: var(--ui-color-danger);`,
		`--ui-error-state-title-color: var(--ui-color-fg);`,
		`--ui-error-state-body-color: var(--ui-color-fg-muted);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source error-state CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("error-state.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("error-state.css must not declare transitions or animations")
	}
}

func TestErrorStateContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-error-state`,
		`.ui-error-state-code`,
		`.ui-error-state-title`,
		`.ui-error-state-body`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled error-state CSS is missing %q", contract)
		}
	}
}

func TestErrorStateClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "error-state.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "error-state.css"), " ")

	for _, cls := range []string{
		"ui-error-state",
		"ui-error-state-code",
		"ui-error-state-title",
		"ui-error-state-body",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("error-state.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("error-state.css is missing selector .%s", cls)
		}
	}

	// The status code is the canonical attribute: the template must interpolate
	// it server-side, and the code element stays decorative (aria-hidden).
	if !strings.Contains(tmpl, `{{.StatusCode}}`) {
		t.Error("error-state.html must interpolate the status code")
	}

	if strings.Contains(tmpl, "ui-error-state-demo") {
		t.Error("error-state.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-error-state-demo") {
		t.Error("error-state.css must not define .ui-error-state-demo selectors")
	}
}

func TestErrorStateUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "error-state.css"), " ")

	// The whole primitive resolves through core tokens (danger for the code
	// accent, fg/fg-muted for heading/body, the space scale, the typescale)
	// and the scoped aliases over them, never a raw color literal.
	for _, contract := range []string{
		`--ui-error-state-code-color: var(--ui-color-danger);`,
		`--ui-error-state-title-color: var(--ui-color-fg);`,
		`--ui-error-state-body-color: var(--ui-color-fg-muted);`,
		`var(--ui-space-8)`,
		`var(--ui-space-2)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("error-state.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("error-state.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
	if strings.Contains(css, "rgb(") {
		t.Error("error-state.css must not contain a raw rgb() literal")
	}
}
