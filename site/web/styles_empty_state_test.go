package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestEmptyStatePrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "empty-state.css"), " ")

	for _, contract := range []string{
		`.ui-empty-state {`,
		`display: flex;`,
		`flex-direction: column;`,
		`place-content: center;`,
		`gap: var(--ui-empty-state-gap);`,
		`padding: var(--ui-empty-state-padding);`,
		`text-align: center;`,
		`.ui-empty-state--compact {`,
		`align-items: flex-start;`,
		`text-align: start;`,
		`.ui-empty-state-icon {`,
		`width: var(--ui-empty-state-icon-size);`,
		`height: var(--ui-empty-state-icon-size);`,
		`.ui-empty-state-title {`,
		`font: var(--ui-type-title-md);`,
		`color: var(--ui-empty-state-title-color);`,
		`.ui-empty-state-body {`,
		`font: var(--ui-type-body-sm);`,
		`color: var(--ui-empty-state-body-color);`,
		`--ui-empty-state-padding: var(--ui-space-6);`,
		`--ui-empty-state-gap: var(--ui-space-2);`,
		`--ui-empty-state-icon-size: var(--ui-size-icon);`,
		`--ui-empty-state-title-color: var(--ui-color-fg);`,
		`--ui-empty-state-body-color: var(--ui-color-fg-muted);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source empty-state CSS is missing contract %q", contract)
		}
	}

	// The primitive is deliberately static: no motion means no reduced-motion
	// block and no transition/animation to disable.
	if strings.Contains(css, `prefers-reduced-motion`) {
		t.Error("empty-state.css must not declare a reduced-motion block (no animation to disable)")
	}
	if strings.Contains(css, `transition:` ) || strings.Contains(css, `animation:`) {
		t.Error("empty-state.css must not declare transitions or animations")
	}
}

func TestEmptyStateContractCSSWired(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-empty-state`,
		`.ui-empty-state--compact`,
		`.ui-empty-state-icon`,
		`.ui-empty-state-title`,
		`.ui-empty-state-body`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled empty-state CSS is missing %q", contract)
		}
	}
}

func TestEmptyStateClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "empty-state.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "empty-state.css"), " ")

	for _, cls := range []string{
		"ui-empty-state",
		"ui-empty-state--compact",
		"ui-empty-state-icon",
		"ui-empty-state-title",
		"ui-empty-state-body",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("empty-state.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("empty-state.css is missing selector .%s", cls)
		}
	}

	if strings.Contains(tmpl, "ui-empty-state-demo") {
		t.Error("empty-state.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-empty-state-demo") {
		t.Error("empty-state.css must not define .ui-empty-state-demo selectors")
	}
}
