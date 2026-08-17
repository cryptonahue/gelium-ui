package app

import (
	"bytes"
	"html/template"
	"net/http"

	"geliumui/lib"
)

// landingView is the marketing home composition. Every section dogfoods an
// existing Gelium public/content pattern or component — no parallel landing CSS
// system beyond a thin page frame.
type landingView struct {
	Hero            *heroView
	Claims          []string
	FeaturesHeading sectionHeadingView
	Features        []featureCardView
	Split           *splitView
	Demo            *landingDemoView
	Recipes         *landingRecipesView
	FAQ             *landingFAQView
	CTABand         *landingCTABandView
}

// heroView matches web/templates/hero.html.
type heroView struct {
	Eyebrow  string
	Title    string
	Subtitle string
	CTAs     []buttonView
	Media    template.HTML
}

// featureCardView matches web/templates/feature-card.html.
type featureCardView struct {
	Title string
	Body  string
	CTA   *buttonView
	Media template.HTML
}

// splitView matches web/templates/split.html.
type splitView struct {
	Eyebrow string
	Title   string
	Body    string
	CTA     *buttonView
	Media   template.HTML
}

// sectionHeadingView matches web/templates/section-heading.html.
type sectionHeadingView struct {
	Eyebrow  string
	Title    string
	Centered bool
}

// landingRecipesView is a section of outbound recipe cards.
type landingRecipesView struct {
	Heading sectionHeadingView
	Cards   []featureCardView
}

// landingCTABandView is the closing call-to-action band.
type landingCTABandView struct {
	Heading sectionHeadingView
	CTA     buttonView
	// SecondaryCTA is an optional secondary action rendered beside the primary
	// (e.g. the GitHub source link). It is intentionally NOT chrome-rewritten:
	// external hrefs must never carry ?theme= / ?scheme= query strings.
	SecondaryCTA *buttonView
}

// landingDemoView is the visual demo card (Basecoat pattern): a prominent
// link into the live WhatsApp demo with a token-styled preview.
type landingDemoView struct {
	Heading sectionHeadingView
	Body    string
	CTA     buttonView
	Media   template.HTML
}

// landingFAQView is the closing FAQ section (Base UI pattern): answers
// developer objections with zero-JS native disclosures.
type landingFAQView struct {
	Heading sectionHeadingView
	Items   []faqItemView
}

// faqItemView is a single question/answer pair.
type faqItemView struct {
	Question string
	Answer   string
}

