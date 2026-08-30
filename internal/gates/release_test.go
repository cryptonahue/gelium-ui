package gates

import "testing"

func TestReleaseRequiresPostbuildEvidence(t *testing.T) {
	ledger := validLedger(t)
	ledger.Gates.RenderedAudit = Gate{Status: "pass", Evidence: []string{"audit.md#responsive"}}
	ledger.Release = ReleaseEvidence{
		Detector:        DetectorEvidence{Status: "clean-pass", Evidence: []string{"detector.json"}},
		Checks:          CheckEvidence{Tests: []string{"go test ./..."}, Builds: []string{"npm run build"}},
		AuthorityMatrix: StatusEvidence{Status: "pass", Evidence: []string{"authority.json"}},
	}
	result := EvaluateRelease(ledger, nil, []string{"templates/feed.html"})
	if result.Status != "pass" {
		t.Fatalf("complete release evidence should pass: %#v", result)
	}

	ledger.Release.Detector.Status = "failed"
	result = EvaluateRelease(ledger, nil, nil)
	if result.Status != "failed" || !hasIssue(result.Issues, "detector-not-ready") {
		t.Fatalf("failed detector must block release: %#v", result)
	}
}

func TestReleaseRejectsMissingRenderedEvidence(t *testing.T) {
	ledger := validLedger(t)
	ledger.Gates.RenderedAudit.Status = "pass"
	result := EvaluateRelease(ledger, nil, nil)
	if result.Status != "failed" || !hasIssue(result.Issues, "missing-rendered-evidence") {
		t.Fatalf("release must require rendered evidence: %#v", result)
	}
}
