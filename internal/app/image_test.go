package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestImageDocsRouteDogfoodsMediaFigureContract(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.imageDocs(res, httptest.NewRequest(http.MethodGet, "/components/image", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("image docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Image</h1>`,
		`<figure class="ui-media ui-media-image">`,
		`class="ui-media ui-media-image ui-media--aspect"`,
		`style="--ui-media-aspect:16 / 9"`,
		`class="ui-media ui-media-picture"`,
		`loading="lazy"`,
		`srcset=`,
		`sizes=`,
		`type="image/webp"`,
		`<figcaption>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("image docs are missing %q", contract)
		}
	}
}

func TestImageDocsRendersAltAndIntrinsicDimensions(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.imageDocs(res, httptest.NewRequest(http.MethodGet, "/components/image", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("image docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, img := range []string{
		`alt="A mountain lake at dusk" width="1200" height="800"`,
		`alt="A mountain lake at dusk" width="1200" height="800" loading="lazy" srcset=`,
	} {
		if !strings.Contains(body, img) {
			t.Errorf("image docs render an img missing %q", img)
		}
	}
	if got := strings.Count(body, "<img "); got != 3 {
		t.Errorf("image docs render %d img elements, want 3 (single, responsive, picture)", got)
	}
	if !strings.Contains(body, "decorative") {
		t.Error("image docs must teach the empty-alt rule for decorative images")
	}
}