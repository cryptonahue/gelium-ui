# Gelium UI — Visual Vocabulary

> Vocabulario canónico del sistema Gelium UI (antes Gelium UI).
> Fase 2 del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/handoffs/vocabulary-audit.md` y `docs/handoffs/composition-audit.md`.
> Referencia de nomenclatura: NameThatUI (reconocimiento visual). NO es autoridad normativa de semántica o accesibilidad.

---

## 1. Cómo leer este vocabulario

Cada término define: nombre canónico, aliases, intención, anatomía, semántica HTML, estados, accesibilidad, cuándo usarlo, cuándo no usarlo, mapping a Gelium, necesidad de JavaScript, relación con patterns.

El vocabulario se organiza en **tres capas** (implícitas en la implementación real):

1. **Patrones de datos** — muestran/agrupan entidades: Card, Panel, List, Data table, Table, Detail view.
2. **Patrones de estado** — comunican condición: Empty state, Loading state, Skeleton, Inline alert, Banner, Callout, Toast, Success feedback (persistente).
3. **Patrones de workflow** — representan procesos y transiciones: Queue, Board, Steps, Timeline, Feed.

Además hay **mecanismos** (no patrones): Popover (primitiva web), Multi-select (capacidad sobre un patrón huésped), Form (patrón nativo).

Leyenda de estado en Gelium: **✅ implementado** · **◐ parcial** (ad-hoc) · **✖ no existe**.

---

## 2. Patrones de datos

### Card ✅

- **Aliases**: tile, tarjeta.
- **Intención**: unidad autocontenida, repetida e interactiva de un set que se escanea (grids, feeds, KPIs).
- **Anatomía**: `.ui-card-title`, `.ui-card-body`; elevación elevada/rellena/contorneada.
- **Semántica HTML**: `<article>` para contenido; `<a>` o `<button>` según la acción; nunca `<div>` interactivo.
- **Estados**: rest, hover, focus, pressed (state layers); no tiene disabled propio.
- **Accesibilidad**: si el card completo es clickeable, el control interno (a/button) es el foco; no card como tabstop.
- **Cuándo usarlo**: instancia repetida de un set con entrada/acción propia; comparación rápida entre ítems.
- **Cuándo no**: región única de layout estática → Panel.
- **Mapping a Gelium**: `web/templates/card.html`, `/components/card`, tokens `--ui-card-*`.
- **JS**: 0.
- **Relación con patterns**: base de Grid/Feed/Board; KPI card en Dashboard.

### Panel ✖ (definir en composición, sin componente propio)

- **Aliases**: section, region, block.
- **Intención**: región única de layout que agrupa contexto de una página; estática, no repetida, sin entrada propia.
- **Semántica HTML**: `<section>` con encabezado cuando tiene heading propio; `<aside>` para contexto complementario.
- **Cuándo usarlo**: bloque de configuración, token card del admin, agrupación de contenido relacionado.
- **Cuándo no**: instancia repetida con acción propia → Card.
- **Mapping a Gelium**: hoy solo contenedores ad-hoc (`demo-whatsapp-admin.html:30-127`).
- **JS**: 0.

### List ✅

- **Aliases**: list item, one/two/three-line, selection list.
- **Intención**: contenido vertical continuo donde cada fila es una unidad discreta con acción primaria y jerarquía de 1-3 líneas.
- **Anatomía**: `.ui-list-item` con headline/supporting, iconos leading/trailing, opcional checkbox/radio; variantes navegación, selección, estático.
- **Semántica HTML**: `<ul>`/`<ol>` (o `<nav>` para navegación) → `<li>`.
- **Estados**: rest, hover, focus, pressed, selected (derivado de `:checked`), disabled.
- **Accesibilidad**: selección con controles nativos; nunca color-only.
- **Cuándo usarlo**: catálogo/índice de unidades discretas; sets pequeños/estáticos; el orden NO es semántica de workflow.
- **Cuándo no**: datos densos columnares a comparar → Data table; hay "siguiente a atender" → Queue.
- **Mapping a Gelium**: `web/templates/list.html`, `/components/list`.
- **JS**: 0.
- **Relación con patterns**: base de Queue, Feed, Menu, Navigation, Detail view.

### Data table ✅

- **Aliases**: table, grid de datos.
- **Intención**: datos densos **columnares** que se **comparan** entre columnas; set grande/remoto con sort + filter + pagination server-side.
- **Anatomía**: `.ui-data-table-cell`, checkbox de selección, headers sortables, pagination `<nav>`.
- **Semántica HTML**: `<table>` + `<caption>` + `<thead>` + `<th scope="col">` + `<tbody>` + `<tr>/<td>`.
- **Estados**: rest, hover, focus, pressed, selected (`:has(input:checked)`), sort activo (`aria-sort`), pagination actual (`aria-current`); empty es GAP (hoy `0 rows` sin CTA).
- **Accesibilidad**: `aria-sort` en la columna activa; selección con checkboxes nativos; estado en el control, no solo color.
- **Cuándo usarlo**: comparar atributos en paralelo; set remoto con sort/filter/pagination; registros homogéneos.
- **Cuándo no**: ≤5-8 filas estáticas → List.
- **Mapping a Gelium**: `web/templates/data-table.html`, `/components/data-table`, contrato GET `?q=&sort=&dir=&page=&selection=`.
- **JS**: 0 base (GET); H refresh remoto (Progress + Toast).
- **Relación con patterns**: núcleo de Admin Resource, Search Results.

### Detail view ✖ (screen recipe, no componente)

- **Aliases**: master-detail, resource detail.
- **Intención**: muestra completa de una entidad con su contexto (header, meta, contenido, acciones).
- **Semántica HTML**: composición de Card/List/Divider/Tabs/Badge/Button/Menu/Dialog/Tooltip.
- **Cuándo usarlo**: el usuario necesita profundidad y contexto sobre un recurso.
- **Cuándo no**: acción corta y reversible → Dialog.
- **Mapping a Gelium**: patrón real en el chat WhatsApp (`demo-whatsapp.html:57-158`).
- **JS**: 0.

---

## 3. Patrones de estado

### Empty state ◐ (componente reusable pendiente)

- **Aliases**: empty, no data, vacío.
- **Intención**: comunicar que no hay datos y guiar al usuario (mensaje + opcional CTA).
- **Semántica HTML**: región con heading + texto + opcional acción; hoy ad-hoc `demo-wa-empty`.
- **Cuándo usarlo**: listados server-side sin resultados (gap real: tabla muestra `0 rows` sin guía).
- **JS**: 0.
- **Nota**: requisito bloqueante de Phase 4; convertir en componente reusable.

### Loading state / Skeleton ◐ / ✖

- **Aliases**: loading, placeholder, skeleton.
- **Intención**: indicar carga de datos sin bloquear; Skeleton muestra la forma del contenido futuro.
- **Estado actual**: solo Button loading (`aria-busy`) y Progress determinate/indeterminate; Skeleton no existe.
- **Cuándo usarlo**: feeds, dashboards, carga inicial de listados.
- **JS**: 0 (CSS puro).
- **Nota**: requisito bloqueante de Phase 4.

### Inline alert ◐ (componente reusable pendiente)

- **Aliases**: inline error, field error, alert.
- **Intención**: mensaje **persistente ligado al contexto** de un formulario/sección; sobrevive a la interacción.
- **Semántica HTML**: `<p role="alert">` junto al campo con `aria-invalid` en el control.
- **Cuándo usarlo**: errores de validación (contrato 422), advertencias de sección.
- **Cuándo no**: feedback transitorio → Toast; aviso global → Banner.
- **Mapping a Gelium**: hoy `text-field.html:8`, `select.html:89` (error de campo).
- **JS**: 0.

### Banner ✖

- **Aliases**: site banner, notice.
- **Intención**: aviso **persistente a nivel página/sitio** que exige acción (sesión expirada, mantenimiento, consent).
- **Semántica HTML**: región con rol `alert`/`status`, sin auto-dismiss.
- **Cuándo usarlo**: errores globales de Auth, mantenimiento, consentimiento.
- **Cuándo no**: nota ignorable → Callout; resultado transitorio → Toast.
- **JS**: 0.

### Callout ✖

- **Aliases**: note, tip, info box.
- **Intención**: contenido **informativo/promocional** sin urgencia ni requisito de acción.
- **Semántica HTML**: `<aside>` o bloque con heading opcional.
- **Cuándo usarlo**: contexto, tips, documentos, ayuda.
- **Cuándo no**: requiere acción → Banner; error del campo → Inline alert.
- **JS**: 0.

### Success feedback ✅ (reuso, sin componente nuevo)

- **Aliases**: success message, confirmation persistente.
- **Intención**: confirmación **NO efímera** de una operación exitosa; sobrevive a la navegación.
- **Implementación**: REUSA `inline-alert--success` (éxito de sección/form) y `banner--success` (éxito de página/operación global); NO es componente nuevo.
- **Semántica HTML**: igual que el patrón reusado (`<div>` en ambos partials: `inline-alert.html:1`, `banner.html:1`).
- **Estados/ARIA**: `role="status"` (polite) derivado del tone en ambos partials; `error` → `role="alert"`.
- **Cuándo usarlo**: guardado exitoso de settings, operación global completada (POST + 303 → página destino re-renderiza el success persistente).
- **Cuándo no**: feedback transitorio post-acción → Toast; error → `inline-alert--error` / `banner--error`.
- **Server contract**: POST + 303 → la página destino re-renderiza el success persistente; NUNCA `HX-Trigger loom:toast` para persistente.
- **JS**: 0.
- **Relación con patterns**: contraparte persistente del Toast (transitorio).

### Toast ✅

- **Aliases**: snackbar, notification. (Decisión Gelium: no crear Snackbar separado; usar solo como referencia visual.)
- **Intención**: feedback **transitorio del resultado de una acción**; no bloquea; auto-dismiss.
- **Anatomía**: `.ui-toast-{info|success|warning|error}`, `.ui-toast-message`, `.ui-toast-action`.
- **Semántica HTML**: `role="status"` (info/success) o `role="alert"` (error) en región `#loom-toast-region aria-live="polite"`.
- **Estados**: variantes info/success/warning/error; auto-dismiss 4s/8s pausable; dismiss manual.
- **Cuándo usarlo**: resultado de operaciones server-driven (`HX-Trigger loom:toast`).
- **Cuándo no**: validación de campos (NUNCA, `toast.go:129-133`); feedback persistente/crítico → Inline alert o Banner.
- **Mapping a Gelium**: `web/templates/toast.html`, `/components/toast`.
- **JS**: J* — sin JS toast inline persistente; con JS región aria-live + auto-dismiss.
- **Relación con patterns**: feedback de operaciones en Queue, Data table refresh, Resource Editor.

