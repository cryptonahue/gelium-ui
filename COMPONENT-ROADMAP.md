# Loom UI — Component Roadmap

> Estado operativo del catálogo y orden recomendado de implementación. Este documento complementa la auditoría histórica `MATERIAL-WEB-PROGRESS.md`; no la reemplaza ni modifica.

## Objetivo

Construir una librería open-code, server-rendered y progresivamente mejorada basada en Go `net/http`, `html/template`, `embed`, Markdown, Tailwind CSS 4, themes `--ui-*`, HTML/CSS moderno y HTMX opcional.

Reglas del producto:

- Material Web upstream define por defecto anatomía, geometría, estados, interacción y jerarquía visual.
- Loom cambia la implementación, no simplifica el diseño sin una decisión explícita.
- Antes de agregar JavaScript se auditan HTML Living Standard, CSS moderno, formularios, Popover API, top layer, Invoker Commands y capacidades nativas relacionadas.
- Todo flujo acordado debe funcionar end-to-end sin JavaScript dentro de la matriz de navegadores declarada. HTMX sólo mejora la experiencia.
- No se usan React, Lit, Shadow DOM, Astro, `templ`, Custom Elements obligatorios, CDN ni dependencias runtime de Material Web.
- Cada componente se entrega como slice vertical completo: contrato, TDD, theme, docs dogfoodeadas, reviews, smoke y aceptación manual.

## Estado actual

| Área | Estado | Evidencia principal |
|---|---|---|
| Foundation Go + Markdown + embed | Completado | `cmd/loom`, `internal/app`, `web/assets.go` |
| Tailwind CSS 4 + HTMX local | Completado | `package.json`, `package-lock.json`, `web/static/` |
| Theme Material light/dark | Completado | `themes/theme-material/theme.css` |
| Button | Completado y aceptado | `web/templates/button.html`, `web/content/button.md` |
| Text field | Completado y aceptado | `web/templates/text-field.html`, validación HTTP 422/200 |
| Dialog | Completado y aceptado | `web/templates/dialog.html`, Invoker Commands, `<dialog>` |
| Toast | Completado y aceptado | `web/templates/toast.html`, `HX-Trigger loom:toast`, región `aria-live` en `web/static/app.js` |
| Divider | Completado y aceptado | `web/templates/divider.html`, tokens `--ui-divider-*`, route `/components/divider` |
| Card | Completado y aceptado | `web/templates/card.html`, tokens `--ui-card-*`, route `/components/card` |
| Badge | Completado y aceptado | `web/templates/badge.html`, tokens `--ui-badge-*`, route `/components/badge` |

`Field` upstream no debe publicarse como componente separado: es una primitive interna para Text field y Select.

## Roadmap por waves

### Wave 0 — inmediata

#### Toast — Loom-only, P0 ✅ Completado y aceptado

Objetivo: feedback transitorio reusable para acciones server-driven.

Implementación cerrada con:

- wire contract congelado y testeado: `HX-Trigger: {"loom:toast":{"type":"success","message":"..."}}`;
- región `aria-live` (`role="status"`/`alert` según tipo) en `web/static/app.js` + listener global `loom:toast`;
- auto-dismiss pausable (hover/focus), dismiss manual, closed vocabulary `info|success|warning|error`;
- flujo no-JS real: re-render de página completa con toast inline persistente;
- validación nunca anuncia toast: HTTP 422 sin `HX-Trigger`;
- assets: `web/static/app.js`, `web/static/app.css`, tokens `--ui-toast-*` en el theme.

### Wave P — habilitación de trabajo paralelo

No cuenta como componente. Debe ejecutarse antes de permitir que varios implementadores escriban simultáneamente sobre el árbol canónico.

1. Dividir `web/styles/app.css` en base, docs y un archivo por componente.
2. Dividir `internal/app/server.go` y sus tests por componente manteniendo `package app`.
3. Dividir tests CSS por componente y theme.
4. Separar tokens del theme por familia y mantener `theme.css` como entrypoint de imports.
5. Convertir navegación, rutas y previews en un registro/integración estable.
6. Mantener `web/static/app.css` y `web/static/htmx.min.js` como outputs exclusivos de la build lane: `npm run build` sobrescribe ambos.
7. Congelar la API reusable de Button: link/button, disabled/loading, icon trusted slot, `Command`, `CommandFor`, `Value`, `Autofocus`.

