package app

import (
	"encoding/json"
	"net/http"
	"strings"
)

// componentJSONEntry is one machine-readable component record served by
// GET /components.json. It is derived from the same registries the docs site
// renders — componentRoutes() for identity and docsSections for category — so
// the JSON can never drift from the human-facing docs.
type componentJSONEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	DocURL      string `json:"doc_url"`
	Description string `json:"description"`
}

// componentRegistryJSON is the envelope served at /components.json.
type componentRegistryJSON struct {
	Name       string               `json:"name"`
	Version    string               `json:"version"`
	Count      int                  `json:"count"`
	Components []componentJSONEntry `json:"components"`
}

// componentCategoryByPath maps each registered component route to its docs
// IA category (the docsSections title), so the JSON category matches what a
// human sees grouped on /docs.
func componentCategoryByPath() map[string]string {
	categories := make(map[string]string)
	for _, section := range docsSections {
		for _, link := range section.Links {
			categories[link.Path] = section.Title
		}
	}
	return categories
}

// buildComponentRegistryJSON derives the runtime registry from the single
// route registry. Components are emitted in navigation/registration order;
// unregistered paths are skipped rather than fabricated with an empty
// category.
func buildComponentRegistryJSON(version string) componentRegistryJSON {
	categories := componentCategoryByPath()
	out := componentRegistryJSON{
		Name:       "Gelium UI",
		Version:    version,
		Components: make([]componentJSONEntry, 0, len(componentRoutes())),
	}
	for _, r := range componentRoutes() {
		category, ok := categories[r.Path]
		if !ok {
			continue
		}
		slug := strings.TrimPrefix(r.Path, "/components/")
		out.Components = append(out.Components, componentJSONEntry{
			ID:          slug,
			Name:        r.Label,
			Category:    category,
			DocURL:      siteBaseURL + r.Path,
			Description: componentDescription(r.Label),
		})
	}
	out.Count = len(out.Components)
	return out
}

// componentsJSON serves GET /components.json — the component registry as
// structured JSON for agent/machine consumption. Same source of truth as the
// docs index; Content-Type is application/json.
func (s *server) componentsJSON(w http.ResponseWriter, r *http.Request) {
	body, err := json.MarshalIndent(buildComponentRegistryJSON(docsShellVersion), "", "  ")
	if err != nil {
		http.Error(w, "registry unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(append(body, '\n'))
}
