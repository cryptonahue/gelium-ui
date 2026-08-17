package app

import (
	"net/http"
	"strings"
	"testing"
)

// Docs-page contract tests for Avatar. The companion recipe-primitive tests
// (renderPartial + avatarView) live in avatar_test.go; this file only covers
// the /components/avatar documentation page.

func TestAvatarDocsRouteDogfoodsLiveSpecimenAndShell(t *testing.T) {
	res := renderDocsPage(t, "/components/avatar", (*server).avatarDocs)

	if res.Code != http.StatusOK {
		t.Fatalf("avatar docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Avatar</h1>`,
		`class="docs-shell"`,
		`class="ui-avatar ui-avatar--sm"`,
		`class="ui-avatar ui-avatar--md"`,
		`class="ui-avatar ui-avatar--lg"`,
		`class="ui-avatar-initials"`,
		`class="ui-avatar-image"`,
		`aria-hidden="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("avatar docs are missing %q", contract)
		}
	}
}

func TestAvatarDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := renderDocsMethod(t, "/components/avatar", http.MethodPost, (*server).avatarDocs)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST avatar status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}