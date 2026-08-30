package lib

import (
	"regexp"
	"strings"
	"testing"
)

// TestIconSetIsTrustedDecorativeMarkup pins the curated Material Symbols
// contract: every glyph in internal/app/icons.go is a trusted inline SVG that
// is aria-hidden, unfocusable, and currentColor-filled (themeable), and the
// generated file stays in sync with its source script.
func TestIconSetIsTrustedDecorativeMarkup(t *testing.T) {
	icons := repositoryFile(t, "internal", "app", "icons.go")
	script := repositoryFile(t, "scripts", "copy-icons.mjs")

	// Every generated glyph must carry the decorative-accessibility contract.
	// Each map row has the shape: 	"name": `<svg ...>...</svg>`,
	type iconRow struct {
		name  string
		attrs string
	}
	var rows []iconRow
	for _, line := range strings.Split(icons, "\n") {
		line = strings.TrimSpace(line)
		marker := "\": `"
		if !strings.HasPrefix(line, "\"") || !strings.Contains(line, marker) {
			continue
		}
		nameEnd := strings.Index(line[1:], "\"") + 1
		name := line[1:nameEnd]
		svgStart := strings.Index(line, marker) + len(marker)
		svg := line[svgStart:]
		end := strings.Index(svg, ">`")
		if end < 0 {
			t.Fatalf("malformed icon row for %q: %q", name, line)
		}
		rows = append(rows, iconRow{name: name, attrs: svg[:end]})
	}
	if len(rows) < 36 {
		t.Fatalf("generated icon rows = %d, want at least 36 (the curated allowlist)", len(rows))
	}
	seen := map[string]bool{}
	for _, r := range rows {
		seen[r.name] = true
		if !strings.Contains(r.attrs, `aria-hidden="true"`) {
			t.Errorf("icon %q missing aria-hidden=\"true\"", r.name)
		}
		if !strings.Contains(r.attrs, `focusable="false"`) {
			t.Errorf("icon %q missing focusable=\"false\"", r.name)
		}
		if !strings.Contains(r.attrs, `fill="currentColor"`) {
			t.Errorf("icon %q missing fill=\"currentColor\"", r.name)
		}
		if !strings.Contains(r.attrs, `class="ui-icon"`) {
			t.Errorf("icon %q missing class=\"ui-icon\"", r.name)
		}
		if strings.Contains(r.attrs, `width="`) || strings.Contains(r.attrs, `height="`) {
			t.Errorf("icon %q must not hardcode width/height (the .ui-icon class owns sizing)", r.name)
		}
	}

	// The script's allowlist and the generated map must not drift: every
	// ICONS entry in copy-icons.mjs resolves to a generated row.
	allowlist := regexp.MustCompile(`(?m)^\s*"([a-z_]+)",?$`)
	for _, m := range allowlist.FindAllStringSubmatch(script, -1) {
		if !seen[m[1]] {
			t.Errorf("script allowlist icon %q missing from generated icons.go", m[1])
		}
	}
}

func TestIconCSSSupportsTablerStroke(t *testing.T) {
	css := repositoryFile(t, "lib", "styles", "icon.css")
	for _, want := range []string{
		`[data-gelium-set="tabler"]`,
		"stroke: currentColor",
		"fill: none",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("icon.css missing Tabler stroke support %q", want)
		}
	}
}
