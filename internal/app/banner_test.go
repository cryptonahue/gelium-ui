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

// newDocsServer builds the same server pipeline New() uses so a docs handler
// can be exercised directly. The component routes are not registered yet
// (registration, nav and sitemap integration lands in main), so tests invoke
// the handler method on the constructed server instead of going through the
// mux. The markdown converter replicates New()'s config exactly: GFM,
// highlighting, auto heading IDs and — critically — WithUnsafe so the raw
// HTML specimens embedded in the content markdown render as live markup.
func newDocsServer(t *testing.T) *server {
	t.Helper()
	return &server{
		templates: parseTestTemplates(t, "templates/*.html"),
		markdown: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(html.WithClasses(true)),
				),
			),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
		),
		assets: webassets.Assets,
	}
}

func TestBannerDocsRendersShellTitleAndLiveToneSpecimens(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.bannerDocs(res, httptest.NewRequest(http.MethodGet, "/components/banner", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("banner docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Banner</h1>`,
		`<title>Banner · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/banner"`,
		`class="docs-shell"`,
		`class="ui-banner ui-banner--error"`,
		`role="alert"`,
		`class="ui-banner ui-banner--info"`,
		`role="status"`,
		`class="ui-banner-title"`,
		`class="ui-banner-body"`,
		`class="ui-banner-dismiss"`,
		`aria-label="Dismiss"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("banner docs are missing %q", contract)
		}
	}
}

func TestBannerDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.bannerDocs(res, httptest.NewRequest(http.MethodGet, "/components/banner", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("banner docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use a banner when",
		"not a transient toast",
		"no component JavaScript",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("banner answer-first lead is missing %q", contract)
		}
	}
}