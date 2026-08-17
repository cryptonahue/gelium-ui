package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestLanguageSwitcherPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "language-switcher.css"), " ")

	for _, contract := range []string{
		`.ui-language-switcher {`,
		`display: flex;`,
		`flex-wrap: wrap;`,
		`align-items: center;`,
		`.ui-language-switcher-label {`,
		`.ui-language-switcher-control {`,
		`.ui-language-switcher-select {`,
		`.ui-language-switcher-select:focus-visible {`,
		`--ui-language-switcher-gap: var(--ui-space-2);`,
		`--ui-language-switcher-type: var(--ui-type-body-sm);`,
		`--ui-language-switcher-color: var(--ui-color-fg-muted);`,
		`--ui-language-switcher-border: var(--ui-color-border);`,
		`--ui-language-switcher-radius: var(--ui-radius-sm);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source language-switcher CSS is missing contract %q", contract)
		}
	}

	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("language-switcher.css must not declare transitions or animations")
	}
}

func TestLanguageSwitcherContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-language-switcher`,
		`.ui-language-switcher-select`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled language-switcher CSS is missing %q", contract)
		}
	}
}

func TestLanguageSwitcherClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "site", "web", "templates", "language-switcher.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "language-switcher.css"), " ")

	for _, cls := range []string{
		"ui-language-switcher",
		"ui-language-switcher-label",
		"ui-language-switcher-control",
		"ui-language-switcher-select",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("language-switcher.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("language-switcher.css is missing selector .%s", cls)
		}
	}

	// Language change is a GET navigation, never a POST mutation: the form is
	// method=get with a visible submit (no auto-submit JS), the select posts
	// name=lang, and the current locale is marked with the selected attribute.
	for _, contract := range []string{
		`<form class="ui-language-switcher" method="get" action="{{.Action}}">`,
		`name="lang"`,
		`<option value="{{.Value}}"{{if eq $.Current .Value}} selected{{end}}>`,
		`{{template "button" .Submit}}`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("language-switcher.html is missing contract %q", contract)
		}
	}
	if strings.Contains(tmpl, "method=\"post\"") {
		t.Error("language-switcher.html must never be a POST form (language change is GET navigation)")
	}
	if strings.Contains(tmpl, "ui-language-switcher-demo") {
		t.Error("language-switcher.html must not ui-prefix demo scaffolding")
	}
}

func TestLanguageSwitcherUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "language-switcher.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-border)`,
		`var(--ui-type-body-sm)`,
		`var(--ui-radius-sm)`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("language-switcher.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("language-switcher.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
