# Gelium UI — Validation Summary Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón VALIDATION SUMMARY (patrón 4 de Phase D, después de Empty State `eba1c4c`, Skeleton `0688020`, Inline Alert `43c0dac`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :408-415), `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-composition-rules.md` (§4.8, §8 state matrix, §9 server-driven), `docs/handoffs/{state-patterns-audit,inline-alert-audit}.md`, `internal/app/{text_field,select,toast,server,routes,data_table}.go`, `web/templates/{text-field,select,toast,inline-alert,empty-state,skeleton,layout}.html`, `web/styles/{tokens,inline-alert,empty-state,app}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/styles_inline_alert_test.go`, `web/static/app.js`.

---

## 1. ESTADO ACTUAL — cómo se reportan errores hoy

### 1.1 No existe ningún resumen de errores de formulario

- Cero referencias `validation-summary` / `ValidationSummary` en código (solo en docs). El audit previo lo marca ✖ (`state-patterns-audit.md:20,41,56,73,91`); la matriz del roadmap lo marca como "Crear contrato" (`roadmap.md:414`); la state matrix de composición marca Form como "422 inline" sin resumen (`composition-rules.md:160`).

### 1.2 Errores por campo (única señal de validación existente)

- `web/templates/text-field.html:5` — input/textarea con `.Error && !.Disabled`: `aria-invalid="true"` + `aria-describedby="{ID}-error"` (o `-help` si solo hay helper).
- `web/templates/text-field.html:8` — `<p class="ui-text-field-message" id="{ID}-error" role="alert"><strong>Error:</strong> {{.Error}}</p>`. El mismo `<p>` reusa el helper con `role="{{.MessageRole}}"` (`"status"` en success).
- `web/templates/select.html:89` — `<p class="ui-select-menu-error" role="alert">{{.Error}}</p>` dentro del form del menú (contenido, **sin id**).
- CSS: `web/styles/text-field.css` (`.ui-text-field-error`, `.ui-text-field-message`) y `web/styles/select-menu.css:67-72` (`.ui-select-menu-error`). Sin surface contenedora, sin variantes, específicos de control.

### 1.3 Ids reales de los campos y de los mensajes de error (anclas posibles)

| Form | Campo | id del campo | id del error (target de ancla) |
|---|---|---|---|
| Text-field validation demo | Name | `validation-name` | `validation-name-error` (`text-field.html:8`) |
| Text-field docs (estático) | Username | `text-error` | `text-error-error` |
| Select menu demo | (value oculto) | form `select-menu-field`, dialog `select-menu` | **NO existe id** en `select.html:89` — gap |
| Toast demo | Message | `toast-message` | `toast-message-error` (vía text-field) |

- Los mensajes de error del text-field SÍ tienen id real `{ID}-error` → un `<a href="#{ID}-error">` apunta a cada campo correctamente. El error del select menu (`select.html:89`) NO tiene id: si un summary quisiera enlazarlo habría que añadirle `id="{ID}-error"`.

### 1.4 Autofocus / scroll al error

- `internal/app/text_field.go:67` — `field.Autofocus = !isHX`: el re-render **no-JS** de página completa con 422 añade `autofocus` al campo con error (recuperación de foco, testeado en `text_field_test.go:216`). El fragmento **HX** NO lleva autofocus (`text_field_test.go:183-185` — "HX 422 fragment must not add autofocus"). Mismo patrón en `toast.go:149`.
- No existe scroll programático al error; los `autofocus` nativos cubren el caso no-JS. Para HX el summary con anclas reales es el mecanismo de navegación (click → salto nativo al campo, sin JS).

### 1.5 ¿Hay algún formulario con 2+ campos validados?

**No.** Inventario de forms con validación server-side:

- `validation-form` (text-field demo, `text-field.html:12-17`): 1 campo validado (`name`).
- `select-menu-demo` (`select.html:73-92`): 1 campo oculto validado (`value`).
- `toast-demo-form` (`toast.html:13-22`): 1 campo validado (`message`); el `<select>` de tipo NO se valida.
- El resto de forms son GET de filtro/selección sin validación server-side: `data-table.html:7,14,83`, `menu.html:94`, `list.html:72`, `segmented-button.html:4`, `demo-whatsapp.html:16,136,146`, `demo-whatsapp-admin.html:72`, `chips.html:50`.

Conclusión: hoy **ningún form del repo justifica un summary por conteo de campos**. El summary es un patrón de composición listo para Phase G (Resource Editor/Auth/Steps).

---

## 2. CONTRATO 422 + SUMMARY

