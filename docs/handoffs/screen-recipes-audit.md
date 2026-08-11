# Gelium UI — Screen Recipes Audit (Phase G, handoff)

> **Alcance**: auditoría read-only del estado real del sistema para implementar las 3 primeras screen recipes de Phase G (Admin Resource, Ops Queue, Public/Social Feed). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff (`docs/handoffs/screen-recipes-audit.md`).
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase G :216-251), `docs/gelium-ui-composition-rules.md` (screen/surface grammar, criterios 4.1-4.8, state matrix :164-178, server-driven :179-190, anti-rules :136-149), `docs/gelium-ui-ux-patterns.md`, `docs/gelium-ui-vocabulary.md` (Queue :170, Feed :208, Pagination :232, Empty :94, Skeleton :103, Inline alert :112, Banner :122, Success :140), `docs/gelium-ui-content-rules.md`, `docs/gelium-ui-seo-contract.md`, `docs/gelium-ui-geo-contract.md`, `docs/handoffs/{composition-audit,state-patterns-audit,success-feedback-audit}.md`, `internal/app/{server,routes,data_table,demo_whatsapp,toast,text_field,list,badge,docs}.go`, `web/templates/{layout,data-table,demo-whatsapp,demo-whatsapp-admin,list,badge,empty-state,skeleton,inline-alert,banner,callout,error-state,validation-summary,toast}.html`, `web/styles_contract_test.go`, `web/static/app.js`, `web/styles/{badge,list,demo-whatsapp}.css`.

---

## 1. RESUMEN EJECUTIVO

- **Phase D está ENTREGADA al 100%** (commits `eba1c4c`→`94340dc`, verificado en `git log`): Empty state, Skeleton, Inline alert, Validation summary, Banner, Callout, Error state, Success feedback (por reuso) + Toast. Las 3 recipes NO tienen bloqueo por patrones de estado. El gap señalado en `state-patterns-audit.md` (6 de 10 sin entregar) quedó cerrado después de esa auditoría.
- **El feedback de transporte HTMX (G5 de Phase E) también está cerrado**: `app.js:88-94` maneja `htmx:responseError`/`sendError` con toast error transitorio ("We couldn't reach the server. Try again.").
- **El Data table ya tiene Empty state integrado** y el bug G4 (select-all checked con 0 filas) quedó corregido: el checkbox "Select all" está oculto cuando `Total == 0` (`data-table.html:42`) y la fila vacía renderiza `empty-state` con CTA real (`data-table.html:67-71`).
- **No existe `docs/gelium-ui-screen-recipes.md`** (el entregable de Phase G). Nada que actualizar, todo por crear.
- **Gaps reales que SÍ bloquean/complican recipes**: (1) **Avatar** — no existe como componente; el demo lo hace ad-hoc (`.demo-wa-avatar` + inicial, `demo-whatsapp.html:35,62`). Bloquea Ops Queue y Feed. (2) **Indicador de tono reusable** — Badge es solo error-tinted (`badge.css:2-5`, token único `--ui-badge-container`), sin variantes success/warning/info; los tones viven solo a nivel token (`--ui-color-success/warning/info/danger-container` + `-fg`) y a nivel demo ad-hoc (`.demo-wa-window--{tone}`, `.demo-wa-quality--{tone}`, `demo-whatsapp.css:217,501-502`). Bloquea Ops Queue. (3) **Pagination standalone** — solo vive dentro del Data table (`data-table.html:74-78`, `vocabulary.md:232` lo marca pendiente). Ops Queue y Feed lo necesitan. (4) **View models Go de los state patterns** — banner/error/empty tienen view model en producción (`server.go:41-66`, `data_table.go:87-95`); skeleton, validation-summary, inline-alert y callout SOLO tienen template+CSS y view model en tests (`server_test.go:242`). Las recipes necesitan modelos de producción para render server-driven. (5) **Slot de "trailing meta/acción de fila" en List** — la List entregada tiene headline/supporting/leading y trailing icon, pero no badge ni acción por fila; se puede componer dentro del `<li>` (como el demo), pero conviene decidirlo antes.
- **Admin Resource es la recipe con menos gaps** (cero primitivas nuevas; todo lo demás es wiring Go): es la candidata natural a implementarse primero.
- **Discrepancias detectadas**: header de validación en código es `X-Loom-Validation` (`app.js:4`, `text_field.go:88`, `toast.go:178`, `server_test.go:109`) mientras el roadmap/docs escriben `X-Gelium-Validation` (`gelium-ui-system-roadmap.md:60`); y el conteo de campos de la recipe: el roadmap lista 18 (`roadmap.md:229-247`) pero su DoD dice "17 campos" (`roadmap.md:251`) y este encargo exige 19 (agrega SECONDARY_TASKS). Las recipes deben usar el header real del código y el set de 19 campos de este encargo.

---

## 2. BASELINE VERIFICADO (lo que existe hoy, con evidencia)

### 2.1 Phase D — estado real (actualización de `state-patterns-audit.md`)

