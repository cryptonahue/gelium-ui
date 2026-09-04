# Vocabulary Audit — Gelium UI (ex Gelium UI)

> Handoff read-only de investigación (Phase 2 del system roadmap). No modifica código, templates, CSS ni tests. Única salida: este documento.
> Alcance: mapear el vocabulario visual objetivo (`docs/gelium-ui-system-roadmap.md:127`) contra el estado real de Gelium para alimentar `docs/gelium-ui-vocabulary.md`.
> Baseline leído: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `lib/skills/14-component-implementation.md`, `internal/app/*.go`, `web/templates/*.html` (todos), `web/content/*.md` (muestra), `web/static/app.js`, `themes/theme-material/theme.css`, handoffs previos (`core-audit.md`, `theme-architecture-audit.md`, `composition-audit.md`).
> Referencia de nomenclatura: NameThatUI (solo reconocimiento visual; NO autoridad normativa de semántica/accesibilidad).

---

## 1. Resumen ejecutivo

- Gelium entrega **26 rutas de componente** + demos, todos con semántica HTML nativa, tokens `--ui-*`, cero JS de componente y HTMX solo como enhancement (`internal/app/routes.go:18-45`).
- Del vocabulario objetivo (29 términos, `gelium-ui-system-roadmap.md:127`), **8 ya existen** como componentes entregados: Card, List, Data table, Tabs, Dialog, Drawer (navigation drawer), Tooltip, Toast (Menu es componente real y cubre parcialmente el término "Popover", pero no es término del vocabulario). **8 parciales** (Pagination, Multi-select, Inline alert, Empty state, Loading state, Form, Popover, Detail view) y **13 sin equivalente ni plan**: Panel, Queue, Feed, Timeline, Board, Steps, Breadcrumbs, Combobox, Date picker, Banner, Callout, Skeleton, Resource editor.
- La mayor brecha estructural: el vocabulario objetivo cubre **patrones de estado de pantalla y de workflow** (Empty, Loading/Skeleton, Banner, Callout, Steps, Queue, Breadcrumbs) que Gelium hoy **no tiene**; la auditoría de composición ya los marcó como bloqueantes de Phase 4 (`composition-audit.md:230`).
- Semántica HTML base ya sólida: `<article>/<a>/<button>` (Card), `<ul>/<ol>/<nav>` (List, Tabs, Menus, Drawer), `<table>` + `<caption>` + `<th scope>` (Data table), `<dialog>` (Dialog, Drawer modal, Select menu), `role="tooltip"` (Tooltip), `aria-live` (Toast). Esto es el cimiento del vocabulario; no hay que inventar raíces.
- JS por término: **todos los entregados son zero-JS o JS-optativo**. Los únicos JS/HTMX reales son Toast (auto-dismiss + `loom:toast`), Slider (fill WebKit), Data table refresh (HTMX) y los `<dialog>`/Popover que dependen de primitives declarativas (`command`, `popovertarget`, `closedby`) con fallback server-rendered.
- Tres conflictos de naming dominantes para resolver en Phase 2: **Popover vs Menu**, **Multi-select (ambigüedad en 4 implementaciones)** y **Drawer vs Navigation drawer** (ver §3).
- Datos densos de referencia: los términos de status/feedback transitorio ya están bien delimitados por contrato (Toast: `loom:toast` con vocabulario cerrado `info|success|warning|error`, `internal/app/toast.go:45`); los que faltan (Inline alert, Banner, Callout) son **persistentes**, no transitorios, y no compiten con Toast.

---

## 2. Mapeo término → estado en Gelium

Leyenda: **✅ Implementado** (ruta real) · **◐ Parcial** (existe como patrón/ad-hoc, no como componente con contrato) · **✖ No existe** (ni implementación ni plan).

