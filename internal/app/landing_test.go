package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLandingPreservesThemeAndSchemeOnCTAs(t *testing.T) {
	body := getOKBody(t, "/?theme=basecoat&scheme=dark")
	if !strings.Contains(body, `theme-basecoat`) || !strings.Contains(body, `theme-dark`) {
		t.Fatalf("landing root must apply theme+scheme; got %s", htmlClassSnippet(body))
	}
	// Primary CTA keeps both query keys.
	if !strings.Contains(body, `href="/docs?`) {
		t.Fatal("missing docs CTA")
	}
	// Scan Get started / docs button href for both keys.
	found := false
	rest := body
	for {
		i := strings.Index(rest, `href="/docs?`)
		if i < 0 {
			break
		}
		rest = rest[i:]
		j := strings.Index(rest[6:], `"`)
		href := rest[:6+j+1]
		if strings.Contains(href, `theme=basecoat`) && strings.Contains(href, `scheme=dark`) {
			found = true
			break
		}
		rest = rest[6:]
	}
	if !found {
		t.Error("landing docs CTA must preserve theme=basecoat and scheme=dark")
	}
}

func TestLandingPrimaryNavIsCompact(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	body := res.Body.String()
	// Compact marketing nav — not the full 28-component dump.
	for _, label := range []string{">Docs<", ">Components<", ">Recipes<", ">Agents<", ">Demo<"} {
		if !strings.Contains(body, label) {
			t.Errorf("landing header missing nav label marker %q", label)
		}
	}
	// Inspect only the primary header nav (footer still lists the catalog).
	start := strings.Index(body, `aria-label="Primary"`)
	if start < 0 {
		t.Fatal("missing primary nav")
	}
	end := strings.Index(body[start:], `</nav>`)
	if end < 0 {
		t.Fatal("unclosed primary nav")
	}
	primary := body[start : start+end]
	if strings.Count(primary, `href="/components/`) > 3 {
		t.Errorf("primary nav should be compact; got %d component hrefs", strings.Count(primary, `href="/components/`))
	}
	if !strings.Contains(primary, `href="/docs"`) {
		t.Error("primary nav must link Docs")
	}
	if !strings.Contains(primary, `href="/docs/agent-workflow"`) {
		t.Error("primary nav must link Agents workflow")
	}
}

// TestLandingSinglePrimaryCTA proves Persuade hierarchy: one filled primary
// button on the marketing home (hero Get started).
func TestLandingSinglePrimaryCTA(t *testing.T) {
	body := getOKBody(t, "/")
	// Class token is "ui-button-primary" as a whole word in class attrs.
	n := strings.Count(body, "ui-button-primary")
	if n != 1 {
		t.Fatalf("landing must render exactly one ui-button-primary, got %d", n)
	}
	if !strings.Contains(body, "Get started") {
		t.Error("landing primary must remain Get started")
	}
}
