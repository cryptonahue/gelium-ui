// Package gates validates versioned Gelium UI workflow evidence.
package gates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const SchemaVersionV1 = 1

// Route controls the amount of workflow evidence required for a UI change.
type Route string

const (
	RouteDirectExempt Route = "direct-exempt"
	RouteDesignGated  Route = "design-gated"
	RouteEscalate     Route = "escalate"
	RouteFullSDD      Route = "full-sdd"
)

// Ledger is the JSON v1 record supplied explicitly to the preflight command.
type Ledger struct {
	SchemaVersion int             `json:"schema_version"`
	Route         Route           `json:"route"`
	Scope         Scope           `json:"scope"`
	Reading       []Reading       `json:"reading"`
	Gates         Gates           `json:"gates"`
	Exceptions    []Exception     `json:"exceptions"`
	Release       ReleaseEvidence `json:"release"`
}

type Scope struct {
	Routes      []string `json:"routes"`
	OwnedPaths  []string `json:"owned_paths"`
	SharedPaths []string `json:"shared_paths"`
}

type Reading struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type Gates struct {
	Plan          Gate `json:"plan"`
	Architecture  Gate `json:"architecture"`
	CriteriaPlan  Gate `json:"criteria_plan"`
	Approval      Gate `json:"approval"`
	RenderedAudit Gate `json:"rendered_audit"`
}

type Gate struct {
	Status   string   `json:"status"`
	Evidence []string `json:"evidence"`
	Packet   string   `json:"packet"`
	Approver string   `json:"approver"`
	Channel  string   `json:"channel"`
	Date     string   `json:"date"`
}

type Exception struct {
	ID                   string `json:"id"`
	Rule                 string `json:"rule"`
	Path                 string `json:"path"`
	Reason               string `json:"reason"`
	Risk                 string `json:"risk"`
	Owner                string `json:"owner"`
	Evidence             string `json:"evidence"`
	ExpiresAt            string `json:"expires_at"`
	ExpiresBeforeVersion string `json:"expires_before_version"`
}

// Issue is a machine-readable validation failure. It reports a record problem;
// it never claims a worker failed to read or a human failed to make a decision.
type Issue struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateLedger parses exactly one JSON v1 ledger and validates its structural
// semantics at now. The caller owns phase-specific preflight decisions.
func ValidateLedger(data []byte, now time.Time) (Ledger, []Issue) {
	var ledger Ledger
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&ledger); err != nil {
		return ledger, []Issue{{Code: "malformed-json", Message: fmt.Sprintf("parse ledger JSON: %v", err)}}
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return ledger, []Issue{{Code: "malformed-json", Message: "ledger must contain one JSON value"}}
		}
		return ledger, []Issue{{Code: "malformed-json", Message: fmt.Sprintf("parse trailing ledger JSON: %v", err)}}
	}

	issues := make([]Issue, 0)
	if ledger.SchemaVersion != SchemaVersionV1 {
		issues = append(issues, issue("unsupported-schema", "schema_version", "must equal 1"))
	}
	if !isRoute(ledger.Route) {
		issues = append(issues, issue("unknown-route", "route", "must be direct-exempt, design-gated, escalate, or full-sdd"))
	}
	issues = append(issues, validateScope(ledger.Scope)...)
	issues = append(issues, validateReading(ledger)...)
	issues = append(issues, validateGates(ledger.Gates)...)
	issues = append(issues, validateExceptions(ledger.Exceptions, now)...)
	return ledger, issues
}

func issue(code, field, message string) Issue {
	return Issue{Code: code, Field: field, Message: message}
}

func isRoute(route Route) bool {
	switch route {
	case RouteDirectExempt, RouteDesignGated, RouteEscalate, RouteFullSDD:
		return true
	default:
		return false
	}
}

