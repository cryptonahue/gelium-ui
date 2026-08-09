package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListDocsRouteDogfoodsListSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/list", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("list docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>List</h1>`,
		`href="/components/list"`,
		`aria-label="List examples"`,
		`class="ui-list"`,
		`<ul`,
		`<li`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("list docs are missing %q", contract)
		}
	}
}

func TestListDocsRouteDistinguishesNavSelectAndStatic(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/list", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("list docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Navigation items must use real <a href> semantics.
	if !strings.Contains(body, `href="`) {
		t.Error("list docs must include navigation items using real <a href>")
	}
	// Selection items must use native checkboxes (no JS) for multi-select.
	if !strings.Contains(body, `type="checkbox"`) {
		t.Error("list docs must include selection items with native checkboxes")
	}
	// Static content must be plain <li> items.
	if !strings.Contains(body, `<ol`) {
		t.Error("list docs must include an ordered list of static content")
	}
	// The no-JS selection demo must be wrapped in a real <form> so checkbox
	// state submission works without JavaScript.
	if !strings.Contains(body, `<form`) {
		t.Error("list docs selection demo must be wrapped in a real <form>")
	}
}

func TestListDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/list", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST list status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestListDocsDemoRendersAllVariantsAndContentTypes(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/list", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("list docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Line variants: one-line, two-line, three-line heights.
	for _, contract := range []string{
		`class="ui-list-item"`,
		`class="ui-list-item ui-list-item--two-line"`,
		`class="ui-list-item ui-list-item--three-line"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("list demo is missing %q", contract)
		}
	}

	// Content types from the roadmap: static (<ol>), navigation (<a href>),
	// selection (native checkbox inside a form), icon list.
	for _, contract := range []string{
		`<ol class="ui-list">`,
		`<nav class="list-demo-nav"`,
		`<a class="ui-list-item-link" href="`,
		`<form method="get" action="/components/list"`,
		`<input type="checkbox" name="selection"`,
		`class="ui-list-item-icon"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("list demo is missing %q", contract)
		}
	}
}
