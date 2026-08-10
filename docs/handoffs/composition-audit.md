# Gelium UI — Composition Rules Audit (handoff para Phase 3)

> **Alcance**: investigación read-only sobre el estado real del catálogo para que `docs/gelium-ui-composition-rules.md` (Phase 3) y `docs/gelium-ui-screen-recipes.md` (Phase 4) reutilicen lo que ya existe. No modifica código, templates, CSS ni tests.
>
> **Baseline**: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `internal/app/*`, `web/templates/*`, `web/styles/*`, `web/static/app.js`.

---

## 1. Estado real de composición

### 1.1 Inventario de componentes entregados (ruta real → template → handler)

Registro de rutas: `internal/app/routes.go:16-47`; handlers montados en `internal/app/server.go:67-89`.

| Capa | Componente | Template | Handler | Evidencia extra |
|---|---|---|---|---|
| Foundation | Elevation, Focus ring, Icon, Divider | `elevation.html`, `focus-ring.html`, `icon.html`, `divider.html` | `elevation.go`, `focus_ring.go`, `icon.go`, `divider.go` | tokens en `themes/theme-material/theme.css` |
| Actions | Button | `button.html` | `button.go` | disabled/loading/aria-busy, link vs button (`README.md:107-116`) |
| Actions | Icon button | `icon-button.html` | `icon_button.go` | toggle `aria-pressed` |
| Actions | FAB | `fab.html` | `fab.go` | extended label / icon-only con AriaLabel obligatorio |
| Actions | Chips | `chips.html` | `chips.go` | input chip con remoción server-side (`chips.go:62-89`) |
| Actions | Segmented buttons | `segmented-button.html` | `segmented_button.go` | radio/checkbox/button nativos, fieldset |
| Actions | Menu | `menu.html` | `menu.go` | Popover API declarativa, items nativos, zero JS |
| Input | Text field | `text-field.html` | `text_field.go` | 422 + `X-Gelium-Validation` |
| Input | Checkbox / Radio / Switch / Select / Slider | `checkbox.html`, `radio.html`, `switch.html`, `select.html`, `slider.html` | `checkbox.go`, `radio.go`, `switch.go`, `select.go`, `slider.go` | controles nativos estilizados |
| Input | List | `list.html` | `list.go` | one/two/three-line, nav/selection/static |
| Feedback | Dialog | `dialog.html` | `dialog.go` | `<dialog>` + Invoker Commands declarativos |
| Feedback | Toast | `toast.html` | `toast.go` | `loom:toast`, región aria-live, fallback inline no-JS |
| Feedback | Progress | `progress.html` | `progress.go` | determinate/indeterminate |
| Feedback | Badge | `badge.html` | `badge.go` | no color-only |
| Feedback | Card | `card.html` | `card.go` | article/a/button |
| Feedback | Tooltip | `tooltip.html` | `tooltip.go` | Interest Invokers declarativos, zero JS |
| Navigation | Tabs | `tabs.html` | `tabs.go` | links reales + `aria-current`, sin tablist falso |
| Navigation | Navigation bar / tab / drawer | `navigation-bar.html`, `navigation-tab.html`, `navigation-drawer.html` | `navigation_bar.go`, `navigation_tab.go`, `navigation_drawer.go` | activo server-side, drawer modal/permanente |
| Data | Data table | `data-table.html` | `data_table.go` | sort/filter/pagination server-side por GET |

**28 rutas de componente** + `/docs` (índice generado, `docs.go:83-98`) + demos.

### 1.2 La única "screen real" que existe hoy: el demo WhatsApp

`internal/app/demo_whatsapp.go` + `web/templates/demo-whatsapp.html` / `demo-whatsapp-admin.html` es la referencia más cercana a una screen recipe y debe usarse como evidencia de composición:

