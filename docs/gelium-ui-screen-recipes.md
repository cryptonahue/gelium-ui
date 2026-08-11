# Gelium UI — Screen Recipes (Phase G)

> **Objetivo**: definir cada pantalla como una *recipe* de 19 campos que ancla
> componentes y contratos **reales** del sistema (no inventados). Cada recipe es
> la composición de primitivas existentes sobre un contrato server explícito,
> con flujo no-JS completo y una mejora HTMX opcional que no introduce contratos
> nuevos.
>
> **Estado**: las tres recipes de Phase G están **implementadas** — **Admin
> Resource** (`/recipes/admin-resource`, `recipe_admin_resource.go`,
> `recipe-admin-resource.html`), **Ops Queue** (`/recipes/ops-queue`,
> `recipe_ops_queue.go`, `recipe-ops-queue.html`) y **Public/Social Feed**
> (`/recipes/public-feed`, `recipe_public_feed.go`, `recipe-public-feed.html`).
> Las dos últimas consumen las primitivas introducidas junto a ellas: **Avatar**
> (`avatar.html`/`avatar.css`/`avatar.go`), **tone variants de Badge**
> (`ui-badge--error|success|warning|info` en `badge.css`, con el on-color nuevo
> `--ui-color-info-fg`) y **pagination standalone** (`pagination.html`/
> `pagination.css`/`pagination.go`, la extracción reutilizable del footer de la
> Data table).
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

## 2. Ops Queue — IMPLEMENTADA

Rutas vivas (`internal/app/server.go`):

| Operación | Ruta | Contrato |
|---|---|---|
| Listar/filtrar/paginar | `GET /recipes/ops-queue` | GET params estables `?status=&kind=&page=` (vocabularios cerrados, sanitizados); `HX-Request` bifurca el fragmento `#queue-panel` |
| Avanzar al siguiente estado | `POST /recipes/ops-queue/{id}/advance` | 303 a la cola + banner success persistente (`pending→in_progress→done`; `blocked→in_progress`; `done` terminal → banner info) |
| Sacar de la cola | `POST /recipes/ops-queue/{id}/dequeue` | 303 a la cola + banner success; 404 con `error-state` si el id no existe |
| Refresh remoto | `POST /recipes/ops-queue/refresh` | POST-only; HX → fragmento + `HX-Trigger loom:toast`; no-JS → página con toast inline + progress |

### SURFACE

App shell server-rendered (appbar + columna de contenido fluida `max-width`, sin breakpoints) con slot de banner global debajo del appbar. La cola es una **página** — unidad de URL — no un overlay; el orden operativo es estado del servidor, nunca presentación cliente.

### USER

Operador/agente que atiende ítems de una cola en triaje secuencial: lee estado y SLA de cada ítem, avanza el siguiente, saca lo terminado y filtra por categoría. Copy en inglés, localizable por contrato.

### PRIMARY_TASK

Atender el "siguiente a atender": ver el estado y SLA de cada ítem (con avatar del requester y badge de tono) y **avanzarlo al siguiente estado** con un POST+303 real.

### SECONDARY_TASKS

Filtrar la cola por estado y categoría (GET con vocabularios cerrados), paginar (partial standalone), remover ítems de la cola, refresco remoto con progress + toast.

### UX_PATTERN

Queue (vocabulario `gelium-ui-vocabulary.md`) — composición List two-line + Avatar + Badge tone + Button + Toast + Empty state + Banner success sobre POST+303 — sin drag ni board (FIFO, anti-regla de `composition-rules.md`).

### VISUAL_VOCABULARY

List two-line (fila de la cola), Avatar (requester, sm, decorativo con nombre visible), Badge tone (estado/SLA: `--info/--warning/--error/--success`), Button (Acción por fila), Banner success persistente post-303, Empty state, Skeleton (placeholder documentado), Toast (refresh). Los tones viven en el primitivo (`ui-badge--*`), nunca en CSS de la recipe.

### COMPONENTS

List (`ui-list` two-line), **Avatar** (`avatar.html`), **Badge tone** (`ui-badge ui-badge-large ui-badge--{tone}`), Button / Icon button, Select (filtros status/kind), **Pagination standalone** (`pagination.html`), Empty state, Banner, Toast, Progress (refresh), `error-state` (404). Primitivas nuevas consumidas: Avatar, tone variants y pagination standalone.

