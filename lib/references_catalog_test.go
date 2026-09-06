package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPortableReferenceCatalogQuality pins the 0.6.6 contract: npm consumers
// get substantive Markdown fichas (not stub one-liners), every catalog entry
// resolves on disk, required feed/shell IDs exist, and no unlicensed
// screenshots are bundled under lib/references.
func TestPortableReferenceCatalogQuality(t *testing.T) {
	root := filepath.Join(repositoryRoot(t), "lib", "references")
	catalogPath := filepath.Join(root, "catalog.json")
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		t.Fatalf("read catalog.json: %v", err)
	}
	var catalog struct {
		SchemaVersion int `json:"schema_version"`
		References    []struct {
			ID    string `json:"id"`
			Kind  string `json:"kind"`
			Type  string `json:"type"`
			Ficha string `json:"ficha"`
			Tags  []string `json:"tags"`
		} `json:"references"`
	}
	if err := json.Unmarshal(raw, &catalog); err != nil {
		t.Fatalf("parse catalog.json: %v", err)
	}
	if catalog.SchemaVersion < 1 {
		t.Fatalf("schema_version must be >= 1, got %d", catalog.SchemaVersion)
	}
	if len(catalog.References) < 9 {
		t.Fatalf("expected at least 9 catalog entries, got %d", len(catalog.References))
	}

	required := map[string]bool{
		"REF-ARTICLE":            false,
		"REF-404":                false,
		"REF-AUTH":               false,
		"REF-FAQ":                false,
		"REF-HERO":               false,
		"REF-PRICING":            false,
		"REF-SOCIAL-FEED":        false,
		"REF-SHELL":              false,
		"REF-CARD-DETAIL-ENTRY":  false,
	}
	minBytes := map[string]int{
		"component-references/social-feed.md":  2000,
		"component-references/shell.md":        800,
		"component-references/detail-entry.md": 800,
		"section-references/article.md":        500,
	}
	for _, ref := range catalog.References {
		if ref.ID == "" || ref.Ficha == "" {
			t.Fatalf("catalog entry missing id or ficha: %+v", ref)
		}
		if _, ok := required[ref.ID]; ok {
			required[ref.ID] = true
		}
		path := filepath.Join(root, filepath.FromSlash(ref.Ficha))
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("catalog %s ficha missing: %v", ref.ID, err)
		}
		if len(strings.TrimSpace(string(body))) < 200 {
			t.Errorf("%s ficha %s is too short (%d bytes) — portable refs must be substantive", ref.ID, ref.Ficha, len(body))
		}
		if !strings.Contains(string(body), ref.ID) && !strings.HasPrefix(strings.TrimSpace(string(body)), "#") {
			t.Errorf("%s ficha should name the reference or start with a heading", ref.ID)
		}
		if min, ok := minBytes[ref.Ficha]; ok && len(body) < min {
			t.Errorf("%s must be at least %d bytes (got %d) so npm agents are not left with stubs", ref.Ficha, min, len(body))
		}
	}
	for id, seen := range required {
		if !seen {
			t.Errorf("catalog missing required id %s", id)
		}
	}

	social := repositoryFile(t, "lib", "references", "component-references", "social-feed.md")
	for _, needle := range []string{
		"metadata",
		"title",
		"actions",
		"Gelium product filter",
		"bottom nav",
		"fake avatar",
		"REF-CARD-DETAIL-ENTRY",
		"single-column",
	} {
		if !strings.Contains(social, needle) {
			t.Errorf("social-feed.md must document %q for npm-only agents", needle)
		}
	}

	// No binary captures under the publishable references tree.
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".mp4", ".webm":
			t.Errorf("references must not ship media captures (copyright boundary): %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk references: %v", err)
	}
}
