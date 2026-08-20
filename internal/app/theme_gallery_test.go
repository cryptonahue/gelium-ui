package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestThemeGalleryRendersLiveSpecimens proves /docs/themes/gallery renders the
// real component demos (buttons, fields, card, badge, toast, dialog) under the
// active theme — the live kitchen-sink the handbook describes in prose.
func TestThemeGalleryRendersLiveSpecimens(t *testing.T) {
	srv := New()
	req := httptest.NewRequest(http.MethodGet, "/docs/themes/gallery?theme=alden", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("/docs/themes/gallery?theme=alden = %d, want 200", res.Code)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Fonts", "Self-hosted", "theme-gallery-swatch-chip--canvas", "theme-gallery-swatch-chip--primary",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("Alden gallery missing %q", contract)
		}
	}
	if strings.Contains(body, "ZgotmplZ") || strings.Contains(body, `style="background:`) {
		t.Error("gallery swatches must use safe static classes, not sanitized or inline dynamic CSS")
	}

	// The active theme must be applied to the document root and the preloads
	// emitted (Alden ships fonts).
	if !strings.Contains(body, `class="theme-alden"`) {
		t.Error("gallery must render under the active theme (theme-alden)")
	}
	if !strings.Contains(body, `rel="preload" as="font"`) {
		t.Error("gallery must emit the active theme's font preloads")
	}

	// Live specimens from the demo slots must be present.
	for _, contract := range []string{
		`ui-button ui-button-primary`,
		`ui-text-field ui-text-field-outlined`,
		`ui-card ui-card-elevated`,
		`ui-badge`,
		`ui-toast ui-toast-info`,
		`Open dialog`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("gallery missing live specimen %q", contract)
		}
	}

	// The intro heading and a link back to the themes handbook.
	if !strings.Contains(body, "live kitchen-sink") {
		t.Error("gallery must render its intro copy")
	}
	if !strings.Contains(body, `/docs/themes`) {
		t.Error("gallery must link to the themes handbook")
	}
}

// TestThemeGalleryDefaultsToMaterial proves the gallery defaults to theme-material
// when no theme is requested (same default as every page), and stays a no-extra-font
// lean head (Material ships no fonts).
func TestThemeGalleryDefaultsToMaterial(t *testing.T) {
	srv := New()
	req := httptest.NewRequest(http.MethodGet, "/docs/themes/gallery", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("/docs/themes/gallery = %d, want 200", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, `class="theme-material"`) {
		t.Error("gallery must default to theme-material")
	}
	if strings.Contains(body, `/static/fonts/`) {
		t.Error("material gallery must emit no font preload (ships none)")
	}
}

// TestThemeGalleryDarkRendersClass proves ?scheme=dark applies the explicit dark
// class route to the gallery so the dark direction shows.
func TestThemeGalleryDarkRendersClass(t *testing.T) {
	srv := New()
	req := httptest.NewRequest(http.MethodGet, "/docs/themes/gallery?theme=alden&scheme=dark", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("/docs/themes/gallery?theme=alden&scheme=dark = %d, want 200", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, `class="theme-alden theme-dark"`) && !strings.Contains(body, `data-theme="dark"`) {
		t.Errorf("dark gallery must carry the explicit dark class or data-theme, got %q", body)
	}
}

func TestThemeGalleryShowsDesignSystemSectionsAndNavigation(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/themes/gallery?theme=linear&scheme=dark", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.Code)
	}
	body := res.Body.String()
	for _, marker := range []string{"Fonts", "system fallback", "Color roles", "Type specimens", "Spacing, shape, and anatomy", "Reference → Gelium", "Quote / testimonial", "#08090a", "--ui-type-display-lg", "theme-gallery-type-card--Display", "Theme Gallery", `href="/docs/themes/gallery?scheme=dark&amp;theme=linear"`} {
		if !strings.Contains(body, marker) {
			t.Errorf("gallery missing %q", marker)
		}
	}
}
