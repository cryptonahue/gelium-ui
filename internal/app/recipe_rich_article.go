package app

import (
	"encoding/json"
	"geliumui/lib"
	"html/template"
	"net/http"
)

// richArticleView is a fixed, server-rendered integration fixture. It deliberately
// composes existing media, list, alert, and data-table contracts without adding
// generic frontmatter or a new runtime abstraction.
type richArticleView struct {
	Meta          metaView
	ThemeClass    string
	DataTheme     string
	AssetsVersion string
	Title         string
	Description   string
}

func richArticleJSONLD() template.JS {
	value := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      "Designing a resilient content surface",
		"description":   "A Gelium UI integration fixture for rich, accessible, server-rendered content.",
		"url":           siteBaseURL + "/recipes/rich-article",
		"inLanguage":    "en",
		"datePublished": "2026-08-17",
		"author":        map[string]string{"@type": "Person", "name": "Gelium UI team"},
		"publisher":     map[string]string{"@type": "Organization", "name": "Gelium UI"},
	}
	b, err := json.Marshal(value)
	if err != nil {
		return template.JS("")
	}
	return template.JS(b) // #nosec G203 -- fixed, system-generated JSON.
}

func (s *server) recipeRichArticle(w http.ResponseWriter, r *http.Request) {
	view := &richArticleView{
		AssetsVersion: lib.AssetsVersion,
		Title:         "Designing a resilient content surface",
		Description:   "A complete Gelium UI article fixture: semantic content, native media, data, and recoverable states.",
	}
	view.Meta = metaView{
		Title: view.Title, Description: view.Description,
		Canonical: siteBaseURL + "/recipes/rich-article", Robots: "noindex, nofollow", OGType: "article", OGTitle: view.Title + " · Gelium UI", OGDescription: view.Description, TwitterURL: siteBaseURL + "/recipes/rich-article", JSONLD: richArticleJSONLD(), Lang: "en",
	}
	applyRequestChrome(r, view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-rich-article", view)
}

// Keep the library import visible to the generated template registry boundary.
var _ = lib.LibAssets
