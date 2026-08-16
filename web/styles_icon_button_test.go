package web

import (
	"embed"
	"regexp"
	"strings"
	"testing"
)

//go:embed styles/icon-button.css
var iconButtonStyles embed.FS

func sourceIconButtonCSS(t *testing.T) string {
	t.Helper()
	css, err := iconButtonStyles.ReadFile("styles/icon-button.css")
	if err != nil {
		t.Fatalf("read icon-button source css: %v", err)
	}
	return string(css)
}

func TestIconButtonInteractiveStateSelectorsExcludeAriaDisabled(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceIconButtonCSS(t), " ")

	selectors := []string{
		`.ui-icon-button:hover:not(:disabled):not([aria-disabled="true"])`,
		`.ui-icon-button:active:not(:disabled):not([aria-disabled="true"])`,
	}
	for _, selector := range selectors {
		if !strings.Contains(css, selector+" {") {
			t.Errorf("icon-button source CSS is missing inactive-safe selector %q", selector)
		}
	}

	disabledSelector := `.ui-icon-button:disabled, .ui-icon-button[aria-disabled="true"]`
	disabledIndex := strings.Index(css, disabledSelector+" {")
	if disabledIndex < 0 {
		t.Fatalf("icon-button source CSS is missing neutral disabled rule %q", disabledSelector)
	}
	for _, selector := range selectors {
		if stateIndex := strings.Index(css, selector+" {"); stateIndex >= disabledIndex && stateIndex > 0 {
			t.Errorf("neutral disabled rule must follow interactive state selector %q", selector)
		}
	}
}

func TestIconButtonVariantSelectorsMapToMaterialTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceIconButtonCSS(t), " ")
	for _, contract := range []string{
		`.ui-icon-button-standard { background: transparent; color: var(--ui-color-fg-muted);`,
		`.ui-icon-button-filled { background: var(--ui-color-primary); color: var(--ui-color-primary-fg);`,
		`.ui-icon-button-filled-tonal { background: var(--ui-color-secondary); color: var(--ui-color-secondary-fg);`,
		`.ui-icon-button-outlined { border-color: var(--ui-color-border-strong); background: transparent; color: var(--ui-color-fg-muted);`,
		`.ui-icon-button:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
		`.ui-icon-button:disabled, .ui-icon-button[aria-disabled="true"] { opacity: var(--ui-state-disabled-opacity); cursor: not-allowed; box-shadow: none;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("icon-button variant/state contract is missing %q", contract)
		}
	}
}

func TestIconButtonToggleSelectedRaisesEmphasisToPrimary(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceIconButtonCSS(t), " ")
	for _, contract := range []string{
		`.ui-icon-button[aria-pressed="true"] { color: var(--ui-color-primary);`,
		`.ui-icon-button-outlined[aria-pressed="true"] { border-color: var(--ui-color-primary); background: var(--ui-color-primary); color: var(--ui-color-primary-fg);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("icon-button selected-state contract is missing %q", contract)
		}
	}
}

func TestIconButtonKeepsFixedTouchTargetDimensions(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceIconButtonCSS(t), " ")
	for _, contract := range []string{
		`.ui-icon-button { position: relative; display: inline-flex; width: var(--ui-touch-target); height: var(--ui-touch-target); align-items: center; justify-content: center; border: var(--ui-border-width-1) var(--ui-border-style-solid) transparent; border-radius: var(--ui-radius-full); padding: 0; flex: none; cursor: pointer; text-decoration: none;`,
		`.ui-icon-button .ui-icon { width: var(--ui-size-icon); height: var(--ui-size-icon);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("icon-button touch-target dimension contract is missing %q", contract)
		}
	}
}
