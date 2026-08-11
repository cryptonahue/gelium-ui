# Gelium UI — Component Registry

> Inventario canónico de TODOS los componentes del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Regla de gate: los registries se crean SOLO porque los contratos de A–G están estables ("no agregar un registry hasta que los patterns tengan contratos estables").
> Fuentes de autoridad: el código real (`web/templates/*.html`, `web/styles/*.css`, `internal/app/*.go`), `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-core.md`, `docs/handoffs/{composition,vocabulary,state-patterns,public-patterns,screen-recipes}-audit.md`.

---

## 1. Cómo leer este registry

Cada entrada de la tabla maestra ($2) declara:

- **Nombre canónico** y **clase raíz `.ui-*`** — la API pública. Las clases son la única API estable; no hay prefijos de terceros (se renombró `m3-*` → `ui-*`, `theme-contract.md:91`).
- **Template / CSS** — el partial `{{define "…"}}` y su hoja. La convención de partials es `web/templates/<x>.html` + `web/styles/<x>.css` + `@import` en `web/styles/app.css` + `web/styles_<x>_test.go` (`public-content-patterns.md:81`).
- **Tokens que consume** — los `var(--ui-*)` leídos por el CSS del componente (inventario extraído de los archivos fuente; el contrato de cobertura vive en `theme-contract.md` §3).
- **Estados** — la superficie de estados del componente (rest/hover/focus/pressed/disabled/selected/error/loading/empty), modelada en el CSS del componente, nunca en el theme.
- **Contrato server** — si aplica (422, `loom:toast`, GET params, POST+303, none). Componentes sin contrato propio son partials inertes o con contrato nativo de form.
- **Variantes** — clases/variantes reales, no inventadas.
- **Categoría** — foundation/action/input/feedback/navigation/data/public/state-pattern/recipe-primitive. Los **state patterns** (Phase D) y los **public patterns** (Phase F) son categorías de *patterns*, no de componentes: se documentan como tales (ver `gelium-ui-pattern-registry.md`) y aparecen aquí porque cada uno tiene al menos un partial real.

**Reglas de mantenimiento**:

1. Un componente nuevo entra al registry solo tras pasar el rationale de `gelium-ui-composition-rules.md` §11.
2. Toda fila nueva requiere: partial + CSS + `@import` en `app.css` + test de contrato (`styles_<x>_test.go`) + ruta en `internal/app/routes.go` (si tiene docs page) — mismo checklist que el DoD de Phase D (`state-patterns-audit.md:102`).
3. Los tokens se leen del código, no de memoria: si el CSS cambia, esta tabla cambia.
4. **El JSON de este registry NO se sirve desde el server todavía** — ver `gelium-ui-dependency-metadata.md` §6 (tooling pendiente documentado).

---

## 2. Tabla maestra de componentes

> Categorías: **F** = foundation · **A** = action · **I** = input · **D** = data · **B** = feedback · **N** = navigation · **P** = public · **SP** = state-pattern · **RP** = recipe-primitive.
> "Route" = docs page en `/components/<slug>` (driver: `internal/app/routes.go:16-47`). "Handler" = view model Go en `internal/app/<x>.go`.