### STATES

- **Rest**: cola con ítems ordenados operativamente (pending → in_progress → blocked → done, FIFO dentro del mismo estado).
- **Empty**: dos variantes server-driven — con filtros activos (`No matching items` + CTA "Clear filters") y sin filtros (`Queue is clear` + CTA "View completed").
- **Loading**: no aplica a la carga inicial (server-rendered); el refresh usa `.ui-progress` determinate; el placeholder de carga está documentado en la recipe Feed.
- **Tone por ítem**: derivado server-side del estado + SLA (overdue → error, cerca del deadline → warning, demás → info; blocked → error, done → success) con el label de estado siempre visible (nunca color-only).
- **Error**: íd de ítem inválido → `error-state` 404 con retry; global persistente → `banner--error` (nunca toast).
- **Success**: persistente post-303 con `banner--success` (`role="status"`); transitorio solo en refresh con toast.

### ACCESSIBILITY

Avatar `aria-hidden="true"` (decorativo, el nombre del requester está en texto visible); estado del ítem = badge con label textual + tone (nunca color-only); acciones por fila = `<button>` reales en forms POST; `aria-current="page"` en paginación standalone; banner deriva `role="status"`; `#loom-toast-region` `aria-live` para transitorio; forced-colors cubre avatar, badge y pagination; `aria-label` en los selects del filtro vía `<label>` asociado.

### CONTENT_RULES

Botón de avance = verbo del estado destino ("Advance", nunca "OK"); SLA vencido = label factual ("SLA overdue"); empty = qué + CTA real; el copy de la cola es terso (subject + requester + kind + tiempo + SLA en una línea supporting); copy en inglés por contrato, localizable.

### SEO_REQUIREMENTS

Ruta `/recipes/ops-queue` → `robots: noindex, nofollow` (superficie de demo). `lang="en"`, `<title>` y `<meta name="description">` propios de la recipe, canonical limpio sin query (`siteBaseURL + /recipes/ops-queue`), un solo `<h1>` por página. Los parámetros `?status=&kind=&page=` nunca entran al canonical.

### GEO_REQUIREMENTS

Contenido factual citado del vocabulario del sistema (estados pending/in_progress/done/blocked, contrato 303/422/`loom:toast`); URLs estables y deep-linkables por `{id}`; sin claim universal no verificado; Gelium UI es la entidad única.

### SERVER_CONTRACT

`GET ?status=&kind=&page=` con vocabularios cerrados sanitizados (`recipeQueueStatuses`, `recipeQueueKinds`, page ≥ 1 → defaults seguros); `POST + 303 SeeOther` para cada transición (advance/dequeue) con banner flash consumido en el siguiente render; el orden operativo y el tone se derivan **server-side** (`recipeQueueRank`, `recipeQueueItemTone`); `HX-Trigger: {"loom:toast":{...}}` solo para refresh; `HX-Request` bifurca `#queue-panel` vs página completa; `/recipes/ops-queue/refresh` en `postOnlyPaths()` → GET responde 405 con `Allow: POST`.

### NO_JS_FLOW

Completo: filtro = form GET real (URL es el estado), paginación = links reales, avance/remoción = form POST + 303 con banner success persistente, refresh = reload con toast inline (`toast.go` fallback no-JS), 404 = `error-state` con retry real.

### HTMX_ENHANCEMENT

Swap de `#queue-panel` (`outerHTML`) al filtrar (hx-get en el form de filtro); refresh devuelve el fragmento del form + `HX-Trigger loom:toast`. Las mutaciones permanecen POST+303 (sin `hx-post`) — el contrato de mutación no cambia.

### RESPONSIVE_BEHAVIOR

Columna fluida `max-width` sin breakpoints; filas full-width con acciones inline (nunca overflow horizontal); paginación server-side si la cola crece; el filtro envuelve en narrow (flex-wrap) con selects a ancho fluido.

### THEME_REQUIREMENTS

Solo tokens `--ui-*` existentes + el on-color nuevo `--ui-color-info-fg` (requerido por la variante info del Badge, definido en core + theme en los 3 esquemas); cero literales de color en el CSS de la recipe (guard `TestNoColorLiteralsInComponents` cubre `recipe-ops-queue.css`); light/dark vía theme, forced-colors y reduced-motion por primitivos; sin tokens de densidad/breakpoints nuevos.