| # | Término vocabulario | Estado | Evidencia Gelium (ruta/template/handler) |
|---|---|---|---|
| 1 | Card | ✅ | `/components/card` · `card.html` · `card.go` · tokens `--ui-card-*` (`theme.css`) · `<article>/<a>/<button>` |
| 2 | Panel | ✖ | Sin componente. Solo contenedores ad-hoc: `tabs-demo-panel` (`tabs.html:19`), `data-table-panel`, secciones `demo-wa-admin-panel` (`demo-whatsapp-admin.html:30,51,70,103,115`) |
| 3 | List | ✅ | `/components/list` · `list.html` · `list.go` · `<ul>/<ol>/<nav>`; one/two/three-line, nav, selection, icons (`list.md:3-41`) |
| 4 | Data table | ✅ | `/components/data-table` · `data-table.html` · `data_table.go` · `<table>/<caption>/<th scope>`; sort/filter/pagination/selection server-side por GET (`data-table.html:34-74`) |
| 5 | Queue | ✖ | Sin componente. Patrón real ad-hoc: sidebar del demo WhatsApp (avatar + unread + window tone ok/warning/expired, `demo-whatsapp.html:26-55`) |
| 6 | Feed | ✖ | Sin componente. Patrón real ad-hoc: thread de mensajes (`demo-whatsapp.html:76-130`) y grid de Cards (`card.html:2-21`) |
| 7 | Timeline | ✖ | Sin componente ni plan |
| 8 | Board | ✖ | Sin componente ni plan |
| 9 | Steps | ✖ | Sin componente ni plan (composition-audit lo marca gap de Resource Editor/Booking, `composition-audit.md:55,60`) |
| 10 | Tabs | ✅ | `/components/tabs` · `tabs.html` · `tabs.go` · `<nav>`+`<ul>`+`<a>` con `aria-current="page"`, activo server-side, sin tablist falso (`tabs.md:53-59`) |
| 11 | Pagination | ◐ | Solo dentro de Data table: `ui-data-table-pagination` prev/next + páginas numeradas (`data-table.html:68-72`), slice 5 (`data_table.go:26`). No es componente standalone |
| 12 | Breadcrumbs | ✖ | Sin componente ni plan |
| 13 | Multi-select | ◐ | Patrón distribuido en 4 componentes: List selection con checkboxes (`list.html:70-104`), Chips filter (`chips.html:16-35`), Segmented multi (`segmented-button.html:47-85`), filas de tabla (`data-table.html:56-60`). Sin widget "multi-select" unificado |
| 14 | Combobox | ✖ | Sin equivalente. El más cercano: Select menu (listbox en `<dialog>` con `command="show-modal"`, `select.html:73-92`), pero NO tiene typeahead/filtrado |
| 15 | Date picker | ✖ | Sin componente ni plan |
| 16 | Dialog | ✅ | `/components/dialog` · `dialog.html` · `dialog.go` · `<dialog>` + Invoker Commands declarativos + `closedby="any"` (`README.md:120`) |
| 17 | Popover | ◐ | No es componente: es la **primitiva** detrás de Menu (`popover` + `popovertarget`, `menu.html:44,48`) y fue auditada y rechazada para Tooltip (`tooltip.md:43-44`). Conflicto de naming con Menu (§3.1) |
| 18 | Drawer | ✅ | `/components/navigation-drawer` · `navigation-drawer.html` · `navigation_drawer.go` · modal = `<dialog>`, permanente = `<nav>` (`navigation-drawer.md:50-57`) |
| 19 | Tooltip | ✅ | `/components/tooltip` · `tooltip.html` · `tooltip.go` · `role="tooltip"` + `aria-describedby`, reveal por `:hover`/`:focus-within` CSS, plain/rich (`tooltip.md:31-44`) |
| 20 | Toast | ✅ | `/components/toast` · `toast.html` · `toast.go` · `HX-Trigger loom:toast`, región `aria-live` (`toast.html:9-11`), fallback inline no-JS |
| 21 | Inline alert | ◐ | No es componente. Patrón real: error de campo con `role="alert"` (`text-field.html:8`, `select.html:89`) y toast inline no-JS (`toast.html:21`). Falta componente genérico persistente |
| 22 | Banner | ✖ | Sin componente ni plan (composition-audit gap §5.2: error state de página) |
| 23 | Callout | ✖ | Sin componente ni plan |
| 24 | Empty state | ◐ | Único ejemplo real no-reusable: `demo-wa-empty` (`demo-whatsapp.html:51-53`). Data table con 0 filas muestra caption `0 rows` sin mensaje/CTA (`data_table.go:239`) |
| 25 | Loading state | ◐ | Cubierto puntualmente: Button loading `aria-busy` (`button.html:4,9`), Progress determinate/indeterminate (`progress.html:5-24`), refresh con Progress+Toast (`data-table.html:76-85`). Sin estado de carga de pantalla ni Skeleton |
| 26 | Skeleton | ✖ | Sin componente ni plan |
| 27 | Detail view | ◐ | Sin componente. Patrón real ad-hoc: master-detail del chat WhatsApp (header window bar + thread + composer, `demo-whatsapp.html:57-158`) |
| 28 | Form | ◐ | No es componente (correcto: es patrón nativo). Formularios reales en todos los demos (`validation-form` `text-field.html:12-16`, selection list `list.html:72`, table `data-table.html:14-23`, menu `menu.html:94-125`). `Field` es primitive interna, no publicable (`COMPONENT-ROADMAP.md:33`) |
| 29 | Resource editor | ✖ | Sin componente ni plan (screen recipe Phase 4, `composition-audit.md:55`) |

