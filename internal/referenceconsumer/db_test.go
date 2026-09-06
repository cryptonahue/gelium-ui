package referenceconsumer

import (
	"context"
	"testing"
)

func TestSQLiteReferenceConsumerScopesRepositoriesByTenant(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(ctx); err != nil {
		t.Fatal(err)
	}
	if err := db.Seed(ctx); err != nil {
		t.Fatal(err)
	}

	projects := NewProjectRepository(db)
	project, err := projects.Get(ctx, "tenant-a", "project-a")
	if err != nil || project.Name != "Tenant A Project" {
		t.Fatalf("project = %+v, err = %v", project, err)
	}
	if _, err := projects.Get(ctx, "tenant-a", "project-b"); !IsNotFound(err) {
		t.Fatalf("cross-tenant project error = %v, want not found", err)
	}

	tasks := NewTaskRepository(db)
	rows, err := tasks.ForProject(ctx, "tenant-a", "project-a")
	if err != nil || len(rows) != 1 || rows[0].TenantID != "tenant-a" {
		t.Fatalf("tasks = %+v, err = %v", rows, err)
	}
	rows, err = tasks.ForProject(ctx, "tenant-a", "project-b")
	if err != nil || len(rows) != 0 {
		t.Fatalf("cross-tenant tasks = %+v, err = %v", rows, err)
	}
}
