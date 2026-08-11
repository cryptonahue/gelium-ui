# Gelium UI — Empty State Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón EMPTY STATE en Gelium UI. No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/handoffs/state-patterns-audit.md` (fila 1, gap bloqueante), `docs/handoffs/ux-accessibility-audit.md` (G4, líneas 123/146), `docs/gelium-ui-vocabulary.md` (§3 Empty state, líneas 94-101), `internal/app/data_table.go`, `web/templates/data-table.html`, `web/styles/data-table.css`, `internal/app/demo_whatsapp.go`, `web/templates/demo-whatsapp.html`, `web/styles/demo-whatsapp.css`, `web/styles/tokens.css`, `web/styles/app.css`, `web/styles_contract_test.go`, `internal/app/data_table_test.go`.

---

## 1. EMPTY ACTUAL

### 1.1 Data table (server-rendered, GET-only, sin JS)

**Contrato GET**: form `web/templates/data-table.html:7-12` — `method="get" action="/components/data-table#data-table-demo"` + `hx-get="/components/data-table" hx-target="#data-table-panel" hx-swap="outerHTML"`. Filtro: `<input type="search" name="q">` (línea 10). El handler `dataTableDocs` (`internal/app/data_table.go:116-136`) lee `r.URL.Query().Get("q")` (línea 118) y bifurca por header `HX-Request` (líneas 120-130): HTMX devuelve solo el fragmento `data-table-panel`; no-JS devuelve página completa.

**Detección de vacío**: `newDataTableDemo` filtra `dataTableItems` por substring case-insensitive sobre Name O Status (`data_table.go:159-165`); `total := len(items)` (línea 175); `Rows` se construye del slice `items[start:end]` (líneas 199-208) → **vacío = `Total == 0`** (equivalente a `len(Rows) == 0`).

**Estado vacío renderizado hoy** (evidencia con `?q=zzz`):
- **Caption**: `data_table.go:239` — `fmt.Sprintf("%d rows · page %d of %d", total, pageNum, totalPages)` → `0 rows · page 1 of 1`. Render: `data-table.html:37` (`<caption class="ui-data-table-caption">`). Es un dato, no una guía.
- **tbody**: `data-table.html:53-66` — `{{range .Rows}}` sobre slice vacío → `<tbody>` vacío. Sin mensaje, sin CTA.
- **Paginación**: `data-table.html:68-72` — con 0 filas: `HasPrev=false` y `HasNext=false` → spans `--disabled` con `aria-disabled="true"` "Previous"/"Next"; `PageLinks` = [1] con `Current` → `<span ... aria-current="page">1</span>`. Se ve "Previous 1 Next" muerto.
- **Selección**: checkbox select-all `data-table.html:42` — `{{if eq .SelectedCount .Total}} checked{{end}}`; filas `data-table.html:54-64`; `SelectedCount` solo se computa si `len(selection) > 0` (`data_table.go:210-221`).

### 1.2 Demo WhatsApp (referencia visual del único empty ad-hoc)

**Markup** (`web/templates/demo-whatsapp.html:51-53`), dentro del `<ul class="demo-wa-conversations">`:
```html
{{if eq (len .Conversations) 0}}
<li class="demo-wa-empty">Sin resultados para tu búsqueda.</li>
{{end}}
```

**Estilos** (`web/styles/demo-whatsapp.css:199-203`):
```css
.demo-wa-empty {
  padding: 1rem;
  color: var(--ui-color-fg-muted);
  font: var(--ui-type-body-sm);
}
```
Single line, muted, sin icono, sin título, sin CTA. Handler: `demo_whatsapp.go:493-539` — `q` (línea 495), filtro por ContactName/Phone (503-508), view `Conversations` (534).

**Relacionado (NO es empty, pero es la otra cara visual)**: `demo-wa-placeholder` (`demo-whatsapp.html:153-158`, CSS `demo-whatsapp.css:425-434`) — placeholder de selección (no hay chat activo): grid centrado `place-content: center`, `text-align: center`, título `--ui-type-title-lg` color fg + body `--ui-type-body-sm` muted. Referencia visual de centrado para el empty de página.

---

