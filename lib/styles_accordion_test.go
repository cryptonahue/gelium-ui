package lib

import (
	"strings"
	"testing"
)

func TestAccordionStylesCoverNativeStates(t *testing.T) {
	css := sourceComponentCSS(t, "accordion.css")
	for _, want := range []string{
		`.ui-accordion-item`,
		`.ui-accordion-trigger:focus-visible`,
		`.ui-accordion-item[open]`,
		`@media (prefers-reduced-motion: reduce)`,
		`@media (forced-colors: active)`,
		`--ui-touch-target`,
		`--ui-focus-thickness`,
		`--ui-motion-short`,
		`stroke: currentColor`,
		`padding-inline`,
	} {
		if !strings.Contains(css, want) {
			t.Errorf("accordion.css missing native/state contract %q", want)
		}
	}
	for _, forbidden := range []string{"[disabled]", `role=`, `aria-expanded`} {
		if strings.Contains(css, forbidden) {
			t.Errorf("accordion.css must not invent disabled/ARIA behavior %q", forbidden)
		}
	}
}

func TestAccordionSkinTokensCaptureDistinctReferenceDirections(t *testing.T) {
	basecoat := themeCSS(t, "theme-basecoat")
	for _, want := range []string{"--ui-accordion-root-border", "--ui-accordion-item-border", "--ui-accordion-trigger-min-height", "--ui-accordion-icon-rotation"} {
		if !strings.Contains(basecoat, want) {
			t.Errorf("basecoat accordion skin missing token %q", want)
		}
	}
	material := themeCSS(t, "theme-material")
	for _, want := range []string{"--ui-accordion-shadow", "--ui-accordion-radius: 4px", "--ui-accordion-icon-size: 1.5rem"} {
		if !strings.Contains(material, want) {
			t.Errorf("material accordion skin missing token %q", want)
		}
	}
	for _, theme := range []string{"theme-baseui", "theme-alden", "theme-linear", "theme-vercel"} {
		if !strings.Contains(themeCSS(t, theme), "--ui-accordion-") {
			t.Errorf("unrelated skin %s lost its accordion token set", theme)
		}
	}
}

func TestAccordionTokensReachEveryThemeAndCompiledBundle(t *testing.T) {
	for _, theme := range availableThemes(t) {
		css := themeCSS(t, theme)
		if !strings.Contains(css, "--ui-accordion-") {
			t.Errorf("%s must define accordion theme tokens", theme)
		}
	}
	compiled := repositoryFile(t, "site", "web", "static", "app.css")
	for _, want := range []string{".ui-accordion", ".ui-accordion-trigger:focus-visible", ".ui-accordion-item[open]"} {
		if !strings.Contains(compiled, want) {
			t.Errorf("compiled app.css missing accordion selector %q (run npm run build)", want)
		}
	}
}

func TestAccordionBehaviorClassesDoNotOwnVisualCSS(t *testing.T) {
	css := sourceComponentCSS(t, "accordion.css")
	if strings.Contains(css, ".ui-accordion--behavior-") {
		t.Fatalf("accordion behavior classes must be diagnostics/policy only; visual selectors remain:\n%s", css)
	}
	for _, forbidden := range []string{"min-height: 0", "--ui-accordion-icon-plus"} {
		if strings.Contains(css, forbidden) {
			t.Errorf("accordion core retains behavior-era visual ownership %q", forbidden)
		}
	}
	if !strings.Contains(css, "max(var(--ui-touch-target)") {
		t.Error("accordion trigger must keep the 44px core touch floor even for dense presets")
	}
	if !strings.Contains(css, "--ui-accordion-chevron-display") || !strings.Contains(css, "--ui-accordion-plus-display") {
		t.Error("accordion icon choice must be controlled by visual tokens, not behavior selectors")
	}
}
