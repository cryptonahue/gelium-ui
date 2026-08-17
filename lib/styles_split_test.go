package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestSplitPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "split.css"), " ")

	for _, contract := range []string{
		`.ui-split {`,
		`display: grid;`,
		`grid-template-columns: repeat(2, minmax(0, 1fr));`,
		`gap: var(--ui-split-gap);`,
		`.ui-split-media {`,
		`.ui-split-media img,`,
		`.ui-split-body {`,
		`.ui-split-eyebrow {`,
		`.ui-split-title {`,
		`.ui-split-copy {`,
		`.ui-split-action {`,
		`--ui-split-gap: var(--ui-space-6);`,
		`--ui-split-media-radius: var(--ui-radius-sm);`,
		`--ui-split-media-border: var(--ui-color-border);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source split CSS is missing contract %q", contract)
		}
	}

	// Narrow screens stack the columns: a single column grid replaces the pair.
	if !strings.Contains(css, `@media (max-width: 47.99rem)`) {
		t.Error("split.css must declare the narrow responsive breakpoint")
	}
	if !strings.Contains(css, `.ui-split-body:only-child { grid-column: 1 / -1; }`) {
		t.Error("split.css must stretch the body across both columns when media is absent")
	}
	// No literal left/right: RTL bidi must come from the grid direction order.
	if strings.Contains(css, `left:`) || strings.Contains(css, `right:`) {
		t.Error("split.css must not use literal left/right (RTL bidi comes from grid direction)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("split.css must not declare transitions or animations")
	}
}

func TestSplitContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-split`,
		`.ui-split-media`,
		`.ui-split-body`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled split CSS is missing %q", contract)
		}
	}
}

func TestSplitClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "split.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "split.css"), " ")

	for _, cls := range []string{
		"ui-split",
		"ui-split-media",
		"ui-split-body",
		"ui-split-eyebrow",
		"ui-split-title",
		"ui-split-copy",
		"ui-split-action",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("split.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("split.css is missing selector .%s", cls)
		}
	}

	// Media, eyebrow and CTA are optional slots; the body column is the only
	// always-rendered child, and the CTA reuses the real button partial.
	for _, contract := range []string{
		`{{if .Media}}<div class="ui-split-media">`,
		`<section class="ui-split">`,
		`{{template "button" .CTA}}`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("split.html is missing contract %q", contract)
		}
	}
	if strings.Contains(tmpl, "ui-split-demo") {
		t.Error("split.html must not ui-prefix demo scaffolding")
	}
}

func TestSplitUsesCoreTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "split.css"), " ")

	for _, contract := range []string{
		`var(--ui-color-fg-muted)`,
		`var(--ui-color-fg)`,
		`var(--ui-type-headline-sm)`,
		`var(--ui-type-body-lg)`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("split.css is missing core token reference %q", contract)
		}
	}

	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	if m := hexLiteral.FindString(css); m != "" {
		t.Errorf("split.css must not contain the color literal %s (use a --ui-color-* token)", m)
	}
}
