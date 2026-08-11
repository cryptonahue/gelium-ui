# Gelium UI — Skeleton Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón SKELETON en Gelium UI (repo físico `loom-ui`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/handoffs/state-patterns-audit.md` (fila 2, líneas 35-36 y 54, 76, 89), `docs/handoffs/empty-state-audit.md` (plantilla de convenciones del patrón ya entregado en `eba1c4c`), `docs/handoffs/{composition-audit,ux-accessibility-audit,vocabulary-audit}.md`, `docs/gelium-ui-vocabulary.md` (103-110), `docs/gelium-ui-system-roadmap.md` (147-172 Phase D; 216-224, 436 Phase G), `docs/gelium-ui-composition-rules.md` (149-162 state matrix; 175-182 accesibilidad), `internal/app/{data_table,server,routes}.go`, `web/templates/{empty-state,button,progress,data-table,demo-whatsapp,layout}.html`, `web/styles/{empty-state,toast,button,progress,list,tokens,app,data-table,demo-whatsapp}.css`, `web/styles_contract_test.go`, `themes/theme-material/theme.css`.

---

## 1. CONTEXTO CARGA ACTUAL — qué existe hoy para "carga" (no skeleton)

Confirmado por grep: **cero referencias a `skeleton` en templates y styles** (solo en docs: `state-patterns-audit.md:15,36`, `vocabulary-audit.md:53,199`). Lo que existe hoy es el "contexto de carga" sobre el que se apoya el patrón:

| Qué existe | Evidencia | Detalle |
|---|---|---|
| **Button loading** | `web/templates/button.html:4,9` | `aria-busy="true"` + `.ui-button-spinner` (animación `ui-spin` definida en `app.css:52`) + `<span class="sr-only">Loading {Label}</span>` + label visible envuelto en `aria-hidden="true"`. Variante link: `role="link" aria-disabled="true" tabindex="-1"` sin `href`. Spinner: `button.css:41-47` con duración literal `.8s` |
| **Progress determinate/indeterminate** | `web/templates/progress.html:5-23` | `<progress>` nativo: `value`+`max` = determinate; sin `value` = indeterminate (animación nativa del browser). `progress.css`: track 4px, radius full |
| **Refresh de operación (demo)** | `data-table.html:82-91`, `data_table.go:388-423` | El único "carga diferida" del sistema hoy es un POST simulado: no-JS re-renderiza página con `.ui-progress` determinate `value=100` + toast inline; HTMX swapea el fragmento `data-table-refresh-form` + `HX-Trigger loom:toast`. Es progreso de operación, no carga de datos |
| **`.sr-only`** | `toast.css:70-80` | Única definición global del repo (clip 1px, `white-space: nowrap`, `border: 0`). Sin prefijo `ui-`. Es la clase que reutiliza el skeleton para el texto "Loading…" |
| **Visually-hidden alternativo** | `segmented-button.css:42-53` | `.ui-segmented-button-legend` con la misma técnica clip (scoped del componente, no reutilizable) |
| **`role="status"` precedente** | `empty-state.html:1`, `data-table.html:21` (notice), `text_field.go:70-71` | Live region implícita `polite`; el swap HTMX `outerHTML` inserta el nodo y el SR lo anuncia. Precedente directo del skeleton |
| **Animación demo (no componente)** | `demo-whatsapp.css:327-337,539-546` | Typing dots con `@keyframes demo-wa-bounce` (`1.2s`) + bloque `@media (prefers-reduced-motion: reduce) { animation: none; }`. Referencia de convención reduced-motion |
| **Reduced-motion consolidado** | `app.css:54-71` (entrada) | Bloque global histórico; los componentes B3+ declaran además bloque propio local (convención verificada por `TestCheckboxRadioReducedMotionCoverage`, `styles_contract_test.go:811-825`). El skeleton debe declarar bloque propio |
| **Forced-colors por componente** | `empty-state.css:52-60`, `data-table.css:340-359` | Cada primitiva declara su propio bloque `@media (forced-colors: active)` |

**Carga síncrona hoy**: todo el sistema es server-rendered síncrono en la primera render (Go + `html/template`, `server.go:79` parsea todos los templates; no existe fetch de datos en cliente — `app.js` solo maneja toast, slider y el swap 422 de validación; cero `hx-indicator`/`hx-disabled-elt` en templates). **Implicación documentada**: el skeleton NO es un estado que aparezca solo por una carga asíncrona real — hoy no existe carga asíncrona. Es un estado de UI que **el servidor renderiza** cuando: (a) una región va a ser llenada por un swap HTMX posterior (el GET inicial responde con skeleton, el fragmento siguiente con datos), o (b) Phase G: el handler sirve el skeleton como primer render de un listado que tarda (feed con datos remotos del lado servidor), reemplazado por el siguiente GET server-rendered. Cero contrato nuevo (coincide con `state-patterns-audit.md:76`: "sin server contract").

