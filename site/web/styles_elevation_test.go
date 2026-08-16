package web

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestElevationPrimitiveCSSMapsLevelsToShadowTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for level := 0; level <= 5; level++ {
		contract := fmt.Sprintf(".ui-elevation-%d { box-shadow: var(--ui-shadow-%d);", level, level)
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-elevation-0, .ui-elevation-1, .ui-elevation-2, .ui-elevation-3, .ui-elevation-4, .ui-elevation-5 { box-shadow: none;") {
		t.Error("elevation is visual-only and must not keep decorative shadows in forced colors")
	}

	reduced := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	if !strings.Contains(reduced, ".ui-elevation-0, .ui-elevation-1, .ui-elevation-2, .ui-elevation-3, .ui-elevation-4, .ui-elevation-5 { transition: none;") {
		t.Error("reduced-motion CSS must disable the elevation box-shadow transition")
	}
}

func TestEmbeddedCompiledCSSIncludesElevationContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-elevation-0{box-shadow:var(--ui-shadow-0)}`,
		`.ui-elevation-1{box-shadow:var(--ui-shadow-1)}`,
		`.ui-elevation-5{box-shadow:var(--ui-shadow-5)}`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled elevation CSS is missing %q", contract)
		}
	}
}
