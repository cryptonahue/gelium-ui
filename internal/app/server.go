package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/yuin/goldmark"

	webassets "loomui/web"
)

// defaultThemeClass is the theme applied when none is requested. The value must
// match a class owned by a theme that ships on disk (themes/*/theme.css).
const defaultThemeClass = "theme-material"

// themeClass resolves a requested theme to a safe CSS class. Theme identity is
// server-driven and validated against an allowlist of themes that exist on
// disk; unknown values fall back to the default so no page can inject an
// arbitrary class or depend on a theme that does not exist. Extend the
// allowlist as new themes (e.g. theme-basecoat) land.
func themeClass(theme string) string {
	for _, allowed := range []string{defaultThemeClass} {
		if theme == allowed {
			return theme
		}
	}
	return defaultThemeClass
}

// bannerView is the server-driven view model for the page-level BANNER slot
// (Phase D pattern 5). The partial is a primitive: no handler wires it yet —
// Phase G will inject it when a page/site condition (expired session,
// maintenance, pending consent) is detected. A nil Banner on pageView renders
// no banner. Tone is a closed vocabulary (info|success|warning|error), the
// CTA is a real link and the dismiss a POST+303 action, both 0 JS.
type bannerView struct {
	Tone        string
	Icon        template.HTML
	Title       string
	Body        string
	CTA         bool
	CTAHref     string
	CTALabel    string
	DismissHref string
	DismissIcon template.HTML
}

// inlineAlertView is the server-driven view model for the section/form-level
// INLINE ALERT slot (Phase D pattern 6). A nil InlineAlert on pageView renders
// no alert. Tone is a closed vocabulary (info|success|warning|error) and the
// role derives from it (error → alert, everything else → status). It is a
// persistent server-rendered signal, never a transient toast.
type inlineAlertView struct {
	Tone  string
	Icon  template.HTML
	Title string
	Body  string
}

// errorStateView is the server-driven view model for the page-level ERROR
// STATE slot (Phase D pattern 7). A nil Error on pageView renders no error
// state and the normal content instead. The status code is the canonical
// attribute: the handler picks the real HTTP status (404/500) and the copy
// per status. Retry is optional — a real GET link back to a known URL —
// and everything works with 0 JS.
type errorStateView struct {
	StatusCode int
	Title      string
	Body       string
	Retry      bool
	Href       string
	Label      string
}

// siteBaseURL is the canonical origin for absolute URLs (canonical, og:url,
// JSON-LD). Gelium UI ships with no production host yet — the docs are
// embedded and self-hosted — so this is a documented placeholder origin
// matching the SEO contract examples. Canonical URLs are built from it plus
// the clean route path and never from the request Host, so metadata stays
// stable across every deployment. Swap this single constant for the real
// origin when a public domain exists.
const siteBaseURL = "https://gelium-ui.example"

const (
	// homeDescription is the system description for the landing page.
	homeDescription = "Gelium UI — themeable, open-code UI components for server-rendered applications. Native HTML semantics, zero component JavaScript, Material 3 design."
	// docsDescription is the system description for the /docs index.
	docsDescription = "Gelium UI component library, organized by category. Every page is dogfooded: it renders the real component it documents."
	// defaultMetaDescription is the fallback description for any route without
	// its own. Component pages derive a per-route description from their label
	// so every URL stays unique (contract: one description per URL).
	defaultMetaDescription = "Gelium UI — open-code, server-rendered UI components with native HTML semantics."
)

// metaView is the server-driven metadata model emitted into <head>. It is
// theme-agnostic by contract: switching themes must not change a single byte
// of SEO-relevant markup. A zero metaView renders a valid minimal head (the
// layout guards every tag behind {{if}}), so error pages and partials keep
// working without metadata.
type metaView struct {
	Title         string // mirrors pageView.Title (single source in the handler)
	Description   string
	Canonical     string
	Robots        string // default "index, follow"
	OGTitle       string
	OGDescription string
	OGType        string
	JSONLD        template.JS // trusted structured data; emitted inside <script type="application/ld+json">
	Lang          string      // default "en"
}

