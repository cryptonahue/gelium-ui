# Gelium UI — Phase F Audit: Public Content Patterns (handoff)

> **Alcance**: inventario **read-only** de los 14 patrones public/content de Phase F (roadmap `gelium-ui-system-roadmap.md:191-214`) contra el código real. No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase F :191-214, matriz :423-433), `docs/gelium-ui-vocabulary.md` (:122-138, :239-243), `docs/handoffs/{mozilla-protocol-audit,callout-audit,seo-geo-audit,state-patterns-audit}.md`, `docs/gelium-ui-seo-patterns.md` (P1/P2), `docs/gelium-ui-geo-contract.md` (§9/§14), `web/templates/{layout,card,button,banner,callout,empty-state,demo-whatsapp-admin}.html`, `web/styles/{base,card,banner,callout,tokens,app}.css`, `themes/theme-material/theme.css`, `internal/app/{server,routes,docs}.go`, `web/styles_contract_test.go`, `COMPONENT-ROADMAP.md`, git log.

---

## 1. Resumen ejecutivo

- **Estado real: 4 de 14 existen como primitivas o equivalentes** (Article ≈ `.prose`, CTA Link ≈ Button link, Callout, Banner ≈ Notification Bar), **Card existe sin los slots públicos** (media/tag/aspect/CTA) y **9 son gaps totales**: Breadcrumb, Billboard/Hero, Feature Card, Footer, Language Switcher, Newsletter, Section Heading, Split, Video. Evidencia file:line en §2.
- **Los 9 gaps son 100% estáticos o con contrato server-driven canónico** (GET form / POST + 422), **cero JS de componente**. Coincide con la conclusión del mozilla-protocol-audit (`:142`): los 3 que usan JS en Protocol (Footer, Language Switcher, Newsletter) tienen alternativa nativa (`<details>/<summary>`, GET form con submit visible, POST server con 422).
- **Dos colisiones de naming ya resueltas en fases previas**: Callout (Phase D → tip box) y Notification Bar (≈ Banner, Phase D). Las dos decisiones que quedan para Phase F: **Section Heading = utilidad tipográfica (no componente)** y **Video = contenedor responsive (no componente de contenido)**. Ver §3.
- **El markup canónico del Breadcrumb YA está fijado** en el contrato SEO (`seo-patterns.md:50-64`, patrón P1) y el GEO exige `BreadcrumbList` JSON-LD (`geo-contract.md:94-99,151`). No hay que inventar contrato: solo implementarlo.
- **Hallazgo**: el CSS `.landing-hero` (`base.css:37-38`) es **huérfano** — ningún template lo consume; el home usa `.prose` + `.hero-action` (`layout.html:26-27`). Es la base para el patrón Hero/Billboard de Phase F.
- **Hallazgo**: el `<footer>` que la matriz asigna a Phase E (`roadmap.md:423`) **NO se entregó** en el commit `9504216` (Phase E branding/metadata/a11y); sigue pendiente → **lo toma Phase F**, y es **bloqueante de Phase G** (todas las recipes y del contrato SEO §3).
- **Orden sugerido**: Tier 1 estáticos puros (Breadcrumb, Section Heading, Video, Article, CTA Link) → Tier 2 composiciones (Feature Card, Split, Hero) → Tier 3 con server contract (Footer, Newsletter, Language Switcher) → Tier 4 alias documental (Notification Bar). Detalle en §4.
- **Deliverable de Phase F pendiente**: `docs/gelium-ui-public-content-patterns.md` (salida del roadmap `:212`) **no existe todavía**; este audit es la base para redactarlo.

---

## 2. Los 14 patrones — estado real, contrato propuesto, JS, orden

Leyenda: **✅ existe** · **◐ parcial/composición** · **✖ gap**. Tier de implementación: **1** estático puro · **2** composición · **3** server contract · **4** alias/documental.

