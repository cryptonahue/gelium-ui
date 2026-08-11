# Gelium UI — Inline Alert Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón INLINE ALERT (patrón 3 de Phase D, después de Empty State `eba1c4c` y Skeleton `0688020`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D), `docs/gelium-ui-vocabulary.md` (§In line alert), `docs/gelium-ui-composition-rules.md` (§4.8), `docs/handoffs/state-patterns-audit.md`, `internal/app/{text_field,select,toast,demo_whatsapp,server,routes,data_table}.go`, `web/templates/{text-field,select,empty-state,skeleton,demo-whatsapp-admin,demo-whatsapp}.html`, `web/styles/{text-field,select-menu,empty-state,skeleton,toast,demo-whatsapp,tokens,app}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/static/app.js`.

---

## 1. ESTADO ACTUAL

### 1.1 Error de campo (text-field) — específico del control, NO reusable

- `web/templates/text-field.html:5` — input/textarea condicionado a `.Error && not .Disabled`: `aria-invalid="true"` + `aria-describedby="{ID}-error"` (o `-help` si solo hay helper).
- `web/templates/text-field.html:6` — `<svg class="ui-text-field-error-icon" aria-hidden="true" focusable="false">` (glyph decorativo de error, trusted interno).
- `web/templates/text-field.html:8` — `<p class="ui-text-field-message" id="{ID}-error" role="alert"><strong>Error:</strong> {{.Error}}</p>`. El mismo `<p class="ui-text-field-message">` reusa el helper con `role="{{.MessageRole}}"` (`"status"` en success).
- CSS: `web/styles/text-field.css:60-70` (`.ui-text-field-error-icon`), `:116-133` (estados `.ui-text-field-error` con `var(--ui-field-error)`), `:134` (`.ui-text-field-message`: `padding: var(--ui-space-1) var(--ui-space-4) 0`, `color: var(--ui-color-fg-muted)`, `font: var(--ui-type-body-sm)`). Sin surface, sin variantes de tone, sin contenedor.
- Handler: `internal/app/text_field.go:55-92` — `validateTextField`. Vacío → `field.Error = "Name is required"` (66), `status = 422` (68). HX: fragmento `validation-form` + `X-Loom-Validation: true` (87-89). No-HX: página completa vía `renderMarkdownPageStatus(..., 422)` (77). Success: `field.Helper = "Name accepted"` + `field.MessageRole = "status"` (70-71).

### 1.2 Error de campo (select menu)

- `web/templates/select.html:89` — `<p class="ui-select-menu-error" role="alert">{{.Error}}</p>` dentro del form del menú.
- CSS: `web/styles/select-menu.css:67-72` — `.ui-select-menu-error`: `color: var(--ui-select-error)` (alias de theme → `--ui-field-error`), `font: var(--ui-type-body-sm)`. También específico, sin surface.
- Handler: `internal/app/select.go:72-124` — `selectMenu`. Valor desconocido → `status = 422` (95) + `demo.Error = "Select a valid option"` (96); mismo bifurcación HX/fragmento + `X-Loom-Validation` (112-113).

### 1.3 Notices ad-hoc de sección (sin role, no reusables)

- `web/templates/demo-whatsapp-admin.html:48` — `<p class="demo-wa-notice demo-wa-notice--warn">⚠️ ... calidad YELLOW ...` (advertencia de sección, **sin role ARIA**).
- `web/templates/demo-whatsapp-admin.html:81` — `<p class="demo-wa-notice">Asignatura activada/desactivada ... HMAC-SHA256 ...` (info de sección, **sin role ARIA**).
- CSS: `web/styles/demo-whatsapp.css:504-510` — `.demo-wa-notice` (padding .5rem/.75rem, radius-sm, body-sm) + `.demo-wa-notice--warn` (`background: var(--ui-color-warning-container); color: var(--ui-color-warning-fg)`). **Este es el único precedente con surface+tone del repo**, pero vive en el demo, no en la librería.
- Relacionado (borde de Banner, no inline alert): `web/templates/demo-whatsapp.html:132-144` — `.demo-wa-expired` con `role="note"`, CSS `demo-whatsapp.css:396-404` (`--ui-color-error` en título).

