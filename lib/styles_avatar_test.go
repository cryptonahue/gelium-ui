package lib

import (
	"regexp"
	"strings"
	"testing"
)

// TestAvatarPrimitiveCSSMapsSizesAndTokens proves the avatar paints from the
// semantic surface tokens, maps sm/md/lg onto the core size scale and stays
// circular with the full radius.
func TestAvatarPrimitiveCSSMapsSizesAndTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "avatar.css"), " ")

	for _, contract := range []string{
		`.ui-avatar {`,
		`--ui-avatar-container: var(--ui-color-surface-container);`,
		`--ui-avatar-fg: var(--ui-color-fg);`,
		`--ui-avatar-size-sm: var(--ui-size-control);`,
		`--ui-avatar-size-md: var(--ui-size-item);`,
		`--ui-avatar-size-lg: var(--ui-size-item-lg);`,
		`background: var(--ui-avatar-container);`,
		`color: var(--ui-avatar-fg);`,
		`border-radius: var(--ui-radius-full);`,
		`.ui-avatar--sm {`,
		`.ui-avatar--md {`,
		`.ui-avatar--lg {`,
		`.ui-avatar-image {`,
		`object-fit: cover;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("avatar.css is missing contract %q", contract)
		}
	}

	// The primitive must consume the --ui-size-* scale (three distinct steps).
	if !strings.Contains(css, "var(--ui-size-control)") ||
		!strings.Contains(css, "var(--ui-size-item)") ||
		!strings.Contains(css, "var(--ui-size-item-lg)") {
		t.Error("avatar.css must consume the core size scale for sm/md/lg")
	}

	// Forced-colors boundary: the circle must stay visible when color is removed.
	if !strings.Contains(sourceComponentCSS(t, "avatar.css"), "border: 1px solid CanvasText") {
		t.Error("avatar must keep a visible boundary in forced colors")
	}
}

// TestAvatarCompiledCSSIncluded proves the compiled app.css embeds the avatar
// primitive so an unthemed page paints it from the core tokens.
func TestAvatarCompiledCSSIncluded(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{`.ui-avatar{`, `var(--ui-avatar-size-md)`, `--ui-avatar-container`} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled avatar CSS is missing %q", contract)
		}
	}
}