| # | Patrón | Estado | Evidencia hoy | Contrato propuesto (resumen) | JS | Tier |
|---|---|---|---|---|---|---|
| 1 | **Article** | ✅ (≈ `.prose`) | `layout.html:26` `<article class="prose">`; `base.css:29-32`; P2 `seo-patterns.md:56-68` | Formalizar contrato tipográfico: `.prose` (h1 display-sm / h2 headline-sm / p body-md muted) + **intro opcional** (párrafo lead `body-lg`). Sin template nuevo | 0 | 1 |
| 2 | **Billboard/Hero** | ◐ (ad-hoc home) | `.hero-action` `layout.html:27` + `base.css:33`; CTA home `server.go:420-424`; **`.landing-hero` huérfano** `base.css:37-38` (sin consumidor) | Composición: `<section class="ui-hero">` + heading `--ui-type-display-lg` + desc + CTA Link + media opcional. Formaliza `.landing-hero` + `.hero-action`. Sin partial obligatorio (slots de layout) | 0 | 2 |
| 3 | **Breadcrumb** | ✖ | Cero matches; gap `vocabulary.md:239-243`, `geo-contract.md:94-99,188`; **contrato ya fijado** `seo-patterns.md:50-64` (P1) | Partial `{{define "breadcrumb"}}`: `<nav aria-label="Breadcrumb"><ol>` + `<li>` con `<a>` / actual = `<span aria-current="page">` (markup P1). Datos: `[]crumb{Path,Label,Current}` derivados de `componentRoutes()`/`navLinks()` (`routes.go:16-58`) | 0 | 1 |
| 4 | **Callout** | ✅ | `callout.html:1-6`; `callout.css` (`app.css:34`); `styles_callout_test.go`; naming resuelto Phase D | **Cerrado** — tip box Gelium (`<aside>`, sin role, variantes default/info/tip, sin tones de estado) | 0 | ✅ |
| 5 | **Card** | ✅ (sin slots públicos) | `card.html:1-22`; `card.css:8-27`; tokens `theme.css:111-115`; `app.css:13` | Extender con **slots públicos**: media con `aspect-ratio`, tag (`--ui-badge` reuso), meta, CTA (`--ui-card-action` ya existe `card.css:27`). Conservar raíz `<article>/<a>/<button>` como control (`vocabulary.md:33`) | 0 | 2 |
| 6 | **CTA Link** | ✅ (≈ Button link) | `button.html:6` variante `Href` → `<a class="ui-button">`; reuso como CTA en `empty-state.html:5`, `banner.html:7`, `callout.html:5` | **Cerrado** por Button link. Opcional (bajo valor): clase inline `.ui-cta-link` con icono para composiciones editoriales | 0 | 1 |
| 7 | **Feature Card** | ✖ | Cero matches; mozilla-audit `:105` (horizontales deprecadas upstream) | **Composición Card + CTA Link** (no primitiva): Card elevada con media + título + desc + `<a class="ui-button">`. NO portar variante horizontal | 0 | 2 |
| 8 | **Footer** | ✖ | Cero `<footer>`; `layout.html:60-61` termina en `main` + `toast-region`; `seo-geo-audit.md:22,150`; matriz `roadmap.md:423` (asignado a E, **no entregado**) | Partial `{{define "footer"}}` + **slot en layout** (`{{if .Footer}}` tras `main`, antes de `toast-region`). `<footer>` + nav secundaria (`h3`+`ul`) + legal; móvil con `<details>/<summary>` nativos. Zero JS | 0 | 3 |
| 9 | **Language Switcher** | ✖ | Cero matches | `<form class="ui-lang-switcher" method="get">` + `<label>` + `<select>` + submit siempre visible. GET `?lang=` → 303 a URL locale; swap `<html lang>` (`layout.html:2`) + RTL. Requiere modelo de locales | 0 (form GET; HTMX opcional) | 3 |
| 10 | **Newsletter** | ✖ | Cero matches (un falso positivo en `checkbox.html`) | `<aside class="ui-newsletter">` + `<form method="post">` email `required type=email` + country/lang + checkbox privacidad `required`. **POST + 422 `X-Gelium-Validation`** + success view (reusa text-field/select/checkbox/button; contratos de `server.go:295-306`). i18n copy | 0 (POST server; HTMX opcional) | 3 |
| 11 | **Notification Bar** | ✅ (≈ Banner) | `banner.html:1-11`; `banner.css` (`app.css:33`); Phase D commit `5787799`; `vocabulary.md:122-129` | **Cerrado** — alias documental: Notification Bar = Banner Gelium (dismiss POST+303 en `banner.html:8-10`). Variantes sticky/scripted **diferidas** | 0 | 4 |
| 12 | **Section Heading** | ✖ | Cero matches; regla h2 en `seo-patterns.md:90` | **Utilidad tipográfica, NO componente**: `.ui-section-heading` (o slot de composición) sobre `--ui-type-headline-sm`/`title-lg` + variante centered. Sin partial, sin handler | 0 | 1 |
| 13 | **Split** | ✖ | Cero matches; "split" en `switch.css`/`segmented-button.css` = clases CSS internas (false positive) | Composición de 2 columnas responsive (stack en mobile): `.ui-split` + `.ui-split-body` + `.ui-split-media`; bidi RTL. CSS + partial mínimo `split.html` opcional | 0 | 2 |
| 14 | **Video** | ✖ | Cero `<video>`; **cero `aspect-ratio`** en todo `web/styles/` + theme | **Contenedor responsive, NO componente de contenido**: `.ui-video` con `aspect-ratio: 16/9` + `<video controls poster loading="lazy">` + `<source>`; variante 4:3. Sin tokens nuevos (aspect no se tokeniza — regla Phase B `roadmap.md:125`); usa `--ui-radius-*` | 0 (controles nativos) | 1 |