### 1.4 Helper con role="status"

- `internal/app/text_field.go:70-71` — success: `field.Helper = "Name accepted"` + `field.MessageRole = "status"` → renderiza en `text-field.html:8` como `.ui-text-field-message` con `role="status"`. Es el único `role="status"` persistente del sistema (fuera de toast/empty/skeleton).

### 1.5 Contrato 422 (base de integración)

- `internal/app/text_field.go:64-68,87-89` y `internal/app/select.go:94-96,112-113`: `422` + `X-Loom-Validation: true` solo en HX. Hook cliente `web/static/app.js:1-9`: `htmx:beforeSwap` swapea 422 con header, sin marcarlo error. No-JS: re-render de página completa con status 422. Validación NUNCA toast (`internal/app/toast.go:129-133`, `toast.go` comentario).

**Conclusión**: todo el error/success actual es campo-specific (`ui-text-field-message`, `ui-select-menu-error`) o demo ad-hoc (`demo-wa-notice`) — sin clases genéricas reusables, sin tones, sin surface, sin role consistente fuera del campo.

---

## 2. DISTINCIÓN CANÓNICA

| Señal | Persistencia | Scope | Rol ARIA | Origen | Estado hoy |
|---|---|---|---|---|---|
| **Inline alert** (patrón 3) | Persistente | Sección/formulario | `alert` (error) / `status` (info/success/warning) | Validación 422, advertencia de sección, success persistente | ✖ genérico (solo campo + notices ad-hoc) |
| **Error de campo** (existente) | Persistente | Un control | `alert` + `aria-invalid` + `aria-describedby` en el control | Validación 422 | ✅ (text-field, select) |
| **Toast** | Transitorio | Región `aria-live` global | `status`/`alert` según tipo | Resultado de acción (`loom:toast`) | ✅ completo |
| **Banner** (patrón 5) | Persistente página/sitio | Página/layout | `alert`/`status`, dismiss POST+303 | Sesión, mantenimiento, aviso global | ✖ (ad-hoc `.demo-wa-expired`) |
| **Callout** (patrón 6) | Persistente ignorable | Contenido | Sin rol especial | Nota informativa | ✖ |

Reglas (`composition-rules.md:110-119`, `toast.go:129-133`): validación nunca toast; feedback persistente/crítico nunca toast; persistente-contextual ≠ transitorio-de-acción (`roadmap.md:162-168`).

**Decisión de diseño**: el Inline alert NO reemplaza al error de campo. Coexisten por capas:
- Error de campo = señal del control (`aria-invalid` + `aria-describedby` + mensaje bajo el campo, sin surface).
- Inline alert genérico = señal de la sección/formulario completo (con surface contenedora + tone), p.ej. summary de errores arriba del form, advertencia de sección, success persistente post-guardado.
- Cuando hay un control asociado, el inline alert de sección NO repite `aria-describedby` (eso lo mantiene el campo); opcionalmente puede llevar un `id` para que la región/fieldset lo referencie.

---

## 3. CONTRATO PROPUESTO

### 3.1 Markup (siguiendo empty-state.html — partial server-rendered)

```html
{{define "inline-alert"}}<div class="ui-inline-alert ui-inline-alert--{{.Tone}}" role="{{if eq .Tone "error"}}alert{{else}}status{{end}}">
  {{if .Icon}}<span class="ui-inline-alert-icon" aria-hidden="true">{{.Icon}}</span>{{end}}
  {{if .Title}}<p class="ui-inline-alert-title">{{.Title}}</p>{{end}}
  <p class="ui-inline-alert-body">{{.Body}}</p>
</div>{{end}}
```

