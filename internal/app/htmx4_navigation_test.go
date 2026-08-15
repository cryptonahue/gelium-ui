package app

import (
	"strings"
	"testing"

	webassets "geliumui/web"
)

func TestHTMX4NavigationContract(t *testing.T) {
	layout := readAsset(t, "templates/layout.html")
	for _, contract := range []string{
		`hx-boost:inherited="true"`,
		`hx-swap:inherited="innerMorph"`,
		`hx-push-url:inherited="true"`,
		`hx-history-elt`,
		`name="htmx-config"`,
		`&quot;history&quot;:true`,
		`&quot;morphSkip&quot;`,
		`&quot;morphIgnore&quot;`,
		`/static/morph-afterswap.js?v=0.4.0`,
	} {
		if !strings.Contains(layout, contract) {
			t.Errorf("layout is missing HTMX 4 contract %q", contract)
		}
	}
}

func TestHTMX4RuntimeAndEnhancementsArePresent(t *testing.T) {
	runtime := readAsset(t, "static/htmx.min.js")
	if !strings.Contains(runtime, "4.0.0-beta6") || !strings.Contains(runtime, "innerMorph") {
		t.Fatal("embedded runtime is not the official HTMX 4.0.0-beta6 build")
	}
	app := readAsset(t, "static/app.js")
	search := readAsset(t, "static/search.js")
	for name, content := range map[string]string{"app.js": app, "search.js": search} {
		for _, contract := range []string{"htmx:after:swap", "htmx:before:history:restore", "data-gelium-"} {
			if !strings.Contains(content, contract) {
				t.Errorf("%s is missing idempotent post-swap contract %q", name, contract)
			}
		}
	}
}

func TestHTMX4NavigationPreservesNativeLinkBoundaries(t *testing.T) {
	layout := readAsset(t, "templates/layout.html")
	if !strings.Contains(layout, `href="#main-content"`) {
		t.Error("layout fixture lost the native skip-link anchor")
	}
	external := readAsset(t, "templates/docs-topbar.html")
	if !strings.Contains(external, "https://github.com") {
		t.Error("external navigation fixture is missing")
	}
	for _, path := range []string{"templates/recipe-public-feed.html", "templates/recipe-admin-resource.html", "templates/recipe-ops-queue.html"} {
		body := readAsset(t, path)
		if !strings.Contains(body, `hx-swap="outerHTML"`) || !strings.Contains(body, `hx-target=`) {
			t.Errorf("%s no longer declares explicit partial recipe swaps", path)
		}
	}
}

func readAsset(t *testing.T, path string) string {
	t.Helper()
	b, err := webassets.Assets.ReadFile(path)
	if err != nil {
		t.Fatalf("read embedded asset %s: %v", path, err)
	}
	return string(b)
}
