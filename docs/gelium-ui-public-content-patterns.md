# Gelium UI — Public Content Patterns (Phase F)

> Salida de la Fase F del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Clasifican como **public/content patterns, NO como theme**. Ninguno exige JS para su esencia.
> Base: `docs/handoffs/public-patterns-audit.md`, `docs/gelium-ui-seo-patterns.md` (P1/P2), `docs/gelium-ui-vocabulary.md`.
>
> Leyenda de estado: **✅ implementado** · **◐ parcial** (composición ad-hoc) · **✖ no existe**.

---

## 1. Matriz de los 14 patrones

| # | Patrón | Estado | Tier | Implementación |
|---|--------|--------|------|----------------|
| 1 | **Article** | ✅ (≈ `.prose`) | 1 | `base.css` `.prose` + `layout.html:26`; contrato tipográfico + intro opcional (P2) |
| 2 | **Billboard/Hero** | ✅ | 2 | `hero.html`/`hero.css` + `styles_hero_test.go`; composición: `h1` `--ui-type-display-lg` + subtitle + CTA(s) (Button) + media de fondo opcional con scrim |
| 3 | **Breadcrumb** | ✅ | 1 | `breadcrumb.html`/`breadcrumb.css` + `styles_breadcrumb_test.go`; markup canónico P1; GEO §9/§14 |
| 4 | **Callout** | ✅ | — | `callout.html`/`callout.css` (Phase D); naming resuelto (tip box) |
| 5 | **Card (slots públicos)** | ✅ | 2 | Primitiva `{{define "card"}}` en `card.html`: slots opcionales media (`aspect-ratio` literal 16/9, `object-fit: cover`), tag (pill `--ui-badge`), meta (`--ui-space-*`), CTA (`--ui-card-action` + Button); `card-demo` sin cambios; + `styles_card_test.go` |
| 6 | **CTA Link** | ✅ (formalizado = Button link) | 1 | `button.html:6` variante `Href` → `<a class="ui-button">`; **no componente propio**; reusado por Empty state/Banner/Callout/Hero/Feature Card/Split |
| 7 | **Feature Card** | ✅ | 2 | `feature-card.html`/`feature-card.css` + `styles_feature_card_test.go`; composición Card + media + CTA Link (no primitiva); sin variante horizontal (deprecada upstream) |
| 8 | **Footer** | ✅ | 3 | `footer.html`/`footer.css` + slot en `layout.html` + `pageView.Footer`/`defaultFooter()` en `server.go`; `<details>/<summary>`; bloqueante Phase G |
| 9 | **Language Switcher** | ✅ | 3 | `language-switcher.html`/`language-switcher.css` + `styles_language_switcher_test.go`; GET `?lang=` → 303 a URL localizada; submit visible (no auto-submit JS) |
| 10 | **Newsletter** | ✅ | 3 | `newsletter.html`/`newsletter.css` + `styles_newsletter_test.go` + handler `internal/app/newsletter.go`; POST + 422 `X-Gelium-Validation` + success view (reusa inline-alert) |
| 11 | **Notification Bar** | ✅ (≈ Banner) | 4 | Alias documental del Banner Gelium (`banner.html`, Phase D); variantes sticky/scripted diferidas |
| 12 | **Section Heading** | ✅ | 1 | Utilidad tipográfica: `section-heading.html`/`section-heading.css`; siempre `h2`, nunca `h1` |
| 13 | **Split** | ✅ | 2 | `split.html`/`split.css` + `styles_split_test.go`; grid 2 col (`.ui-split` + `.ui-split-body` + `.ui-split-media`); stack en narrow; bidi RTL automático |
| 14 | **Video** | ✅ | 1 | Contenedor responsive: `video.html`/`video.css`; `aspect-ratio` literal 16:9 (no se tokeniza); `<video controls>` nativo |

**Resumen**: 14 de 14 ✅ — Phase F completo. Esta entrega añade los 5 restantes (**Billboard/Hero, Feature Card, Language Switcher, Newsletter, Split**) + formaliza **CTA Link** (Button link, sin componente propio) + cierra **Card slots públicos** (primitiva `{{define "card"}}` con media/tag/meta/CTA opcionales).

---

## 2. Orden de implementación (estado: completa)

### Tier 1 — 100% estáticos, cero server contract

1. **Breadcrumb** ✅ — partial + CSS + tests; markup canónico (P1). Desbloquea GEO §9/§14.
2. **Section Heading** ✅ — utilidad CSS + partial mínimo.
3. **Video** ✅ — contenedor + partial mínimo; cero JS.
4. **Article** ✅ — formalizado como `.prose` (contrato tipográfico); sin partial nuevo.
5. **CTA Link** ✅ — **formalizado = Button link** (`button.html` variante `Href`); variante inline con icono opcional (no implementada, no requerida).

### Tier 2 — composiciones de existentes

