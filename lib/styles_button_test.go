package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestButtonInteractiveStateSelectorsExcludeAriaDisabled(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")

	selectors := []string{
		`.ui-button:hover:not(:disabled):not([aria-disabled="true"])`,
		`.ui-button:active:not(:disabled):not([aria-disabled="true"])`,
	}
	for _, selector := range selectors {
		if !strings.Contains(compact, selector+" {") {
			t.Errorf("source CSS is missing inactive-safe selector %q", selector)
		}
	}

	disabledSelector := `.ui-button:disabled, .ui-button[aria-disabled="true"]`
	disabledIndex := strings.Index(compact, disabledSelector+" {")
	if disabledIndex < 0 {
		t.Fatalf("source CSS is missing neutral disabled rule %q", disabledSelector)
	}
	for _, selector := range selectors {
		if stateIndex := strings.Index(compact, selector+" {"); stateIndex >= disabledIndex {
			t.Errorf("neutral disabled rule must follow interactive state selector %q", selector)
		}
	}
}

func TestTextButtonUsesPrimaryStateLayersWithoutInactiveInteraction(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-button-text { --ui-button-bg: var(--ui-button-text-bg, transparent); --ui-button-fg: var(--ui-button-text-fg, var(--ui-color-primary));`,
		`--ui-button-text-hover-shadow, inset 0 0 0 999px color-mix(in srgb, var(--ui-color-primary) calc(var(--ui-state-hover-opacity) * 100%), transparent)`,
		`.ui-button:focus-visible { outline: var(--ui-button-focus-outline, var(--ui-focus-thickness) solid var(--ui-color-focus-ring)); outline-offset: var(--ui-button-focus-offset, var(--ui-focus-offset));`,
		`opacity: var(--ui-button-disabled-opacity, var(--ui-state-disabled-opacity)); cursor: not-allowed; box-shadow: none;`,
		`min-height: var(--ui-button-resolved-min-height, max(var(--ui-touch-target), var(--ui-button-min-height, var(--ui-touch-target))))`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("text button state contract is missing %q", contract)
		}
	}
}

func TestReducedMotionDisablesButtonSpinnerAnimation(t *testing.T) {
	css := sourceAppCSS(t)
	reducedMotionCSS := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	spinnerAnimationNone := regexp.MustCompile(`(?s)\.ui-button-spinner\s*\{[^}]*animation:\s*none\s*;?[^}]*\}`)
	if !spinnerAnimationNone.MatchString(reducedMotionCSS) {
		t.Error("reduced-motion CSS must disable the spinner with animation: none")
	}
}
func TestButtonConsumesTokenizedStateTransforms(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(readSourceStyle(t, "button.css"), " ")
	for _, contract := range []string{
		`transform: var(--ui-button-transform, none);`,
		`transform var(--ui-motion-short) var(--ui-easing-standard);`,
		`.ui-button:hover:not(:disabled):not([aria-disabled="true"]) {`,
		`transform: var(--ui-button-hover-transform, var(--ui-button-transform, none));`,
		`.ui-button:active:not(:disabled):not([aria-disabled="true"]) {`,
		`transform: var(--ui-button-active-transform, var(--ui-button-hover-transform, var(--ui-button-transform, none)));`,
		`.ui-button-primary { --ui-button-shadow: var(--ui-button-primary-shadow, none); --ui-button-transform: var(--ui-button-primary-transform, none);`,
		`.ui-button-secondary { --ui-button-shadow: var(--ui-button-secondary-shadow, none); --ui-button-transform: var(--ui-button-secondary-transform, none);`,
		`.ui-button-outline { --ui-button-shadow: var(--ui-button-outline-shadow, none); --ui-button-transform: var(--ui-button-outline-transform, none);`,
		`@media (prefers-reduced-motion: reduce) { .ui-button { transform: none; } }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("button transform contract is missing %q", contract)
		}
	}
}

func TestNeubrutalismDefinesButtonDepthEffects(t *testing.T) {
	css := themeCSS(t, "theme-neubrutalism")
	for _, contract := range []string{
		`--ui-button-primary-shadow: var(--ui-shadow-2);`,
		`--ui-button-primary-hover-shadow: none;`,
		`--ui-button-primary-hover-transform: translate(4px, 4px);`,
		`--ui-button-primary-active-shadow: none;`,
		`--ui-button-primary-active-transform: translate(5px, 5px);`,
		`--ui-button-secondary-shadow: var(--ui-shadow-2);`,
		`--ui-button-outline-shadow: var(--ui-shadow-2);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("Neubrutalism button effect contract is missing %q", contract)
		}
	}
}
