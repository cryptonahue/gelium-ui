package app

import (
	"geliumui/lib"
	"net/http"
	"strings"
)

type recipeActivity struct {
	ID        string
	ProjectID string
	Type      string
	Summary   string
	Actor     string
	CreatedAt string
}

const recipeActivityReadAction = "recipes.admin-resource.activity.read"

var recipeActivityTypes = []string{"status", "comment", "system"}

var activityDemoStore = []recipeActivity{
	{ID: "activity-alpha-1", ProjectID: "alpha", Type: "status", Summary: "Project moved to Active", Actor: "Alicia R.", CreatedAt: "2026-01-03 09:20"},
	{ID: "activity-alpha-2", ProjectID: "alpha", Type: "comment", Summary: "Release checklist reviewed", Actor: "Bob T.", CreatedAt: "2026-01-04 14:10"},
	{ID: "activity-beta-1", ProjectID: "beta", Type: "system", Summary: "Rollout window scheduled", Actor: "System", CreatedAt: "2026-01-22 08:00"},
}

type recipeAdminActivityView struct {
	AssetsVersion string
	Meta          metaView
	ThemeClass    string
	Project       recipeResource
	Activities    []recipeActivity
	Type          string
}

func (s *server) recipeAdminResourceActivity(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	project, ok := resourceDemoStore.get(projectID)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project activity you are trying to view does not exist.")
		return
	}
	if s.recipeAdminAuthorize != nil && !s.recipeAdminAuthorize(r, recipeActivityReadAction, &project) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	activityType := strings.TrimSpace(r.URL.Query().Get("type"))
	if activityType != "" && !recipeActivityTypeAllowed(activityType) {
		activityType = ""
	}
	activities := make([]recipeActivity, 0)
	for _, activity := range activityDemoStore {
		if activity.ProjectID == projectID && (activityType == "" || activity.Type == activityType) {
			activities = append(activities, activity)
		}
	}
	view := recipeAdminActivityView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta(project.Name+" activity · Admin Resource recipe", "Read-only project activity in the Admin Resource recipe.", "/recipes/admin-resource/"+projectID+"/activity"),
		ThemeClass:    themeClass(""), Project: project, Activities: activities, Type: activityType,
	}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-activity", view)
}

func recipeActivityTypeAllowed(value string) bool {
	for _, allowed := range recipeActivityTypes {
		if value == allowed {
			return true
		}
	}
	return false
}