- Root `<div>` (contenedor de superficie). Icono opcional `aria-hidden` (glyph trusted, igual que empty-state/text-field). Título opcional como `<p>` con `font: var(--ui-type-title-md)` (precedente `ui-empty-state-title`, evita saltos de heading).
- Tones: `.ui-inline-alert--error|--success|--info|--warning` (vocabulario cerrado, misma familia que toast `KNOWN_TYPES` en `app.js:19` y `sanitizeToastType`).
- Nada de color-only: texto + icono opcional + forced-colors.

### 3.2 ARIA

- `role="alert"` (assertive) solo en `--error`. `role="status"` (polite) en `--info` / `--success` (y `--warning` salvo que la advertencia exija acción — en ese caso el consumidor decide; contract base: status).
- El control con error mantiene su propio `aria-invalid` + `aria-describedby` (el inline alert NO se relaciona con el control por `aria-describedby`; es contenedor, no control).
- Opcional: `id` en el root para que un `fieldset`/región lo referencie; `aria-live` ya lo aporta el rol.

### 3.3 Tokens

Core (`web/styles/tokens.css`) — verificados:
- `--ui-color-danger` (:32), `--ui-color-danger-container` (:51), `--ui-color-warning` (:47), `--ui-color-warning-fg` (:48), `--ui-color-warning-container` (:49), `--ui-color-info` (:50), `--ui-color-success` (:45), `--ui-color-success-fg` (:46). `--ui-color-error` = alias de danger (:38).
- **NO existen**: `success-container`, `info-container`. El theme (`themes/theme-material/theme.css:25-33`) tampoco (solo `warning-container` :30 y `danger-container` :32). `danger-fg` existe pero es blanco para botones — sobre `danger-container` claro el fg correcto es `--ui-color-danger` (precedente `.demo-wa-quality--RED`, `demo-whatsapp.css:503`).
- Precedente de surface+tone con color-mix: `.demo-wa-quality--GREEN` (`demo-whatsapp.css:501`): `color-mix(in srgb, var(--ui-color-success) 12%, var(--ui-color-surface))` + `var(--ui-color-success)`.
- Espacio/forma/type: `--ui-space-2/3/4`, `--ui-radius-sm`, `--ui-type-title-md`, `--ui-type-body-sm`, `--ui-size-icon` — todos en core.

Tokens scoped `--ui-inline-alert-*` (declarados en el root, patrón empty-state/skeleton):

```css
.ui-inline-alert {
  --ui-inline-alert-padding: var(--ui-space-3) var(--ui-space-4);
  --ui-inline-alert-gap: var(--ui-space-2);
  --ui-inline-alert-radius: var(--ui-radius-sm);
  --ui-inline-alert-bg: var(--ui-color-surface-container);
  --ui-inline-alert-fg: var(--ui-color-fg);
  --ui-inline-alert-icon-size: var(--ui-size-icon);
  --ui-inline-alert-title-color: var(--ui-color-fg);
  --ui-inline-alert-body-color: var(--ui-color-fg-muted);
}
.ui-inline-alert--error   { --ui-inline-alert-bg: var(--ui-color-danger-container);  --ui-inline-alert-fg: var(--ui-color-danger); }
.ui-inline-alert--warning { --ui-inline-alert-bg: var(--ui-color-warning-container); --ui-inline-alert-fg: var(--ui-color-warning-fg); }
.ui-inline-alert--success { --ui-inline-alert-bg: color-mix(in srgb, var(--ui-color-success) 12%, var(--ui-color-surface)); --ui-inline-alert-fg: var(--ui-color-success); }
.ui-inline-alert--info    { --ui-inline-alert-bg: color-mix(in srgb, var(--ui-color-info) 12%, var(--ui-color-surface));    --ui-inline-alert-fg: var(--ui-color-info); }
```