- **Master-detail** (chat): sidebar = lista tipo queue (avatar, nombre, última hora, unread badge, window tone ok/warning/expired) + panel de detalle (header con window bar, thread de mensajes, composer). `demo-whatsapp.html:25-160`.
- **Search por GET** `?q=` y selección activa por `?c=` — estado del detalle derivado de la URL (`demo_whatsapp.go:491-536`), no de estado cliente.
- **Empty state** del listado filtrado: `demo-whatsapp.html:51-53` (único empty state real del sistema).
- **Admin** = sección de tablas nativas (sin contrato Data table) + forms + token card + rate bar: `demo-whatsapp-admin.html:21-128`.
- Contrato de acciones: POST + **303 SeeOther redirect** (full reload) en `demo_whatsapp.go:555,569,574,584` — patrón "server-driven simple", sin fragmentos HTMX.

### 1.3 Mapeo de screen recipes objetivo → componentes existentes (Phase 4)

| Receta (Phase 4) | Componentes que ya existen | Gaps reales |
|---|---|---|
| Admin Resource | Data table (sort/filter/pagination/selection), Button, Icon button, Menu, Dialog, Toast, Badge, Text field (search), Chips | Empty state de tabla; confirmar borrado → dialog existe |
| Resource Detail | Card, List (read-only), Divider, Tabs, Badge, Button, Menu (overflow), Dialog, Tooltip | Detail view / meta list con label+valor (vocabulario pendiente) |
| Resource Editor | Text field, Select, Checkbox, Radio, Switch, Slider, Segmented buttons, Chips, Button, contrato 422 | Pasos de wizard (Steps) |
| Ops Queue | List two-line + Badge + Chip + Button + Toast; POST+redirect del demo WhatsApp | Queue como patrón propio (vocabulario pendiente) |
| Public/Social Feed | Card, List three-line, avatar inicial (demo), Badge, Button | **Loading state** de feed; **skeleton**; empty state de feed |
| Dashboard | Card (KPIs), Progress, Data table (recent), Tabs, Badge, Divider | Skeleton; KPI card con delta |
| Settings Page | List (nav), Text field, Select, Switch, Radio, Checkbox, Slider, Divider, Tabs, Button, Toast | Inline alert/banner para éxito persistente |
| Booking Flow | Text field, Select, Radio/Segmented, Progress (determinate), Dialog (confirm), Toast | **Steps**, **Date picker** |
| Search Results | Text field (search GET), Data table o List, Chips (filters), Empty state (demo WhatsApp) | Empty state genérico reusable |
| Authentication Flow | Text field, Button, Checkbox, Radio, Dialog, Toast, contrato 422 | **Inline alert/banner** para errores de credenciales; page-level error |

**Conclusión 1.3**: el catálogo cubre los *patrones de acción y de datos* (wave 2–6 del roadmap), pero los patrones **de estado de pantalla** (Empty state, Loading/Skeleton, Inline alert, Banner, Callout, Steps, Breadcrumbs) están en el vocabulario de Phase 2 (`gelium-ui-system-roadmap.md:127`) y **no tienen componente entregado**. Toda pantalla compuesta va a necesitarlos.

---

## 2. Data table: sort/filter/pagination server-side (patrón de referencia)

### 2.1 Implementación

`internal/app/data_table.go` + `web/templates/data-table.html` + `web/styles/data-table.css`.

- Semántica nativa `<table>/<thead>/<th scope>/<caption>`; header 56px, filas 52px, padding 16px (`web/styles/data-table.css:14-16`, `data-table.md:29`).
- **Sort**: headers son `<a href>` reales que llevan `?sort=&dir=`; columna activa con `aria-sort="ascending|descending"` (`data-table.html:45-50`, `data_table.go:269-290`). La columna activa togglea dirección; las otras linkean asc.
- **Filter**: `<form method="get">` real con `?q=...` que filtra por name/status case-insensitive (`data-table.html:7-12`, `data_table.go:160-165`).
- **Pagination**: links reales `?page=2`, página actual como `<span aria-current="page">` (`data-table.html:68-72`, `data_table.go:223-230`).
- **Selection**: checkboxes nativos en form real; `selection=all` para la página (`data-table.html:40-60`, `data_table.go:184-208`); round-trip re-renderiza checked + notice (`data_table.go:326-345`).
- **Refresh remoto**: POST `/examples/data-table/refresh` con `.ui-progress` determinate + toast inline (no-JS) o fragmento + `HX-Trigger loom:toast` (HTMX) (`data_table.go:354-389`, `data-table.html:76-85`).

