# Gelium UI — Banner Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón BANNER (patrón 5 de Phase D, después de Empty State `eba1c4c`, Skeleton `0688020`, Inline Alert `43c0dac`, Validation Summary `c9afd49`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :411, :433), `docs/gelium-ui-vocabulary.md` (:122-129 Banner), `docs/gelium-ui-composition-rules.md` (§4.8, §8 state matrix, §9 server-driven :171), `docs/handoffs/{state-patterns-audit,inline-alert-audit,validation-summary-audit,mozilla-protocol-audit}.md`, `internal/app/{demo_whatsapp,server,routes,toast}.go`, `web/templates/{demo-whatsapp,layout,inline-alert,empty-state,validation-summary}.html`, `web/styles/{demo-whatsapp,inline-alert,empty-state,validation-summary,tokens,app}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/styles_inline_alert_test.go`, `web/styles_demo_whatsapp_test.go`.

---

## 1. ESTADO ACTUAL

**No existe ningún banner.** Cero coincidencias `ui-banner`/`.banner` en templates, styles ni en el compilado `web/static/app.css`. El vocabulario lo marca ✖ (`vocabulary.md:122-129`); la matriz del roadmap lo marca bloqueante "Crear componente" (`roadmap.md:411`).

### 1.1 Ad-hoc `demo-wa-expired` (única referencia visual)

- `web/templates/demo-whatsapp.html:132-144` — `{{if $chat.WindowExpired}}` → `<div class="demo-wa-expired" role="note">` con título (`.demo-wa-expired-title`), cuerpo (`.demo-wa-expired-body`) y un form interno `POST /demo/whatsapp/send-template` (selector de mensaje prediseñado + botón `.ui-button`). **`role="note"`**: semánticamente flojo para un aviso que exige acción (debería ser `alert`/`status`); es el gap que el patrón cierra.
- CSS: `web/styles/demo-whatsapp.css:395-404` — `.demo-wa-expired` (grid, gap .5rem, padding 1rem, `background: var(--ui-color-surface)`, `border-block-start`), título con `color: var(--ui-color-error)` (alias de danger, tokens.css:38), cuerpo `font: var(--ui-type-body-sm)` + `color: var(--ui-color-fg-muted)`. Sin surface de tone, sin icono, sin dismiss.
- Pinneado por tests: `web/styles_demo_whatsapp_test.go:25,63` (`.demo-wa-expired` presente) e `internal/app/demo_whatsapp_test.go:63-82` (`TestWhatsAppExpiredConversationBlocksComposer`, contrato `demo-wa-expired`). Cualquier migración del ad-hoc obliga a tocar estos tests.
- Nota: no es page-level — vive dentro del chat pane por conversación. Es referencia **visual** (candidato de tono/título/cuerpo/CTA), no un caso exacto de banner global.

### 1.2 Aviso global en layout

**No existe.** `web/templates/layout.html:12-15` — el `<header class="site-header">` solo tiene brand + `<nav aria-label="Primary">`. No hay notice, alert ni slot de banner en el layout; el `<main class="docs-shell">` arranca directo tras el header (:16). No hay `<footer>` (gap Phase E, `mozilla-protocol-audit.md:15`). `pageView` (`internal/app/server.go:32-69`) no tiene campo Banner.

### 1.3 Precedentes de convención (patrones 1-4 entregados)

- **Empty state** (`eba1c4c`): partial + CSS + tests; único con consumidor real (fila `<td colspan>` en `data-table.html:69`).
- **Skeleton** (`0688020`): partial + CSS + tests; sin consumidor.
- **Inline alert** (`43c0dac`): partial + CSS + tests; sin consumidor (primitiva lista).
- **Validation summary** (`c9afd49`): partial + CSS + tests; sin consumidor (primitiva lista, "ningún form del repo justifica un summary hoy").
- **Convención exacta de archivos por commit** (verificado `git show --stat`): `web/templates/<x>.html` + `web/styles/<x>.css` + `web/styles_<x>_test.go` + `@import` en `web/styles/app.css` + `web/static/app.css` regenerado + `web/styles_contract_test.go` (`sourceAppCSS` :24-58 + `TestComponentSizeTokensDeclaredScoped` :405-429). **Sin handler Go, sin ruta, sin docs page** — el partial se entrega como primitiva.

---

## 2. NAMING DECISION

