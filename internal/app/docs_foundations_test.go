package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFoundationDocsRoutesRenderTokenShowcases(t *testing.T) {
	cases := []struct {
		path  string
		title string
		must  []string
	}{
		{"/docs/typography", "Typography", []string{"--ui-type-display-lg-size", "--ui-type-body-md-line-height", "65ch", "https://www.w3.org/TR/WCAG22/#text-spacing", "href=\"/docs/typography\""}},
		{"/docs/spacing", "Spacing", []string{"--ui-space-1", "--ui-space-8", "docs-space-composition", "https://design-system.service.gov.uk/styles/spacing/", "href=\"/docs/spacing\""}},
		{"/docs/colors", "Colors", []string{"--ui-color-canvas", "--ui-color-focus-ring", "--ui-color-success", "theme-dark", "https://www.w3.org/TR/WCAG22/#contrast-minimum", "href=\"/docs/colors\""}},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200", tc.path, res.Code)
			}
			body := res.Body.String()
			if !strings.Contains(body, ">"+tc.title+"</h1>") {
				t.Errorf("missing H1 %q", tc.title)
			}
			for _, want := range tc.must {
				if !strings.Contains(body, want) {
					t.Errorf("missing %q", want)
				}
			}
		})
	}
}

func TestFoundationDocsAreInHubNavAndLLMSUX(t *testing.T) {
	for _, path := range []string{"/docs", "/components/button"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
		if res.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", path, res.Code)
		}
		body := res.Body.String()
		for _, href := range []string{"/docs/typography", "/docs/spacing", "/docs/colors"} {
			if !strings.Contains(body, `href="`+href+`"`) {
				t.Errorf("%s missing %s", path, href)
			}
		}
	}
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/llms-ux.txt", nil))
	body := res.Body.String()
	for _, id := range []string{"DOC-TYPE", "DOC-SPACE", "DOC-COLOR", "TYPE-ROLE", "SPACE-SCALE", "COLOR-FOCUS"} {
		if !strings.Contains(body, id) {
			t.Errorf("llms-ux.txt missing ID %s", id)
		}
	}
}