---

## 4. Patrones de workflow

### Queue ✖ (composición sobre List)

- **Aliases**: work queue, cola.
- **Intención**: workflow con **orden y posición operativos** (FIFO/pending) + estado por ítem + acción para avanzar.
- **Composición**: List two-line + Badge/Chip (tone) + Button + Toast; POST + redirect para mover estados.
- **Semántica HTML**: `<ul>` con items que exponen estado (badge) y acción.
- **Cuándo usarlo**: hay un "siguiente a atender" y estados que transicionar.
- **Cuándo no**: índice/referencia sin workflow → List; múltiples vías paralelas → Board.
- **Mapping a Gelium**: patrón real en sidebar WhatsApp (unread + window tone, `demo-whatsapp.html:26-55`).
- **JS**: 0 (POST + 303); H opcional.

### Board ✖

- **Aliases**: kanban, columnas.
- **Intención**: **múltiples vías paralelas** donde mover ítems entre estados es la tarea.
- **Composición**: columna = `<section>` con Cards + form por ítem (sin drag; POST + redirect + Toast).
- **Cuándo usarlo**: transferencia libre entre estados, comparación entre vías.
- **Cuándo no**: cola FIFO estricta → Queue.
- **JS**: 0.

### Steps ✖

- **Aliases**: wizard, stepper, proceso.
- **Intención**: progreso lineal por pasos con validación por paso.
- **Composición**: Progress determinate + navegación server-driven + validación 422 por paso + Inline alert.
- **Cuándo usarlo**: Resource Editor, Booking Flow.
- **Cuándo no**: flujo que cabe en una página → page simple.
- **JS**: 0.