---

## 2.1 Contratos detallados para los 9 gaps (siguen las convenciones de partials Phase D)

> Convención Phase D verificada en `callout-audit.md:34`: `web/templates/<x>.html` (`{{define "x"}}`) + `web/styles/<x>.css` (`@layer components`, tokens scoped en root, forced-colors, sin animación → sin reduced-motion) + `@import` en `app.css:33-40` + `web/styles_<x>_test.go` + `sourceAppCSS` (`styles_contract_test.go:22-78`) + `TestComponentSizeTokensDeclaredScoped` + `npm run build` regenera `static/app.css`. Sin handler Go hasta que haya consumidor (primitiva lista para Phase G).

### Breadcrumb (Tier 1 — cero JS, zero server new contract)

```html
{{define "breadcrumb"}}<nav aria-label="Breadcrumb">
  <ol class="ui-breadcrumb">
    {{range .Crumbs}}{{if .Current}}<li><span aria-current="page">{{.Label}}</span></li>
    {{else}}<li><a href="{{.Path}}">{{.Label}}</a></li>{{end}}{{end}}
  </ol>
</nav>{{end}}
```

- **Semántica**: markup canónico **ya fijado** por P1 (`seo-patterns.md:50-64`) y `vocabulary.md:242`; coincide con el de Protocol 1:1 (`mozilla-audit:78`). El actual usa `<span aria-current="page">` (mismo patrón que pagination `data-table.html:68-72`).
- **Tokens**: `--ui-type-body-sm`/`label-sm` (theme :68-70), `--ui-color-primary`/`border-strong` (`tokens.css:26,31`), `--ui-space-*` (:102-107). Separador CSS (no texto "›") para i18n.
- **Datos**: `Crumbs []crumb` construido desde `componentRoutes()` (`routes.go:16-47`) → Home → Docs → Component; misma fuente que el JSON-LD `BreadcrumbList` (P4, `geo-contract.md:151`).
- **Variantes**: dark opcional (`mzp-t-dark`) → diferida hasta consumidor.
- **JS**: 0. **Tests**: ver §5.

### Section Heading (Tier 1 — utilidad, NO partial)

