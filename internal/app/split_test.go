package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSplitDocsRouteDogfoodsTwoColumnComposition(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.splitDocs(res, httptest.NewRequest(http.MethodGet, "/components/split", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("split docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Split</h1>`,
		`<section class="ui-split">`,
		`class="ui-split-media"`,
		`class="ui-split-body"`,
		`class="ui-split-eyebrow">Field notes</p>`,
		`class="ui-split-title">A slower way to read</h2>`,
		`class="ui-split-copy">Long-form pages read better`,
		`class="ui-split-action">`,
		`class="ui-button ui-button-outline"`,
		`class="ui-button ui-button-primary"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("split docs are missing %q", contract)
		}
	}
}

func TestSplitDocsShowBodyOnlyLayout(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.splitDocs(res, httptest.NewRequest(http.MethodGet, "/components/split", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("split docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// Both specimens render, and the media slot is genuinely optional: the
	// body-only section has no media div.
	if !strings.Contains(body, `class="ui-split-body">`+"\n    <p class=\"ui-split-eyebrow\">Announcement</p>") {
		t.Error("split docs are missing the body-only specimen (media-less layout)")
	}
	if got := strings.Count(body, `class="ui-split-media"`); got != 1 {
		t.Errorf("split docs render %d media slots, want exactly 1 (media is optional)", got)
	}
}