**Resumen**: 8 ✅ · 8 ◐ (Pagination, Multi-select, Popover, Inline alert, Empty state, Loading state, Detail view, Form) · 13 ✖ (Panel, Queue, Feed, Timeline, Board, Steps, Breadcrumbs, Combobox, Date picker, Banner, Callout, Skeleton, Resource editor).

---

## 3. Aliases y conflictos de naming

### 3.1 Alias existentes en Gelium (evidencia file:line)

| Alias (término alternativo) | Canónico real en Gelium | Evidencia / estado |
|---|---|---|
| Snackbar | Toast | Decisión explícita: no crear Snackbar separado, usar como referencia visual (`COMPONENT-ROADMAP.md:255`, `MATERIAL-WEB-PROGRESS.md:159`) |
| Navigation drawer | Drawer (vocab.) / Navigation drawer (Material) | Componente entregado como "navigation drawer" (`docs.go:68`) |
| Tab / Navigation tab | Tabs | "Navigation tab" reusa contrato Tabs/Nav bar (`COMPONENT-ROADMAP.md:138`) |
| Selection list | List (selection variant) | Implementado como checkboxes en `<li>` (`list.md:33`) |
| Menu popup / dropdown | Menu | Implementado sobre Popover API (`menu.html:44-125`) |
| Popover | Menu (hoy) | Primitiva popover usada por Menu; término del vocabulario pendiente |
| Table | Data table | Gelium solo tiene Data table (Labs). Tabla nativa simple existe ad-hoc en admin (`demo-whatsapp-admin.html:33-47`) |
| Alert (ARIA) | Inline error / Toast error | `role="alert"` ya se usa para error de campo (`text-field.html:8`) y toast error (`toast.go:56-61`); el término "Inline alert" compite con ese rol |
| Snackbar anatomy en test | Toast | `TestToastSourceCSSImplementsSnackbarAnatomyAndAccessibleStates` (`web/styles_toast_test.go:50`) — el alias sobrevive en el nombre del test |

### 3.2 Los 3 conflictos de naming más importantes

1. **Popover vs Menu**. "Popover" es a la vez un término del vocabulario objetivo (`gelium-ui-system-roadmap.md:127`) y el nombre de la primitiva web (`popover`/`popovertarget`). Hoy **Menu es la única superficie popover entregada** y Tooltip la descartó por no-Baseline (`tooltip.md:43-44`). Si Phase 2 define "Popover" como patrón canónico separado de Menu, choca con: (a) el nombre de la API nativa, (b) el componente Menu ya entregado, (c) la regla "no imitar con CSS". **Propuesta**: el vocabulario NO debe listar "Popover" como patrón de UI; debe ser un *mecanismo* (top layer/overlay) que Menu y futuro Context menu usan, documentado como tal.

2. **Multi-select: 4 implementaciones, 1 término**. El término del vocabulario "Multi-select" hoy está materializado en List (checkboxes), Chips filter, Segmented multi-select y selección de filas en Data table. Son patrones distintos (listado vs filtro vs control segmentado vs grid) con la misma etiqueta. **Propuesta**: "Multi-select" debe definirse como *capacidad de selección múltiple sobre un patrón huésped* (List, Data table, Menus) con regla de elección, no como componente. El único candidato a widget es el Combobox multi (gap, §6).

