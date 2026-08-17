package app

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webassets "geliumui/site/web"
)

// TestDocsIndexRendersExplainer proves the /docs hub keeps the deep-dive essay
// (quick start, themes, Base UI) after Start here orientation.
func TestDocsIndexRendersExplainer(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">Start here</h2>",
		">Deep dive</h2>",
		"Quick start",
		"npm install gelium-ui",
		"Themes and the layer model",
		"Base UI",
		"BASE_URL",
		"?theme=basecoat",
		"https://github.com/cryptonahue/gelium-ui",
		"gelium-ui/themes/theme-material.css",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("docs index is missing explainer content %q", contract)
		}
	}
	start := strings.Index(body, ">Start here</h2>")
	deep := strings.Index(body, ">Deep dive</h2>")
	if start < 0 || deep < 0 || start > deep {
		t.Fatalf("Start here must precede Deep dive (idxs %d, %d)", start, deep)
	}
	// Exactly one H1: the shell hub title.
	if count := strings.Count(body, "<h1"); count != 1 {
		t.Errorf("docs index must render exactly one h1, got %d", count)
	}
}

// TestContentIndexExplainsLibrary proves the embedded docs root file itself
// carries the explainer sections (npm-first quick start / themes / Base UI).
func TestContentIndexExplainsLibrary(t *testing.T) {
	source, err := fs.ReadFile(webassets.Assets, "content/index.md")
	if err != nil {
		t.Fatalf("read content/index.md: %v", err)
	}
	content := string(source)
	for _, section := range []string{
		"Quick start",
		"npm install gelium-ui",
		"Themes and the layer model",
		"Base UI",
		"BASE_URL",
		"?theme=basecoat",
		"https://github.com/cryptonahue/gelium-ui",
		"gelium-ui/themes/theme-material.css",
	} {
		if !strings.Contains(content, section) {
			t.Errorf("content/index.md is missing %q", section)
		}
	}
	if strings.Contains(content, "themes/theme-material/theme.css") {
		t.Error("content/index.md must not use legacy nested theme path")
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
