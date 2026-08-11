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
| 2 | **Billboard/Hero** | ✖ pendiente | 2 | Composición: formaliza `.landing-hero` huérfano + `.hero-action`; heading `--ui-type-display-lg` + desc + CTA Link + media opcional |
| 3 | **Breadcrumb** | ✅ | 1 | `breadcrumb.html`/`breadcrumb.css` + `styles_breadcrumb_test.go`; markup canónico P1; GEO §9/§14 |
| 4 | **Callout** | ✅ | — | `callout.html`/`callout.css` (Phase D); naming resuelto (tip box) |
| 5 | **Card (slots públicos)** | ◐ pendiente | 2 | Extender `card.html` con media (`aspect-ratio`), tag (`--ui-badge`), meta, CTA (`--ui-card-action`) |
| 6 | **CTA Link** | ✅ (≈ Button link) | 1 | `button.html:6` variante `Href` → `<a class="ui-button">`; reusado por Empty state/Banner/Callout |
| 7 | **Feature Card** | ✖ pendiente | 2 | Composición Card + CTA Link (no primitiva); sin variante horizontal (deprecada upstream) |
| 8 | **Footer** | ✅ | 3 | `footer.html`/`footer.css` + slot en `layout.html` + `pageView.Footer`/`defaultFooter()` en `server.go`; `<details>/<summary>`; bloqueante Phase G |
| 9 | **Language Switcher** | ✖ pendiente | 3 | GET `?lang=` → 303 a URL localizada; submit visible; requiere modelo de locales (`lang`/RTL) |
| 10 | **Newsletter** | ✖ pendiente | 3 | POST + 422 `X-Gelium-Validation` + success view; reusa text-field/select/checkbox/button |
| 11 | **Notification Bar** | ✅ (≈ Banner) | 4 | Alias documental del Banner Gelium (`banner.html`, Phase D); variantes sticky/scripted diferidas |
| 12 | **Section Heading** | ✅ | 1 | Utilidad tipográfica: `section-heading.html`/`section-heading.css`; siempre `h2`, nunca `h1` |
| 13 | **Split** | ✖ pendiente | 2 | Composición grid 2 col (`.ui-split` + `.ui-split-body` + `.ui-split-media`); stack en narrow; bidi RTL |
| 14 | **Video** | ✅ | 1 | Contenedor responsive: `video.html`/`video.css`; `aspect-ratio` literal 16:9 (no se tokeniza); `<video controls>` nativo |

**Resumen**: 4 existentes ✅ (Article, Callout, CTA Link, Notification Bar) + **4 nuevos ✅ en esta entrega** (Breadcrumb, Section Heading, Video, Footer) + 6 pendientes ✖/◐ (Billboard/Hero, Card slots, Feature Card, Language Switcher, Newsletter, Split).

---

## 2. Orden de implementación

### Tier 1 — 100% estáticos, cero server contract (paralelizables)

1. **Breadcrumb** ✅ — partial + CSS + tests; markup ya canónico (P1). Desbloquea GEO §9/§14.
2. **Section Heading** ✅ — utilidad CSS + partial mínimo.
3. **Video** ✅ — contenedor + partial mínimo; cero JS.
4. **Article** — formalizar contrato `.prose` (tipografía + intro opcional); doc + posible utilidad CSS.
5. **CTA Link** ✅ — cerrado por Button link; variante inline con icono opcional.

### Tier 2 — composiciones de existentes (después de Card slots)

6. **Card → slots públicos** (media + aspect-ratio, tag `--ui-badge`, meta, CTA `--ui-card-action`).
7. **Feature Card** = Card + CTA Link (composición, CSS mínimo).
8. **Split** = grid 2 col (body + media), bidi RTL; partial opcional.
9. **Hero/Billboard** = composición que formaliza `.landing-hero` + `.hero-action`; útil para Landing (Phase G).

### Tier 3 — server contract (necesitan handler Go + datos)

10. **Footer** ✅ — partial + slot en layout + modelo de datos (`footerView`); `<details>/<summary>`. **Bloqueante Phase G.**
11. **Newsletter** — POST + 422 `X-Gelium-Validation` + success view (reusa text-field/select/checkbox/button).
12. **Language Switcher** — GET form + 303; requiere modelo de locales (`lang`/RTL); se compone dentro del Footer.

### Tier 4 — alias documental (sin código)

13. **Notification Bar** ✅ → alias de Banner en `gelium-ui-vocabulary.md`.

---

## 3. Reglas

- **Zero-JS end-to-end**: Footer → `<details>/<summary>`; Language Switcher → GET form con submit visible; Newsletter → POST + 422. HTMX solo como enhancement.
- **No-tokenizar**: `aspect-ratio` (Video), breakpoints y z-index NO se convierten en tokens públicos (geometría estructural, no escala de theme).
- **Naming**: secciones NO son theme; no portar grid `mzp-*` ni naming `mzp-*`; Feature Card horizontal descartada.
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