### 2.2 Estados que cubre

- Rest / hover (state layer) / focus (`:focus-visible` en link o checkbox) / pressed — `data-table.css:82-99`.
- Selected (derivado de `:checked`, nunca color-only) — `data-table.css:104-107`.
- Disabled: explícitamente NO parte del contrato de la tabla (`data-table.md:50`); cada control interno mantiene su propio disabled.
- Paginación prev/next deshabilitados por `aria-disabled` — `data-table.html:69-71`.
- **Empty**: NO cubierto. Con `q=zzz` la caption renderiza `0 rows · page 1 of 1` y `<tbody>` vacío (`data_table.go:239`), sin mensaje ni CTA. Es el gap más visible.

### 2.3 Patrón de interacción server-driven (a reutilizar por las composition rules)

1. **Todo estado navegable es una URL**: `?q=&sort=&dir=&page=` con orden estable y `url.QueryEscape` (`data_table.go:302-322`). La URL es el estado; no-JS = reload completo; HTMX = el mismo link lleva `hx-get` con target `#data-table-panel` + `outerHTML` (`data-table.html:8,48,69-70`).
2. **Vocabularios cerrados**: sort keys `name|status|date` (`data_table.go:23`) y statuses `Active|Pending|Done` (`data_table.go:29`) sanitizados contra defaults (`data_table.go:142-157`). Nada inventado por el cliente sobrevive.
3. **HX-Request bifurca en el handler**: con `HX-Request: true` devuelve SOLO el fragmento `data-table-panel` (sin documento completo, sin forms externos); sin la header devuelve página completa (`data_table.go:120-130`). Test: `data_table_test.go:135-154` y `156-166`.
4. **Los links preservan el estado**: al cambiar sort se conserva `q`; al paginar se conserva `q`+`sort`+`dir` (`data_table.go:269-289`, test `data_table_test.go:285-298`).
5. **Operaciones remotas** (refresh) reusan `.ui-progress` + `.ui-toast` + `loom:toast` (`data_table.go:365-385`), nunca un spinner ad-hoc.

---

## 3. Contratos server-driven reales (evidencia file:line)

### 3.1 HTTP 422 + `X-Gelium-Validation: true` (validación server-side)

| Evidencia | Archivo:línea |
|---|---|
| 422 en texto vacío (Text field) | `internal/app/text_field.go:64-68` |
| Header en rama HX únicamente | `internal/app/text_field.go:87-89` |
| 422 en mensaje vacío (Toast) | `internal/app/toast.go:147-150` |
| Header en rama HX (Toast) | `internal/app/toast.go:177-179` |
| 422 en valor desconocido (Select menu) | `internal/app/select.go:94-96` |
| Header en rama HX (Select menu) | `internal/app/select.go:112-114` |
| Hook `htmx:beforeSwap` (swapea SOLO 422 con header) | `web/static/app.js:1-9` |
| Contrato documentado | `web/content/text-field.md:19` |

Reglas del contrato:
- **422 ≠ error de transporte**: el hook `app.js:1-9` sólo hace `shouldSwap=true` cuando status=422 Y header=`true`. Otros 422 siguen siendo errores.
- Sin HX el servidor re-renderiza la **página completa con status 422** (`text_field.go:76-78`), preservando valor, `aria-invalid="true"` y mensaje (`text-field.html:5,8`).
- Validación **nunca dispara toast**: `toast.go:129-133` ("Validation failures are never reported as toasts") y `COMPONENT-ROADMAP.md:49`.

### 3.2 `HX-Trigger: {"loom:toast":{...}}` (feedback server-driven)