### 2.1 Cómo funciona el 422 hoy

- Handler de validación (patrón `text_field.go:55-92`, `select.go:72-124`, `toast.go:134-189`): `ParseForm` → valida → si falla: `status = 422` + errores por campo inlineados en el view model.
- HX (`HX-Request: true`): re-render del **fragmento del form** (`ExecuteTemplate` "validation-form"/"select-menu-demo"/"toast-demo-form") + header `X-Loom-Validation: true` (`text_field.go:87-89`).
- No-HX: página completa vía `renderMarkdownPageStatus(..., 422)` (`server.go:155-190`).
- Hook cliente `web/static/app.js:1-9`: `htmx:beforeSwap` swapea el 422 solo si `X-Loom-Validation === "true"`, marcando `shouldSwap = true` / `isError = false`.

### 2.2 Cómo encajaría el summary

- El summary vive **dentro del mismo fragmento que hoy se re-renderiza**, arriba del form (primera línea del `<form>` en `text-field.html:12-17` o un bloque previo dentro del partial del form).
- **El mismo 422 trae el summary + los errores por campo**: el handler recolecta los errores de todos los campos → construye `validationSummaryView{Title, Items[]}` → el partial `validation-form` renderiza `{{template "validation-summary" .Summary}}` antes de los campos.
- **El fragmento HX incluye el summary** (está dentro del `<form>` que swapea `outerHTML`); la página completa no-JS también (422 + `autofocus` de recuperación en el primer campo con error se mantienen).
- El hook `app.js:1-9` **no cambia**; sin header nuevo, sin contrato nuevo (regla `state-patterns-audit.md` §4: "sin inventar contratos"). Validación sigue sin disparar toast (`toast.go:129-133`).

---

## 3. CONTRATO PROPUESTO

### 3.1 Markup (siguiendo inline-alert.html / empty-state.html — partial server-rendered)

```html
{{define "validation-summary"}}<div class="ui-validation-summary" role="alert">
  <h{{.HeadingLevel}} class="ui-validation-summary-title">{{.Title}}</h{{.HeadingLevel}}>
  <ul class="ui-validation-summary-list">
    {{range .Items}}<li class="ui-validation-summary-item"><a href="{{.Href}}">{{.Message}}</a></li>{{end}}
  </ul>
</div>{{end}}
```

- Root `<div>` con `role="alert"` (assertive) **fijo** — el summary ES de errores de validación, no hay tones.
- Título como heading real `<h2>` por defecto, `<h3>` configurable vía `HeadingLevel` (default 2; el `<h1>` de la página queda reservado). Distinto de inline-alert/empty-state que usan `<p>` porque el summary es un landmark navegable propio.
- `<ul>` de `<li>` con `<a href="#{ID}-error">` **reales** (salto nativo al campo, 0 JS). Los mensajes se escapan con el escaping de html/template.

### 3.2 Estructura de datos (Go)

```go
type validationSummaryItem struct {
    Href    string // "#validation-name-error" — id del mensaje de error del campo
    Message string
}

type validationSummaryView struct {
    Title        string                   // p.ej. "2 fields need your attention"
    HeadingLevel int                      // 2 default; 3 para Steps/contextos anidados
    Items        []validationSummaryItem  // orden estable (orden del form)
}
```

- Campo nuevo en el view model del form (`validationFormView`) y en `pageView` si se expone standalone.
- El handler solo construye el summary cuando `len(errors) > 0`.

### 3.3 Tokens (verificados en `web/styles/tokens.css` y `themes/theme-material/theme.css`)

Core:
- `--ui-color-danger` (tokens.css:32; theme :21, :245/:306 dark) ✅
- `--ui-color-danger-container` (tokens.css:51; theme :32, :256/:317 dark) ✅ — misma superficie que `.ui-inline-alert--error` (inline-alert.css:29)
- `--ui-radius-sm` (tokens.css:80) ✅
- `--ui-space-2/3/4` (tokens.css:102-105) ✅
- `--ui-type-title-md` (tokens.css:145 core; theme :65) ✅ — título del summary
- `--ui-type-body-sm` (theme :68) ✅ — mensajes de los items
- Enlace: hereda `color: var(--ui-color-danger)` + `text-decoration: underline` (nativo) — nunca color-only: el enlace real + role="alert" + título textual portan el significado.

Tokens scoped `--ui-validation-summary-*` (declarados en el root, patrón inline-alert.css:12-18):

