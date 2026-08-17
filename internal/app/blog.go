package app

import (
	"bytes"
	"html/template"
	"io/fs"
	"math"
	"net/http"
	"strings"

	"geliumui/lib"
)

// blogPostSummary is one card on the blog index: the metadata that makes a
// post browsable before opening it.
type blogPostSummary struct {
	Slug    string
	Title   string
	Date    string // ISO-8601 (2006-01-02), mirrors the docs date table format
	Summary string
}

// blogPostView is the rendered single-post page: header metadata plus the
// long-form markdown content rendered by the server.
type blogPostView struct {
	Title       string
	Date        string
	Author      string
	ReadingTime int // minutes, derived from the source word count
	Content     template.HTML
}

// blogView is the blog space composition. The index sets Posts; post pages
// set Post. The layout renders the blog template when Blog is non-nil.
type blogView struct {
	Title string
	Intro string
	Posts []blogPostSummary
	Post  *blogPostView
}

// blogPost is one registry entry for a blog post. Slug is the URL segment
// (GET /blog/<slug>), ContentPath points at the embedded markdown file.
type blogPost struct {
	Slug        string
	Title       string
	Date        string
	Author      string
	Summary     string
	ContentPath string
}

// blogPosts is the ordered blog registry: registration order is the card
// order on the index and the lookup table for /blog/<slug>. The registry is
// the single source of post metadata — the content files carry no frontmatter.
var blogPosts = []blogPost{
	{
		Slug:        "genesis",
		Title:       "How Gelium UI was born",
		Date:        "2026-08-15",
		Author:      "Gelium UI team",
		Summary:     "Gelium UI started as a side project and grew into a contract-tested component library for server-rendered Go apps.",
		ContentPath: "content/blog-genesis.md",
	},
	{
		Slug:        "roadmap",
		Title:       "Where Gelium UI is going",
		Date:        "2026-08-15",
		Author:      "Gelium UI team",
		Summary:     "Phases A–J shipped with contract tests. Docs and DX are done. Next: truth sync, demos, SEO, and a release.",
		ContentPath: "content/blog-roadmap.md",
	},
}

// readingTimeMinutesFor estimates reading time from markdown source: one
// minute per 200 words, rounded up, at least one minute. Pure function.
func readingTimeMinutesFor(source string) int {
	words := len(strings.Fields(source))
	return max(1, int(math.Ceil(float64(words)/200)))
}

// blogIndex is GET /blog — the blog space index: every registered post as a
// Gelium card grid (dogfooding ui-card markup on the blog's own surface).
func (s *server) blogIndex(w http.ResponseWriter, r *http.Request) {
	posts := make([]blogPostSummary, 0, len(blogPosts))
	for _, p := range blogPosts {
		posts = append(posts, blogPostSummary{
			Slug:    p.Slug,
			Title:   p.Title,
			Date:    p.Date,
			Summary: p.Summary,
		})
	}
	data := pageView{
		Title: "Gelium blog",
		Meta:  metaView{Description: "The Gelium UI blog — how the library is built, what shipped, and where it is going."},
		Blog: &blogView{
			Title: "Gelium blog",
			Intro: "The Gelium UI blog — how the library is built, what shipped, and where it is going.",
			Posts: posts,
		},
	}
	s.renderBlog(w, r, data)
}

// blogPost is GET /blog/{slug} — one post rendered from its embedded
// markdown. Slugs outside the registry get the styled 404, never a ghost page.
func (s *server) blogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	var entry *blogPost
	for i := range blogPosts {
		if blogPosts[i].Slug == slug {
			entry = &blogPosts[i]
			break
		}
	}
	if entry == nil {
		s.renderErrorPage(w, r, http.StatusNotFound, "Post not found", "This blog post does not exist. Browse the blog index for published posts.", true, "/blog", "Back to blog", "/blog/"+slug)
		return
	}
	source, err := fs.ReadFile(s.assets, entry.ContentPath)
	if err != nil {
		s.renderErrorPage(w, r, http.StatusInternalServerError, "Something went wrong", "This post could not be loaded. Please try again later.", true, "/blog", "Back to blog", "/blog/"+slug)
		return
	}
	// The post header renders the title; drop the file's leading H1 so the
	// page keeps exactly one heading level 1 (same helper as the docs root).
	var rendered bytes.Buffer
	if err := s.markdown.Convert([]byte(stripDocsRootH1(string(source))), &rendered); err != nil {
		s.renderErrorPage(w, r, http.StatusInternalServerError, "Something went wrong", "This post could not be rendered. Please try again later.", true, "/blog", "Back to blog", "/blog/"+slug)
		return
	}
	data := pageView{
		Title: entry.Title,
		Meta:  metaView{Description: entry.Summary},
		Blog: &blogView{
			Post: &blogPostView{
				Title:       entry.Title,
				Date:        entry.Date,
				Author:      entry.Author,
				ReadingTime: readingTimeMinutesFor(string(source)),
				Content:     template.HTML(rendered.String()), // #nosec G203 -- markdown is embedded and trusted.
			},
		},
	}
	s.renderBlog(w, r, data)
}

// renderBlog renders the blog space (index or post) through the shared layout
// without the Markdown prose wrapper. The blog is a separate space like the
// landing: legacy site-header chrome with theme/scheme switchers, no docs
// shell and no sidebar. Meta resolves against the real request path so
// canonical URLs stay clean per route.
func (s *server) renderBlog(w http.ResponseWriter, r *http.Request, data pageView) {
	routePath := requestPath(r)
	data.Meta = resolveMeta(data, routePath)
	data.AssetsVersion = lib.AssetsVersion
	if data.Footer == nil {
		data.Footer = defaultFooter()
	}
	data.Nav = homeLandingNav()

	themeSlug := ""
	if q := themeFromRequest(r); q != "" {
		data.ThemeClass = q
		themeSlug = themeSlugFromClass(q)
	} else {
		data.ThemeClass = themeClass(data.ThemeClass)
	}
	scheme := schemeFromRequest(r)
	applyDocumentRootScheme(&data, scheme)

	// Same 0-JS theme + appearance controls as the docs, in the site header.
	data.ThemeSwitcher = themeSwitcherFor(r, data.ThemeClass, themeSlug, scheme)
	data.SchemeSwitcher = schemeSwitcherFor(r, themeSlug, scheme)

	// Rewrite compact nav hrefs so theme/scheme survive header clicks.
	if themeSlug != "" || normalizeScheme(scheme) != "" {
		for i := range data.Nav {
			data.Nav[i].Path = chromeHref(data.Nav[i].Path, themeSlug, scheme)
		}
	}

	var page bytes.Buffer
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(page.Bytes())
}
