package referenceconsumer

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	DB       *DB
	Projects ProjectRepository
	Tasks    TaskRepository
	Policy   Policy
	Audit    AuditSink
}

func NewHandler(db *DB, policy Policy, audit AuditSink) Handler {
	return Handler{DB: db, Projects: NewProjectRepository(db), Tasks: NewTaskRepository(db), Policy: policy, Audit: audit}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/reference/projects" && r.Method == http.MethodGet {
		h.projects(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/reference/projects/") {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) == 4 && parts[3] == "tasks" {
			if r.Method == http.MethodGet {
				h.tasks(w, r, parts[2])
				return
			}
			if r.Method == http.MethodPost {
				h.createTask(w, r, parts[2])
				return
			}
		}
	}
	http.NotFound(w, r)
}

func (h Handler) projects(w http.ResponseWriter, r *http.Request) {
	// Listing is intentionally omitted from this first reference slice; a real
	// consumer would add a tenant-scoped repository query here.
	http.Error(w, "reference project listing is not implemented", http.StatusNotImplemented)
}

func (h Handler) tasks(w http.ResponseWriter, r *http.Request, projectID string) {
	project, err := h.Projects.Get(r.Context(), tenantFrom(r), projectID)
	if err != nil || !h.allowed(r, ActionReadTasks, projectID) {
		http.NotFound(w, r)
		return
	}
	tasks, err := h.Tasks.ForProject(r.Context(), tenantFrom(r), projectID)
	if err != nil {
		http.Error(w, "tasks unavailable", http.StatusInternalServerError)
		return
	}
	renderReferenceTasks(w, project, tasks, "")
}

func (h Handler) createTask(w http.ResponseWriter, r *http.Request, projectID string) {
	project, err := h.Projects.Get(r.Context(), tenantFrom(r), projectID)
	if err != nil || !h.allowed(r, ActionCreateTask, projectID) {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	task := Task{ID: fmt.Sprintf("task-%d", time.Now().UnixNano()), TenantID: tenantFrom(r), ProjectID: projectID, Title: strings.TrimSpace(r.FormValue("title")), Status: strings.TrimSpace(r.FormValue("status")), DueDate: strings.TrimSpace(r.FormValue("due_date")), Assignee: strings.TrimSpace(r.FormValue("assignee"))}
	if task.Title == "" || !validTaskStatus(task.Status) || task.DueDate == "" {
		tasks, _ := h.Tasks.ForProject(r.Context(), tenantFrom(r), projectID)
		renderReferenceTasksStatus(w, project, tasks, "Enter a title, valid status, and due date.", http.StatusUnprocessableEntity)
		return
	}
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "task unavailable", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(r.Context(), `INSERT INTO tasks(id, tenant_id, project_id, title, status, due_date, assignee) VALUES (?, ?, ?, ?, ?, ?, ?)`, task.ID, task.TenantID, task.ProjectID, task.Title, task.Status, task.DueDate, task.Assignee); err != nil {
		_ = tx.Rollback()
		http.Error(w, "task unavailable", http.StatusInternalServerError)
		return
	}
	if h.Audit != nil {
		if _, err := tx.ExecContext(r.Context(), `INSERT INTO audit_events(id, tenant_id, actor, action, resource, resource_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, fmt.Sprintf("audit-%d", time.Now().UnixNano()), tenantFrom(r), actorFrom(r), string(ActionCreateTask), "task", task.ID, time.Now().UTC().Format(time.RFC3339)); err != nil {
			_ = tx.Rollback()
			http.Error(w, "audit unavailable", http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "task unavailable", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reference/projects/"+projectID+"/tasks", http.StatusSeeOther)
}

func (h Handler) allowed(r *http.Request, action Action, projectID string) bool {
	return h.Policy == nil || h.Policy.Allowed(r.Context(), tenantFrom(r), actorFrom(r), action, projectID)
}
func tenantFrom(r *http.Request) string { return r.Header.Get("X-Consumer-Tenant") }
func actorFrom(r *http.Request) string  { return r.Header.Get("X-Consumer-Actor") }
func validTaskStatus(status string) bool {
	return status == "Active" || status == "Pending" || status == "Done"
}

type tasksPageData struct {
	Project Project
	Tasks   []Task
	Error   string
}

var tasksPage = template.Must(template.New("reference-tasks").Parse(`<!doctype html><html><head><title>{{.Project.Name}} tasks</title></head><body><main><p><a href="/reference/projects/{{.Project.ID}}/tasks">Tasks</a></p><h1>{{.Project.Name}} tasks</h1>{{if .Error}}<p role="alert">{{.Error}}</p>{{end}}<ul>{{range .Tasks}}<li>{{.Title}} — {{.Status}} — {{.DueDate}}</li>{{else}}<li>No tasks.</li>{{end}}</ul><form method="post" action="/reference/projects/{{.Project.ID}}/tasks"><label>Title <input name="title" required></label><label>Status <select name="status"><option>Pending</option><option>Active</option><option>Done</option></select></label><label>Due date <input name="due_date" type="date" required></label><label>Assignee <input name="assignee"></label><button type="submit">Add task</button></form></main></body></html>`))

func renderReferenceTasks(w http.ResponseWriter, project Project, tasks []Task, message string) {
	renderReferenceTasksStatus(w, project, tasks, message, http.StatusOK)
}
func renderReferenceTasksStatus(w http.ResponseWriter, project Project, tasks []Task, message string, status int) {
	w.WriteHeader(status)
	_ = tasksPage.Execute(w, tasksPageData{Project: project, Tasks: tasks, Error: message})
}

var _ http.Handler = Handler{}