```css
.ui-validation-summary {
  --ui-validation-summary-padding: var(--ui-space-3) var(--ui-space-4);
  --ui-validation-summary-gap: var(--ui-space-2);
  --ui-validation-summary-radius: var(--ui-radius-sm);
  --ui-validation-summary-bg: var(--ui-color-danger-container);
  --ui-validation-summary-fg: var(--ui-color-danger);
  --ui-validation-summary-title-color: var(--ui-color-danger);
  --ui-validation-summary-item-color: var(--ui-color-danger);
}
```

### 3.4 Variantes

**Una sola variante (error)**: el summary es exclusivamente de errores de validación de formulario. Sin tones, sin `role` condicional (`alert` fijo), sin icono obligatorio. Si en el futuro hiciera falta un resumen informativo, ese es el dominio de Inline alert / Banner — no de este patrón.

### 3.5 CSS

- `web/styles/validation-summary.css` nuevo, `@layer components`, anatomía de `inline-alert.css`/`empty-state.css`: tokens scoped en el root, título con `font: var(--ui-type-title-md)`, lista sin bullets (`list-style: none`, `margin: 0`, `padding: 0`), items con `font: var(--ui-type-body-sm)` y links subrayados, bloque `@media (forced-colors: active)` (borde `CanvasText`, links `LinkText`, texto `Mark` en error — precedente inline-alert.css:54-69). Sin animación → sin bloque reduced-motion.
- `web/styles/app.css`: `@import "./validation-summary.css";` (tras `inline-alert.css` :34, junto a los patrones de estado).
- `web/styles_contract_test.go`: agregar `"styles/validation-summary.css"` a `sourceAppCSS` (:24-58) y un token scoped a `TestComponentSizeTokensDeclaredScoped` (:405-429). El guard `TestNoColorLiteralsInComponents` (`checked < 25`, :742) pasa a 26 archivos.

### 3.6 ¿Reusa `.ui-inline-alert` o es independiente?

**Independiente**, pero con la misma superficie visual que `.ui-inline-alert--error` (danger-container + danger). El audit previo (`inline-alert-audit.md:127`) describe que el summary "se compone SOBRE este inline alert" — es decir, ocupa el mismo slot de sección, pero su anatomía difiere: heading real + `<ul>` de anclas (inline-alert tiene `<p>` title/body sin lista). Anidar `.ui-inline-alert` dentro del summary forzaría roles duplicados (`alert` dentro de `alert`) y estructura ajena. Decisión: primitiva propia con tokens scoped compartiendo los mismos tokens core; ambos coexisten por capas (summary = form-level navegable, inline alert = sección genérica).

---

## 4. INTEGRACIÓN

### 4.1 Formulario candidato

**Ningún form del repo tiene 2+ campos validados hoy** (ver §1.5). Dos opciones, siendo la segunda la recomendada:

1. **Extender el demo de text-field** (`validation-form`) a 2 campos (p.ej. `name` + `email`): es un cambio real y testeado de `text_field.go`/`text-field.html`/`text_field_test.go` — convierte el demo en el primer consumidor del summary. Costo: tocar el contrato 422 que hoy es mono-campo.
2. **Dejar el summary como patrón listo** para Phase G (Resource Editor / Auth / Steps — `roadmap.md:414`, `vocabulary.md:181`): el partial + CSS + tests se entregan como primitiva; ningún demo cambia. Recomendado: **la opción 2**, coherente con `state-patterns-audit.md` §5 orden 4 ("form-level sobre el contrato 422 existente; desbloquea Resource Editor/Auth").

### 4.2 Patrón de uso (cuando haya un form multi-campo)

```text
POST /resource → handler valida TODOS los campos
  └─ errores: []fieldError (campo, mensaje)
       ├─ len > 0 → 422
       │    ├─ por campo: Error = msg (aria-invalid + aria-describedby, text-field.html:5-8)
       │    ├─ summary: Items = [{Href: "#"+id+"-error", Message: msg}...] → {{template "validation-summary"}}
       │    ├─ HX: fragmento "validation-form" (con summary arriba) + X-Loom-Validation: true
       │    └─ no-HX: página completa 422 + autofocus de recuperación (text_field.go:67)
       └─ len == 0 → 200 (o POST + 303) — sin summary, sin toast
```

- Mismo flujo que `text_field.go:55-92` / `select.go:72-124`; `app.js:1-9` intacto; validación nunca toast.

---

## 5. TESTS PROPUESTOS

