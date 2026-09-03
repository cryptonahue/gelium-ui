package app

import (
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

// sectionRefType is a closed filter vocabulary for GET /docs/section-references?type=.
var sectionRefTypes = []string{"article", "hero", "pricing", "auth", "faq", "404"}

// sectionRefEntry is one published or planned ficha. Only Published rows appear
// in the index. ContentPath is the embedded markdown for the detail route.
type sectionRefEntry struct {
	ID          string
	Type        string
	Title       string
	Cited       string
	Summary     string
	Path        string
	ContentPath string
	Published   bool
}

// sectionRefCatalog is the first-party collection. Add a row + markdown file +
// tests when shipping the next ficha. Do not publish without a live remake route.
var sectionRefCatalog = []sectionRefEntry{
	{
		ID:          "article",
		Type:        "article",
		Title:       "Article / rich post",
		Cited:       "Vercel",
		Summary:     "Header, prose, table, related — on-page Gelium reconstruction",
		Path:        "/docs/section-references/article",
		ContentPath: "content/handbook-section-references-article.md",
		Published:   true,
	},
	{
		ID:          "404-vercel",
		Type:        "404",
		Title:       "404 / missing route",
		Cited:       "Vercel",
		Summary:     "Minimal not-found state with a clear path back",
		Path:        "/docs/section-references/404-vercel",
		ContentPath: "content/handbook-section-references-404-vercel.md",
		Published:   true,
	},
	{
		ID:          "auth-register",
		Type:        "auth",
		Title:       "Auth / register",
		Cited:       "Vercel",
		Summary:     "Registration form with consent, validation, and recovery",
		Path:        "/docs/section-references/auth-register",
		ContentPath: "content/handbook-section-references-auth-register.md",
		Published:   true,
	},
	{
		ID:          "faq-vercel",
		Type:        "faq",
		Title:       "FAQ / pricing questions",
		Cited:       "Vercel",
		Summary:     "Native disclosures for pricing uncertainty and decisions",
		Path:        "/docs/section-references/faq-vercel",
		ContentPath: "content/handbook-section-references-faq-vercel.md",
		Published:   true,
	},
	{
		ID:          "hero-linear",
		Type:        "hero",
		Title:       "Hero / product direction",
		Cited:       "Linear",
		Summary:     "Promise, context, and one primary landing action",
		Path:        "/docs/section-references/hero-linear",
		ContentPath: "content/handbook-section-references-hero-linear.md",
		Published:   true,
	},
	{
		ID:          "pricing-linear",
		Type:        "pricing",
		Title:       "Pricing / plan comparison",
		Cited:       "Linear",
		Summary:     "Plan ladder and mobile-safe feature comparison",
		Path:        "/docs/section-references/pricing-linear",
		ContentPath: "content/handbook-section-references-pricing-linear.md",
		Published:   true,
	},
}

func closedSectionRefType(raw string) string {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" || v == "all" {
		return "all"
	}
	for _, t := range sectionRefTypes {
		if v == t {
			return t
		}
	}
	return "all"
}

func lookupSectionRef(id string) (sectionRefEntry, bool) {
	for _, e := range sectionRefCatalog {
		if e.ID == id && e.Published {
			return e, true
		}
	}
	return sectionRefEntry{}, false
}

func publishedSectionRefPaths() []string {
	out := make([]string, 0, len(sectionRefCatalog))
	for _, e := range sectionRefCatalog {
		if e.Published {
			out = append(out, e.Path)
		}
	}
	return out
}

func sectionRefCatalogMarkdown(selected string) string {
	var b strings.Builder
	b.WriteString("\n## Browse by type\n\n")
	b.WriteString(`<form method="get" action="/docs/section-references" aria-label="Filter section references">`)
	b.WriteString("\n<div class=\"ui-select ui-select-outlined\">\n")
	b.WriteString(`<select id="section-ref-type" name="type">`)
	b.WriteByte('\n')
	writeTypeOption(&b, "all", "all", selected)
	for _, t := range sectionRefTypes {
		writeTypeOption(&b, t, t, selected)
	}
	b.WriteString("</select>\n")
	b.WriteString(`<label class="ui-select-label" for="section-ref-type">Type</label>`)
	b.WriteString("\n<span class=\"ui-select-caret\" aria-hidden=\"true\">")
	b.WriteString(string(iconSVG("arrow_downward")))
	b.WriteString("</span>\n</div>\n")
	b.WriteString(`<button type="submit" class="ui-button ui-button-outline"><span>Apply</span></button>`)
	b.WriteString("\n</form>\n\n")
	b.WriteString(sectionRefResultsMarkdown(selected, sectionRefEntriesFor(selected)))
	return b.String()
}

