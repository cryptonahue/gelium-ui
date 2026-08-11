# Gelium UI — Error State Audit (Phase D, pattern 7, handoff)

> **Alcance**: inventario read-only del patrón ERROR STATE en Gelium UI. No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/handoffs/state-patterns-audit.md` (fila 7 y §4), `docs/gelium-ui-system-roadmap.md` (Phase D y E), `docs/gelium-ui-vocabulary.md` (§3), `internal/app/{server,routes,docs,data_table}.go`, `web/templates/{layout,banner,empty-state,callout}.html`, `web/styles/{tokens,app,banner,empty-state,callout}.css`, `web/static/app.js`, `web/styles_contract_test.go`, `web/styles_{banner,callout,empty_state}_test.go`, `internal/app/server_test.go`.

---

## 1. ESTADO ACTUAL — manejo de errores hoy

**No existe ningún handler de 404/500 personalizado.** El sistema es el `http.NewServeMux` de Go 1.22+ con patrones de método.

### Rutas registradas (`internal/app/server.go:105-128`)
- `GET /healthz` (:106), `GET /{$}` home (:107), `GET /docs` (:108), rutas de componentes (:109-114), POSTs de ejemplos (:115-119), demos WhatsApp (:120-125), `GET /static/{name}` (:126).
- **No hay catch-all 404**: una ruta inexistente cae al `NotFoundHandler` default de `net/http`.

### Qué pasa hoy con una ruta inexistente
- **404 default de net/http**: el `ServeMux` responde texto plano `404 page not found\n` sin layout, sin CSS, sin `<h1>`, sin retry. Confirmado por test: `internal/app/text_field_test.go:333-337` (`GET /examples/text-field/missing` → 404) y `dialog_test.go:78` (`/components/dialog/missing` → 404). Solo validan el **status**, no el body.
- `http.NotFound` explícito solo para assets estáticos: `server.go:145,150` (`staticAsset` con nombre desconocido o archivo inexistente).

### Errores 500 actuales (texto plano, sin estado de UI)
- `server.go:177` — `fs.ReadFile` del contenido markdown falla → `http.Error(w, "documentation unavailable", 500)`.
- `server.go:193` — falla la conversión goldmark → `http.Error(w, "documentation unavailable", 500)`.
- `server.go:202` — falla `templates.ExecuteTemplate("layout")` → `http.Error(w, "documentation unavailable", 500)`.
- Idem fragmentos: `data_table.go:140,390,401,408`, `select.go:118`, `toast.go:157,183`, `text_field.go:83` (todos `http.Error` texto plano 500/400).
- **Nota crítica**: el 500 de `server.go:202` es el *template exec* del propio layout — NO puede renderizar la página de error vía layout (es el mismo template que falló). Ese fallback final debe seguir siendo mínimo/plano; los 500 de :177 y :193 sí pueden usar el layout (el template ya está parseado y funcionando).

### El pipeline que el patrón debe aprovechar
- `renderMarkdownPageStatus(w, data, contentPath, status)` (`server.go:174-181`) y `renderMarkdownStatus(..., status)` (`server.go:190-208`) ya aceptan un **status HTTP real** y lo escriben con `w.WriteHeader(status)` (:207). Es el gancho canónico: el estado de error se re-renderiza server-side con el status real, exactamente como indica `state-patterns-audit.md:72`.
- El layout (`web/templates/layout.html`) no tiene `<h1>` propio: las páginas obtienen su h1 del markdown (`# ...` dentro de `.prose`, layout.html:18). Un `.ui-error-state` con `<h1>` en un slot del layout no compite con ningún otro h1.

### Referencia: patrón Banner (el análogo más cercano — patrón de página, ya entregado)
- `bannerView` struct: `server.go:32-48`; slot `Banner *bannerView` en `pageView`: `server.go:55`; conditional en layout: `layout.html:16` (`{{if .Banner}}{{template "banner" .Banner}}{{end}}`).
- **Banner está deliberadamente sin wirear** (`server.go:32-37`): "no handler wires it yet — Phase G will inject it". Solo se entregó la primitiva + slot + tests (`server_test.go:229-342`: `renderBanner` helper, `TestLayoutOmitsBannerWhenNil`, `TestLayoutRendersBannerSlotBetweenHeaderAndMain`, `TestBannerRoleIsDerivedFromTone`, `TestBannerMarkupContracts`).
- Convención de contrato CSS: `styles_<pattern>_test.go` (3 tests: tokens mapeados / CSS compilado wireado / vocabulario cerrado) + entrada en `sourceAppCSS` (`styles_contract_test.go:22-77`) + `@import` en `app.css:33-39`.