**Canónico: Banner.** "Notification Bar" queda como **alias documentado**, no como patrón separado — exactamente el precedente de "Snackbar → alias de Toast" (`vocabulary.md:142`).

Evidencia de la colisión resuelta:
- `roadmap.md:58` — "Notification Bar (≈ Banner) requieren resolución contra Protocol"; `roadmap.md:205` — "Notification Bar (≈ Banner Gelium)" en Phase F; `roadmap.md:433` — "Feature Card / Notification Bar | pattern | F | … ≈ Banner" (composición Card+CTA, scope distinto, no conflict).
- `mozilla-protocol-audit.md:59,86,113,197` — `<aside class="mzp-c-notification-bar">` + `<p>` + variantes tone + dismiss + cta ≈ **Banner Gelium** (`vocabulary.md:122-129`, persistente página/sitio); base estática, dismiss = POST+303, sticky/scripted diferido (JS).
- `vocabulary.md:122` — el término canónico del vocabulario Gelium ya es **Banner** (aliases: site banner, notice).

Naming del contrato:
- Partial: `{{define "banner"}}` → `web/templates/banner.html`.
- Clase raíz: `.ui-banner` + modificador de tone `.ui-banner--{error|warning|info|success}` (propuesta original `state-patterns-audit.md:38`).
- En docs, "Notification Bar" se menciona como alias (Phase F lo usará para la composición pública sobre el mismo patrón).

---

## 3. CONTRATO PROPUESTO

### 3.1 Markup (patrón inline-alert.html / empty-state.html — partial server-rendered)

```html
{{define "banner"}}<div class="ui-banner ui-banner--{{.Tone}}" role="{{if eq .Tone "error"}}alert{{else}}status{{end}}">
  {{if .Icon}}<span class="ui-banner-icon" aria-hidden="true">{{.Icon}}</span>{{end}}
  <div class="ui-banner-content">
    {{if .Title}}<p class="ui-banner-title">{{.Title}}</p>{{end}}
    <p class="ui-banner-body">{{.Body}}</p>
  </div>
  {{if .CTA}}<a class="ui-button" href="{{.CTAHref}}">{{.CTALabel}}</a>{{end}}
  {{if .DismissHref}}<form class="ui-banner-dismiss" method="post" action="{{.DismissHref}}">
    <button class="ui-icon-button" type="submit" aria-label="Dismiss">{{.DismissIcon}}</button>
  </form>{{end}}
</div>{{end}}
```

- Root `<div>` full-width nivel página (flex row, `align-items: center`), distinto del flex column de inline-alert: el banner es una barra horizontal al tope, no un bloque de sección.
- Icono opcional `aria-hidden="true"` (glyph trusted, precedente inline-alert.html:2 / text-field.html:6). Título opcional como `<p>` con `font: var(--ui-type-title-md)` (precedente `ui-empty-state-title` — evita saltos de heading; el `<h1>` de la página queda reservado). Cuerpo con `font: var(--ui-type-body-sm)`.
- CTA (acción primaria, p.ej. "Reautenticarse", "Ver estado") = **`<a class="ui-button">` real** con href (navegación, precedente empty-state.html:5). Nunca div-span como control.
- Sin auto-dismiss (contrato `vocabulary.md:126`). El banner persiste hasta que el servidor limpia el estado.

### 3.2 Dismiss: `<button>` en `<form method="post">` (recomendado) vs `<a>`

| Opción | Mecánica | Veredicto |
|---|---|---|
| **`<button type="submit">` dentro de `<form method="post" action="{ruta-dismiss}">`** | POST a ruta de dismiss → handler limpia el flag → `http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)` (303). Re-render sin banner. 0 JS. | **RECOMENDADA** — es literalmente el patrón del sistema (POST+303, `composition-rules.md:171`, `demo_whatsapp.go:559,573`) y el único que puede **limpiar estado server-side**. Un dismiss sin efecto en servidor dejaría el banner vivo en el siguiente render (anti-regla "sin URL = sin no-JS", `composition-rules.md:173`). |
| `<a href="{ruta-dismiss}">` (GET) | Navegación GET con efecto colateral. | **NO recomendado** — GET con mutación rompe idempotencia/semántica HTTP; si el estado vive en servidor, un GET no debería borrarlo y un back/refresh lo reviviría. Viable solo para dismiss puramente cosmético cliente, que el sistema no usa. |

