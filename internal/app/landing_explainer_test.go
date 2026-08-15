package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLandingLinksToGitHubSource proves the landing surface links the open
// source repository (task: add GitHub repo link to landing).
func TestLandingLinksToGitHubSource(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `href="https://github.com/cryptonahue/gelium-ui"`) {
		t.Error("landing must link https://github.com/cryptonahue/gelium-ui")
	}
	if !strings.Contains(body, "View source") {
		t.Error("landing GitHub link must be labelled View source")
	}
}

// TestLandingGitHubLinkStaysClean proves the external source link is never
// chrome-rewritten: ?theme= / ?scheme= must not leak onto the GitHub href.
func TestLandingGitHubLinkStaysClean(t *testing.T) {
	body := getOKBody(t, "/?theme=basecoat&scheme=dark")
	if !strings.Contains(body, `href="https://github.com/cryptonahue/gelium-ui"`) {
		t.Fatal("landing GitHub link must survive ?theme=basecoat&scheme=dark")
	}
	if strings.Contains(body, "github.com/cryptonahue/gelium-ui?") {
		t.Error("external GitHub href must not carry the chrome query string")
	}
}

// TestLandingCodeSampleDocumentsBaseURL proves the landing code sample tells
// deployers about the BASE_URL environment variable (canonical origin).
func TestLandingCodeSampleDocumentsBaseURL(t *testing.T) {
	body := getOKBody(t, "/")
	if !strings.Contains(body, "BASE_URL=") {
		t.Error("landing code sample must document BASE_URL for deployers")
	}
}