| Componente | Clase raíz | Template | CSS | Route | Handler | Tokens clave | Estados | Variantes | Categoría |
|---|---|---|---|---|---|---|---|---|---|
| Elevation | `.ui-elevation-{0..5}` | `elevation.html` | `elevation.css` | `/components/elevation` | `elevation.go` | `--ui-shadow-0..5`, `--ui-radius-md` | rest | 6 niveles | F |
| Focus ring | `.ui-focus-ring` | `focus-ring.html` | `focus-ring.css` | `/components/focus-ring` | `focus_ring.go` | `--ui-focus-*`, `--ui-color-focus-ring` | focus-visible | none | F |
| Icon | (SVG inline) | `icon.html` | `icon.css` | `/components/icon` | `icon.go` | `--ui-color-fg-muted`, `--ui-color-primary`, `--ui-size-icon` | none | decorative `aria-hidden` | F |
| Divider | `.ui-divider` | `divider.html` | `divider.css` | `/components/divider` | `divider.go` | `--ui-divider-{color,thickness}`, `--ui-color-fg-muted` | none | none | F |
| Tokens core | — | — | `tokens.css` | — | — | define `--ui-color-*`, `--ui-size-*`, `--ui-state-*`, `--ui-focus-*`, `--ui-border-*`, `--ui-radius-*`, `--ui-shadow-*`, `--ui-motion-*`, `--ui-space-*` | — | — | F (capa, no componente) |
| Button | `.ui-button` | `button.html` | `button.css` | `/components/button` | `button.go` | `--ui-color-primary[-fg]`, `--ui-color-secondary[-fg]`, `--ui-color-border-strong`, `--ui-state-*-opacity`, `--ui-size-control`, `--ui-radius-full`, `--ui-type-label-lg` | rest/hover/focus/pressed/disabled/loading (`aria-busy`) | `ui-button-{primary,secondary,outline,text}`; link (`Href`) vs button; spinner | A |
| Icon button | `.ui-icon-button` | `icon-button.html` | `icon-button.css` | `/components/icon-button` | `icon_button.go` | `--ui-color-primary`, `--ui-size-control`, `--ui-size-icon`, `--ui-state-*` | rest/hover/focus/pressed/disabled/toggle (`aria-pressed`) | `ui-icon-button-{standard,outlined,filled,filled-tonal}` | A |
| FAB | `.ui-fab` | `fab.html` | `fab.css` | `/components/fab` | `fab.go` | `--ui-fab-{primary,surface,secondary}-*`, `--ui-shadow-{1,3,4}`, `--ui-size-*` | rest/hover/focus/pressed/disabled (`aria-disabled`) | `ui-fab-{primary,surface,secondary,lowered}`; `ui-fab-{small,medium,large}`; extended/icon-only | A |
| Chips | `.ui-chip` | `chips.html` | `chips.css` | `/components/chips` | `chips.go` | `--ui-chip-height`, `--ui-color-primary`, `--ui-color-secondary[-fg]`, `--ui-state-*` | rest/hover/focus/pressed/selected/disabled/removable | `ui-chip-{filter,input,suggestion}`; remove button | A |
| Segmented buttons | `.ui-segmented-button-set` / `.ui-segmented-button` | `segmented-button.html` | `segmented-button.css` | `/components/segmented-button` | `segmented_button.go` | `--ui-segmented-button-{selected-container,selected-fg,outline}`, `--ui-color-primary` | rest/hover/focus/pressed/selected/disabled | radio (single) / checkbox (multi) / button; `--action` | A |
| Menu | `.ui-menu` | `menu.html` | `menu.css` | `/components/menu` | `menu.go` | `--ui-menu-*` (scoped), `--ui-shadow-2`, `--ui-color-primary`, `--ui-size-*` | rest/hover/focus/pressed/selected (`aria-selected`)/disabled | Popover API (zero-JS); items nativos; dividers | A |
| Text field | `.ui-text-field` | `text-field.html` | `text-field.css` | `/components/text-field` | `text_field.go` | `--ui-field-*`, `--ui-size-field`, `--ui-color-primary`, `--ui-radius-xs`, `--ui-border-*` | rest/hover/focus/disabled/error (`aria-invalid`)/helper | `ui-text-field-{filled,outlined}`; input/textarea | I |
| Checkbox | `.ui-checkbox` | `checkbox.html` | `checkbox.css` | `/components/checkbox` | `checkbox.go` | `--ui-checkbox-*`, `--ui-state-*`, `--ui-color-focus-ring` | rest/hover/focus/pressed/checked/indeterminate/disabled/error | none (nativo `input[type=checkbox]`) | I |
| Radio | `.ui-radio` | `radio.html` | `radio.css` | `/components/radio` | `radio.go` | `--ui-radio-*`, `--ui-color-danger` | rest/hover/focus/pressed/checked/disabled/error | none (nativo) | I |
| Switch | `.ui-switch` | `switch.html` | `switch.css` | `/components/switch` | `switch.go` | `--ui-switch-*` (17 tokens), `--ui-state-*` | rest/hover/focus/pressed/checked/disabled | none (nativo) | I |
| Select | `.ui-select` | `select.html` | `select.css` | `/components/select` | `select.go` | `--ui-select-*` (scoped), `--ui-field-*`, `--ui-radius-*` | rest/hover/focus/disabled/error (`aria-invalid`)/selected | `ui-select-{filled,outlined}` | I |
| Slider | `.ui-slider` | `slider.html` | `slider.css` | `/components/slider` | `slider.go` | `--ui-slider-*` (10), `--ui-color-focus-ring` | rest/hover/focus/pressed/disabled; fill por instancia (`--ui-slider-fill`) | none (nativo `input[type=range]`) | I |
| List | `.ui-list` | `list.html` | `list.css` | `/components/list` | `list.go` | `--ui-list-*` (scoped), `--ui-size-item-*`, `--ui-state-*`, `--ui-color-primary` | rest/hover/focus/pressed/selected (`:checked`)/disabled | `ui-list-item--{two-line,three-line}`; navegación / selección / estático | D |
| Card | `.ui-card` | `card.html` | `card.css` | `/components/card` | `card.go` | `--ui-card-*`, `--ui-shadow-1`, `--ui-radius-*` | rest/hover/focus/pressed (state layer) | `ui-card-{elevated,filled,outlined}`; slots `ui-card-title/body` | D |
| Data table | `.ui-data-table` | `data-table.html` | `data-table.css` | `/components/data-table` | `data_table.go` | `--ui-data-table-*` (scoped), `--ui-color-primary`, `--ui-state-*` | rest/hover/focus/pressed/selected (`:has(input:checked)`)/sort (`aria-sort`)/pagination (`aria-current`)/empty/error | sortable columns; row selection; pagination footer | D |
| Avatar | `.ui-avatar` | `avatar.html` | `avatar.css` | — | `avatar.go` | `--ui-avatar-*`, `--ui-color-surface-container`, `--ui-radius-full` | none (decorativo, `aria-hidden`) | `ui-avatar--{sm,md,lg}`; image/initials | RP |
| Pagination | `.ui-pagination` | `pagination.html` | `pagination.css` | — | `pagination.go` | `--ui-pagination-{page-color,active-color}`, `--ui-radius-full` | current (`aria-current`)/disabled boundary | standalone (extracción del footer del Data table) | RP |
| Dialog | `.ui-dialog` / `.ui-dialog-page` | `dialog.html` | `dialog.css` | `/components/dialog` | `dialog.go` | `--ui-dialog-*`, `--ui-color-*`, `--ui-motion-long` | open/closed (top layer); light dismiss; confirm/cancel | confirm; page variant (`<dialog open>` server-rendered) | B |
| Toast | `.ui-toast` | `toast.html` | `toast.css` | `/components/toast` | `toast.go` | `--ui-toast-*`, `--ui-shadow-3`, `--ui-color-*` | transitorio: auto-dismiss 4s/8s, pausable; dismiss manual | `ui-toast-{info,success,warning,error}`; action | B |
| Progress | `.ui-progress` | `progress.html` | `progress.css` | `/components/progress` | `progress.go` | `--ui-progress-{track,indicator,radius,track-height}` | determinate / indeterminate (`aria-busy`) | determinate/indeterminate | B |
| Badge | `.ui-badge` | `badge.html` | `badge.css` | `/components/badge` | `badge.go` | `--ui-badge-*`, `--ui-color-danger|success|warning|info[-fg]`, `--ui-radius-full` | none (decorativo, nunca color-only) | `ui-badge-large`; tones `ui-badge--{error,success,warning,info}`; dot/count | B |
| Tooltip | `.ui-tooltip` | `tooltip.html` | `tooltip.css` | `/components/tooltip` | `tooltip.go` | `--ui-tooltip-*` (scoped), `--ui-shadow-2`, `--ui-color-*` | visible `:hover`/`:focus-within` | plain/rich; `--top` placement; host | B |
| Tabs | `.ui-tabs` | `tabs.html` | `tabs.css` | `/components/tabs` | `tabs.go` | `--ui-tabs-*`, `--ui-color-primary`, `--ui-size-*` | active (`aria-current="page"`), rest/hover/focus/pressed | `<nav>` links; stacked | N |
| Navigation bar | `.ui-nav-bar` | `navigation-bar.html` | `navigation-bar.css` | `/components/navigation-bar` | `navigation_bar.go` | `--ui-nav-bar-*` (scoped), `--ui-size-nav-*` | active (`aria-current`)/hover/focus/pressed; badge count | destinations + indicator | N |
| Navigation tab | `.ui-nav-tab` | `navigation-tab.html` | `navigation-tab.css` | `/components/navigation-tab` | `navigation_tab.go` | `--ui-nav-tab-*` (scoped), `--ui-size-nav-*` | active/hover/focus/pressed; badge | primary/secondary tabs | N |
| Navigation drawer | `.ui-navigation-drawer` | `navigation-drawer.html` | `navigation-drawer.css` | `/components/navigation-drawer` | `navigation_drawer.go` | `--ui-navigation-drawer-*` (scoped), `--ui-dialog-*` | active (`aria-current`)/hover/focus/pressed; open/closed | modal (`<dialog>`) / permanente (`<nav>`) | N |
| Breadcrumb | `.ui-breadcrumb` | `breadcrumb.html` | `breadcrumb.css` | — | (server.go crumbs) | `--ui-breadcrumb-*`, `--ui-color-fg[-muted]` | current (`aria-current="page"`) | separador por CSS (no texto literal) | P |
| Footer | `.ui-footer` | `footer.html` | `footer.css` | — | (server.go `footerView`) | `--ui-footer-*`, `--ui-color-*`, `--ui-space-*` | none | `<details>/<summary>` plegable; grid→stack | P |
| Section heading | `.ui-section-heading` | `section-heading.html` | `section-heading.css` | — | — | `--ui-section-heading-*`, `--ui-type-headline-sm` | none (siempre `h2`) | `--centered`; eyebrow opcional | P |
| Video | `.ui-video` | `video.html` | `video.css` | — | — | `--ui-video-*`, `--ui-radius-sm` | none | `--aspect-4-3`; `<video controls>` nativo | P |
| Feature card | `.ui-feature-card` | `feature-card.html` | `feature-card.css` | — | — | `--ui-color-surface-container`, `--ui-space-*` | none (reusa `.ui-card`) | media 16:9 (aspect literal, no token) | P |
| Hero | `.ui-hero` | `hero.html` | `hero.css` | — | — | `--ui-hero-*`, `--ui-type-display-lg`, `--ui-color-scrim` | none | `--has-media`; scrim | P |
| Language switcher | `.ui-language-switcher` | `language-switcher.html` | `language-switcher.css` | — | — | `--ui-language-switcher-*`, `--ui-color-border` | none | GET form + submit visible (cero auto-submit) | P |
| Newsletter | `.ui-newsletter` | `newsletter.html` | `newsletter.css` | — | `newsletter.go` (handler) | `--ui-newsletter-*`, `--ui-color-surface-container`, `--ui-color-error` | success (POST → `role="status"`)/error (422 inline) | none | P |
| Split | `.ui-split` | `split.html` | `split.css` | — | — | `--ui-split-*`, `--ui-radius-sm` | none | media + body; stack narrow; bidi RTL | P |
| Banner | `.ui-banner` | `banner.html` | `banner.css` | — | — | `--ui-banner-*`, `--ui-color-danger|success|info`, `--ui-color-warning-container` | persistente nivel página; dismiss (`POST+303`) | `ui-banner--{error,success,warning,info}` | SP |
| Callout | `.ui-callout` | `callout.html` | `callout.css` | — | — | `--ui-callout-*`, `--ui-color-info`, `--ui-color-secondary[-fg]` | none (nota ignorable) | `ui-callout--{info,tip}` | SP |
| Empty state | `.ui-empty-state` | `empty-state.html` | `empty-state.css` | — | (server-rendered en data-table/recipes) | `--ui-empty-state-*`, `--ui-color-fg[-muted]` | vacío (`role="status"`) | `--compact`; title+body+CTA | SP |
| Error state | `.ui-error-state` | `error-state.html` | `error-state.css` | — | (server 404/500) | `--ui-error-state-*`, `--ui-color-danger`, `--ui-type-display-lg` | error (`role="alert"` + status real) | code + retry `.ui-button` | SP |
| Inline alert | `.ui-inline-alert` | `inline-alert.html` | `inline-alert.css` | — | — | `--ui-inline-alert-*`, `--ui-color-danger-container`, `--ui-color-*` | persistente contextual (`role="alert"`/`status`) | `ui-inline-alert--{error,success,warning,info}` | SP |
| Skeleton | `.ui-skeleton` | `skeleton.html` | `skeleton.css` | — | — | `--ui-skeleton-*`, `--ui-color-surface-container` | loading (`role="status"` + sr-only) | `--avatar`; blocks `--title/--line/--short/--circle` | SP |
| Validation summary | `.ui-validation-summary` | `validation-summary.html` | `validation-summary.css` | — | — | `--ui-validation-summary-*`, `--ui-color-danger-container` | error (422) | links a `#campo-error` | SP |
| Success feedback | (reuso) | reusa `inline-alert--success` / `banner--success` | — | — | — | — | success persistente (`role="status"`) | — | SP |

