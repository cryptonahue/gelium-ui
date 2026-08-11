# Gelium UI — Callout Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón CALLOUT (patrón 6 de Phase D, después de Empty State `eba1c4c`, Skeleton `0688020`, Inline Alert `43c0dac`, Validation Summary `c9afd49`, Banner `5787799`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :412), `docs/gelium-ui-vocabulary.md` (:131-138 Callout), `docs/gelium-ui-composition-rules.md` (§4.8 :110-119), `docs/handoffs/{state-patterns-audit,mozilla-protocol-audit,inline-alert-audit,banner-audit,ux-accessibility-audit,vocabulary-audit}.md`, `internal/app/{server,docs,routes,demo_whatsapp}.go`, `web/templates/{demo-whatsapp-admin,demo-whatsapp,layout,inline-alert,empty-state,banner,tooltip}.html`, `web/styles/{demo-whatsapp,inline-alert,empty-state,banner,base,tokens,app}.css`, `web/styles_contract_test.go`, `web/styles_banner_test.go`.

---

## 1. ESTADO ACTUAL

**No existe ningún callout.** Cero coincidencias `callout`/`ui-callout` en `web/` ni `internal/` (grep: sin resultados), ni en el compilado `web/static/app.css`. El vocabulario lo marca ✖ (`vocabulary.md:131-138`); la matriz del roadmap lo marca Phase D "Crear componente (resolver naming Protocol)" (`roadmap.md:412`). Commits de referencia de la serie Phase D: `eba1c4c`, `0688020`, `43c0dac`, `c9afd49`, `5787799` (git log, orden de patrones 1→5).

### 1.1 Ad-hoc `demo-wa-notice` (referencias visuales)

- `web/templates/demo-whatsapp-admin.html:48` — `{{if eq .Quality "YELLOW"}}<p class="demo-wa-notice demo-wa-notice--warn">⚠️ {phone} con calidad YELLOW — puede bajar el throughput.</p>` — **sin role**. Tono warning (`demo-whatsapp.css:510`: `warning-container`/`warning-fg`). Es señal de estado → candidato de **inline-alert/banner**, no de callout.
- `web/templates/demo-whatsapp-admin.html:81` — `<p class="demo-wa-notice">Asignatura activada/desactivada · HMAC-SHA256 · headers … · responder 200 en ≤ 20 s · valida sha256=…</p>` — **sin role**, sin modificador de tono. Es **el caso callout puro del repo**: contenido informativo de ayuda ignorable (explica el contrato del webhook), no es un estado. CSS base `demo-whatsapp.css:504-509` (padding `.5rem .75rem`, `border-radius: var(--ui-radius-sm)`, `font: var(--ui-type-body-sm)`).
- CSS: `web/styles/demo-whatsapp.css:504-510` — `.demo-wa-notice` (base) + `.demo-wa-notice--warn` (tone). Sin surface por defecto (el base solo da typography + radius), sin icono, sin heading.

### 1.2 `role="note"` ad-hoc (relacionado, NO callout)

- `web/templates/demo-whatsapp.html:133` — `<div class="demo-wa-expired" role="note">` (ventana de chat expirada, form interno send-template). Ya clasificado como candidato **banner** (`banner-audit.md:15`, `inline-alert-audit.md:30`): exige acción, `role="note"` flojo. No compite con callout.

### 1.3 Tooltips (referencia de contraste, no callout)

- `web/templates/tooltip.html:9-14` — `.ui-tooltip-host` + `<span class="ui-tooltip" role="tooltip" id="…">` referenciado por `aria-describedby` en el control. Transitorio (hover/focus, `tooltip.css:53`), ligado a un control. CSS `web/styles/tooltip.css:20` + reduced-motion (:127) + forced-colors (:131).

### 1.4 Prose de documentación (contenedor donde viviría)

- `web/styles/base.css:29-40` — `.prose` (max-width 48rem), `h1/h2`, `p` con `color: var(--ui-color-fg-muted)`, `code` con `background: var(--ui-color-surface-container)`. Es el contenedor de contenido de las docs pages (`layout.html:18`: `<article class="prose">{{.Content}}</article>`). Sin estilo de nota/callout hoy.
- Docs pages: server-rendered desde Markdown (`internal/app/server.go:186-208`, `goldmark.New()` con defaults → **sin raw HTML en el Markdown**; `internal/app/docs.go:83-97` genera el índice programáticamente; `content/*.md` embebidos vía `web/assets.go`).

### 1.5 Precedentes de convención (patrones 1-5 entregados)

