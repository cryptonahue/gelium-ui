package app

import (
	"fmt"
	"geliumui/lib"
	"net/http"
	"strings"
	"sync"
)

const recipeTaskCreateAction = "recipes.admin-resource.tasks.create"

type recipeTask struct {
	ID        string
	ProjectID string
	Title     string
	Status    string
	DueDate   string
	Assignee  string
}

type recipeTaskStore struct {
	mu    sync.Mutex
	seq   int
	items []recipeTask
}

var taskDemoStore = newRecipeTaskStore()

func newRecipeTaskStore() *recipeTaskStore {
	return &recipeTaskStore{items: []recipeTask{
		{ID: "task-alpha-1", ProjectID: "alpha", Title: "Confirm release checklist", Status: "Active", DueDate: "2026-01-08", Assignee: "Alicia R."},
		{ID: "task-alpha-2", ProjectID: "alpha", Title: "Publish release notes", Status: "Pending", DueDate: "2026-01-11", Assignee: "Bob T."},
		{ID: "task-beta-1", ProjectID: "beta", Title: "Review rollout metrics", Status: "Done", DueDate: "2026-01-30", Assignee: "Carla M."},
	}}
}

func resetRecipeTaskStore() { taskDemoStore = newRecipeTaskStore() }

func (s *recipeTaskStore) forProject(projectID string) []recipeTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []recipeTask
	for _, task := range s.items {
		if task.ProjectID == projectID {
			out = append(out, task)
		}
	}
	return out
}

func (s *recipeTaskStore) create(projectID, title, status, dueDate, assignee string) recipeTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	task := recipeTask{ID: fmt.Sprintf("task-%s-%d", projectID, s.seq), ProjectID: projectID, Title: title, Status: status, DueDate: dueDate, Assignee: assignee}
	s.items = append(s.items, task)
	return task
}

type recipeAdminTasksView struct {
	AssetsVersion string
	Meta          metaView
	ThemeClass    string
	Project       recipeResource
	Tasks         []recipeTask
	Status        string
	Error         string
}

func (s *server) recipeAdminResourceTasks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	project, ok := resourceDemoStore.get(projectID)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project whose tasks you are trying to view does not exist.")
		return
	}
	view := recipeAdminTasksView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta(project.Name+" tasks · Admin Resource recipe", "Nested project tasks in the Admin Resource recipe.", "/recipes/admin-resource/"+projectID+"/tasks"),
		ThemeClass:    themeClass(""), Project: project, Tasks: taskDemoStore.forProject(projectID),
		Status: "",
	}
	if filter := strings.TrimSpace(r.URL.Query().Get("status")); filter != "" {
		view.Status = filter
		filtered := view.Tasks[:0]
		for _, task := range view.Tasks {
			if task.Status == filter {
				filtered = append(filtered, task)
			}
		}
		view.Tasks = filtered
	}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-tasks", view)
}

func (s *server) recipeAdminResourceTaskCreate(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	project, ok := resourceDemoStore.get(projectID)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project for this task does not exist.")
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	title, status, dueDate, assignee := strings.TrimSpace(r.FormValue("title")), strings.TrimSpace(r.FormValue("status")), strings.TrimSpace(r.FormValue("due_date")), strings.TrimSpace(r.FormValue("assignee"))
	if title == "" || !recipeTaskStatus(status) || dueDate == "" {
		view := recipeAdminTasksView{AssetsVersion: lib.AssetsVersion, Meta: recipeAdminMeta(project.Name+" tasks · Admin Resource recipe", "Nested project tasks in the Admin Resource recipe.", r.URL.Path), ThemeClass: themeClass(""), Project: project, Tasks: taskDemoStore.forProject(projectID), Status: status, Error: "Enter a title, choose a valid status, and provide a due date."}
		applyRequestChrome(r, &view)
		s.renderRecipeTemplate(w, http.StatusUnprocessableEntity, "recipe-admin-resource-tasks", view)
		return
	}
	taskDemoStore.create(projectID, title, status, dueDate, assignee)
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Task created", fmt.Sprintf("%q was added to %s.", title, project.Name)))
	http.Redirect(w, r, "/recipes/admin-resource/"+projectID+"/tasks", http.StatusSeeOther)
}

func recipeTaskStatus(status string) bool {
	return status == "Active" || status == "Pending" || status == "Done"
}