Hasta completar esta wave, la investigación puede ser paralela, pero la integración al checkout canónico debe ser serial.

### Wave 1 — foundations y componentes estáticos

#### 1A. Primitives de sistema

| Primitive | Estrategia | JavaScript |
|---|---|---|
| Focus ring | `:focus-visible`, tokens y forced colors | No |
| Elevation | tokens y sombras CSS | No |
| Icon | SVG inline trusted y contrato accesible | No |

No son tres widgets complejos, sino contratos compartidos que desbloquean el catálogo.

#### 1B. Estáticos

| Componente | Base semántica | Dependencias |
|---|---|---|
| Divider ✅ | `<hr>` o separador semántico contextual | Theme |
| Badge (Labs) ✅ | texto/contador asociado; no color-only | Theme, Icon opcional |
| Card (Labs) ✅ | `<article>`, `<a>` o `<button>` según acción | Theme, Elevation, Focus |

Wave 1B cerrada: los tres componentes están entregados con contrato, TDD, theme, docs dogfoodeadas, smoke y aceptación.

Estos componentes pueden prepararse en paralelo en workspaces aislados después de congelar foundations.

### Wave 2 — controles de formulario y acciones

#### 2A. Forms nativos

| Componente | Base nativa | Alcance inicial |
|---|---|---|
| Checkbox | `input[type="checkbox"]` | checked, unchecked, disabled, error; indeterminate sólo si se justifica JS |
| Radio | `input[type="radio"]` + fieldset/legend | grupo, selected, disabled, error |
| Switch | checkbox nativo con semántica de switch cuando corresponda | on/off, disabled, labels |
| Select | `<select>` nativo primero | filled/outlined, helper/error, no select enriquecido inicial |
| Slider | `input[type="range"]` | single value; dual range diferido |

Submit y validación deben funcionar con HTML normal. HTMX es opcional.

#### 2B. Estado y acciones

| Componente | Dependencias | Estrategia | Estado |
|---|---|---|---|
| Progress | Theme | `<progress>` para determinate; indeterminate accesible; SSE/polling sólo para progreso remoto | ✅ Entregado y aceptado |
| Icon button | Icon + Focus + Button contracts | `<button>`/`<a>` con nombre accesible obligatorio | ✅ Entregado y aceptado |
| FAB | Icon + Focus + Button contracts | acción primaria flotante, semántica nativa | ✅ Entregado y aceptado |

### Wave 3 — composición

#### List

- `<ul>`, `<ol>`, `<nav>` o colección semántica según contenido.
- Depende de Icon y Focus.
- Debe distinguir navegación, selección y contenido estático.
- ✅ Entregado y aceptado (Wave 3, commit 2352f53): one/two/three-line, navegación con `<a href>`, selección multi con checkboxes nativos no-JS, iconos leading/trailing.

#### Chips

Familia upstream: assist, filter, input y suggestion chips. Debe implementarse por variantes aprobadas, no como un único `div` genérico.

- links, buttons o checkboxes según semántica;
- remove/selection puede requerir round-trip server-side o JS mínimo;
- depende de Icon, Icon button, Button y Focus.
- ✅ Entregado y aceptado (Wave 3, commit 2352f53): assist/suggestion como buttons/links, filter como checkboxes nativos, input chip con remoción server-side en `POST /examples/chips/remove`.

### Wave 4 — navegación y selección agrupada

#### 4A. Navegación

| Componente | Dependencias | Base/fallback | Estado |
|---|---|---|---|
| Tabs | Focus | links/páginas como fallback; tablist ARIA sólo si teclado completo está resuelto | ✅ Entregado (Wave 4A, commit 8e7c741): links reales `<a href>` con selección server-side `aria-current="page"`, primary/secondary, sin roving focus (sin brecha JS demostrada) |
| Navigation bar (Labs) | Icon, Focus, Badge opcional | navegación real con links | ✅ Entregado (Wave 4A, commit 8e7c741): `<nav>` + destinos `<a href>` con activo server-side, reusa `.ui-badge` |
| Navigation tab (Labs) | contrato Tabs/Nav bar | link semántico, no tab falso | ✅ Entregado (commit 8ca17a9): `<a href>` real con activo server-side, reusa `.ui-nav-bar` y `.ui-badge` |