- **Decisión**: utilidad tipográfica, no componente (ver §3). Regla ya fijada: **siempre `h2`, nunca `h1`** (`seo-patterns.md:90`; P2 exige un único h1 por página).
- **Contrato CSS**: `.ui-section-heading` sobre `--ui-type-headline-sm` (`theme.css:63`) + `margin` con `--ui-space-*`; modificador `.ui-section-heading--centered` (`text-align: center`). Sin template — se compone como `<h2 class="ui-section-heading">` en markup de página/landing.

### Video (Tier 1 — contenedor, cero JS)

```html
{{define "video"}}<div class="ui-video{{if eq .Aspect "4:3"}} ui-video--aspect-4-3{{end}}">
  <video controls{{if .Poster}} poster="{{.Poster}}"{{end}} loading="lazy"{{if .Track}} crossorigin="anonymous"{{end}}>
    {{range .Sources}}<source src="{{.Src}}" type="{{.Type}}">{{end}}
    {{if .Track}}<track kind="captions" src="{{.Track.Src}}" srclang="{{.Track.SrcLang}}" label="{{.Track.Label}}">{{end}}
  </video>
</div>{{end}}
```

- **Semántica**: `<video controls>` nativo + `<track kind="captions">` (a11y). Best-used-in: Split, Card, Hero.
- **Tokens**: `aspect-ratio: 16/9` **literal** (no se tokeniza, regla `roadmap.md:125`), `border-radius: var(--ui-radius-sm)` (`tokens.css:80`), `--ui-color-border` para borde del marco.
- **Variantes**: default 16:9, `.ui-video--aspect-4-3`. Sin autoplay (no-JS y a11y).

### Hero / Billboard (Tier 2 — composición; formaliza lo ad-hoc)

- **Estado hoy**: el home ya tiene CTA (`server.go:420-424` → `layout.html:27` `.hero-action`) y `.landing-hero` CSS **huérfano** (`base.css:37-38`). El patrón es: contenido Markdown en `.prose` + CTA.
- **Contrato**: composición sobre el layout o sección dedicada: `h1` (`--ui-type-display-lg`, ya core `tokens.css:144`) + desc (`body-lg`, máx 48ch — `base.css:38`) + CTA Link + media opcional. No crear primitiva: **es composición** (`mozilla-audit:162`). Sin embargo, un partial mínimo `{{define "hero"}}` justifica la reutilización en Landing/Public-Feed (Phase G).
- **RTL**: no se invierte (copy length sensible al idioma, no contrato).

### Feature Card (Tier 2 — composición Card + CTA Link)

- Composición documentada (no primitiva): Card elevada (`--ui-card-container-elevated` `theme.css:112`) + media + título + desc + CTA Link. NO portar variante horizontal (deprecada upstream, `mozilla-audit:105,186`). Requiere primero los slots públicos de Card (§2.1 Card).

### Split (Tier 2 — CSS de composición, partial mínimo opcional)

- **Contrato**: `.ui-split` grid `grid-template-columns: repeat(2, minmax(0,1fr))` → colapsa a 1 col en narrow (stack); `.ui-split-body` (tipografía NO aplicada por defecto, precedente Protocol `mozilla-audit:117`) + `.ui-split-media` (img/video). **Bidi RTL automático** (requisito del sistema; usará `direction`/flexbox, no `left/right` literales).
- **No-nos** (Protocol `mozilla-audit:117,187`): no anidar en contenedor de contenido propio → regla de composición.

### Footer (Tier 3 — server contract de datos, cero JS)

```html
{{define "footer"}}<footer class="ui-footer">
  <nav aria-label="Footer" class="ui-footer-nav">
    {{range .Sections}}<section class="ui-footer-section">
      <details class="ui-footer-details">
        <summary class="ui-footer-heading">{{.Title}}</summary>
        <ul class="ui-footer-list">{{range .Links}}<li><a href="{{.Path}}">{{.Label}}</a></li>{{end}}</ul>
      </details>
    </section>{{end}}
  </nav>
  <div class="ui-footer-legal">{{.Legal}}</div>
</footer>{{end}}
```

