package app

import (
	"fmt"
	"geliumui/lib"
	"net/http"
	"sort"
)

type recipeAdminDashboardMetric struct {
	Label string
	Value string
}

type recipeAdminDashboardView struct {
	AssetsVersion string
	Meta          metaView
	ThemeClass    string
	DataTheme     string
	State         string
	Metrics       []recipeAdminDashboardMetric
	HasData       bool
	ListHref      string
	Recent        []recipeResource
}

func (s *server) recipeAdminDashboard(w http.ResponseWriter, r *http.Request) {
	items := resourceDemoStore.snapshot()
	recent := append([]recipeResource(nil), items...)
	sort.SliceStable(recent, func(i, j int) bool { return recent[i].Date > recent[j].Date })
	if len(recent) > 5 {
		recent = recent[:5]
	}
	state := r.URL.Query().Get("state")
	if state != "ready" && state != "empty" && state != "loading" && state != "error" {
		state = "ready"
	}
	active, pending, done := 0, 0, 0
	for _, item := range items {
		switch item.Status {
		case "Active":
			active++
		case "Pending":
			pending++
		case "Done":
			done++
		}
	}
	view := recipeAdminDashboardView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta("Projects dashboard", "Server-rendered project metrics in the Admin Resource dashboard recipe.", "/recipes/admin-dashboard"),
		ThemeClass:    themeClass(""),
		State:         state,
		Metrics: []recipeAdminDashboardMetric{
			{Label: "Total projects", Value: fmt.Sprintf("%d", len(items))},
			{Label: "Active projects", Value: fmt.Sprintf("%d", active)},
			{Label: "Pending projects", Value: fmt.Sprintf("%d", pending)},
			{Label: "Done projects", Value: fmt.Sprintf("%d", done)},
		},
		HasData:  len(items) > 0,
		ListHref: "/recipes/admin-resource",
		Recent:   recent,
	}
	if state == "empty" {
		view.HasData = false
	}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-dashboard", view)
}
