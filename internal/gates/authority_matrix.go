package gates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const authorityMatrixPath = "docs/gelium-authority-matrix.json"

// AuthorityDrift identifies one current authority surface that differs from its
// declared canonical value. The checker only reports drift; it never edits a
// contract or repository file.
type AuthorityDrift struct {
	GroupID  string `json:"group_id"`
	Path     string `json:"path"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

// AuthorityMatrixResult reports the coherence of all current authority surfaces.
type AuthorityMatrixResult struct {
	Status string           `json:"status"`
	Drifts []AuthorityDrift `json:"drifts"`
}

type authorityMatrix struct {
	SchemaVersion int                    `json:"schema_version"`
	Groups        []authorityMatrixGroup `json:"groups"`
}

type authorityMatrixGroup struct {
	ID             string          `json:"id"`
	Kind           string          `json:"kind"`
	Canonical      authorityTarget `json:"canonical"`
	CanonicalValue string          `json:"canonical_value"`
	Current        json.RawMessage `json:"current"`
}

type authorityTarget struct {
	Path    string `json:"path"`
	Extract string `json:"extract"`
}

// CheckAuthorityMatrix reads root/docs/gelium-authority-matrix.json and compares
// only its declared current authority surfaces. Historical references are never
// scanned.
func CheckAuthorityMatrix(root string) (AuthorityMatrixResult, error) {
	data, err := os.ReadFile(filepath.Join(root, authorityMatrixPath))
	if err != nil {
		return AuthorityMatrixResult{}, fmt.Errorf("read authority matrix: %w", err)
	}

	var matrix authorityMatrix
	if err := json.Unmarshal(data, &matrix); err != nil {
		return AuthorityMatrixResult{}, fmt.Errorf("parse authority matrix: %w", err)
	}
	if matrix.SchemaVersion != SchemaVersionV1 {
		return AuthorityMatrixResult{}, fmt.Errorf("unsupported authority matrix schema: %d", matrix.SchemaVersion)
	}

	result := AuthorityMatrixResult{Status: "pass", Drifts: []AuthorityDrift{}}
	for _, group := range matrix.Groups {
		var drifts []AuthorityDrift
		switch group.Kind {
		case "version":
			drifts, err = checkVersionAuthority(root, group)
		case "wire":
			drifts, err = checkWireAuthority(root, group)
		default:
			err = fmt.Errorf("group %q has unsupported kind %q", group.ID, group.Kind)
		}
		if err != nil {
			return AuthorityMatrixResult{}, err
		}
		result.Drifts = append(result.Drifts, drifts...)
	}
	if len(result.Drifts) != 0 {
		result.Status = "failed"
	}
	return result, nil
}

func checkVersionAuthority(root string, group authorityMatrixGroup) ([]AuthorityDrift, error) {
	if group.ID == "" || group.Canonical.Path == "" || group.Canonical.Extract == "" {
		return nil, fmt.Errorf("version group requires id and canonical path/extract")
	}
	canonical, err := extractAuthorityValue(root, group.Canonical)
	if err != nil {
		return nil, fmt.Errorf("group %q canonical value: %w", group.ID, err)
	}
	var current []authorityTarget
	if err := json.Unmarshal(group.Current, &current); err != nil {
		return nil, fmt.Errorf("group %q current version targets: %w", group.ID, err)
	}

	drifts := make([]AuthorityDrift, 0)
	for _, target := range current {
		actual, err := extractAuthorityValue(root, target)
		if err != nil {
			return nil, fmt.Errorf("group %q current path %q: %w", group.ID, target.Path, err)
		}
		if actual != canonical {
			drifts = append(drifts, AuthorityDrift{GroupID: group.ID, Path: target.Path, Expected: canonical, Actual: actual})
		}
	}
	return drifts, nil
}

func checkWireAuthority(root string, group authorityMatrixGroup) ([]AuthorityDrift, error) {
	if group.ID == "" || group.CanonicalValue == "" {
		return nil, fmt.Errorf("wire group requires id and canonical_value")
	}
	var current []authorityTarget
	if err := json.Unmarshal(group.Current, &current); err != nil {
		return nil, fmt.Errorf("group %q current wire targets: %w", group.ID, err)
	}

	drifts := make([]AuthorityDrift, 0)
	for _, target := range current {
		actual, err := extractAuthorityValue(root, target)
		if err != nil {
			return nil, fmt.Errorf("group %q current path %q: %w", group.ID, target.Path, err)
		}
		if actual != group.CanonicalValue {
			drifts = append(drifts, AuthorityDrift{GroupID: group.ID, Path: target.Path, Expected: group.CanonicalValue, Actual: actual})
		}
	}
	return drifts, nil
}

func extractAuthorityValue(root string, target authorityTarget) (string, error) {
	if target.Path == "" || target.Extract == "" {
		return "", fmt.Errorf("path and extract are required")
	}
	data, err := os.ReadFile(filepath.Join(root, target.Path))
	if err != nil {
		return "", err
	}
	content := string(data)

	switch {
	case target.Extract == "json:/version":
		var document struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(data, &document); err != nil {
			return "", err
		}
		if document.Version == "" {
			return "", fmt.Errorf("json pointer /version is empty")
		}
		return document.Version, nil
	case strings.HasPrefix(target.Extract, "go-const:"):
		name := strings.TrimPrefix(target.Extract, "go-const:")
		match := regexp.MustCompile(`(?m)\bconst\s+` + regexp.QuoteMeta(name) + `\s*=\s*"([^"]+)"`).FindStringSubmatch(content)
		if len(match) != 2 {
			return "", fmt.Errorf("Go constant %q not found", name)
		}
		return match[1], nil
	case strings.HasPrefix(target.Extract, "line-prefix:"):
		prefix := strings.TrimPrefix(target.Extract, "line-prefix:")
		for _, line := range strings.Split(content, "\n") {
			if strings.HasPrefix(line, prefix) {
				value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				if value != "" {
					return value, nil
				}
			}
		}
		return "", fmt.Errorf("line with prefix %q not found", prefix)
	case target.Extract == "release-heading":
		match := regexp.MustCompile(`(?m)^Current release:\s+\*\*v?([^*\s]+)\*\*`).FindStringSubmatch(content)
		if len(match) != 2 {
			return "", fmt.Errorf("current release heading not found")
		}
		return match[1], nil
	case target.Extract == "first-release-heading":
		match := regexp.MustCompile(`(?m)^## \[?v?([0-9]+\.[0-9]+\.[0-9]+[^\]\s]*)\]?`).FindStringSubmatch(content)
		if len(match) != 2 {
			return "", fmt.Errorf("first release heading not found")
		}
		return match[1], nil
	case strings.HasPrefix(target.Extract, "go-http-header:"):
		return extractMatchingLiteral(content, regexp.MustCompile(`\.Header\(\)\.Set\(\s*"([^"]+)"`), strings.TrimPrefix(target.Extract, "go-http-header:"), "HTTP header setter")
	case strings.HasPrefix(target.Extract, "js-response-header:"):
		return extractMatchingLiteral(content, regexp.MustCompile(`response\.headers\.get\("([^"]+)"\)`), strings.TrimPrefix(target.Extract, "js-response-header:"), "response header lookup")
	case strings.HasPrefix(target.Extract, "js-event-listener:"):
		return extractMatchingLiteral(content, regexp.MustCompile(`document\.addEventListener\("([^"]+)"`), strings.TrimPrefix(target.Extract, "js-event-listener:"), "event listener")
	case strings.HasPrefix(target.Extract, "go-json-tag:"):
		return extractMatchingLiteral(content, regexp.MustCompile(`json:"([^"]+)"`), strings.TrimPrefix(target.Extract, "go-json-tag:"), "JSON tag")
	default:
		return "", fmt.Errorf("unsupported extract %q", target.Extract)
	}
}

func extractMatchingLiteral(content string, pattern *regexp.Regexp, expected, label string) (string, error) {
	matches := pattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("%s not found", label)
	}
	for _, match := range matches {
		if match[1] == expected {
			return expected, nil
		}
	}
	return matches[0][1], nil
}

func closestWireValue(canonical, content string) string {
	var pattern *regexp.Regexp
	switch {
	case strings.HasPrefix(canonical, "X-"):
		pattern = regexp.MustCompile(`\bX-[A-Za-z0-9-]+\b`)
	case strings.Contains(canonical, ":"):
		pattern = regexp.MustCompile(`\b[A-Za-z][A-Za-z0-9-]*:[A-Za-z][A-Za-z0-9-]*\b`)
	default:
		return ""
	}
	for _, candidate := range pattern.FindAllString(content, -1) {
		if candidate != canonical {
			return candidate
		}
	}
	return ""
}