### Timeline ✖

- **Aliases**: activity timeline, history.
- **Intención**: narrativa cronológica con **eje temporal explícito** (hitos + fechas + contexto).
- **Semántica HTML**: `<ol>` con items que exponen fecha y marker.
- **Cuándo usarlo**: audit trail, booking history, reconstruir un proceso con fechas.
- **Cuándo no**: actividad reciente tersa → Activity list (List one/two-line).
- **JS**: 0.

### Feed ✖ (composición)

- **Aliases**: activity feed, social feed.
- **Intención**: orientado a **evento/tiempo**; el ítem ES el evento; la novedad es el valor.
- **Composición**: List three-line o Card + avatar + Badge; loading/empty states críticos.
- **Cuándo usarlo**: enterarse de lo nuevo.
- **Cuándo no**: encontrar/administrar el set → Collection (Data table/List + filtros).
- **Mapping a Gelium**: thread de mensajes del demo WhatsApp (`demo-whatsapp.html:76-130`).
- **JS**: 0; H para refresh.

---

## 5. Navegación y control

### Tabs ✅

- **Aliases**: tab navigation, pestañas.
- **Intención**: alternar entre vistas del mismo contexto.
- **Semántica HTML**: `<nav>` → `<ul>` → `<li>` → `<a href>` con `aria-current="page"`; activo server-side.
- **Cuándo usarlo**: contenido del mismo nivel que compite por la misma superficie.
- **Cuándo no**: reclamar `role="tablist"`/roving focus sin teclado completo.
- **Mapping a Gelium**: `web/templates/tabs.html`, `/components/tabs`.
- **JS**: 0.

