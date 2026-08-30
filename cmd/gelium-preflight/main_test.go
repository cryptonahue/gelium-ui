package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunPrebuildWritesMachineReadableResult(t *testing.T) {
	ledger := `{"schema_version":1,"route":"direct-exempt","scope":{"owned_paths":["templates/feed.html"]}}`
	path := writeLedger(t, ledger)
	var output bytes.Buffer
	code := run([]string{"prebuild", "--ledger", path, "--changed", "templates/feed.html", "--format", "json"}, &output, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if code != 0 || output.String() != "{\"status\":\"pass\",\"issues\":[]}\n" {
		t.Fatalf("unexpected prebuild result: code=%d output=%q", code, output.String())
	}
}

func TestRunReleaseWritesEvidenceStatus(t *testing.T) {
	ledger := `{"schema_version":1,"route":"design-gated","scope":{"owned_paths":["templates/feed.html"]},"reading":[{"path":"PRODUCT.md","status":"attested"}],"gates":{"plan":{"status":"pass"},"architecture":{"status":"pass"},"criteria_plan":{"status":"pass"},"approval":{"status":"approved"},"rendered_audit":{"status":"pass","evidence":["audit.md"]}},"release":{"detector":{"status":"clean-pass","evidence":["detector.json"]},"checks":{"tests":["go test ./..."],"builds":["npm run build"]},"authority_matrix":{"status":"pass","evidence":["authority.json"]}}}`
	path := writeLedger(t, ledger)
	var output bytes.Buffer
	code := run([]string{"release", "--ledger", path, "--format", "json"}, &output, time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if code != 0 || output.String() != "{\"status\":\"pass\",\"issues\":[]}\n" {
		t.Fatalf("unexpected release result: code=%d output=%q", code, output.String())
	}
}

func writeLedger(t *testing.T, ledger string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ledger.json")
	if err := os.WriteFile(path, []byte(ledger), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
