# Gelium UI — Screen Recipes (Phase G)

> **Objetivo**: definir cada pantalla como una *recipe* de 19 campos que ancla
> componentes y contratos **reales** del sistema (no inventados). Cada recipe es
> la composición de primitivas existentes sobre un contrato server explícito,
> con flujo no-JS completo y una mejora HTMX opcional que no introduce contratos
> nuevos.
>
> **Estado**: la primera recipe — **Admin Resource** — está implementada
> (`/recipes/admin-resource`, handlers en `internal/app/recipe_admin_resource.go`,
> template `web/templates/recipe-admin-resource.html`, layout `web/styles/recipe-admin-resource.css`,
> tests `internal/app/recipe_admin_resource_test.go`). **Ops Queue** y
> **Public/Social Feed** quedan **pendientes** (requieren Avatar, tone variants y
> pagination standalone, según `docs/handoffs/screen-recipes-audit.md`).
>
> **Fuentes de autoridad**: `docs/gelium-ui-composition-rules.md` (gramática de
> pantalla, criterios 4.1-4.8, state matrix, anti-rules), `docs/gelium-ui-ux-patterns.md`,
> `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-content-rules.md`,
> `docs/gelium-ui-seo-contract.md`, `docs/gelium-ui-geo-contract.md`,
> `docs/handoffs/screen-recipes-audit.md`.

---

## 1. Admin Resource — IMPLEMENTADA

Rutas vivas (`internal/app/server.go`):

| Operación | Ruta | Contrato |
|---|---|---|
| Listar/filtrar/ordenar/paginar/seleccionar | `GET /recipes/admin-resource` | GET params estables `?q=&sort=&dir=&page=&selection=`; `HX-Request` bifurca el fragmento `#resource-panel` |
| Crear (form) | `GET /recipes/admin-resource/new` | página completa no-JS |
| Crear (mutación) | `POST /recipes/admin-resource` | 303 a la lista (éxito, banner persistente) o 422 + `X-Loom-Validation` (fallo) |
| Editar (form) | `GET /recipes/admin-resource/{id}/edit` | página completa; 404 con `error-state` si el id no existe |
| Editar (mutación) | `POST /recipes/admin-resource/{id}/edit` | 303 o 422 |
| Borrar (confirmación) | `GET /recipes/admin-resource/{id}/delete` | Dialog page variant (`<dialog open>` nativo) |
| Borrar (mutación) | `POST /recipes/admin-resource/{id}/delete` | 303 a la lista + banner success |
| Refresh remoto | `POST /recipes/admin-resource/refresh` | POST-only; HX → fragmento + `HX-Trigger loom:toast`; no-JS → página con toast inline |

### SURFACE

Página de app shell server-rendered: appbar + columna de contenido (`max-width` fluida, sin breakpoints). El listado es una **página** — unidad de URL (gramática de superficie, `composition-rules.md`), no un overlay ni un dialog. El banner de estado vive en el slot global de la página (debajo del appbar), no dentro de la tabla.

### USER

Administrador/operador que gestiona un set de recursos (proyectos): encuentra, compara, selecciona y actúa por fila. Copy en inglés, localizable por contrato.

### PRIMARY_TASK

Gestionar el set de recursos: **encontrar** (filtro/orden/paginación server-side), **seleccionar** y **actuar por fila** (editar/borrar), **crear** nuevos recursos con validación.

### SECONDARY_TASKS

Filtrar/buscar por nombre, estado u owner; ordenar por columnas; paginar; crear recurso; editar recurso; borrar individual con confirmación; refresh remoto de la lista.

### UX_PATTERN