func validateScope(scope Scope) []Issue {
	issues := make([]Issue, 0)
	if len(scope.OwnedPaths) == 0 {
		issues = append(issues, issue("missing-owned-path", "scope.owned_paths", "must declare at least one owned path"))
	}
	owned := map[string]struct{}{}
	for _, path := range scope.OwnedPaths {
		if strings.TrimSpace(path) == "" {
			issues = append(issues, issue("invalid-ownership-boundary", "scope.owned_paths", "paths must not be empty"))
			continue
		}
		if _, duplicate := owned[path]; duplicate {
			issues = append(issues, issue("invalid-ownership-boundary", "scope.owned_paths", "owned paths must be unique"))
		}
		owned[path] = struct{}{}
	}
	shared := map[string]struct{}{}
	for _, path := range scope.SharedPaths {
		if strings.TrimSpace(path) == "" {
			issues = append(issues, issue("invalid-ownership-boundary", "scope.shared_paths", "paths must not be empty"))
			continue
		}
		if _, duplicate := shared[path]; duplicate {
			issues = append(issues, issue("invalid-ownership-boundary", "scope.shared_paths", "shared paths must be unique"))
		}
		if _, overlap := owned[path]; overlap {
			issues = append(issues, issue("invalid-ownership-boundary", "scope", "a path cannot be both owned and shared"))
		}
		shared[path] = struct{}{}
	}
	return issues
}

func validateReading(ledger Ledger) []Issue {
	if ledger.Route != RouteDesignGated {
		return nil
	}
	issues := make([]Issue, 0)
	if len(ledger.Reading) == 0 {
		return append(issues, issue("missing-reading-attestation", "reading", "design-gated work requires at least one reading attestation"))
	}
	for i, reading := range ledger.Reading {
		field := fmt.Sprintf("reading[%d]", i)
		if strings.TrimSpace(reading.Path) == "" {
			issues = append(issues, issue("invalid-reading-attestation", field+".path", "path is required"))
		}
		switch reading.Status {
		case "attested", "absent-with-user-record", "not-applicable":
		default:
			issues = append(issues, issue("unknown-reading-status", field+".status", "status is not recognized"))
		}
	}
	return issues
}

func validateGates(gates Gates) []Issue {
	checks := []struct {
		field   string
		gate    Gate
		allowed map[string]bool
	}{
		{"gates.plan", gates.Plan, set("pass", "pending", "fail", "exception")},
		{"gates.architecture", gates.Architecture, set("pass", "pending", "fail", "exception")},
		{"gates.criteria_plan", gates.CriteriaPlan, set("pass", "pass-with-escalation", "pending", "fail", "exception")},
		{"gates.approval", gates.Approval, set("approved", "exempt", "changes-requested", "declined", "draft", "exception", "pending")},
		{"gates.rendered_audit", gates.RenderedAudit, set("pass", "pass-with-escalation", "not-applicable", "pending", "fail", "exception")},
	}
	issues := make([]Issue, 0)
	for _, check := range checks {
		if check.gate.Status != "" && !check.allowed[check.gate.Status] {
			issues = append(issues, issue("unknown-gate-status", check.field+".status", "status is not recognized"))
		}
	}
	return issues
}

func validateExceptions(exceptions []Exception, now time.Time) []Issue {
	issues := make([]Issue, 0)
	for i, exception := range exceptions {
		field := fmt.Sprintf("exceptions[%d]", i)
		for name, value := range map[string]string{
			"id": exception.ID, "rule": exception.Rule, "path": exception.Path,
			"reason": exception.Reason, "risk": exception.Risk, "owner": exception.Owner,
			"evidence": exception.Evidence,
		} {
			if strings.TrimSpace(value) == "" {
				issues = append(issues, issue("invalid-exception", field+"."+name, "field is required"))
			}
		}
		if exception.ExpiresAt == "" && exception.ExpiresBeforeVersion == "" {
			issues = append(issues, issue("invalid-exception-expiry", field, "one deterministic expiry is required"))
		}
		if exception.ExpiresAt != "" {
			expiresAt, err := time.Parse(time.RFC3339, exception.ExpiresAt)
			if err != nil {
				issues = append(issues, issue("invalid-exception-expiry", field+".expires_at", "must use RFC3339"))
			} else if !expiresAt.After(now) {
				issues = append(issues, issue("expired-exception", field+".expires_at", "exception expiry has passed"))
			}
		}
		if exception.ExpiresBeforeVersion != "" && !semverPattern.MatchString(exception.ExpiresBeforeVersion) {
			issues = append(issues, issue("invalid-exception-expiry", field+".expires_before_version", "must use semantic version syntax"))
		}
	}
	return issues
}

func set(values ...string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

var semverPattern = regexp.MustCompile(`^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
