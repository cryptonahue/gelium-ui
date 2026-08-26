package app

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestComponentsJSONServesRegistry proves GET /components.json serves the
// component registry as structured JSON derived from the same registries the
// docs render — one entry per registered component with a docs category.
func TestComponentsJSONServesRegistry(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest("GET", "/components.json", nil))
	if res.Code != 200 {
		t.Fatalf("GET /components.json = %d, want 200", res.Code)
	}
	if ct := res.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body componentRegistryJSON
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("/components.json is not valid JSON: %v", err)
	}
	if body.Count != len(body.Components) {
		t.Errorf("count = %d, want %d (len of components)", body.Count, len(body.Components))
	}
	if body.Version == "" {
		t.Error("version must be set (docs shell version)")
	}
	// Every route that carries a docs IA category must appear exactly once.
	want := 0
	categories := componentCategoryByPath()
	for _, r := range componentRoutes() {
		if categories[r.Path] != "" {
			want++
		}
	}
	if got := len(body.Components); got != want {
		t.Errorf("components = %d, want %d (one per route with a docs category)", got, want)
	}
	seen := map[string]bool{}
	for _, c := range body.Components {
		if c.ID == "" || c.Name == "" {
			t.Errorf("component entry %#v has empty id or name", c)
		}
		if seen[c.ID] {
			t.Errorf("duplicate component id %q", c.ID)
		}
		seen[c.ID] = true
		if want := siteBaseURL + "/components/" + c.ID; c.DocURL != want {
			t.Errorf("%s doc_url = %q, want %q", c.ID, c.DocURL, want)
		}
		if categories["/components/"+c.ID] == "" {
			t.Errorf("%s has no docs category but was served in JSON", c.ID)
		}
		if c.Category == "" {
			t.Errorf("%s has empty category", c.ID)
		}
	}
	// Spot-check identity against the single route registry (no drift).
	for _, r := range componentRoutes() {
		slug := strings.TrimPrefix(r.Path, "/components/")
		if categories[r.Path] == "" {
			continue // recipe/primitive routes without a docs IA category
		}
		if !seen[slug] {
			t.Errorf("registry JSON is missing %q (route %s)", slug, r.Path)
		}
	}
}
