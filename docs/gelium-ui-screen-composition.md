# Gelium UI — Screen Composition

> Cómo se compone una pantalla desde los registries del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Unifica: `docs/gelium-ui-component-registry.md` (qué componentes existen), `docs/gelium-ui-pattern-registry.md` (qué patterns componen), `docs/gelium-ui-dependency-metadata.md` (qué desbloquea cada pieza) y `docs/gelium-ui-screen-recipes.md` (las 3 recipes implementadas).
> Fuentes de autoridad para el ejemplo real: `internal/app/recipe_admin_resource.go`, `web/templates/recipe-admin-resource.html`, `web/styles/recipe-admin-resource.css`.

---

## 1. Cómo se compone una pantalla (flujo)

1. **Elegir el patrón** con la gramática de pantalla (`composition-rules.md` §2): `superficie × usuario × patrón × estado × acción × contrato server × fallback no-JS + enhancement HTMX`.
2. **Seleccionar el pattern del vocabulario** (`gelium-ui-pattern-registry.md`) con la regla de selección de `composition-rules.md` §4 (table vs list, list vs queue, feed vs collection, dialog vs page, toast vs inline vs banner vs callout…).
3. **Descomponer en componentes del registry** (`gelium-ui-component-registry.md`) — SOLO piezas existentes; la recipe es wiring, cero primitivas nuevas salvo aprobación.
4. **Resolver estados** — todo pattern de datos declara los 4 estados (rest/vacío/carga/error) con los state patterns (D): Empty state, Skeleton, Inline alert, Banner, Error state, Validation summary, Success feedback.
5. **Fijar el contrato server** (`composition-rules.md` §9): GET params estables para listados, POST+303 para mutaciones, 422 + `X-Gelium-Validation` para validación, `HX-Trigger gelium:toast` solo transitorio.
6. **Completar los 19 campos** (§3) y **verificar la matriz de consumidores** (§4).

## 2. Ejemplo real: Admin Resource descompuesto

> Recipe implementada: `/recipes/admin-resource` (`recipe_admin_resource.go`, `recipe-admin-resource.html`). Descomposición completa con referencias a registries y contratos.

| Capa | Pieza | Referencia | Rol en la pantalla |
|---|---|---|---|
| **Pattern** | Resource list (E3) + Search (E4) + Filters (E5) + Pagination (E6) + Destructive (E10) + Bulk (E11) + Confirmation (E18) + Error recovery (E9) + Notifications (E15) | `pattern-registry.md` §3, `ux-patterns.md` | el set de recursos con find/compare/select/act-per-row |
| **Data component** | Data table (`data-table.html`, `data_table.go`) | `component-registry.md` §2 | vehículo de lista: sort/filter/pagination/selection server-side |
| **State patterns** | Empty state (2 variantes: "No results"/"No projects yet"), Banner (`--success` persistente post-303), Inline alert + Validation summary (422), Error state (404), Toast (refresh) | `pattern-registry.md` §2 | todos los estados de la pantalla |
| **Action/input** | Text field (search + form), Select (estado), Button/Icon button, Dialog (confirm page variant) | `component-registry.md` §2 | búsqueda, filtros, crear/editar/borrar |
| **Recipe primitives** | Skeleton (refresh), Progress (refresh) | `component-registry.md` §2 | loading |
| **Server contract** | `GET ?q=&sort=&dir=&page=&selection=`; `POST+303`; `422 + X-Gelium-Validation`; `gelium:toast` solo refresh; `HX-Request` bifurca `#resource-panel`; refresh en `postOnlyPaths()` | `screen-recipes.md` §1, `composition-rules.md` §9 | el esqueleto de la interacción |
| **No-JS flow** | links GET reales + checkboxes nativos + form POST+303 + `<dialog open>` server-rendered + reload con toast inline | `screen-recipes.md` §1 NO_JS_FLOW | funciona sin JS |
| **HTMX enhancement** | `hx-get` swap de `#resource-panel`; refresh → fragmento + `HX-Trigger gelium:toast` | `screen-recipes.md` §1 HTMX_ENHANCEMENT | mejora sin cambiar contrato |
| **SEO/GEO** | `/recipes/*` = `noindex, nofollow`; title/description por ruta; canonical limpio; un solo `<h1>`; entidad única Gelium UI | `screen-recipes.md` §1 SEO/GEO_REQUIREMENTS | indexabilidad correcta |
| **Theme** | solo tokens `--ui-*` existentes; cero literales de color (`TestNoColorLiteralsInComponents` sobre `recipe-admin-resource.css`); light/dark/forced-colors por primitivos | `theme-registry.md`, `theme-contract.md` | apariencia sin trabajo extra |

**Lectura de la tabla**: la recipe es la suma de ~15 piezas del registry + 1 patrón de lista + 1 contrato server. Ninguna pieza es nueva; esto es lo que significa "recipe = 100% wiring".

## 3. Plantilla de composición (19 campos)

> Los 19 campos de una recipe, en orden de redacción (`screen-recipes.md` los implementa para Admin Resource, Ops Queue y Public Feed). Este es el **contrato de contenido** de `docs/gelium-ui-screen-recipes.md`.