---

## 2. OPCIONES DE SERVICIO

### (a) Página standalone de error (404/500 a nivel sitio)
Handler que captura rutas inexistentes o recursos fallidos y responde la página completa con el layout + `.ui-error-state`, status HTTP real (404/500), sin contenido markdown. Requiere un slot `Error *errorStateView` en `pageView` + conditional en el layout (patrón Banner).

### (b) Estado dentro de una página existente (recurso que falla)
El handler de un recurso detecta el error y re-renderiza **esa misma página** con `.ui-error-state` reemplazando el contenido + status real (404/500/503). Es exactamente el flujo de `renderMarkdownPageStatus`: hoy los handlers lo derivan del markdown; con el patrón lo derivan del view model.

### Recomendación: (b) es la integración natural; el patrón se sirve como slot de página (a+b)
- La arquitectura **ya tiene el gancho de status real** en `renderMarkdownPageStatus` (`server.go:174`); (b) solo cambia *qué* se renderiza en el body (slot en lugar de markdown), reusando el mismo mecanismo.
- El estado de error canónico (task: "página/recurso con error") es **el mismo primitivo en ambos casos**: página completa para 404/500 de sitio y fragmento/recurso para errores de handler. No son dos componentes, es un slot de página reutilizable (igual que Banner).
- El transport (500/network en HTMX) queda **fuera de Phase D**: `app.js:1-9` solo maneja el swap 422; el transporte es Phase E (`roadmap.md:419`, `state-patterns-audit.md:40,72`). No se toca `app.js`.
- Por lo tanto: entregar `.ui-error-state` como **primitiva + slot `Error` en `pageView`** (patrón Banner exacto). El wireing del handler se decide en §4.

---

## 3. CONTRATO PROPUESTO — `.ui-error-state`

### Markup (partial `web/templates/error-state.html`)
Siguiendo `empty-state.html` / `banner.html` (partial de una línea, interpolación server-driven):

```html
{{define "error-state"}}<div class="ui-error-state" role="alert">
  {{if .StatusCode}}<p class="ui-error-state-code" aria-hidden="true">{{.StatusCode}}</p>{{end}}
  <h1 class="ui-error-state-title">{{.Title}}</h1>
  <p class="ui-error-state-body">{{.Body}}</p>
  {{if .Retry}}<a class="ui-button" href="{{.RetryHref}}">{{.RetryLabel}}</a>{{end}}
</div>{{end}}
```

- **Status code visible**: `<p class="ui-error-state-code">` con el número grande (font `--ui-type-display-lg`), **`aria-hidden="true"`** — es énfasis decorativo del h1; el significado lo carga el heading ("Page not found"), y `role="alert"` ya anuncia el h1. Precedente: el icono de empty-state/banner es `aria-hidden` (`empty-state.html:2`, `banner.html:2`).
- **Heading**: `<h1>` **único de página** (canónico, `state-patterns-audit.md:40`). El layout no tiene otro h1, así que no hay conflicto. Para uso inline/fragmento en Phase G el nivel se revisa (igual que empty-state discutió `<p>` vs heading), pero Phase D entrega página-level con h1.
- **Body**: `<p class="ui-error-state-body">` muted.
- **Retry**: SOLO elemento real — `<a class="ui-button" href>` para GET (retry del recurso: recargar URL / volver al home). `<button>` dentro de form POST queda para acciones de Phase G; la canónica es GET. `Retry` booleano omite el link (un 404 puede no tener retry útil y apuntar a "/").

### Variantes: **genérico con `StatusCode + Title + Body + RetryHref`** — NO variantes por status
- La diferencia canónica 404/500 es **datos** (status HTTP), no un tono visual distinto: ambos son "el recurso no se pudo entregar" y comparten el tono danger. `StatusCode` es un campo del view model (`errorStateView{StatusCode, Title, Body, Retry, RetryHref, RetryLabel}`); el handler elige el número y el título por status.
- Los patrones Phase D que sí usan variantes lo hacen por **tono con diferencia visual/semántica real** (banner `--error/--warning/--success/--info`, callout `--info/--tip`). Error state tiene UN solo tono (error/danger): no hay `--404`/`--500` que pintar distinto.
- Evita CSS muerto: un `.ui-error-state--404` y `--500` idénticos son ruido de contrato (viola la filosofía "solo tokens con consumidores reales", `tokens.css:38-44`).
- Los títulos por status viven en el handler (helper tipo `errorStateView(404)` → "Page not found" / `errorStateView(500)` → "Something went wrong"), no en el template ni en CSS.