---

## 2. CONTRATO PROPUESTO — `.ui-skeleton`

### 2.1 Markup (partial `web/templates/skeleton.html`, patrón `empty-state.html`)

```html
{{define "skeleton"}}<div class="ui-skeleton{{if .Avatar}} ui-skeleton--avatar{{end}}" role="status">
  <span class="sr-only">{{if .Label}}{{.Label}}{{else}}Loading{{end}}</span>
  <div class="ui-skeleton-blocks" aria-hidden="true">
    {{if .Avatar}}<span class="ui-skeleton-block ui-skeleton-block--circle"></span>{{end}}
    <span class="ui-skeleton-block ui-skeleton-block--line ui-skeleton-block--title"></span>
    <span class="ui-skeleton-block ui-skeleton-block--line"></span>
    {{if .Avatar}}<span class="ui-skeleton-block ui-skeleton-block--line ui-skeleton-block--short"></span>{{end}}
  </div>
</div>{{end}}
```

Decisiones alineadas al repo:
- **Sin heading**: el skeleton es decorativo + announce; el texto vive en `.sr-only`, nunca visible (a diferencia del empty-state que muestra título/body).
- **`.Label`** opcional para anuncio contextual ("Loading messages", "Loading table") — default "Loading".
- **`.Avatar` bool** (no clase de variante suelta): la única composición con consumidor real hoy es fila con avatar (Feed/List three-line, ver §3).
- **`aria-hidden="true"`** en el contenedor visual `.ui-skeleton-blocks` (todos los bloques internos son decorativos); el texto sr-only es lo único en el árbol de accesibilidad.
- **`role="status"`** en la raíz → announce polite al insertarse por swap HTMX (mismo mecanismo que el empty row de la tabla).

### 2.2 Variantes — solo las que tienen consumidor real

| Variante | Bloques | Consumidor real |
|---|---|---|
| base (default) | 2 líneas (`--title` + `--line`) | Texto/listado simple, search results |
| `--avatar` | circle + 3 líneas (`--title`, `--line`, `--short`) | **Public/Social Feed** (recipe G: "Composición Card/List+Skeleton", `roadmap.md:436`; `composition-audit.md:57`) y listas tipo conversación |
| rect (KPI) | bloque `--rect` | **Dashboard** (diferido post-3, `roadmap.md:224`; state matrix `composition-rules.md:158`) |

Las alturas de fila imitan el contrato de List: `--ui-list-item-one-line-height: var(--ui-size-item)` 56px, `--two-line` 72px, `--three-line` 88px (`list.css:15-17`) — el feed es una List three-line, así que el skeleton de feed mide 88px por fila.

### 2.3 Tokens (todos verificados como existentes en `tokens.css` + theme)

- **`--ui-color-surface-container`** — **EXISTE** en core (`tokens.css:23`) y en theme en las 3 rutas (light + dark class + dark media, verificado por `TestSurfaceContainerTokenClosedAcrossCoreAndEveryScheme`). Es el color de los bloques.
- **`--ui-radius-sm` (.5rem)** para líneas, **`--ui-radius-full`** para el círculo, **`--ui-radius-md`** para rect — EXISTEN (`tokens.css:80-83`).
- **Spacing**: `--ui-space-1/2/3` para gap entre bloques — EXISTEN.
- **Tokens scoped** (convención `TestComponentSizeTokensDeclaredScoped`, `styles_contract_test.go:403-425`): declarar en `.ui-skeleton`:
  `--ui-skeleton-line-height: var(--ui-size-item);` (o 1rem para línea de texto), `--ui-skeleton-avatar-size: var(--ui-size-icon-lg…)` — **NOTA**: no existe `--ui-size-icon-lg`; usar `--ui-size-icon` (1.5rem) o un valor scoped propio. `--ui-skeleton-radius: var(--ui-radius-sm);`, `--ui-skeleton-color: var(--ui-color-surface-container);`.
- **Motion**: `--ui-motion-short/long` + `--ui-easing-standard` EXISTEN (`tokens.css:97-99`), pero **ninguna duración sirve para shimmer/pulse** (150ms/500ms son demasiado cortas). Decisión (ver §2.4): duración scoped literal, con precedente en el repo (`.ui-button-spinner` usa `.8s` literal, `button.css:47`; `TestDialogDrawerNoMotionLiterals` solo restringe `dialog.css`/`navigation-drawer.css`). **No agregar token de motion nuevo** — evita tocar core + theme + `motionCoreTokens` y no hay test que lo exija.

