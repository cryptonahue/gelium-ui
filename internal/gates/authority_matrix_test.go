package gates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAuthorityMatrixIsVersionedAndNamesCurrentWireGroups(t *testing.T) {
	b, err := os.ReadFile("../../docs/gelium-authority-matrix.json")
	if err != nil {
		t.Fatal(err)
	}
	var v struct {
		SchemaVersion int `json:"schema_version"`
		Groups        []struct {
			ID string `json:"id"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if v.SchemaVersion != 1 {
		t.Fatalf("schema=%d", v.SchemaVersion)
	}
	seen := map[string]bool{}
	for _, g := range v.Groups {
		seen[g.ID] = true
	}
	for _, id := range []string{"gelium-ui-release-version", "validation-header", "toast-event"} {
		if !seen[id] {
			t.Errorf("missing %s", id)
		}
	}
}

func TestCheckAuthorityMatrixReadsCurrentRepositoryMatrix(t *testing.T) {
	result, err := CheckAuthorityMatrix("../..")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "pass" {
		t.Fatalf("status=%q drifts=%+v", result.Status, result.Drifts)
	}
	if len(result.Drifts) != 0 {
		t.Fatalf("drifts=%+v", result.Drifts)
	}
}

func TestCheckAuthorityMatrixReportsVersionDriftByPathAndValue(t *testing.T) {
	root := t.TempDir()
	writeAuthorityTestFile(t, root, "docs/gelium-authority-matrix.json", `{
  "schema_version": 1,
  "groups": [{
    "id": "release-version",
    "kind": "version",
    "canonical": {"path": "lib/package.json", "extract": "json:/version"},
    "current": [{"path": "lib/version.go", "extract": "go-const:AssetsVersion"}]
  }]
}`)
	writeAuthorityTestFile(t, root, "lib/package.json", `{"version":"0.6.0"}`)
	writeAuthorityTestFile(t, root, "lib/version.go", `package lib
const AssetsVersion = "0.5.0"
`)

	result, err := CheckAuthorityMatrix(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []AuthorityDrift{{GroupID: "release-version", Path: "lib/version.go", Expected: "0.6.0", Actual: "0.5.0"}}
	if result.Status != "failed" || !reflect.DeepEqual(result.Drifts, want) {
		t.Fatalf("result=%+v want status=failed drifts=%+v", result, want)
	}
}

func TestCheckAuthorityMatrixReportsWireDriftByPathAndValue(t *testing.T) {
	root := t.TempDir()
	writeAuthorityTestFile(t, root, "docs/gelium-authority-matrix.json", `{
  "schema_version": 1,
  "groups": [{
    "id": "validation-header",
    "kind": "wire",
    "canonical_value": "X-Gelium-Validation",
    "current": [{"path": "internal/app/text_field.go", "extract": "go-http-header:X-Gelium-Validation"}]
  }]
}`)
	writeAuthorityTestFile(t, root, "internal/app/text_field.go", `// X-Gelium-Validation remains the documented contract.
w.Header().Set("X-Loom-Validation", "true")`)

	result, err := CheckAuthorityMatrix(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []AuthorityDrift{{GroupID: "validation-header", Path: "internal/app/text_field.go", Expected: "X-Gelium-Validation", Actual: "X-Loom-Validation"}}
	if result.Status != "failed" || !reflect.DeepEqual(result.Drifts, want) {
		t.Fatalf("result=%+v want status=failed drifts=%+v", result, want)
	}
}

func writeAuthorityTestFile(t *testing.T, root, path, content string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