### ARIA: `role="alert"` — **sí, apropiado para página/recurso**
- Decisión: error de página es un fallo que requiere acción (retry/volver) → anuncio **assertive** con `role="alert"`. Difiere de empty-state (`role="status"`, `empty-state.html:1`) porque el vacío **no es un fallo**, es ausencia de datos: anuncio polite, sin urgencia.
- Precedente en el repo: banner usa `role="alert"` para el tone error y `role="status"` para el resto (`banner.html:1`, testeado en `server_test.go:296-307`); toast error usa `role="alert"` (`app.js:45`). La regla repo: error → alert, no-error → status.
- Matiz: en carga inicial server-rendered el alert no se anuncia como tal (el SR anuncia el título de página); su valor real aparece en Phase G cuando un swap HTMX inserta el nodo de error. Consistente con la justificación de empty-state (`empty-state-audit.md:77-79`). Sin `aria-live` extra en el contenedor (evita doble anuncio, misma razón).
- Sin auto-dismiss, sin foco automático: es estado persistente de página.

### Tokens (todos existen — `web/styles/tokens.css`)
| Rol | Token | Línea |
|---|---|---|
| Code color (acento danger) | `--ui-color-danger` | `tokens.css:32` |
| Title | `--ui-color-fg` | `tokens.css:24` |
| Body | `--ui-color-fg-muted` | `tokens.css:25` |
| Spacing | `--ui-space-2/4/6/8` | `tokens.css:102-107` |
| Code font (big) | `--ui-type-display-lg` (core) / `--ui-type-display-sm` | `tokens.css:144`, theme `theme.css:61` |
| Title font | `--ui-type-title-lg` | theme `theme.css:64` |
| Body font | `--ui-type-body-sm` | theme (tiene `body-sm/md/lg`) |
| Scrim/fondo opcional | `--ui-color-surface-container` (como callout/banner) | `tokens.css:23` |

**Nota danger vs error**: `--ui-color-danger` es el canónico; `--ui-color-error` es alias temporal de compatibilidad (`tokens.css:38`, "never the other way around"). El patrón usa `--ui-color-danger*`, nunca `--ui-color-error`. Banner ya sienta precedente (`banner.css:29`).

**Tokens scoped `--ui-error-state-*`** (declarados en la raíz `.ui-error-state`, convención `TestComponentSizeTokensDeclaredScoped`):
```
--ui-error-state-padding: var(--ui-space-8);
--ui-error-state-gap: var(--ui-space-2);
--ui-error-state-code-color: var(--ui-color-danger);
--ui-error-state-title-color: var(--ui-color-fg);
--ui-error-state-body-color: var(--ui-color-fg-muted);
```
- **Sin token de size scoped** (no hay icono en el markup, el ancla visual es el número): no requiere entrada en `TestComponentSizeTokensDeclaredScoped`. `TestSpaceTokensConsumedByComponents` cubre el archivo automáticamente por los `var(--ui-space-*)`.
- Sin literales hex (guarda `TestNoColorLiteralsInComponents`); sin `color-mix` necesario (danger directo sobre canvas, contraste alto en display).

### CSS (`web/styles/error-state.css`, nuevo)
```css
@layer components {
  .ui-error-state {
    --ui-error-state-padding: var(--ui-space-8);
    --ui-error-state-gap: var(--ui-space-2);
    --ui-error-state-code-color: var(--ui-color-danger);
    --ui-error-state-title-color: var(--ui-color-fg);
    --ui-error-state-body-color: var(--ui-color-fg-muted);

    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: var(--ui-error-state-gap);
    padding: var(--ui-error-state-padding);
  }

  .ui-error-state-code {
    margin: 0;
    font: var(--ui-type-display-lg);
    color: var(--ui-error-state-code-color);
  }

  .ui-error-state-title {
    margin: 0;
    font: var(--ui-type-title-lg);
    color: var(--ui-error-state-title-color);
  }

  .ui-error-state-body {
    margin: 0;
    font: var(--ui-type-body-sm);
    color: var(--ui-error-state-body-color);
  }
}

@media (forced-colors: active) {
  .ui-error-state-code,
  .ui-error-state-title,
  .ui-error-state-body { color: CanvasText; forced-color-adjust: auto; }
  .ui-error-state a.ui-button { color: LinkText; }
  .ui-error-state a.ui-button:focus-visible { outline-color: Highlight; }
}
```
- Centrado en página (misma geometría que empty-state, `empty-state.css:16-23`). **Sin animación → sin bloque reduced-motion propio** (convención empty-state/callout; si se agrega transición, es `transition: none`).
- **Wiring obligatorio**: `@import "./error-state.css";` en `web/styles/app.css` (tras la línea 39, junto a los demás state patterns) **y** entrada `"styles/error-state.css"` en `sourceAppCSS` (`styles_contract_test.go:24-64`) — sync obligatoria (`state-patterns-audit.md:102`).

