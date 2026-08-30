package lib

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractUsedIconsEmbedsOnlyReferencedGlyphs(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	scan := filepath.Join(dir, "settings.templ")
	if err := os.WriteFile(scan, []byte(`<a class="ui-list-item-link" href="/deepfilter/settings/appearance" data-gelium-icon="chevron_right"><span>Tema</span></a>
const _ = icons.SVG("settings")
`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "icons.go")
	cmd := exec.Command("node", script, "--scan", dir, "--out", out, "--package", "icons")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extract-used-icons: %v\n%s", err, got)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, want := range []string{
		"package icons",
		`"chevron_right":`,
		`"settings":`,
		`class="ui-icon"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`fill="currentColor"`,
		"func SVG(name string)",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("generated icons missing %q\n%s", want, text)
		}
	}
	if strings.Contains(text, `"delete":`) {
		t.Errorf("unused catalog glyph leaked into the consumer embed")
	}
	if strings.Count(text, `"chevron_right":`) != 1 {
		t.Errorf("duplicate chevron_right rows")
	}
}

func TestExtractUsedIconsDefaultSetIsMaterial(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(`<span data-gelium-icon="home"></span>`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "icons.go")
	cmd := exec.Command("node", script, "--scan", dir, "--out", out, "--package", "icons")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extract-used-icons: %v\n%s", err, got)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, `fill="currentColor"`) {
		t.Fatalf("default set should emit Material fill glyphs:\n%s", text)
	}
	if strings.Contains(text, `stroke="currentColor"`) {
		t.Fatalf("default set should not emit Tabler stroke glyphs:\n%s", text)
	}
}

func TestExtractUsedIconsTablerSetUsesStrokeCatalog(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(`<span data-gelium-icon="chevron-right"></span>
const _ = icons.SVG("settings")
`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "icons.go")
	cmd := exec.Command("node", script, "--scan", dir, "--out", out, "--package", "icons", "--set", "tabler")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extract-used-icons tabler: %v\n%s", err, got)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, want := range []string{
		`"chevron-right":`,
		`"settings":`,
		`class="ui-icon"`,
		`data-gelium-set="tabler"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`stroke="currentColor"`,
		`fill="none"`,
	} {
		if !strings.Contains(text, want) {
			t.Errorf("tabler embed missing %q\n%s", want, text)
		}
	}
	// stroke-width is required for Tabler outline; reject only layout size attrs.
	if strings.Contains(text, ` width="`) || strings.Contains(text, ` height="`) ||
		strings.Contains(text, "\nwidth=\"") || strings.Contains(text, "\nheight=\"") {
		t.Errorf("tabler embed must not hardcode layout width/height")
	}
	if strings.Contains(text, `"delete":`) {
		t.Errorf("unused tabler glyph leaked into the consumer embed")
	}
}

func TestExtractUsedIconsPrefixedNamesPickCatalog(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(`<span data-gelium-icon="ms:home"></span>
<span data-gelium-icon="tabler:chevron-right"></span>
`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "icons.go")
	cmd := exec.Command("node", script, "--scan", dir, "--out", out, "--package", "icons")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extract-used-icons mixed: %v\n%s", err, got)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, `"ms:home":`) && !strings.Contains(text, `"home":`) {
		t.Fatalf("missing material home key:\n%s", text)
	}
	if !strings.Contains(text, `"tabler:chevron-right":`) {
		t.Fatalf("missing tabler prefixed key:\n%s", text)
	}
	if !strings.Contains(text, `fill="currentColor"`) {
		t.Errorf("material glyph should keep fill=currentColor")
	}
	if !strings.Contains(text, `stroke="currentColor"`) {
		t.Errorf("tabler glyph should keep stroke=currentColor")
	}
}

func TestExtractUsedIconsRejectsUnknownGlyph(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(`<span data-gelium-icon="not_a_real_gelium_icon"></span>`), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", script, "--scan", dir, "--out", filepath.Join(dir, "icons.go"), "--package", "icons")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected unknown glyph to fail, output:\n%s", got)
	}
	if !strings.Contains(string(got), "not_a_real_gelium_icon") {
		t.Fatalf("error should name the missing glyph, got:\n%s", got)
	}
}

func TestExtractUsedIconsRejectsUnknownTablerGlyph(t *testing.T) {
	root := repositoryRoot(t)
	script := filepath.Join(root, "lib", "scripts", "extract-used-icons.mjs")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(`<span data-gelium-icon="tabler:not-a-real-tabler-icon"></span>`), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", script, "--scan", dir, "--out", filepath.Join(dir, "icons.go"), "--package", "icons")
	cmd.Dir = root
	got, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected unknown tabler glyph to fail, output:\n%s", got)
	}
	if !strings.Contains(string(got), "not-a-real-tabler-icon") {
		t.Fatalf("error should name the missing tabler glyph, got:\n%s", got)
	}
}