> Alternativa para success/info: añadir `--ui-color-success-container` / `--ui-color-info-container` al core + theme (matriz light + dark-class + dark-media). La regla del core "solo tokens con consumidores reales" (`tokens.css:40-44`, `styles_contract_test.go:592-601`) lo permite si el tone tiene consumidor real; color-mix evita tocar tokens.css/theme/tests de contrato existentes. Recomendación: **color-mix** (precedente GREEN, cero cambios de contrato); dejar los tokens container nuevos para cuando un consumidor real los exija.

### 3.4 CSS

- `web/styles/inline-alert.css` nuevo, `@layer components`, siguiendo la anatomía de `empty-state.css`: tokens scoped en el root, modificadores de tone, icono/título/cuerpo, bloque `@media (forced-colors: active)` (borde `CanvasText`, `Mark` en error, icono `forced-color-adjust: auto`). Sin animación → sin bloque reduced-motion.
- `web/styles/app.css`: `@import "./inline-alert.css";` (entre `skeleton.css` :34 y `tooltip.css` :35).
- `web/styles_contract_test.go`: agregar `"styles/inline-alert.css"` a `sourceAppCSS` (:24-58, tras `skeleton.css`), y `"inline-alert.css": "--ui-inline-alert-icon-size:"` a `TestComponentSizeTokensDeclaredScoped` (:405-421). El contador de `TestNoColorLiteralsInComponents` (`checked < 25`, :740) sigue verde (32 → 33 archivos).
- Rebuild: `npm run build` regenera `web/static/app.css` (minificado; `package.json` `build`).

---

## 4. INTEGRACIÓN

Encaje con el contrato 422 existente (sin inventar contratos, `state-patterns-audit.md` §4):

1. **Formulario multi-campo (summary de errores)**: el handler de validación (patrón `text_field.go:55-92` / `select.go:72-124`) agrega arriba del form un `.ui-inline-alert--error` con `role="alert"` (p.ej. "N campos requieren atención") y mantiene los errores por campo intactos. HX: fragmento + 422 + `X-Loom-Validation` (el hook `app.js:1-9` no cambia). No-HX: página completa con status 422. Los links `#campo-error` quedan para el patrón 4 (Validation summary), que se compone SOBRE este inline alert.
2. **Success persistente (settings guardado)**: POST + 303 (patrón `demo_whatsapp.go:559,573`) → la página destino re-renderiza `.ui-inline-alert--success` con `role="status"` (polite, persistente). En HX, el fragmento post-submit incluye el success. NUNCA `loom:toast` (anti-regla `composition-rules.md:126`).
3. **Advertencia de sección**: reemplaza los notices ad-hoc de `demo-whatsapp-admin.html:48` (`--warning`, calidad YELLOW) y `:81` (`--info`, firma HMAC) — ambos sin role hoy, ganan `role="status"`. Migración: quitar `.demo-wa-notice`/`.demo-wa-notice--warn` de `demo-whatsapp.css:504-510` (no hay test que los pinne; `.demo-wa-expired` SÍ está testeado en `web/styles_demo_whatsapp_test.go:25,63` e `internal/app/demo_whatsapp_test.go:76`, y es candidato a Banner/patrón 5, no a inline alert).
4. **El error de campo existente NO cambia**: `text-field.html:5-8`, `select.html:89`, `text_field.go` y `select.go` quedan igual. El inline alert es una capa de sección que compone con ellos.

0 JS end-to-end (server-rendered, `app.js` intacto). Todo lo persistente es output del servidor.

---

## 5. TESTS PROPUESTOS