---

## 4. INTEGRACION

### Slot de página (patrón Banner, imprescindible para el DoD)
- `internal/app/server.go`: struct `errorStateView {StatusCode int; Title string; Body string; Retry bool; RetryHref string; RetryLabel string}` + campo `Error *errorStateView` en `pageView` (junto a `Banner`, `server.go:55`).
- `web/templates/layout.html`: envolver el contenido:
  ```html
  {{if .Error}}{{template "error-state" .Error}}{{else}}<article class="prose">…{{end}}
  ```
  (el else abarca `.prose` y los previews; una página de error no lleva previews de componentes).

### Wireing (handler detecta error → render con status real)
1. **404 catch-all**: `mux.HandleFunc("/{path...}", s.notFound)` en `server.go` (`New()`, tras las rutas específicas — el `ServeMux` prioriza el patrón más específico, `GET /static/{name}` y `GET /{$}` ganan sobre el catch-all). El handler construye `errorStateView{StatusCode: 404, Title: "Page not found", Body: …, Retry: true, RetryHref: "/", RetryLabel: "Back to home"}` y renderiza el layout **sin markdown**, con `WriteHeader(404)`. Hoy es el default `404 page not found` plano (`text_field_test.go:333-337`).
2. **500 de recurso**: reemplazar los `http.Error` de `server.go:177,193` por el mismo render con `StatusCode: 500`. **El de :202 (template exec) NO puede usar el layout** — mantener un fallback mínimo plano (el template que falló es el layout mismo).
3. **Errores de handler/fragmento** (data_table/select/toast/text_field): el patrón de bifurcación `HX-Request` (`data_table.go:137-147`) re-renderiza el fragmento con `.ui-error-state`/`role="alert"` — Phase G.
4. **Transport (500/network HTMX)**: NO en Phase D (`state-patterns-audit.md:72`, `roadmap.md:419`). `app.js:1-9` intacto.

### ¿Wirear ahora o primitiva lista para Phase G?
- **Opción A — primitiva + slot, sin wireing** (espejo exacto de Banner): 0 handlers; Phase G registra el 404/500. Más consistente con cómo se entregó Banner (`server.go:32-37`), cero riesgo de mux.
- **Opción B — primitiva + slot + 404 catch-all mínimo ahora**: el patrón de error es el ÚNICO de Phase D cuyo atributo canónico es el **status HTTP real**, y hoy una ruta inexistente devuelve texto plano sin UI (el único gap de Phase D visible inmediatamente, sin esperar a Phase G). El wireing es barato: un catch-all + reemplazo de dos `http.Error` en :177/:193 (el :202 queda plano por definición) + tests de status/markup.

**Recomendación: Opción B** (primitiva + slot + catch-all 404 + swap de los 500 :177/:193). Razones: (1) solo así el patrón es verificable de verdad (test puede afirmar `404` + `.ui-error-state` + `h1` — un primitivo sin status real no demuestra la diferencia canónica); (2) Banner se difirió porque su trigger es una condición de app (sesión/mantenimiento) que no existe hasta Phase G, mientras que el trigger del error es **inmediato** (cualquier ruta desconocida hoy); (3) no contradice Phase E — el transporte HTMX sigue fuera. Si el equipo prefiere disciplina estricta tipo Banner, la Opción A es la alternativa segura; se entrega el slot listo y Phase G solo agrega handlers.

---

## 5. TESTS PROPUESTOS

**`web/styles_error_state_test.go`** (espejo de `styles_callout_test.go` / `styles_empty_state_test.go`):
1. `TestErrorStatePrimitiveCSSMapsTokens` — `.ui-error-state {` + `display: flex; flex-direction: column;` + tokens scoped (padding/gap/code-color/title-color/body-color) + type tokens (`--ui-type-display-lg`, `--ui-type-title-lg`, `--ui-type-body-sm`). Sin `prefers-reduced-motion`, sin `transition:`/`animation:`.
2. `TestErrorStateContractCSSWired` — presencia en el CSS compilado `static/app.css` (`.ui-error-state`, `.ui-error-state-code`, `.ui-error-state-title`, `.ui-error-state-body`, `@media (forced-colors:active)`).
3. `TestErrorStateClassVocabularyIsClosed` — clases del partial existen en template y CSS; interpolación del status code; prefijo ui reservado; sin `ui-error-state-demo`.
4. `TestErrorStateUsesCoreDangerTokens` — `code-color: var(--ui-color-danger)`, `title-color: var(--ui-color-fg)`, `body-color: var(--ui-color-fg-muted)`; sin literales hex (refuerza `TestNoColorLiteralsInComponents`).

