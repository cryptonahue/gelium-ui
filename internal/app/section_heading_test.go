package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"

	webassets "geliumui/site/web"
)

// docsTestServer builds a server identical to New() (same template merge, same
// goldmark config with WithUnsafe) without registering routes. Partial-branch
// component pages are not wired into routes.go yet, so their contract tests
// exercise the handler directly against the real embedded assets.
func docsTestServer(t *testing.T) *server {
	t.Helper()
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(html.WithClasses(true)),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
	)
	return &server{
		templates: parseTestTemplates(t, "templates/*.html"),
		markdown:  md,
		assets:    webassets.Assets,
	}
}

func TestSectionHeadingDocsRouteDogfoodsH2OnlySemantics(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.sectionHeadingDocs(res, httptest.NewRequest(http.MethodGet, "/components/section-heading", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("section heading docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Section heading</h1>`,
		`<p class="ui-section-heading-eyebrow">Eyebrow</p>`,
		`<h2 class="ui-section-heading">Features</h2>`,
		`<h2 class="ui-section-heading ui-section-heading--centered">Features</h2>`,
		`ui-section-heading-eyebrow`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("section heading docs are missing %q", contract)
		}
	}
}

func TestSectionHeadingDocsExposesNoH1InSpecimens(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.sectionHeadingDocs(res, httptest.NewRequest(http.MethodGet, "/components/section-heading", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("section heading docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// The page must lead with exactly one h1 and every specimen heading must
	// stay an h2 (the utility never renders h1).
	if got := strings.Count(body, "<h1"); got != 1 {
		t.Errorf("section heading docs render %d h1 elements, want exactly 1 (the page title)", got)
	}
}