### Pagination ◐ (solo en Data table; standalone pendiente)

- **Intención**: navegar páginas de un set paginado.
- **Semántica HTML**: `<nav>` con links reales; actual como `<span aria-current="page">`.
- **Mapping a Gelium**: `data-table.html:68-72`.
- **JS**: 0 / H.

### Breadcrumb ✅

- **Intención**: ubicación jerárquica con navegación atrás; patrón public/content de Phase F (bloquea GEO §9/§14).
- **Anatomía**: `.ui-breadcrumb` (`<ol>`), `.ui-breadcrumb-item` (`<li>`); separador por CSS (`--ui-breadcrumb-separator`), nunca texto literal en markup (i18n).
- **Semántica HTML**: `<nav aria-label="Breadcrumb">` → `<ol>` → `<li>` → `<a>`; el crumb actual es `<span aria-current="page">`, NUNCA un link (markup canónico P1, `seo-patterns.md:50-64`).
- **Tokens**: scoped `--ui-breadcrumb-*` sobre `--ui-space-*`, `--ui-color-fg-muted`/`fg`, `--ui-type-label-sm`; forced-colors.
- **Mapping a Gelium**: `web/templates/breadcrumb.html`, `web/styles/breadcrumb.css`, datos `[]crumb{Href,Label,Current}` derivados de `componentRoutes()`/`navLinks()` (misma fuente que el JSON-LD `BreadcrumbList`).
- **JS**: 0.

### Section Heading ✅ (utilidad tipográfica, NO componente)

- **Decisión de Phase F**: es una **utilidad tipográfica**, no un componente: `.ui-section-heading` sobre `--ui-type-headline-sm` + `margin` con `--ui-space-*`; variante `.ui-section-heading--centered`; eyebrow/kicker opcional.
- **Semántica HTML**: SIEMPRE `<h2>`, nunca `<h1>` — la página posee un único h1 (P2, `seo-patterns.md:90`).
- **Mapping a Gelium**: `web/templates/section-heading.html`, `web/styles/section-heading.css`.
- **JS**: 0.

### Video ✅ (contenedor responsive, NO componente de contenido)

