package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLibPackageJsonIsPublishable pins the npm packaging contract (S4):
// gelium-ui is the publishable package (private:false), ships the styles,
// templates, js and dist entries, and declares an exports map so consumers
// can import subpaths (styles/*, templates/*, js/*, dist/*) without deep
// package internals.
func TestLibPackageJsonIsPublishable(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "package.json"))
	if err != nil {
		t.Fatalf("read lib/package.json: %v", err)
	}
	var pkg struct {
		Name        string            `json:"name"`
		Version     string            `json:"version"`
		Private     bool              `json:"private"`
		Files       []string          `json:"files"`
		SideEffects []string          `json:"sideEffects"`
		Exports     map[string]string `json:"exports"`
	}
	if err := json.Unmarshal(b, &pkg); err != nil {
		t.Fatalf("parse lib/package.json: %v", err)
	}
	if pkg.Name != "gelium-ui" {
		t.Errorf("package name must be gelium-ui, got %q", pkg.Name)
	}
	if pkg.Private {
		t.Error("lib package must be publishable (private must be false)")
	}
	if len(pkg.Files) == 0 {
		t.Error("files must list published entries (styles, templates, js, dist)")
	}
	hasSkill := false
	for _, entry := range pkg.Files {
		if entry == "SKILL.md" {
			hasSkill = true
			break
		}
	}
	if !hasSkill {
		t.Error("files must include SKILL.md for agent skill registries")
	}
	hasReferences := false
	for _, entry := range pkg.Files {
		if entry == "references" {
			hasReferences = true
			break
		}
	}
	if !hasReferences {
		t.Error("files must include the portable agent reference catalog")
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot(t), "lib", "references", "catalog.json")); err != nil {
		t.Fatalf("portable reference catalog missing: %v", err)
	}
	if len(pkg.SideEffects) == 0 {
		t.Error("sideEffects must declare *.css so bundlers keep the CSS imports")
	}
	// Exports must cover the entries the consumer contract relies on.
	for _, sub := range []string{".", "./styles/*", "./themes/*", "./templates/*", "./js/*", "./dist/*"} {
		if _, ok := pkg.Exports[sub]; !ok {
			t.Errorf("exports missing %q", sub)
		}
	}
}