- El botón reusa visual `.ui-icon-button` (existe, `icon-button.html:9`) pero con `type="submit"` real (la primitiva icon-button renderiza `type="button"`; el dismiss es el caso que exige submit). `aria-label="Dismiss"` en el botón (icono decorativo `aria-hidden`/glyph).
- El form va como sibling del content dentro del root flex (no anida el CTA ni el texto).

### 3.3 Tones

**Los 4, todos con consumidor** (misma decisión que inline-alert, que entrega 4 tones y los pinnea en test):

- `error` → sesión expirada, aviso crítico que exige acción (consumidor inmediato: migración de `demo-wa-expired`).
- `warning` → mantenimiento programado.
- `info` → consentimiento/aviso global (cookie consent, `vocabulary.md:127`).
- `success` → success persistente a nivel página (patrón 7 lo reusa — `state-patterns-audit.md:42`; NUNCA toast, anti-regla `composition-rules.md:126`).

### 3.4 ARIA

- `role="alert"` (assertive) **solo** en `--error`; `role="status"` (polite) en warning/info/success — derivado del tone en el template, idéntico a `inline-alert.html:1`.
- **Sin `aria-live` adicional**: `role="alert"`/`status` ya portan el anuncio y el banner es server-rendered (presente al cargar). Añadir `aria-live` al root sería redundante. Si un swap HTMX futuro reemplaza el banner, el role ya anuncia el cambio.
- Sin auto-dismiss → sin `aria-live` de región; el dismiss es una acción explícita del usuario.
- Nada color-only: icono + texto + role + bloque forced-colors (convención del sistema, `ux-accessibility-audit.md`).

### 3.5 Tokens

Core (`web/styles/tokens.css`) — **verificados**:

- `--ui-color-danger` (:32), `--ui-color-danger-container` (:51), `--ui-color-warning` (:47), `--ui-color-warning-fg` (:48), `--ui-color-warning-container` (:49), `--ui-color-info` (:50), `--ui-color-success` (:45), `--ui-color-success-fg` (:46). `--ui-color-error` = alias de danger (:38).
- **NO existen** `success-container` ni `info-container` (ni en core ni en `themes/theme-material/theme.css`) → success/info usan `color-mix(in srgb, var(--ui-color-*) 12%, var(--ui-color-surface))`, precedente `.ui-inline-alert--success/--info` (`inline-alert.css:31-32`) y `.demo-wa-quality--GREEN` (`demo-whatsapp.css:501`). Cero cambios de contrato en tokens.css/theme.
- Espacio: `--ui-space-2/3/4/6` (escala core contract `styles_contract_test.go:263-270`). Radius: `--ui-radius-sm`/`--ui-radius-md`. Type: `--ui-type-title-md` (título), `--ui-type-body-sm` (cuerpo), `--ui-type-label-lg` (opcional CTA). Size: `--ui-size-icon`/`--ui-size-icon-sm` (icono). Todos en core.

Tokens scoped `--ui-banner-*` (declarados en el root, patrón `inline-alert.css:12-18`):

```css
.ui-banner {
  --ui-banner-padding: var(--ui-space-3) var(--ui-space-4);
  --ui-banner-gap: var(--ui-space-3);
  --ui-banner-radius: var(--ui-radius-sm);
  --ui-banner-bg: var(--ui-color-surface-container);
  --ui-banner-fg: var(--ui-color-fg);
  --ui-banner-title-color: var(--ui-color-fg);
  --ui-banner-body-color: var(--ui-color-fg-muted);
  --ui-banner-icon-size: var(--ui-size-icon);
}
.ui-banner--error   { --ui-banner-bg: var(--ui-color-danger-container);  --ui-banner-fg: var(--ui-color-danger); }
.ui-banner--warning { --ui-banner-bg: var(--ui-color-warning-container); --ui-banner-fg: var(--ui-color-warning-fg); }
.ui-banner--success { --ui-banner-bg: color-mix(in srgb, var(--ui-color-success) 12%, var(--ui-color-surface)); --ui-banner-fg: var(--ui-color-success); }
.ui-banner--info    { --ui-banner-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface));    --ui-banner-fg: var(--ui-color-info); }
```

### 3.6 CSS