- **Semántica**: `<footer>` + nav secundaria + legal; `<details>/<summary>` nativos para el accordion móvil (sin `open` → collapsed por defecto; desktop lo muestra abierto por CSS). Alternativa nativa a `MzpFooter`+`MzpDetails` (JS en Protocol).
- **Datos**: `footerView{Sections []{Title, Links []navLink}, Legal string}` — reusa `navLink` (`routes.go:55-58`); puede derivar secciones de `docsSections` (`docs.go:9-78`) o un modelo propio del consumidor.
- **Slot en layout**: `{{if .Footer}}{{template "footer" .Footer}}{{end}}` tras `</main>` (`layout.html:61`) y antes de `{{template "toast-region"}}`.
- **Tokenización**: puede reusar `--ui-color-surface`/`border`, `--ui-type-label-lg`/`body-sm` (`theme.css:69,68`). Sin tokens nuevos obligatorios.
- **Server contract**: solo datos (GET); el `<footer>` vive en todas las páginas del layout → **bloqueante de Phase G** y del contrato SEO §3 (`seo-geo-audit.md:22,150`).

### Newsletter (Tier 3 — POST + 422, el único form transaccional del set)

```html
{{define "newsletter"}}<aside class="ui-newsletter" aria-labelledby="newsletter-heading">
  <h2 id="newsletter-heading" class="ui-newsletter-title">{{.Title}}</h2>
  {{if .Success}}<p class="ui-newsletter-success" role="status">{{.Success}}</p>
  {{else}}<form class="ui-newsletter-form" method="post" action="{{.Action}}">
    {{template "text-field" .EmailField}}
    {{template "select" .CountryField}}
    {{template "checkbox" .PrivacyField}}
    {{template "button" .SubmitButton}}
  </form>{{end}}
</aside>{{end}}
```

- **Contrato server (canónico, no inventar)**: POST → valida → **422 + `X-Gelium-Validation`** (contrato a, `composition-audit` §9; precedente `text_field.go`) con re-render del aside mostrando inline-alert; éxito → **POST + 303 → GET re-renderiza vista success persistente** (contrato d; precedente `banner.html:8-10`, `chipsRemoveDemo`). HTMX opcional.
- **Composición**: reusa text-field/select/checkbox/button reales (`{{template "text-field" .}}` etc.) — no duplica markup.
- **i18n**: copy, países/idiomas, privacy (accesible), `required` + `type=email`.
- **Tokens**: `--ui-color-surface`, `--ui-space-*`, reuso `--ui-text-field-*`/`--ui-button-*`.

### Language Switcher (Tier 3 — GET form + 303; requiere modelo de locales)

```html
{{define "lang-switcher"}}<form class="ui-lang-switcher" method="get" action="{{.Action}}">
  <label class="ui-lang-switcher-label" for="{{.SelectID}}">{{.Label}}</label>
  <select class="ui-lang-switcher-select" id="{{.SelectID}}" name="lang">{{range .Locales}}
    <option value="{{.Code}}"{{if .Current}} selected{{end}}>{{.Name}}</option>{{end}}</select>
  <button class="ui-button ui-button-outlined" type="submit">{{.SubmitLabel}}</button>
</form>{{end}}
```

- **Contrato server**: GET `?lang=<code>` → **303 a la URL localizada** (contrato c: GET con params estables; 303 redirección). El `<html lang>` (`layout.html:2`) y RTL se resuelven server-side con el locale.
- **No-JS**: submit siempre visible (fallback nativo de Protocol `mozilla-audit:109`); auto-submit JS **diferido**.
- **i18n**: label localizado, primer path de URL como locale (convención `mozilla-audit:109`).
- **Ubicación**: dentro del `<footer>` (nav secundaria, precedente Protocol) o en el header.

