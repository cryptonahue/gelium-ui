package referenceconsumer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReferenceConsumerHTTPScopesTenantAndAuditsTaskCreation(t *testing.T) {
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
	h := NewHandler(db, SQLitePolicy{}, NewSQLiteAuditSink(db))

	get := httptest.NewRequest(http.MethodGet, "/reference/projects/project-a/tasks", nil)
	get.Header.Set("X-Consumer-Tenant", "tenant-a")
	get.Header.Set("X-Consumer-Actor", "actor-a")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, get)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Tenant A task") {
		t.Fatalf("get = %d, body missing task", res.Code)
	}

	foreign := httptest.NewRequest(http.MethodGet, "/reference/projects/project-b/tasks", nil)
	foreign.Header.Set("X-Consumer-Tenant", "tenant-a")
	foreign.Header.Set("X-Consumer-Actor", "actor-a")
	foreignRes := httptest.NewRecorder()
	h.ServeHTTP(foreignRes, foreign)
	if foreignRes.Code != http.StatusNotFound {
		t.Fatalf("foreign get = %d, want 404", foreignRes.Code)
	}

	post := httptest.NewRequest(http.MethodPost, "/reference/projects/project-a/tasks", strings.NewReader("title=New+task&status=Pending&due_date=2026-02-01&assignee=Alicia"))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("X-Consumer-Tenant", "tenant-a")
	post.Header.Set("X-Consumer-Actor", "actor-a")
	postRes := httptest.NewRecorder()
	h.ServeHTTP(postRes, post)
	if postRes.Code != http.StatusSeeOther {
		t.Fatalf("post = %d, want 303", postRes.Code)
	}
	var auditCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE tenant_id = 'tenant-a' AND action = 'tasks.create'`).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 1 {
		t.Fatalf("audit count = %d, want 1", auditCount)
	}
}
