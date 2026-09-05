# Gelium UI — Composition Rules

> Reglas de composición del sistema Gelium UI.
> Fase 3 del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/handoffs/composition-audit.md`, `docs/handoffs/vocabulary-audit.md` y el catálogo real (`internal/app/*`, `web/templates/*`).

---

## 1. Propósito

Estas reglas deciden **qué patrón usar y cómo combinarlo** en una pantalla server-rendered. No describen componentes: describen la elección entre ellos y su integración con los contratos server-driven. Toda screen recipe (Phase 4) y todo componente nuevo DEBE justificarse contra estas reglas (rationale obligatorio).

## 2. Referencias de contratos (Phase E)

Toda composición opera dentro de los contratos de Phase E. Estos son los documentos de referencia y qué aporta cada uno a la composición:

| Contrato | Documento | Aporta a la composición |
|---|---|---|
| UX principles | `docs/gelium-ui-ux-principles.md` | Los principios de UX que toda pantalla respeta (HTML-first, no-JS end-to-end, native before ARIA). |
| UX patterns | `docs/gelium-ui-ux-patterns.md` | Patrones de UX reutilizables para combinar superficies y componentes. |
| Content rules | `docs/gelium-ui-content-rules.md` | Reglas de contenido (copy, etiquetas, microcopy) que toda superficie debe satisfacer. |
| Accessibility contract | `docs/gelium-ui-accessibility-contract.md` | Contrato de accesibilidad WCAG 2.1 AA y el estado de los gaps G1-G11. |
| SEO contract | `docs/gelium-ui-seo-contract.md` | Metadata server-driven por ruta: description, canonical, robots, OG, JSON-LD. |
| SEO patterns | `docs/gelium-ui-seo-patterns.md` | Patrones de composición SEO: head server-driven, breadcrumbs, JSON-LD por tipo de página. |
| GEO contract | `docs/gelium-ui-geo-contract.md` | Contrato GEO: entidad única **Gelium UI**, resúmenes citables, provenance visible. |
| GEO patterns | `docs/gelium-ui-geo-patterns.md` | Patrones de contenido GEO: answer-first, headings, definiciones, citas visibles. |

## 2. Screen grammar (gramática de pantalla)

Toda pantalla se describe como:

```text
superficie (nivel de inmersión)
  × usuario (tarea primaria)
  × patrón (qué muestra/organiza)
  × estado (vacío/carga/error/rest)
  × acción (qué puede hacer)
  × contrato server (cómo se comunica con el servidor)
  × fallback no-JS + enhancement HTMX
```

- **Superficie**: contexto inmersivo (app shell), sección de página, panel, card, overlay.
- **Usuario**: la tarea primaria del usuario en esa superficie, no el rol del sistema.
- **Patrón**: uno de los términos del vocabulario (`docs/gelium-ui-vocabulary.md`).
- **Estado**: todo patrón de datos declara sus 4 estados (rest/vacío/carga/error); el estado vacío y de carga no son opcionales.
- **Acción**: cada acción se resuelve con un componente de acción y un contrato (POST+redirect, GET params, 422, toast).

## 3. Surface grammar (gramática de superficie)

1. **Una tarea primaria por superficie**: si una superficie tiene dos tareas primarias, se divide.
2. **La página es la unidad de URL**: estado navegable = URL; todo lo que necesita deep-link/back es page, no overlay.
3. **Overlays** (Dialog, Menu, Popover, Drawer modal, Tooltip) para tareas cortas y reversibles con contexto retenido; nunca para flujos largos.
4. **Paneles para agrupar contexto** (estáticos); Cards para instancias repetidas; el patrón se elige por la naturaleza del contenido, no por el gusto.

## 4. Pattern selection (criterios de decisión)

### 4.1 table vs list

| Señal | Table | List |
|---|---|---|
| Comparación | Atributos en paralelo entre columnas | Escaneo de unidades verticales |
| Set | Grande/remoto, sort+filter+pagination | Pequeño/estático |
| Filas | Registros homogéneos | Unidades discretas con jerarquía 1-3 líneas |
| Acción | En fila (row action) o selección | Acción primaria por ítem |

- **Anti-regla**: no usar Data table para ≤5-8 filas estáticas → List.

### 4.2 list vs queue

| Señal | List | Queue |
|---|---|---|
| Orden | Presentacional (índice/catálogo) | Operativo (FIFO/pending) |
| Estado por ítem | Opcional (badge decorativo) | Semántico (tone, unread, window) |
| Acción | Navegar/seleccionar | Avanzar al siguiente estado |
| Concepto | Referencia | "Siguiente a atender" |

- **Anti-regla**: un índice sin workflow no es queue; una cola no se modela con tabla.

### 4.3 queue vs board

| Señal | Queue | Board |
|---|---|---|
| Vías | Una sola, ordenada | Múltiples paralelas |
| Tarea dominante | Tomar el siguiente | Mover entre estados |
| Movimiento | Avance secuencial | Transferencia libre entre columnas |

- **Anti-regla**: no usar Board para colas FIFO estrictas (fricción sin semántica). Board server-side = columna `<section>` con Cards + form por ítem, sin drag.

### 4.4 card vs panel

| Señal | Card | Panel |
|---|---|---|
| Repetición | Instancia de un set | Región única de layout |
| Acción | Entrada propia (a/button) | Contexto de la página |
| Uso | Grids, feeds, KPIs | Secciones de configuración |

### 4.5 feed vs collection

| Señal | Feed | Collection |
|---|---|---|
| Orientación | Evento/tiempo | Entidad/set |
| Valor | La novedad | La totalidad con filtros |
| Tarea | Enterarse | Encontrar/administrar |
| Estados | Loading/empty críticos | Sort/filter/pagination |

### 4.6 timeline vs activity list

| Señal | Timeline | Activity list |
|---|---|---|
| Eje temporal | Explícito (hitos + fechas) | No formal |
| Orden | Ascendente semántico (proceso) | Inverso (reciente primero) |
| Uso | Audit trail, history | Actividad reciente tersa |
| Implementación | Markers + fechas | List one/two-line |

### 4.7 dialog vs page

| Señal | Dialog | Page |
|---|---|---|
| Duración | Corta, enfocada | Profunda/completa |
| Contexto | Retenido (quick edit, confirm, picker) | URL + back + deep-link |
| Reversibilidad | Alta | Media (navegación) |

- **Anti-regla**: no abrir flujo largo en Dialog (editor/booking multi-paso → Steps/page).

### 4.8 toast vs inline alert vs banner vs callout

| Señal | Toast | Inline alert | Banner | Callout |
|---|---|---|---|---|
| Persistencia | Transitorio (auto-dismiss) | Persistente junto al campo | Persistente nivel página/sitio | Persistente ignorable |
| Origen | Resultado de acción | Validación/error de sección | Aviso global que exige acción | Contenido informativo |
| Rol | `status`/`alert` (aria-live) | `alert` junto al campo | `alert`/`status` sin dismiss | Nota, sin urgencia |
| Contrato | `HX-Trigger gelium:toast` | 422 + `X-Gelium-Validation` | — | — |

- **Anti-reglas**: validación nunca toast; feedback persistente/crítico nunca toast.

## 5. Entity-to-pattern anti-rules

1. **No usar Data table para ≤5-8 filas estáticas** → List (costo de sort/filter sin beneficio).
2. **No usar Board para colas FIFO estrictas** → Queue.
3. **No reportar errores de validación como Toast** → 422 + Inline alert (`toast.go:129-133`).
4. **No usar Toast para feedback persistente/crítico** → Inline alert o Banner.
5. **No inventar estados de listado fuera de la URL** (sort/filter/page/selection en cliente) → GET con params estables.
6. **No reclamar `role="tablist"`/roving focus sin teclado completo** → links reales con `aria-current` (Tabs).
7. **No abrir un flujo largo en Dialog** → page/steps.
8. **No usar spinner cuando existe Progress determinate/indeterminate** → `.ui-progress` + `aria-busy`/`role="status"`.
9. **No dejar empty states sin mensaje/CTA** → todo listado server-side define vacío; Data table ya compone `Empty state` con mensaje y CTA.
10. **No usar color como único portador de estado** → estado en el control nativo + forced colors.
11. **No introducir JS para lo que un form GET ya resuelve** → platform-first.

## 6. Density rules

- La densidad es responsabilidad del core/theme (Phase 1: `--ui-density-*`), NO de las recetas.
- No definir densidades por pantalla hasta que existan tokens de densidad.
- La geometría de controles se consume vía tokens `--ui-size-*`; nunca literales en recetas.

## 7. Responsive rules

1. **Fluido primero**: layouts por `auto-fit/minmax` y `min()/clamp()` antes de breakpoints (patrón real del sistema: `card.css:30`, `dialog.css:3-5`, `menu.css:53`).
2. **Breakpoints solo cuando el layout fluido no alcanza** (tokens `--ui-breakpoint-*`).
3. **Tablas grandes se resuelven server-side** (paginación), nunca con scroll horizontal de componente.
4. **Overlays ya son fluidos** (`calc(100vw - n)`): dialog/menu/tooltip/drawer se usan sin breakpoints.
5. **Drawer** resuelve responsive por variantes (modal vs permanente), no por media query.

## 8. State matrix (matriz de estados por patrón)

| Patrón | Rest | Hover/Focus/Pressed | Selected | Disabled | Empty | Loading | Error |
|---|---|---|---|---|---|---|---|
| Card | ✅ | ✅ | ✅ (control interno) | n/a | n/a | n/a | n/a |
| List | ✅ | ✅ | ✅ (`:checked`) | ✅ | GAP | GAP | n/a |
| Data table | ✅ | ✅ | ✅ (`:has(input:checked)`) | por control | ✅ (Empty state) | H | ✅ (Error state / Inline alert) |
| Queue | ✅ | ✅ | ✅ | por control | GAP | GAP | GAP |
| Feed | ✅ | ✅ | n/a | n/a | **GAP** | **GAP** | GAP |
| Dashboard | ✅ | ✅ | n/a | n/a | GAP | **GAP** | GAP |
| Dialog | ✅ | ✅ | n/a | acciones | n/a | n/a | inline |
| Form | ✅ | ✅ | ✅ | ✅ | n/a | submit | 422 inline |

`GAP` en esta matriz significa que la composición todavía no tiene una cobertura de recipe o contrato suficiente; no significa que falte necesariamente un primitive. Empty state, Skeleton, Inline alert, Banner y Error state son primitives existentes. La cobertura de loading para Data table queda condicionada a una espera real de región, no a la primera respuesta server-rendered.

## 9. Server-driven rules

Toda composición reutiliza los contratos existentes; NO se inventan otros:

1. **Validación**: HTTP 422 + `X-Gelium-Validation: true`; sin HX página completa 422; con HX fragmento. La validación nunca dispara toast.
2. **Feedback transitorio**: `HX-Trigger: {"gelium:toast":{"type":"info|success|warning|error","message":"…"}}`; sin JS toast inline persistente. Los errores de transporte HTMX (red/5xx) muestran un toast genérico transitorio (`app.js`), nunca un mensaje de validación.
3. **Estado de listados**: GET con params estables (`?q=&sort=&dir=&page=&selection=`), vocabularios cerrados sanitizados, `HX-Request` bifurca fragmento vs página completa, links preservan estado.
4. **Mover estados (workflow)**: POST + 303 SeeOther redirect (patrón WhatsApp) — simple, sin fragmentos, funciona sin JS.
5. **Indexabilidad (SEO/GEO)**: la metadata se resuelve server-driven en el choke point de render. Los demos (`/demo/*`) y ejemplos (`/examples/*`) son `noindex, nofollow`; el resto `index, follow`. Toda página indexable emite canonical sin query (`siteBaseURL` + ruta limpia), `lang` server-driven (`es` en demos WhatsApp) y JSON-LD básico en home. Restricciones completas en `docs/gelium-ui-seo-contract.md` (§4, §16) y `docs/gelium-ui-geo-contract.md` (entidad única **Gelium UI**).

**Regla de oro**: si un estado es navegable (listado, filtro, sort, paginación, selección), es una URL. Sin URL = sin no-JS, sin deep-link, sin back.

## 10. Accessibility rules (composición)

El contrato completo de accesibilidad (WCAG 2.1 AA, invariantes y estado de los gaps G1-G11) vive en `docs/gelium-ui-accessibility-contract.md`. Para composición, las reglas mínimas son:

1. Estructura de landmarks: `<header>`, `<nav>`, `<main>`, `<aside>`, `<footer>` según el rol de cada superficie.
2. Orden del documento = orden visual; no reorganizar con CSS para lectura.
3. Overlays: foco gestionado por la primitiva nativa (`<dialog>`, popover); fallback server-rendered si la primitiva no está disponible.
4. Todo patrón de datos declara estados vacío/carga/error accesibles (`aria-busy`, `role="status"`, mensajes no-color-only).
5. Keyboard: navegación nativa; roving focus solo cuando el teclado completo está resuelto.
6. Reduced motion y forced colors por diseño (estrategia central del core).

## 11. Rationale obligatorio

Antes de implementar cualquier componente o receta nueva, el PRD/issue debe responder:

1. ¿Qué patrón del vocabulario cubre y por qué ese patrón (regla de selección §4)?
2. ¿Qué screen recipe desbloquea?
3. ¿Qué estados de la matriz cubre y cuáles son GAP?
4. ¿Qué contrato server-driven usa (422/toast/GET/POST+redirect)?
5. ¿Cómo funciona sin JS y cómo mejora con HTMX?
6. ¿Qué anti-regla de §5 podría violar y por qué se justifica la excepción?

Sin rationale aprobado, no se implementa.

---

**Definición de done (Phase 3)**: reglas escritas, diferenciaciones explícitas con evidencia del catálogo real, anti-rules con justificación, state matrix completa, server-driven rules ancladas a los contratos existentes, aprobadas antes de Phase 4.