- Convención exacta por commit (verificada en banner-audit.md:30): `web/templates/<x>.html` + `web/styles/<x>.css` + `web/styles_<x>_test.go` + `@import` en `web/styles/app.css` + `web/static/app.css` regenerado + `web/styles_contract_test.go` (`sourceAppCSS` :24-58 + `TestComponentSizeTokensDeclaredScoped` :405-429). **Sin handler Go, sin ruta, sin docs page** — el partial se entrega como primitiva.

---

## 2. NAMING RESOLUTION

**Canónico: Callout Gelium = tip box (nota contextual informativa ignorable).** "Callout" queda reservado al patrón Gelium; el "Callout" de Protocol (hero full-width promocional) NO se implementa con este nombre — si llega, será composición "Promo Section/Hero" (recomendación de `mozilla-protocol-audit.md:199`).

Evidencia de la colisión resuelta:
- `vocabulary.md:131-138` — el término canónico Gelium: aliases *note, tip, info box*; intención "contenido informativo/promocional sin urgencia ni requisito de acción"; semántica "`<aside>` o bloque con heading opcional"; cuándo usarlo "contexto, tips, documentos, ayuda"; JS 0.
- `mozilla-protocol-audit.md:16,79,99,196` — "Callout" de Protocol = `<section class="mzp-c-callout">` full-width centrado, title+desc+CTA+media, hero-like, variantes fondo/ancho, no para long-form; anatomía opuesta a la nota Gelium.
- `state-patterns-audit.md:39` — contrato propuesto ya fijado: `<aside class="ui-callout">` (contexto complementario, `vocabulary.md:135`) con heading opcional `<h3>` + `<p>` cuerpo + CTA opcional; tono neutral/info preferente; **sin role especial**.
- `roadmap.md:58,78,363,412` — colisión listada y asignada a la resolución en implementación; matriz Phase D :412 "Crear componente (resolver naming Protocol)" — **resuelta en este handoff**.

Naming del contrato:
- Partial: `{{define "callout"}}` → `web/templates/callout.html`.
- Clase raíz: `.ui-callout` + modificador de variante `.ui-callout--{info|tip}`.
- En docs, "Callout" se documenta como el término Gelium; el patrón Protocol se menciona como "Promo Section" (alias futuro, composición, no primitiva).

---

## 3. DISTINCION CANONICA

| Señal | Persistencia | Scope | Rol ARIA | Origen | ¿Ignorable? | Estado hoy |
|---|---|---|---|---|---|---|
| **Callout** (patrón 6) | Persistente | Contenido (inline en docs/sección) | **Sin rol especial** (`<aside>`, heading nativo) | Nota informativa: contexto, tips, ayuda | **Sí, se puede ignorar** | ✖ (ad-hoc `.demo-wa-notice` admin:81, sin role) |
| **Inline alert** (patrón 3) | Persistente | Sección/formulario | `alert` (error) / `status` (info/success/warning) | Validación 422, advertencia de sección | No (señala estado del form) | ✖ genérico (solo campo ✅ + notices ad-hoc) |
| **Banner** (patrón 5) | Persistente | Página/sitio, full-width | `alert`/`status`, dismiss POST+303 | Sesión, mantenimiento, consent | No (exige acción) | ✖ (ad-hoc `.demo-wa-expired` :133) |
| **Tooltip** | Transitorio | Hover/focus en un control | `tooltip` + `aria-describedby` en el control | Contexto del control | — | ✅ completo |

Reglas (`composition-rules.md:110-119`, `vocabulary.md:128,137`): "cuándo no" callout = requiere acción → Banner; error del campo → Inline alert. Nota ignorable → Callout. Anti-regla general: nada persistente se anuncia con `loom:toast` (`state-patterns-audit.md:45`).

**Decisión de diseño**: el Callout NO es un patrón de estado — es contenido. Por eso: sin role, sin tones de estado (error/warning/success quedan en inline-alert/banner, que sí portan role), sin dismiss, sin CTA obligatoria. Comparte con inline-alert/banner la superficie y la tipografía, pero la semántica es de documento, no de señal.

---

## 4. CONTRATO PROPUESTO

### 4.1 Markup (patrón banner.html / inline-alert.html — partial server-rendered)