> Nota de categoría: **Toast** es transitorio-de-acción y está listado en feedback; los 8 **state patterns** de Phase D son Banner, Callout, Empty state, Error state, Inline alert, Skeleton, Validation summary y Success feedback (reuso) — ver `state-patterns-audit.md:45` y `gelium-ui-pattern-registry.md` §2.
> **Avatar, Pagination y Badge tone** son primitivas de recipe (introducidas en Phase G junto a Ops Queue y Public Feed) — cross-listed: Avatar/Pagination en RP, Badge tone en B.

---

## 3. Componentes con contrato server (resumen)

El contrato completo de cada recipe vive en `docs/gelium-ui-screen-recipes.md`. Por componente:

| Componente | Contrato server | Dónde |
|---|---|---|
| Text field | 422 + `X-Loom-Validation: true`; valor preservado (`text_field.go:62`), focus al campo fallido (`text_field.go:67`) | `internal/app/text_field.go` |
| Select | 422 + `X-Loom-Validation` (campo error); `:has(select:disabled)` | `select.css:88` |
| Data table | `GET ?q=&sort=&dir=&page=&selection=` (vocabularios cerrados); `HX-Request` bifurca fragmento vs página; refresh `POST` + `HX-Trigger loom:toast` | `data_table.go`, `data-table.html` |
| Toast | `HX-Trigger: {"loom:toast":{"type":"info\|success\|warning\|error","message":"…"}}`; fallback no-JS inline (`toast.go`); validación NUNCA toast (`toast.go:129-133`) | `toast.go:13-14,45` |
| Banner | dismiss = `POST + 303` | `banner.html` |
| Newsletter | `POST` + 422 `X-Loom-Validation`; success `role="status"` persistente | `internal/app/newsletter.go` |
| Dialog | `GET/POST /components/dialog/confirm`; page variant `<dialog open>` server-rendered | `dialog.go` |
| Pagination | `GET ?page=` con clamping al rango válido | `pagination.go` |
| Error state | status HTTP real (404/500/503) + re-render server | `server.go` (render 404) |
| Recipe primitives | Avatar/Pagination son partials sin contrato propio; el contrato lo posee la recipe que los consume | `recipe_*.go` |

