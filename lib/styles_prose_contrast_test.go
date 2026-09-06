package lib

// Prose contrast contract (WCAG AA).
//
// The docs prose column is text over the docs-shell surface
// (--ui-color-surface), colored by --ui-color-fg-muted (asserted by
// TestProseColorUsesFgMutedToken). Body text must meet WCAG 2.1 AA contrast
// (>= 4.5:1) in every served color context: the unthemed core defaults plus
// each theme on disk, in both the light route and the dark class route.
// Contrast is a semantic-token responsibility: the theme owns the values,
// this test proves the ratio.

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// tokenHex extracts the hex value of a --ui-* token from a theme or token
// file. When the token is redefined by the dark class route (which always
// comes after the light route in these files), dark=true returns the last
// definition.
func tokenHex(t *testing.T, css, name string, dark bool) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)` + regexp.QuoteMeta(name) + `:\s*(#[0-9a-fA-F]{6})\s*;`)
	matches := re.FindAllStringSubmatch(css, -1)
	if len(matches) == 0 {
		t.Fatalf("token %s not found in CSS", name)
	}
	idx := 0
	if dark {
		idx = len(matches) - 1
	}
	return matches[idx][1]
}

// channelLinear converts one sRGB channel (0-255) to linear light per the
// WCAG 2.1 relative-luminance formula.
func channelLinear(c uint8) float64 {
	v := float64(c) / 255
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// relativeLuminance computes WCAG 2.1 relative luminance for a #rrggbb hex.
func relativeLuminance(t *testing.T, hex string) float64 {
	t.Helper()
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		t.Fatalf("token value %q is not #rrggbb", hex)
	}
	ch := func(off int) float64 {
		v, err := strconv.ParseUint(hex[off:off+2], 16, 8)
		if err != nil {
			t.Fatalf("parse %q: %v", hex, err)
		}
		return channelLinear(uint8(v))
	}
	return 0.2126*ch(0) + 0.7152*ch(2) + 0.0722*ch(4)
}

// contrastRatio returns the WCAG 2.1 contrast ratio between two colors
// (>= 1.0; 21:1 is the maximum).
func contrastRatio(t *testing.T, a, b string) float64 {
	t.Helper()
	la := relativeLuminance(t, a)
	lb := relativeLuminance(t, b)
	if lb > la {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// TestProseColorUsesFgMutedToken pins the prose color chain: .prose p must
// color itself with the --ui-color-fg-muted semantic token (never a literal),
// so the contrast contract below applies to the actual served prose text.
func TestProseColorUsesFgMutedToken(t *testing.T) {
	p := cssBlock(t, readSourceStyle(t, "docs-chrome.css"), ".prose p")
	if !strings.Contains(p, "color: var(--ui-color-fg-muted)") {
		t.Errorf(".prose p must color via var(--ui-color-fg-muted), got block: %s", p)
	}
}

// TestCodeHighlightingContrastMeetsWcagAA enforces the code-block AA
// contract: chroma string color (--ui-color-code-string) on the code block
// background (--ui-color-surface-container) must reach >= 4.5:1 in every
// served color context and both routes. Guards the regression where
// --ui-color-secondary (a surface token) was used as the string color,
// failing at 1.0-1.1:1 in light themes.
func TestCodeHighlightingContrastMeetsWcagAA(t *testing.T) {
	type context struct {
		name string
		css  string
	}
	contexts := []context{{name: "core", css: readSourceStyle(t, "tokens.css")}}
	for _, theme := range availableThemes(t) {
		contexts = append(contexts, context{name: "theme-" + theme, css: themeCSS(t, theme)})
	}
	for _, ctx := range contexts {
		for _, dark := range []bool{false, true} {
			route := "light"
			if dark {
				route = "dark"
			}
			t.Run(ctx.name+"/"+route, func(t *testing.T) {
				str := tokenHex(t, ctx.css, "--ui-color-code-string", dark)
				bg := tokenHex(t, ctx.css, "--ui-color-surface-container", dark)
				ratio := contrastRatio(t, str, bg)
				if ratio < 4.5 {
					t.Errorf("code string %s on surface-container %s = %.2f:1, want >= 4.5:1 (WCAG AA body text)", str, bg, ratio)
				}
			})
		}
	}
}

// TestProseContrastMeetsWcagAA enforces the prose AA contract: in every
// served color context — core defaults plus each theme on disk, light and
// dark routes — --ui-color-fg-muted on --ui-color-surface must reach
// >= 4.5:1 (WCAG 2.1 AA for body text).
func TestProseContrastMeetsWcagAA(t *testing.T) {
	type context struct {
		name string
		css  string
	}
	contexts := []context{{name: "core", css: readSourceStyle(t, "tokens.css")}}
	for _, theme := range availableThemes(t) {
		contexts = append(contexts, context{name: "theme-" + theme, css: themeCSS(t, theme)})
	}
	for _, ctx := range contexts {
		for _, dark := range []bool{false, true} {
			route := "light"
			if dark {
				route = "dark"
			}
			t.Run(ctx.name+"/"+route, func(t *testing.T) {
				fg := tokenHex(t, ctx.css, "--ui-color-fg-muted", dark)
				bg := tokenHex(t, ctx.css, "--ui-color-surface", dark)
				ratio := contrastRatio(t, fg, bg)
				if ratio < 4.5 {
					t.Errorf("prose fg-muted %s on surface %s = %.2f:1, want >= 4.5:1 (WCAG AA body text)", fg, bg, ratio)
				}
			})
		}
	}
}

// TestNeubrutalismPrimaryContrastMeetsWcagAA guards both semantic uses of the
// primary pair: primary is rendered directly as interactive text on a surface,
// while primary-fg is rendered on primary-filled controls. Both combinations
// must remain readable in light and dark routes.
func TestNeubrutalismPrimaryContrastMeetsWcagAA(t *testing.T) {
	css := themeCSS(t, "theme-neubrutalism")
	for _, dark := range []bool{false, true} {
		route := "light"
		if dark {
			route = "dark"
		}
		t.Run(route, func(t *testing.T) {
			primary := tokenHex(t, css, "--ui-color-primary", dark)
			primaryFG := tokenHex(t, css, "--ui-color-primary-fg", dark)

			for _, backgroundToken := range []string{"--ui-color-canvas", "--ui-color-surface", "--ui-color-surface-container"} {
				background := tokenHex(t, css, backgroundToken, dark)
				if ratio := contrastRatio(t, primary, background); ratio < 4.5 {
					t.Errorf("primary text %s on %s %s = %.2f:1, want >= 4.5:1", primary, backgroundToken, background, ratio)
				}
			}
			if ratio := contrastRatio(t, primaryFG, primary); ratio < 4.5 {
				t.Errorf("primary foreground %s on primary %s = %.2f:1, want >= 4.5:1", primaryFG, primary, ratio)
			}
		})
	}
}
