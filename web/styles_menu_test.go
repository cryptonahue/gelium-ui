package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestMenuPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "menu.css"), " ")

	for _, contract := range []string{
		`.ui-menu {`,
		`min-width: 112px;`,
		`border-radius: var(--ui-menu-container-radius);`,
		`box-shadow: var(--ui-menu-container-elevation);`,
		`--ui-menu-container-color:`,
		`--ui-menu-container-radius:`,
		`--ui-menu-container-elevation:`,
		`--ui-menu-item-height: 48px;`,
		`--ui-menu-item-leading-space: 12px;`,
		`--ui-menu-item-trailing-space: 12px;`,
		`--ui-menu-item-icon-size: 24px;`,
		`.ui-menu-item {`,
		`min-height: var(--ui-menu-item-height);`,
		`.ui-menu-item-button {`,
		`.ui-menu-item-link {`,
		`.ui-menu-item-label {`,
		`.ui-menu-divider {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source menu CSS is missing contract %q", contract)
		}
	}
}

// TestMenuStateLayersAndNativeSelection guards the platform-first contract: the
// state layers use the shared --ui-state-* opacities and every selected/disabled
// style derives from the native checkbox/radio :checked / :disabled pseudo
// classes, exactly like the List and Segmented button components.
func TestMenuStateLayersAndNativeSelection(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "menu.css"), " ")

	for _, contract := range []string{
		`.ui-menu-item::before {`,
		`.ui-menu-item:hover::before {`,
		`.ui-menu-item:focus-within::before {`,
		`.ui-menu-item:active::before {`,
		`.ui-menu-item:has(input:checked) {`,
		`.ui-menu-item:has(input:disabled) {`,
		`.ui-menu-item-button:disabled`,
		`.ui-menu-item-link:focus-visible {`,
		`.ui-menu-item-button:focus-visible {`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`--ui-menu-item-hover-opacity: var(--ui-state-hover-opacity);`,
		`--ui-menu-item-pressed-opacity: var(--ui-state-pressed-opacity);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source menu CSS is missing state/selection contract %q", contract)
		}
	}
}

// TestMenuPopoverDeclarativeOpenClose guards the roadmap's "top layer/open-close
// puede ser declarativo": the surface opens through the native popover attribute
// and anchor positioning is a progressive enhancement inside @supports (the
// popovertarget control supplies the implicit anchor, so no explicit
// anchor-name is needed), never a JS requirement.
func TestMenuPopoverDeclarativeOpenClose(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "menu.css"), " ")

	for _, contract := range []string{
		`.ui-menu[popover] {`,
		`position: fixed;`,
		`@supports (anchor-name:`,
		`inset-block-start: anchor(bottom);`,
		`inset-inline-start: anchor(left);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source menu CSS is missing popover/positioning contract %q", contract)
		}
	}
}

func TestMenuDemoClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "menu.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "menu.css"), " ")

	for _, cls := range []string{
		"menu-demo-grid",
		"menu-demo-group",
		"menu-demo-heading",
		"menu-demo-anchor",
		"menu-demo-caret",
		"menu-demo-form",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("menu.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("menu.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-menu",
		"ui-menu-item",
		"ui-menu-item--select",
		"ui-menu-item-button",
		"ui-menu-item-link",
		"ui-menu-item-label",
		"ui-menu-item-label-text",
		"ui-menu-item-icon",
		"ui-menu-divider",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("menu.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("menu.css is missing primitive selector .%s", cls)
		}
	}

	if strings.Contains(tmpl, "ui-menu-demo") {
		t.Error("menu.html must not ui-prefix the demo scaffolding classes")
	}
	if strings.Contains(css, ".ui-menu-demo") {
		t.Error("menu.css must not define .ui-menu-demo selectors")
	}
}

func TestMenuReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "menu.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("menu.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `transition: none;`) {
		t.Error("menu reduced-motion must disable transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("menu.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-menu { border: 1px solid CanvasText;`) {
		t.Error("menu surface must stay discernible in forced colors")
	}
	if !strings.Contains(css, `.ui-menu-item:has(input:checked) { background: Highlight; color: HighlightText;`) {
		t.Error("menu selected item must switch to Highlight/HighlightText in forced colors")
	}
	if !strings.Contains(css, `.ui-menu-item-button:disabled`) {
		t.Error("menu disabled items must switch to GrayText in forced colors")
	}
	if !strings.Contains(css, `.ui-menu-item-link:focus-visible { outline-color: Highlight;`) {
		t.Error("menu link focus ring must switch to Highlight in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesMenuContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-menu`,
		`.ui-menu-item`,
		`.ui-menu-divider`,
		`.menu-demo-grid`,
		`.menu-demo-anchor`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled menu CSS is missing %q", contract)
		}
	}
}
