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

Los patrones de estado comunican la condición de la interfaz: qué está pasando, qué falló, qué falta. La regla del roadmap es **persistente-contextual ≠ transitorio-acción**: el feedback que debe sobrevivir a la interacción (errores, validaciones, confirmaciones, ausencia de datos) es **persistente** y vive anclado a su contexto; el feedback del resultado inmediato de una acción es **transitorio** y desaparece solo.

**No confundir:**
- **Persistente (contextual)**: Empty state, Inline alert, Banner, Callout, Error state, Validation summary, Success feedback.
- **Transitorio (acción)**: Toast, Loading state.

| Pattern | Partial | CSS | Test |
|---|---|---|---|
| Empty state | empty-state.html | empty-state.css | styles_empty_state_test.go |
| Loading state | button.html `aria-busy` | button.css | styles_button_test.go |
| Skeleton | skeleton.html | skeleton.css | styles_skeleton_test.go |
| Inline alert | inline-alert.html | inline-alert.css | styles_inline_alert_test.go |
| Banner | banner.html | banner.css | styles_banner_test.go |
| Callout | callout.html | callout.css | styles_callout_test.go |
| Error state | error-state.html | error-state.css | styles_error_state_test.go |
| Validation summary | validation-summary.html | validation-summary.css | styles_validation_summary_test.go |
| Success feedback | reuso inline-alert--success / banner--success | — | — |
| Toast | toast.html | toast.css | styles_toast_test.go |

### Empty state ✅ · Persistente

- **Aliases**: empty, no data, vacío.
- **Intención**: comunicar que no hay datos y guiar al usuario (mensaje + opcional CTA).
- **Semántica HTML**: región con heading + texto + opcional acción (`<div role="status">` en `empty-state.html:1`).
- **Estados**: vacío (estado único; opcional compact).
- **Accesibilidad**: `role="status"` (polite) en el contenedor; el CTA es un control real con foco propio.
- **Server contract**: none.
- **Cuándo usarlo**: listados server-side sin resultados (gap real: tabla muestra `0 rows` sin guía).
- **Cuándo no**: feedback del resultado de una acción → Toast; aviso global → Banner.
- **Mapping a Gelium**: `web/templates/empty-state.html`, `web/styles/empty-state.css`, `web/styles_empty_state_test.go`.
- **JS**: 0.

### Loading state / Skeleton ✅ · Transitorio

- **Aliases**: loading, placeholder, skeleton.
- **Intención**: indicar carga de datos sin bloquear; Skeleton muestra la forma del contenido futuro.
- **Composición**: Button con `aria-busy` (`button.html`) + Progress determinate/indeterminate (`progress.html`) + Skeleton (`skeleton.html`); NO es un partial nuevo.
- **Semántica HTML**: estado en el control (`aria-busy`) o región de progreso (`role="progressbar"`); Skeleton como placeholders estáticos.
- **Estados**: loading (Button `aria-busy`, Progress determinate/indeterminate); skeleton (forma del contenido futuro).
- **Accesibilidad**: `aria-busy` en el control durante la carga; Progress anuncia progreso; Skeleton estático, no interactivo.
- **Server contract**: none (la fase de carga se renderiza desde el servidor).
- **Cuándo usarlo**: feeds, dashboards, carga inicial de listados.
- **Cuándo no**: feedback del resultado de una acción → Toast; ausencia de datos → Empty state.
- **Mapping a Gelium**: `web/templates/button.html` (`aria-busy`), `web/templates/progress.html`, `web/templates/skeleton.html`; tests `styles_button_test.go`, `styles_progress_test.go`, `styles_skeleton_test.go`.
- **JS**: 0 (CSS puro).

### Inline alert ✅ · Persistente