// jsonLDPublisher and webSiteLD are the home-page WebSite structured-data
// types. They are marshaled with encoding/json so escaping and validity are
// guaranteed (contract §12: no JSON built by string concatenation).
type jsonLDPublisher struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type webSiteLD struct {
	Context    string          `json:"@context"`
	Type       string          `json:"@type"`
	Name       string          `json:"name"`
	URL        string          `json:"url"`
	InLanguage string          `json:"inLanguage"`
	Publisher  jsonLDPublisher `json:"publisher"`
}

// websiteJSONLD is the home-page WebSite block, built once at startup. The
// value is template.JS on purpose: inside a <script type="application/ld+json">
// element, html/template escapes template.HTML as a JSON string, while
// template.JS (a valid JS expression; JSON is valid JS) is emitted verbatim.
var websiteJSONLD = func() template.JS {
	b, err := json.Marshal(webSiteLD{
		Context:    "https://schema.org",
		Type:       "WebSite",
		Name:       "Gelium UI",
		URL:        siteBaseURL + "/",
		InLanguage: "en",
		Publisher:  jsonLDPublisher{Type: "Organization", Name: "Gelium UI"},
	})
	if err != nil {
		return template.JS("")
	}
	return template.JS(b) // #nosec G203 -- trusted, system-generated JSON.
}()

// routePathForContent maps an embedded markdown file to its public route path.
// Content files live at content/<route-segment>.md and render at
// /components/<route-segment>; content/index.md renders at /. The convention
// is what keeps the canonical URL clean and query-free.
func routePathForContent(contentPath string) string {
	base := path.Base(contentPath)
	if base == "index.md" {
		return "/"
	}
	return "/components/" + strings.TrimSuffix(base, ".md")
}

// componentLabel resolves the navigation label for a /components/* route from
// the single route registry, so metadata never drifts from the nav.
func componentLabel(routePath string) string {
	for _, r := range componentRoutes() {
		if r.Path == routePath {
			return r.Label
		}
	}
	return routePath
}

// componentDescription derives a unique per-route description for a component
// page from its registry label (contract: one description per URL).
func componentDescription(label string) string {
	return "The " + label + " component in Gelium UI — server-rendered docs with native HTML semantics, zero component JavaScript."
}

// resolveMeta computes the page metadata for a route path at the single render
// choke point. Fields already set by the handler (data.Meta) win; resolution
// only fills what is still empty. Demos and examples are never indexed; every
// other route defaults to "index, follow". The canonical is always the clean
// route path — no query string can leak into it (contract §16).
func resolveMeta(data pageView, routePath string) metaView {
	meta := data.Meta
	meta.Title = data.Title
	if meta.Lang == "" {
		meta.Lang = "en"
	}
	if meta.Robots == "" {
		meta.Robots = "index, follow"
	}
	if strings.HasPrefix(routePath, "/demo/") || strings.HasPrefix(routePath, "/examples/") {
		meta.Robots = "noindex, nofollow"
	}

	meta.Canonical = siteBaseURL + routePath

	if meta.Description == "" {
		switch {
		case routePath == "/":
			meta.Description = homeDescription
		case routePath == "/docs":
			meta.Description = docsDescription
		case strings.HasPrefix(routePath, "/components/"):
			meta.Description = componentDescription(componentLabel(routePath))
		default:
			meta.Description = defaultMetaDescription
		}
	}
	if meta.OGType == "" {
		switch {
		case routePath == "/" || routePath == "/docs":
			meta.OGType = "website"
		case strings.HasPrefix(routePath, "/components/"):
			meta.OGType = "article"
		default:
			meta.OGType = "website"
		}
	}
	if meta.OGTitle == "" {
		if routePath == "/" {
			meta.OGTitle = "Gelium UI"
		} else {
			meta.OGTitle = meta.Title + " · Gelium UI"
		}
	}
	if meta.OGDescription == "" {
		meta.OGDescription = meta.Description
	}
	if routePath == "/" && meta.JSONLD == "" {
		meta.JSONLD = websiteJSONLD
	}
	return meta
}

