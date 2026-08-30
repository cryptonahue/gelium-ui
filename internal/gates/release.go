package gates

// ReleaseEvidence records postbuild evidence. It is evidence status only and
// never grants permission to commit, push, publish, or deploy.
type ReleaseEvidence struct {
	Detector        DetectorEvidence `json:"detector"`
	Checks          CheckEvidence    `json:"checks"`
	AuthorityMatrix StatusEvidence   `json:"authority_matrix"`
}

type DetectorEvidence struct {
	Status   string   `json:"status"`
	Evidence []string `json:"evidence"`
}

type CheckEvidence struct {
	Tests  []string `json:"tests"`
	Builds []string `json:"builds"`
}

type StatusEvidence struct {
	Status   string   `json:"status"`
	Evidence []string `json:"evidence"`
}

// EvaluateRelease validates a candidate's postbuild evidence after prebuild
// gates have passed. Its result is intentionally limited to evidence status.
func EvaluateRelease(ledger Ledger, validationIssues []Issue, changedPaths []string) PreflightResult {
	prebuild := EvaluatePrebuild(ledger, validationIssues, changedPaths)
	if prebuild.Status == "invalid-configuration" {
		return prebuild
	}
	issues := append([]Issue{}, prebuild.Issues...)
	if ledger.Gates.RenderedAudit.Status != "pass" && ledger.Gates.RenderedAudit.Status != "pass-with-escalation" && ledger.Gates.RenderedAudit.Status != "exception" {
		issues = append(issues, issue("rendered-audit-not-ready", "gates.rendered_audit", "rendered audit is not ready for Release"))
	}
	if len(ledger.Gates.RenderedAudit.Evidence) == 0 {
		issues = append(issues, issue("missing-rendered-evidence", "gates.rendered_audit.evidence", "release requires rendered audit evidence"))
	}
	if ledger.Release.Detector.Status != "clean-pass" && ledger.Release.Detector.Status != "pass-with-exceptions" {
		issues = append(issues, issue("detector-not-ready", "release.detector.status", "release requires clean-pass or pass-with-exceptions detector output"))
	}
	if len(ledger.Release.Detector.Evidence) == 0 {
		issues = append(issues, issue("missing-detector-evidence", "release.detector.evidence", "release requires detector output evidence"))
	}
	if len(ledger.Release.Checks.Tests) == 0 || len(ledger.Release.Checks.Builds) == 0 {
		issues = append(issues, issue("missing-check-evidence", "release.checks", "release requires test and build evidence"))
	}
	if ledger.Release.AuthorityMatrix.Status != "pass" || len(ledger.Release.AuthorityMatrix.Evidence) == 0 {
		issues = append(issues, issue("authority-matrix-not-ready", "release.authority_matrix", "release requires a passing authority-matrix result"))
	}
	if len(issues) != 0 {
		return PreflightResult{Status: "failed", Issues: issues}
	}
	return PreflightResult{Status: "pass", Issues: []Issue{}}
}