- **Aliases**: inline error, field error, alert.
- **Intención**: mensaje **persistente ligado al contexto** de un formulario/sección; sobrevive a la interacción.
- **Semántica HTML**: `<p role="alert">` (tone error) o `role="status"` (resto) junto al campo con `aria-invalid` en el control; partial `inline-alert.html:1` (`<div role="alert|status">` según tone).
- **Estados**: error, warning, info, success (tone).
- **Accesibilidad**: `role="alert"` en tone error; `role="status"` en el resto; el campo conserva `aria-invalid` + `aria-describedby`.
- **Server contract**: 422 con errores de campo (campo + mensaje); re-render del partial con el valor preservado.
- **Cuándo usarlo**: errores de validación (contrato 422), advertencias de sección.
- **Cuándo no**: feedback transitorio → Toast; aviso global → Banner.
- **Mapping a Gelium**: `web/templates/inline-alert.html`, `web/styles/inline-alert.css`, `web/styles_inline_alert_test.go`; hoy `text-field.html:8`, `select.html:89` (error de campo).
- **JS**: 0.

### Banner ✅ · Persistente

- **Aliases**: site banner, notice.
- **Intención**: aviso **persistente a nivel página/sitio** que exige acción (sesión expirada, mantenimiento, consent).
- **Semántica HTML**: región con `role="alert"` (tone error) o `role="status"` (resto), sin auto-dismiss (`banner.html:1`).
- **Estados**: error, warning, info, success (tone); dismiss explícito por POST.
- **Accesibilidad**: `role="alert"` en tone error; `role="status"` en el resto; dismiss por POST explícito, nunca temporizado.
- **Server contract**: re-render del partial desde el servidor; el dismiss es un POST explícito.
- **Cuándo usarlo**: errores globales de Auth, mantenimiento, consentimiento.
- **Cuándo no**: nota ignorable → Callout; resultado transitorio → Toast.
- **Mapping a Gelium**: `web/templates/banner.html`, `web/styles/banner.css`, `web/styles_banner_test.go`.
- **JS**: 0.

### Callout ✅ · Persistente

- **Aliases**: note, tip, info box.
- **Intención**: contenido **informativo/promocional** sin urgencia ni requisito de acción.
- **Semántica HTML**: `<aside>` con heading opcional (`callout.html:1`).
- **Estados**: estático (variantes informativas).
- **Accesibilidad**: contenido estático leído en orden de documento; sin roles dinámicos.
- **Server contract**: none.
- **Cuándo usarlo**: contexto, tips, documentos, ayuda.
- **Cuándo no**: requiere acción → Banner; error del campo → Inline alert.
- **Mapping a Gelium**: `web/templates/callout.html`, `web/styles/callout.css`, `web/styles_callout_test.go`.
- **JS**: 0.

### Error state ✅ · Persistente

- **Aliases**: error page, fatal error, pantalla de error.
- **Intención**: comunicar un error **no recuperable en contexto** (HTTP 4xx/5xx) y ofrecer una vía de salida (retry / inicio).
- **Semántica HTML**: región `role="alert"` con código + heading + mensaje + retry CTA opcional (`error-state.html:1-5`).
- **Estados**: error (código + título + cuerpo); retry opcional.
- **Accesibilidad**: `role="alert"` en la región; el código es decorativo (`aria-hidden="true"`); el CTA es un enlace real con foco propio.
- **Server contract**: respuesta de error 4xx/5xx renderiza el partial — nunca Toast.
- **Cuándo usarlo**: 404, 500, fallos de operación crítica sin contexto de formulario.
- **Cuándo no**: error recuperable de campo → Inline alert / Validation summary; resultado transitorio → Toast.
- **Mapping a Gelium**: `web/templates/error-state.html`, `web/styles/error-state.css`, `web/styles_error_state_test.go`.
- **JS**: 0.

### Validation summary ✅ · Persistente

- **Aliases**: error summary, resumen de errores.
- **Intención**: resumir al inicio del formulario los campos que requieren atención, con enlaces a cada uno.
- **Semántica HTML**: región `role="alert"` + heading ("N campos requieren atención") + `<ul>` de `<a href="#{campo}-error">` (`validation-summary.html:1-6`).
- **Estados**: error de validación (1+ campos).
- **Accesibilidad**: `role="alert"` al renderizar; cada enlace lleva al campo, que conserva `aria-invalid` + `aria-describedby`; heading jerárquico (`h{level}`).
- **Server contract**: 422 renderiza el summary con los items de error.
- **Cuándo usarlo**: formularios con validación server-side (contrato 422).
- **Cuándo no**: un solo campo con error → Inline alert; resultado transitorio → Toast.
- **Mapping a Gelium**: `web/templates/validation-summary.html`, `web/styles/validation-summary.css`, `web/styles_validation_summary_test.go`.
- **JS**: 0.

