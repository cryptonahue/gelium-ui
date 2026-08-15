package app

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webassets "geliumui/web"
)

// TestDocsIndexRendersExplainer proves the /docs hub leads with the expanded
// docs root: what Gelium UI is, how to use it, and the theme/layer model.
func TestDocsIndexRendersExplainer(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Quick start",
		"Themes and the layer model",
		"Base UI",
		"BASE_URL",
		"?theme=basecoat",
		"https://github.com/cryptonahue/gelium-ui",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("docs index is missing explainer content %q", contract)
		}
	}
	// Exactly one H1: the shell hub title. The embedded docs root must not
	// bring a second heading level 1 into the page.
	if count := strings.Count(body, "<h1"); count != 1 {
		t.Errorf("docs index must render exactly one h1, got %d", count)
	}
}

// TestContentIndexExplainsLibrary proves the embedded docs root file itself
// carries the explainer sections (WHAT / HOW / themes / Base UI).
func TestContentIndexExplainsLibrary(t *testing.T) {
	source, err := fs.ReadFile(webassets.Assets, "content/index.md")
	if err != nil {
		t.Fatalf("read content/index.md: %v", err)
	}
	content := string(source)
	for _, section := range []string{
		"Quick start",
		"Themes and the layer model",
		"Base UI",
		"BASE_URL",
		"?theme=basecoat",
		"https://github.com/cryptonahue/gelium-ui",
	} {
		if !strings.Contains(content, section) {
			t.Errorf("content/index.md is missing %q", section)
		}
	}
}

// TestStripDocsRootH1 proves the H1-strip helper removes only the leading
// heading line, keeping the body intact for embedding under the hub H1.
func TestStripDocsRootH1(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{name: "leading h1 stripped", source: "# Gelium UI\n\nBody", want: "Body"},
		{name: "no leading h1", source: "Body", want: "Body"},
		{name: "h1 mid-file kept", source: "Lead\n# Not first\n", want: "Lead\n# Not first\n"},
		{name: "blank after h1", source: "# Gelium UI\n\n\nBody", want: "Body"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripDocsRootH1(tt.source); got != tt.want {
				t.Errorf("stripDocsRootH1(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}