3. **Drawer vs Navigation drawer**. El componente entregado se llama "Navigation drawer" (Material upstream, `docs.go:68`), el vocabulario objetivo dice "Drawer". Además coexisten dos variantes con semántica distinta: modal (`<dialog>`) y permanente (`<nav>`) (`navigation-drawer.html:8-19`). **Propuesta**: mantener "Navigation drawer" como nombre canónico de Gelium (evita colisión genérica con "drawer" de overlays tipo bottom-sheet no-Material) y alias "Drawer" → Navigation drawer.

Otros conflictos menores a resolver en Phase 2: **Alert** (rol ARIA vs componente Inline alert vs Banner), **Table** vs **Data table**, **Select** vs **Select menu** vs **Combobox**, **Queue** vs **List two-line** (evidencia demo WhatsApp), **Steps** vs **Pagination** (ambos con paginación/avance, intenciones distintas).

---

## 4. Distinciones canónicas para pares críticos (intención/estructura)

Basado en `composition-audit.md` §4 (criterios ya propuestos) y validado contra los componentes reales:

### 4.1 Card vs Panel
- **Card** (`card.html`): unidad **autocontenida, repetida e interactiva** (`<article>`, `<a>`, `<button>`); cada instancia tiene su propia entrada/acción.
- **Panel**: región **única de layout** que agrupa contexto de una página; estática, no repetida, sin entrada propia (bloques admin `demo-whatsapp-admin.html:30-127`).
- Criterio: ¿es un elemento repetido de un set con acción propia (card) o un bloque contextual de la página (panel)?

### 4.2 List vs Queue
- **List** (`list.html`): catálogo/índice; el **orden no es semántica de workflow**; no hay "próximo a atender".
- **Queue**: workflow con **orden y posición operativos** (FIFO/pending) + estado por ítem (badge/tone) + acción para avanzar al siguiente estado. Evidencia: sidebar WhatsApp (unread + window tone, `demo-whatsapp.html:43-46`).
- Criterio: ¿hay un siguiente a procesar y estados que transicionar? → Queue (List two-line + Badge + Button + Toast).

### 4.3 Queue vs Board
- **Queue**: **una sola vía ordenada**; la posición es la señal primaria.
- **Board**: **múltiples vías paralelas** donde mover ítems entre estados es la tarea (en este stack: form POST + redirect + toast, sin drag; columna = `<section>` con Cards + form por ítem).
- Criterio: ¿flujo en una secuencia (queue) o transferencia libre entre estados (board)? ¿"tomar el siguiente" o "mover de columna"?

### 4.4 Feed vs Collection
- **Feed**: orientado a **evento/tiempo**, orden cronológico, el ítem ES el evento; la novedad es el valor; loading/empty críticos.
- **Collection**: orientado a **entidad/set**; catálogo curado, filtrable/ordenable/paginable; el set es el producto (evidencia: Data table, search).
- Criterio: ¿enterarme de lo nuevo (feed) o encontrar/administrar el set (collection)?

### 4.5 Timeline vs Activity list
- **Timeline**: narrativa cronológica con **eje temporal explícito** (hitos + fechas + contexto); orden ascendente es semántico (audit trail, booking history).
- **Activity list**: agregación de eventos recientes en orden inverso, filas tersas (List one/two-line), sin eje formal.
- Criterio: ¿reconstruir un proceso con fechas (timeline, necesita markers) o ver actividad reciente (activity list, List ya sirve)?

### 4.6 Table vs List
- **Table** (`data-table.html`): datos densos **columnares** que se **comparan** entre columnas; set grande/remoto con sort+filter+pagination server-side; fila = registro homogéneo.
- **List** (`list.html`): unidades verticales discretas con acción primaria y jerarquía 1-3 líneas; sets pequeños/estáticos. Regla: no usar tabla para ≤5 filas estáticas (`composition-audit.md:263`).
- Criterio: ¿comparar atributos en paralelo o escanear unidades? ¿set remoto con sort/filter (tabla) o índice estático (lista)?

### 4.7 Dialog vs Page
- **Dialog** (`dialog.html`): tarea **corta, enfocada, reversible** con contexto retenido (confirmar, quick edit, picker); no cambia la URL.
- **Page**: tarea **profunda/completa**; dirección URL, back, deep-link (booking de pasos = page/steps, no modal).
- Criterio: ¿se resuelve en una interacción sin perder contexto (dialog) o necesita superficie completa + URL (page)? Anti-regla: no abrir flujo largo en Dialog (`composition-audit.md:269`).