Regla de oro (`composition-rules.md:189`): si un estado es navegable (listado, filtro, sort, paginación, selección), es una URL. Componentes con contrato GET params (Data table, Pagination, recipes) cumplen esto; el resto es markup inerte o contrato nativo de form.

---

## 4. Estados vs categorías (matriz)

La state matrix de `composition-rules.md` §8 cubre los **patterns** (Card, List, Data table, Queue, Feed, Dashboard, Dialog, Form) con sus GAP. Los **componentes** de esta tabla cubren su propia superficie de estados en CSS; un theme nuevo verifica la matriz theme × component × variant × state (`theme-contract.md` §9.4). Los 8 componentes con tokens scoped (List, Menu, Data table, Nav bar, Nav tab, Navigation drawer, Segmented button, Tooltip) son los que un theme NO puede sobreescribir salvo declaración global — ver `gelium-ui-theme-registry.md` §4.

---

## 5. Checklist de ingreso de un componente nuevo

1. Rationale aprobado (`composition-rules.md` §11).
2. Partial `web/templates/<x>.html` + `web/styles/<x>.css` + `@import` en `app.css` + `web/styles_<x>_test.go` + `sourceAppCSS` en sync (`styles_contract_test.go`).
3. Tokens `--ui-<x>-*` cerrados (scoped en el CSS o en el theme).
4. States + variantes + forced-colors/reduced-motion (el core centraliza ambos bloques).
5. Si requiere docs page: ruta en `routes.go` + `docsSections` en `docs.go` + contenido markdown en `web/content/<x>.md`.
6. Esta tabla actualizada.

---

**Definición de done (Phase J)**: tabla maestra generada desde el código real, categorías sin colisión con los registries de patterns, cada componente con contrato server resuelto, checklist de ingreso referenciado desde `composition-rules.md`.
