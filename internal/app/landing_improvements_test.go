package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLandingShowsFAQ proves the landing closes with an FAQ section (Base UI
// pattern) answering the four developer objections, rendered as zero-JS
// native <details> disclosures.
func TestLandingShowsFAQ(t *testing.T) {
	body := getOKBody(t, "/")

	// The section itself must exist with a real heading.
	for _, want := range []string{
		`class="ui-landing-section ui-landing-faq"`,
		"Frequently asked questions",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("landing FAQ section missing %q", want)
		}
	}

	// Every question must render as a disclosure with its answer fragment.
	qa := []struct{ question, answer string }{
		{
			question: "Is Gelium UI a framework?",
			answer:   "no SPA runtime to install",
		},
		{
			question: "Does it work with my existing stack?",
			answer:   "HTML partials and CSS",
		},
		{
			question: "How does it differ from a React component library?",
			answer:   "server contracts instead of props",
		},
		{
			question: "Can I switch themes?",
			answer:   "a class on the html element",
		},
	}
	for _, item := range qa {
		if !strings.Contains(body, item.question) {
			t.Errorf("landing FAQ missing question %q", item.question)
		}
		if !strings.Contains(body, item.answer) {
			t.Errorf("landing FAQ missing answer fragment %q", item.answer)
		}
	}

	// Exactly four zero-JS disclosures must be rendered inside the FAQ section;
	// the marketing shell's responsive menu has its own native summary.
	start := strings.Index(body, `ui-landing-faq`)
	if start < 0 {
		t.Fatal("landing FAQ section missing before disclosure count")
	}
	end := strings.Index(body[start:], `</section>`)
	faq := body[start:]
	if end >= 0 {
		faq = faq[:end]
	}
	if got := strings.Count(faq, "<summary>"); got != 4 {
		t.Errorf("landing FAQ rendered %d <summary> disclosures, want 4", got)
	}
}

// TestLandingClaimsStrip proves the Naive-UI-style checkmark strip renders
// under the hero with all four claims and their checkmarks.
func TestLandingClaimsStrip(t *testing.T) {
	body := getOKBody(t, "/")

	if !strings.Contains(body, `class="ui-landing-claims"`) {
		t.Fatal("landing claims strip container missing")
	}

	claims := []string{
		"HTML-first",
		"No-JS baseline",
		"Server-first state",
		"Themes without forks",
	}
	// html/template escapes '+' as &#43; in text; normalize so the assertion
	// checks the claim copy, not the escaping mechanism.
	normalized := strings.ReplaceAll(body, "&#43;", "+")
	for _, claim := range claims {
		if !strings.Contains(normalized, claim) {
			t.Errorf("landing claims strip missing claim %q", claim)
		}
	}

	// One checkmark glyph per claim — no other ✓ on the landing.
	if got := strings.Count(body, "✓"); got != 4 {
		t.Errorf("landing claims strip rendered %d checkmarks, want 4", got)
	}
}

// TestLandingDemoCard proves the visual demo link card (Basecoat pattern)
// links the live WhatsApp demo from a token-styled preview, no image file.
func TestLandingDemoCard(t *testing.T) {
	body := getOKBody(t, "/")

	for _, want := range []string{
		`class="ui-landing-demo-card"`,
		`class="ui-landing-demo-phone"`,
		"Launch live demo",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("landing demo card missing %q", want)
		}
	}
	if !strings.Contains(body, `class="ui-button ui-button-secondary" href="/demo/whatsapp"`) {
		t.Error("landing demo card CTA must be secondary and link /demo/whatsapp")
	}
}

// TestLandingDemoCardKeepsChrome proves the demo card CTA is an internal link
// and therefore participates in chrome rewriting (unlike the external GitHub
// CTA, which must stay clean). The scan anchors on the demo card's primary
// button so the marketing nav's own chrome-rewritten Demo link cannot satisfy
// the assertion.
func TestLandingDemoCardKeepsChrome(t *testing.T) {
	body := getOKBody(t, "/?theme=basecoat&scheme=dark")

	anchor := `class="ui-button ui-button-secondary" href="/demo/whatsapp?`
	i := strings.Index(body, anchor)
	if i < 0 {
		t.Fatalf("demo card CTA missing (anchor %q)", anchor)
	}
	href := body[i+len(anchor) : i+len(anchor)+160]
	if !strings.Contains(href, `theme=basecoat`) || !strings.Contains(href, `scheme=dark`) {
		t.Errorf("demo card CTA must carry chrome query like other internal landing CTAs; got %q", href)
	}
}

// TestLandingUnchangedStatus proves the augmented landing still renders 200
// with the single-h1 hero untouched.
func TestLandingUnchangedStatus(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := strings.Count(res.Body.String(), "<h1"); got != 1 {
		t.Errorf("landing rendered %d <h1> elements, want 1", got)
	}
}