---

## 3. Decisiones de naming (pendientes del mozilla-protocol-audit `:192-199`)

| Término | Decisión | Estado | Fuente |
|---|---|---|---|
| **Callout** | **Resuelto (Phase D)**: canónico Gelium = tip box (`<aside>`, contenido ignorable). El "Callout" de Protocol (hero full-width) NO se implementa con ese nombre → queda cubierto por Hero/Billboard (composición). | ✅ | `callout-audit.md:38-51`, `callout.html`, `mozilla-audit.md:196` |
| **Notification Bar** | **Resuelto**: ≈ **Banner Gelium** (persistente página/sitio, `banner.html`). Notification Bar queda como **alias documental** del Banner; sticky/scripted diferidos. | ✅ | `vocabulary.md:122-129`, Phase D commit `5787799`, `mozilla-audit.md:197` |
| **Section Heading** | **Es utilidad tipográfica, NO componente**: `.ui-section-heading` sobre `--ui-type-*`, variante centered; sin partial ni handler. Regla: siempre `h2`, nunca `h1` (`seo-patterns.md:90`). | 🆕 decidir en F | `mozilla-audit.md:115`, `seo-patterns.md:90` |
| **Video** | **Es contenedor responsive, NO componente de contenido**: `.ui-video` + `aspect-ratio` + `<video controls>` nativo; "best used inside another component" (Split/Card/Hero). | 🆕 decidir en F | `mozilla-audit.md:119` |

**Otras decisiones fijadas** (no reabrir): CTA Link = Button link (`button.html:6`); Article = `.prose` formalizado; Feature Card = composición Card+CTA (no primitiva); Split = composición de grid; Hero = composición.

---

## 4. Orden de implementación sugerido

### Tier 1 — 100% estáticos, cero server contract, paralelizables (primero)
1. **Breadcrumb** — partial + CSS + tests; markup ya canónico (P1). **Desbloquea GEO §9/§14 y seo-contract §11.**
2. **Section Heading** — utilidad CSS (sin partial); 30 min.
3. **Video** — contenedor + partial mínimo; cero JS.
4. **Article** — formalizar contrato `.prose` (tipografía + intro opcional); doc + posible utilidad CSS.
5. **CTA Link** — solo si hace falta la variante inline con icono (opcional; Button link ya cubre).

### Tier 2 — composiciones de existentes (después de Card slots)
6. **Card → slots públicos** (media + aspect-ratio, tag `--ui-badge`, meta, CTA `--ui-card-action` ya existe). Requisito de Feature Card y Split-media.
7. **Feature Card** = Card + CTA Link (composición, CSS mínimo).
8. **Split** = grid 2 col (body + media), RTL bidi; partial opcional.
9. **Hero/Billboard** = composición que formaliza `.landing-hero` huérfano + `.hero-action`; útil para Landing (Phase G).

### Tier 3 — server contract (necesitan handler Go + datos)
10. **Footer** — partial + slot en layout + modelo de datos; `<details>/<summary>`. **Bloqueante Phase G.**
11. **Newsletter** — POST + 422 `X-Gelium-Validation` + success view (reusa text-field/select/checkbox/button).
12. **Language Switcher** — GET form + 303; requiere modelo de locales (`lang`/RTL). Se compone dentro del Footer.

### Tier 4 — alias documental (sin código)
13. **Notification Bar** → alias de Banner en `gelium-ui-vocabulary.md` (+ doc en el deliverable F).

**Regla de paralelismo** (`COMPONENT-ROADMAP.md:281-307`): Tier 1 es seguro en paralelo (archivos exclusivos); Tier 3 toca `layout.html` + `server.go` + `routes.go` → lane serial del integrador. Un solo agente por archivo compartido.

---

## 5. Tests necesarios por patrón

Convención Phase D (`callout-audit.md:153-162`, `styles_banner_test.go` como plantilla):