func sectionRefEntriesFor(selected string) []sectionRefEntry {
	entries := make([]sectionRefEntry, 0, len(sectionRefCatalog))
	for _, e := range sectionRefCatalog {
		if e.Published && (selected == "all" || e.Type == selected) {
			entries = append(entries, e)
		}
	}
	return entries
}

func sectionRefResultsMarkdown(selected string, entries []sectionRefEntry) string {
	var b strings.Builder
	b.WriteString("## Results\n\n")
	b.WriteString(`<p id="section-ref-results-status" role="status">Showing `)
	b.WriteString(strconv.Itoa(len(entries)))
	if len(entries) == 1 {
		b.WriteString(" result")
	} else {
		b.WriteString(" results")
	}
	if selected != "all" {
		b.WriteString(" for ")
		b.WriteString(selected)
	}
	b.WriteString(".</p>\n\n")
	if len(entries) == 0 {
		b.WriteString(`<div class="ui-empty-state ui-empty-state--compact" role="status">`)
		b.WriteString("\n<p class=\"ui-empty-state-title\">No references found</p>\n")
		b.WriteString("<p class=\"ui-empty-state-body\">There are no published references for this type yet.</p>\n")
		b.WriteString(`<a class="ui-button ui-button-outline" href="/docs/section-references">Clear filter</a>`)
		b.WriteString("\n</div>\n")
		return b.String()
	}

	b.WriteString(`<nav aria-label="Section reference results">`)
	b.WriteString("\n<ul class=\"ui-list\">\n")
	for _, e := range entries {
		b.WriteString("<li class=\"ui-list-item ui-list-item--two-line\">\n")
		b.WriteString(`<a class="ui-list-item-link" href="`)
		b.WriteString(e.Path)
		b.WriteString(`">`)
		b.WriteString("\n<span class=\"ui-list-item-text\">\n")
		b.WriteString(`<span class="ui-list-item-headline">`)
		b.WriteString(e.Title)
		b.WriteString("</span>\n")
		b.WriteString(`<span class="ui-list-item-supporting">`)
		b.WriteString(e.Type)
		b.WriteString(" · ")
		b.WriteString(e.Cited)
		b.WriteString(" cited · ")
		b.WriteString(e.Summary)
		b.WriteString("</span>\n</span>\n</a>\n")
		b.WriteString(`<span class="ui-list-item-icon ui-list-item-icon--end" aria-hidden="true">`)
		b.WriteString(string(iconSVG("chevron_right")))
		b.WriteString("</span>\n</li>\n")
	}
	b.WriteString("</ul>\n</nav>\n")
	return b.String()
}

func writeTypeOption(b *strings.Builder, value, label, selected string) {
	b.WriteString(`<option value="`)
	b.WriteString(value)
	b.WriteString(`"`)
	if selected == value {
		b.WriteString(" selected")
	}
	b.WriteString(">")
	b.WriteString(label)
	b.WriteString("</option>\n")
}

// docsSectionReferences is GET /docs/section-references — catalog index with a
// closed GET ?type= filter (0-JS). Unknown type values fall back to all.
func (s *server) docsSectionReferences(w http.ResponseWriter, r *http.Request) {
	source, err := fs.ReadFile(s.assets, "content/handbook-section-references.md")
	if err != nil {
		s.renderErrorPage(w, r, http.StatusInternalServerError, "Something went wrong", "This page could not be loaded. Please try again later.", true, "/", "Back to home", "/docs/section-references")
		return
	}
	selected := closedSectionRefType(r.URL.Query().Get("type"))
	md := string(source) + sectionRefCatalogMarkdown(selected)
	s.renderMarkdown(w, r, pageView{Title: "Section references"}, md, "/docs/section-references")
}

// docsSectionReferenceDetail is GET /docs/section-references/{id}. Unknown or
// unpublished ids render the docs 404.
func (s *server) docsSectionReferenceDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, ok := lookupSectionRef(id)
	if !ok {
		s.notFound(w, r)
		return
	}
	s.renderMarkdownPageAt(w, r, pageView{Title: e.Title}, e.ContentPath, e.Path)
}
