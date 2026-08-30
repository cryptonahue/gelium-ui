package gates

import "strings"

// PreflightResult is evidence status only. It never grants delivery authority.
type PreflightResult struct {
	Status string  `json:"status"`
	Issues []Issue `json:"issues"`
}

// EvaluatePrebuild validates the facts that must exist before Build. Rendered
// audit evidence is deliberately excluded: it cannot exist before a candidate.
func EvaluatePrebuild(ledger Ledger, validationIssues []Issue, changedPaths []string) PreflightResult {
	if len(validationIssues) != 0 {
		return PreflightResult{Status: "invalid-configuration", Issues: validationIssues}
	}

	issues := scopeIssues(ledger.Scope, changedPaths)
	if ledger.Route == RouteDesignGated {
		for _, requirement := range []struct {
			name string
			gate Gate
			ok   map[string]bool
		}{
			{"plan", ledger.Gates.Plan, set("pass", "exception")},
			{"architecture", ledger.Gates.Architecture, set("pass", "exception")},
			{"criteria plan", ledger.Gates.CriteriaPlan, set("pass", "pass-with-escalation", "exception")},
			{"approval", ledger.Gates.Approval, set("approved", "exempt", "exception")},
		} {
			if !requirement.ok[requirement.gate.Status] {
				issues = append(issues, issue("approval-blocked", "gates."+strings.ReplaceAll(requirement.name, " ", "_"), requirement.name+" gate is not ready for Build"))
			}
		}
	}
	if len(issues) != 0 {
		return PreflightResult{Status: "failed", Issues: issues}
	}
	return PreflightResult{Status: "pass", Issues: []Issue{}}
}

func scopeIssues(scope Scope, changedPaths []string) []Issue {
	declared := make(map[string]struct{}, len(scope.OwnedPaths)+len(scope.SharedPaths))
	for _, path := range scope.OwnedPaths {
		declared[path] = struct{}{}
	}
	for _, path := range scope.SharedPaths {
		declared[path] = struct{}{}
	}
	issues := make([]Issue, 0)
	for _, changed := range changedPaths {
		if _, covered := declared[changed]; !covered {
			issues = append(issues, issue("uncovered-changed-path", "changed_paths", changed+" is outside declared owned/shared scope"))
		}
	}
	return issues
}
