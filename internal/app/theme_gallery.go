package app

import "net/http"

type themeGalleryColor struct{ Name, Token, Value, Role, SwatchClass string }
type themeGalleryFont struct{ Family, Source, Weights, Subset, Status, Tokens string }
type themeGalleryType struct{ Name, Token, Sample, Size, Line, Tracking string }
type themeGalleryMeasure struct{ Name, Token, Value, Note string }
type themeGalleryMapping struct{ Group, Items string }
type themeGalleryView struct {
	Theme, Scheme, FontSource, FontWeights string
	Behavior                               accordionBehavior
	Execution                              accordionExecution
	Reference, Direction                   string
	Colors                                 []themeGalleryColor
	Fonts                                  []themeGalleryFont
	Types                                  []themeGalleryType
	Measures                               []themeGalleryMeasure
	Mappings                               []themeGalleryMapping
	QuoteText, QuoteCite                   string
	Accordion                              *accordionView
}

// Theme metadata mirrors values declared in lib/themes/*.css. It is deliberately
// small and labeled as implementation data: the browser still resolves the live
// tokens on <html>, while this gives every value an accessible explanation.
var galleryThemeData = map[string]themeGalleryView{
	"material": {Theme: "Material", Reference: "Material 3-inspired semantic roles", Direction: "Expressive purple surfaces, generous controls, and tonal elevation.", FontSource: "System fallback stack (no webfont pack shipped)", FontWeights: "400 / 500 / 600", Colors: []themeGalleryColor{{"Canvas", "--ui-color-canvas", "#fff7ff", "page background", "theme-gallery-swatch-chip--canvas"}, {"Primary", "--ui-color-primary", "#6750a4", "primary action", "theme-gallery-swatch-chip--primary"}, {"Surface", "--ui-color-surface", "#f7f2fa", "alternate surface", "theme-gallery-swatch-chip--surface"}, {"Foreground", "--ui-color-fg", "#1d1b20", "primary text", "theme-gallery-swatch-chip--fg"}, {"Border", "--ui-color-border", "#cac4d0", "hairline", "theme-gallery-swatch-chip--border"}}, Fonts: gallerySystemFonts("Material", "Inter / system sans; SFMono-Regular / Consolas", "system fallback stacks; no dedicated webfont pack"), Measures: galleryMeasures("Material"), Mappings: galleryMappings("Material")},
	"basecoat": {Theme: "Basecoat", Reference: "Refero Vega / Basecoat UI", Direction: "Neutral utility surfaces with compact, restrained controls.", FontSource: "System fallback stack (no webfont pack shipped)", FontWeights: "400 / 500 / 600", Colors: []themeGalleryColor{{"Canvas", "--ui-color-canvas", "#ffffff", "page background", "theme-gallery-swatch-chip--canvas"}, {"Primary", "--ui-color-primary", "#171717", "primary action", "theme-gallery-swatch-chip--primary"}, {"Surface", "--ui-color-surface", "#f5f5f5", "alternate surface", "theme-gallery-swatch-chip--surface"}, {"Foreground", "--ui-color-fg", "#0a0a0a", "primary text", "theme-gallery-swatch-chip--fg"}, {"Border", "--ui-color-border", "#e5e5e5", "hairline", "theme-gallery-swatch-chip--border"}}, Fonts: gallerySystemFonts("Basecoat", "Geist Sans / Geist Mono names with platform fallbacks", "system fallback stacks; no dedicated webfont pack"), Measures: galleryMeasures("Basecoat"), Mappings: galleryMappings("Basecoat")},
	"linear":   {Theme: "Linear", Reference: "Refero Linear", Direction: "Dark-first product UI: dense geometry, hairlines, and restrained contrast.", FontSource: "System fallback stack only (Geist Sans/Mono names; no dedicated webfont pack)", FontWeights: "400 / 500 / 600", Colors: []themeGalleryColor{{"Void", "--ui-color-canvas", "#08090a", "dark-first canvas", "theme-gallery-swatch-chip--canvas"}, {"Primary", "--ui-color-primary", "#e4f222", "accent action (dark)", "theme-gallery-swatch-chip--primary"}, {"Carbon", "--ui-color-surface", "#0f1011", "raised surface", "theme-gallery-swatch-chip--surface"}, {"Mist", "--ui-color-fg", "#d0d6e0", "primary text (dark)", "theme-gallery-swatch-chip--fg"}, {"Graphite", "--ui-color-border", "#23252a", "hairline (dark)", "theme-gallery-swatch-chip--border"}}, Fonts: gallerySystemFonts("Linear", "Geist Sans / Geist Mono, then platform sans/mono stacks", "system fallback only; no dedicated webfont pack"), Measures: galleryMeasures("Linear"), Mappings: galleryMappings("Linear")},
	"vercel":   {Theme: "Vercel", Reference: "Refero Vercel · design-bites DESIGN.md", Direction: "Light-first typeset terminal on white paper: crisp hairlines, black actions, blue interaction, and compact status dots.", FontSource: "Inter/system fallback; Geist Sans/Mono reference names only", FontWeights: "400 / 500 / 600", Colors: []themeGalleryColor{{"Canvas", "--ui-color-canvas", "#fafafa", "paper canvas", "theme-gallery-swatch-chip--canvas"}, {"Elevated", "--ui-color-surface", "#ffffff", "raised surface", "theme-gallery-swatch-chip--surface"}, {"Recessed", "--ui-color-surface-container", "#f2f2f2", "recessed surface", "theme-gallery-swatch-chip--secondary"}, {"Ink", "--ui-color-fg", "#171717", "primary text / black action", "theme-gallery-swatch-chip--fg"}, {"Interactive", "--ui-color-primary", "#0072f5", "interactive blue", "theme-gallery-swatch-chip--primary"}, {"Hairline", "--ui-color-border", "#ebebeb", "shadow-as-border outline", "theme-gallery-swatch-chip--border"}}, Fonts: gallerySystemFonts("Vercel", "Inter / system sans; Geist Mono name with platform mono fallbacks", "Geist pack not verified or bundled; fallback stacks only"), Measures: []themeGalleryMeasure{{"Base unit", "--ui-space-1", "4px", "4px-derived compact rhythm"}, {"Control", "--ui-size-control", "2.75rem", "touch-safe; reference 32/40px heights are not copied"}, {"Field", "--ui-size-field", "3rem", "floating-label and 44px touch contract"}, {"Card radius", "--ui-card-radius", "12px", "reference card panel anatomy"}, {"Container", "--ui-container-max", "80rem", "about 1280px; not the reference breakpoint list"}}, Mappings: []themeGalleryMapping{{"Implemented directly", "semantic colors, black primary/outline controls, fields, cards, badges/status dots, type direction, spacing, dark class route, focus-visible and live Gelium components"}, {"Adapted / approximated", "Geist becomes Inter/system fallback; Refero anatomy maps to Gelium tokens; status palette remains compact dots/icons; card panels use 12px while functional controls use 6px"}, {"Missing or intentionally out of scope", "verified Geist font files, 32/40px form heights, 45px breakpoint list, brand triangle/wordmark, marketing gradients, bespoke product-stage and avatar-stack primitives, and a core focus double-ring token"}}},
	"alden":    {Theme: "Alden", Reference: "Refero Alden", Direction: "Serene clinic on warm parchment: achromatic with sky and sage punctuation.", FontSource: "Self-hosted Inter + Source Serif 4 (SIL OFL); commercial Stk bureau fonts unavailable", FontWeights: "Inter 400 / 500 / 600; Source Serif 4 400", Colors: []themeGalleryColor{{"Paper", "--ui-color-canvas", "#ffffff", "dominant canvas", "theme-gallery-swatch-chip--canvas"}, {"Parchment", "--ui-color-surface", "#f3f1eb", "second surface", "theme-gallery-swatch-chip--surface"}, {"Ink", "--ui-color-fg", "#28262a", "primary text", "theme-gallery-swatch-chip--fg"}, {"Sky", "--ui-color-secondary", "#97cde5", "highlight", "theme-gallery-swatch-chip--secondary"}, {"Sage", "--ui-color-primary", "#c8dfaa", "action", "theme-gallery-swatch-chip--primary"}}, Fonts: []themeGalleryFont{{"Inter", "Self-hosted WOFF2", "400 / 500 / 600", "latin + latin-ext", "Loaded from the bundled Alden font pack; 400 latin/latin-ext preloaded", "--ui-font-sans; body, title, label, and UI type families"}, {"Source Serif 4", "Self-hosted WOFF2", "400", "latin + latin-ext", "Loaded from the bundled Alden font pack; display mapping", "--ui-font-display; display and headline type families"}}, Measures: galleryMeasures("Alden"), Mappings: galleryMappings("Alden")},
	"baseui":   {Theme: "Base UI-inspired", Reference: "Base UI composition model (headless/state-driven; not official Base UI code)", Direction: "Neutral, low-ornament surfaces that make open/closed state and focus visible.", FontSource: "System fallback stack (no webfont pack shipped)", FontWeights: "400 / 500 / 600", Colors: []themeGalleryColor{{"Canvas", "--ui-color-canvas", "#ffffff", "page background", "theme-gallery-swatch-chip--canvas"}, {"Primary", "--ui-color-primary", "#18181b", "state/action foreground", "theme-gallery-swatch-chip--primary"}, {"Surface", "--ui-color-surface", "#fafafa", "quiet surface", "theme-gallery-swatch-chip--surface"}, {"Foreground", "--ui-color-fg", "#18181b", "primary text", "theme-gallery-swatch-chip--fg"}, {"Focus", "--ui-color-focus-ring", "#2563eb", "keyboard focus", "theme-gallery-swatch-chip--primary"}}, Fonts: gallerySystemFonts("Base UI-inspired", "System sans / mono fallback", "system fallback stacks; no dedicated webfont pack"), Measures: galleryMeasures("Base UI-inspired"), Mappings: []themeGalleryMapping{{"Implemented directly", "native details/summary Accordion, open/closed state styling, focus-visible, reduced motion, forced colors, and theme tokens"}, {"Adapted / approximated", "Base UI Root/Item/Header/Trigger/Panel composition and state vocabulary without React"}, {"Missing or intentionally out of scope", "disabled summary and hiddenUntilFound polyfills; native browser support is not faked"}}},
}