type pageView struct {
	Meta                 metaView
	Title                string
	Content              template.HTML
	ThemeClass           string
	Nav                  []navLink
	Banner               *bannerView
	Error                *errorStateView
	InlineAlert          *inlineAlertView
	CTA                  *buttonView
	Buttons              []buttonView
	TextFields           []textFieldView
	ValidationForm       *validationFormView
	Dialog               *dialogView
	DialogConfirm        *dialogConfirmView
	Toasts               []toastView
	ToastDemo            *toastDemoView
	ElevationDemo        *elevationDemo
	FocusRingDemo        *focusRingDemo
	IconDemo             *iconDemo
	DividerDemo          *dividerDemo
	CardDemo             *cardDemo
	BadgeDemo            *badgeDemo
	CheckboxDemo         *checkboxDemo
	RadioDemo            *radioDemo
	SwitchDemo           *switchDemo
	SelectDemo           *selectDemo
	SelectMenuDemo       *selectMenuDemo
	SliderDemo           *sliderDemo
	ProgressDemo         *progressDemo
	IconButtonDemo       *iconButtonDemo
	FabDemo              *fabDemo
	ListDemo             *listDemo
	ChipsDemo            *chipsDemo
	TabsDemo             *tabsDemo
	NavigationBarDemo    *navigationBarDemo
	NavigationTabDemo    *navigationTabDemo
	SegmentedButtonDemo  *segmentedButtonDemo
	MenuDemo             *menuDemo
	NavigationDrawerDemo *navigationDrawerDemo
	DataTableDemo        *dataTableDemo
	TooltipDemo          *tooltipDemo
}

type server struct {
	templates *template.Template
	markdown  goldmark.Markdown
	assets    fs.FS
}

