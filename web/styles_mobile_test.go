package web

import (
	"regexp"
	"strings"
	"testing"
)

// Layer 1 (mobile foundations) contracts: touch-target and container-cap
// tokens, the documented breakpoint ladder, and the step-based type scale
// finding. Assertions are source-level (web/styles via the shared embed) plus
// compiled-bundle presence checks (npm run build regenerates
// web/static/app.css, committed with the change).

// TestMobileTokensDefinedInCore proves the mobile foundation tokens exist in
// the core contract: --ui-touch-target (interactive hit-area floor) and
// --ui-container-max (consumer layout cap).
func TestMobileTokensDefinedInCore(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	for _, token := range []string{"--ui-touch-target:", "--ui-container-max:"} {
		if !strings.Contains(string(css), token) {
			t.Errorf("core tokens.css must define %s (mobile foundation)", token)
		}
	}
}

// TestTouchTargetPinnedDefaultsAndThemeOwnership pins the ONE documented
// choice: the core floor is 44px (clears the GOV.UK 40px minimum), the
// Material theme overrides it to 48px (Material 3 48x48, matching USWDS's
// 48px menu button), and frozen theme-basecoat stays untouched so it
// inherits the core floor.
func TestTouchTargetPinnedDefaultsAndThemeOwnership(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(core), " ")
	if !strings.Contains(compact, "--ui-touch-target: 44px;") {
		t.Error("core --ui-touch-target default must be pinned at 44px (>=44px floor, GOV.UK 40px minimum)")
	}
	for _, ref := range []string{"M3", "USWDS", "GOV.UK"} {
		if !strings.Contains(string(core), ref) {
			t.Errorf("core tokens.css must document the touch-target references (%s)", ref)
		}
	}

	material := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(material, "--ui-touch-target: 48px;") {
		t.Error("theme-material must override --ui-touch-target at 48px (Material 3 48x48 touch target)")
	}

	basecoat := themeCSS(t, "theme-basecoat")
	if strings.Contains(basecoat, "--ui-touch-target") {
		t.Error("theme-basecoat is frozen and must not override --ui-touch-target (inherits the core 44px floor)")
	}
}

// TestMobileTokensWiredIntoComponents proves the foundation tokens have real
// consumers, not orphan definitions: button and icon-button reach the
// hit-area floor via the touch-target token, and the recipe consumer layout
// resolves the container cap.
func TestMobileTokensWiredIntoComponents(t *testing.T) {
	for file, needle := range map[string]string{
		"styles/button.css":                "min-height: var(--ui-touch-target);",
		"styles/icon-button.css":           "width: var(--ui-touch-target); height: var(--ui-touch-target);",
		"styles/recipe-admin-resource.css": "max-width: var(--ui-container-max);",
	} {
		css, err := sourceStyles.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(css), " ")
		if !strings.Contains(compact, needle) {
			t.Errorf("%s must consume the mobile foundation token via %q", file, needle)
		}
	}
}

// TestCoreTokensStayStepBased pins the GOV.UK/USWDS scale finding: the type
// and spacing scales stay step-based, so tokens.css must never introduce
// clamp() fluid steps.
func TestCoreTokensStayStepBased(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if strings.Contains(string(css), "clamp(") {
		t.Error("core tokens.css must stay step-based (GOV.UK/USWDS finding): no clamp() fluid steps")
	}
}

// TestBreakpointLadderDocumented proves the four-rung breakpoint ladder is
// documented in a tokens.css comment (GOV.UK/USWDS reference ladder). The
// marker anchors the scan so the rungs must come from the comment, not from
// token values.
func TestBreakpointLadderDocumented(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	const marker = "Breakpoint ladder"
	text := string(css)
	i := strings.Index(text, marker)
	if i < 0 {
		t.Fatalf("tokens.css must document the breakpoint ladder with a %q comment marker", marker)
	}
	ladder := text[i:]
	for _, rung := range []string{"320px", "480px", "48rem", "64rem"} {
		if !strings.Contains(ladder, rung) {
			t.Errorf("breakpoint ladder comment must include rung %q", rung)
		}
	}
}

// TestMobileTokensReachCompiledBundle proves npm run build regenerates
// web/static/app.css carrying the mobile foundation tokens alongside their
// source definitions.
func TestMobileTokensReachCompiledBundle(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read compiled bundle: %v", err)
	}
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(compiled), " ")
	for _, token := range []string{"--ui-touch-target:", "--ui-container-max:"} {
		if !strings.Contains(compact, token) {
			t.Errorf("compiled static/app.css must carry %s (regenerate with npm run build)", token)
		}
	}
}