#### 4B. Segmented buttons

| Componente | Dependencias | Base nativa | Estado |
|---|---|---|---|
| Segmented button (Labs) | Button, Focus | radio/checkbox/button según selección | ✅ Entregado (commit 8ca17a9): radio/checkbox/button nativos, `:checked` sin JS |
| Segmented button set (Labs) | Segmented button | fieldset o grupo accesible | ✅ Entregado (commit 8ca17a9): `<fieldset>`+`<legend>` para selección, `role="group"` para acciones |

Single/multi select debe preferir radios/checkboxes sin JS.

### Wave 5 — overlays y navegación compleja

#### Menu

- depende de List, Icon y Focus;
- auditar Popover API, `popovertarget`, anchor positioning e Invoker Commands;
- top layer/open-close puede ser declarativo;
- posicionamiento, typeahead y teclado completo sólo con JS mínimo si la plataforma no alcanza;
- fallback no-JS: navegación o formulario real, no imitación CSS inaccesible.
- ✅ Entregado (Wave 5, commit 4e1631c): Popover API declarativa (`popover` + `popovertarget`, light-dismiss/Escape nativos), anchor positioning con fallback, items como links/buttons/checkbox/radio nativos, zero JS.

#### Navigation drawer (Labs)

- depende de Dialog y primitives de navegación;
- variante modal sobre `<dialog>`;
- variante permanente como `<nav>` en layout;
- HTMX sólo para contenido remoto, no para semántica básica.
- ✅ Entregado (Wave 5, commit 4e1631c): modal sobre `<dialog>` nativo, permanente como `<nav>`, destinos `<a href>` con activo server-side, reusa `.ui-badge`.

### Wave 6 — Loom-only posteriores

#### Data table

- tabla HTML nativa;
- filtros, sort y paginación server-side;
- HTMX como enhancement;
- Progress y Toast para operaciones remotas;
- SSE opcional para actualización realtime.

#### Tooltip

- auditar primero `title`, Popover API, interest invokers, anchor positioning y soporte actual;
- nunca esconder información esencial exclusivamente en tooltip;
- fallback accesible visible o `aria-describedby`;
- JS sólo para una brecha real de posicionamiento/interacción.

### Deferred — Ripple

Ripple es Core upstream, pero queda explícitamente diferido.

- Dependencias: Button, Focus ring y state-layer tokens.
- Loom ya cubre hover/focus/pressed mediante state layers CSS sin JavaScript.
- Un ripple posicional fiel necesita coordenadas del puntero, lifecycle y JavaScript; no aporta semántica adicional.
- Sólo se reconsiderará si un requisito de producto exige ripple posicional y la auditoría platform-first demuestra que CSS moderno no puede reproducirlo correctamente.
- Si se aprueba, será enhancement visual framework-free: nunca requisito para activación, foco o feedback accesible.

## Inventario upstream pendiente

### Core distribuibles

- Checkbox
- Chips
- Elevation
- FAB
- Focus ring
- Icon
- Icon button
- List
- Menu
- Progress
- Radio
- Ripple
- Select
- Slider
- Switch
- Tabs

### Core completados

- Button
- Dialog
- Divider
- Text field

### Core interno/no publicar

- Field

### Labs candidatos

- Navigation bar
- Navigation drawer
- Navigation tab
- Segmented button
- Segmented button set

### Labs completados

- Badge
- Card

### Labs diferidos/no portar directamente

- Item: contrato inestable/interno.
- ARIA primitives: referencia de comportamiento, no componente visual distribuible.
- Behaviors: infraestructura upstream, no runtime Loom.
- GB: experimental; no portar en bloque.

### Loom-only

- Toast — P0.
- Data table — posterior.
- Tooltip — posterior y platform-first.

