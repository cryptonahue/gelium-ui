package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestSkeletonPrimitiveCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "skeleton.css"), " ")

	for _, contract := range []string{
		`.ui-skeleton {`,
		`display: flex;`,
		`flex-direction: column;`,
		`gap: var(--ui-skeleton-gap);`,
		`padding: var(--ui-skeleton-padding);`,
		`.ui-skeleton-blocks {`,
		`gap: var(--ui-skeleton-block-gap);`,
		`.ui-skeleton-block {`,
		`height: var(--ui-skeleton-block-height);`,
		`border-radius: var(--ui-skeleton-block-radius);`,
		`background: var(--ui-skeleton-block-color);`,
		`animation: ui-skeleton-pulse`,
		`.ui-skeleton-block--title {`,
		`.ui-skeleton-block--short {`,
		`.ui-skeleton-block--circle {`,
		`border-radius: var(--ui-radius-full);`,
		`--ui-skeleton-padding: var(--ui-space-6);`,
		`--ui-skeleton-gap: var(--ui-space-2);`,
		`--ui-skeleton-block-gap: var(--ui-space-1);`,
		`--ui-skeleton-block-height: var(--ui-size-label-height);`,
		`--ui-skeleton-block-radius: var(--ui-radius-sm);`,
		`--ui-skeleton-block-color: var(--ui-color-surface-container);`,
		`--ui-skeleton-line-width: 100%;`,
		`--ui-skeleton-pulse-duration: 1.2s;`,
		`@keyframes ui-skeleton-pulse {`,
		`50% { opacity: .5; }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source skeleton CSS is missing contract %q", contract)
		}
	}

	// The pulse duration is deliberately scoped: a 1.2s opacity loop has no
	// core motion token (150ms/500ms), and the repo precedent is a literal
	// spinner duration (button.css). No core/theme motion token may be added.
	if !strings.Contains(css, `--ui-skeleton-pulse-duration: 1.2s;`) {
		t.Error("skeleton.css must scope the pulse duration with --ui-skeleton-pulse-duration")
	}
	if strings.Contains(css, `var(--ui-motion-`) {
		t.Error("skeleton.css must not consume core motion tokens for the pulse loop")
	}
}

func TestSkeletonContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-skeleton`,
		`.ui-skeleton--avatar`,
		`.ui-skeleton-blocks`,
		`.ui-skeleton-block`,
		`.ui-skeleton-block--circle`,
		`ui-skeleton-pulse`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled skeleton CSS is missing %q", contract)
		}
	}
}

func TestSkeletonClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "skeleton.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "skeleton.css"), " ")

	for _, cls := range []string{
		"ui-skeleton",
		"ui-skeleton--avatar",
		"ui-skeleton-blocks",
		"ui-skeleton-block",
		"ui-skeleton-block--line",
		"ui-skeleton-block--title",
		"ui-skeleton-block--short",
		"ui-skeleton-block--circle",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("skeleton.html is missing class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("skeleton.css is missing selector .%s", cls)
		}
	}

	if strings.Contains(tmpl, "ui-skeleton-demo") {
		t.Error("skeleton.html must not ui-prefix demo scaffolding")
	}
	if strings.Contains(css, ".ui-skeleton-demo") {
		t.Error("skeleton.css must not define .ui-skeleton-demo selectors")
	}
}

func TestSkeletonReducedMotionDisablesAnimation(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "skeleton.css"), " ")
	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("skeleton.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-skeleton-block { animation: none; }`) {
		t.Error("skeleton.css reduced-motion must disable the pulse with animation: none")
	}
}