```html
{{define "callout"}}<aside class="ui-callout{{if .Variant}} ui-callout--{{.Variant}}{{end}}">
  {{if .Icon}}<span class="ui-callout-icon" aria-hidden="true">{{.Icon}}</span>{{end}}
  <div class="ui-callout-content">
    {{if .Heading}}<h3 class="ui-callout-heading">{{.Heading}}</h3>{{end}}
    <p class="ui-callout-body">{{.Body}}</p>
  </div>
  {{if .CTA}}<a class="ui-button" href="{{.CTAHref}}">{{.CTALabel}}</a>{{end}}
</aside>{{end}}
```

- Root **`<aside>`** (semántica HTML nativa de contenido complementario, `vocabulary.md:135`; `state-patterns-audit.md:39`). Distinto del `<div>` de banner/inline-alert porque callout no es un slot de estado sino contenido contextual.
- Icono opcional `aria-hidden="true"` (glyph trusted, precedente inline-alert.html:2 / empty-state.html:2).
- Heading opcional **`<h3>`** real (nativo, aporta semántica; `state-patterns-audit.md:39`). Nota de tensión: banner/inline-alert usan `<p>` para título y evitan saltos de jerarquía (`banner-audit.md:69`); en callout el `<h3>` es correcto porque vive en contenido (docs) donde la jerarquía importa. Si un consumidor no quiere salto de heading, omite Heading y el body se sostiene solo.
- Cuerpo `<p class="ui-callout-body">` con `font: var(--ui-type-body-sm)`.
- CTA opcional = **`<a class="ui-button">` real** con href (precedente empty-state.html:5, banner.html:7). Sin dismiss (no es señal, no hay estado que limpiar).

### 4.2 Variantes (sin error/warning — los cubren inline-alert y banner)

| Variante | Mecánica | Consumidor real |
|---|---|---|
| **default (neutra)** | `--ui-callout-bg: var(--ui-color-surface-container)` / `--ui-callout-fg: var(--ui-color-fg)` | Nota genérica de contexto; hoy: `.demo-wa-notice` admin:81 (neutro) |
| `--info` | `color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface))` / fg `--ui-color-info` | Nota informativa en docs/componentes (misma receta que `.ui-inline-alert--info`, `inline-alert.css:32`) |
| `--tip` (opcional) | `color-mix(in srgb, var(--ui-color-secondary) 20%, var(--ui-color-surface))` / fg `--ui-color-fg` | Highlight de tip/ayuda **fuera de la paleta de estado** (refuerza "no es una señal") |

- `--tip` usa `--ui-color-secondary` (existe en core, `tokens.css:29`) precisamente para NO reusar la paleta de estado: el callout no anuncia nada. `--info` es el único tone de la familia de estado justificable (información ≠ estado), y se ofrece por el precedente del roadmap ("tono neutral/info preferente", `state-patterns-audit.md:39`).
- **Intencionalmente ausentes**: `error`/`warning`/`success` → ya portados por inline-alert (sección) y banner (página) con sus roles ARIA. Variantes solo con texto, nunca color-only (convención del sistema, `ux-accessibility-audit.md`).

### 4.3 ARIA

- **Sin role especial** en el root. `<aside>` sin `role="complementary"` NO crea landmark en el a11y tree (es elemento genérico), y eso es lo correcto: el callout es contenido ignorable, no una región que anunciar. `role="alert"`/`status` quedan para inline-alert/banner.
- **`aria-labelledby` NO** cuando hay heading: sobre un `<aside>` sin role el atributo es inerte (solo aplica a elementos con rol/landmark); el `<h3>` nativo interior ya porta el nombre y la estructura. HTML nativo antes que ARIA.
- Icono `aria-hidden="true"` (decorativo). Nada color-only: heading/cuerpo de texto + bloque forced-colors (convención `ux-accessibility-audit.md`).

### 4.4 Tokens

Core (`web/styles/tokens.css`) — **verificados**:
- `--ui-color-surface-container` (:23), `--ui-color-fg` (:24), `--ui-color-fg-muted` (:25), `--ui-color-secondary` (:29), `--ui-color-info` (:50).
- **NO existe** `info-container` (ni core ni `themes/theme-material/theme.css`) → info vía `color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface))`, precedente `.ui-inline-alert--info` (`inline-alert.css:32`), `.ui-banner--info` (`banner.css:32`), `.demo-wa-quality--GREEN` (`demo-whatsapp.css:501`). Cero cambios en tokens.css/theme.
- Espacio: `--ui-space-2/3/4` (escala core, `styles_contract_test.go:263-270`). Radius: `--ui-radius-sm` (core :80). Type: `--ui-type-body-sm` (cuerpo, en theme-material; usado por empty-state/banner/inline-alert) y heading `<h3>` (hereda la tipografía del contenido — sin token de heading específico, precedente de composición en prose). Size: `--ui-size-icon-sm` (core :123, mismo que banner/inline-alert).

