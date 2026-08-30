package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunReportsJSONForScopedFinding(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "x.css")
	os.WriteFile(p, []byte(".topbar { color:#123456 }"), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", p, "--format", "json", p}, &o); c != 1 {
		t.Fatal(c)
	}
}
func TestRunRejectsMissingScopeAndUnreadableInputs(t *testing.T) {
	var o bytes.Buffer
	if c := run([]string{"--format", "json"}, &o); c != 2 {
		t.Fatal(c)
	}
	if c := run([]string{"--owned", "missing.css", "missing.css"}, &o); c != 2 {
		t.Fatal(c)
	}
}
func TestRunReportsExceptionWithoutHidingFinding(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "x.css")
	os.WriteFile(p, []byte(".topbar {}"), 0600)
	m := filepath.Join(d, "e.json")
	os.WriteFile(m, []byte(`[{"id":"x","finding_id":"custom-shell","rule":"shell-contract","path":"`+p+`","reason":"r","risk":"r","owner":"o","evidence":"e","expires_at":"2099-01-01T00:00:00Z"}]`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", p, "--exceptions", m, "--format", "json", p}, &o); c != 0 {
		t.Fatal(c)
	}
}

func TestRunEvaluatesEachFormAndRestrictsDarkMediaToCSS(t *testing.T) {
	d := t.TempDir()
	forms := filepath.Join(d, "forms.html")
	os.WriteFile(forms, []byte(`<form><div class="validation-summary"></div></form><form><input name="email"></form>`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", forms, forms}, &o); c != 1 {
		t.Fatalf("mixed valid/invalid forms must fail, code=%d output=%s", c, o.String())
	}

	comment := filepath.Join(d, "note.go")
	os.WriteFile(comment, []byte(`// prefers-color-scheme: dark is documentation, not CSS`), 0600)
	o.Reset()
	if c := run([]string{"--owned", comment, "--format", "json", comment}, &o); c != 0 {
		t.Fatalf("non-CSS dark media mention must not fail, code=%d output=%s", c, o.String())
	}

	darkCSS := filepath.Join(d, "dark.css")
	os.WriteFile(darkCSS, []byte(`@media (prefers-color-scheme: dark) { .x { color: red; } }`), 0600)
	o.Reset()
	if c := run([]string{"--owned", darkCSS, darkCSS}, &o); c != 1 {
		t.Fatalf("CSS dark media must fail, code=%d output=%s", c, o.String())
	}
}

func TestRunRejectsOverlappingOrUncoveredScopeAndEncodesEmptyFindingsAsArray(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "safe.html")
	os.WriteFile(p, []byte(`<main>safe</main>`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", p, "--shared", p, p}, &o); c != 2 {
		t.Fatalf("overlapping scope must be invalid configuration, code=%d", c)
	}
	if c := run([]string{"--owned", filepath.Join(d, "declared.html"), p}, &o); c != 2 {
		t.Fatalf("uncovered scan input must be invalid configuration, code=%d", c)
	}
	o.Reset()
	if c := run([]string{"--owned", p, "--format", "json", p}, &o); c != 0 {
		t.Fatalf("safe input must pass, code=%d output=%s", c, o.String())
	}
	var result struct {
		Findings []json.RawMessage `json:"findings"`
	}
	if err := json.Unmarshal(o.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Findings == nil || len(result.Findings) != 0 {
		t.Fatalf("findings must encode as an empty array: %#v", result.Findings)
	}
}

func TestRunDoesNotTreatIDSelectorsAsColorLiterals(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "fragment.html")
	os.WriteFile(p, []byte(`<section id="feed-panel" hx-target="#feed-panel"></section>`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", p, "--format", "json", p}, &o); c != 0 {
		t.Fatalf("ID selectors are not hex colors, code=%d output=%s", c, o.String())
	}
}

func TestRunReportsSharedFindingWithoutHidingIt(t *testing.T) {
	d := t.TempDir()
	owned := filepath.Join(d, "owned.html")
	shared := filepath.Join(d, "shared.css")
	os.WriteFile(owned, []byte(`<main>owned</main>`), 0600)
	os.WriteFile(shared, []byte(`.topbar { display: flex; }`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", owned, "--shared", shared, "--format", "json", owned, shared}, &o); c != 0 || !bytes.Contains(o.Bytes(), []byte(`"status":"pass-with-exceptions"`)) || !bytes.Contains(o.Bytes(), []byte(`"attribution":"shared"`)) {
		t.Fatalf("shared finding must be visible and non-clean without failing owned scope, code=%d output=%s", c, o.String())
	}
}

func TestRunReportsUnknownInformativeMediaWithoutInventingDimensions(t *testing.T) {
	d := t.TempDir()
	known := filepath.Join(d, "known.html")
	os.WriteFile(known, []byte(`<img src="cover.jpg" alt="Product cover" width="1200" height="800">`), 0600)
	var o bytes.Buffer
	if c := run([]string{"--owned", known, "--format", "json", known}, &o); c != 0 {
		t.Fatalf("known dimensions must pass, code=%d output=%s", c, o.String())
	}

	decorative := filepath.Join(d, "decorative.html")
	os.WriteFile(decorative, []byte(`<img src="flourish.svg" alt="">`), 0600)
	o.Reset()
	if c := run([]string{"--owned", decorative, "--format", "json", decorative}, &o); c != 0 {
		t.Fatalf("decorative media may omit dimensions, code=%d output=%s", c, o.String())
	}

	unknown := filepath.Join(d, "unknown.html")
	os.WriteFile(unknown, []byte(`<img src="https://cdn.example/image" alt="Campaign photo" class="ui-media">`), 0600)
	o.Reset()
	if c := run([]string{"--owned", unknown, "--format", "json", unknown}, &o); c != 1 {
		t.Fatalf("unknown informative media must be visible, code=%d output=%s", c, o.String())
	}
	if !bytes.Contains(o.Bytes(), []byte(`"id":"media-metadata-unknown"`)) {
		t.Fatalf("unknown media output must use the honest vocabulary: %s", o.String())
	}
}

func TestRunDistinguishesMalformedExpiredAndUnmatchedExceptions(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "shell.css")
	os.WriteFile(p, []byte(`.topbar { display: flex; }`), 0600)
	manifest := func(name, body string) string {
		path := filepath.Join(d, name)
		if err := os.WriteFile(path, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	invalid := manifest("invalid.json", `[{"id":"missing-fields"}]`)
	var o bytes.Buffer
	if c := run([]string{"--owned", p, "--exceptions", invalid, p}, &o); c != 2 {
		t.Fatalf("malformed exception must be invalid configuration, code=%d output=%s", c, o.String())
	}
	o.Reset()
	if c := run([]string{"--owned", p, "--exceptions", invalid, "--format", "json", p}, &o); c != 2 || !bytes.Contains(o.Bytes(), []byte(`"status":"invalid-configuration"`)) || !bytes.Contains(o.Bytes(), []byte(`"findings":[]`)) {
		t.Fatalf("malformed manifest must have a machine-readable invalid result, code=%d output=%s", c, o.String())
	}
	expired := manifest("expired.json", `[{"id":"expired","finding_id":"custom-shell","rule":"shell-contract","path":"`+p+`","reason":"migration","risk":"visual","owner":"maintainer","evidence":"ledger","expires_at":"2000-01-01T00:00:00Z"}]`)
	o.Reset()
	if c := run([]string{"--owned", p, "--exceptions", expired, "--format", "json", p}, &o); c != 1 || !bytes.Contains(o.Bytes(), []byte(`"expired_exception_id":"expired"`)) {
		t.Fatalf("expired exception must remain a visible unresolved finding, code=%d output=%s", c, o.String())
	}
	unmatched := manifest("unmatched.json", `[{"id":"other","finding_id":"other","rule":"shell-contract","path":"`+p+`","reason":"migration","risk":"visual","owner":"maintainer","evidence":"ledger","expires_at":"2099-01-01T00:00:00Z"}]`)
	o.Reset()
	if c := run([]string{"--owned", p, "--exceptions", unmatched, p}, &o); c != 1 {
		t.Fatalf("unmatched exception must fail, code=%d output=%s", c, o.String())
	}
}
