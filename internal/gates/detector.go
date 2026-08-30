package gates

import "time"

// DetectorFinding is a raw mechanical finding retained in every result.
type DetectorFinding struct {
	ID                 string `json:"id"`
	Rule               string `json:"rule"`
	Path               string `json:"path"`
	Attribution        string `json:"attribution"`
	ExceptionID        string `json:"exception_id,omitempty"`
	ExpiredExceptionID string `json:"expired_exception_id,omitempty"`
}

type DetectorException struct {
	ID        string `json:"id"`
	FindingID string `json:"finding_id"`
	Rule      string `json:"rule"`
	Path      string `json:"path"`
	Reason    string `json:"reason"`
	Risk      string `json:"risk"`
	Owner     string `json:"owner"`
	Evidence  string `json:"evidence"`
	ExpiresAt string `json:"expires_at"`
}

type DetectorResult struct {
	Status   string            `json:"status"`
	Findings []DetectorFinding `json:"findings"`
}

// ClassifyDetectorFindings retains every finding, attributes shared scope, and
// only produces pass-with-exceptions for bounded, unexpired exact matches.
func ClassifyDetectorFindings(findings []DetectorFinding, scope Scope, exceptions []DetectorException, now time.Time) DetectorResult {
	for _, ex := range exceptions {
		if ex.ID == "" || ex.FindingID == "" || ex.Rule == "" || ex.Path == "" || ex.Reason == "" || ex.Risk == "" || ex.Owner == "" || ex.Evidence == "" {
			return DetectorResult{Status: "invalid-configuration", Findings: findings}
		}
		if _, err := time.Parse(time.RFC3339, ex.ExpiresAt); err != nil {
			return DetectorResult{Status: "invalid-configuration", Findings: findings}
		}
	}
	owned, shared := map[string]bool{}, map[string]bool{}
	for _, p := range scope.OwnedPaths {
		owned[p] = true
	}
	for _, p := range scope.SharedPaths {
		shared[p] = true
	}
	unresolved := false
	matched := false
	sharedBoundary := false
	for i := range findings {
		if shared[findings[i].Path] {
			findings[i].Attribution = "shared"
		} else if owned[findings[i].Path] {
			findings[i].Attribution = "owned"
		} else {
			findings[i].Attribution = "uncovered"
		}
		for _, ex := range exceptions {
			expires, err := time.Parse(time.RFC3339, ex.ExpiresAt)
			if err != nil || ex.FindingID != findings[i].ID || ex.Rule != findings[i].Rule || ex.Path != findings[i].Path {
				continue
			}
			if expires.After(now) {
				findings[i].ExceptionID = ex.ID
				matched = true
				break
			}
			findings[i].ExpiredExceptionID = ex.ID
		}
		if findings[i].ExceptionID == "" {
			if findings[i].ExpiredExceptionID != "" || findings[i].Attribution != "shared" {
				unresolved = true
			} else {
				sharedBoundary = true
			}
		}
	}
	if unresolved {
		return DetectorResult{Status: "failed", Findings: findings}
	}
	if matched || sharedBoundary {
		return DetectorResult{Status: "pass-with-exceptions", Findings: findings}
	}
	return DetectorResult{Status: "clean-pass", Findings: findings}
}