### 4.8 Toast vs Inline alert vs Banner vs Callout
- **Toast** (`toast.go`): feedback **transitorio del resultado de una acción**; auto-dismiss, `aria-live`; vocabulario `info|success|warning|error`; NUNCA para validación de campos (`toast.go:129-133`).
- **Inline alert**: mensaje **persistente ligado al contexto** del formulario/sección; `role="alert"` junto al campo (hoy `text-field.html:8`, `select.html:89`). Sobrevive a la interacción.
- **Banner**: aviso **persistente a nivel página/sitio** que exige acción (sesión expirada, mantenimiento); rol alert/status, sin auto-dismiss.
- **Callout**: contenido **informativo/promocional** sin urgencia ni requisito de acción (contexto, tips); se puede ignorar.
- Criterio: ¿efímero resultado (toast) vs error persistente ligado al campo (inline) vs aviso global que exige acción (banner) vs nota contextual ignorable (callout)? Anti-regla: no usar toast para feedback persistente/crítico (`composition-audit.md:266`).

---

## 5. Anatomía y semántica HTML de los términos existentes (base del vocabulario)

| Término | Raíz semántica | Anatomía interna (clases `ui-*`) | Evidencia |
|---|---|---|---|
| Card | `<article>` / `<a>` / `<button>` (según acción) | `.ui-card-title`, `.ui-card-body` | `card.html:4-18` |
| List | `<ul>` / `<ol>` / `<nav>` → `<li>` | `.ui-list-item`, `-headline`, `-supporting`, `-icon` (leading/trailing), `-link`, `-label` (checkbox) | `list.html:6-120` |
| Data table | `<table>` + `<caption>` + `<thead>` + `<th scope="col">` + `<tbody>` + `<tr>/<td>` | `.ui-data-table-cell`, `-checkbox`, `-sortable`, `-pagination` (`<nav>`) | `data-table.html:35-72` |
| Tabs | `<nav>` → `<ul>` → `<li>` → `<a href>` | `.ui-tab`, `-icon`, `-label`, `-indicator`; activo con `aria-current="page"` | `tabs.html:6-18` |
| Dialog | `<dialog closedby="any">` | `-headline` (h2 + `aria-labelledby`), `-content` (`aria-describedby`), `-actions` | `dialog.html:3-10` |
| Navigation drawer | modal: `<dialog>`; permanente: `<nav>` → `<ul>` → `<a href>` | `-list`, `-item`, `-destination` (+`--active`), `-indicator`, `-glyph`, `-label` | `navigation-drawer.html:8-38` |
| Navigation bar | `<nav>` → `<ul>` → `<a href>` | `-destination`, `-icon`, `-label`, `-indicator`, reusa `.ui-badge` | `navigation-bar.html:6-21` |
| Menu | `<ul>` (con `popover`) → `<li>` | `-item`, `-item-button` (`<button>`), `-item-link` (`<a>`), `-item-label` (checkbox/radio), `-divider` (`role="separator"`) | `menu.html:6-125` |
| Toast | `<div>` con `role="status"|"alert"` + región `#loom-toast-region aria-live="polite"` | `.ui-toast-{type}`, `-message`, `-action` | `toast.html:2-11` |
| Tooltip | `role="tooltip"` en `<span>` + `aria-describedby` en el control | `.ui-tooltip` (+`--rich`, `--top`), `-subhead`, `-supporting-text`, `-action` | `tooltip.html:10-66` |
| Text field | `<label>` + `<input>`/`<textarea>`; error `<p role="alert">` + `aria-invalid` | `.ui-text-field-control`, `-error-icon`, `-message` | `text-field.html:2-9` |
| Checkbox/Radio/Switch | `<label>` + `<input type="checkbox|radio">`; radio con `<fieldset>`+`<legend>` | `.ui-checkbox/radio/switch-*` mark/track/handle | `checkbox.html`, `radio.html`, `switch.html` |
| Slider | `<div>` + `<input type="range">` + `--ui-slider-fill` | `.ui-slider` | `slider.html` |
| Select / Select menu | `<select>` nativo; menu: `<dialog>` con `<button>` opciones | `.ui-select`, `.ui-select-menu-item` | `select.html:5-92` |
| Chips | `<button>`/`<a>`/`<label>`+checkbox | `.ui-chip-assist/filter/suggestion/input`, `-remove` (submit) | `chips.html:6-63` |
| Segmented | `<fieldset>`+`<legend>` (select) o `role="group"` (acciones); radio/checkbox/button | `.ui-segmented-button-*` | `segmented-button.html:7-108` |
| Button | `<button>`/`<a>` (+ `command`/`commandfor`, `aria-busy`) | `.ui-button-{variant}`, `-spinner` | `button.html:1-10` |
| Badge | `<span>` + `aria-hidden="true"` (nunca color-only) | `.ui-badge` (+`-large`) | `badge.html:4-14` |
| Progress | `<progress>` nativo (value/max o indeterminado) | `.ui-progress` | `progress.html:4-25` |
| Divider | `<hr>` | `.ui-divider` (+`-inset`, `-inset-start`, `-inset-end`) | `divider.html:5-18` |
| FAB / Icon button | `<button>`/`<a>` con nombre accesible obligatorio | `.ui-fab-*`, `.ui-icon-button-*` (+`aria-pressed` toggle) | `fab.html:1-23`, `icon-button.html:1-26` |