// marketingLanding builds the home page composition from Gelium primitives.
// Persuade mode: exactly one primary button on the page (hero Get started).
func marketingLanding() landingView {
	return landingView{
		Hero: &heroView{
			Eyebrow:  "Gelium UI",
			Title:    "Server-rendered components. Zero required JS.",
			Subtitle: "Open-code HTML and tokens for server apps. Install npm gelium-ui, pick a theme class, copy partials. HTMX when you want it.",
			CTAs: []buttonView{
				{Label: "Get started", Variant: "primary", Href: "/docs"},
				{Label: "Browse components", Variant: "secondary", Href: "/components/button"},
			},
		},
		Claims: []string{
			"Zero required JS",
			"Two themes, one bundle",
			"735+ contract tests",
			"Open code, MIT",
		},
		FeaturesHeading: sectionHeadingView{
			Eyebrow:  "Why Gelium",
			Title:    "Built for server-rendered apps",
			Centered: true,
		},
		Features: []featureCardView{
			{
				Title: "Native HTML first",
				Body:  "Real buttons, dialogs, selects, and tables. ARIA only when the platform has no equivalent.",
				CTA:   &buttonView{Label: "Read principles", Variant: "outline", Href: "/docs/principles"},
			},
			{
				Title: "Tokens, not forks",
				Body:  "Themes map aesthetics onto --ui-* variables. Markup stays stable when you switch Material or Basecoat.",
				CTA:   &buttonView{Label: "Themes", Variant: "outline", Href: "/docs/themes"},
			},
			{
				Title: "Server contracts",
				Body:  "GET with stable query params, POST + 303, 422 validation, and gelium:toast for transient feedback — no parallel APIs.",
				CTA:   &buttonView{Label: "See a recipe", Variant: "outline", Href: "/recipes/admin-resource"},
			},
		},
		Split: &splitView{
			Eyebrow: "How it fits",
			Title:   "Install the package. Copy the open code.",
			Body:    "Consumers use npm gelium-ui (CSS, themes, templates, JS helpers). This docs site is a Go dogfood app — not the install path for product UI. Progressive enhancement stays optional.",
			CTA:     &buttonView{Label: "Open the docs", Variant: "secondary", Href: "/docs"},
			Media: template.HTML(
				`<pre class="ui-landing-code" tabindex="0"><code>npm install gelium-ui

/* CSS */
@import "gelium-ui/dist/gelium.css";

/* Theme on &lt;html&gt; */
&lt;html class="theme-material"&gt;

/* Optional */
// gelium.js — toast + 422 helper</code></pre>`,
			),
		},
		Demo: &landingDemoView{
			Heading: sectionHeadingView{
				Eyebrow:  "See it live",
				Title:    "A real app, running on Gelium",
				Centered: true,
			},
			Body: "The WhatsApp manager demo — chat list, search, send, and templates — is a complete server-rendered app built from the same public components and tokens you install.",
			CTA:  buttonView{Label: "Launch live demo", Variant: "secondary", Href: "/demo/whatsapp"},
			Media: template.HTML(
				`<div class="ui-landing-demo-phone" aria-hidden="true">
  <div class="ui-landing-demo-phone-bar"></div>
  <div class="ui-landing-demo-chat">
    <div class="ui-landing-demo-bubble ui-landing-demo-bubble--in"></div>
    <div class="ui-landing-demo-bubble ui-landing-demo-bubble--out"></div>
    <div class="ui-landing-demo-bubble ui-landing-demo-bubble--in ui-landing-demo-bubble--short"></div>
  </div>
</div>`,
			),
		},
		Recipes: &landingRecipesView{
			Heading: sectionHeadingView{
				Eyebrow:  "Screen recipes",
				Title:    "Full flows composed from the catalog",
				Centered: true,
			},
			Cards: []featureCardView{
				{
					Title: "Admin Resource",
					Body:  "Table, filters, create/edit forms, and delete confirm — server-driven end to end.",
					CTA:   &buttonView{Label: "Open recipe", Variant: "outline", Href: "/recipes/admin-resource"},
				},
				{
					Title: "Ops Queue",
					Body:  "FIFO work queue with advance/dequeue actions and refresh without a SPA.",
					CTA:   &buttonView{Label: "Open recipe", Variant: "outline", Href: "/recipes/ops-queue"},
				},
				{
					Title: "Public Feed",
					Body:  "Social-style feed composition with reactions and progressive enhancement.",
					CTA:   &buttonView{Label: "Open recipe", Variant: "outline", Href: "/recipes/public-feed"},
				},
			},
		},
		FAQ: &landingFAQView{
			Heading: sectionHeadingView{
				Eyebrow:  "FAQ",
				Title:    "Frequently asked questions",
				Centered: true,
			},
			Items: []faqItemView{
				{
					Question: "Is Gelium UI a framework?",
					Answer:   "No — the UI package is npm gelium-ui (HTML/CSS/tokens). This documentation site is a Go dogfood app; there is no SPA runtime to install and no service to run for components.",
				},
				{
					Question: "Does it work with my existing stack?",
					Answer:   "Yes — components are HTML partials and CSS, so any server-rendered stack can use them. HTMX stays optional: the core contract is plain forms, links, and progressive enhancement.",
				},
				{
					Question: "How does it differ from a React component library?",
					Answer:   "No component JavaScript and no framework runtime for consumers — Gelium ships server contracts instead of props, and the server renders the HTML.",
				},
				{
					Question: "Can I switch themes?",
					Answer:   "Yes — Material is the default and Basecoat ships in the same bundle. Switch with ?theme=basecoat on any URL or a class on the html element; no rebuild required.",
				},
			},
		},
		CTABand: &landingCTABandView{
			Heading: sectionHeadingView{
				Eyebrow:  "Ready when you are",
				Title:    "Ship UI that works without JavaScript",
				Centered: true,
			},
			// Secondary only — hero already owns the single primary (Persuade / Screens).
			CTA: buttonView{Label: "Read the docs", Variant: "secondary", Href: "/docs"},
			SecondaryCTA: &buttonView{
				Label:   "View source",
				Variant: "outline",
				Href:    "https://github.com/cryptonahue/gelium-ui",
			},
		},
	}
}