## 2. BUG SELECT-ALL (G4, `ux-accessibility-audit.md:123,146`)

**Causa raíz**: `data-table.html:42` — `eq .SelectedCount .Total` con 0 filas y sin selección enviada: `0 == 0` → true → el checkbox "Select all rows" se renderiza `checked` con cero filas. Además, si el usuario envía igual, `selection=all` produce "All rows selected." (`data_table.go:338-339`) con 0 filas — mensaje engañoso.

**Fix propuesto** (server-driven, 0 JS):
1. Guardar en `Total`, no en la igualdad: el input select-all solo se renderiza/activa si `gt .Total 0`, y `checked` requiere `and (gt .Total 0) (eq .SelectedCount .Total)`.
2. `disabled` (o ausente) cuando `Total == 0`; el `<th>` checkbox queda vacío (el empty row con `colspan` cubre el ancho).
3. El botón "Submit selection" (`data-table.html:20`) también se deshabilita/omite con 0 filas.
4. Opcional: guardar el notice en `data-table.html:21` con `gt .Total 0` (evitar "All rows selected." con 0 filas).

---

## 3. CONTRATO PROPUESTO — `.ui-empty-state`

### Markup (partial `web/templates/empty-state.html`)
```html
{{define "empty-state"}}
<div class="ui-empty-state{{if .Compact}} ui-empty-state--compact{{end}}" role="status">
  {{if .Icon}}<span class="ui-empty-state-icon" aria-hidden="true">{{.Icon}}</span>{{end}}
  <p class="ui-empty-state-title">{{.Title}}</p>
  <p class="ui-empty-state-body">{{.Body}}</p>
  {{if .CTA}}<a class="ui-button" href="{{.CTAHref}}">{{.CTALabel}}</a>{{end}}
</div>
{{end}}
```
- **Título**: `<p class="ui-empty-state-title">` con `font: var(--ui-type-title-md)` y `color: var(--ui-color-fg)` — NO `<h2>/<h3>` fijo: dentro de un `<td>` de tabla el heading es semánticamente incómodo y el nivel depende del contexto (la demo ya tiene `h1` + `h3` de sección). Si se quiere heading real (feed/página), el partial acepta un nivel configurable; el default es `<p>` con peso visual de título.
- **Body**: `<p class="ui-empty-state-body">`, muted.
- **CTA**: SOLO elemento real — `<a href>` para GET/navegación (ej. "Clear filter" → `?sort=name&dir=asc`), `<button>` para acciones. Nunca div/span como control (anti-regla de `composition-rules.md`). Reusa `.ui-button`.
- **Icono**: opcional, decorativo, `aria-hidden="true"`.

### ARIA
- `role="status"` en `.ui-empty-state` → implícito `aria-live="polite"`, announce al insertarse. El swap HTMX `outerHTML` del panel (`data-table.html:8,48,70,71`) inserta el nodo → el SR lo anuncia. Precedente: notice de selección `data-table.html:21` y `chips.html:63` (`role="status"`).
- NO agregar `aria-live` al contenedor del panel: el swap reemplaza todo el panel y un aria-live en él re-anunciaría todo. `role="status"` en el nodo insertado es suficiente y evita doble anuncio. El patrón `aria-relevant="additions text"` del toast (`toast.html:10`) es para la región del toast, no se replica aquí.
- Vacío→no-vacío: el nodo se remueve (remoción no se anuncia por polite). Si se quisiera anunciar la recuperación, el patrón existente es un notice `role="status"` server-rendered como `SelectionNotice` (ej. "N rows shown"). Opcional, NO inventar en Phase D.
- Mensaje nunca color-only.

### Tokens (todos existen en `web/styles/tokens.css` / theme)
- Color: `--ui-color-fg` (título), `--ui-color-fg-muted` (body), `--ui-color-primary` (CTA/icono opcional).
- Spacing: `--ui-space-2/3/4/6` (padding/gap).
- Size: `--ui-size-icon` o `--ui-size-icon-sm` (glyph).
- Type: `--ui-type-title-md` (título), `--ui-type-body-sm` (body); CTA usa el type del propio `.ui-button`.
- **Tokens scoped** (convención repo, `TestComponentSizeTokensDeclaredScoped`): declarar en `.ui-empty-state`:
  `--ui-empty-state-padding: var(--ui-space-6); --ui-empty-state-gap: var(--ui-space-2); --ui-empty-state-icon-size: var(--ui-size-icon); --ui-empty-state-title-color: var(--ui-color-fg); --ui-empty-state-body-color: var(--ui-color-fg-muted);`