### Success feedback ✅ · Persistente

- **Aliases**: success message, confirmation persistente.
- **Intención**: confirmación **NO efímera** de una operación exitosa; sobrevive a la navegación.
- **Implementación**: REUSA `inline-alert--success` (éxito de sección/form) y `banner--success` (éxito de página/operación global); NO es componente nuevo.
- **Semántica HTML**: igual que el patrón reusado (`<div>` en ambos partials: `inline-alert.html:1`, `banner.html:1`).
- **Estados**: success (persistente); error → tone error del patrón reusado.
- **Accesibilidad**: `role="status"` (polite) derivado del tone en ambos partials; `error` → `role="alert"`.
- **Cuándo usarlo**: guardado exitoso de settings, operación global completada (POST + 303 → página destino re-renderiza el success persistente).
- **Cuándo no**: feedback transitorio post-acción → Toast; error → `inline-alert--error` / `banner--error`.
- **Server contract**: POST + 303 → la página destino re-renderiza el success persistente; NUNCA `HX-Trigger loom:toast` para persistente.
- **Mapping a Gelium**: reuso de `web/templates/inline-alert.html` (tone success) y `web/templates/banner.html` (tone success); sin partial propio.
- **JS**: 0.
- **Relación con patterns**: contraparte persistente del Toast (transitorio).

### Toast ✅ · Transitorio

- **Aliases**: snackbar, notification. (Decisión Gelium: no crear Snackbar separado; usar solo como referencia visual.)
- **Intención**: feedback **transitorio del resultado de una acción**; no bloquea; auto-dismiss.
- **Anatomía**: `.ui-toast-{info|success|warning|error}`, `.ui-toast-message`, `.ui-toast-action`.
- **Semántica HTML**: `role="status"` (info/success) o `role="alert"` (error) en región `#loom-toast-region aria-live="polite"`.
- **Estados**: variantes info/success/warning/error; auto-dismiss 4s/8s pausable; dismiss manual.
- **Accesibilidad**: `role="status"`/`role="alert"` según variante + `aria-live="polite"` en la región (`toast.html:10`).
- **Server contract**: operación server-driven dispara el toast vía HX-Trigger (prefijo congelado — referencia únicamente; nunca para persistente).
- **Cuándo usarlo**: resultado de operaciones server-driven (`HX-Trigger loom:toast`).
- **Cuándo no**: validación de campos (NUNCA, `toast.go:129-133`); feedback persistente/crítico → Inline alert o Banner.
- **Mapping a Gelium**: `web/templates/toast.html`, `web/styles/toast.css`, `web/styles_toast_test.go`, `/components/toast`.
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

### Hero ✅ (Billboard; NO se llama Callout)

- **Decisión de naming**: el "Hero" promocional de Protocol (full-width) NO se implementa como Callout — ese nombre ya es el tip box de Phase D. Es una **composición**: `h1` `--ui-type-display-lg` (la página posee un único h1, P2) + subtitle + CTA(s) (Button link) + media de fondo opcional con scrim.
- **Anatomía**: `.ui-hero` → `.ui-hero-media` (fondo absoluto, variante `.ui-hero--has-media`) + `.ui-hero-content` (`.ui-hero-eyebrow`, `.ui-hero-title`, `.ui-hero-subtitle`, `.ui-hero-actions`).
- **Semántica HTML**: `<section>` con `h1` de página; media decorativa fuera del flujo de lectura.
- **Tokens**: scoped `--ui-hero-*` sobre `--ui-color-surface`/`scrim`, `--ui-type-display-lg`/`body-lg`, `--ui-space-*`; forced-colors; padding generoso en wide (media query).
- **Mapping a Gelium**: `web/templates/hero.html`, `web/styles/hero.css`; desbloquea Landing/Public-Feed (Phase G).
- **JS**: 0.

### Split ✅