### 2.4 Animación — **pulse** (M3), no shimmer

- **Pulse** (opacity fade entre 0.5 y 1 sobre `surface-container`): fiel a Material 3 (el theming del repo es Material), simple, sin gradientes que repintar por bloque (mejor para listas largas de feed), y se apaga limpio con `animation: none`.
- **Shimmer** (gradiente barrido) descartado por: sin consumidor que lo justifique, más costoso en muchos bloques, y complejidad de `background-size` animado.
- Keyframes scoped en `skeleton.css` (convención: `ui-spin` vive en `app.css:52` a nivel entrada; las keyframes de componente viven en su archivo, p. ej. `demo-wa-bounce` en `demo-whatsapp.css:539`).
- **Reduced-motion: bloque propio obligatorio**:
  ```css
  @media (prefers-reduced-motion: reduce) {
    .ui-skeleton-block { animation: none; }
  }
  ```
  (convención B3+ con bloque local, `TestCheckboxRadioReducedMotionCoverage`; el estado resultante es bloques estáticos grises — correcto, el anuncio lo da `role="status"` + sr-only, no la animación).

### 2.5 ARIA

- `role="status"` en `.ui-skeleton` → `aria-live="polite"` implícito; el swap HTMX inserta el nodo y se anuncia. Precedente: empty-state (`empty-state.html:1`) y notice de tabla (`data-table.html:21`).
- `.sr-only` con el texto de carga (única definición global, `toast.css:70-80`).
- Bloques visuales `aria-hidden="true"` (contenedor `.ui-skeleton-blocks`).
- **`aria-busy="true"` va en la región contenedora** (feed/lista), NO en el skeleton: el skeleton ES el contenido del placeholder; la región en carga es la que el servidor marca. Lo pone el template que compone la región (documentado en `state-patterns-audit.md:36`), no el partial. Si el partial se usa standalone sin región, el campo `.Busy` del view puede renderizarlo en la raíz — decisión de implementación.
- **Mensaje nunca color-only**: el estado no se comunica por color (los bloques grises son decorativos); la semántica está en `role="status"` + texto.

### 2.6 CSS

- **NUEVO `web/styles/skeleton.css`**: `@layer components`, scoped tokens `--ui-skeleton-*`, keyframes `ui-skeleton-pulse`, bloque reduced-motion propio, bloque forced-colors (bloques → `background: CanvasText` con `forced-color-adjust: auto` para que el placeholder siga distinguiéndose en high contrast, patrón `empty-state.css:52-60`).
- `@import "./skeleton.css";` en `web/styles/app.css` (después de `empty-state.css`, línea 33) **y** en la lista `sourceAppCSS` de `web/styles_contract_test.go:24-58` (sync obligatoria, `state-patterns-audit.md:102`).
- **Sin demo propio**: la página docs standalone (ruta + `pageView` + `server.go`) es opcional igual que en empty-state (ver FILES); la integración primaria es Phase G.

---

## 3. CONSUMIDORES REALES — dónde se usaría

| Consumidor | Justificación (evidencia) | Timing |
|---|---|---|
| **Public/Social Feed recipe** | Recipe G #3 (`roadmap.md:222`); definición explícita "Composición Card/List+Skeleton" (`roadmap.md:436`); state matrix Feed Loading = **GAP** (`composition-rules.md:157`); composition audit: "Loading state de feed; skeleton" (`composition-audit.md:57`) | Phase G — consumidor primario de `--avatar` |
| **Dashboard** | state matrix Dashboard Loading = **GAP** (`composition-rules.md:158`); composition audit: "Skeleton; KPI card con delta" (`composition-audit.md:58`) | Diferido post-3 (`roadmap.md:224`) — consumidor de `--rect` |
| **Data table (refresh HTMX)** | Contrato GET existente `hx-get="/components/data-table" hx-target="#data-table-panel" hx-swap="outerHTML"` (`data-table.html:7-12`); hoy el refresh usa `.ui-progress` determinate (`data-table.html:87`) — el skeleton sería el placeholder de la región durante un refresh real con latencia | Opcional; el contrato actual es síncrono |

**Integración con los contratos existentes (no implementar, solo documentar)**:
- **Sin contrato server nuevo**: el skeleton es output del servidor (`state-patterns-audit.md:76`). En un flujo HTMX el servidor puede responder el primer GET de la región con skeleton (con `aria-busy` en la región) y el siguiente request/fragmento trae los datos reales vía el mismo swap `outerHTML` — cero cambio en `app.js` ni en los contratos (a) 422, (b) loom:toast, (c) GET params, (d) POST+303.
- **No-JS**: el skeleton solo aparece si el servidor decide servirlo en la primera render; el siguiente GET de página completa (link/filter real) lo reemplaza. El reemplazo por datos es responsabilidad del servidor, nunca del cliente.

