package app

import (
	"geliumui/lib"
	"net/http"
	"strings"
)

type recipeGlobalSearchResult struct {
	Resource string
	Title    string
	Detail   string
	Href     string
}

type recipeGlobalSearchView struct {
	AssetsVersion string
	Meta          metaView
	ThemeClass    string
	Query         string
	Results       []recipeGlobalSearchResult
}

func (s *server) recipeGlobalSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	needle := strings.ToLower(query)
	results := make([]recipeGlobalSearchResult, 0)
	if needle != "" {
		for _, project := range resourceDemoStore.snapshot() {
			if containsAny(needle, project.ID, project.Name, project.Owner, project.Status) {
				results = append(results, recipeGlobalSearchResult{Resource: "Admin Resource", Title: project.Name, Detail: project.Status + " · " + project.Owner, Href: "/recipes/admin-resource/" + project.ID})
			}
		}
		for _, item := range queueDemoStore.snapshot() {
			if containsAny(needle, item.ID, item.Subject, item.Requester, item.Kind, item.Status) {
				results = append(results, recipeGlobalSearchResult{Resource: "Ops Queue", Title: item.Subject, Detail: item.Status + " · " + item.Requester, Href: "/recipes/ops-queue/" + item.ID})
			}
		}
	}
	view := recipeGlobalSearchView{AssetsVersion: lib.AssetsVersion, Meta: recipeAdminMeta("Global search · Gelium UI", "Search across the Admin Resource and Ops Queue demos.", "/recipes/search"), ThemeClass: themeClass(""), Query: query, Results: results}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-global-search", view)
}

func containsAny(needle string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), needle) {
			return true
		}
	}
	return false
}