- **Intención**: composición editorial de 2 columnas (media + cuerpo) que apila en narrow.
- **Anatomía**: `.ui-split` (grid `repeat(2, minmax(0,1fr))`) → `.ui-split-media` (slot img/video) + `.ui-split-body` (`.ui-split-eyebrow`, `.ui-split-title`, `.ui-split-copy`, `.ui-split-action`).
- **Semántica HTML**: `<section>`; la tipografía del cuerpo NO se aplica por defecto (el consumidor conserva el contrato `.prose`).
- **Bidi**: RTL automático — las columnas fluyen en orden de dirección del documento (media primero → derecha en RTL), sin `left/right` literales.
- **Tokens**: scoped `--ui-split-*` sobre `--ui-color-fg`/`fg-muted`, `--ui-type-headline-sm`/`body-lg`, `--ui-radius-sm`, `--ui-space-*`; forced-colors.
- **Mapping a Gelium**: `web/templates/split.html`, `web/styles/split.css`; best-used-in: Landing/Public-Feed (Phase G), junto a Video.
- **JS**: 0.

### Feature Card ✅ (composición, NO primitiva)

- **Decisión**: es **composición de Card + CTA Link**, no un componente nuevo: el wrapper `ui-feature-card` reusa `.ui-card` (elevated) + media + `.ui-card-title`/`.ui-card-body` + `.ui-card-action` (Button link). Variante horizontal descartada (deprecada upstream).
- **Anatomía**: `.ui-feature-card` → `.ui-feature-card-media` (aspect-ratio literal 16:9, no tokenizado) + `.ui-feature-card-body`.
- **CSS**: mínimo — solo geometría (media aspect, spacing); superficie/sombra/foco vienen de `.ui-card`.
- **Mapping a Gelium**: `web/templates/feature-card.html`, `web/styles/feature-card.css`.
- **JS**: 0.

### Language Switcher ✅

- **Intención**: cambiar el idioma del sitio; es **navegación GET, nunca POST**.
- **Semántica HTML**: `<form method="get" action="{{.Action}}">` + `<label>` + `<select name="lang">` + **submit siempre visible** (cero auto-submit JS). El server responde al `?lang=<code>` con **303 a la URL localizada** y resuelve `<html lang>`/RTL server-side.
- **Alcance**: el patrón es la **primitiva lista** (contrato); el modelo de locales (`?lang=` → 303) queda **fuera de alcance** — no hay i18n real todavía; el server debe resolverlo cuando exista.
- **Tokens**: scoped `--ui-language-switcher-*` sobre `--ui-color-border`/`fg-muted`, `--ui-type-body-sm`, `--ui-radius-sm`, `--ui-space-*`; forced-colors.
- **Mapping a Gelium**: `web/templates/language-switcher.html`, `web/styles/language-switcher.css`; se compone dentro del Footer (nav secundaria).
- **JS**: 0.

### Newsletter ✅

- **Intención**: conversión de audiencia (suscripción); formulario **zero-JS** con contrato server completo.
- **Semántica HTML**: `<aside class="ui-newsletter" aria-labelledby>` → título `h2` + descripción + `<form method="post" action="{{.Action}}">` con email (`type="email"` + `required`) + submit (Button). Success reemplaza el form por un `<p role="status">` persistente.
- **Contrato server**: **POST + 422 + `X-Loom-Validation: true`** (header real del código; el roadmap lo escribe `X-Gelium-Validation`, ver `screen-recipes-audit.md:17`) en email inválido, re-render con `inline-alert--error` y valor preservado; **POST → 200 success** (ejemplo) o **POST + 303 → GET success** (producción, contrato d). HTMX opcional (swap del aside).
- **Tokens**: scoped `--ui-newsletter-*` sobre `--ui-color-surface-container`/`error`, `--ui-type-headline-sm`, `--ui-radius-*`, `--ui-space-*`; forced-colors.
- **Mapping a Gelium**: `web/templates/newsletter.html`, `web/styles/newsletter.css`, `internal/app/newsletter.go`, ejemplo `GET/POST /examples/newsletter` (noindex).
- **JS**: 0 (HTMX opcional).

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
- **Contrato**: HTTP 422 + `X-Loom-Validation`; GET para estado de listados.
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
