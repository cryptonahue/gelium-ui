package lib

import (
	"html/template"
	"regexp"
	"strings"
	"testing"
)

func TestCardPrimitiveCSSMapsVariantsToThemeTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-card {`,
		`border-radius: var(--ui-card-radius);`,
		`.ui-card-elevated { background: var(--ui-card-container-elevated); box-shadow: var(--ui-shadow-1);`,
		`.ui-card-filled { background: var(--ui-card-container-filled);`,
		`.ui-card-outlined { background: var(--ui-card-container-outlined); border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-card-outline-color);`,
		`.ui-card:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing card contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-card { border: 1px solid CanvasText;") {
		t.Error("card must keep a visible boundary in forced colors")
	}
}

func TestCardThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-card-radius:",
		"--ui-card-container-elevated:",
		"--ui-card-container-filled:",
		"--ui-card-container-outlined:",
		"--ui-card-outline-color:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing card token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesCardContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-card{`,
		`.ui-card-elevated{`,
		`.ui-card-filled{`,
		`.ui-card-outlined{`,
		`@media (forced-colors:active)`,
		`var(--ui-card-container-filled)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled card CSS is missing %q", contract)
		}
	}
}

// TestCardSlotSourceCSSContracts locks the Phase F public slot CSS contracts
// in the source component file: media (literal 16/9 aspect-ratio + cover fill),
// tag badge (reusing --ui-badge-* tokens) and meta line (--ui-space-* spacing).
func TestCardSlotSourceCSSContracts(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "card.css"), " ")

	for _, contract := range []string{
		`.ui-card-media {`,
		`aspect-ratio: 16 / 9;`,
		`object-fit: cover;`,
		`.ui-card-tag {`,
		`var(--ui-badge-`,
		`.ui-card-meta {`,
		`var(--ui-space-`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source card.css is missing slot contract %q", contract)
		}
	}
}

// TestCardAspectRatioNotTokenized is the load-bearing anti-tokenization guard:
// card media aspect-ratio is structural geometry (same rule as Video and
// Feature Card media) and must stay literal 16 / 9, never var(--ui-*).
func TestCardAspectRatioNotTokenized(t *testing.T) {
	css := sourceComponentCSS(t, "card.css")
	if strings.Contains(css, `aspect-ratio: var(--ui-`) {
		t.Error("card.css must keep aspect-ratio literal 16 / 9 (structural geometry, never tokenized)")
	}
}

// TestCardSlotCompiledCSSContracts proves the slot rules survive the Tailwind
// build and are present in the embedded compiled app.css.
func TestCardSlotCompiledCSSContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-card-media{`,
		`.ui-card-tag{`,
		`.ui-card-meta{`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled app CSS is missing slot contract %q", contract)
		}
	}
}

// TestCardTemplateSlotGuards locks the {{define "card"}} primitive's slot
// guards and their source order (media → tag → meta → title → body → CTA):
// media first, the action wrapper last (R1 full-slot scenario).
func TestCardTemplateSlotGuards(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "card.html")

	for _, contract := range []string{
		`{{if .Media}}`,
		`{{if .Tag}}`,
		`{{if .Meta}}`,
		`{{if .CTA}}`,
		`{{template "button" .CTA}}`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("card.html is missing slot guard contract %q", contract)
		}
	}

	media := strings.Index(tmpl, `{{if .Media}}`)
	tag := strings.Index(tmpl, `{{if .Tag}}`)
	meta := strings.Index(tmpl, `{{if .Meta}}`)
	title := strings.Index(tmpl, `{{if .Title}}`)
	body := strings.Index(tmpl, `{{if .Body}}`)
	action := strings.Index(tmpl, `{{if .CTA}}`)
	for _, p := range []struct {
		name string
		a, b int
	}{
		{"media before tag", media, tag},
		{"tag before meta", tag, meta},
		{"meta before title", meta, title},
		{"title before body", title, body},
		{"body before CTA", body, action},
	} {
		if p.a < 0 || p.b < 0 {
			t.Errorf("card.html is missing a slot guard, order check %q incomplete (media=%d tag=%d meta=%d title=%d body=%d action=%d)", p.name, media, tag, meta, title, body, action)
			continue
		}
		if p.a > p.b {
			t.Errorf("slot order violated in card.html: %s", p.name)
		}
	}
}

// cardSlotData mirrors the {{define "card"}} primitive's public slot fields.
// Media is pre-rendered HTML (feature-card precedent); CTA carries the Button
// partial's data.
type cardSlotData struct {
	Variant string
	Media   template.HTML
	Tag     string
	Meta    string
	Title   string
	Body    string
	CTA     *cardCTA
}

type cardCTA struct {
	Variant  string
	Label    string
	Href     string
	IconSVG  string
	Disabled bool
	Loading  bool
}

// renderCardPrimitive parses card.html + button.html from the repository and
// executes the "card" template with the given slot data.
func renderCardPrimitive(t *testing.T, data cardSlotData) string {
	t.Helper()
	cardHTML := repositoryFile(t, "lib", "templates", "card.html")
	buttonHTML := repositoryFile(t, "lib", "templates", "button.html")
	tmpl, err := template.New("card").Parse(cardHTML)
	if err != nil {
		t.Fatalf("parse card.html: %v", err)
	}
	tmpl, err = tmpl.Parse(buttonHTML)
	if err != nil {
		t.Fatalf("parse button.html: %v", err)
	}
	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, "card", data); err != nil {
		t.Fatalf("execute card primitive: %v", err)
	}
	return buf.String()
}

// TestCardTemplateRendersSlots proves the render behavior of the primitive:
// full slot data renders media first, the CTA action wrapper last, with
// tag/meta/title/body between; no slot data renders only title/body with zero
// empty wrappers and the elevated default variant (R1 no-slot output).
func TestCardTemplateRendersSlots(t *testing.T) {
	full := cardSlotData{
		Variant: "filled",
		Media:   template.HTML(`<img src="/cover.jpg" alt="Cover">`),
		Tag:     "Recipe",
		Meta:    "Updated 2026-08-14",
		Title:   "Card with slots",
		Body:    "Body text.",
		CTA:     &cardCTA{Variant: "primary", Label: "Read more", Href: "/recipe/1"},
	}
	out := renderCardPrimitive(t, full)

	if !strings.Contains(out, `<img src="/cover.jpg" alt="Cover">`) {
		t.Error("card media slot must render the raw media HTML")
	}
	for _, contract := range []string{
		`class="ui-card ui-card-filled"`,
		`class="ui-card-media"`,
		`class="ui-card-tag">Recipe<`,
		`class="ui-card-meta">Updated 2026-08-14<`,
		`class="ui-card-title">Card with slots<`,
		`class="ui-card-body">Body text.<`,
		`class="ui-card-action"`,
	} {
		if !strings.Contains(out, contract) {
			t.Errorf("full-slot card render is missing %q", contract)
		}
	}

	media := strings.Index(out, `class="ui-card-media"`)
	tag := strings.Index(out, `class="ui-card-tag"`)
	meta := strings.Index(out, `class="ui-card-meta"`)
	title := strings.Index(out, `class="ui-card-title"`)
	body := strings.Index(out, `class="ui-card-body"`)
	action := strings.Index(out, `class="ui-card-action"`)
	for _, p := range []struct {
		name string
		a, b int
	}{
		{"media first", media, tag},
		{"tag before meta", tag, meta},
		{"meta before title", meta, title},
		{"title before body", title, body},
		{"body before action", body, action},
	} {
		if p.a < 0 || p.b < 0 {
			t.Errorf("full-slot render is missing elements for order check %q (media=%d tag=%d meta=%d title=%d body=%d action=%d)", p.name, media, tag, meta, title, body, action)
			continue
		}
		if p.a > p.b {
			t.Errorf("full-slot render violates slot order: %s", p.name)
		}
	}
	if action < 0 {
		t.Fatal("full-slot render must include the .ui-card-action wrapper")
	}
	actionHTML := out[action:]
	for _, contract := range []string{
		`ui-button ui-button-primary`,
		`href="/recipe/1"`,
		`>Read more<`,
	} {
		if !strings.Contains(actionHTML, contract) {
			t.Errorf(".ui-card-action must wrap the rendered Button partial, missing %q", contract)
		}
	}

	none := renderCardPrimitive(t, cardSlotData{Title: "Only title", Body: "Only body"})
	for _, contract := range []string{
		`<article class="ui-card ui-card-elevated">`,
		`class="ui-card-title">Only title<`,
		`class="ui-card-body">Only body<`,
	} {
		if !strings.Contains(none, contract) {
			t.Errorf("no-slot card render is missing %q", contract)
		}
	}
	for _, absent := range []string{
		`.ui-card-media`,
		`.ui-card-tag`,
		`.ui-card-meta`,
		`.ui-card-action`,
		`ui-card-filled`,
	} {
		if strings.Contains(none, absent) {
			t.Errorf("no-slot card render must not contain %q (empty wrappers are forbidden)", absent)
		}
	}
}