### ALTERNATIVES_REJECTED

**Board**: FIFO estricto = queue, no multi-vía (anti-regla `composition-rules.md`). **Data table**: cada fila es una unidad discreta con acción, no comparación columnar (criterio 4.1/4.2). **Drag & drop**: server-first, POST+303. **Toast para SLA persistente**: el SLA es estado persistente, va en label+badge, nunca transitorio. **Badge sin label**: tone nunca color-only — el label de estado siempre acompaña.

### RATIONALE

Es la composición de workflow más simple del sistema (List+badge+button+POST+303, ya probada en la sidebar WhatsApp) y demuestra que el orden y el tono son **estado server-side**. Ejercita las tres primitivas nuevas del slice (Avatar, tone, pagination) con el mínimo de contrato nuevo y desbloquea el patrón Queue para recursos reales.

---

## 3. Public/Social Feed — IMPLEMENTADA

Rutas vivas (`internal/app/server.go`):

| Operación | Ruta | Contrato |
|---|---|---|
| Ver el feed | `GET /recipes/public-feed` | GET params estables `?view=&page=` (vocabulario cerrado `for-you|following|new`); `HX-Request` bifurca el fragmento `#feed-panel` |
| Reaccionar (like) | `POST /recipes/public-feed/{id}/react` | 303 al feed + toast flash transitorio en el siguiente render; 404 con `error-state` si el id no existe |
| Refresh remoto | `POST /recipes/public-feed/refresh` | POST-only; HX → fragmento + `HX-Trigger loom:toast`; no-JS → página con toast inline + progress |

### SURFACE

Página pública del feed (app shell server-rendered, columna fluida `max-width`) orientada a la **novedad** (criterio 4.5): el ítem ES el evento, no una entidad administrable. El estado de carga y el vacío son output del servidor, nunca flash de cliente.

### USER

Usuario público/consumidor que se entera de lo nuevo: escanea rápido, cambia de vista, reacciona y sigue leyendo. Copy en inglés, localizable por contrato.

### PRIMARY_TASK

Enterarse de la novedad: escaneo rápido de posts en orden cronológico inverso (los más nuevos primero) con avatar del autor, tiempo relativo y marcador "New".

### SECONDARY_TASKS

Cambiar de vista (For you / Following / New vía Tabs server-side con `aria-current`), reaccionar (like), ver conteos (likes/comments), cargar más (paginación server-side), refresh remoto.

### UX_PATTERN

