package lib

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Layout utilities contract: .ui-container (page column cap), measure/prose
// ≤75ch, stack, and *-from-desktop row escape hatches (min-width: 48rem).

func TestLayoutSourceFileExists(t *testing.T) {
	if _, err := sourceStyles.ReadFile("styles/layout.css"); err != nil {
		t.Fatalf("styles/layout.css must exist in the lib styles embed: %v", err)
	}
}

func TestLayoutIndexManifestImportsLayout(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/index.css")
	if err != nil {
		t.Fatalf("read index.css: %v", err)
	}
	if !strings.Contains(string(css), `@import "./layout.css";`) {
		t.Error(`styles/index.css must import "./layout.css"`)
	}
}

func TestLayoutContainerAndFromDesktopContracts(t *testing.T) {
	raw, err := sourceStyles.ReadFile("styles/layout.css")
	if err != nil {
		t.Fatalf("read layout.css: %v", err)
	}
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(raw), " ")

	for _, contract := range []string{
		`.ui-container`,
		`max-width: var(--ui-container-max)`,
		`width: 100%`,
		`margin-inline: auto`,
		`padding-inline: var(--ui-space-4)`,
		`min-width: 0`,
		`box-sizing: border-box`,
		`.ui-container-prose`,
		`.ui-measure`,
		`max-width: 65ch`,
		`.ui-stack`,
		`.ui-row-from-desktop`,
		`from-desktop`,
		`@media (min-width: 48rem)`,
		`flex-direction: column`,
		`flex-direction: row`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("layout.css is missing contract %q", contract)
		}
	}

	// Header documents the stack-default / enhance-from-desktop reflow line.
	if !strings.Contains(string(raw), "stack by default") && !strings.Contains(string(raw), "enhance from-desktop") {
		t.Error("layout.css header must document stack-by-default and enhance from-desktop")
	}
}

func TestLayoutUtilitiesReachCompiledBundle(t *testing.T) {
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(compiledAppCSS(t), " ")
	for _, needle := range []string{
		`.ui-container`,
		`max-width:var(--ui-container-max)`,
		`from-desktop`,
		`.ui-row-from-desktop`,
		`.ui-stack`,
		`.ui-measure`,
	} {
		// Minifier may drop spaces around :; accept either form.
		alt := strings.ReplaceAll(needle, ": ", ":")
		if !strings.Contains(compact, needle) && !strings.Contains(compact, alt) {
			// Also try fully compacted (no spaces).
			full := regexp.MustCompile(`\s+`).ReplaceAllString(needle, "")
			if !strings.Contains(regexp.MustCompile(`\s+`).ReplaceAllString(compact, ""), full) {
				t.Errorf("compiled static/app.css must carry layout utility %q (regenerate with npm run build)", needle)
			}
		}
	}
}

// TestProseMeasureStillWithin75ch keeps the readability ceiling: docs .prose
// and the public measure utilities must use a ch cap of at most 75.
func TestProseMeasureStillWithin75ch(t *testing.T) {
	const maxAllowed = 75

	prose := cssBlock(t, readSourceStyle(t, "docs-chrome.css"), ".prose")
	if !strings.Contains(prose, "max-width: 65ch") {
		t.Errorf(".prose must keep max-width: 65ch, got block: %s", prose)
	}
	assertChMeasureAtMost(t, prose, maxAllowed, ".prose")

	layout, err := sourceStyles.ReadFile("styles/layout.css")
	if err != nil {
		t.Fatalf("read layout.css: %v", err)
	}
	layoutText := string(layout)
	if !strings.Contains(layoutText, "max-width: 65ch") {
		t.Error("layout.css measure utilities must set max-width: 65ch")
	}
	assertChMeasureAtMost(t, layoutText, maxAllowed, "layout.css")

	// Explicit numeric contract: 65 <= 75.
	const measureCh = 65
	if measureCh > maxAllowed {
		t.Fatalf("prose/measure ch %d exceeds ceiling %d", measureCh, maxAllowed)
	}
	if !(measureCh <= maxAllowed) {
		t.Fatalf("assert 65 <= 75 failed: %d <= %d", measureCh, maxAllowed)
	}
}

func assertChMeasureAtMost(t *testing.T, css string, max int, label string) {
	t.Helper()
	re := regexp.MustCompile(`max-width:\s*([0-9]+)ch`)
	matches := re.FindAllStringSubmatch(css, -1)
	if len(matches) == 0 {
		t.Errorf("%s: expected at least one max-width: Nch declaration", label)
		return
	}
	for _, m := range matches {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			t.Errorf("%s: bad ch measure %q: %v", label, m[1], err)
			continue
		}
		if n > max {
			t.Errorf("%s: max-width %dch exceeds ≤%dch contract", label, n, max)
		}
	}
}