| Evidencia | Archivo:línea |
|---|---|
| Wire contract `{"loom:toast":{"type":"success","message":"Saved"}}` | `internal/app/toast.go:108-127` |
| Seteo de la header en rama HX | `internal/app/toast.go:154-160` |
| Reuso desde Data table refresh | `internal/app/data_table.go:365-370` |
| Listener `loom:toast` → `makeToast` | `web/static/app.js:69-76` |
| Región `#loom-toast-region` con `aria-live="polite"` | `web/templates/toast.html:9-11` |
| Región incluida en layout y demos | `web/templates/layout.html:51`, `demo-whatsapp.html:162`, `demo-whatsapp-admin.html:130` |
| Vocabulario cerrado `info|success|warning|error` | `internal/app/toast.go:45`; `error`→`role="alert"` (`toast.go:56-61`) |
| Auto-dismiss 4s / 8s para error, pausable | `web/static/app.js:17-18,61-65` |
| Fallback no-JS: toast inline persistente | `toast.html:21`, `toast.go:161-164` |

### 3.3 Data table GET params (`q`, `sort`, `dir`, `page`, `selection`)

| Evidencia | Archivo:línea |
|---|---|
| Lectura de los 5 params | `internal/app/data_table.go:118` |
| Vocabularios cerrados (sort keys / statuses) | `data_table.go:23,29` |
| Sanitización a defaults | `data_table.go:142-157` (test `data_table_test.go:246-268`) |
| Orden estable de params + escape | `data_table.go:302-322` |
| Fragmento HX exclusivo | `data_table.go:120-130` |
| Round-trip de selección | `data_table.go:184-220` |
| Preservación de estado en links | `data_table.go:269-289` |

### 3.4 Contratos menores existentes

| Contrato | Evidencia |
|---|---|
| Chips remove: POST `/examples/chips/remove` + re-render full page + notice inline (sin rama HX) | `internal/app/chips.go:66-89` |
| Tabs server-side: `?tab=` / `?sub=` validados contra vocabulario cerrado, activo con `aria-current="page"` | `internal/app/tabs.go` (patrón), `web/content/tabs.md:53-58` |
| WhatsApp actions: POST + 303 redirect | `internal/app/demo_whatsapp.go:555,569,574,584` |
| Regla de enriquecimiento local: 422+header swap | `web/static/app.js:1-9` |

**Conclusión 3**: los composition rules DEBEN reusar estos tres contratos sin inventar otros: (a) 422+header para validación de campos/valores, (b) `loom:toast` para feedback transitorio de operaciones, (c) GET con parámetros estables para estado de listados (sort/filter/page/selection). Para "mover estados" de workflow el patrón ya probado es POST + redirect (WhatsApp), no fragmentos.

---

## 4. Diferenciaciones críticas (criterios propuestos)

Criterios de decisión basados en lo que existe y en la tarea primaria. Vocabulario referencia: `gelium-ui-system-roadmap.md:127`.

### 4.1 table vs list
- **Table** (evidencia `data-table.md:86`): datos densos y **columnares** donde se **compara** entre columnas, el set es grande/remoto y necesita **sort + filter + pagination server-side**; cada fila es un registro con atributos homogéneos.
- **List** (`list.md:41`): contenido vertical continuo donde **cada fila es una unidad discreta** con acción primaria (navegar, seleccionar) y jerarquía de 1-3 líneas; sets pequeños/estáticos, sin comparación entre columnas.
- **Criterio**: ¿el usuario compara atributos en paralelo o escanea unidades? ¿el set supera la paginación visible y requiere sort/filter? ¿los atributos son homogéneos (→ tabla) o heterogéneos con énfasis (→ lista)?

### 4.2 list vs queue
- **List**: catálogo/índice; el **orden no es semántica de workflow**; sin noción de "próximo ítem a procesar".
- **Queue**: workflow con **orden y posición operativos** (FIFO/pending), estado de cada ítem (badge/tone) y **acción para avanzar** el ítem al siguiente estado. Evidencia del sistema: la sidebar del demo WhatsApp es una cola (unread + window tone ok/warning/expired, `demo-whatsapp.html:43-46`).
- **Criterio**: ¿hay un "siguiente a atender" y estados que transicionar? ¿el orden importa al usuario? → queue (List two-line + Badge/Chip + Button + toast). Si es referencia/índice → list.