func gallerySystemFonts(theme, family, status string) []themeGalleryFont {
	return []themeGalleryFont{{family, "System fallback", "400 / 500 / 600", "Platform-provided; no bundled subset", status, "--ui-font-sans; all sans type families"}, {theme + " monospace", "System fallback", "400 / 500", "Platform-provided; no bundled subset", "Used only when the mono token is requested", "--ui-font-mono; code and mono roles"}}
}

func galleryMeasures(theme string) []themeGalleryMeasure {
	return []themeGalleryMeasure{{"Base unit", "--ui-space-1", ".25rem", "Gelium spacing vocabulary"}, {"Control", "--ui-size-control", "theme-defined", "live control geometry"}, {"Field", "--ui-size-field", "theme-defined", "floating-label headroom"}, {"Card radius", "--ui-card-radius", "theme-defined", "component anatomy token"}, {"Container", "--ui-container-max", "64rem", "core max-width contract"}}
}
func galleryMappings(theme string) []themeGalleryMapping {
	return []themeGalleryMapping{{"Implemented directly", "semantic color roles, typescale, controls, fields, cards, badges, status, focus, scheme switching, and a semantic blockquote/figure quote"}, {"Adapted / approximated", "Refero spacing and shape are mapped to Gelium --ui-space/--ui-size/--ui-radius vocabulary; component anatomy keeps Gelium markup"}, {"Missing or intentionally out of scope", "blue wash, logo strip, avatar stack, and product stage remain decorative style-specific patterns, not global primitives"}}
}