Tokens scoped `--ui-callout-*` (declarados en el root, patrón `banner.css:11-18`):

```css
.ui-callout {
  --ui-callout-padding: var(--ui-space-3) var(--ui-space-4);
  --ui-callout-gap: var(--ui-space-2);
  --ui-callout-radius: var(--ui-radius-sm);
  --ui-callout-icon-size: var(--ui-size-icon-sm);
  --ui-callout-bg: var(--ui-color-surface-container);
  --ui-callout-fg: var(--ui-color-fg);
  --ui-callout-heading-color: var(--ui-color-fg);
  --ui-callout-body-color: var(--ui-color-fg-muted);
}
.ui-callout--info { --ui-callout-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface)); --ui-callout-fg: var(--ui-color-info); }
.ui-callout--tip  { --ui-callout-bg: color-mix(in srgb, var(--ui-color-secondary) 20%, var(--ui-color-surface)); --ui-callout-fg: var(--ui-color-fg); }
```

### 4.5 CSS

- `web/styles/callout.css` nuevo, `@layer components`, anatomía de `banner.css`/`inline-alert.css`: tokens scoped en el root, 2 modificadores de variante (`--info`, `--tip`), `.ui-callout-icon` (width/height `--ui-callout-icon-size`, `flex: none`), `.ui-callout-content` (`flex: 1` o column), `.ui-callout-heading` (`color: var(--ui-callout-heading-color)`), `.ui-callout-body` (`font: var(--ui-type-body-sm)`, `color: var(--ui-callout-body-color)`). Bloque `@media (forced-colors: active)` (borde `CanvasText`; sin tone específico — no hay tones de estado que preservar; icon/heading/body → `CanvasText`, precedente `banner.css:67-82`). **Sin animación → sin bloque reduced-motion**.
- `web/styles/app.css`: `@import "./callout.css";` (tras `validation-summary.css` :38, junto a los patrones de estado).
- `web/styles_contract_test.go`: agregar `"styles/callout.css"` a `sourceAppCSS` (tras `"styles/validation-summary.css"` :60, antes de `"styles/demo-whatsapp.css"` :61) y `"callout.css": "--ui-callout-icon-size:"` a `TestComponentSizeTokensDeclaredScoped` (:405-429).
- Rebuild: `npm run build` regenera `web/static/app.css` (precedente commits 43c0dac/c9afd49/5787799).

---

## 5. INTEGRACION

- **Es contenido, server-rendered como los demás partials**: 0 JS end-to-end, se compone con `{{template "callout" .Callout}}` desde cualquier template de página o slot del layout (`layout.html:16` ya tiene el precedente `{{if .Banner}}…{{end}}`).
- **Dónde se usaría**:
  1. **Docs de componentes** (páginas `/components/*` vía `content/*.md` + `.prose`): notas de "cuándo usar / cuándo no", tips de composición. El `<aside>` encaja dentro de `<article class="prose">` (`layout.html:18`). **Ojo**: las docs pages se generan con `goldmark.New()` en defaults (`server.go:101,190-195`) → **sin raw HTML en el Markdown**; un callout inline en un `.md` NO renderizaría hoy. La vía limpia es un slot del layout (`pageView.Callout` → `{{if .Callout}}{{template "callout" .Callout}}{{end}}` dentro del article) o el generador de `docs.go`, no HTML crudo en Markdown. Documentar como decisión de integración.
  2. **Contenido de ayuda en pantallas**: el caso real de hoy es `demo-whatsapp-admin.html:81` (explicación del contrato de webhook) — migración opcional del ad-hoc `.demo-wa-notice` neutro al partial.
  3. **Settings/paneles de admin** (Phase G): ayudas de configuración, explicaciones de campos sensibles.
- **No es primitiva con ruta propia** (mismo criterio que banner/validation-summary — "primitiva lista para Phase G", `banner-audit.md:153`): sin handler Go, sin ruta `/components/callout`, sin docs page. `pageView` ganaría el campo solo cuando un consumidor lo wirree (opcional en Phase D).
- El ad-hoc admin:48 (`--warn`) NO migra a callout — es señal de estado → inline-alert o banner.

---

## 6. TESTS PROPUESTOS

