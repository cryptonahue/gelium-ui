package referenceconsumer

import (
	"context"
	"time"
)

type Action string

const (
	ActionReadProject  Action = "project.read"
	ActionReadTasks    Action = "tasks.read"
	ActionCreateTask   Action = "tasks.create"
	ActionReadActivity Action = "activity.read"
	ActionTransition   Action = "project.transition"
)

type Policy interface {
	Allowed(ctx context.Context, tenantID, actorID string, action Action, projectID string) bool
}

type AuditEvent struct {
	ID, TenantID, ActorID, Action, Resource, ResourceID string
	CreatedAt                                           time.Time
}

type AuditSink interface {
	Record(ctx context.Context, event AuditEvent) error
}

type SQLitePolicy struct{}

func (SQLitePolicy) Allowed(_ context.Context, tenantID, actorID string, action Action, projectID string) bool {
	return tenantID != "" && actorID != "" && projectID != "" && action != ""
}

type SQLiteAuditSink struct{ db *DB }

func NewSQLiteAuditSink(db *DB) SQLiteAuditSink { return SQLiteAuditSink{db: db} }

func (s SQLiteAuditSink) Record(ctx context.Context, event AuditEvent) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events(id, tenant_id, actor, action, resource, resource_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, event.ID, event.TenantID, event.ActorID, event.Action, event.Resource, event.ResourceID, event.CreatedAt.UTC().Format(time.RFC3339))
	return err
}
