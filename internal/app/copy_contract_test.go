package app

import (
	"io/fs"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	webassets "geliumui/site/web"
)

// Copy contract (content style guide, web/content/handbook-content-style.md):
// every recipe error message follows the action pattern — the fix is the
// message ("Enter the project name."), never a consequence-only "Name is
// required." — and every recipe empty state pairs the reason with a real next
// step. The style guide is enforced here and by the Handbook route tests
// (handbook_test.go).
func TestRecipeErrorCopyUsesActionPattern(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource", url.Values{
		"name":   {"  "},
		"status": {"On hold"},
		"date":   {""},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	body := res.Body.String()
	for _, want := range []string{
		`href="#recipe-ar-name-error">Enter the project name.</a>`,
		`href="#recipe-ar-status-error">Choose a status from the list.</a>`,
		`href="#recipe-ar-date-error">Enter the project date.</a>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("422 form is missing the action-pattern error %q", want)
		}
	}
	for _, banned := range []string{
		"Name is required", "Choose a valid status", "Date is required",
		"is required.", "Invalid input", "is invalid",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("422 form must not contain consequence-only copy %q", banned)
		}
	}
}

// TestRecipeEmptyStatesCarryActionLanguage proves the recipe empty states
// answer "what is here / why it is empty / what to do next" with real action
// language, never a bare "No data" (content style guide §Empty states).
func TestRecipeEmptyStatesCarryActionLanguage(t *testing.T) {
	t.Run("admin search no match", func(t *testing.T) {
		resetRecipeResourceStore()
		body := getOKBody(t, "/recipes/admin-resource?q=zzzznothing")
		for _, want := range []string{`No results`, "Try adjusting the filters."} {
			if !strings.Contains(body, want) {
				t.Errorf("admin search empty state is missing %q", want)
			}
		}
	})

	t.Run("admin empty dataset", func(t *testing.T) {
		resetRecipeResourceStore()
		for _, it := range resourceDemoStore.snapshot() {
			resourceDemoStore.delete(it.ID)
		}
		body := getOKBody(t, "/recipes/admin-resource")
		for _, want := range []string{`No projects yet`, "Create your first project to get started."} {
			if !strings.Contains(body, want) {
				t.Errorf("admin empty dataset is missing %q", want)
			}
		}
	})

	t.Run("feed following view", func(t *testing.T) {
		resetRecipeFeedStore()
		for _, it := range feedDemoStore.snapshot() {
			feedDemoStore.delete(it.ID)
		}
		body := getOKBody(t, "/recipes/public-feed?view=following")
		for _, want := range []string{`No posts from people you follow`, "Follow more people to fill this feed."} {
			if !strings.Contains(body, want) {
				t.Errorf("following empty state is missing %q", want)
			}
		}
	})
}

// componentContentSlugs mirrors web/content_name_that_ui_test.go: every served
// component page under web/content. Handbook pages (handbook-*), the docs root
// (index.md) and the design principles page (principles.md) are NOT components
// and are excluded: the length contract applies to components only.
var componentContentSlugs = []string{
	"badge", "button", "card", "checkbox", "chips", "data-table", "dialog",
	"divider", "elevation", "fab", "focus-ring", "icon", "icon-button", "list",
	"menu", "navigation-bar", "navigation-drawer", "navigation-tab", "progress",
	"radio", "segmented-button", "select", "slider", "switch", "tabs",
	"text-field", "toast", "tooltip",
}

// sentenceSplitter splits prose on sentence-final punctuation followed by
// whitespace, per the copy contract (content style guide §Paragraphs and
// sentences: no sentence over 25 words).
var sentenceSplitter = regexp.MustCompile(`[.!?]\s+`)

// componentProse returns the prose sentences of a component page, with code
// excluded: fenced code blocks are dropped, inline `code` spans are replaced
// by a space (they inflate word counts), and table rows (lines starting with
// "|") are skipped entirely.
func componentProse(t *testing.T, slug string) []string {
	t.Helper()
	source, err := fs.ReadFile(webassets.Assets, "content/"+slug+".md")
	if err != nil {
		t.Fatalf("read content/%s.md: %v", slug, err)
	}
	var kept []string
	inFence := false
	for _, line := range strings.Split(string(source), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence || strings.HasPrefix(trimmed, "|") {
			continue
		}
		kept = append(kept, line)
	}
	doc := strings.Join(kept, "\n")
	doc = regexp.MustCompile("`[^`]*`").ReplaceAllString(doc, " ")
	var out []string
	for _, part := range sentenceSplitter.Split(doc, -1) {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// TestComponentPagesKeepSentencesUnder25Words is the copy length contract:
// every sentence in every served component page is at most 25 words, with 20
// as the target (content style guide §Paragraphs and sentences). Readers scan
// (NNG F-pattern); long sentences bury the answer, so the contract is
// mechanical: split on [.!?] + whitespace, count words, fail on > 25.
func TestComponentPagesKeepSentencesUnder25Words(t *testing.T) {
	offenders := 0
	for _, slug := range componentContentSlugs {
		for _, sentence := range componentProse(t, slug) {
			if n := len(strings.Fields(sentence)); n > 25 {
				offenders++
				t.Errorf("%s.md has a %d-word sentence (> 25): %q", slug, n, sentence)
			}
		}
	}
	if offenders > 0 {
		t.Fatalf("found %d sentences over 25 words across component pages", offenders)
	}
}
