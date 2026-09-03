package gates

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateLedgerV1Fixtures(t *testing.T) {
	now := time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name     string
		fixture  string
		wantCode string
	}{
		{name: "valid", fixture: "valid-v1.json"},
		{name: "unsupported schema", fixture: "unsupported-schema.json", wantCode: "unsupported-schema"},
		{name: "missing reading", fixture: "missing-reading.json", wantCode: "missing-reading-attestation"},
		{name: "unknown status", fixture: "unknown-status.json", wantCode: "unknown-gate-status"},
		{name: "ownership overlap", fixture: "ownership-overlap.json", wantCode: "invalid-ownership-boundary"},
		{name: "expired exception", fixture: "expired-exception.json", wantCode: "expired-exception"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			ledger, issues := ValidateLedger(data, now)
			if tc.wantCode == "" {
				if len(issues) != 0 {
					t.Fatalf("valid ledger has issues: %#v", issues)
				}
				if ledger.SchemaVersion != 1 || ledger.Route != RouteDesignGated {
					t.Fatalf("parsed unexpected ledger: %#v", ledger)
				}
				return
			}
			if !hasIssue(issues, tc.wantCode) {
				t.Fatalf("want issue %q, got %#v", tc.wantCode, issues)
			}
		})
	}
}

func TestValidateLedgerRejectsMalformedJSON(t *testing.T) {
	_, issues := ValidateLedger([]byte(`{"schema_version":`), time.Now())
	if !hasIssue(issues, "malformed-json") {
		t.Fatalf("want malformed-json issue, got %#v", issues)
	}
}

func TestValidateLedgerAcceptsDelegatedDirectRoute(t *testing.T) {
	data := []byte(`{"schema_version":1,"route":"delegated-direct","scope":{"owned_paths":["internal/app/feed.go"]}}`)
	ledger, issues := ValidateLedger(data, time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC))
	if len(issues) != 0 {
		t.Fatalf("delegated-direct ledger has issues: %#v", issues)
	}
	if ledger.Route != Route("delegated-direct") {
		t.Fatalf("parsed route = %q, want delegated-direct", ledger.Route)
	}
}

func hasIssue(issues []Issue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}