### Variantes
- **Una sola variante base centrada** (default): `display: grid; place-content: center; text-align: center;` — para feed/search-results a nivel página.
- **`--compact`** (start-aligned, inline): texto alineado a inicio, padding menor — para el caso dentro de tabla/listas (row height de tabla es 52px; el empty no debe inflar filas).
- Icono on/off: elemento opcional, NO clase de variante.

### CSS
- Archivo nuevo `web/styles/empty-state.css` con `@layer components { .ui-empty-state {...} }`, bloque `@media (forced-colors: active)` (título/body → CanvasText; CTA → LinkText/Highlight, patrón `data-table.css:340-359`). Sin animación → sin bloque reduced-motion propio (si se agrega, convención `transition: none`).
- `@import "./empty-state.css";` en `web/styles/app.css` (después de data-table.css, línea 32) **y** en la lista `sourceAppCSS` de `web/styles_contract_test.go:24-58` (sync obligatoria, `state-patterns-audit.md:102`).

---

## 4. INTEGRACION TABLA

**tbody** (`data-table.html:53-66`):
```html
<tbody>
  {{if .Rows}}
    {{range .Rows}}<tr class="ui-data-table-row">…{{end}}
  {{else}}
    <tr>
      <td colspan="{{.Colspan}}" class="ui-data-table-cell">
        {{template "empty-state" .EmptyState}}
      </td>
    </tr>
  {{end}}
</tbody>
```

- **colspan**: 4 = `len(Columns)` (3: name/status/date) + 1 columna checkbox. No existe `FuncMap` con `add` en el repo (grep sin resultados) → opciones: (a) hardcodear `colspan="4"` (vocabulario cerrado, `dataTableSortKeys` es fijo); (b) campo de view `Colspan int` computado en `newDataTableDemo` (`len(dataTableColumns(...)) + 1`), server-driven y auto-mantenible. **Recomendado (b)**; un test de contrato fija `colspan="4"`.
- **Caption**: mantener `0 rows · page 1 of 1` (`data_table.go:239`) — es el dato del `<caption>` (contrato de tabla); el empty row agrega la guía. La fila `role="status"` anuncia el mensaje.
- **Select-all**: ver §2. Con 0 filas el `<th>` checkbox queda sin input o con input `disabled`; el empty row con colspan preserva la grilla.
- **HTMX**: cero contrato nuevo — el empty row vive dentro del fragmento `data-table-panel`, el swap `outerHTML` existente lo trae.
- **Demo WhatsApp**: el `demo-wa-empty` queda como referencia visual (es demo, clase sin prefijo ui). Migrar el demo al componente es opcional y fuera del alcance Phase D.

---

## 5. TESTS PROPUESTOS

**Nuevo `web/styles_empty_state_test.go`**:
1. `TestEmptyStatePrimitiveCSSMapsTokens` — `.ui-empty-state {` + consumo de tokens scoped (padding/gap/icon-size/title-color/body-color) y type tokens.
2. `TestEmptyStateContractCSSWired` — presencia en el CSS compilado `static/app.css` (espejo de `TestEmbeddedCompiledCSSIncludesDataTableContracts`) + forced-colors.
3. `TestEmptyStateClassVocabularyIsClosed` — clases del partial (`ui-empty-state`, `--compact`, `--title`, `--body`, `--icon`) existen en template y CSS; prefijo ui reservado; sin `ui-empty-state-demo`.

**`web/styles_contract_test.go`**:
4. Agregar `"styles/empty-state.css"` a `sourceAppCSS` (líneas 24-58).
5. (Opcional) entrada en `TestComponentSizeTokensDeclaredScoped` si se declara token de tamaño scoped. `TestNoColorLiteralsInComponents` y `TestSpace/SizeTokensConsumedByComponents` cubren el archivo automáticamente.

