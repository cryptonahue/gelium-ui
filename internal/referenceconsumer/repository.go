package referenceconsumer

import (
	"context"
	"database/sql"
)

type Project struct {
	ID, TenantID, Name, Status, Date, Owner string
}

type Task struct {
	ID, TenantID, ProjectID, Title, Status, DueDate, Assignee string
}

type Activity struct {
	ID, TenantID, ProjectID, Type, Summary, Actor, CreatedAt string
}

type ProjectRepository struct{ db *DB }
type TaskRepository struct{ db *DB }
type ActivityRepository struct{ db *DB }

func NewProjectRepository(db *DB) ProjectRepository   { return ProjectRepository{db: db} }
func NewTaskRepository(db *DB) TaskRepository         { return TaskRepository{db: db} }
func NewActivityRepository(db *DB) ActivityRepository { return ActivityRepository{db: db} }

func (r ProjectRepository) Get(ctx context.Context, tenantID, id string) (Project, error) {
	var project Project
	err := r.db.QueryRowContext(ctx, `SELECT id, tenant_id, name, status, date, owner FROM projects WHERE tenant_id = ? AND id = ?`, tenantID, id).
		Scan(&project.ID, &project.TenantID, &project.Name, &project.Status, &project.Date, &project.Owner)
	return project, err
}

func (r TaskRepository) ForProject(ctx context.Context, tenantID, projectID string) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, tenant_id, project_id, title, status, due_date, assignee FROM tasks WHERE tenant_id = ? AND project_id = ? ORDER BY due_date, id`, tenantID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.TenantID, &task.ProjectID, &task.Title, &task.Status, &task.DueDate, &task.Assignee); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r TaskRepository) Create(ctx context.Context, tenantID string, task Task) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO tasks(id, tenant_id, project_id, title, status, due_date, assignee) VALUES (?, ?, ?, ?, ?, ?, ?)`, task.ID, tenantID, task.ProjectID, task.Title, task.Status, task.DueDate, task.Assignee)
	return err
}

func (r ActivityRepository) ForProject(ctx context.Context, tenantID, projectID, activityType string) ([]Activity, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, tenant_id, project_id, type, summary, actor, created_at FROM activity WHERE tenant_id = ? AND project_id = ? AND (? = '' OR type = ?) ORDER BY created_at DESC, id`, tenantID, projectID, activityType, activityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []Activity
	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.ID, &activity.TenantID, &activity.ProjectID, &activity.Type, &activity.Summary, &activity.Actor, &activity.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

func IsNotFound(err error) bool { return err == sql.ErrNoRows }