- **Decisión de Phase F**: es un **contenedor responsive**, no un componente de contenido: `.ui-video` con `aspect-ratio` literal y `<video controls>` nativo; "best used inside another component" (Split/Card/Hero).
- **Semántica HTML**: `<div class="ui-video">` → `<video controls poster loading="lazy">` → `<source>` + `<track kind="captions">` (a11y) + fallback. Sin autoplay.
- **Regla**: `aspect-ratio` NO se tokeniza (geometría estructural, como breakpoints/z-index); variante `.ui-video--aspect-4-3`.
- **Mapping a Gelium**: `web/templates/video.html`, `web/styles/video.css`.
- **JS**: 0 (controles nativos).

### Footer ✅

- **Intención**: chrome de sitio (brand + nav secundaria + legal) en todas las páginas; bloqueante de Phase G y del contrato SEO §3.
- **Anatomía**: `.ui-footer-brand`, `.ui-footer-nav` (grid → stack en narrow), `.ui-footer-section` con `.ui-footer-heading` (`<summary>`) + `.ui-footer-list`, `.ui-footer-legal`.
- **Semántica HTML**: `<footer>` + nav secundaria + legal; secciones plegables con `<details>/<summary>` nativos (sin `open` → collapsed por defecto; desktop fuerza open por CSS).
- **Tokens**: scoped `--ui-footer-*` sobre `--ui-color-surface`/`fg-muted`, `--ui-type-label-lg`/`body-sm`, `--ui-space-*`; forced-colors.
- **Mapping a Gelium**: `web/templates/footer.html`, `web/styles/footer.css`, slot `{{if .Footer}}` en `layout.html` tras `</main>`, datos `footerView{Brand, Sections[]{Title, Links[]navLink}, Legal}` en `internal/app/server.go`.
- **JS**: 0.

### Dialog ✅

- **Aliases**: modal, alert dialog (confirm).
- **Intención**: tarea **corta, enfocada y reversible** con contexto retenido; no cambia la URL.
- **Semántica HTML**: `<dialog closedby="any">` con `.ui-dialog-headline` (h2 + `aria-labelledby`), `.ui-dialog-content` (`aria-describedby`), `.ui-dialog-actions`.
- **Estados**: abierto/cerrado (top layer), light dismiss, confirm/cancel.
- **Accesibilidad**: foco al abrir, `aria-labelledby`/`aria-describedby`, acciones explícitas Cancel/Confirm.
- **Cuándo usarlo**: confirmar, quick edit, picker.
- **Cuándo no**: flujo largo o profundo → Page/Steps.
- **Mapping a Gelium**: `web/templates/dialog.html`, `/components/dialog`.
- **JS**: D — `<dialog>` + Invoker Commands declarativos + `closedby="any"`; fallback server-rendered en navegadores previos.

### Popover (mecanismo, NO patrón)

- **Decisión de vocabulario**: "Popover" es el nombre de la **primitiva web** (`popover`/`popovertarget`) y NO se define como patrón canónico de UI. Es el **mecanismo** de top layer/overlay que Menu y futuros menús contextuales usan.
- **Mapping a Gelium**: `menu.html:44,48` (Popover API declarativa, zero JS); Tooltip la descartó por no-Baseline.

### Drawer → Navigation drawer ✅

- **Aliases**: drawer (alias genérico → canónico "Navigation drawer").
- **Intención**: superficie lateral de navegación; modal (tarea enfocada) o permanente (layout).
- **Semántica HTML**: modal = `<dialog>`; permanente = `<nav>` → `<ul>` → `<a href>`.
- **Estados**: activo server-side (`aria-current`), abierto/cerrado.
- **Mapping a Gelium**: `/components/navigation-drawer`, `navigation-drawer.html:8-38`.
- **JS**: D (modal) / 0 (permanente).

### Tooltip ✅

- **Aliases**: tip, hover hint.
- **Intención**: texto corto de ayuda al hover/focus; nunca esconde información esencial.
- **Semántica HTML**: `role="tooltip"` + `aria-describedby` en el control.
- **Estados**: visible en `:hover`/`:focus-within`; variantes plain/rich.
- **Accesibilidad**: información esencial nunca solo en tooltip.
- **Mapping a Gelium**: `web/templates/tooltip.html`, `/components/tooltip`.
- **JS**: 0 (CSS reveal; Interest Invokers rechazados por no-Baseline).

