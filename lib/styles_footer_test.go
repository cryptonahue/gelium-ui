package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestFooterPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "footer.css"), " ")

	for _, contract := range []string{
		`.ui-footer {`,
		`display: flex;`,
		`flex-direction: column;`,
		`background: var(--ui-footer-surface);`,
		`color: var(--ui-footer-fg);`,
		`.ui-footer-brand {`,
		`.ui-footer-nav {`,
		`display: grid;`,
		`grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));`,
		`.ui-footer-section {`,
		`.ui-footer-heading {`,
		`.ui-footer-list {`,
		`.ui-footer-list a {`,
		`.ui-footer-legal {`,
		`border-top: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-footer-border);`,
		`--ui-footer-surface: var(--ui-color-surface);`,
		`--ui-footer-fg: var(--ui-color-fg-muted);`,
		`--ui-footer-heading-color: var(--ui-color-fg);`,
		`--ui-footer-border: var(--ui-color-border);`,
		`--ui-footer-heading-type: var(--ui-type-label-lg);`,
		`--ui-footer-type: var(--ui-type-body-sm);`,
		`--ui-footer-gap: var(--ui-space-6);`,
		`--ui-footer-section-gap: var(--ui-space-2);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source footer CSS is missing contract %q", contract)
		}
	}

	// Desktop breakpoint forces every collapsed section open by overriding the
	// UA details:not([open]) display:none — the zero-JS accordion escape hatch.
	if !strings.Contains(css, `@media (min-width: 48rem)`) {
		t.Error("footer.css must declare the desktop responsive breakpoint")
	}
	if !strings.Contains(css, `.ui-footer-details > summary`) {
		t.Error("footer.css must hide the disclosure chrome on desktop")
	}
	if !strings.Contains(css, `.ui-footer-list`) {
		t.Error("footer.css must force footer lists open on desktop")
	}
	if !strings.Contains(css, `.ui-footer-details:not([open]) > .ui-footer-list { display: block; }`) {
		t.Error("footer.css must override closed details for desktop footer lists")
	}
	if !strings.Contains(css, `.ui-footer-details { display: contents; }`) {
		t.Error("footer.css must keep desktop details in the layout")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("footer.css must not declare transitions or animations")
	}
}

func TestFooterContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-footer`,
		`.ui-footer-brand`,
		`.ui-footer-nav`,
		`.ui-footer-section`,
		`.ui-footer-legal`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled footer CSS is missing %q", contract)
		}
	}
}

func TestFooterClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "footer.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "footer.css"), " ")

	for _, cls := range []string{
		"ui-footer",
		"ui-footer-brand",
		"ui-footer-nav",
		"ui-footer-section",
		"ui-footer-heading",
		"ui-footer-list",
		"ui-footer-legal",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("footer.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("footer.css is missing selector .%s", cls)
		}
	}

	// Groups are open initially so desktop always exposes their links. On narrow
	// screens the native summary remains available to collapse each group.
	if !strings.Contains(tmpl, "<details class=\"ui-footer-details\" open>") {
		t.Error("footer.html must render footer groups open initially")
	}
	if strings.Contains(tmpl, "<details open") {
		t.Error("footer.html must not render <details open> (collapsed by default)")
	}
	if strings.Contains(tmpl, "ui-footer-demo") {
		t.Error("footer.html must not ui-prefix demo scaffolding")
	}
}

func TestFooterUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "footer.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-surface)`,
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-fg)`,
		`var(--ui-type-body-sm)`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("footer.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("footer.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
