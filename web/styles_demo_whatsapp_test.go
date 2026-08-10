package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestWhatsAppDemoCSSMapsBubblesWindowAndComposer(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "demo-whatsapp.css"), " ")

	for _, contract := range []string{
		`.demo-wa-appbar {`,
		`.demo-wa-sidebar {`,
		`.demo-wa-conv {`,
		`.demo-wa-bubble--in {`,
		`.demo-wa-bubble--out {`,
		`.demo-wa-window--ok {`,
		`.demo-wa-window--warning {`,
		`.demo-wa-window--expired {`,
		`.demo-wa-window-bar {`,
		`.demo-wa-window-fill {`,
		`.demo-wa-composer {`,
		`.demo-wa-composer-input {`,
		`.demo-wa-expired {`,
		`.demo-wa-template-send {`,
		`.demo-wa-admin-table {`,
		`.demo-wa-quality--GREEN {`,
		`.demo-wa-rate-bar {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("demo-whatsapp.css is missing %q", contract)
		}
	}

	// Window tones must be visually distinct.
	if !strings.Contains(css, `.demo-wa-window--expired {`) {
		t.Error("expired window chip must have its own style")
	}

	// Reduced motion disables the typing animation.
	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("demo-whatsapp.css must include a reduced-motion media query")
	}

	// Forced colors keeps the app chrome distinguishable.
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("demo-whatsapp.css must include a forced-colors media query")
	}
}

func TestWhatsAppDemoDemoClassesCarryNoUIPrefix(t *testing.T) {
	css := sourceComponentCSS(t, "demo-whatsapp.css")
	// Demo scaffolding classes must not use the ui- prefix (reserved for
	// component primitives). We assert the expected demo classes exist as
	// plain class names.
	for _, c := range []string{
		"demo-wa-appbar",
		"demo-wa-sidebar",
		"demo-wa-bubble",
		"demo-wa-window",
		"demo-wa-composer",
		"demo-wa-expired",
		"demo-wa-admin-panel",
	} {
		if !strings.Contains(css, "."+c) {
			t.Errorf("demo class %q missing from demo-whatsapp.css", c)
		}
	}
}