var galleryTypeSpecs = []themeGalleryType{{"Display", "--ui-type-display-lg", "Design systems should be legible", "3.5rem", "4rem", "normal"}, {"Headline", "--ui-type-headline-sm", "A live headline specimen", "1.5rem", "2rem", "normal"}, {"Body", "--ui-type-body-md", "The active theme supplies this real body role.", "1rem", "1.5rem", "normal"}, {"Label", "--ui-type-label-md", "TOKEN LABEL", ".875rem", "1.25rem", "normal"}}

func (s *server) docsThemeGallery(w http.ResponseWriter, r *http.Request) {
	slug := themeSlugFromClass(themeClass(themeFromRequest(r)))
	if slug == "" {
		slug = "material"
	}
	g := galleryThemeData[slug]
	g.QuoteText, g.QuoteCite = "A calm interface gives important work room to breathe.", "Gelium design-system audit"
	g.Behavior = accordionBehaviorFromRequest(r)
	g.Execution = accordionExecutionFromRequest(r)
	g.Accordion = accordionForRequest(r)
	g.Types = galleryTypeSpecs
	g.Scheme = schemeFromRequest(r)
	if g.Scheme == "" {
		g.Scheme = "light or OS preference"
	}
	data := pageView{Title: "Theme gallery", ThemeGallery: &g,
		Buttons:    []buttonView{{Label: "Primary", Variant: "primary"}, {Label: "Secondary", Variant: "secondary"}, {Label: "Outline", Variant: "outline"}, {Label: "Disabled", Variant: "primary", Disabled: true}},
		TextFields: []textFieldView{{ID: "gallery-name", Label: "Name", Variant: "outlined", Placeholder: "e.g. Ana López"}, {ID: "gallery-email", Label: "Email", Variant: "filled", Type: "email", Helper: "Native field with floating label."}},
		CardDemo:   &cardDemo{Static: cardDemoCard{Title: "Live component card", Body: "This is the same Gelium card markup used by every theme."}}, BadgeDemo: &badgeDemo{},
		Toasts: []toastView{{ID: "g-toast-info", Type: "info", Role: "status", Message: "New information is available.", IconSVG: toastIcons["info"], Dismiss: true}, {ID: "g-toast-success", Type: "success", Role: "status", Message: "Your changes were saved.", IconSVG: toastIcons["success"], Dismiss: true}},
		Dialog: &dialogView{Trigger: buttonView{Label: "Open dialog", Variant: "primary", Href: "/components/dialog/confirm"}}}
	intro := "# Theme gallery\n\nA live design-system card for the active theme. This live kitchen-sink compares the Refero direction with Gelium’s resolved implementation. Use the **Theme Gallery** link in navigation; the native theme selector and light/dark switcher remain available above.\n"
	s.renderMarkdownStatus(w, r, data, intro, "/docs/themes/gallery", http.StatusOK)
}