- `web/styles/banner.css` nuevo, `@layer components`, anatomía de `inline-alert.css`/`empty-state.css`: tokens scoped en el root, 4 modificadores de tone, `.ui-banner-icon` (width/height `--ui-banner-icon-size`, `flex: none`), `.ui-banner-content` (`flex: 1`), `.ui-banner-title` (`font: var(--ui-type-title-md)`), `.ui-banner-body` (`font: var(--ui-type-body-sm)`), `.ui-banner-dismiss` (margin-inline-start auto). Bloque `@media (forced-colors: active)` (borde `CanvasText`, `Mark` en error, precedente `inline-alert.css:54-69`). **Sin animación → sin bloque reduced-motion**.
- `web/styles/app.css`: `@import "./banner.css";` (tras `validation-summary.css` :37, junto a los patrones de estado).
- `web/styles_contract_test.go`: agregar `"styles/banner.css"` a `sourceAppCSS` (:24-58) y `"banner.css": "--ui-banner-icon-size:"` a `TestComponentSizeTokensDeclaredScoped` (:405-429). Guard `TestNoColorLiteralsInComponents` (checked ≥ 25, :743) sigue verde (26 → 27 archivos).
- Rebuild: `npm run build` regenera `web/static/app.css` (precedente commits 43c0dac/c9afd49).

---

## 4. INTEGRACION

### 4.1 Dónde vive en el layout

`web/templates/layout.html` — slot global **fuera del `<main>`**, entre el `</header>` (línea 15) y `<main>` (línea 16):

```html
<header class="site-header">…</header>
{{if .Banner}}{{template "banner" .Banner}}{{end}}
<main class="docs-shell">…
```

- A nivel sitio (mantenimiento, consent) el banner puede ir antes del header; la recomendación base es **después del header, antes de main** (slot persistente de página). Landmarks intactos: header → banner → main.
- `pageView` (`internal/app/server.go:32-69`) gana un campo `Banner *bannerView` (pointer nil = sin banner, como los demos actuales).

### 4.2 Cómo se dispara

- **Output del servidor, nunca estado cliente**: el handler que detecta la condición (sesión expirada, maintenance flag, consent pendiente) setea `pageView.Banner` y el layout lo renderiza en toda página. Si la condición es global (mantenimiento), un middleware/helper central lo inyecta; si es por-ruta (sesión), el handler propio lo setea.
- El partial es una primitiva: se entrega sin ruta de componente propia (mismo criterio que validation-summary `c9afd49` — "primitiva lista para Phase G").
- Vista de datos (Go, patrón de los view models existentes):

```go
type bannerView struct {
    Tone        string // info|success|warning|error (vocabulario cerrado, sanitize como toast.go:45-54)
    Icon        string // glyph trusted opcional
    Title       string
    Body        string
    CTAHref     string // CTA opcional (link real, .ui-button)
    CTALabel    string
    DismissHref string // ruta POST de dismiss (opcional; ausente = banner sin dismiss)
}
```

### 4.3 Dismiss (POST + 303, contrato existente (d), `state-patterns-audit.md:75`)

```text
POST {DismissHref} → handler: limpia el flag en el store → http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
```

- 303 SeeOther a la URL actual → re-render sin banner. Patrón exacto `demo_whatsapp.go:559,573`. 0 JS; `app.js` intacto. Sin contrato nuevo.

### 4.4 Migración del ad-hoc `demo-wa-expired`

- Opcional en Phase D (criterio validation-summary §4.1 opción 2: primitiva primero, wiring cuando haya consumidor). Si se migra: el bloque `demo-whatsapp.html:132-144` usaría el partial con tone error + CTA (mantiene el form send-template como CTA o contenido) y `role="note"` → `alert`/`status`. Toca `internal/app/demo_whatsapp_test.go:76` y `web/styles_demo_whatsapp_test.go:25,63` (pinnean `.demo-wa-expired`). Como es conversation-level (no page-level), la migración no es obligatoria: el demo queda como referencia visual legítima.

---

## 5. TESTS PROPUESTOS

1. **`web/styles_banner_test.go`** (nuevo, patrón `styles_inline_alert_test.go`):
   - `TestBannerPrimitiveCSSMapsTokens` — contratos exactos: `.ui-banner {`, `display: flex`, `align-items: center`, `gap: var(--ui-banner-gap)`, `padding`, `border-radius`, `background`, `color`, los 4 modificadores de tone, icono/título/cuerpo, todos los tokens scoped. Sin `transition:`/`animation:`/`prefers-reduced-motion`.
   - `TestBannerContractCSSWired` — compilado embebido (`Assets.ReadFile("static/app.css")`) contiene `.ui-banner`, los 4 modificadores, `@media (forced-colors:active)`.
   - `TestBannerClassVocabularyIsClosed` — cada clase en template ↔ selector en CSS (`ui-banner`, `ui-banner-icon`, `ui-banner-content`, `ui-banner-title`, `ui-banner-body`, `ui-banner-dismiss`); interpolación `ui-banner--`; prohibido `ui-banner-demo`.
   - `TestBannerTonesUseCoreTokens` — cada tone resuelve solo tokens core (`--ui-color-danger*`, `--ui-color-warning*`, `--ui-color-success`, `--ui-color-info`) o color-mix sobre ellos; nunca hex literal (guard del blanket `TestNoColorLiteralsInComponents`).