### 4.3 queue vs board
- **Queue**: **una sola vía ordenada** (dispatch FIFO, triaje secuencial); la posición es la señal primaria.
- **Board**: **múltiples vías paralelas** donde mover ítems **entre estados** es la tarea; requiere comparación entre vías y acción de movimiento por ítem (en este stack: form POST + redirect + toast, sin drag).
- **Criterio**: ¿el trabajo fluye en una sola secuencia (queue) o hay varios estados paralelos con transferencia libre (board)? ¿la acción dominante es "tomar el siguiente" o "mover a otra columna"? Un board sin drag server-side = columna = `<section>` con Cards + form por ítem.

### 4.4 card vs panel
- **Card** (`card.html`): **unidad autocontenida e interactiva** (`<article>`, `<a>`, `<button>`) para grids/feeds que se escanean; cada card es un elemento repetido.
- **Panel**: **contenedor de sección de una página** que agrupa contenido relacionado, típicamente estático y no repetido (bloque de configuración, token card del admin `demo-whatsapp-admin.html:117-126`).
- **Criterio**: ¿es una instancia repetida de un set (card) o una región única de layout (panel)? ¿tiene su propia acción de entrada (card) o es contexto de la página (panel)?

### 4.5 feed vs collection
- **Feed**: **orientado a evento y tiempo** (orden cronológico, actualización frecuente, el ítem ES el evento); importa la novedad, loading/empty states críticos.
- **Collection**: **orientado a entidad y set** (catálogo curado, filtrable/ordenable/paginable; el set es el producto). Evidencia: search results / admin tables.
- **Criterio**: ¿el valor está en lo que llegó nuevo (feed) o en la totalidad del set con filtros (collection)? ¿la tarea es "enterarme" o "encontrar/administrar"?

### 4.6 timeline vs activity list
- **Timeline**: **narrativa cronológica con eje temporal explícito** (hitos, fechas, contexto de proceso: booking history, audit trail); orden ascendente es semántico.
- **Activity list**: **agregación de eventos recientes** en orden inverso, filas tersas (feed ligero), sin eje formal.
- **Criterio**: ¿el usuario reconstruye un proceso con fechas (timeline) o ve actividad reciente (activity list)? Timeline necesita markers + fechas; activity list es una List one/two-line.

### 4.7 dialog vs page
- **Dialog** (`dialog.html`): tarea **corta, enfocada y reversible** con contexto retenido (confirmar, quick edit, picker); no cambia la URL.
- **Page**: tarea **profunda o de contexto completo** (recurso completo, edición larga); dirección URL, back, deep-link. Evidencia: el admin usa páginas/secciones; el confirm borrado usa dialog.
- **Criterio**: ¿la tarea se resuelve en una interacción sin perder contexto (dialog) o necesita superficie completa, URL y navegación atrás (page)? Booking de varios pasos = page (o steps), confirmación final = dialog.

### 4.8 toast vs inline alert vs banner
- **Toast** (`toast.go:13-14,45`): feedback **transitorio del resultado de una acción**, no bloquea, auto-dismiss, región aria-live, vocabulario info/success/warning/error. NUNCA para validación de campos (contrato 3.1).
- **Inline alert**: error/mensaje **persistente ligado al contexto del formulario o sección** (422 + `role="alert"` junto al campo, `text-field.html:8`).
- **Banner**: aviso **persistente a nivel página/sitio** (sesión expirada, mantenimiento, consent) que exige acción, con rol alert/status y sin auto-dismiss.
- **Criterio**: ¿efímero resultado (toast) vs error que debe sobrevivir a la interacción junto al campo (inline) vs aviso global persistente que bloquea o advierte a nivel sitio (banner)? Hoy el sistema sólo tiene toast + inline (error de campo). Banner y callout = gap Phase 2.

---

## 5. Estados: inventario y gaps

### 5.1 Estados cubiertos por componentes entregados

