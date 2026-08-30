package gates

import (
	"testing"
	"time"
)

func TestClassifyDetectorFindingsPreservesExceptionsAndSharedScope(t *testing.T) {
	result := ClassifyDetectorFindings([]DetectorFinding{
		{ID: "shell", Rule: "shell-contract", Path: "templates/feed.html"},
		{ID: "logout", Rule: "form-contract", Path: "templates/logout.html"},
	}, Scope{OwnedPaths: []string{"templates/feed.html"}, SharedPaths: []string{"templates/logout.html"}}, []DetectorException{
		{ID: "legacy-shell", FindingID: "shell", Rule: "shell-contract", Path: "templates/feed.html", Reason: "migration", Risk: "visual", Owner: "maintainer", Evidence: "ledger.json", ExpiresAt: "2026-09-01T00:00:00Z"},
		{ID: "shared-logout", FindingID: "logout", Rule: "form-contract", Path: "templates/logout.html", Reason: "shared owner", Risk: "validation", Owner: "maintainer", Evidence: "ledger.json", ExpiresAt: "2026-09-01T00:00:00Z"},
	}, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if result.Status != "pass-with-exceptions" {
		t.Fatalf("want visible exception result, got %#v", result)
	}
	if len(result.Findings) != 2 || result.Findings[0].ExceptionID != "legacy-shell" || result.Findings[1].Attribution != "shared" {
		t.Fatalf("findings must remain visible and attributed: %#v", result.Findings)
	}

	invalid := ClassifyDetectorFindings([]DetectorFinding{{ID: "shell", Rule: "shell-contract", Path: "templates/feed.html"}}, Scope{OwnedPaths: []string{"templates/feed.html"}}, []DetectorException{{
		ID: "missing-owner", FindingID: "shell", Rule: "shell-contract", Path: "templates/feed.html", Reason: "migration", Risk: "visual", Evidence: "ledger.json", ExpiresAt: "2026-09-01T00:00:00Z",
	}}, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if invalid.Status != "invalid-configuration" {
		t.Fatalf("incomplete exception must be invalid configuration: %#v", invalid)
	}

	expired := ClassifyDetectorFindings([]DetectorFinding{{ID: "shell", Rule: "shell-contract", Path: "templates/feed.html"}}, Scope{OwnedPaths: []string{"templates/feed.html"}}, []DetectorException{{
		ID: "expired", FindingID: "shell", Rule: "shell-contract", Path: "templates/feed.html", Reason: "migration", Risk: "visual", Owner: "maintainer", Evidence: "ledger.json", ExpiresAt: "2026-08-30T00:00:00Z",
	}}, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if expired.Status != "failed" || expired.Findings[0].ExceptionID != "" {
		t.Fatalf("expired exception must remain unresolved: %#v", expired)
	}
}

func TestClassifyDetectorFindingsReportsSharedBoundaryWithoutException(t *testing.T) {
	result := ClassifyDetectorFindings([]DetectorFinding{{ID: "logout", Rule: "form-contract", Path: "templates/logout.html"}}, Scope{SharedPaths: []string{"templates/logout.html"}}, nil, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if result.Status != "pass-with-exceptions" || len(result.Findings) != 1 || result.Findings[0].Attribution != "shared" {
		t.Fatalf("declared shared boundary must stay visible without failing owned scope: %#v", result)
	}
}
