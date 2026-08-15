package web

// Prose readability contract tests (readability-improvements).
//
// The docs text column is the primary reading surface of the Gelium UI
// documentation, so its typography is a contract, not a taste:
//
//  1. Measure — .prose caps line length at 65ch (the 65-75ch typographic
//     standard; USWDS uses ~68ch, GOV.UK ~30em). ch is the correct unit for
//     measure because it scales with the rendered font: the column always
//     carries roughly the same number of characters regardless of type size.
//  2. Text wrapping — body copy uses text-wrap: pretty (avoids widows and
//     orphaned final lines), headings use text-wrap: balance (even line
//     breaks on short heading blocks). Both are Baseline 2024.
//  3. Hyphenation — .prose sets hyphens: auto and inherits the document
//     language from <html lang>, so justified-ish narrow columns hyphenate
//     per language instead of overflowing.
//  4. Progressive leading trim — headings set text-box-trim: trim-both with
//     text-box-edge: cap alphabetic where supported (Chromium); unsupported
//     browsers ignore the pair, so this is a pure enhancement.
//  5. Readable line height — prose body copy keeps line-height >= 1.6.

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// readSourceStyle reads one file from the embedded web/styles source tree.
func readSourceStyle(t *testing.T, name string) string {
	t.Helper()
	css, err := sourceStyles.ReadFile("styles/" + name)
	if err != nil {
		t.Fatalf("read styles/%s: %v", name, err)
	}
	return string(css)
}

// cssBlock returns the rule block for an exact selector prefix (e.g. ".prose"
// or ".prose p"), from the opening brace to the matching close brace. The
// style files keep these rules flat (no nested braces), so brace counting is
// sufficient.
func cssBlock(t *testing.T, css, selector string) string {
	t.Helper()
	idx := strings.Index(css, selector+" {")
	if idx < 0 {
		t.Fatalf("CSS is missing rule %q", selector)
	}
	brace := strings.Index(css[idx:], "{")
	depth := 0
	for i := idx + brace; i < len(css); i++ {
		switch css[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return css[idx : i+1]
			}
		}
	}
	t.Fatalf("rule %q is unterminated", selector)
	return ""
}

// TestProseMeasureUsesChUnits is the readability contract: the prose text
// column must cap line length at 65ch, and must not fall back to a rem-based
// measure (the old 48rem ≈ 90 chars at the body size — too wide for comfort).
func TestProseMeasureUsesChUnits(t *testing.T) {
	prose := cssBlock(t, readSourceStyle(t, "base.css"), ".prose")
	if !strings.Contains(prose, "max-width: 65ch") {
		t.Errorf(".prose must cap line length at 65ch, got block: %s", prose)
	}
	if strings.Contains(prose, "48rem") {
		t.Errorf(".prose must use a ch measure, not rem: %s", prose)
	}
}

// TestProseBodyTextWrapPretty is the widows/orphans contract: body copy must
// set text-wrap: pretty so the final line of a paragraph never dangles alone.
func TestProseBodyTextWrapPretty(t *testing.T) {
	p := cssBlock(t, readSourceStyle(t, "base.css"), ".prose p")
	if !strings.Contains(p, "text-wrap: pretty") {
		t.Errorf(".prose p must set text-wrap: pretty, got block: %s", p)
	}
}

// TestProseBodyLineHeightStaysReadable pins the readable line-height floor:
// prose body copy must keep line-height >= 1.6 (current: 1.7).
func TestProseBodyLineHeightStaysReadable(t *testing.T) {
	p := cssBlock(t, readSourceStyle(t, "base.css"), ".prose p")
	m := regexp.MustCompile(`line-height:\s*([0-9]+(?:\.[0-9]+)?)`).FindStringSubmatch(p)
	if m == nil {
		t.Fatalf(".prose p must set a numeric line-height >= 1.6, got block: %s", p)
	}
	lh, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		t.Fatalf("parse line-height %q: %v", m[1], err)
	}
	if lh < 1.6 {
		t.Errorf(".prose p line-height = %v, want >= 1.6", lh)
	}
}

// TestProseHeadingsTextWrapBalance is the balanced-headings contract: h1-h3
// must set text-wrap: balance so multi-line headings break evenly.
func TestProseHeadingsTextWrapBalance(t *testing.T) {
	heads := cssBlock(t, readSourceStyle(t, "base.css"), ".prose h1, .prose h2, .prose h3")
	if !strings.Contains(heads, "text-wrap: balance") {
		t.Errorf(".prose h1/h2/h3 must set text-wrap: balance, got block: %s", heads)
	}
}

// TestProseHyphensAuto is the hyphenation contract: .prose must set
// hyphens: auto and inherit the document language from <html lang>.
func TestProseHyphensAuto(t *testing.T) {
	prose := cssBlock(t, readSourceStyle(t, "base.css"), ".prose")
	if !strings.Contains(prose, "hyphens: auto") {
		t.Errorf(".prose must set hyphens: auto, got block: %s", prose)
	}
}

// TestProseHeadingTextBoxTrimProgressive is the progressive-enhancement
// contract: headings trim visual leading with text-box-trim/text-box-edge
// where supported (Chromium only for now). The pair must be declared on the
// heading rule so a future browser ships it automatically.
func TestProseHeadingTextBoxTrimProgressive(t *testing.T) {
	heads := cssBlock(t, readSourceStyle(t, "base.css"), ".prose h1, .prose h2, .prose h3")
	if !strings.Contains(heads, "text-box-trim: trim-both") {
		t.Errorf(".prose h1/h2/h3 must set text-box-trim: trim-both (progressive), got block: %s", heads)
	}
	if !strings.Contains(heads, "text-box-edge: cap alphabetic") {
		t.Errorf(".prose h1/h2/h3 must set text-box-edge: cap alphabetic (progressive), got block: %s", heads)
	}
}
