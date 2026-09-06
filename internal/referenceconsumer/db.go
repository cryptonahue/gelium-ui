package referenceconsumer

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

// DB is the reference consumer's persistence boundary. Gelium packages do not
// depend on this type; the consumer owns its connection and lifecycle.
type DB struct{ *sql.DB }

func Open(ctx context.Context, dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	consumerDB := &DB{DB: db}
	if err := consumerDB.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return consumerDB, nil
}

func (db *DB) Close() error { return db.DB.Close() }

func (db *DB) CreateSchema(ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS tenants (id TEXT PRIMARY KEY, name TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS projects (
  id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(id),
  name TEXT NOT NULL, status TEXT NOT NULL, date TEXT NOT NULL, owner TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(id), project_id TEXT NOT NULL REFERENCES projects(id),
  title TEXT NOT NULL, status TEXT NOT NULL, due_date TEXT NOT NULL, assignee TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS activity (
  id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(id), project_id TEXT NOT NULL REFERENCES projects(id),
  type TEXT NOT NULL, summary TEXT NOT NULL, actor TEXT NOT NULL, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS audit_events (
  id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(id), actor TEXT NOT NULL,
  action TEXT NOT NULL, resource TEXT NOT NULL, resource_id TEXT NOT NULL, created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS projects_tenant_idx ON projects(tenant_id);
CREATE INDEX IF NOT EXISTS tasks_project_tenant_idx ON tasks(project_id, tenant_id);
CREATE INDEX IF NOT EXISTS activity_project_tenant_idx ON activity(project_id, tenant_id);
`)
	return err
}

func (db *DB) Seed(ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
INSERT OR IGNORE INTO tenants(id, name) VALUES ('tenant-a', 'Tenant A'), ('tenant-b', 'Tenant B');
INSERT OR IGNORE INTO projects(id, tenant_id, name, status, date, owner) VALUES
 ('project-a', 'tenant-a', 'Tenant A Project', 'Active', '2026-01-10', 'Alicia'),
 ('project-b', 'tenant-b', 'Tenant B Project', 'Active', '2026-01-10', 'Bob');
INSERT OR IGNORE INTO tasks(id, tenant_id, project_id, title, status, due_date, assignee) VALUES
 ('task-a', 'tenant-a', 'project-a', 'Tenant A task', 'Pending', '2026-01-20', 'Alicia'),
 ('task-b', 'tenant-b', 'project-b', 'Tenant B task', 'Pending', '2026-01-20', 'Bob');
INSERT OR IGNORE INTO activity(id, tenant_id, project_id, type, summary, actor, created_at) VALUES
 ('activity-a', 'tenant-a', 'project-a', 'status', 'Project created', 'Alicia', '2026-01-01 09:00'),
 ('activity-b', 'tenant-b', 'project-b', 'status', 'Project created', 'Bob', '2026-01-01 09:00');
`)
	return err
}