No crear un componente Snackbar separado: usarlo sólo como referencia visual para Toast.

## Dependencias

```text
Theme/foundations
├─ Focus ring
├─ Elevation
├─ Icon
│  ├─ Icon button
│  ├─ FAB
│  ├─ Chips
│  ├─ Menu
│  └─ Navigation components
├─ Text field ── Select visual contract
├─ Button ────── Dialog actions / Toast dismiss / Chips
│  └─ Ripple (deferred; state layers CSS remain default)
├─ List ──────── Menu
├─ Tabs ──────── Navigation tab contract
├─ Segmented button ── Segmented button set
├─ Dialog ────── Navigation drawer modal
└─ Toast + Progress + Forms ── Data table operations
```

## Qué puede hacerse en paralelo

### Seguro hoy

- auditorías upstream independientes;
- auditorías HTML/CSS/Baseline;
- extracción de contratos y matrices de estados;
- preparación en copias físicas aisladas;
- archivos nuevos exclusivos del componente;
- reviews read-only.

### No seguro hoy sobre `D:\repos\loom-ui`

- dos agentes editando `server.go`, `server_test.go` o `styles_contract_test.go`;
- dos agentes editando `web/styles/app.css` o el theme;
- un build Tailwind mientras otro agente edita CSS;
- Button y un consumidor de Button en paralelo;
- cualquier componente en paralelo con foundations/theme;
- dos agentes reiniciando servidores o reemplazando `loom.exe`.

### Modelo operativo recomendado

1. Un integrador único posee el árbol canónico.
2. Workers paralelos usan `SHARED_HANDOFF` o copias físicas aisladas.
3. Cada worker entrega artifacts exclusivos y solicitudes acotadas para archivos compartidos.
4. El integrador aplica una lane por vez.
5. Cada integración publica un baseline SHA-256 nuevo.
6. El coordinador anota `STALE_REBASE_REQUIRED` sobre cualquier handoff previamente emitido cuyo baseline cambió; no es un estado final del worker. Un worker que descubre drift durante la ejecución termina `ABORTED_ON_DRIFT`.
7. Build, tests, release y smoke son lanes seriales.

## Definition of Done por componente

Un componente no está terminado hasta cumplir todo:

1. upstream Material concreto inspeccionado con rutas literales;
2. auditoría platform-first actual y matriz de compatibilidad;
3. contrato visual/interactivo/accesible aprobado antes de código;
4. divergencias explícitas aprobadas;
5. tests RED observados antes de producción;
6. implementación GREEN mínima y refactor seguro;
7. estados y combinaciones probados;
8. tokens públicos exclusivamente `--ui-*`;
9. light/dark, narrow/wide, RTL cuando aplica;
10. reduced motion y forced colors;
11. flujo principal sin JS;
12. HTMX sólo como enhancement;
13. docs dogfoodeadas usando el componente real;
14. build reproducible y asset versionado;
15. spec review PASS;
16. quality review APPROVED;
17. smoke en `:8788`, HTTP y browser real;
18. consola sin errores;
19. `:8787` y `loom.exe` intactos durante implementación;
20. checklist manual aceptada por Nahuel antes de avanzar.

## Trabajo de sistema posterior

Después de Toast y primitives:

1. registry reusable de componentes/docs;
2. SQLite WAL + migrations;
3. broker Go y SSE con heartbeat/cancelación;
4. demo realtime en dos pestañas;
5. themes adicionales;
6. CLI y adapters.

Estos elementos son infraestructura; no deben contarse como componentes terminados.

## Fuentes locales de autoridad

- Auditoría histórica: `MATERIAL-WEB-PROGRESS.md` — read-only salvo tarea explícita.
- Prompt de auditoría: `PROMPT-MATERIAL-WEB-INVENTORY.md` — read-only.
- Upstream: `D:\repos\material-web-upstream` — read-only.
- Fork de referencia: `D:\repos\material-web-tailwind` — read-only.
- Estado real actual: código, tests y docs en `D:\repos\loom-ui`.
- Prompt operativo de workers: `AI-COMPONENT-IMPLEMENTER-PROMPT.md`.
