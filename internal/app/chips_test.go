package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChipsDocsRouteDogfoodsVariants(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/chips", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("chips docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Chips</h1>`,
		`href="/components/chips"`,
		`aria-label="Chips examples"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("chips docs are missing %q", contract)
		}
	}
}

func TestChipsDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/chips", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST chips status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestChipsDocsRouteUsesPrimitiveSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/chips", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("chips docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-chip ui-chip-assist"`,
		`class="ui-chip ui-chip-suggestion"`,
		`class="ui-chip ui-chip-input"`,
		`type="checkbox"`,
		`action="/examples/chips/remove"`,
		`aria-label="Remove Star Wars"`,
		`class="ui-chip-remove"`,
		`class="ui-chip ui-chip-filter"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("chips docs are missing primitive semantics %q", contract)
		}
	}
}

func TestChipsRemoveRoundTripDropsChipWithoutJS(t *testing.T) {
	form := strings.NewReader("remove=star-wars")
	req := httptest.NewRequest(http.MethodPost, "/examples/chips/remove", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("chips remove status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `class="ui-chip-none" data-removed="star-wars"`) {
		t.Error("removed chip must be re-rendered absent (non-interactive stand-in)")
	}
	if strings.Contains(body, `aria-label="Remove Star Wars"`) {
		t.Error("the removed chip's remove button must no longer render")
	}
	if !strings.Contains(body, `role="status"`) {
		t.Error("the removal must surface a live status notice")
	}
	if !strings.Contains(body, `data-chip="star-trek"`) {
		t.Error("an unremoved chip must still render")
	}
}

func TestChipsRemoveRoundTripUnknownValueStaysIntact(t *testing.T) {
	form := strings.NewReader("remove=does-not-exist")
	req := httptest.NewRequest(http.MethodPost, "/examples/chips/remove", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("chips remove status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`aria-label="Remove Star Wars"`,
		`aria-label="Remove Star Trek"`,
		`Nothing to remove.`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("unknown removal must keep both chips and notify; missing %q", contract)
		}
	}
}
