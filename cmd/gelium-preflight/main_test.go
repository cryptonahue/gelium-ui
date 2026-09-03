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

func TestRunRouteWritesDelegatedDirectPlan(t *testing.T) {
	var output bytes.Buffer
	code := run([]string{"route", "--route", "delegated-direct", "--format", "json"}, &output, time.Now())
	if code != 0 {
		t.Fatalf("route command failed: code=%d output=%q", code, output.String())
	}
	want := `{"route":"delegated-direct","status":"pass","next":"delegate-bounded-action","requires_design_gate":false}` + "\n"
	if output.String() != want {
		t.Fatalf("unexpected route result: got=%q want=%q", output.String(), want)
	}
}

func TestRunRouteCoversEveryRoute(t *testing.T) {
	for _, tc := range []struct {
		route string
		next  string
		code  int
	}{
		{route: "direct-exempt", next: "implement-bounded-action"},
		{route: "delegated-direct", next: "delegate-bounded-action"},
		{route: "design-gated", next: "prepare-visible-architecture-packet"},
		{route: "full-sdd", next: "start-openspec"},
		{route: "escalate", next: "request-concrete-decision", code: 1},
	} {
		t.Run(tc.route, func(t *testing.T) {
			var output bytes.Buffer
			code := run([]string{"route", "--route", tc.route, "--format", "json"}, &output, time.Now())
			if code != tc.code {
				t.Fatalf("route exit code = %d, want %d; output=%q", code, tc.code, output.String())
			}
			if !bytes.Contains(output.Bytes(), []byte(`"next":"`+tc.next+`"`)) {
				t.Fatalf("route output = %q, want next=%q", output.String(), tc.next)
			}
		})
	}
}

func TestRunRouteRejectsUnknownRoute(t *testing.T) {
	var output bytes.Buffer
	code := run([]string{"route", "--route", "unknown", "--format", "json"}, &output, time.Now())
	if code != 2 || !bytes.Contains(output.Bytes(), []byte(`"status":"invalid-configuration"`)) {
		t.Fatalf("unknown route = code %d, output %q; want configuration failure", code, output.String())
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