---

## 4. TESTS PROPUESTOS

**NUEVO `web/styles_skeleton_test.go`** (patrón `styles_empty_state_test.go`):
1. `TestSkeletonPrimitiveCSSMapsTokens` — `.ui-skeleton {` + consumo de tokens scoped (`--ui-skeleton-*`) y core (`--ui-color-surface-container`, `--ui-radius-sm/full`) + keyframes de pulse.
2. `TestSkeletonContractCSSWired` — presencia en el CSS compilado `static/app.css` (espejo de `TestEmptyStateContractCSSWired`) + forced-colors.
3. `TestSkeletonClassVocabularyIsClosed` — clases del partial (`ui-skeleton`, `ui-skeleton--avatar`, `ui-skeleton-blocks`, `ui-skeleton-block`, `--line`, `--circle`, `--title`, `--short`) existen en template y CSS; prefijo ui reservado; sin `ui-skeleton-demo`.
4. `TestSkeletonReducedMotionDisablesAnimation` — `@media (prefers-reduced-motion: reduce)` con `animation: none` (espejo de `TestCheckboxRadioReducedMotionCoverage`).

**`web/styles_contract_test.go`**:
5. Agregar `"styles/skeleton.css"` a `sourceAppCSS` (líneas 24-58).
6. Entrada en `TestComponentSizeTokensDeclaredScoped` para el token scoped de tamaño (si se declara). `TestNoColorLiteralsInComponents` y `TestSpace/SizeTokensConsumedByComponents` cubren el archivo automáticamente.

**`internal/app/` (si se agrega página docs standalone)**:
7. Test de render del partial: `role="status"`, `.sr-only` con texto, bloques `aria-hidden="true"` (espejo de cómo `data_table_test.go` verifica el empty state en el HTML renderizado).

---

## 5. FILES IMPACTADOS (read-only — NO modificados por este audit)

| Archivo | Cambio |
|---|---|
| `web/styles/skeleton.css` | **NUEVO** — primitiva `.ui-skeleton` (+ variantes `--avatar`, bloques `--line/--circle/--rect`, keyframes pulse, reduced-motion, forced-colors) |
| `web/templates/skeleton.html` | **NUEVO** — partial `{{define "skeleton"}}` |
| `web/styles_skeleton_test.go` | **NUEVO** — tests de contrato CSS |
| `web/styles/app.css` | `@import "./skeleton.css";` (después de línea 33) |
| `web/styles_contract_test.go` | lista `sourceAppCSS` + (opcional) mapa de tokens scoped |
| `web/static/app.css` | regenerado por el build (compilado embebido, verificado por test de contrato) |
| `docs/gelium-ui-vocabulary.md` | Loading/Skeleton ✖→✅ (DoD, `state-patterns-audit.md:102`) |
| `docs/gelium-ui-composition-rules.md` | state matrix §8 (Feed/Dashboard Loading GAP→✅) |
| `COMPONENT-ROADMAP.md` | registro del componente (hoy sin entrada: grep sin resultados) |
| Opcional (no requerido para integración): `internal/app/routes.go`, `internal/app/server.go`, `internal/app/skeleton.go` | solo si se agrega página docs standalone del skeleton |
| Phase G (fuera de Phase D): templates/handlers de Feed y Dashboard | consumo real de `--avatar`/`--rect` |

**DoD completo** (`state-patterns-audit.md:102`): `styles_skeleton_test.go` + partial + entrada en `sourceAppCSS` + actualización vocabulario/matrix. **No requiere** token de motion nuevo, ni cambio en `tokens.css`/theme (duración scoped literal con precedente `button.css:47`), ni página standalone para la integración.

---

## 6. Fuentes de autoridad usadas

`docs/handoffs/{state-patterns-audit,empty-state-audit,composition-audit,ux-accessibility-audit,vocabulary-audit}.md`, `docs/gelium-ui-vocabulary.md` (103-110), `docs/gelium-ui-system-roadmap.md` (147-172, 216-224, 409, 436-437), `docs/gelium-ui-composition-rules.md` (149-162, 175-182), `internal/app/{data_table,server,routes}.go`, `web/templates/{empty-state,button,progress,data-table,demo-whatsapp,layout}.html`, `web/styles/{empty-state,toast,button,progress,list,tokens,app,data-table,demo-whatsapp}.css`, `web/styles_contract_test.go` (24-58, 403-425, 811-825), `web/styles_empty_state_test.go`, `themes/theme-material/theme.css`, `web/static/app.js`, commit `eba1c4c` (diff del Empty State como plantilla).