| Estado | Dónde (evidencia) |
|---|---|
| disabled | `text_field.go:39-41`; `switch.html:21,29`; `slider.html:17`; `select.html:24,53`; `radio.html:21,28`; `list.html:94`; `menu.html:26,63`; `segmented-button.html:74,103`; `icon-button.html`, `fab.html` (aria-disabled + tabindex=-1) |
| error (campo) | `text_field.go:38` + `aria-invalid` + `role="alert"` (`text-field.html:5,8`); `select.html:32,61`; `radio.html:35`; contrato 422 |
| loading (botón) | `button.html` aria-busy + nombre dinámico `Loading {Label}` (`README.md:113`) |
| loading (operación) | Progress determinate en refresh `data-table.html:81`; indeterminate (`progress.go`) |
| selected | checkbox list (`list.html:74-104`), data table (`data-table.html:40-60`), segmented, tabs `aria-current` |
| hover/focus/pressed | state layers `--ui-state-*` (`theme.css:192-196`), `data-table.css:82-99`, `list.css:64-66` |
| empty (listado) | ÚNICO: `demo-whatsapp.html:51-53` (no reusable) |
| transient feedback | toast inline + región (`toast.html:1-11`) |

### 5.2 Gaps para pantallas completas

| Estado de pantalla | Gap | Impacto |
|---|---|---|
| **Empty state** de tabla/listado | Data table con 0 filas renderiza caption `0 rows` y tbody vacío (`data_table.go:239`) sin mensaje/CTA | Toda pantalla de admin/search va a mostrar vacíos sin guía |
| **Loading state** / Skeleton | No existe Skeleton; sólo Progress y loading de botón | Feeds y dashboards no tienen estado de carga de datos (sólo `progress` en refresh puntual) |
| **Error state** de página/recurso | Sin banner ni callout; sólo error de campo y toast transitorio | Auth flow y resource detail no pueden mostrar error persistente global |
| **Success persistente** | Toast es transitorio; no hay banner/inline alert de éxito | Settings/editor necesitan confirmación no efímera |
| **Offline / no-data** | Sin contrato | — |

**Recomendación 5**: los nuevos patrones de estado (Empty, Loading/Skeleton, Inline alert, Banner, Callout, Steps) son **requisito bloqueante de Phase 4**; sin ellos las screen recipes quedan incompletas.

---

## 6. Responsive: evidencia existente y reglas inferibles

### 6.1 Qué existe

- **No hay breakpoints de ancho** en ningún CSS de componente. Los únicos media queries son `prefers-reduced-motion` (`app.css:52-69`) y `forced-colors` (`app.css:71-213`, más uno por componente).
- **Layouts fluidos por grid**, no reflow por query:
  - Cards: `grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr))` (`card.css:30`).
  - Elevation demo: `repeat(3, minmax(0,1fr))` (`elevation.css:15`); text-field preview: `repeat(2, minmax(0,1fr))` (`base.css:55`).
- **Anchos fluidos por `min()`/`clamp()`/`calc()`**:
  - Dialog: `max-width: min(560px, calc(100% - 48px))`, `min-width: 280px` (`dialog.css:3-5`).
  - Drawer modal: `max-width: calc(100vw - 56px)` (`navigation-drawer.css:67`).
  - Menu: `max-width: calc(100vw - 2rem)` (`menu.css:53`).
  - Tooltip: `max-width: min(18rem, calc(100vw - 2rem))` (`tooltip.css:26`).
  - Header: `padding: 1rem clamp(1rem, 4vw, 4rem)` (`base.css:21`); shell: `width: min(68rem, calc(100% - 2rem))` (`base.css:28`).
- **Data table**: filas de altura fija con ellipsis en celdas label (`data-table.css:62-69`); el sistema confía en **paginación server-side** para no desbordar (slice de 5, `data_table.go:26`), no en scroll horizontal/container query.
- **Sin tokens de densidad**: el theme expone `--ui-state-*` opacidades (`theme.css:192-196`) pero **no** `--ui-density-*` ni escala de spacing. El roadmap de sistema lista "Density (compacta, cómoda, holgada)" y "Responsive behavior (breakpoints…)" como trabajo pendiente de Phase 1 (`gelium-ui-system-roadmap.md:101,111`).
- **Drawer** resuelve responsive por variantes (modal vs permanente, `navigation-drawer.md:50-56`), no por media query.

### 6.2 Reglas inferibles para las composition rules