| Patrón | Tests requeridos |
|---|---|
| **Breadcrumb** | `web/styles_breadcrumb_test.go`: contrato CSS (`.ui-breadcrumb`, `ol/li/a`, separador sin texto, tokens core, forced-colors, sin animation). Render Go: `<nav aria-label="Breadcrumb">` → `<ol>` → `<li>`; actual = `<span aria-current="page">` sin `<a>`; cero links para `Current`. `sourceAppCSS` + `TestComponentSizeTokensDeclaredScoped`. |
| **Section Heading** | Sin test de template (no hay partial); test CSS mínimo si se crea `.ui-section-heading` (presencia en compiled `app.css`). |
| **Video** | `styles_video_test.go`: `aspect-ratio: 16/9`, variante 4:3, radius core, forced-colors. Render: `<video controls>`, `poster`, `loading="lazy"`, `<source>`, `<track kind="captions">`, sin autoplay. |
| **Hero** | Test de composición (reuso): home renderiza CTA `<a class="ui-button">` (ya cubierto por `server_test.go`); aserción de que `.landing-hero` CSS existe en compiled (o se elimina el huérfano). |
| **Feature Card / Split** | Test de composición documental (no primitiva): clases CSS mínimas si se crean; pin de reuso Card+CTA. |
| **Footer** | `styles_footer_test.go`: `<details>` sin `open` por defecto, heading `h3`/`summary`, listas, forced-colors. Render Go: `<footer>`, nav secundaria, legal. Layout test: slot `{{if .Footer}}` no rompe páginas sin Footer (nil-safe). |
| **Newsletter** | Contrato server (Go): POST sin email → **422 + `X-Gelium-Validation`** + re-render con inline-alert; POST válido → 303 → GET success view `role="status"`; `type=email` + `required` en markup. CSS: `styles_newsletter_test.go`. |
| **Language Switcher** | Contrato server: GET `?lang=` → 303 a URL locale; submit `<button type="submit">` siempre presente (no-JS); `select` con `selected` por locale; swap de `<html lang>`. |
| **Card slots** | Extender `styles_card_test.go` + `card_test.go`: media `aspect-ratio`, tag, meta, CTA presente/ausente según datos. |
| **Common (todos)** | `sourceAppCSS` sincronizado (`styles_contract_test.go:22-78`); `TestComponentSizeTokensDeclaredScoped` para tokens scoped de tamaño (si aplica); `npm run build` + `go test ./...` + `go vet ./...` (`roadmap.md:446-466`); light/dark, narrow/wide, RTL (Split/LangSwitcher), forced-colors, teclado, no-JS end-to-end, HTMX solo enhancement. |

---

## 6. Archivos impactados (plan read-only — NO modificados en este audit)

**Nuevos**:
- `web/templates/breadcrumb.html`, `web/templates/video.html`, `web/templates/footer.html` (y opcional `split.html`, `hero.html`, `newsletter.html`, `lang-switcher.html`).
- `web/styles/breadcrumb.css`, `web/styles/video.css`, `web/styles/footer.css` (y opcional `split.css`, `hero.css`, `newsletter.css`, `lang-switcher.css`, `section-heading` utilidad).
- `web/styles_breadcrumb_test.go`, `web/styles_video_test.go`, `web/styles_footer_test.go`, `web/styles_newsletter_test.go`, etc.
- `docs/gelium-ui-public-content-patterns.md` — **deliverable de Phase F** (roadmap `:212`), hoy inexistente.
- `docs/handoffs/public-patterns-audit.md` — este handoff.

