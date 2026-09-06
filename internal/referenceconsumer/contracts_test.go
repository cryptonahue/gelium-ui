package referenceconsumer

import (
	"context"
	"testing"
	"time"
)

func TestReferenceConsumerActivityPolicyAndAuditContracts(t *testing.T) {
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

	activities, err := NewActivityRepository(db).ForProject(ctx, "tenant-a", "project-a", "")
	if err != nil || len(activities) != 1 || activities[0].TenantID != "tenant-a" {
		t.Fatalf("activities = %+v, err = %v", activities, err)
	}
	foreignActivities, err := NewActivityRepository(db).ForProject(ctx, "tenant-a", "project-b", "")
	if err != nil || len(foreignActivities) != 0 {
		t.Fatalf("cross-tenant activities = %+v, err = %v", foreignActivities, err)
	}
	policy := SQLitePolicy{}
	if !policy.Allowed(ctx, "tenant-a", "actor-a", ActionReadTasks, "project-a") {
		t.Fatal("valid consumer policy context should be allowed")
	}
	if policy.Allowed(ctx, "tenant-a", "", ActionReadTasks, "project-a") {
		t.Fatal("missing actor should be denied")
	}

	sink := NewSQLiteAuditSink(db)
	if err := sink.Record(ctx, AuditEvent{ID: "audit-1", TenantID: "tenant-a", ActorID: "actor-a", Action: string(ActionCreateTask), Resource: "task", ResourceID: "task-a", CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE tenant_id = 'tenant-a' AND id = 'audit-1'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("audit count = %d, want 1", count)
	}
}