| Patrón | Estado | Evidencia | Nota para recipes |
|---|---|---|---|
| Empty state | ✅ Entregado | `empty-state.html`, `emptyStateView` (`data_table.go:87-95`), integrado en Data table (`data-table.html:67-71`) | CTA es link real; `role="status"`; variante `Compact` para `<td colspan>` |
| Skeleton | ✅ Entregado | `skeleton.html` (variant `Avatar`, `Lines`), `skeleton.css`, token `--ui-skeleton-avatar-size` | Sin view model Go de producción; con `role="status"` + `.sr-only` |
| Inline alert | ✅ Entregado | `inline-alert.html` (tone→`role` derivado: error→alert, resto→status), `inline-alert.css` | Reemplaza los notices ad-hoc; sin view model Go de producción |
| Validation summary | ✅ Entregado | `validation-summary.html` (`HeadingLevel`, `Items[].Href` — links reales a `#campo-error`), `validation-summary.css` | Sin view model Go de producción; contrato 422 listo para usarlo |
| Banner | ✅ Entregado | `banner.html` + `bannerView` (`server.go:41-51`) + slot en `layout.html:24` | Slot sin handler que lo inyecte aún; dismiss = POST+303 |
| Callout | ✅ Entregado | `callout.html`, `callout.css` | No requerido por las 3 recipes (informativo); no bloquea |
| Error state | ✅ Entregado | `error-state.html` + `errorStateView` (`server.go:59-66`), 404 catch-all wired (`server.go:355-364`) | h1 único, retry = link real |
| Success persistente | ✅ Por reuso | `inline-alert--success` / `banner--success` (`role="status"`), pinneado en `server_test.go:296-307,334-345` | NUNCA `loom:toast` (guard `TestPersistentSuccessPartialsNeverToast`, `styles_contract_test.go:842-851`) |
| Toast | ✅ Completo | contrato `{"loom:toast":{...}}`, `#loom-toast-region` `aria-live`, fallback no-JS | Único transitorio-de-acción |
| Transporte HTMX (G5) | ✅ Cerrado | `app.js:88-94` (`htmx:responseError`/`sendError` → toast error) | Ya no es gap de Phase E |

### 2.2 Contratos server canónicos (reusar, no inventar)

1. **422 + `X-Loom-Validation: true`** para validación de campos/valores (nunca toast): `text_field.go:64-68,87-89`, `select.go:94-96`, hook `app.js:1-9`. **OJO**: el nombre real en el código es `X-Loom-Validation` (ver §1 discrepancias).
2. **`HX-Trigger: {"loom:toast":{...}}`** para feedback transitorio de acción: `toast.go:108-127,154-160`, reuso en `data_table.go:365-385`. Vocabulario cerrado `info|success|warning|error` (`toast.go:45`).
3. **GET params estables** para estado de listados (`q`, `sort`, `dir`, `page`, `selection`) con vocabularios cerrados sanitizados: `data_table.go:23,29,142-157,302-322`. URL es el estado.
4. **POST + 303 SeeOther** para mutaciones/workflow: `demo_whatsapp.go:559,573,584`. Sin fragmentos para mover estados.
5. **HX-Request bifurca** fragmento vs página completa: `data_table.go:120-130`.
6. **Metadata server-driven**: `layout.html:3-14` + `resolveMeta` (`server.go:175-226`); demos/examples → `noindex, nofollow` (`server.go:184-186`).

### 2.3 Qué existe para composición (inventario verificado)

- **28 rutas de componente** (`routes.go:16-47`) cubriendo acciones, input, feedback, navegación, datos. Nada falta para las 3 recipes a nivel de catálogo Material.
- **Demo WhatsApp** = única "screen real": master-detail (`demo-whatsapp.html`), POST+303, search por `?q=`, selección por `?c=`, sidebar tipo queue con unread + tone de ventana (`demo-whatsapp.html:43-46`). Es la referencia viva para Ops Queue y Feed.
- **Data table** = patrón de referencia server-driven completo (sort/filter/page/selection/empty/refresh con `.ui-progress` + toast, `data_table.go:354-389`).
- **Tokens**: `--ui-space-*`, `--ui-size-*`, `--ui-border-*`, `--ui-motion-*`, semánticos `--ui-color-{success,warning,info,danger-container}-*` en core (`tokens.css`), override por theme Material. Sin tokens de densidad/breakpoints (decisión Phase B: no tokenizar sin consumidor).

---

## 3. RECIPE 1 — ADMIN RESOURCE

### 3.1 Componentes: existentes / faltantes

| Necesidad | Estado | Evidencia |
|---|---|---|
| Data table (sort/filter/page/selection) | ✅ Existe | `data-table.html`, `data_table.go` |
| Empty state de tabla (mensaje + CTA) | ✅ Existe | `data-table.html:67-71` + `empty-state.html` |
| Button / Icon button | ✅ | `button.go`, `icon_button.go` |
| Menu (overflow row actions) | ✅ | `menu.html`/`menu.go` (Popover declarativo, zero JS) |
| Dialog (confirmar borrado) | ✅ | `dialog.html` (`closedby="any"`, Cancel `autofocus` — `dialog.go:18`) |
| Toast (resultado de acción) | ✅ | `toast.go` |
| Badge (counts/labels) | ✅ | `badge.go` (sin tones, ver §6 gap) |
| Text field (search) + GET form | ✅ | `text-field.html`, patrón `data-table.html:7-12` |
| Chips (filters) | ✅ | `chips.go` |
| Banner (error global persistente) | ✅ slot | `layout.html:24`, `bannerView` — sin handler que lo inyecte (wiring en Phase G) |
| Skeleton (carga inicial) | ✅ | `skeleton.html` (falta view model Go) |
| Avatar | ➖ No necesario | — |
| Pagination standalone | ➖ No necesario | Reusa la del Data table |
| **Faltante**: view models Go de skeleton/inline-alert/validation-summary para el form de alta/edición | ◐ | Solo templates + test view models |

**Veredicto**: cero primitivas nuevas. Todo lo que la recipe necesita existe; el trabajo es 100% wiring (handler + template de pantalla + view models de estado).

### 3.2 Esquema de datos — entidad `Resource`

