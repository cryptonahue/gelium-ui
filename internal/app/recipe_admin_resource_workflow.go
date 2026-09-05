package app

import (
	"fmt"
	"geliumui/lib"
	"net/http"
	"strings"
)

const recipeAdminTransitionAction = "recipes.admin-resource.transition"

func (s *recipeResourceStore) transition(id, target string) (recipeResource, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		item := &s.items[i]
		if item.ID != id || !recipeStatusTransitionAllowed(item.Status, target) {
			continue
		}
		item.Status = target
		return *item, true
	}
	return recipeResource{}, false
}

func recipeStatusTransitionAllowed(current, target string) bool {
	switch {
	case current == "Pending" && target == "Active":
		return true
	case current == "Active" && target == "Done":
		return true
	default:
		return false
	}
}

func (s *server) recipeAdminResourceTransition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := resourceDemoStore.get(id)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project whose status you are trying to change does not exist.")
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	target := strings.TrimSpace(r.FormValue("status"))
	if _, changed := resourceDemoStore.transition(id, target); !changed {
		s.renderRecipeResourceDetailError(w, r, item, "That status transition is not allowed. Refresh the project and choose the next valid status.")
		return
	}
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Project status updated", fmt.Sprintf("%q moved to %s.", item.Name, target)))
	http.Redirect(w, r, "/recipes/admin-resource/"+id, http.StatusSeeOther)
}

func (s *server) renderRecipeResourceDetailError(w http.ResponseWriter, r *http.Request, item recipeResource, message string) {
	view := recipeAdminResourceDetailView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta(item.Name+" · Admin Resource recipe", "Project details in the Admin Resource screen recipe.", "/recipes/admin-resource/"+item.ID),
		ThemeClass:    themeClass(""), Item: item, StatusTone: recipeStatusTone(item.Status),
		ListHref: "/recipes/admin-resource", EditHref: "/recipes/admin-resource/" + item.ID + "/edit", DeleteHref: "/recipes/admin-resource/" + item.ID + "/delete",
		TransitionError: message,
	}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusUnprocessableEntity, "recipe-admin-resource-detail", view)
}
