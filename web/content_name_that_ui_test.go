package web

// Content contract tests for the Name That UI additions:
//
//  1. Every state pattern in §3 of docs/gelium-ui-vocabulary.md must carry the
//     two new fields — "Nombres alternativos" (alternate names, Name That UI
//     style) and "Prompt para agentes" (a paste-into-your-agent prompt).
//  2. Every served component page listed here must open with the two new
//     English sections — "## Alternative names" and "## Agent prompt".
//
// The sections exist to make the Gelium vocabulary and the served component
// docs usable by AI coding agents (namethatui.com pattern: alternate names +
// copyable prompt). These tests fail if either field/section is missing or
// left empty, so the contract cannot rot silently.

import (
	"strings"
	"testing"
)

// statePatternBlocks splits the §3 "Patrones de estado" section of the
// vocabulary into per-pattern blocks keyed by the "### " heading.
func statePatternBlocks(t *testing.T, vocab string) map[string]string {
	t.Helper()
	start := strings.Index(vocab, "## 3. Patrones de estado")
	end := strings.Index(vocab, "## 4. Patrones de workflow")
	if start < 0 || end < 0 || end <= start {
		t.Fatalf("cannot locate §3 state patterns in docs/gelium-ui-vocabulary.md")
	}
	blocks := map[string]string{}
	current := ""
	for _, line := range strings.Split(vocab[start:end], "\n") {
		if strings.HasPrefix(line, "### ") {
			current = strings.TrimSpace(strings.TrimPrefix(line, "### "))
			blocks[current] = ""
			continue
		}
		if current != "" {
			blocks[current] += line + "\n"
		}
	}
	return blocks
}

// fieldValue returns the non-empty value of a "- **Field**: value" bullet,
// or "" when the field is missing or empty.
func fieldValue(block, field string) string {
	prefix := "- **" + field + "**:"
	for _, line := range strings.Split(block, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
	}
	return ""
}

func wordCount(s string) int {
	return len(strings.Fields(s))
}

// sectionBody returns the text between heading and the next "## " heading.
func sectionBody(t *testing.T, doc, heading string) string {
	t.Helper()
	start := strings.Index(doc, heading)
	if start < 0 {
		t.Fatalf("document is missing section %q", heading)
	}
	rest := doc[start+len(heading):]
	end := strings.Index(rest, "\n## ")
	if end < 0 {
		return strings.TrimSpace(rest)
	}
	return strings.TrimSpace(rest[:end])
}

func TestVocabularyStatePatternsCarryNameThatUIFields(t *testing.T) {
	vocab := repositoryFile(t, "docs", "gelium-ui-vocabulary.md")
	blocks := statePatternBlocks(t, vocab)
	patterns := []string{
		"Empty state ✅ · Persistente",
		"Loading state / Skeleton ✅ · Transitorio",
		"Inline alert ✅ · Persistente",
		"Banner ✅ · Persistente",
		"Callout ✅ · Persistente",
		"Error state ✅ · Persistente",
		"Validation summary ✅ · Persistente",
		"Success feedback ✅ · Persistente",
		"Toast ✅ · Transitorio",
	}
	if len(blocks) != len(patterns) {
		t.Fatalf("§3 state pattern blocks = %d, want %d (%v)", len(blocks), len(patterns), patterns)
	}
	for _, pattern := range patterns {
		block, ok := blocks[pattern]
		if !ok {
			t.Errorf("state pattern %q is missing from §3", pattern)
			continue
		}
		names := fieldValue(block, "Nombres alternativos")
		if wordCount(names) < 2 || !strings.Contains(names, ",") {
			t.Errorf("%q is missing a real 'Nombres alternativos' list, got %q", pattern, names)
		}
		prompt := fieldValue(block, "Prompt para agentes")
		if wordCount(prompt) < 20 {
			t.Errorf("%q is missing a substantial 'Prompt para agentes' (want >= 20 words, got %d)", pattern, wordCount(prompt))
		}
		// The prompt must anchor the pattern contract: when to use it, its
		// server contract, and what NOT to use it for. Spanish negation is
		// sentence-initial ("No lo uses..."), so match case-insensitively.
		lower := strings.ToLower(prompt)
		stated := false
		for _, anchor := range []string{"no lo uses", "nunca", "no uses", "no usar"} {
			if strings.Contains(lower, anchor) {
				stated = true
				break
			}
		}
		if !stated {
			t.Errorf("%q agent prompt does not state what NOT to use the pattern for", pattern)
		}
	}
}

func TestServedComponentPagesCarryNameThatUISections(t *testing.T) {
	for _, slug := range []string{"button", "card", "dialog", "toast", "text-field"} {
		doc := repositoryFile(t, "web", "content", slug+".md")
		altIndex := strings.Index(doc, "## Alternative names")
		promptIndex := strings.Index(doc, "## Agent prompt")
		if altIndex < 0 {
			t.Errorf("%s.md is missing '## Alternative names'", slug)
		}
		if promptIndex < 0 {
			t.Errorf("%s.md is missing '## Agent prompt'", slug)
		}
		if altIndex >= 0 && promptIndex >= 0 && altIndex > promptIndex {
			t.Errorf("%s.md must list '## Alternative names' before '## Agent prompt'", slug)
		}
		// Both sections must carry real content, not empty headings.
		if altIndex >= 0 {
			body := sectionBody(t, doc, "## Alternative names")
			if wordCount(body) < 3 {
				t.Errorf("%s.md '## Alternative names' body is empty", slug)
			}
		}
		if promptIndex >= 0 {
			body := sectionBody(t, doc, "## Agent prompt")
			if wordCount(body) < 20 {
				t.Errorf("%s.md '## Agent prompt' body is too short (want >= 20 words, got %d)", slug, wordCount(body))
			}
		}
	}
}
