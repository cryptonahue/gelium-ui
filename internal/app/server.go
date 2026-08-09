package app

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/yuin/goldmark"

	webassets "loomui/web"
)

type pageView struct {
	Title          string
	Content        template.HTML
	Nav            []navLink
	CTA            *buttonView
	Buttons        []buttonView
	TextFields     []textFieldView
	ValidationForm *validationFormView
	Dialog         *dialogView
	Toasts         []toastView
	ToastDemo      *toastDemoView
	ElevationDemo  *elevationDemo
	FocusRingDemo  *focusRingDemo
	IconDemo       *iconDemo
	DividerDemo    *dividerDemo
	CardDemo       *cardDemo
	BadgeDemo      *badgeDemo
	CheckboxDemo   *checkboxDemo
	RadioDemo      *radioDemo
	SwitchDemo     *switchDemo
	SelectDemo     *selectDemo
	SelectMenuDemo *selectMenuDemo
	SliderDemo     *sliderDemo
	ProgressDemo   *progressDemo
	IconButtonDemo *iconButtonDemo
	FabDemo        *fabDemo
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
	for _, r := range componentRoutes() {
		r := r
		mux.HandleFunc("GET "+r.Path, func(w http.ResponseWriter, req *http.Request) {
			r.Handler(s, w, req)
		})
	}
	mux.HandleFunc("POST /examples/text-field/validate", s.validateTextField)
	mux.HandleFunc("POST /examples/toast/demo", s.toastDemo)
	mux.HandleFunc("POST /examples/select/menu", s.selectMenu)
	mux.HandleFunc("GET /static/{name}", s.staticAsset)
	return mux
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
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
		Title: "Loom UI",
		CTA: &buttonView{
			Label:   "Explore Button",
			Variant: "primary",
			Href:    "/components/button",
		},
	}, "content/index.md")
}

func (s *server) renderMarkdownPage(w http.ResponseWriter, data pageView, contentPath string) {
	s.renderMarkdownPageStatus(w, data, contentPath, http.StatusOK)
}

func (s *server) renderMarkdownPageStatus(w http.ResponseWriter, data pageView, contentPath string, status int) {
	source, err := fs.ReadFile(s.assets, contentPath)
	if err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	var rendered bytes.Buffer
	if err := s.markdown.Convert(source, &rendered); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	var page bytes.Buffer
	data.Nav = navLinks()
	data.Content = template.HTML(rendered.String()) // #nosec G203 -- source is trusted, embedded Markdown.
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(page.Bytes())
}
