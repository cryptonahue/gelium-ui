package web

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestComponentRegistrySyncsWithTemplates is the Phase J closure guard: every
// template partial in web/templates/*.html must appear in
// docs/gelium-ui-component-registry.md (the registry rule is "code is the
// source of truth", component-registry.md §1). Presence-based only — no prose
// or hex assertions. Fails if a template partial disappears from the registry.
func TestComponentRegistrySyncsWithTemplates(t *testing.T) {
	registryPath := filepath.Join(repositoryRoot(t), "docs", "gelium-ui-component-registry.md")
	registry, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("read component registry: %v", err)
	}
	registryText := string(registry)

	partials, err := filepath.Glob(filepath.Join(repositoryRoot(t), "web", "templates", "*.html"))
	if err != nil {
		t.Fatalf("glob templates: %v", err)
	}

	// Non-component templates: shell/layout, docs shell chrome, demos,
	// recipes, and theme switchers are NOT registry components — they are
	// pages/shell and are excluded from the component registry contract.
	shellPrefixes := []string{"demo-", "docs-", "recipe-"}
	var missing []string
	for _, p := range partials {
		name := filepath.Base(p)
		if name == "layout.html" || name == "landing.html" || strings.Contains(name, "switcher") {
			continue
		}
		excluded := false
		for _, pre := range shellPrefixes {
			if strings.HasPrefix(name, pre) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		if !strings.Contains(registryText, name) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Errorf("template partials missing from docs/gelium-ui-component-registry.md: %v", missing)
	}
}

// TestRegistryReferencedFilesExist guards the reverse direction: every file
// path referenced by the three registries (component, pattern, theme) must
// exist on disk. Presence/existence only.
func TestRegistryReferencedFilesExist(t *testing.T) {
	root := repositoryRoot(t)
	registryDocs := []string{
		filepath.Join(root, "docs", "gelium-ui-component-registry.md"),
		filepath.Join(root, "docs", "gelium-ui-pattern-registry.md"),
		filepath.Join(root, "docs", "gelium-ui-theme-registry.md"),
	}

	for _, doc := range registryDocs {
		content, err := os.ReadFile(doc)
		if err != nil {
			t.Fatalf("read %s: %v", filepath.Base(doc), err)
		}
		text := string(content)
		for _, ext := range []string{".html", ".css", ".go", ".md"} {
			for _, ref := range referencedFileNames(t, text, ext) {
				// refs look like "name.ext" — resolve against the dirs where
				// registry docs point (web/templates, web/styles,
				// internal/app, themes, docs).
				if !registryFileExists(root, ref) {
					t.Errorf("%s references %q but no such file exists", filepath.Base(doc), ref)
				}
			}
		}
	}
}

// referencedFileNames extracts bare filename references (e.g. "button.html")
// from a doc body for a given extension. Tokens with placeholders (<x>) are
// literals, not files — skipped.
func referencedFileNames(t *testing.T, text, ext string) []string {
	t.Helper()
	var out []string
	for _, tok := range strings.Fields(text) {
		clean := strings.Trim(tok, "`|,;()[]{}\"'")
		if strings.HasSuffix(clean, ext) && !strings.Contains(clean, "/") && !strings.Contains(clean, "*") && !strings.Contains(clean, "<") {
			out = append(out, clean)
		}
	}
	return out
}

// registryFileExists resolves a bare filename across the directories the
// registries point at. Docs may reference the canonical prefixed name
// (gelium-ui-*.md, gelium-ui-theme-*.md) or the short name
// (theme-contract.md) — all resolve. handoffs/ lives under docs/.
func registryFileExists(root, name string) bool {
	candidates := []string{name}
	if !strings.HasPrefix(name, "gelium-ui-theme-") {
		candidates = append(candidates, "gelium-ui-theme-"+name)
	}
	if !strings.HasPrefix(name, "gelium-ui-") {
		candidates = append(candidates, "gelium-ui-"+name)
	}
	for _, dir := range []string{"web/templates", "web/styles", "internal/app", "themes/theme-material", "themes/theme-basecoat", "docs", "docs/handoffs", "web"} {
		for _, c := range candidates {
			if _, err := os.Stat(filepath.Join(root, dir, c)); err == nil {
				return true
			}
		}
	}
	return false
}
