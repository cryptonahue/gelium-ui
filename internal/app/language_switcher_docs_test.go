package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLanguageSwitcherDocsRouteDogfoodsGetFormContract(t *testing.T) {
	s := docsTestServer(t)
	res := httptest.NewRecorder()
	s.languageSwitcherDocs(res, httptest.NewRequest(http.MethodGet, "/components/language-switcher", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("language switcher docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`>Language switcher</h1>`,
		`<form class="ui-language-switcher" method="get"`,
		`class="ui-language-switcher-label"`,
		`class="ui-language-switcher-control"`,
		`class="ui-language-switcher-select"`,
		`name="lang"`,
		`type="submit"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("language switcher docs are missing %q", contract)
		}
	}
}

func TestLanguageSwitcherDocsRegisterRouteAndNav(t *testing.T) {
	// The route registry and the docs IA must both carry the page so the
	// sidebar, /docs index, and /components.json stay in sync.
	foundRoute := false
	for _, r := range componentRoutes() {
		if r.Path == "/components/language-switcher" && r.Label == "Language switcher" {
			foundRoute = true
		}
	}
	if !foundRoute {
		t.Error("componentRoutes() is missing /components/language-switcher")
	}

	categories := componentCategoryByPath()
	if got := categories["/components/language-switcher"]; got == "" {
		t.Error("language switcher has no docs IA category (docsSections)")
	}
}