// New builds the Gelium UI documentation HTTP handler from embedded assets.
func New() http.Handler {
	templates := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	s := &server{
		templates: templates,
		markdown:  goldmark.New(),
		assets:    webassets.Assets,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /docs", s.docsIndex)
	for _, r := range componentRoutes() {
		r := r
		mux.HandleFunc("GET "+r.Path, func(w http.ResponseWriter, req *http.Request) {
			r.Handler(s, w, req)
		})
	}
	mux.HandleFunc("POST /examples/text-field/validate", s.validateTextField)
	mux.HandleFunc("POST /examples/toast/demo", s.toastDemo)
	mux.HandleFunc("POST /examples/select/menu", s.selectMenu)
	mux.HandleFunc("GET /components/dialog/confirm", s.dialogConfirm)
	mux.HandleFunc("POST /components/dialog/confirm", s.dialogConfirmPost)
	mux.HandleFunc("POST /examples/chips/remove", s.chipsRemoveDemo)
	mux.HandleFunc("POST /examples/data-table/refresh", s.dataTableRefreshDemo)
	mux.HandleFunc("GET /demo/whatsapp", s.whatsAppDemo)
	mux.HandleFunc("GET /demo/whatsapp/admin", s.whatsAppAdmin)
	mux.HandleFunc("POST /demo/whatsapp/send", s.whatsAppSend)
	mux.HandleFunc("POST /demo/whatsapp/send-template", s.whatsAppSendTemplate)
	mux.HandleFunc("POST /demo/whatsapp/typing", s.whatsAppTyping)
	mux.HandleFunc("POST /demo/whatsapp/read", s.whatsAppRead)
	mux.HandleFunc("POST /demo/whatsapp/admin/webhook", s.whatsAppWebhookSave)
	mux.HandleFunc("GET /static/{name}", s.staticAsset)
	// 404 catch-all: any unknown GET path falls back to the styled ERROR STATE
	// page (the mux gives the more specific patterns above priority). Post-only
	// routes register a GET 405 companion below so a GET to them keeps the
	// mux's modern method-mismatch semantics instead of being swallowed by the
	// catch-all.
	for _, path := range postOnlyPaths() {
		mux.HandleFunc("GET "+path, methodNotAllowed)
	}
	mux.HandleFunc("GET /{path...}", s.notFound)
	return mux
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

// postOnlyPaths are the routes that only accept POST. A GET to any of them
// must stay a 405 (the mux's modern method-mismatch semantics) even after the
// catch-all 404, so each path also registers a GET methodNotAllowed companion.
func postOnlyPaths() []string {
	return []string{
		"/examples/text-field/validate",
		"/examples/toast/demo",
		"/examples/select/menu",
		"/examples/chips/remove",
		"/examples/data-table/refresh",
		"/demo/whatsapp/send",
		"/demo/whatsapp/send-template",
		"/demo/whatsapp/typing",
		"/demo/whatsapp/read",
		"/demo/whatsapp/admin/webhook",
	}
}

// methodNotAllowed answers a GET to a POST-only route with the net/http 405
// Method Not Allowed response, preserving the Allow header the ServeMux would
// otherwise emit.
func methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Allow", "POST")
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

// notFound serves the styled ERROR STATE page for unknown routes, replacing
// the plain-text net/http default. It renders the full docs layout with the
// Error slot set and the real 404 status.
func (s *server) notFound(w http.ResponseWriter, _ *http.Request) {
	s.renderErrorPage(w,
		http.StatusNotFound,
		"Page not found",
		"The page you are looking for does not exist or has moved.",
		true,
		"/",
		"Back to home",
	)
}

// renderErrorPage renders the full docs layout with the ERROR STATE slot set
// and the real HTTP status. It is the page-level failure path for the 404
// catch-all and resource 500s; the layout template is already parsed, so a
// failed template exec still falls back to a minimal plain response.
func (s *server) renderErrorPage(w http.ResponseWriter, status int, title, body string, retry bool, href, label string) {
	var page bytes.Buffer
	data := pageView{
		Title:      title,
		ThemeClass: themeClass(""),
		Nav:        navLinks(),
		Error: &errorStateView{
			StatusCode: status,
			Title:      title,
			Body:       body,
			Retry:      retry,
			Href:       href,
			Label:      label,
		},
	}
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(page.Bytes())
}

func (s *server) staticAsset(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	contentTypes := map[string]string{
		"app.css":     "text/css; charset=utf-8",
		"htmx.min.js": "text/javascript; charset=utf-8",
		"app.js":      "text/javascript; charset=utf-8",
	}
	contentType, ok := contentTypes[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	asset, err := fs.ReadFile(s.assets, "static/"+name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(asset)
}

func (s *server) home(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Gelium UI",
		CTA: &buttonView{
			Label:   "Read the docs",
			Variant: "primary",
			Href:    "/docs",
		},
	}, "content/index.md")
}

func (s *server) renderMarkdownPage(w http.ResponseWriter, data pageView, contentPath string) {
	s.renderMarkdownPageStatus(w, data, contentPath, http.StatusOK)
}

func (s *server) renderMarkdownPageStatus(w http.ResponseWriter, data pageView, contentPath string, status int) {
	source, err := fs.ReadFile(s.assets, contentPath)
	if err != nil {
		s.renderErrorPage(w, http.StatusInternalServerError, "Something went wrong", "This page could not be loaded. Please try again later.", true, "/", "Back to home")
		return
	}
	s.renderMarkdownStatus(w, data, string(source), routePathForContent(contentPath), status)
}

// renderMarkdown converts an in-memory Markdown string and renders it into the
// docs layout. Used by pages that build their content programmatically (the
// /docs index) as well as by pages served from embedded files. routePath is
// the public route identity used to resolve metadata.
func (s *server) renderMarkdown(w http.ResponseWriter, data pageView, source, routePath string) {
	s.renderMarkdownStatus(w, data, source, routePath, http.StatusOK)
}

func (s *server) renderMarkdownStatus(w http.ResponseWriter, data pageView, source, routePath string, status int) {
	var rendered bytes.Buffer
	if err := s.markdown.Convert([]byte(source), &rendered); err != nil {
		s.renderErrorPage(w, http.StatusInternalServerError, "Something went wrong", "This page could not be rendered. Please try again later.", true, "/", "Back to home")
		return
	}

	var page bytes.Buffer
	data.Nav = navLinks()
	data.Meta = resolveMeta(data, routePath)
	data.Content = template.HTML(rendered.String()) // #nosec G203 -- markdown is trusted (embedded or generated).
	data.ThemeClass = themeClass(data.ThemeClass)
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(page.Bytes())
}