**`web/styles_contract_test.go`**:
5. Entrada `"styles/error-state.css"` en `sourceAppCSS` (líneas 24-64). Sin entrada en `TestComponentSizeTokensDeclaredScoped` (no declara token de tamaño).

**`internal/app/server_test.go`** (espejo de los tests de Banner, `server_test.go:229-342`):
6. `TestErrorStateRoleIsAlert` — render del partial → `role="alert"`.
7. `TestErrorStateMarkupContracts` — `StatusCode` visible, `<h1>` único, body, retry `<a class="ui-button" href>`, `Retry=false` omite el link.
8. `TestLayoutOmitsErrorStateWhenNil` + `TestLayoutRendersErrorStateSlot` — slot opcional, reemplaza `.prose` cuando está set.
9. Si se wirea (Opción B): `TestUnknownRouteRendersErrorStatePage` — `GET /does-not-exist` → status 404 + `.ui-error-state` + `Page not found` + ausencia del body plano `404 page not found`.

---

## 6. FILES IMPACTADOS (read-only — NO modificados por este audit)

| Archivo | Cambio |
|---|---|
| `web/styles/error-state.css` | **NUEVO** — primitiva `.ui-error-state` |
| `web/templates/error-state.html` | **NUEVO** — partial `{{define "error-state"}}` |
| `web/styles_error_state_test.go` | **NUEVO** — tests de contrato CSS (5.1-5.4) |
| `web/styles/app.css` | `@import "./error-state.css";` |
| `web/styles_contract_test.go` | lista `sourceAppCSS` (:24-64) |
| `internal/app/server.go` | struct `errorStateView` + campo `Error` en `pageView` (junto a `Banner`, :55) + (Opción B) catch-all `/{path...}` + handler `notFound` + reemplazo de `http.Error` :177/:193 |
| `web/templates/layout.html` | conditional `{{if .Error}}…{{else}}` sobre `.prose`/previews |
| `internal/app/server_test.go` | tests 5.6-5.9 (espejo Banner) |
| `docs/gelium-ui-vocabulary.md` | Error state ✖→✅ (DoD Phase D, `roadmap.md:157-172`) |
| `docs/gelium-ui-composition-rules.md` | state matrix §4.8/§5 |
| `COMPONENT-ROADMAP.md` | registro del componente |
| `docs/handoffs/state-patterns-audit.md` | fila 7 ✖→✅ |
| **No** | `internal/app/routes.go` (catch-all va en `New()`, no en `componentRoutes`), `internal/app/docs.go`, `web/static/app.js` (transporte = Phase E), `web/templates/*` de demos |

**DoD completo** (`state-patterns-audit.md:102`): `styles_error_state_test.go` + partial + slot `Error` en `pageView` + entrada `sourceAppCSS` + vocabulario/matrix actualizados. Sin página de docs standalone (ningún state pattern de Phase D tiene ruta propia; se demuestran por tests y por las pantallas de Phase G).

---

## 7. Fuentes de autoridad usadas

`internal/app/server.go` (32-48, 55, 105-128, 174-208), `internal/app/routes.go` (16-47), `internal/app/docs.go` (83-97), `internal/app/data_table.go` (137-147, 240-253), `internal/app/server_test.go` (229-342), `internal/app/text_field_test.go` (326-338), `internal/app/dialog_test.go` (78), `web/templates/{layout,banner,empty-state,callout}.html`, `web/styles/{tokens,app,banner,empty-state,callout}.css`, `web/static/app.js` (1-9, 45), `web/styles_contract_test.go` (22-77, 408-433, 715-750), `web/styles_{banner,callout,empty_state}_test.go`, `themes/theme-material/theme.css` (61-69, 14-32, 238-256, 299-317), `docs/handoffs/{state-patterns-audit,empty-state-audit}.md`, `docs/gelium-ui-system-roadmap.md` (157-172, 413, 419), `docs/gelium-ui-vocabulary.md` (§3).
