package lib

import (
	"testing"
)

// TestLibSiteTemplateNamesAreDisjoint is the collision guard for the dual
// template set: site/web owns the shell/chrome (layout, docs-*, demo-*,
// recipe-*, switchers, landing, blog) and lib owns the component partials.
// A template name defined in BOTH sets would silently shadow depending on
// ParseFS order — this guard fails loudly instead.
func TestLibSiteTemplateNamesAreDisjoint(t *testing.T) {
	libNames := map[string]bool{}
	libEntries, err := LibAssets.ReadDir("templates")
	if err != nil {
		t.Fatalf("read lib templates dir: %v", err)
	}
	for _, e := range libEntries {
		if !e.IsDir() {
			libNames[e.Name()] = true
		}
	}

	siteEntries := readSiteTemplateNames(t)
	var overlap []string
	for name := range libNames {
		if siteEntries[name] {
			overlap = append(overlap, name)
		}
	}
	if len(overlap) > 0 {
		t.Fatalf("template name collision between lib and site (silent shadowing risk): %v", overlap)
	}
}