### Multi-select (capacidad, NO widget)

- **Decisión de vocabulario**: "Multi-select" es la **capacidad de selección múltiple sobre un patrón huésped**, con regla de elección, no un componente único. Hoy vive en 4 implementaciones: List con checkboxes, Chips filter, Segmented multi, filas de Data table.
- **Regla**: se aplica sobre List (checkboxes nativos), Data table (selección de filas), Menus/Chips (filtros) — siempre con controles nativos, server round-trip.

### Combobox ✖

- **Intención**: campo de texto con selección filtrable (typeahead).
- **Candidato nativo**: `<input>` + `<datalist>`/autocomplete GET server-driven; Select menu es listbox en `<dialog>`, no combobox.
- **JS**: requiere auditoría platform-first antes de comprometerlo; probablemente H/JS mínimo.

### Date picker ✖

- **Intención**: selección de fecha.
- **Candidato nativo**: `<input type="date">` + calendario server-driven.
- **JS**: auditoría platform-first antes de comprometer.

### Form ◐ (patrón nativo, no componente)

- **Intención**: agrupar controles con submit/validación server-side.
- **Semántica HTML**: `<form>` nativo + `<fieldset>`/`<legend>` para grupos; `Field` es primitive interna, no publicable.
- **Contrato**: HTTP 422 + `X-Gelium-Validation`; GET para estado de listados.
- **JS**: 0; H opcional.

---

## 6. Componentes de acción e input (referencia)

Completan el vocabulario operativo, ya implementados y aceptados: Button (primary/secondary/outline/text; link vs button; disabled/loading con `aria-busy`), Icon button (toggle `aria-pressed`), FAB (extended/icon-only con AriaLabel obligatorio), Chips (assist/filter/input/suggestion), Segmented buttons (radio/checkbox/button nativos en `<fieldset>`), Menu (Popover API, items nativos, zero JS), Text field (floating label filled/outlined, 422), Checkbox/Radio/Switch/Select/Slider (controles nativos), Progress (determinate/indeterminate), Badge (dot/count, nunca color-only), Divider, Icon, Elevation, Focus ring.

---

## 7. Necesidad de JavaScript por capa

| Capa | JS |
|---|---|
| Patrones de datos (Card, Panel, List, Table, Detail view) | 0 |
| Patrones de estado (Empty, Loading/Skeleton, Inline alert, Banner, Callout, Success feedback) | 0 |
| Toast | J* (fallback no-JS real) |
| Workflow (Queue, Board, Steps, Timeline, Feed) | 0; H opcional para refresh |
| Overlays (Dialog, Drawer modal, Menu, Select menu) | D (primitivas declarativas + fallback server-rendered) |
| Combobox, Date picker | Auditoría platform-first pendiente; H/JS mínimo probable |

Regla: **ningún término entregado exige JS para el flujo principal.**

---

## 8. Conflictos de naming resueltos

1. **Popover** → mecanismo web, no patrón de UI (Menu es el patrón).
2. **Multi-select** → capacidad sobre patrón huésped, no widget.
3. **Drawer** → alias de Navigation drawer (canónico Gelium).
4. **Snackbar** → alias visual de Toast; no crear componente.
5. **Alert** → rol ARIA (`role="alert"`) vs componente Inline alert vs Banner; resuelto por capa de persistencia.
6. **Table vs Data table** → Gelium solo tiene Data table; tabla nativa simple existe ad-hoc en admin.
7. **Select vs Select menu vs Combobox** → Select nativo; Select menu = listbox en dialog; Combobox = gap con filtrado.
8. **Steps vs Pagination** → Steps es proceso con validación por paso; Pagination es navegación de un set.
9. **Queue vs List two-line** → Queue añade orden operativo + estado + acción de avance.

---

**Definición de done (Phase 2)**: vocabulario completo con anatomía/estados/accesibilidad/cuándo usarlo/cuándo no; conflictos de naming resueltos; mapeo a Gelium con rutas reales; aprobado antes de implementar componentes nuevos.
