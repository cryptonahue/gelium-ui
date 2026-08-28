package lib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDistBundleIsCommittedAndComplete pins the S5 dist contract: the
// prebuilt bundle (lib/dist/gelium.css) must exist, be non-trivial, and carry
// the component primitives an unthemed consumer relies on.
func TestDistBundleIsCommittedAndComplete(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "dist", "gelium.css"))
	if err != nil {
		t.Fatalf("read lib/dist/gelium.css: %v", err)
	}
	if len(b) < 50_000 {
		t.Errorf("dist bundle suspiciously small (%d bytes); expected the full component set", len(b))
	}
	for _, contract := range []string{`.ui-button`, `.ui-toast`, `.ui-text-field`, `--ui-touch-target`} {
		if !strings.Contains(string(b), contract) {
			t.Errorf("dist bundle missing %q", contract)
		}
	}
	// The dist bundle MUST carry the theme roots: themes ship WITH the
	// package (lib/themes/*.css) so a consumer installing gelium-ui gets a
	// working theme pair, not bare components.
	if !strings.Contains(string(b), ".theme-material") || !strings.Contains(string(b), ".theme-basecoat") {
		t.Error("dist bundle must embed both theme roots (.theme-material, .theme-basecoat)")
	}
}

// TestDistBundleAppliesThemesAfterCoreDefaults pins the public bundle cascade:
// each theme declaration must follow the corresponding core default so it can
// override it on the same element. The parser intentionally checks rule and
// declaration structure instead of Tailwind's exact serialization.
func TestDistBundleAppliesThemesAfterCoreDefaults(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "dist", "gelium.css"))
	if err != nil {
		t.Fatalf("read lib/dist/gelium.css: %v", err)
	}

	rules := parseCSSRules(string(b))
	core, ok := findDistCSSRuleWithDeclaration(rules, ":root", "--ui-color-canvas")
	if !ok {
		t.Fatal("dist bundle missing the core :root rule")
	}

	// Surface-container values differ from the core default in both themes,
	// making the check prove that representative theme declarations are used.
	for _, themeName := range []string{".theme-material", ".theme-basecoat"} {
		theme, ok := findDistCSSRule(rules, themeName)
		if !ok {
			t.Fatalf("dist bundle missing the %s rule", themeName)
		}
		for _, property := range []string{"--ui-color-canvas", "--ui-color-surface-container"} {
			if _, ok := core.declarations[property]; !ok {
				t.Fatalf("core :root rule missing representative %s declaration", property)
			}
			if _, ok := theme.declarations[property]; !ok {
				t.Fatalf("%s rule missing representative %s declaration", themeName, property)
			}
		}
		if core.declarations["--ui-color-surface-container"] == theme.declarations["--ui-color-surface-container"] {
			t.Fatalf("%s representative surface-container declaration does not differ from the core default", themeName)
		}
		if theme.order <= core.order {
			t.Fatalf("%s rule appears before the core :root rule (theme rule %d, core rule %d)", themeName, theme.order, core.order)
		}
	}
}

type parsedCSSRule struct {
	selector     string
	declarations map[string]string
	order        int
}

// parseCSSRules extracts ordinary rules in source order, recursing through
// at-rules such as @layer and @media. It is deliberately a small CSS parser:
// this contract only needs selectors, declaration names, and rule ordering.
func parseCSSRules(css string) []parsedCSSRule {
	rules := make([]parsedCSSRule, 0)
	var parse func(string)
	parse = func(input string) {
		for pos := 0; pos < len(input); {
			open := strings.IndexByte(input[pos:], '{')
			if open < 0 {
				return
			}
			open += pos
			close := matchingBrace(input, open)
			if close < 0 {
				return
			}
			selector := normalizeDistCSSWhitespace(input[pos:open])
			body := input[open+1 : close]
			if strings.HasPrefix(selector, "@") {
				parse(body)
			} else {
				rules = append(rules, parsedCSSRule{
					selector:     selector,
					declarations: parseCSSDeclarations(body),
					order:        len(rules),
				})
			}
			pos = close + 1
		}
	}
	parse(css)
	return rules
}

func matchingBrace(css string, open int) int {
	depth := 0
	var quote byte
	for i := open; i < len(css); i++ {
		if quote != 0 {
			if css[i] == '\\' {
				i++
				continue
			}
			if css[i] == quote {
				quote = 0
			}
			continue
		}
		if css[i] == '\'' || css[i] == '"' {
			quote = css[i]
			continue
		}
		if css[i] == '/' && i+1 < len(css) && css[i+1] == '*' {
			if end := strings.Index(css[i+2:], "*/"); end >= 0 {
				i += end + 3
				continue
			}
			return -1
		}
		switch css[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func parseCSSDeclarations(body string) map[string]string {
	declarations := make(map[string]string)
	for _, declaration := range strings.Split(body, ";") {
		colon := strings.IndexByte(declaration, ':')
		if colon < 0 {
			continue
		}
		property := strings.TrimSpace(declaration[:colon])
		if property == "" || strings.HasPrefix(property, "@") {
			continue
		}
		declarations[property] = strings.TrimSpace(declaration[colon+1:])
	}
	return declarations
}

func findDistCSSRuleWithDeclaration(rules []parsedCSSRule, selector, property string) (parsedCSSRule, bool) {
	selector = normalizeDistCSSWhitespace(selector)
	for _, rule := range rules {
		if _, ok := rule.declarations[property]; !ok {
			continue
		}
		for _, candidate := range strings.Split(rule.selector, ",") {
			if normalizeDistCSSWhitespace(candidate) == selector {
				return rule, true
			}
		}
	}
	return parsedCSSRule{}, false
}

func findDistCSSRule(rules []parsedCSSRule, selector string) (parsedCSSRule, bool) {
	selector = normalizeDistCSSWhitespace(selector)
	for _, rule := range rules {
		for _, candidate := range strings.Split(rule.selector, ",") {
			if normalizeDistCSSWhitespace(candidate) == selector {
				return rule, true
			}
		}
	}
	return parsedCSSRule{}, false
}

func normalizeDistCSSWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// TestLibJsShipsConsumerEnhancements pins the S5 JS contract: lib/js/gelium.js
// carries the consumer enhancements (toast region, 422 contract, VT flag,
// slider fill) and the docs chrome (theme/scheme optimistic toggle) does NOT.
func TestLibJsShipsConsumerEnhancements(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "js", "gelium.js"))
	if err != nil {
		t.Fatalf("read lib/js/gelium.js: %v", err)
	}
	js := string(b)
	for _, contract := range []string{`gelium:toast`, `X-Gelium-Validation`, `startViewTransition`, `--ui-slider-fill`} {
		if !strings.Contains(js, contract) {
			t.Errorf("gelium.js missing consumer enhancement %q", contract)
		}
	}
	if strings.Contains(js, "refreshChromeHrefs") {
		t.Error("gelium.js must not contain docs chrome (refreshChromeHrefs is site-only)")
	}
}
