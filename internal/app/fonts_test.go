package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"geliumui/lib"
)

// TestThemePreloadFontsEmptyByDefault proves the font mechanism is a no-op
// until a theme ships fonts: the reference themes (material, basecoat) emit no
// preload <link> in the <head>, so existing pages are byte-identical apart from
// the (empty) helper call. This is the non-regression guarantee.
func TestThemePreloadFontsEmptyByDefault(t *testing.T) {
	got := string(themePreloadFonts("theme-material", lib.AssetsVersion))
	if got != "" {
		t.Fatalf("theme-material must ship no fonts yet, got %q", got)
	}
	got = string(themePreloadFonts("theme-basecoat", lib.AssetsVersion))
	if got != "" {
		t.Fatalf("theme-basecoat must ship no fonts yet, got %q", got)
	}
	// Unknown class falls back to the default theme (material = no fonts).
	if got := string(themePreloadFonts("theme-nope", lib.AssetsVersion)); got != "" {
		t.Fatalf("unknown theme must fall back to no fonts, got %q", got)
	}
}

// TestThemePreloadFontsRendersLinks injects a temporary theme that ships fonts
// and proves the layout helper renders correct cache-busted <link rel="preload">
// markup. availableThemes is restored after the test.
func TestThemePreloadFontsRendersLinks(t *testing.T) {
	orig := availableThemes
	availableThemes = []themeDirection{
		{
			Class: "theme-material", Slug: "material", Label: "Material",
			Fonts: []themeFont{{File: "theme-material-sans.woff2", Name: "Material Sans", Preload: true}},
		},
		{
			Class: "theme-alden", Slug: "alden", Label: "Alden",
			Fonts: []themeFont{
				{File: "theme-alden-sans.woff2", Name: "Alden Sans", Preload: true},
				{File: "theme-alden-display.woff2", Name: "Alden Display"}, // not preloaded
			},
		},
	}
	t.Cleanup(func() { availableThemes = orig })

	got := string(themePreloadFonts("theme-alden", "9.9.9"))
	// Only the Preload-marked font emits a <link>; the display serif (Preload
	// false) is servable but not preloaded.
	for _, want := range []string{
		`<link rel="preload" as="font" href="/static/fonts/theme-alden-sans.woff2?v=9.9.9" type="font/woff2" crossorigin>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("themePreloadFonts missing %q in %q", want, got)
		}
	}
	if strings.Contains(got, "theme-alden-display.woff2") {
		t.Errorf("themePreloadFonts must not preload non-Preload fonts, got %q", got)
	}
	// Unknown class falls back to the default theme (material) fonts.
	got = string(themePreloadFonts("theme-nope", lib.AssetsVersion))
	if !strings.Contains(got, "theme-material-sans.woff2") {
		t.Errorf("unknown theme must fall back to default theme fonts, got %q", got)
	}
}

// TestFontAssetServesOnlyAllowlistedWoff2 proves the /static/fonts/{file} route
// serves an embedded font with the WOFF2 content type and rejects any filename
// no allowlisted theme ships (closed namespace, no arbitrary reads).
func TestFontAssetServesOnlyAllowlistedWoff2(t *testing.T) {
	orig := availableThemes
	availableThemes = []themeDirection{
		{
			Class: "theme-alden", Slug: "alden", Label: "Alden",
			Fonts: []themeFont{{File: "theme-alden-sans.woff2", Name: "Alden Sans"}},
		},
	}
	t.Cleanup(func() { availableThemes = orig })

	srv := New()

	cases := []struct {
		path   string
		status int
	}{
		{"/static/fonts/theme-alden-sans.woff2", http.StatusNotFound}, // allowlisted filename, missing on disk
		{"/static/fonts/theme-nope.woff2", http.StatusNotFound},        // not allowlisted
		{"/static/fonts/x.jpg", http.StatusNotFound},                   // wrong extension
		{"/static/fonts/", http.StatusNotFound},                        // no file
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		res := httptest.NewRecorder()
		srv.ServeHTTP(res, req)
		if res.Code != c.status {
			t.Errorf("%s = %d, want %d", c.path, res.Code, c.status)
		}
	}
}

// TestFontAssetServesRealAldenWoff2 proves the real shipped Alden font files
// (embedded in lib.Assets via lib/fonts/*) are served with the WOFF2 content
// type through the closed /static/fonts/ namespace, using the real
// availableThemes catalog (no override).
func TestFontAssetServesRealAldenWoff2(t *testing.T) {
	srv := New()
	for _, file := range []string{
		"theme-alden-inter-400-latin.woff2",
		"theme-alden-inter-400-latin-ext.woff2",
		"theme-alden-source-serif-4-400-latin.woff2",
	} {
		req := httptest.NewRequest(http.MethodGet, "/static/fonts/"+file, nil)
		res := httptest.NewRecorder()
		srv.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("%s = %d, want 200 (real Alden font must serve)", file, res.Code)
			continue
		}
		if ct := res.Header().Get("Content-Type"); ct != "font/woff2" {
			t.Errorf("%s content-type = %q, want font/woff2", file, ct)
		}
		if res.Body.Len() == 0 {
			t.Errorf("%s served empty body", file)
		}
	}

	// A non-preloaded but allowlisted font (heavier weight) is still servable:
	// @font-face loads it on demand, it just isn't preloaded in <head>.
	req := httptest.NewRequest(http.MethodGet, "/static/fonts/theme-alden-inter-500-latin.woff2", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("inter-500 latin (allowlisted, non-preloaded) = %d, want 200 (servable on demand)", res.Code)
	}

	// A real font file that NO theme declares must be rejected.
	req = httptest.NewRequest(http.MethodGet, "/static/fonts/theme-alden-nonexistent.woff2", nil)
	res = httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusNotFound {
		t.Errorf("undeclared font = %d, want 404 (closed namespace)", res.Code)
	}
}

// TestLayoutEmitsAldenPreloads proves the rendered layout <head> carries the
// preload links for the Alden theme (real catalog, no override) and none for a
// theme that ships no fonts (material), keeping pages lean.
func TestLayoutEmitsAldenPreloads(t *testing.T) {
	srv := New()

	// Alden page → two preload links (Inter 400 latin + latin-ext).
	req := httptest.NewRequest(http.MethodGet, "/docs?theme=alden", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("/docs?theme=alden = %d, want 200", res.Code)
	}
	body := res.Body.String()
	for _, want := range []string{
		`<link rel="preload" as="font" href="/static/fonts/theme-alden-inter-400-latin.woff2?v=` + lib.AssetsVersion,
		`<link rel="preload" as="font" href="/static/fonts/theme-alden-inter-400-latin-ext.woff2?v=` + lib.AssetsVersion,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("Alden layout missing preload %q", want)
		}
	}

	// Material page → no font preloads (ships none): lean head.
	req = httptest.NewRequest(http.MethodGet, "/docs", nil)
	res = httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	body = res.Body.String()
	if strings.Contains(body, `/static/fonts/`) {
		t.Errorf("material layout must emit no font preload, got %q", body)
	}
	if !strings.Contains(body, `class="theme-material"`) {
		t.Errorf("material layout must default to theme-material")
	}
}