1. **`web/styles_callout_test.go`** (nuevo, patrón `styles_banner_test.go`):
   - `TestCalloutPrimitiveCSSMapsTokens` — contratos exactos: `.ui-callout {`, `display: flex`, `gap: var(--ui-callout-gap)`, `padding`, `border-radius`, `background`, `color`, modificadores `--info`/`--tip`, icon/heading/body, todos los tokens scoped. Sin `transition:`/`animation:`/`prefers-reduced-motion`.
   - `TestCalloutContractCSSWired` — compilado embebido (`Assets.ReadFile("static/app.css")`) contiene `.ui-callout`, `.ui-callout--info`, `@media (forced-colors:active)`.
   - `TestCalloutClassVocabularyIsClosed` — clases template ↔ selectores CSS (`ui-callout`, `ui-callout-icon`, `ui-callout-content`, `ui-callout-heading`, `ui-callout-body`); interpolación `ui-callout--`; prohibido `ui-callout-demo`.
   - `TestCalloutTonesUseCoreTokens` — `--info`/`--tip` resuelven solo tokens core o color-mix sobre ellos; nunca hex literal (guard del blanket `TestNoColorLiteralsInComponents`).
2. **`web/styles_contract_test.go`** (modificar): `sourceAppCSS` + `TestComponentSizeTokensDeclaredScoped` (+`--ui-callout-icon-size:`).
3. **(Opcional, render Go)** `internal/app/callout_test.go` — root `<aside>` **sin role**; heading `<h3>` presente/ausente según `.Heading`; icono `aria-hidden`; CTA como `<a class="ui-button" href>`; body siempre presente; sin `aria-labelledby` (contrato: heading nativo).
4. **(Opcional, si se migra el ad-hoc)** ajustes en `internal/app/demo_whatsapp_test.go` y `web/styles_demo_whatsapp_test.go` por el pin de `.demo-wa-notice` (admin:48,81).

---

## 7. FILES IMPACTADOS (solo read-only)

**Nuevos**:
- `web/templates/callout.html` — partial `{{define "callout"}}`.
- `web/styles/callout.css` — primitiva + variantes `--info`/`--tip` + tokens scoped + forced-colors.
- `web/styles_callout_test.go` — tests de contrato (patrón banner/inline-alert).
- `docs/handoffs/callout-audit.md` — este handoff.

**Modificados**:
- `web/styles/app.css` — `@import "./callout.css";` (+1, tras `validation-summary.css` :38).
- `web/static/app.css` — regenerado por `npm run build` (precedente commits 43c0dac/c9afd49/5787799).
- `web/styles_contract_test.go` — `sourceAppCSS` (:24-58) + `TestComponentSizeTokensDeclaredScoped` (:405-429).

**Opcionales (migración del ad-hoc / wiring Phase G)**:
- `web/templates/demo-whatsapp-admin.html` — `.demo-wa-notice` :81 → partial callout (solo el caso neutro; :48 va a inline-alert/banner).
- `web/styles/demo-whatsapp.css` — remover/adaptar `.demo-wa-notice` (:504-510).
- `internal/app/demo_whatsapp_test.go` y `web/styles_demo_whatsapp_test.go` — ajustes por el pin de `.demo-wa-notice`.
- `internal/app/server.go` — campo `Callout *calloutView` en `pageView` (:50-88) solo cuando un consumidor lo wirree.
- Docs (opcional, NO en precedente de commits): `docs/gelium-ui-vocabulary.md` (Callout ✖→✅), `docs/gelium-ui-composition-rules.md` (state matrix §8), `docs/gelium-ui-system-roadmap.md` (matriz :412).

**No tocados**: `web/static/app.js` (0 JS), `web/styles/tokens.css` y `themes/theme-material/theme.css` (info/tip vía color-mix, cero tokens nuevos), `internal/app/{toast,text_field,select}.go` y sus templates.

---

## 8. Fuentes de autoridad

`docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :412), `docs/gelium-ui-vocabulary.md` (:131-138), `docs/gelium-ui-composition-rules.md` (:110-119), `docs/handoffs/{state-patterns-audit,mozilla-protocol-audit,inline-alert-audit,banner-audit,ux-accessibility-audit,vocabulary-audit}.md`, commits `eba1c4c`/`0688020`/`43c0dac`/`c9afd49`/`5787799` (patrón de archivos verificado con git log), `internal/app/{server,docs,routes,demo_whatsapp}.go`, `web/templates/{demo-whatsapp-admin,demo-whatsapp,layout,inline-alert,empty-state,banner,tooltip}.html`, `web/styles/{demo-whatsapp,inline-alert,empty-state,banner,base,tokens,app}.css`, `web/styles_contract_test.go`, `web/styles_banner_test.go`, `web/static/app.css`.