Reglas transversales del vocabulario (derivadas de la implementación): elementos nativos antes que ARIA; sin `div`/`span` como controles; SVG decorativo con `aria-hidden` + `focusable="false"`; estado nunca color-only; foco en el control nativo; root semántico elegido por la acción (article/a/button, ul/ol/nav/table/dialog).

---

## 6. Necesidad de JavaScript por término

Leyenda: **0 = cero JS** (nativo/declarativo server-rendered) · **H = HTMX enhancement** · **J = JS mínimo del framework** · **J* = JS + fallback no-JS real** · **D = depende de primitiva declarativa** (requiere navegador compatible + fallback server-rendered).

| Término | JS | Detalle / evidencia |
|---|---|---|
| Card, List, Tabs, Divider, Badge, Elevation, Focus ring, Icon, Button, Icon button, FAB, Chips, Checkbox, Radio, Switch, Progress, Segmented, Text field, Select | 0 | Nativo: `:checked`, `:focus-visible`, forms GET/POST, `progress` (`list.md:66`, `switch.md:3`, `select.md:62`) |
| Data table (base) | 0 | sort/filter/pagination/selection por GET con params estables (`data_table.go:302-322`); sin JS la página recarga |
| Data table (refresh) | H | `hx-post` + `hx-target` swap del panel + `loom:toast`; fallback full-page + toast inline (`data-table.html:76-85`) |
| Menu | D | Popover API declarativa (`popover` + `popovertarget`), light-dismiss/Escape nativos, zero JS (`menu.html:44-125`) |
| Dialog | D | `<dialog>` + Invoker Commands (`command`/`commandfor`), `closedby="any"`; en navegadores previos el trigger no hace nada → fallback server-rendered del consumidor (`README.md:120`, `dialog.md`) |
| Navigation drawer (modal) | D | sobre `<dialog>` nativo + invokers; variante permanente = `<nav>` sin JS (`navigation-drawer.md:54-57,101`) |
| Tooltip | 0 | reveal `:hover`/`:focus-within` CSS; Interest Invokers rechazados por no-Baseline (`tooltip.md:36-44`) |
| Select menu | D | `<dialog>` + `command="show-modal"`/`request-close` + submit server (`select.html:77-87`) |
| Toast | J* | Sin JS: toast inline persistente completo (`toast.html:21`); con JS: región `aria-live` + auto-dismiss pausable + `loom:toast` (`app.js:11-77`) |
| Slider | J | Solo fill en WebKit (`--ui-slider-fill`, `app.js:79-101`); input nativo operable sin JS |
| Pagination (en tabla) | 0 / H | Links reales + `hx-get` opcional (`data-table.html:69-71`) |

**Conclusión JS**: ningún término entregado exige JS para el flujo principal. Los términos faltantes que SÍ necesitarán JS o enhancement HTMX (por naturaleza): **Combobox** (typeahead/filtrado en cliente o server-driven), **Date picker** (puede resolverse en parte con `<input type="date">` nativo + calendario server-driven), **Board** (mover ítems entre vías = POST + redirect + toast, sin drag obligatorio), **Steps/Resource editor** (avance server-driven con validación 422 por paso). Skeleton, Empty, Banner, Callout, Breadcrumbs, Panel, Detail view son 100% estáticos.

---