**Modificados** (cuando se implemente, no ahora):
- `web/templates/layout.html` — slot `{{if .Footer}}` tras `</main>` (`:61`); opcional slot breadcrumb/hero.
- `web/styles/app.css` — `@import` de los nuevos CSS (tras `validation-summary.css` `:40`, junto a la capa public/content; NO tocar tokens/theme).
- `web/static/app.css` — regenerado por `npm run build` (precedente commits Phase D).
- `web/styles_contract_test.go` — `sourceAppCSS` (`:22-78`) + `TestComponentSizeTokensDeclaredScoped`.
- `internal/app/server.go` — `pageView` (`:228-268`) + slots `Footer *footerView`, `Breadcrumb []crumb`; handlers Newsletter/LangSwitcher (POST/GET).
- `internal/app/routes.go` — rutas `POST /newsletter/subscribe` (o similar) + GET 405 companion (`postOnlyPaths()` `server.go:329-342`).
- `docs/gelium-ui-vocabulary.md` — marcar aliases: Notification Bar ≈ Banner; Section Heading utilidad; Video contenedor; Breadcrumb ✖→✅.
- `docs/gelium-ui-seo-patterns.md` / `geo-contract.md` — ya referencian P1/P4 (breadcrumb): sin cambios de contrato.

**No tocados**: `web/static/app.js` (cero JS), `web/styles/tokens.css` + `themes/theme-material/theme.css` (tipografía y colores existentes; Video usa `aspect-ratio` literal por regla `roadmap.md:125`), partials Phase D existentes.

---

## 7. Desbloqueo de Phase G (screen recipes)

| Patrón Phase F | Desbloquea en Phase G | Evidencia |
|---|---|---|
| **Breadcrumb** | **Settings** (recipe G) y TODAS las páginas `/components/*` (GEO §9: breadcrumb visible + `BreadcrumbList` JSON-LD) | `geo-contract.md:94-99,151`; `seo-contract.md:296-301` |
| **Footer** | **Todas las recipes** (chrome de sitio: nav secundaria + legal + locale) y contrato SEO §3 (footer faltante) | `seo-geo-audit.md:22,150`; `roadmap.md:423` |
| **Hero/Billboard** | **Landing/Public-Social Feed** (recipe G-3); formaliza el home actual (`server.go:420-424`, `.landing-hero` huérfano) | `roadmap.md:430` |
| Card slots públicos | **Admin Resource, Public Feed** (cards de contenido editorial con media/tag/CTA) | `roadmap.md:434,436` |
| Video / Split | **Public Feed / Landing** (media promocional); opcional en las 3 primeras recipes | `roadmap.md:430` |

**No bloquean G**: Newsletter y Language Switcher (conversión de audiencia de un sitio público multilingüe; no necesarias para las 3 primeras recipes internas). **Section Heading** desbloquea el encabezado de secciones de cualquier recipe (utilidad, barata).

---

## Fuentes

- `docs/gelium-ui-system-roadmap.md` — Phase F `:191-214`; matriz `:423-433`; dependencias `:315-330`; reglas no-tokenizar `:125`; verificación `:446-466`.
- `docs/handoffs/mozilla-protocol-audit.md` — mapping de los 14 `:74-119`; clasificación `:146-177`; divergencias `:180-199`; top 5 `:215-221`.
- `docs/handoffs/callout-audit.md` — naming resuelto `:38-51`; convención de partials `:34,153-162`.
- `docs/handoffs/seo-geo-audit.md` — footer faltante `:22,150`.
- `docs/gelium-ui-vocabulary.md` — Banner `:122-129`, Callout `:131-138`, Breadcrumbs `:239-243`.
- `docs/gelium-ui-seo-patterns.md` — P1 breadcrumb `:50-64`; P2 single-h1 `:66-92`.
- `docs/gelium-ui-geo-contract.md` — §9 Relations `:94-99`; §14 `BreadcrumbList` `:151`.
- Código: `web/templates/layout.html:26-27,60-61`, `web/templates/{card,button,banner,callout,empty-state}.html`, `web/styles/{base,card,banner,callout,tokens,app}.css`, `themes/theme-material/theme.css:61-70,111-115`, `internal/app/{server,routes,docs}.go`, `web/styles_contract_test.go:22-78`, `COMPONENT-ROADMAP.md`.
- Git: commits Phase D (`eba1c4c`…`9504216`), `9504216` (Phase E — sin footer).
