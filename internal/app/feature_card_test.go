package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeatureCardDocsRouteDogfoodsCardComposition(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.featureCardDocs(res, httptest.NewRequest(http.MethodGet, "/components/feature-card", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("feature card docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Feature card</h1>`,
		`class="ui-card ui-card-elevated ui-feature-card"`,
		`class="ui-feature-card-media"`,
		`class="ui-feature-card-body"`,
		`class="ui-card-title">Plan a weekend escape</h3>`,
		`class="ui-card-body">Three quiet trails`,
		`class="ui-card-action">`,
		`class="ui-button ui-button-primary"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feature card docs are missing %q", contract)
		}
	}
}

func TestFeatureCardDocsRendersOneRealAction(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.featureCardDocs(res, httptest.NewRequest(http.MethodGet, "/components/feature-card", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("feature card docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if got := strings.Count(body, `<a class="ui-button ui-button-primary"`); got != 1 {
		t.Errorf("feature card specimen renders %d action links, want exactly 1", got)
	}
	if strings.Contains(body, "</a></a>") {
		t.Error("feature card specimen must not nest links")
	}
}