## 7. Gaps: términos sin equivalente ni plan en Gelium

| Término | Estado | Nota para Phase 2 / roadmap de componentes |
|---|---|---|
| **Combobox** | ✖ | El más complejo de los faltantes. Select menu es listbox en `<dialog>`, no combobox con filtrado. Requiere decidir: typeahead server-driven (GET) vs JS mínimo. No hay plan |
| **Date picker** | ✖ | Nada. Referencia nativa: `<input type="date">`; la agenda calendario es gap real. Gap en Booking Flow (`composition-audit.md:60`) |
| **Steps** | ✖ | Gap en Resource Editor y Booking Flow (`composition-audit.md:55,60`). Avance server-driven + validación 422 por paso + Progress determinate |
| **Breadcrumbs** | ✖ | No hay nada (solo nav del layout `<header>`, `layout.html:14`). 100% estático, `<nav aria-label="Breadcrumb">` + `<ol>` |
| **Pagination** (standalone) | ◐ | Solo dentro de Data table (`data-table.html:68-72`). Extraer como contrato reusable para listados sin tabla |
| **Skeleton** | ✖ | No hay nada; solo Progress y loading de botón. Requisito de Feed/Dashboard (`composition-audit.md:57-58`) |
| **Banner** | ✖ | Error/success persistente a nivel página/sitio; requisito Auth y Settings (`composition-audit.md:59,62`) |
| **Callout** | ✖ | Nota informativa persistente ignorable; 100% estático |
| **Panel** | ✖ | Contenedor de sección estático (contraste con Card, §4.1); puede definirse solo en vocabulario sin componente propio (regla de composición) |
| **Queue / Feed / Timeline / Board** | ✖ | Patrones de workflow; dependen de List/Card/Badge/Progress existentes. Definir en vocabulario como *composiciones*, no componentes nuevos (evidencia demo WhatsApp) |
| **Detail view** | ✖ | Patrón master-detail (chat WhatsApp). Definir como screen recipe (Phase 4), no componente |
| **Resource editor** | ✖ | Screen recipe (Phase 4), no componente; requiere Steps + Inline alert |
| **Empty state / Loading state / Inline alert** | ◐ | Existen como ad-hoc; convertir en componentes reusables (requisito bloqueante Phase 4, `composition-audit.md:230`) |

---

## 8. Recomendaciones para Phase 2 (`docs/gelium-ui-vocabulary.md`)

1. **Resolver los 3 conflictos de naming de §3.2** antes de fijar nombres canónicos: Popover como mecanismo (no patrón), Multi-select como capacidad, Drawer alias de Navigation drawer.
2. **El vocabulario debe distinguir 3 capas** (ya implícitas en Gelium): *patrones de datos* (Card, Panel, List, Data table, Table), *patrones de estado* (Empty, Loading/Skeleton, Inline alert, Banner, Callout, Toast) y *patrones de workflow* (Queue, Board, Steps, Timeline, Feed) — cada capa con reglas de composición propias.
3. **Usar las raíces HTML de §5 como contrato canónico**; ningún término debe inventar raíz (article/ul/table/nav/dialog ya cubren el 90%).
4. **Los términos faltantes se dividen en**: (a) componentes nuevos que el roadmap de componentes debe planificar (Skeleton, Inline alert, Banner, Callout, Empty state, Steps, Breadcrumbs, Pagination standalone), y (b) patrones que se definen solo en vocabulario/composición sin componente (Panel, Queue, Feed, Timeline, Board, Detail view, Form, Resource editor, Multi-select, Combobox como combinación de Select menu + filtrado).
5. **Combobox y Date picker requieren auditoría platform-first antes de comprometer JS** (regla `lib/skills/14-component-implementation.md`, platform-first): `<input type="date">`, `<datalist>` y GET con autocomplete son los candidatos nativos.

**Fuentes de autoridad usadas**: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `lib/skills/14-component-implementation.md`, `internal/app/{routes,server,toast,data_table,docs}.go`, `web/templates/*.html` (todos), `web/content/{card,list,data-table,dialog,toast,tooltip,tabs,navigation-drawer,menu,segmented-button,switch,slider,radio,select,chips,index}.md`, `web/static/app.js`, `themes/theme-material/theme.css`, `docs/handoffs/{core-audit,theme-architecture-audit,composition-audit}.md`.