Feed (vocabulario `gelium-ui-vocabulary.md`) + Loading (#8, Skeleton server-driven) + Empty (#7) + Pagination (#6, partial standalone) + Tabs server-side. El "load more" es paginación real (`?page=N`) — la URL sigue siendo el estado.

### VISUAL_VOCABULARY

Card (unidad repetida autocontenida, criterio 4.4) con header de autor (Avatar + nombre + kind + tiempo), Badge "New" (tone info) y Badge count (comments); Tabs; Skeleton (placeholder documentado); Empty state; Toast (reacción/refresh). Los tones viven en los primitivos, nunca en CSS de la recipe.

### COMPONENTS

Card (`ui-card ui-card-outlined`), **Avatar** (`avatar.html`, sm, decorativo con nombre visible), Badge (tone "New" + count), Tabs (`ui-tabs`/`ui-tab`), Button (like + refresh), **Pagination standalone**, Skeleton (placeholder de carga), Empty state, Toast (flash + refresh), Progress (refresh), `error-state` (404). Primitivas nuevas consumidas: Avatar y pagination standalone.

### STATES

- **Rest**: feed con posts en orden cronológico inverso.
- **Loading**: la carga inicial es server-rendered (sin fase cliente); la recipe documenta y renderiza el placeholder Skeleton (`role="status"` + sr-only "Loading the feed") como demo servida.
- **Empty**: por vista — "Nothing new yet" / "No posts from people you follow" / "No posts yet", cada una con CTA real (ver todos / refresh).
- **Error**: reacción a íd inválido → `error-state` 404 con retry; transporte → toast error transitorio (`app.js` `htmx:responseError`).
- **New**: marcador con label ("New" en badge tone info) + `aria-label` de la tarjeta — nunca color-only.
- **Success**: reacción → toast flash transitorio post-303; refresh → toast inline/`loom:toast`.

### ACCESSIBILITY

Avatar `aria-hidden` + nombre del autor en texto visible (nunca inicial decorativa sola); tarjetas con `aria-label="Post by <autor>"`; Tabs como links reales con `aria-current="page"` (sin `role="tablist"` ni JS); skeleton con `.sr-only` "Loading"; toast flash en `role="status"`; paginación con `aria-current="page"` y boundary `aria-disabled`; forced-colors cubre avatar/badge/pagination/tabs; nunca color-only.

### CONTENT_RULES

Copy terso y factual; timestamps relativos ("12m ago", "3h ago"); el "like" es verbo de acción ("Like · 24"); empty = mensaje + CTA real; nunca "cargando…" como texto — el loading se muestra con Skeleton/progress (`content-rules.md`); copy en inglés por contrato, localizable.

### SEO_REQUIREMENTS

Ruta `/recipes/public-feed` → `robots: noindex, nofollow` (superficie de demo). `lang="en"`, `<title>` y `<meta name="description">` propios de la recipe, canonical limpio sin query (`siteBaseURL + /recipes/public-feed`), un solo `<h1>` por página. Si el feed real fuera público e indexable, `?view=&page=` quedarían excluidos del canonical y se añadiría JSON-LD — no aplica a esta recipe de demo.

### GEO_REQUIREMENTS

Contenido del feed factual/verificable (items de demo sin claims no verificados); vistas y estados descritos con el vocabulario del sistema; URLs estables deep-linkables por `{id}`; entidad única Gelium UI.

### SERVER_CONTRACT

`GET ?view=&page=` con vocabulario cerrado sanitizado (`recipeFeedViews`, page ≥ 1 → defaults seguros); orden cronológico inverso **server-side**; `POST + 303 SeeOther` para reacciones con toast flash consumido en el siguiente render; `HX-Trigger: {"loom:toast":{...}}` solo para refresh; `HX-Request` bifurca `#feed-panel` vs página completa; `/recipes/public-feed/refresh` en `postOnlyPaths()` → GET responde 405 con `Allow: POST`; empty/loading son output del servidor (nunca flash cliente).

### NO_JS_FLOW

Completo: vistas = links reales (Tabs), paginación = links reales, like = form POST + 303 con toast flash inline, refresh = reload con toast inline, 404 = `error-state` con retry real.

### HTMX_ENHANCEMENT

Swap de `#feed-panel` (`outerHTML`) al cambiar de vista (hx-get en los tabs) y paginar; refresh devuelve el fragmento del form + `HX-Trigger loom:toast`. Las reacciones permanecen POST+303 (sin `hx-post`) — el contrato de mutación no cambia.

### RESPONSIVE_BEHAVIOR

Columna fluida `max-width` sin breakpoints; tarjetas a ancho completo (single-column) con header que envuelve en narrow; paginación server-side en vez de scroll infinito; sin breakpoints — grid fluido/`min()`/`clamp()`.

### THEME_REQUIREMENTS

Solo tokens `--ui-*` existentes (incluido `--ui-color-info-fg` para el badge "New"); cero literales de color en el CSS de la recipe (guard `TestNoColorLiteralsInComponents` cubre `recipe-public-feed.css`); light/dark vía theme, forced-colors y reduced-motion (shimmer del skeleton apagado por primitivo); sin tokens de densidad/breakpoints nuevos.

### ALTERNATIVES_REJECTED

**Collection/Data table**: el feed es evento/tiempo, no entidad/set con filtros (criterio 4.5). **Timeline**: no hay eje temporal de proceso explícito, es actividad reciente tersa (criterio 4.6). **Board**: sin multi-vía. **Scroll infinito cliente**: server-first, la URL es el estado → paginación real. **Spinner ad-hoc**: anti-regla (`composition-rules.md`) — se usa Skeleton para carga y `.ui-progress` para refresh.

### RATIONALE

Es la recipe que obliga a resolver **loading + novedad + avatar/paginación** — los patrones de carga más críticos del sistema según `composition-audit.md`. Valida el contrato de estado de carga server-driven (Skeleton documentado, nunca flash de empty) y el patrón de "load more" por paginación real sobre la composición Card+Avatar+Badge.
