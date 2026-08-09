package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCardDocsRouteDogfoodsSemanticRootNodes(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/card", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("card docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Card</h1>`,
		`href="/components/card"`,
		`aria-label="Card examples"`,
		`<article class="ui-card ui-card-elevated"`,
		`<a class="ui-card ui-card-outlined"`,
		`<button class="ui-card ui-card-filled"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("card docs are missing %q", contract)
		}
	}
}

func TestCardDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/card", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST card status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}
