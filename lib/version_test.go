package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestAssetsVersionCoherence pins the S4.6 cache-busting contract: the single
// version constant (lib.AssetsVersion) must match the npm package version, and
// every template that renders a static asset URL must reference
// {{.AssetsVersion}} instead of a hardcoded ?v= literal. This is the guard
// that prevents the 0.5.2-vs-0.5.0 drift from ever coming back.
func TestAssetsVersionCoherence(t *testing.T) {
	// 1. npm package version matches the constant.
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "package.json"))
	if err != nil {
		t.Fatalf("read lib/package.json: %v", err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(b, &pkg); err != nil {
		t.Fatalf("parse lib/package.json: %v", err)
	}
	if pkg.Version != AssetsVersion {
		t.Errorf("lib.AssetsVersion = %q but lib/package.json version = %q (must match)", AssetsVersion, pkg.Version)
	}

	// 2. No template may hardcode a ?v= literal.
	entries, err := os.ReadDir(filepath.Join(repositoryRoot(t), "site", "web", "templates"))
	if err != nil {
		t.Fatalf("list site templates: %v", err)
	}
	hardcoded := regexp.MustCompile(`\?v=\d`)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		body, err := os.ReadFile(filepath.Join(repositoryRoot(t), "site", "web", "templates", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		if m := hardcoded.Find(body); m != nil {
			t.Errorf("%s hardcodes a cache-buster %q — use {{.AssetsVersion}} instead", e.Name(), m)
		}
	}
}
