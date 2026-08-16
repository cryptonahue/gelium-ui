package web

// Content contract tests for the Name That UI additions:
//
//  1. Every state pattern in §3 of docs/gelium-ui-vocabulary.md must carry the
//     two new fields — "Nombres alternativos" (alternate names, Name That UI
//     style) and "Prompt para agentes" (a paste-into-your-agent prompt).
//  2. Every served component page listed here must open with the two new
//     English sections — "## Alternative names" and "## Agent prompt".
//  3. Every served component page must carry the normalized usage-guidance
//     block — "## Guidance" with the four subsections "When to use" /
//     "When not to use" / "Usability" / "Accessibility" — so the docs always
//     answer "saber qué componente usar en cada situación". The seven input
//     controls must additionally cross-link the "Choose the right control"
//     decision page.
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

// subsectionBody returns the text between a "### " heading and the next
// heading of the same or higher level ("### " or "## "), or "" when the
// heading is absent.
func subsectionBody(doc, heading string) string {
	start := strings.Index(doc, heading)
	if start < 0 {
		return ""
	}
	rest := doc[start+len(heading):]
	end := len(rest)
	if i := strings.Index(rest, "\n## "); i >= 0 && i < end {
		end = i
	}
	if i := strings.Index(rest, "\n### "); i >= 0 && i < end {
		end = i
	}
	return strings.TrimSpace(rest[:end])
}

// guidanceSection returns the text between the "## Guidance" heading and the
// next "## " heading, or "" when the section is absent. Unlike sectionBody it
// does not fatal: the Guidance contract test wants per-page errors, not a
// test abort on the first missing section.
func guidanceSection(doc string) string {
	start := strings.Index(doc, "## Guidance")
	if start < 0 {
		return ""
	}
	rest := doc[start+len("## Guidance"):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return strings.TrimSpace(rest[:end])
	}
	return strings.TrimSpace(rest)
}

// bulletCount counts markdown "- " bullet lines inside a block.
func bulletCount(block string) int {
	n := 0
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- ") {
			n++
		}
	}
	return n
}

// componentContentSlugs lists every served component page under web/content.
// Handbook pages (handbook-*), the docs root (index.md) and the design
// principles page (principles.md) are NOT components and are deliberately
// excluded: the Guidance contract applies to components only.
var componentContentSlugs = []string{
	"badge", "button", "card", "checkbox", "chips", "data-table", "dialog",
	"divider", "elevation", "fab", "focus-ring", "icon", "icon-button", "list",
	"menu", "navigation-bar", "navigation-drawer", "navigation-tab", "progress",
	"radio", "segmented-button", "select", "slider", "switch", "tabs",
	"text-field", "toast", "tooltip",
}

// guidanceLinkSlugs are the input controls that must cross-link the
// "Choose the right control" decision page from their Guidance section.
var guidanceLinkSlugs = []string{
	"radio", "select", "checkbox", "switch", "slider", "text-field", "menu",
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
		doc := repositoryFile(t, "site", "web", "content", slug+".md")
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

// TestComponentPagesCarryGuidanceSections proves every served component page
// carries the normalized usage-guidance block: "## Guidance" with the four
// subsections "### When to use", "### When not to use", "### Usability" and
// "### Accessibility" in that order, each with real content. The block answers
// "saber qué componente usar en cada situación": the right conditions, the
// alternative components (linked), practical usage bullets, and accessibility
// bullets. The block sits right after the answer-first intro — before any
// "## Alternative names" — so the usage story opens the page. The seven input
// controls must also cross-link the "Choose the right control" decision page.
func TestComponentPagesCarryGuidanceSections(t *testing.T) {
	for _, slug := range componentContentSlugs {
		doc := repositoryFile(t, "site", "web", "content", slug+".md")
		guidance := guidanceSection(doc)
		if guidance == "" {
			t.Errorf("%s.md is missing '## Guidance'", slug)
			continue
		}
		order := []string{"### When to use", "### When not to use", "### Usability", "### Accessibility"}
		last := -1
		for _, heading := range order {
			idx := strings.Index(guidance, heading)
			if idx < 0 {
				t.Errorf("%s.md Guidance is missing %q", slug, heading)
				continue
			}
			if idx < last {
				t.Errorf("%s.md Guidance subsections are out of order: %q must follow the previous subsection", slug, heading)
			}
			last = idx
		}
		if whenToUse := subsectionBody(guidance, "### When to use"); wordCount(whenToUse) < 5 {
			t.Errorf("%s.md 'When to use' is empty or trivial (want >= 5 words, got %d)", slug, wordCount(whenToUse))
		}
		if whenNot := subsectionBody(guidance, "### When not to use"); wordCount(whenNot) < 5 {
			t.Errorf("%s.md 'When not to use' is empty or trivial (want >= 5 words, got %d)", slug, wordCount(whenNot))
		}
		if usability := subsectionBody(guidance, "### Usability"); bulletCount(usability) < 2 {
			t.Errorf("%s.md 'Usability' needs at least 2 bullets (got %d)", slug, bulletCount(usability))
		}
		if a11y := subsectionBody(guidance, "### Accessibility"); bulletCount(a11y) < 2 {
			t.Errorf("%s.md 'Accessibility' needs at least 2 bullets (got %d)", slug, bulletCount(a11y))
		}
		// The Guidance block opens the usage story: it comes after the
		// answer-first intro and before any "## Alternative names" section.
		if alt := strings.Index(doc, "## Alternative names"); alt >= 0 && strings.Index(doc, "## Guidance") > alt {
			t.Errorf("%s.md must place '## Guidance' before '## Alternative names'", slug)
		}
	}
	for _, slug := range guidanceLinkSlugs {
		doc := repositoryFile(t, "site", "web", "content", slug+".md")
		if !strings.Contains(doc, "choose-the-right-control") {
			t.Errorf("%s.md must link the 'Choose the right control' decision page from its Guidance", slug)
		}
	}
}
