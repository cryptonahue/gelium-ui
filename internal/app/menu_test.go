package app

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestMenuDocsRouteDogfoodsPopoverAndAnatomy(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("menu docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Menu</h1>`,
		`href="/components/menu"`,
		`aria-label="Menu examples"`,
		`class="ui-menu"`,
		`popover`,
		`popovertarget=`,
		`popovertargetaction=`,
		`aria-expanded="false"`,
		`<ul`,
		`<li`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("menu docs are missing %q", contract)
		}
	}
}

// TestMenuPopoverTriggersDeclareExpandedState closes gap G8: every popover
// trigger must expose aria-expanded so assistive tech knows the surface is
// currently collapsed. Each popover trigger carries the attribute.
func TestMenuPopoverTriggersDeclareExpandedState(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("menu docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// Scope to the popover triggers themselves: the docs chrome now also emits
	// aria-expanded (search input), so a body-wide count would over-count.
	// Each trigger must pair its own aria-haspopup with aria-expanded (G8).
	triggers := regexp.MustCompile(`class="ui-button ui-button-secondary menu-demo-trigger"[^>]*popovertarget="[^"]+"[^>]*>`).FindAllString(body, -1)
	if len(triggers) != 3 {
		t.Fatalf("menu popover triggers = %d, want 3 (Actions, Navigate, Select)", len(triggers))
	}
	for i, tag := range triggers {
		if !strings.Contains(tag, `aria-haspopup="menu"`) || !strings.Contains(tag, `aria-expanded="false"`) {
			t.Errorf("trigger %d must pair aria-haspopup with aria-expanded=false (G8): %s", i, tag)
		}
	}
	if strings.Count(body, `aria-haspopup="menu"`) != len(triggers) {
		t.Error("every aria-haspopup trigger must pair with aria-expanded (G8)")
	}
}

func TestMenuDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/menu", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST menu status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

// TestMenuDocsActionsUseNativeButtons guards the action-item contract: menu
// items that perform an action are real <button type="button"> elements, and a
// disabled item uses the native disabled attribute so the platform removes the
// activation path without any component JavaScript.
func TestMenuDocsActionsUseNativeButtons(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<button type="button" class="ui-menu-item-button">`,
		`<button type="button" class="ui-menu-item-button" disabled`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("action menu is missing %q", contract)
		}
	}
	for _, forbidden := range []string{
		`onclick=`,
		`onchange=`,
		`onkeydown=`,
	} {
		if strings.Contains(body, forbidden) {
			t.Errorf("menu docs must not rely on %q", forbidden)
		}
	}
}

// TestMenuDocsNavigationUsesRealLinks guards the roadmap's "fallback no-JS:
// navegación o formulario real": navigation items must be real <a href> links
// to existing Gelium routes, never inaccessible CSS imitation.
func TestMenuDocsNavigationUsesRealLinks(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<a class="ui-menu-item-link" href="`,
		`href="/components/button"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation menu is missing %q", contract)
		}
	}
}

// TestMenuDocsSelectionUsesNativeCheckboxRadio guards the selection contract:
// checkbox/radio items are native controls inside a real <form>, so the checked
// values submit through a normal GET round-trip with no component JavaScript.
func TestMenuDocsSelectionUsesNativeCheckboxRadio(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<form method="get" action="/components/menu"`,
		`<input type="checkbox" name="menu-multi"`,
		`<input type="radio" name="menu-single"`,
		`type="submit"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("selection menu is missing %q", contract)
		}
	}
}

// TestMenuDocsDeclarativeTopLayerNoJS guards the platform-first decision: the
// menu surface opens through the native popover attribute and the trigger uses
// popovertarget/popovertargetaction — no component JavaScript anywhere.
func TestMenuDocsDeclarativeTopLayerNoJS(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<button type="button" class="ui-button`,
		`popovertarget="menu-`,
		`popovertargetaction="toggle"`,
		`aria-haspopup="menu"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("menu docs are missing declarative trigger %q", contract)
		}
	}
	if strings.Contains(body, "<script>") {
		t.Error("menu docs must not include inline component JavaScript")
	}
}

// TestMenuDocsRendersDividerAndIcons guards the Material menu anatomy: a
// divider separates grouped items and leading icons are decorative and hidden
// from assistive technology.
func TestMenuDocsRendersDividerAndIcons(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/menu", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`role="separator"`,
		`class="ui-menu-item-icon"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("menu anatomy is missing %q", contract)
		}
	}
}