6. **Card → slots públicos** ✅ — primitiva `{{define "card"}}` en `card.html` (Phase F): slots opcionales media (`aspect-ratio` literal, nunca tokenizado), tag (pill `--ui-badge`), meta, CTA (`--ui-card-action`); `card-demo` sin cambios.
7. **Feature Card** ✅ — `feature-card.html` + `feature-card.css`: composición Card + CTA Link (CSS mínimo, solo media aspect + spacing).
8. **Split** ✅ — `split.html` + `split.css`: grid 2 col (body + media), stack en narrow, bidi RTL automático.
9. **Hero/Billboard** ✅ — `hero.html` + `hero.css`: composición display `h1` + subtitle + CTA(s) + media de fondo; formaliza el `h1` de landing.

### Tier 3 — server contract

10. **Footer** ✅ — partial + slot en layout + modelo de datos (`footerView`); `<details>/<summary>`. **Bloqueante Phase G (entregado).**
11. **Newsletter** ✅ — `newsletter.html`/`newsletter.css` + handler `internal/app/newsletter.go`: **POST + 422 `X-Gelium-Validation`** (header real del código) + success view persistente; ejemplo en `GET/POST /examples/newsletter` (noindex).
12. **Language Switcher** ✅ — `language-switcher.html`/`language-switcher.css`: **GET form + submit visible**, cero auto-submit JS. El **modelo de locales** (`?lang=` → 303 a URL localizada, swap `<html lang>`/RTL) queda **fuera de alcance** (no hay i18n real todavía) — el patrón es la primitiva lista; el server debe resolver `?lang=`.

### Tier 4 — alias documental

13. **Notification Bar** ✅ — alias de Banner en `gelium-ui-vocabulary.md`.

---

## 2.1 Decisiones de naming (Phase F)

| Decisión | Estado |
|---|---|
| **Hero NO se llama Callout** — el "Callout" de Protocol (hero full-width) queda cubierto por `ui-hero`; el Callout Gelium (Phase D) es el tip box `<aside>` ignorable. Naming cerrado. | ✅ |
| **Feature Card = composición**, no primitiva: `ui-feature-card` envuelve `.ui-card` + media + `.ui-card-title`/`.ui-card-body` + `.ui-card-action` (CTA Link). Sin variante horizontal (deprecada upstream). | ✅ |
| **Language Switcher = GET navigation form**: `method="get"` + `<select name="lang">` + submit visible; el cambio de idioma es navegación (server responde 303 a la URL localizada), nunca POST. Cero auto-submit JS. | ✅ |
| **Newsletter = POST + 422** con el header real del código **`X-Gelium-Validation: true`** (el header real es `X-Gelium-Validation`; ver `screen-recipes-audit.md:17`). Error reusa `inline-alert--error`; success = texto `role="status"` persistente. | ✅ |
| **CTA Link = Button link** — no se crea componente propio; `{{template "button" .CTA}}` con `Href` es la forma canónica. | ✅ |
| **Card slots públicos** — cerrado: primitiva `{{define "card"}}` con slots media/tag/meta/CTA opcionales (data-driven, `{{if}}` guards); demo sin cambios. | ✅ |

---

## 3. Reglas

- **Zero-JS end-to-end**: Footer → `<details>/<summary>`; Language Switcher → GET form con submit visible; Newsletter → POST + 422 (HTMX solo como enhancement: `hx-post` + swap del fragmento del aside).
- **No-tokenizar**: `aspect-ratio` (Video, Feature Card media, Card media), breakpoints y z-index NO se convierten en tokens públicos (geometría estructural, no escala de theme).
- **Naming**: secciones NO son theme; no portar grid `mzp-*` ni naming `mzp-*`; Feature Card horizontal descartada; el "Hero" de Protocol no se llama Callout.
- **Convención de partials Phase D**: `web/templates/<x>.html` (`{{define "x"}}`) + `web/styles/<x>.css` (`@layer components`, tokens scoped en root, forced-colors) + `@import` en `app.css` + `web/styles_<x>_test.go` + `sourceAppCSS` + `npm run build` regenera `static/app.css`.

---

## 4. Desbloqueo de Phase G (screen recipes)

| Patrón Phase F | Desbloquea en Phase G |
|---|---|
| **Breadcrumb** | Settings y TODAS las páginas `/components/*` (GEO §9: breadcrumb visible + `BreadcrumbList` JSON-LD) |
| **Footer** | TODAS las recipes (chrome de sitio: nav secundaria + legal + locale) y contrato SEO §3 |
| **Hero/Billboard** | Landing/Public-Social Feed (recipe G-3); formaliza el home actual |
| **Card slots públicos** | Admin Resource, Public Feed (cards editoriales con media/tag/CTA) |
| **Video / Split** | Public Feed / Landing (media promocional) |

**No bloquean G**: Newsletter y Language Switcher. **Section Heading** desbloquea el encabezado de secciones de cualquier recipe (utilidad barata).
