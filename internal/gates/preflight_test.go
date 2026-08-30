package gates

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPrebuildEvaluatesDesignGatedLedger(t *testing.T) {
	ledger := validLedger(t)
	result := EvaluatePrebuild(ledger, nil, []string{"templates/feed.html", "templates/logout.html"})
	if result.Status != "pass" {
		t.Fatalf("approved design-gated ledger should pass: %#v", result)
	}

	ledger.Gates.Approval.Status = "changes-requested"
	result = EvaluatePrebuild(ledger, nil, nil)
	if result.Status != "failed" || !hasIssue(result.Issues, "approval-blocked") {
		t.Fatalf("changes requested must block prebuild: %#v", result)
	}
}

func TestPrebuildAcceptsDirectExemptAndRejectsUncoveredPaths(t *testing.T) {
	ledger := Ledger{
		SchemaVersion: SchemaVersionV1,
		Route:         RouteDirectExempt,
		Scope:         Scope{OwnedPaths: []string{"templates/feed.html"}},
	}
	result := EvaluatePrebuild(ledger, nil, []string{"templates/feed.html"})
	if result.Status != "pass" {
		t.Fatalf("direct-exempt ledger should pass: %#v", result)
	}
	result = EvaluatePrebuild(ledger, nil, []string{"templates/other.html"})
	if result.Status != "failed" || !hasIssue(result.Issues, "uncovered-changed-path") {
		t.Fatalf("uncovered path must fail: %#v", result)
	}
}

func TestPrebuildReportsInvalidConfiguration(t *testing.T) {
	result := EvaluatePrebuild(Ledger{}, []Issue{{Code: "malformed-json"}}, nil)
	if result.Status != "invalid-configuration" || !hasIssue(result.Issues, "malformed-json") {
		t.Fatalf("validation issue must remain invalid configuration: %#v", result)
	}
}

func validLedger(t *testing.T) Ledger {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "valid-v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	ledger, issues := ValidateLedger(data, time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC))
	if len(issues) != 0 {
		t.Fatalf("valid fixture issues: %#v", issues)
	}
	return ledger
}
