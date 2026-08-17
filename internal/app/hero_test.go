package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"

	webassets "geliumui/site/web"
)

// renderDocsMethod renders a /components/* docs route through the same
// template + markdown stack the production server uses (buildTemplates +
// the New() goldmark options, including WithUnsafe so the live raw-HTML
// specimens in content/*.md render). Only the given route is registered on a
// fresh mux: on this partial branch the routes are wired in main, so the
// handler is exercised in isolation instead. The method is passed through so
// both the GET contract and the 405 semantics are exercised through real mux
// routing.
func renderDocsMethod(t *testing.T, route, method string, handler func(*server, http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	s := &server{
		templates: buildTemplates(),
		markdown: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
				),
			),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
		),
		assets: webassets.Assets,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+route, func(w http.ResponseWriter, r *http.Request) {
		handler(s, w, r)
	})
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, httptest.NewRequest(method, route, nil))
	return res
}

// renderDocsPage is the GET shorthand of renderDocsMethod.
func renderDocsPage(t *testing.T, route string, handler func(*server, http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	return renderDocsMethod(t, route, http.MethodGet, handler)
}

func TestHeroDocsRouteDogfoodsLiveSpecimenAndShell(t *testing.T) {
	res := renderDocsPage(t, "/components/hero", (*server).heroDocs)

	if res.Code != http.StatusOK {
		t.Fatalf("hero docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Hero</h1>`,
		`class="docs-shell"`,
		`aria-label="Docs"`,
		`class="ui-hero"`,
		`class="ui-hero-content"`,
		`class="ui-hero-eyebrow"`,
		`class="ui-hero-title"`,
		`class="ui-hero-subtitle"`,
		`class="ui-hero-actions"`,
		`class="ui-button ui-button-primary"`,
		`class="ui-hero-media"`,
		`class="ui-hero ui-hero--has-media"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("hero docs are missing %q", contract)
		}
	}
}

func TestHeroDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := renderDocsMethod(t, "/components/hero", http.MethodPost, (*server).heroDocs)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST hero status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}