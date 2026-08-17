package app

import (
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"geliumui/lib"
)

func TestMediaTemplatesExposeAccessibleContracts(t *testing.T) {
	tmpl := parseTestTemplates(t, "templates/media.html", "templates/image.html")
	for _, name := range []string{"image", "picture", "audio", "transcript", "embed"} {
		if tmpl.Lookup(name) == nil {
			t.Fatalf("missing template %q", name)
		}
	}
	contracts := map[string][]string{
		"templates/media.html": {"controls", "Fallback", "preload="},
		"templates/image.html": {"alt=", "width=", "height=", "loading=", "srcset=", "sizes="},
	}
	// The contract is asserted against the embedded source, avoiding browser timing or layout.
	for source, wants := range contracts {
		b, err := fs.ReadFile(lib.LibAssets, source)
		if err != nil {
			t.Fatal(err)
		}
		text := string(b)
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s missing %q", source, want)
			}
		}
	}
}

func TestMediaDocsRouteAndNavigation(t *testing.T) {
	for _, path := range []string{"/docs/media", "/docs"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
		if res.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, res.Code)
		}
		if !strings.Contains(res.Body.String(), "/docs/media") {
			t.Errorf("%s missing media nav", path)
		}
	}
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/llms-ux.txt", nil))
	for _, id := range []string{"MEDIA-IMAGE", "MEDIA-RESPONSIVE", "MEDIA-AUDIO", "MEDIA-CAPTIONS", "MEDIA-EMBED", "MEDIA-STATES"} {
		if !strings.Contains(res.Body.String(), id) {
			t.Errorf("missing %s", id)
		}
	}
	docsRes := httptest.NewRecorder()
	New().ServeHTTP(docsRes, httptest.NewRequest(http.MethodGet, "/docs/media", nil))
	for _, url := range []string{"https://www.w3.org/WAI/WCAG22/Understanding/non-text-content.html", "https://www.w3.org/WAI/WCAG22/Understanding/captions-prerecorded.html", "https://design-system.service.gov.uk/styles/images/"} {
		if !strings.Contains(docsRes.Body.String(), url) {
			t.Errorf("source not found: %s", url)
		}
	}
}

var _ *template.Template

func TestMediaComponentDocsRouteDogfoodsAudioTranscriptEmbed(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.mediaDocs(res, httptest.NewRequest(http.MethodGet, "/components/media", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("media component docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Media</h1>`,
		`<figure class="ui-media ui-media-audio">`,
		`<audio controls preload="metadata">`,
		`type="audio/mpeg"`,
		`class="ui-media-fallback">Your browser does not support audio.`,
		`class="ui-transcript" id="transcript-launch"`,
		`aria-labelledby="transcript-launch-heading"`,
		`<h2 id="transcript-launch-heading">Transcript</h2>`,
		`class="ui-transcript-content"`,
		`class="ui-media ui-media-embed"`,
		`class="ui-embed-boundary"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("media component docs are missing %q", contract)
		}
	}
}

func TestMediaComponentDocsEmbedsNothingByDefault(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.mediaDocs(res, httptest.NewRequest(http.MethodGet, "/components/media", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("media component docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// The embed specimen is NOT on an allowlist, so it must render the consent
	// boundary — never an actual iframe to a third party.
	if strings.Contains(body, "<iframe") {
		t.Error("media component docs must not render a real iframe for an unapproved embed source")
	}
	if !strings.Contains(body, "unavailable until the source is approved") {
		t.Error("media component docs must render the consent fallback copy for unapproved embeds")
	}
	if strings.Contains(body, "autoplay") {
		t.Error("media component docs must never suggest autoplay in the specimens")
	}
}