1. **La responsividad actual es "fluida por grid, no reflow por breakpoint"**: las recetas deben componer con `auto-fit/minmax` y `min()/clamp()` antes de introducir breakpoints.
2. **Densidad es responsabilidad del theme/Phase 1, no de las recetas**: no definir densidades por pantalla hasta que existan tokens `--ui-density-*`.
3. **Tablas grandes se resuelven server-side** (paginación), nunca con scroll horizontal de componente: la receta Admin Resource debe heredar slice + sort/filter de GET.
4. **Overlays (dialog/menu/tooltip/drawer) ya son fluidos** (`calc(100vw - n)`) y pueden usarse sin breakpoints en recetas.

---

## 7. Anti-rules propuestas (qué NO hacer, con justificación)

1. **No usar Data table para ≤ 5 filas estáticas** → List. La tabla justifica costo de sort/filter/pagination y comparación columnar; un set pequeño estático es ruido (regla ya enunciada en `data-table.md:86`). Umbral propuesto: >8-10 filas con sort/filter reales y set remoto.
2. **No usar Board para colas FIFO estrictas** → Queue (List + badges + botón "siguiente"). Un board presupone movimiento libre entre vías; una cola estricta tiene un único próximo paso y el board agrega fricción visual sin semántica.
3. **No reportar errores de validación como Toast** → 422 + inline `role="alert"`. Contrato vigente: "validation never announces toast" (`toast.go:129-133`, `COMPONENT-ROADMAP.md:49`). El toast es para resultado de acciones, no para correcciones de campos.
4. **No usar Toast para feedback persistente/crítico** → inline alert o banner. El toast auto-dismissa en 4s/8s (`app.js:17-18`); la información que debe sobrevivir a la interacción se pierde.
5. **No inventar estados de listado fuera de la URL** (sort/filter/page/selection custom en cliente) → reusar GET con params estables y vocabularios cerrados (`data_table.go:302-322`). Sin URL = sin no-JS, sin deep-link, sin back.
6. **No reclamar `role="tablist"`/roving focus sin teclado completo** → links reales con `aria-current` como Tabs (`tabs.md:58`). Misma regla para menu/drawer.
7. **No abrir un flujo largo en Dialog** → page/steps. Dialog es para tareas cortas y reversibles (4.7); un editor o booking de pasos en modal pierde URL/back y escala mal.
8. **No usar un spinner cuando existe Progress determinate/indeterminate** → `.ui-progress` + `aria-busy`/`role="status"`. El refresh demo ya es el patrón (`data_table.go:365-385`).
9. **No dejar empty states sin mensaje/CTA** → todo listado server-side debe definir vacío (hoy la tabla muestra `0 rows` sin guía, `data_table.go:239`).
10. **No usar color como único portador de estado** → estado siempre en el control nativo (`:checked`, `aria-sort`, `aria-current`, `role`) y forced-colors (`data-table.css:339-359`).
11. **No introducir JS para lo que un form GET ya resuelve** (chips remove, selección, sort) → platform-first (`gelium-ui-system-roadmap.md:330`).

---

## 8. Entregables sugeridos para Phase 3/4 (a partir de esta auditoría)

- `docs/gelium-ui-composition-rules.md` (Phase 3): reglas 4.1-4.8 + sección 7 como anti-rules, ancladas a los contratos de la sección 3.
- `docs/gelium-ui-screen-recipes.md` (Phase 4): tabla 1.3 como punto de partida; cada receta debe declarar los componentes de estado faltantes (sección 5.2) como dependencia.
- Antes de Phase 4: decidir si Empty state / Skeleton / Inline alert / Banner / Steps entran al roadmap de componentes (hoy son vocabulario Phase 2 sin implementación).

**Fuentes de autoridad usadas**: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `internal/app/{routes,server,data_table,data_table_test,toast,text_field,select,chips,demo_whatsapp,docs}.go`, `web/templates/{data-table,list,card,toast,layout,demo-whatsapp,demo-whatsapp-admin}.html`, `web/styles/{app,base,data-table,list,card}.css`, `web/static/app.js`, `web/content/{data-table,list,toast,text-field,tabs}.md`, `themes/theme-material/theme.css`.