1. `web/styles_validation_summary_test.go` (nuevo, patrón `styles_inline_alert_test.go`):
   - `TestValidationSummaryPrimitiveCSSMapsTokens` — contratos exactos: `.ui-validation-summary {`, `display: flex`/`flex-direction: column`, `gap: var(--ui-validation-summary-gap)`, `padding`, `border-radius`, `background`, `color`, título (`font: var(--ui-type-title-md)`), lista/items (sin bullets, `font: var(--ui-type-body-sm)`), todos los tokens scoped. Sin `transition:`/`animation:`/`prefers-reduced-motion`.
   - `TestValidationSummaryContractCSSWired` — compilado embebido (`Assets.ReadFile("static/app.css")`) contiene `.ui-validation-summary`, `.ui-validation-summary-title`, `.ui-validation-summary-list`, `.ui-validation-summary-item`, `@media (forced-colors:active)`.
   - `TestValidationSummaryClassVocabularyIsClosed` — cada clase en template ↔ selector en CSS; prohibido `ui-validation-summary-demo`.
   - `TestValidationSummaryUsesCoreTokens` — solo tokens core (`--ui-color-danger`, `--ui-color-danger-container`, `--ui-space-*`, `--ui-radius-sm`, `--ui-type-*`), nunca hex literal (guard del blanket `TestNoColorLiteralsInComponents`).
2. `web/styles_contract_test.go` (modificar): `sourceAppCSS` + `TestComponentSizeTokensDeclaredScoped` (+`--ui-validation-summary-*`).
3. Render (si se integra en Go, patrón `internal/app/text_field_test.go`): `internal/app/validation_summary_test.go` — assertions de `role="alert"`, heading `<h2>`/`<h3>` configurable, `<ul>` con `<a href="#{campo}-error">` reales por item, mensajes escapados, sin summary cuando `len(Items)==0`. Si se extiende el demo: ajustes en `text_field_test.go` y `select_test.go` (no-JS `autofocus` solo en el primer campo con error).

---

## 6. FILES IMPACTADOS (solo read-only)

**Nuevos**:
- `web/styles/validation-summary.css` — primitiva + tokens scoped + forced-colors.
- `web/styles_validation_summary_test.go` — tests de contrato (patrón inline-alert/empty-state).
- `web/templates/validation-summary.html` — partial `{{define "validation-summary"}}`.
- (opcional, si se integra en Go) `internal/app/validation_summary.go` + `internal/app/validation_summary_test.go` — view model + render tests.
- (opcional) `web/content/validation-summary.md` — página docs si se le da ruta.

**Modificados** (si se integra):
- `web/styles/app.css` — `@import "./validation-summary.css";` (+1, junto a patrones de estado).
- `web/static/app.css` — regenerado por `npm run build` (precedente commits eba1c4c/0688020/43c0dac).
- `web/styles_contract_test.go` — `sourceAppCSS` (:24-58) + `TestComponentSizeTokensDeclaredScoped` (:405-429).
- `web/templates/text-field.html` — (solo opción 1 de §4.1) `validation-form` con 2 campos + summary.
- `internal/app/text_field.go` — (solo opción 1) view model multi-campo + construcción del summary.
- `internal/app/validation-form` consumidores — (opción 1) `text_field_test.go`, `select_test.go`.
- Docs (opcional, NO en precedente de los commits): `docs/gelium-ui-vocabulary.md` (Validation summary ✖→✅), `docs/gelium-ui-composition-rules.md` (state matrix §8 Form: "422 inline" → "422 inline + summary"), `docs/gelium-ui-system-roadmap.md` (matriz :414).

**No tocados**: `web/static/app.js` (0 JS — el hook 422 queda igual), `internal/app/{select,toast}.go` y sus templates (error de campo se mantiene; `select.html:89` solo ganaría `id` si un summary futuro lo enlaza), `web/styles/{tokens,inline-alert}.css`, `themes/theme-material/theme.css` (no se añaden tokens core nuevos).

---

## 7. Fuentes de autoridad

`docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :408-415), `docs/gelium-ui-vocabulary.md` (:181 Steps), `docs/gelium-ui-composition-rules.md` (:110-119, :149-162, :164-173), `docs/handoffs/{state-patterns-audit,inline-alert-audit,ux-accessibility-audit}.md`, commits `eba1c4c`/`0688020`/`43c0dac` (patrón de archivos), `internal/app/{text_field,select,toast,server,routes}.go`, `web/templates/{text-field,select,toast,inline-alert,empty-state,skeleton,layout}.html`, `web/styles/{tokens,inline-alert,empty-state,app,text-field,select-menu}.css`, `themes/theme-material/theme.css`, `web/styles_contract_test.go`, `web/styles_inline_alert_test.go`, `web/static/app.js`.