1. `web/styles_inline_alert_test.go` (nuevo, patrón `styles_empty_state_test.go` + `styles_skeleton_test.go`):
   - `TestInlineAlertPrimitiveCSSMapsTokens` — contratos exactos: `.ui-inline-alert {`, `display: flex;`, `align-items: flex-start;`, `gap: var(--ui-inline-alert-gap);`, `padding: var(--ui-inline-alert-padding);`, `border-radius: var(--ui-inline-alert-radius);`, `background: var(--ui-inline-alert-bg);`, `color: var(--ui-inline-alert-fg);`, los 4 modificadores de tone, icono/título/cuerpo, todos los tokens scoped. Sin `transition:`/`animation:`/`prefers-reduced-motion` (primitiva estática, precedente empty-state).
   - `TestInlineAlertContractCSSWired` — compilado embebido (`Assets.ReadFile("static/app.css")`) contiene `.ui-inline-alert`, los 4 modificadores, `@media (forced-colors:active)`.
   - `TestInlineAlertClassVocabularyIsClosed` — cada clase en template ↔ selector en CSS; prohibido `ui-inline-alert-demo`.
   - `TestInlineAlertTonesUseCoreTokens` — cada tone referencia solo tokens core (`--ui-color-danger*`, `--ui-color-warning*`, `--ui-color-success`, `--ui-color-info`) y nunca un hex literal (guard del blanket `TestNoColorLiteralsInComponents`).
2. `web/styles_contract_test.go` (modificar): `sourceAppCSS` + `TestComponentSizeTokensDeclaredScoped` (+`--ui-inline-alert-icon-size:`).
3. Render (si se integra en Go, patrón `internal/app/text_field_test.go`): `internal/app/inline_alert_test.go` — assertions de `role="alert"` vs `role="status"` por tone, icono `aria-hidden`, título/cuerpo. Y si el admin migra: `internal/app/demo_whatsapp_test.go` + `web/styles_demo_whatsapp_test.go` ajustes por remoción de `.demo-wa-notice` (hoy no está testeado).

---

## 6. FILES IMPACTADOS (solo read-only)

**Nuevos**:
- `web/styles/inline-alert.css` — primitiva + tones + forced-colors.
- `web/styles_inline_alert_test.go` — tests de contrato (patrón empty-state/skeleton).
- `web/templates/inline-alert.html` — partial `{{define "inline-alert"}}`.
- (opcional) `internal/app/inline_alert.go` + `internal/app/inline_alert_test.go` — view model + render tests.
- (opcional) `web/content/inline-alert.md` — página docs si se le da ruta.

**Modificados**:
- `web/styles/app.css` — `@import "./inline-alert.css";` (+1).
- `web/static/app.css` — regenerado por `npm run build` (2 líneas, precedente commits eba1c4c/0688020).
- `web/styles_contract_test.go` — `sourceAppCSS` (:24-58) + `TestComponentSizeTokensDeclaredScoped` (:405-421).
- `web/templates/demo-whatsapp-admin.html` — migrar notices :48 y :81 a `.ui-inline-alert--warning` / `--info`.
- `web/styles/demo-whatsapp.css` — remover `.demo-wa-notice`/`--warn` (:504-510) si se migra el admin.
- `internal/app/demo_whatsapp.go` — (opcional) view data del alert del admin.
- Docs (opcional, NO en precedente de los commits eba1c4c/0688020): `docs/gelium-ui-vocabulary.md` (Inline alert ✖→✅), `docs/gelium-ui-composition-rules.md` (state matrix §4.8/§5), `docs/gelium-ui-system-roadmap.md` (matriz :410).

**No tocados**: `web/static/app.js` (0 JS), `internal/app/{text_field,select,toast}.go` y sus templates (error de campo se mantiene), `web/styles/{tokens,toast}.css` si se elige color-mix.

---

## 7. Fuentes de autoridad

`docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :410), `docs/gelium-ui-vocabulary.md` (:112-120), `docs/gelium-ui-composition-rules.md` (:110-119, :162), `docs/handoffs/state-patterns-audit.md` (:37, §4, §5), commits `eba1c4c` y `0688020` (patrón de archivos), `internal/app/{text_field,select,toast,demo_whatsapp,server,routes}.go`, `web/templates/{text-field,select,empty-state,skeleton,demo-whatsapp-admin}.html`, `web/styles/{text-field,select-menu,empty-state,skeleton,toast,demo-whatsapp,tokens,app}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/static/app.js`, `web/assets.go`, `package.json`.