```text
Resource {
  ID       string   // slug estable, deep-linkable
  Name     string
  Status   string   // vocabulario cerrado: Active | Pending | Done (reuso dataTableStatuses, data_table.go:29)
  Date     string   // ISO-8601 (orden cronológico = orden string, data_table.go:39)
  Owner    string   // opcional
  Tags     []string // para filtro Chips (opcional)
  // vista:
  Selected bool     // derivado de ?selection= (data_table.go:201-225)
}
```

Sort keys cerradas: `name | status | date` (reuso `dataTableSortKeys`, `data_table.go:23`). Set de demo: slice server-side (como `dataTableItems`, `data_table.go:41-54`) con page size 10-15 (el demo usa 5, `data_table.go:26`).

### 3.3 Contrato server

| Operación | Método + ruta | Params / body | Respuesta | Contrato |
|---|---|---|---|---|
| Listar/filtrar/ordenar/paginar | `GET /recipes/admin-resource` | `?q=&sort=&dir=&page=` | página completa (no-JS) o fragmento `resource-panel` (HX) | (c) GET params estables, `data_table.go:302-322` |
| Selección masiva | `GET` con `?selection=` | `selection=all` o IDs | round-trip re-renderiza checked + notice | (c), `data_table.go:184-220` |
| Borrar (individual o bulk) | `POST /recipes/admin-resource/delete` | `id` o `selection[]` | 303 a `GET /recipes/admin-resource` | (d) POST+303; confirmación vía Dialog; feedback transitorio `loom:toast` |
| Crear/editar (form en Dialog o sección) | `POST /recipes/admin-resource/save` | campos del form | 303 (éxito) o 422 + `X-Loom-Validation` + validation-summary + inline-alert por campo (fallo) | (a) 422 + (d) 303 |
| Refresh remoto | `POST /recipes/admin-resource/refresh` | — | fragmento + `HX-Trigger loom:toast` | reuso exacto de `data_table.go:365-385` |
| Error global persistente | — | — | `bannerView{Tone:"error"}` en slot `layout.html:24` | success/error persistente nunca toast |

Fragmento HX: mismo patrón de bifurcación `HX-Request` de `data_table.go:137-147`. Los links llevan `hx-get` + `hx-target="#resource-panel"` + `hx-swap="outerHTML"` (como `data-table.html:8,47-48,75-76`).

### 3.4 Composición — 19 campos (resumida)