// homeLandingNav is the compact primary nav for the marketing site chrome.
func homeLandingNav() []navLink {
	return []navLink{
		{Path: "/docs", Label: "Docs"},
		{Path: "/components/button", Label: "Components"},
		{Path: "/recipes/admin-resource", Label: "Recipes"},
		{Path: "/docs/agent-workflow", Label: "Agents"},
		{Path: "/demo/whatsapp", Label: "Demo"},
	}
}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	landing := marketingLanding()
	data := pageView{
		Title:   "Themeable open-code UI components for server-rendered apps",
		Landing: &landing,
		Nav:     homeLandingNav(),
	}
	s.renderLanding(w, r, data)
}

// renderLanding renders the marketing home (or any pageView with Landing set)
// through the shared layout without Markdown content.
func (s *server) renderLanding(w http.ResponseWriter, r *http.Request, data pageView) {
	routePath := "/"
	data.Meta = resolveMeta(data, routePath)
	data.AssetsVersion = lib.AssetsVersion
	if data.Footer == nil {
		data.Footer = defaultFooter()
	}

	themeSlug := ""
	if q := themeFromRequest(r); q != "" {
		data.ThemeClass = q
		themeSlug = themeSlugFromClass(q)
	} else {
		data.ThemeClass = themeClass(data.ThemeClass)
	}
	scheme := schemeFromRequest(r)
	applyDocumentRootScheme(&data, scheme)

	// Site chrome on the landing: same 0-JS theme + appearance controls as docs.
	data.ThemeSwitcher = themeSwitcherFor(r, data.ThemeClass, themeSlug, scheme)
	data.SchemeSwitcher = schemeSwitcherFor(r, themeSlug, scheme)

	// Rewrite compact nav hrefs so theme/scheme survive header clicks.
	if themeSlug != "" || normalizeScheme(scheme) != "" {
		for i := range data.Nav {
			data.Nav[i].Path = chromeHref(data.Nav[i].Path, themeSlug, scheme)
		}
		// Hero/feature CTAs are absolute docs paths — preserve chrome query too.
		if data.Landing != nil {
			applyLandingChrome(data.Landing, themeSlug, scheme)
		}
	}

	var page bytes.Buffer
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(page.Bytes())
}

// applyLandingChrome appends allowlisted theme/scheme query onto landing CTAs.
func applyLandingChrome(l *landingView, themeSlug, scheme string) {
	if l.Hero != nil {
		for i := range l.Hero.CTAs {
			if l.Hero.CTAs[i].Href != "" {
				l.Hero.CTAs[i].Href = chromeHref(l.Hero.CTAs[i].Href, themeSlug, scheme)
			}
		}
	}
	for i := range l.Features {
		if l.Features[i].CTA != nil && l.Features[i].CTA.Href != "" {
			l.Features[i].CTA.Href = chromeHref(l.Features[i].CTA.Href, themeSlug, scheme)
		}
	}
	if l.Split != nil && l.Split.CTA != nil && l.Split.CTA.Href != "" {
		l.Split.CTA.Href = chromeHref(l.Split.CTA.Href, themeSlug, scheme)
	}
	if l.Demo != nil && l.Demo.CTA.Href != "" {
		l.Demo.CTA.Href = chromeHref(l.Demo.CTA.Href, themeSlug, scheme)
	}
	if l.Recipes != nil {
		for i := range l.Recipes.Cards {
			if l.Recipes.Cards[i].CTA != nil && l.Recipes.Cards[i].CTA.Href != "" {
				l.Recipes.Cards[i].CTA.Href = chromeHref(l.Recipes.Cards[i].CTA.Href, themeSlug, scheme)
			}
		}
	}
	if l.CTABand != nil && l.CTABand.CTA.Href != "" {
		l.CTABand.CTA.Href = chromeHref(l.CTABand.CTA.Href, themeSlug, scheme)
	}
}