2. **`web/styles_contract_test.go`** (modificar): `sourceAppCSS` + `TestComponentSizeTokensDeclaredScoped` (+`--ui-banner-icon-size:`).
3. **(Opcional, render Go)** `internal/app/banner_test.go` — `role="alert"` vs `role="status"` por tone, icono `aria-hidden`, CTA como `<a class="ui-button" href>`, dismiss como `<form method="post">` con `<button type="submit">`, sin dismiss cuando `DismissHref` vacío, sin banner cuando `.Banner == nil`.
4. **(Opcional, si se migra el demo)** ajustes en `internal/app/demo_whatsapp_test.go` y `web/styles_demo_whatsapp_test.go` por remoción/adaptación de `.demo-wa-expired`.

---

## 6. FILES IMPACTADOS (solo read-only)

**Nuevos**:
- `web/templates/banner.html` — partial `{{define "banner"}}`.
- `web/styles/banner.css` — primitiva + tones + tokens scoped + forced-colors.
- `web/styles_banner_test.go` — tests de contrato (patrón inline-alert/validation-summary).
- `docs/handoffs/banner-audit.md` — este handoff.

**Modificados**:
- `web/styles/app.css` — `@import "./banner.css";` (+1, tras `validation-summary.css` :37).
- `web/static/app.css` — regenerado por `npm run build` (precedente commits 43c0dac/c9afd49).
- `web/styles_contract_test.go` — `sourceAppCSS` (:24-58) + `TestComponentSizeTokensDeclaredScoped` (:405-429).
- `web/templates/layout.html` — slot `{{if .Banner}}{{template "banner" .Banner}}{{end}}` entre header (:15) y main (:16).
- `internal/app/server.go` — campo `Banner *bannerView` en `pageView` (:32-69).

**Opcionales (migración del ad-hoc / wiring Phase G)**:
- `web/templates/demo-whatsapp.html` — `demo-wa-expired` (:132-144) → partial banner.
- `web/styles/demo-whatsapp.css` — remover/adaptar `.demo-wa-expired` (:396-404).
- `internal/app/demo_whatsapp.go` — vista de banner + ruta de dismiss POST+303.
- `internal/app/demo_whatsapp_test.go` (:76) y `web/styles_demo_whatsapp_test.go` (:25,63) — ajustes por el pin de `.demo-wa-expired`.
- Docs (opcional, NO en precedente de commits): `docs/gelium-ui-vocabulary.md` (Banner ✖→✅), `docs/gelium-ui-composition-rules.md` (state matrix §8), `docs/gelium-ui-system-roadmap.md` (matriz :411).

**No tocados**: `web/static/app.js` (0 JS — dismiss por POST+303, sin hook nuevo), `web/styles/tokens.css` y `themes/theme-material/theme.css` (success/info vía color-mix, cero tokens nuevos), `internal/app/{toast,text_field,select}.go` y sus templates.

---

## 7. Fuentes de autoridad

`docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :411, :433), `docs/gelium-ui-vocabulary.md` (:122-129, :142 alias), `docs/gelium-ui-composition-rules.md` (:110-119, :149-162, :164-173), `docs/handoffs/{state-patterns-audit,inline-alert-audit,validation-summary-audit,mozilla-protocol-audit,ux-accessibility-audit}.md`, commits `eba1c4c`/`0688020`/`43c0dac`/`c9afd49` (patrón de archivos verificado con `git show --stat`), `internal/app/{demo_whatsapp,server,routes,toast}.go`, `web/templates/{demo-whatsapp,layout,inline-alert,empty-state,validation-summary}.html`, `web/styles/{demo-whatsapp,inline-alert,empty-state,validation-summary,tokens,app}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/styles_inline_alert_test.go`, `web/styles_demo_whatsapp_test.go`, `web/static/app.css`.