**`internal/app/data_table_test.go`**:
6. `TestDataTableEmptyStateRendersMessageAndCTA` — `?q=zzz` → 200, caption `0 rows · page 1 of 1`, `.ui-empty-state`, `colspan="4"`, CTA `<a class="ui-button"`, ausencia de `<tr class="ui-data-table-row">`.
7. `TestDataTableEmptyStateSelectAllDisabledOrAbsent` — `?q=zzz` → el input `selection=all` NO está `checked`; está `disabled` o ausente; "Submit selection" deshabilitado/ausente.
8. `TestDataTableEmptyStateHXFragmentIncludesEmptyRow` — `HX-Request: true` + `?q=zzz` → fragmento contiene `.ui-empty-state` en tbody, sin `<html`.
9. `TestDataTableEmptyStateEscapesQuery` — `?q=<script>` → mensaje/CTA del empty escapan (extensión de `TestDataTableDocsRouteEscapesFilterQuery`).
10. `TestDataTableEmptyStateAllSelectionCountsZero` — `?q=zzz&selection=all` → sin "All rows selected." (notice guardado o vacío), SelectedCount=0.
11. Actualizar `TestDataTableDemoClassVocabularyIsClosed` si aparecen clases demo nuevas.

---

## 6. FILES IMPACTADOS (read-only — NO modificados por este audit)

| Archivo | Cambio |
|---|---|
| `web/styles/empty-state.css` | **NUEVO** — primitiva `.ui-empty-state` (+ `--compact`) |
| `web/styles/app.css` | `@import "./empty-state.css";` |
| `web/styles_contract_test.go` | lista `sourceAppCSS` + (opcional) mapa de tokens scoped |
| `web/styles_empty_state_test.go` | **NUEVO** — tests de contrato CSS |
| `web/templates/empty-state.html` | **NUEVO** — partial `{{define "empty-state"}}` |
| `web/templates/data-table.html` | tbody condicional (líneas 53-66), select-all fix (línea 42), submit guard (línea 20), notice guard (línea 21) |
| `internal/app/data_table.go` | campo `Colspan` y view model del empty state (título/body/CTA) en `newDataTableDemo` (alrededor de :239) |
| `internal/app/data_table_test.go` | tests de comportamiento 6-11 |
| `docs/gelium-ui-vocabulary.md` | Empty state ◐→✅ (DoD, `state-patterns-audit.md:102`) |
| `docs/gelium-ui-composition-rules.md` | state matrix §4.8/§5 |
| `COMPONENT-ROADMAP.md` | registro del componente |
| Opcional (fuera de Phase D): `web/templates/demo-whatsapp.html` + `web/styles/demo-whatsapp.css` | migración del demo al componente (referencia visual, no bloqueante) |
| Opcional (no requerido): `internal/app/routes.go`, `internal/app/server.go` | solo si se agrega una página de docs standalone del empty |

**DoD completo** (`state-patterns-audit.md:102`): `styles_empty_state_test.go` + partial + entrada en `sourceAppCSS` + actualización vocabulario/matrix. Ruta standalone NO requerida para la integración con tabla.

---

## 7. Fuentes de autoridad usadas

`internal/app/data_table.go` (116-136, 142-252, 326-345), `web/templates/data-table.html` (7-12, 20-21, 34-74), `web/styles/data-table.css` (1-360), `internal/app/demo_whatsapp.go` (493-539), `web/templates/demo-whatsapp.html` (30-55, 76, 153-158), `web/styles/demo-whatsapp.css` (199-203, 425-434), `web/styles/tokens.css` (19-146), `web/styles/app.css` (1-34), `web/styles_contract_test.go` (13-71, 402-423, 704-739), `internal/app/data_table_test.go`, `web/styles_data_table_test.go`, `docs/handoffs/{state-patterns-audit,ux-accessibility-audit}.md`, `docs/gelium-ui-vocabulary.md` (94-101), `themes/theme-material/theme.css` (61-70), `web/styles/toast.css` (70-80), `web/templates/toast.html` (10).