| # | Campo | Qué captura | Fuente a consultar |
|---|---|---|---|
| 1 | SURFACE | contexto inmersivo (app shell/section/panel/card/overlay) | `composition-rules.md` §2 |
| 2 | USER | la tarea primaria del usuario, no el rol del sistema | — |
| 3 | PRIMARY_TASK | qué hace el usuario en esa superficie | — |
| 4 | SECONDARY_TASKS | tareas secundarias (filtrar/ordenar/paginar/crear…) | — |
| 5 | UX_PATTERN | patterns del registry que compone | `pattern-registry.md` §3 |
| 6 | VISUAL_VOCABULARY | componentes y su rol visual | `component-registry.md` §2 |
| 7 | COMPONENTS | lista exacta de piezas (cero nuevas) | `component-registry.md` |
| 8 | STATES | rest/empty/loading/error/selected/success | `pattern-registry.md` §2, `composition-rules.md` §8 |
| 9 | ACCESSIBILITY | aria-sort, aria-current, dialog nativo, 422 con focus | `accessibility-contract.md` |
| 10 | CONTENT_RULES | copy, botones = verbos, empty = mensaje + CTA | `content-rules.md` |
| 11 | SEO_REQUIREMENTS | noindex/noindex, title/desc por ruta, canonical, h1 | `seo-contract.md`, `seo-patterns.md` |
| 12 | GEO_REQUIREMENTS | contenido factual, URLs deep-linkables, entidad única | `geo-contract.md`, `geo-patterns.md` |
| 13 | SERVER_CONTRACT | GET params, POST+303, 422, gelium:toast, postOnlyPaths | `composition-rules.md` §9 |
| 14 | NO_JS_FLOW | rama no-HX completa | `lib/skills/14-component-implementation.md` |
| 15 | HTMX_ENHANCEMENT | fragmento swap, sin cambiar contrato de mutación | `composition-rules.md` §9 |
| 16 | RESPONSIVE_BEHAVIOR | columna fluida, min()/clamp(), paginación server-side | `composition-rules.md` §7 |
| 17 | THEME_REQUIREMENTS | solo tokens `--ui-*`, guard de literales | `theme-contract.md` |
| 18 | ALTERNATIVES_REJECTED | qué patterns se rechazaron y por qué | `composition-rules.md` §4-5 |
| 19 | RATIONALE | justificación contra composition-rules | `composition-rules.md` §11 |

## 4. Matriz recipe × componente consumido

> Quién consume qué. Fuente: `docs/gelium-ui-screen-recipes.md` (secciones COMPONENTS de cada recipe).

| Componente | Admin Resource | Ops Queue | Public Feed |
|---|---|---|---|
| Data table | ✅ | — | — |
| List | — | ✅ (two-line filas de cola) | — |
| Card | — | — | ✅ (unidad del feed, `ui-card-outlined`) |
| Avatar | — | ✅ (requester, sm, decorativo) | ✅ (autor, sm, decorativo) |
| Badge (+ tones) | — | ✅ (estado/SLA: error/success/warning/info) | ✅ ("New" info + counts) |
| Pagination standalone | — | ✅ | ✅ |
| Text field | ✅ (search + form) | — | — |
| Select | ✅ (estado) | ✅ (filtros status/kind) | — |
| Button / Icon button | ✅ | ✅ | ✅ (like, refresh) |
| Dialog (confirm) | ✅ (delete page variant) | — | — |
| Banner | ✅ (success post-303) | ✅ (success post-advance/dequeue) | — |
| Toast | ✅ (refresh) | ✅ (refresh) | ✅ (react flash + refresh) |
| Progress | ✅ (refresh) | ✅ (refresh) | ✅ (refresh) |
| Empty state | ✅ (2 variantes) | ✅ (2 variantes) | ✅ (por vista) |
| Skeleton | — (documentado) | — (documentado) | ✅ (placeholder servido) |
| Inline alert + Validation summary | ✅ (form 422) | — | — |
| Error state | ✅ (404) | ✅ (404) | ✅ (404) |
| Tabs | — | — | ✅ (vistas server-side, `aria-current`) |
| Notification/refresh | `POST .../refresh` (postOnly) | `POST .../refresh` (postOnly) | `POST .../refresh` (postOnly) |

**Matriz recipe × patrón UX** (qué pattern consumen, `pattern-registry.md` §3): Admin Resource → E3+E4+E5+E6+E9+E10+E11+E15+E18; Ops Queue → Queue (vocabulario) + E5+E6+E15; Public Feed → Feed (vocabulario) + E6+E7+E8+E15.

## 5. Mapa de desbloqueo inverso (qué pieza habilita qué recipe)

| Pieza | Recipes que desbloquea |
|---|---|
| Empty state (D1) | TODAS (Admin Resource, Ops Queue, Public Feed) |
| Avatar + Pagination standalone + Badge tone | Ops Queue, Public Feed |
| Data table + Dialog confirm + Validation summary | Admin Resource |
| Tabs | Public Feed (vistas) |
| Footer + Breadcrumb | chrome de TODAS las recipes (SEO §3, GEO §9) |
| Error state | 404 de las 3 recipes |
| Banner + Toast + Progress | feedback de las 3 recipes |

El detalle completo (con "sin consumidor directo hoy") está en `gelium-ui-dependency-metadata.md` §4.

---

**Definición de done (Phase J)**: flujo de composición en 6 pasos, ejemplo real (Admin Resource) descompuesto con referencias a los 4 registries, plantilla de 19 campos y matriz recipe × componente + recipe × pattern.
