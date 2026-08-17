package lib

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// readSiteTemplateNames returns the set of template filenames in site/web.
// Site templates are the shell/chrome (layout, docs-*, demo-*, recipe-*,
// switchers, landing, blog) — disjoint from the lib component set by
// construction (registry_sync enforces the boundary on the docs side).
func readSiteTemplateNames(t *testing.T) map[string]bool {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	root := filepath.Dir(filename) + "/../.."
	entries, err := os.ReadDir(filepath.Join(root, "site", "web", "templates"))
	if err != nil {
		return nil
	}
	out := map[string]bool{}
	for _, e := range entries {
		if !e.IsDir() {
			out[e.Name()] = true
		}
	}
	return out
}