Resource list (#3 `ux-patterns.md`) + Pagination (#6) + Bulk action (#11) + Destructive (#10) + Confirmation (#18) + Error recovery (#9). El patrón Data table server-driven es el vehículo de la lista.

### VISUAL_VOCABULARY

Data table (colección, criterio 4.1) con Empty state integrado; Skeleton no aplica al listado inicial (server-rendered, sin fase de carga cliente); Inline alert + Validation summary en el form; Banner persistente post-303; Toast transitorio solo para refresh. Los tones de estado viven en los primitivos (`banner--success`, `inline-alert--error`), nunca en CSS de la recipe.

### COMPONENTS

Data table (`data-table.html`/`data_table.go`, reuso de columnas/paginación/empty), Text field (search + form fields), Select (estado), Button / Icon button, Dialog (page variant de confirmación), Banner (success persistente), Inline alert + Validation summary (422), Empty state, Toast (refresh), Progress (refresh). **Cero primitivas nuevas** — la recipe es 100% wiring.

### STATES

- **Rest**: tabla con filas y selección desmarcada.
- **Empty**: dos variantes server-driven — con búsqueda (`No results` + CTA "Clear search") y sin datos (`No projects yet` + CTA "New project").
- **Loading**: no aplica a la carga inicial (server-rendered); el refresh usa `.ui-progress` determinate.
- **Error**: campo → 422 con `role="alert"` por campo + validation-summary; recurso → `error-state` 404; global persistente → `banner--error` (nunca toast).
- **Selected**: checkboxes nativos `:checked` renderizados server-side desde `?selection=`; select-all con `aria-label` oculto cuando no hay filas.
- **Success**: persistente post-303 con `banner--success` (`role="status"`, nunca `loom:toast`); transitorio solo en refresh con toast.

### ACCESSIBILITY

`aria-sort` en la columna activa; `aria-current="page"` en paginación; checkboxes por fila con `aria-label` (nombre del recurso); Dialog page variant con `<dialog open>` nativo, `aria-labelledby`/`aria-describedby` y Cancel como link real; 422 con `aria-invalid` + `aria-describedby` por campo y validation-summary con links reales a `#campo-error`; los tones de banner/inline/empty derivan `role="status"`/`alert`; nunca color-only (forced-colors del sistema). El select-all queda oculto con 0 filas para no ofrecer un control muerto.

### CONTENT_RULES

Botones = verbos de acción ("New project", "Delete Alpha release", nunca "OK"); confirmación = qué se hace + irreversibilidad explícita ("This action cannot be undone."); errores = qué pasó + cómo resolverlo ("Name is required."); empty = mensaje + CTA real; copy en inglés por contrato, localizable.

### SEO_REQUIREMENTS

Toda ruta `/recipes/*` emite `robots: noindex, nofollow` (superficie de demo, igual que `/demo/*`). `lang="en"`, `<title>` y `<meta name="description">` **por ruta** (lista / new / edit / delete tienen copy distinta), canonical limpio sin query (`siteBaseURL + routePath`), un solo `<h1>` por página. El documento `screen-recipes.md` renderizado vía layout sería la única superficie indexable de Phase G.

### GEO_REQUIREMENTS

Contenido factual citado del vocabulario del sistema (estados Active/Pending/Done, contrato server 422/303/`loom:toast`); URLs estables y deep-linkables por `{id}`; sin claim universal no verificado; Gelium UI es la entidad única (no hay marca externa).

### SERVER_CONTRACT

`GET ?q=&sort=&dir=&page=&selection=` con vocabularios cerrados sanitizados (`dataTableSortKeys`, `dataTableStatuses`, page ≥ 1 → defaults seguros); `POST + 303 SeeOther` para todas las mutaciones; `422 + X-Loom-Validation: true` para validación (nunca toast); `HX-Trigger: {"loom:toast":{...}}` para transitorio (refresh); banner/inline success para persistente; `HX-Request` bifurca fragmento `#resource-panel` vs página completa; `POST /recipes/admin-resource/refresh` registrado en `postOnlyPaths()` → GET responde 405 con `Allow: POST`.

### NO_JS_FLOW

Completo: filtro/orden/página/sort = links GET reales (la URL es el estado); selección = checkboxes en form real; crear/editar/borrar = form POST + 303; 422 = re-render de página completa preservando valores; confirmación = `<dialog open>` server-rendered con form POST real; refresh = reload de página con toast inline (`toast.go` fallback no-JS); banner success persistente post-303.

### HTMX_ENHANCEMENT

Swap de `#resource-panel` (`outerHTML`) en sort/filter/page; refresh devuelve el fragmento del form + `HX-Trigger loom:toast`; focus al validation-summary en 422 (opcional). Las mutaciones de form permanecen POST+303 (sin `hx-post`) — el contrato de mutación no cambia.

### RESPONSIVE_BEHAVIOR

Columna fluida `max-width` sin breakpoints; paginación server-side en vez de scroll horizontal de tabla; grid del form con `min()`/`clamp()` para campos; en narrow las acciones de fila quedan como links inline (nunca overflow horizontal).

### THEME_REQUIREMENTS

Solo tokens `--ui-*` existentes; cero literales de color en el CSS de la recipe (guard `TestNoColorLiteralsInComponents` lo verifica sobre `recipe-admin-resource.css`); light/dark vía theme, forced-colors y reduced-motion cubiertos por los primitivos; sin tokens de densidad/breakpoints nuevos.

### ALTERNATIVES_REJECTED

**Board**: es un set con filtros, no un workflow multi-vía (criterio 4.1). **List** para >8-10 filas: anti-regla (`composition-rules.md`) — la colección columnar pide Data table. **Dialog para el listado completo**: página = URL/back (anti-regla "no flujo largo en dialog"). **Toast para errores de validación**: contrato 422 + `X-Loom-Validation`, persistente nunca transitorio. **Drag & drop**: server-first, POST+303.

### RATIONALE

Es la recipe que ejercita la mayor cantidad de contratos canónicos (GET params, selección, POST+303 destructivo, Dialog confirm, 422, empty, toast, banner/inline persistente) sobre la primitiva más madura del sistema (Data table). Valida el marco de recipes (template de pantalla + handler + store + fragmento HX) con riesgo mínimo y establece el patrón que reusarán Ops Queue y Feed.

---

## 2. Ops Queue — PENDIENTE

> Bloqueada por: **Avatar** (primitiva nueva), **tone variants de Badge/Chip** y
> **pagination standalone**. La composición objetivo está documentada en
> `docs/handoffs/screen-recipes-audit.md` (§4). Template/handlers/rutas a crear:
> `web/templates/recipe-ops-queue.html`, `internal/app/recipe_ops_queue.go`,
> `GET /recipes/ops-queue` + mutaciones POST+303 por transición.

Los 19 campos de esta recipe se completarán en la implementación, siguiendo el
mismo marco que Admin Resource (surfaces, contrato server, no-JS flow y
HTMX enhancement).

---

## 3. Public/Social Feed — PENDIENTE

> Bloqueado por: **Avatar** (igual que Ops Queue) y **pagination/"load more"**.
> La composición objetivo está documentada en `docs/handoffs/screen-recipes-audit.md`
> (§5). Template/handlers/rutas a crear: `web/templates/recipe-social-feed.html`,
> `internal/app/recipe_social_feed.go`, `GET /recipes/social-feed` + reacciones
> POST+303 + refresh.

Los 19 campos de esta recipe se completarán en la implementación, siguiendo el
mismo marco que Admin Resource.
