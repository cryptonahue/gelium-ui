package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestProgressPrimitiveCSSMapsNativeProgressAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-progress {`,
		`progress {`,
		`appearance: none;`,
		`height: var(--ui-progress-track-height);`,
		`border-radius: var(--ui-progress-radius);`,
		`::-webkit-progress-bar {`,
		`::-webkit-progress-value {`,
		`var(--ui-progress-indicator)`,
		`::-moz-progress-bar {`,
		`.ui-progress--circular {`,
		`.ui-progress-circle {`,
		`animation: ui-spin .8s linear infinite;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing progress contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`::-webkit-progress-value`,
		`::-moz-progress-bar`,
		`.ui-progress-circle`,
		`border-top-color: Highlight;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("progress must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestProgressCircularVariantStopsUnderReducedMotion(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "progress.css"), " ")
	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("progress.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-progress-circle { animation: none; }`) {
		t.Error("circular progress reduced-motion must stop the ring sweep")
	}
}

func TestProgressThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-progress-track-height:",
		"--ui-progress-radius:",
		"--ui-progress-track:",
		"--ui-progress-indicator:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing progress token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesProgressContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-progress`,
		`-webkit-progress-value`,
		`-moz-progress-bar`,
		`.ui-progress--circular`,
		`ui-spin`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled progress CSS is missing %q", contract)
		}
	}
}