| Campo | Valor |
|---|---|
| **SURFACE** | Página de app shell (server-rendered) con appbar + slot banner global; no overlay para el listado (surface grammar: página = unidad de URL, `composition-rules.md:52`) |
| **USER** | Admin/operador que administra un recurso (crea, filtra, selecciona, borra) |
| **PRIMARY_TASK** | Gestionar el set de recursos: encontrar, comparar, seleccionar y actuar por fila |
| **SECONDARY_TASKS** | Filtrar/buscar, ordenar, paginar, crear recurso, borrar individual/bulk |
| **UX_PATTERN** | Resource list (#3 `ux-patterns.md`) + Pagination (#6) + Bulk action (#11) + Destructive (#10) + Confirmation (#18) + Error recovery (#9) |
| **VISUAL_VOCABULARY** | Data table (colección: criterio 4.1 `composition-rules.md:57`), Empty state, Skeleton, Inline alert, Validation summary, Banner, Toast |
| **COMPONENTS** | Data table, Text field (search), Button, Icon button, Menu (row overflow), Dialog (confirm), Toast, Badge (counts), Chips (filtros), Navigation (Tabs/drawer para sub-secciones), Skeleton, Empty state, Inline alert, Validation summary, Banner |
| **STATES** | Rest; **Empty** (`empty-state.html` fila `<td colspan>`, CTA "Clear search"); **Loading** (Skeleton en carga inicial, `.ui-progress` determinate en refresh); **Error** (campo: 422 inline; recurso: `error-state` 404/500; global: Banner); **Selected** (checkbox nativo `:checked`); **Success** (toast transitorio + `banner--success`/`inline-alert--success` persistente `role="status"` post-303) |
| **ACCESSIBILITY** | `aria-sort` en columna activa, `aria-current="page"` en paginación, checkboxes nativos con `aria-label` por fila, Dialog `closedby="any"` + focus trap nativo + Cancel `autofocus`, summary `role="alert"` con links reales a `#campo-error`, error nunca color-only (forced-colors en `data-table.css`) |
| **CONTENT_RULES** | Botones = verbos de acción ("Delete record", no "OK"); empty = qué + CTA real; confirmación = qué se hace + si es irreversible (`content-rules.md:22-41,76-84`); errores = qué pasó + qué hacer; copy en inglés por contrato, localizable |
| **SEO_REQUIREMENTS** | Ruta `/recipes/*` → `noindex, nofollow` (como demos, `server.go:184-186`); canonical limpio sin query (§16 `seo-contract.md:399-414`); title server-driven; un solo `<h1>`; si se documenta en `screen-recipes.md` renderizada vía layout, esa página sí es indexable |
| **GEO_REQUIREMENTS** | Contenido factual (contrato server en los mismos términos del código: 422/`loom:toast`/GET, §11 `geo-contract.md:111-119`); URLs estables; sin claim universal no verificado; texto de estado citado del vocabulario (Gelium UI es la entidad única) |
| **SERVER_CONTRACT** | `GET ?q=&sort=&dir=&page=&selection=`; `HX-Request` bifurca; `POST + 303` para mutaciones; `422 + X-Loom-Validation` para validación; `HX-Trigger loom:toast` para transitorio; banner/inline para persistente |
| **NO_JS_FLOW** | Completo: GET links reales (sort/page/filter), checkbox en form real, POST+303, 422 re-render página completa preservando valores, dialog server-rendered fallback, toast inline no-JS (`toast.go:161-164`) |
| **HTMX_ENHANCEMENT** | Swap del panel `#resource-panel` (`outerHTML`), `loom:toast` para resultados, focus al validation-summary en 422 (enhancement opcional) |
| **RESPONSIVE_BEHAVIOR** | Paginación server-side en vez de scroll horizontal (composición fluida, `composition-rules.md:156-163`); en narrow, List sobre Data table (criterio 4.1); sin breakpoints — grid fluido `minmax/auto-fit` + `min()/clamp()` |
| **THEME_REQUIREMENTS** | Solo tokens `--ui-*` existentes; cero literales de color (guard `TestNoColorLiteralsInComponents`); light/dark + forced-colors + reduced-motion; sin tokens de densidad nuevos |
| **ALTERNATIVES_REJECTED** | Board (es un set, no workflow multi-vía); List para >8-10 filas (anti-regla `composition-rules.md:66`); dialog para el listado completo (página = URL/back, anti-regla "no flujo largo en dialog"); toast para errores de validación (contrato 3.1) |
| **RATIONALE** | Es la recipe que ejercita la mayor cantidad de contratos canónicos (GET params, selección, POST+303 destructivo, Dialog confirm, empty, toast, banner/inline) sobre una primitiva (Data table) ya completa — valida el framework de recipes con el mínimo riesgo |

### 3.5 Gaps específicos (Admin Resource)

1. **View models Go de producción** para Skeleton, Inline alert, Validation summary (hoy solo templates + test models). Requeridos para el form de alta/edición.
2. **Slot Banner sin inyección**: `pageView.Banner` existe y el layout lo renderiza, pero ningún handler lo setea; la recipe es la primera oportunidad de usarlo (ej. banner de mantenimiento/error global).
3. **Bulk delete confirm**: el Dialog existe; falta el wiring de selección→dialog→POST `selection[]`→303 (patrón Bulk action #11).

---

## 4. RECIPE 2 — OPS QUEUE

### 4.1 Componentes: existentes / faltantes

| Necesidad | Estado | Evidencia |
|---|---|---|
| List two-line (fila de queue) | ✅ Existe | `list.html` (two-line `ui-list-item--two-line`) |
| Badge (counts/unread) | ✅ | `badge.go` — **sin tones** |
| Chip (categoría/estado) | ✅ | `chips.go` — sin tones |
| Button (acción "avanzar"/"tomar siguiente") | ✅ | `button.go` (aria-busy) |
| Toast (resultado de avanzar/reasignar) | ✅ | `toast.go` |
| Dialog (confirmar acciones irreversibles) | ✅ | `dialog.html` |
| Empty state | ✅ | `empty-state.html` |
| Skeleton (carga inicial) | ✅ | `skeleton.html` (falta view model Go) |
| Select/Text field (filtros de cola) | ✅ | `select.go`, `text_field.go` |
| **Avatar** (requester) | ✖ **NO existe** — demo ad-hoc `.demo-wa-avatar` (`demo-whatsapp.html:35`, `demo-whatsapp.css:154-166`) | Bloquea la fila tipo queue |
| **Indicador de tono por ítem** (ok/warning/expired, prioridad) | ✖ **NO reusable** — Badge solo error-tinted; tones solo token/demo (`demo-whatsapp.css:217,501-502`) | Bloquea el estado de cola |
| **Acción por fila + badge trailing** en List | ◐ Falta slot en `list.html` (solo leading/trailing icon) | Componible dentro del `<li>` (patrón demo), decidir |
| Pagination standalone (cola larga) | ◐ Solo en Data table (`vocabulary.md:232`) | Depende del volumen; si el slice es corto, "cargar más"/links |

### 4.2 Esquema de datos — entidad `QueueItem`

```text
QueueItem {
  ID          string    // slug
  Subject     string    // headline
  Requester   string    // nombre + initial (avatar)
  Kind        string    // vocabulario cerrado: message | order | support | billing (filtro Chips)
  Status      string    // vocabulario cerrado: new | processing | done | blocked (tone + label, nunca color-only)
  ReceivedAt  time.Time // orden FIFO del server
  SLADeadline time.Time // → tone ok|warning|expired (reuso del patrón ventana WhatsApp, demo_whatsapp.go:298-312)
  Unread      int
  Assignee    string    // opcional
  // vista:
  Active      bool      // ?item= activo (patrón ?c= del demo)
}
```

El **orden es operativo** (criterio 4.2 `composition-rules.md:68-77`): el servidor ordena por prioridad/estado/received; la posición NO es presentacional. Tono = derivación server-side (como `windowStatus`, `demo_whatsapp.go:298-312`), nunca estado cliente.

### 4.3 Contrato server

| Operación | Método + ruta | Params / body | Respuesta | Contrato |
|---|---|---|---|---|
| Listar cola (orden operativo) | `GET /recipes/ops-queue` | `?status=&kind=&page=` (vocabularios cerrados) | página o fragmento `queue-panel` | (c) GET params; sanitización a defaults (patrón `data_table.go:142-157`) |
| Avanzar ítem al siguiente estado | `POST /recipes/ops-queue/advance` | `id` | 303 a la cola (+ `?item=` si aplica) + toast "Marked as processing" | (d) POST+303 + `loom:toast` (o banner/inline success persistente) |
| Reasignar / bloquear / reabrir | `POST /recipes/ops-queue/action` | `id` + `action` (vocabulario cerrado) | 303 + toast | (d) |
| "Tomar siguiente" (batch) | `POST /recipes/ops-queue/next` | — | 303 a `?item=<próximo>` | (d); patrón de cola FIFO |
| Filtros | `GET` con `?status=&kind=` | — | fragmento re-render | (c) |
| Error global | — | — | Banner `role="alert"` (slot layout) | persistente nunca toast |

**Regla de queue**: la acción dominante es "avanzar el siguiente" → se resuelve con form POST + 303 + toast, sin drag ni fragmentos para mutaciones (anti-regla: no board para FIFO, `composition-rules.md:87`).

### 4.4 Composición — 19 campos (resumida)

| Campo | Valor |
|---|---|
| **SURFACE** | App shell full-width con appbar + slot banner; master-detail opcional (cola + panel del ítem activo) como el demo WhatsApp |
| **USER** | Operador/agente que atiende ítems de una cola (triaje secuencial) |
| **PRIMARY_TASK** | Atender el "siguiente a atender": ver estado de cada ítem y avanzarlo al siguiente estado |
| **SECONDARY_TASKS** | Filtrar cola por estado/categoría, tomar el próximo ítem, reasignar, ver SLA (tone) |
| **UX_PATTERN** | Queue (vocabulario `:170-179`, composición List two-line + Badge/Chip tone + Button + Toast; POST+redirect) + Notifications (#15) + Loading (#8) |
| **VISUAL_VOCABULARY** | List two-line, Badge (unread) + indicador de tono (SLA/estado), Button, Toast, Dialog (confirmar bloqueo), Empty state, Skeleton, Inline alert, Banner |
| **COMPONENTS** | List, Badge, Chips, Button, Icon button, Select/Text field (filtros), Dialog, Toast, Progress (opcional), Skeleton, Empty state, Inline alert, Banner, **Avatar (nuevo)** |
| **STATES** | Rest; **Empty** ("No hay ítems pendientes" + CTA "Ver completados"); **Loading** (Skeleton por fila); **Error** (422 en acciones con body inválido, inline; global Banner); **Tone por ítem** (ok/warning/expired según SLA — label + tone, nunca color-only); **Success** (toast transitorio o inline/banner success persistente post-303) |
| **ACCESSIBILITY** | Rows con acción = `<a>`/`<button>` reales; estado por ítem con texto visible + tone (no color-only); `aria-current` en fila activa; Dialog `closedby="any"`; `#loom-toast-region` para transitorio; lista `aria-live="polite"` opcional en swaps |
| **CONTENT_RULES** | Botón de acción = verbo del estado destino ("Mark as done", "Take next"); SLA expirado = qué pasó + acción; empty = mensaje + CTA real; tone con label textual (`content-rules.md:11-24`) |
| **SEO_REQUIREMENTS** | `/recipes/*` → noindex; canonical limpio; un h1; si se indexa la versión documentada, title/description por ruta única |
| **GEO_REQUIREMENTS** | Estados/acciones descritos con los términos del contrato; copy factual; entidad única Gelium UI |
| **SERVER_CONTRACT** | `GET ?status=&kind=&page=`; `POST + 303` para cada transición; `422 + X-Loom-Validation` si un body es inválido; `HX-Trigger loom:toast` para resultados transitorios; banner/inline success persistente para confirmaciones que deben sobrevivir |
| **NO_JS_FLOW** | Fila = `<li>` con estado + form POST por acción; 303 full reload; orden y tono calculados server-side; sin JS en el flujo principal |
| **HTMX_ENHANCEMENT** | Swap `#queue-panel` al filtrar/avanzar; `loom:toast`; `hx-swap-oob` opcional para la fila avanzada; sin inventar contrato nuevo |
| **RESPONSIVE_BEHAVIOR** | Master-detail colapsa a lista sola en narrow (patrón demo); filas full-width; sin breakpoints — contenedor fluido; paginación server-side si el slice crece |
| **THEME_REQUIREMENTS** | Indicador de tono DEBE reusar `--ui-color-{success,warning,info,danger-container}-*` + `-fg` (como `demo-whatsapp.css:217,501-502`); cero literales; light/dark/forced-colors/reduced-motion |
| **ALTERNATIVES_REJECTED** | Board (FIFO estricto = queue, no multi-vía, anti-regla `composition-rules.md:87`); Data table (cada fila es unidad discreta con acción, no comparación columnar, criterio 4.1/4.2); drag & drop (server-first, POST+303); toast para SLA expirado persistente |
| **RATIONALE** | Es la composición de workflow más simple del sistema (List+badge+button+POST+303, ya probada en la sidebar WhatsApp) y demuestra que el orden es estado server-side; desbloquea el patrón Queue para recursos reales |

### 4.5 Gaps específicos (Ops Queue)

1. **Avatar** — componente nuevo (o promoción del `.demo-wa-avatar` a primitiva `ui-avatar`). Necesario para la fila requester. **Bloqueante**.
2. **Indicador de tono reusable** — decisión: variantes `--tone` en Badge/Chip o un "status indicator" nuevo; debe exponer label + tone (nunca color-only) y reusar los tokens semánticos. **Bloqueante**.
3. **Slot trailing en List** (badge + acción por fila) — componible dentro del `<li>`; decidir si se formaliza como slot de `list.html`.
4. **Pagination standalone** — si el set de cola es grande; extraer el `<nav>` de `data-table.html:74-78` como partial reutilizable o usar "load more".

---

## 5. RECIPE 3 — PUBLIC / SOCIAL FEED

### 5.1 Componentes: existentes / faltantes

| Necesidad | Estado | Evidencia |
|---|---|---|
| Card (ítem de feed) | ✅ | `card.html` (article/a/button) |
| List three-line (feed terso) | ✅ | `list.html` (`ui-list-item--three-line`) |
| Badge (likes/comments/unread) | ✅ | `badge.go` (count) |
| Button (reaccionar/compartir/load more) | ✅ | `button.go` |
| Skeleton (estado de carga crítico del feed) | ✅ | `skeleton.html` — falta view model Go |
| Empty state | ✅ | `empty-state.html` |
| Error state / transporte | ✅ | `error-state.html`, `app.js:88-94` |
| Toast (acción transitoria) | ✅ | `toast.go` |
| **Avatar** (autor) | ✖ **NO existe** | Igual que Ops Queue — bloquea la fila/card de feed |
| **Pagination / "load more"** | ◐ No standalone | Vocabulario `:232`; feed orientado a novedad → "load more" o links reales `?page=` |
| Tabs (vistas "Para ti / Siguiendo / Nuevos") | ✅ | `tabs.go` (links reales + `aria-current`) |

### 5.2 Esquema de datos — entidad `FeedItem`

```text
FeedItem {
  ID        string
  Author    string   // nombre + initial (avatar)
  Kind      string   // vocabulario cerrado: text | image | link | announcement
  Body      string   // contenido (línea 1 = headline)
  Timestamp time.Time // orden reverso-cronológico (novedad = valor, criterio 4.5)
  Likes     int
  Comments  int
  New       bool     // marcador de novedad (aria-label, no color-only)
}
```

El ítem **ES el evento** (`vocabulary.md:208-216`); orden invertido; loading/empty críticos. Filtro por vistas = Tabs server-side (`?view=for-you|following|new`, vocabulario cerrado como `dataTableSortKeys`).

### 5.3 Contrato server

| Operación | Método + ruta | Params / body | Respuesta | Contrato |
|---|---|---|---|---|
| Ver feed | `GET /recipes/social-feed` | `?view=&page=` (o `?before=` para load-more) | página o fragmento `feed-panel` | (c) GET params estables |
| Vista Tabs | `GET` con `?view=` | vocabulario cerrado, sanitizado | re-render con `aria-current` | patrón Tabs (`tabs.go`) |
| Cargar más / paginar | `GET ?page=` o `?before=<id>` | — | links reales (`hx-get` target feed-panel) o append | (c); si append HTMX, el no-JS es page relink |
| Reaccionar (like) | `POST /recipes/social-feed/react` | `id` | 303 + toast | (d) + `loom:toast` transitorio |
| Marcar leído / ocultar | `POST /recipes/social-feed/dismiss` | `id` | 303 | (d) |
| Refresh | `POST /recipes/social-feed/refresh` | — | fragmento + toast (patrón `data_table.go:354-389`) | (c)+(b) |

**Regla de feed**: "load more" con HTMX hace append al fragmento; sin JS, el mismo link es una página `?page=N` re-render (URL sigue siendo el estado). El loading inicial es Skeleton server-rendered; nunca un flash de empty antes de los datos.

### 5.4 Composición — 19 campos (resumida)

| Campo | Valor |
|---|---|
| **SURFACE** | Página pública del feed (app shell o standalone público) — orientada a novedad, no a administrar (criterio 4.5 `composition-rules.md:97-105`) |
| **USER** | Usuario público/consumidor que se entera de lo nuevo |
| **PRIMARY_TASK** | Enterarse de la novedad (escaneo rápido, orden cronológico) |
| **SECONDARY_TASKS** | Cambiar de vista (Para ti/Siguiendo/Nuevos), reaccionar, compartir, cargar más |
| **UX_PATTERN** | Feed (vocabulario `:208-216`, Card/List three-line + avatar + Badge; loading/empty críticos) + Loading (#8) + Empty (#7) + Pagination (#6) + Tabs server-side |
| **VISUAL_VOCABULARY** | Card (unidad repetida autocontenida, criterio 4.4) o List three-line; Avatar; Badge (counts); Tabs; Skeleton; Empty state; Toast |
| **COMPONENTS** | Card, List (three-line), Badge, Button, Icon button, Tabs, Toast, Progress (refresh), Skeleton, Empty state, **Avatar (nuevo)** |
| **STATES** | Rest; **Loading** (Skeleton `role="status"` + `aria-busy` en región — crítico en feed); **Empty** ("Nada nuevo aún" + CTA "Explorar" o refresh); **Error** (transporte → toast `app.js:88-94`; recurso → `error-state`); **New** (marcador con label, no color-only); **Success** (toast transitorio para reacciones) |
| **ACCESSIBILITY** | Avatar `aria-hidden` + nombre en texto (nunca solo inicial decorativa); tarjetas linkeables como `<a>` reales; tabs como links con `aria-current`; región feed `aria-live="polite"` opcional para swaps; skeleton con `.sr-only` "Loading" |
| **CONTENT_RULES** | Copy factual y terso; empty con CTA real; timestamps relativos; "like" = verbo de acción; nunca "cargando…" como mensaje (mostrar skeleton/progress, `content-rules.md:64-72`) |
| **SEO_REQUIREMENTS** | `/recipes/*` noindex; canonical limpio; si el feed real fuera público e indexable, JSON-LD `BreadcrumbList`/`Article` y `?page=` excluidos del canonical (§16); un h1 por página |
| **GEO_REQUIREMENTS** | Contenido del feed factual/verificable; entidad única; los ítems de demo sin claims no verificados |
| **SERVER_CONTRACT** | `GET ?view=&page=`; `POST + 303` para reacciones/dismiss; `loom:toast` transitorio; `HX-Request` bifurca fragmento `feed-panel`; empty/loading son output del servidor |
| **NO_JS_FLOW** | Feed completo server-rendered: links reales para vista/página, POST+303 para reacciones, Skeleton reemplazado por el siguiente GET, refresh = reload con toast inline (`toast.go:161-164`) |
| **HTMX_ENHANCEMENT** | Swap `#feed-panel` al cambiar vista; append en "load more"; `loom:toast` para reacciones; `aria-live` anuncia el swap |
| **RESPONSIVE_BEHAVIOR** | Cards en grid fluido `repeat(auto-fit, minmax(14rem, 1fr))` (patrón `card.css:30`) o lista single-column en narrow; sin breakpoints; load-more/paginación server-side |
| **THEME_REQUIREMENTS** | Avatar y badges con tokens existentes; cero literales de color; light/dark/forced-colors/reduced-motion (shimmer del skeleton apagado con `prefers-reduced-motion`) |
| **ALTERNATIVES_REJECTED** | Collection/Data table (el feed es evento/tiempo, no entidad/set con filtros — criterio 4.5); Timeline (no hay eje temporal explícito de proceso, es actividad reciente tersa — criterio 4.6); Board; spinner ad-hoc en vez de Skeleton (anti-regla `composition-rules.md:130`) |
| **RATIONALE** | Es la recipe que obliga a resolver loading + novedad + avatar/paginación; su DoD valida los patrones de estado de carga del sistema (los más críticos según `composition-audit.md:225`) |

### 5.5 Gaps específicos (Public/Social Feed)

1. **Avatar** — mismo gap de Ops Queue. **Bloqueante** (autor del ítem).
2. **Pagination/"load more"** — decidir el patrón (page relink no-JS + append HTMX). El vocabulario lo marca pendiente (`vocabulary.md:232`).
3. **View model Go de Skeleton** — para el estado de carga inicial server-rendered.
4. **Marcador "New"** — requiere label accesible (no color-only); Badge count sirve como unread.

---

## 6. GAPS BLOQUEANTES TRANSVERSALES (antes de la recipe)

| Gap | Recipes que bloquea | Estado | Acción propuesta |
|---|---|---|---|
| **Avatar** (inicial del usuario) | Ops Queue, Feed | ✖ No existe (ad-hoc demo) | Nueva primitiva `ui-avatar` (o promoción del `.demo-wa-avatar`, `demo-whatsapp.css:154-166`) con `aria-hidden` + nombre textual; variante `--sm`/`--lg` |
| **Indicador de tono reusable** (estado/SLA con label + tone) | Ops Queue | ◐ Solo tokens + demo | Variantes tone en Badge/Chip reusando `--ui-color-{success,warning,info,danger-container}-*` + `-fg`; label siempre presente |
| **Pagination standalone** | Ops Queue (si cola larga), Feed | ◐ Solo en Data table (`vocabulary.md:232`) | Extraer `<nav aria-label="Table pages">` (`data-table.html:74-78`) como partial `pagination` reutilizable |
| **View models Go de Skeleton / Inline alert / Validation summary / Callout** | Las 3 (wiring server-driven) | ◐ Solo templates + test models | Agregar tipos Go de producción (como `emptyStateView`/`bannerView`); sin esto no hay render server-driven limpio |
| **Slot trailing de List** (badge + acción por fila) | Ops Queue | ◐ Componible | Decidir: formalizar slot en `list.html` o componer dentro del `<li>` (patrón demo WhatsApp) |
| **Naming header validación** (`X-Loom-Validation` vs `X-Gelium-Validation`) | Las 3 (formularios) | ⚠ Discrepancia docs/código | Usar el valor real del código; unificar docs en una pass de naming (no bloquea) |

**NO bloqueantes** (ya entregados): Empty state, Skeleton template, Inline alert template, Validation summary template, Banner+slot, Error state+404, Success persistente por reuso, transporte HTMX G5, Data table con empty. **No hay necesidad de Steps, Breadcrumbs, Date picker ni Callout** en las 3 recipes (quedan para Detail/Editor/Booking/Auth/Settings diferidas).

---

## 7. ORDEN DE IMPLEMENTACIÓN

```text
0. Prerrequisitos compartidos (slice previo a recipes):
   1) Avatar (primitiva)                     — desbloquea Ops Queue y Feed
   2) Tone variants de Badge/Chip            — desbloquea Ops Queue
   3) Partial Pagination standalone          — desbloquea Ops Queue/Feed
   4) View models Go de Skeleton/Inline alert/Validation summary — wiring de las 3

1. Admin Resource   — PRIMERO. Cero primitivas nuevas (todo existe: Data table + empty + dialog
                      + toast + banner). Ejercita la mayor cantidad de contratos canónicos
                      (GET params, selección, POST+303 destructivo, Dialog confirm, 422, toast,
                      banner/inline persistente) sobre la primitiva más madura del sistema.
                      Valida el marco de recipes (template de pantalla + handler + fragmento HX)
                      con riesgo mínimo y establece el patrón que reusarán las otras dos.

2. Ops Queue        — SEGUNDO. Consume Avatar y tone variants (ya construidos en el paso 0).
                      Establece el patrón de workflow (POST+303 por transición, orden operativo
                      server-side) sobre la composición ya probada en la sidebar WhatsApp.

3. Public/Social Feed — TERCERO. Consume Avatar + pagination. Exige loading/skeleton y la
                      decisión load-more vs page; al llegar al final, los prerrequisitos y los
                      dos patrones de pantalla previos ya resolvieron la mayoría de las decisiones.
```

Racional de orden: **dependencia de primitivas y madurez** (Admin no depende de nada nuevo; Queue y Feed dependen de Avatar/tone/pagination). Alternativa rechazada: Feed primero (más visible) — su dependencia de Avatar/load-more y su exigencia de loading lo hacen el más caro sin los prerrequisitos.

---

## 8. DÓNDE VIVIRÍAN (nuevos templates, handlers, rutas)

Basado en cómo viven los demos actuales (`demo_whatsapp.go` + `demo-whatsapp.html` como template de pantalla standalone, registrado en `server.go:300-306`; ejemplos fragmentados en `/examples/*`):

- **Documento**: `docs/gelium-ui-screen-recipes.md` — el entregable de Phase G (`roadmap.md:249`), 3 recipes con los 19 campos, anclado a los componentes/contratos reales.
- **Templates de pantalla** (standalone, como `demo-whatsapp.html` — documento completo con `<head>` + metadata): `web/templates/recipe-admin-resource.html`, `web/templates/recipe-ops-queue.html`, `web/templates/recipe-social-feed.html`. Reutilizan partials (`data-table-panel`, `empty-state`, `skeleton`, `inline-alert`, `banner`, `toast-region`, `validation-summary`) ya existentes.
- **Handlers** (nuevo archivo por recipe, patrón `demo_whatsapp.go`): `internal/app/recipe_admin_resource.go`, `recipe_ops_queue.go`, `recipe_social_feed.go`, con mock stores server-side (mutex + slice/slice ordenado como `whatsAppStore`).
- **Rutas**: registrar en `server.go` `New()` junto al mux (`server.go:285-317`):
  - `GET /recipes/admin-resource`, `GET /recipes/ops-queue`, `GET /recipes/social-feed` (noindex como demos, `server.go:184-186`);
  - POST actions (`/recipes/admin-resource/{delete,save,refresh}`, `/recipes/ops-queue/{advance,action,next}`, `/recipes/social-feed/{react,dismiss,refresh}`) en el mux + lista `postOnlyPaths()` (`server.go:329-342`) para el 405 companion.
  - Enlazar desde `/docs` (sección Demos de `docs.go:95-96`) y desde `screen-recipes.md`.
- **Fragmentos HTMX**: partials tipo `data-table-panel` (`data-table.html:34-80`) por recipe (`resource-panel`, `queue-panel`, `feed-panel`) para el swap `outerHTML`.

---

## 9. ARCHIVOS IMPACTADOS (solo read-only — para implementación futura)

**Docs**
- `docs/gelium-ui-screen-recipes.md` — NUEVO (Phase G, 3 recipes × 19 campos).
- `docs/gelium-ui-system-roadmap.md` — corregir conteo de campos (lista 18 / DoD 17 / encargo 19) y nota de naming de header (opcional).

**Código Go (nuevo)**
- `internal/app/recipe_admin_resource.go`, `recipe_ops_queue.go`, `recipe_social_feed.go` — handlers + mock stores + view models de recipe.
- `internal/app/state_views.go` (o similar) — view models Go de producción para Skeleton, Inline alert, Validation summary (hoy solo en tests).

**Código Go (cambios mínimos)**
- `internal/app/server.go` — registro de rutas en `New()`, `postOnlyPaths()`, uso del slot `Banner`/`Error`.
- `internal/app/routes.go` — (opcional) si las recipes entran en el nav; mejor enlazarlas desde `/docs`.
- `internal/app/docs.go` — sección "Recipes" en el índice.

**Templates (nuevo)**
- `web/templates/recipe-{admin-resource,ops-queue,social-feed}.html` (+ partials de panel por recipe).
- `web/templates/avatar.html` (primitiva nueva); opcional `pagination.html` (extracción).

**Templates (cambios mínimos)**
- `web/templates/list.html` — slot trailing (badge/acción) si se formaliza.

**CSS (nuevo)**
- `web/styles/avatar.css`; variantes tone en `web/styles/badge.css` (o `status.css`); `web/styles/recipe-*.css` para el layout de pantalla (scaffolding, como `demo-whatsapp.css`).

**CSS (cambios mínimos)**
- `web/styles/app.css` — imports; `web/styles_contract_test.go` — lista `sourceAppCSS` (`styles_contract_test.go:24-65`) y guards de contrato (cero literales de color, etc.).

**Tests**
- `internal/app/server_test.go` — render tests de las 3 screens (contratos de los 19 campos que importan: roles, `aria-*`, fragmentos, 303, 422); `web/styles_*_test.go` para Avatar/tone.

**No tocados**: componentes Material existentes, contratos server (`toast`, `text_field`, `data_table`), `web/static/app.js`.

---

## 10. FUENTES DE AUTORIDAD

`docs/gelium-ui-system-roadmap.md` (Phase G :216-251, matriz :434-437), `docs/gelium-ui-composition-rules.md` (:28-54, :57-135, :164-190), `docs/gelium-ui-ux-patterns.md` (#3 Resource list, #6 Pagination, #7 Empty, #8 Loading, #10 Destructive, #11 Bulk, #18 Confirmation), `docs/gelium-ui-vocabulary.md` (:92-166, :170-217, :232-237), `docs/gelium-ui-content-rules.md`, `docs/gelium-ui-seo-contract.md`, `docs/gelium-ui-geo-contract.md`, `docs/handoffs/{composition-audit,state-patterns-audit,success-feedback-audit}.md`, `internal/app/{server,routes,data_table,demo_whatsapp,toast,text_field,list,badge,docs}.go`, `web/templates/{layout,data-table,demo-whatsapp,demo-whatsapp-admin,list,badge,empty-state,skeleton,inline-alert,banner,callout,error-state,validation-summary,toast}.html`, `web/styles_contract_test.go`, `web/styles/{badge,list,demo-whatsapp,tokens}.css`, `web/static/app.js`, `git log` (commits Phase D `eba1c4c`..`94340dc`).
