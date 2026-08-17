package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVideoDocsRouteDogfoodsNativePlayerContract(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.videoDocs(res, httptest.NewRequest(http.MethodGet, "/components/video", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("video docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Video</h1>`,
		`<div class="ui-video">`,
		`class="ui-video ui-video--aspect-4-3"`,
		`<video controls`,
		`loading="lazy"`,
		`crossorigin="anonymous"`,
		`<source src="/static/video/walkthrough.mp4" type="video/mp4">`,
		`<track kind="captions" src="/static/video/walkthrough.en.vtt" srclang="en" label="English">`,
		`class="ui-video-fallback">Your browser does not support HTML video.`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("video docs are missing %q", contract)
		}
	}
}

func TestVideoDocsNeverAutoplays(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.videoDocs(res, httptest.NewRequest(http.MethodGet, "/components/video", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("video docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// The attribute must not appear on any rendered video element; prose
	// mentions writing the spec rule, but never as an attribute value.
	if strings.Contains(body, " autoplay") {
		t.Error("video docs must never render the autoplay attribute on a video element")
	}
	if got := strings.Count(body, "<video "); got != 2 {
		t.Errorf("video docs render %d video elements, want 2 (16:9 and 4:3 specimens)", got)
	}
}