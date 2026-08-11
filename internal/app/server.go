package app

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"

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

type pageView struct {
	Title                string
	Content              template.HTML
	ThemeClass           string
	Nav                  []navLink
	Banner               *bannerView
	Error                *errorStateView
	CTA                  *buttonView
	Buttons              []buttonView
	TextFields           []textFieldView
	ValidationForm       *validationFormView
	Dialog               *dialogView
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

// New builds the Loom UI documentation HTTP handler from embedded assets.
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
	mux.HandleFunc("POST /examples/chips/remove", s.chipsRemoveDemo)
	mux.HandleFunc("POST /examples/data-table/refresh", s.dataTableRefreshDemo)
	mux.HandleFunc("GET /demo/whatsapp", s.whatsAppDemo)
	mux.HandleFunc("GET /demo/whatsapp/admin", s.whatsAppAdmin)
	mux.HandleFunc("POST /demo/whatsapp/send", s.whatsAppSend)
	mux.HandleFunc("POST /demo/whatsapp/send-template", s.whatsAppSendTemplate)
	mux.HandleFunc("POST /demo/whatsapp/typing", s.whatsAppTyping)
	mux.HandleFunc("POST /demo/whatsapp/read", s.whatsAppRead)
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
		Title: "Gelidium UI",
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
	s.renderMarkdownStatus(w, data, string(source), status)
}

// renderMarkdown converts an in-memory Markdown string and renders it into the
// docs layout. Used by pages that build their content programmatically (the
// /docs index) as well as by pages served from embedded files.
func (s *server) renderMarkdown(w http.ResponseWriter, data pageView, source string) {
	s.renderMarkdownStatus(w, data, source, http.StatusOK)
}

func (s *server) renderMarkdownStatus(w http.ResponseWriter, data pageView, source string, status int) {
	var rendered bytes.Buffer
	if err := s.markdown.Convert([]byte(source), &rendered); err != nil {
		s.renderErrorPage(w, http.StatusInternalServerError, "Something went wrong", "This page could not be rendered. Please try again later.", true, "/", "Back to home")
		return
	}

	var page bytes.Buffer
	data.Nav = navLinks()
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
