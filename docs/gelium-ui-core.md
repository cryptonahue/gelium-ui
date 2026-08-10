# Gelium UI — Core

> Contrato del core y foundations del sistema Gelium UI.
> Fase A del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Este documento formaliza lo que el core ES, las capas, las foundations, los tokens y lo que NO pertenece al core. Complementa `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-composition-rules.md` y `docs/gelium-ui-theme-contract.md`.

---

## 1. Principios

1. **HTML-first, server-first**: cada componente se entrega como markup semántico server-rendered. El flujo principal funciona end-to-end **sin JavaScript**.
2. **Progressive enhancement**: HTMX (local, opcional) y JS mínimo solo mejoran la experiencia; nunca son requisito.
3. **Platform-first**: antes de agregar JavaScript se auditan HTML Living Standard, CSS moderno, Popover API, top layer, Invoker Commands y capacidades nativas relacionadas.
4. **Token-only theming**: los componentes leen `var(--ui-*)` en runtime; los themes definen tokens, no markup ni JS.
5. **Core agnóstico**: el core no contiene dirección visual (ni Material, ni Basecoat, ni shadcn). Solo roles, escalas y contratos.
6. **Accesibilidad nativa**: elementos nativos antes que ARIA; estado nunca color-only; foco en el control nativo; forced colors y reduced motion como ciudadanos de primera clase.
7. **Sin dependencias runtime innecesarias**: no React, no Lit, no Shadow DOM, no Custom Elements obligatorios, no CDN.
8. **Un integrador serial**: la integración al árbol canónico es serial; la investigación puede ser paralela.

## 2. Capas

El core es la base de un sistema de capas que NO se confunden:

| Capa | Definición | Dónde vive |
|---|---|---|
| **Core / foundations** | Primitivas de sistema: tokens, escalas, contratos transversales | Este documento + tokens `--ui-*` del core |
| **Visual vocabulary** | Nombres canónicos y anatomía de cada patrón | `docs/gelium-ui-vocabulary.md` |
| **Components** | Implementaciones concretas del vocabulario | `web/templates/*.html`, `web/styles/*.css` |
| **Composition patterns** | Reglas para elegir y combinar patrones | `docs/gelium-ui-composition-rules.md` |
| **Screen recipes** | Recetas de pantallas completas | Phase 4 del roadmap |
| **Server-driven contracts** | Contratos de interacción servidor-cliente | Sección 7 de este documento |
| **Themes** | Direcciones visuales intercambiables sobre el core | `themes/<theme>/theme.css` |
| **Registry / tooling** | Índices y prompts para agentes | Phase 8 del roadmap |

Regla operativa: un componente nuevo debe justificar qué pattern o screen recipe desbloquea; nada se agrega arbitrariamente.

## 3. Responsabilidades

**El core es responsable de:**

- Definir y documentar las familias de tokens `--ui-*` (roles, no valores de marca).
- Proveer escalas: color semántico, tipografía, spacing, radius, elevation, border, focus, motion, z-index, breakpoints, densidad.
- Definir contratos transversales: estado (hover/focus/pressed/disabled/selected/error), reduced motion, forced colors, dark mode, responsive.
- Ser el único lugar donde se declaran los roles semánticos que los themes pueblan.
- Mantener el contrato de clases `ui-*` y la semántica HTML de cada componente.

**El core NO es responsable de:**

- Valores de paleta, personalidad de shape, typescale o densidad elegida por un theme (eso es dirección visual).
- Implementar componentes; los componentes consumen el core.
- Decidir qué patrón usar (eso es composition rules).
- Renderizar pantallas (eso son screen recipes).

## 4. Foundations

| Primitiva | Estado actual (baseline) | Contrato del core |
|---|---|---|
| Color semántico | Existe `--ui-color-*` parcial (`theme.css:10-22`), faltan `success/warning/info/outline/scrim` | Roles: `canvas, surface, surface-container, primary, secondary, error, warning, success, info, outline, outline-strong, scrim, focus-ring` + variantes `-fg` (on-) |
| Tipografía | Existe `--ui-type-*` shorthand bundleada (`theme.css:48-58`), con huecos (`display-lg`, `title-md`) | Componer en `--ui-font-family-*`, `--ui-font-size-*`, `--ui-font-weight-*`, `--ui-line-height-*`, `--ui-letter-spacing-*` por step (display/headline/title/body/label) y reconstruir `--ui-type-*` como alias |
| Spacing | NO existe (`--ui-space-*` faltante; paddings en rem suelto) | Escala `--ui-space-0..N` (base 0.25rem) y uso en todos los paddings/gaps |
| Radius | Existe `--ui-radius-*` (`theme.css:61-67`), pero componentes usan px propios (dialog 28px, card 12px, fab, select 4px) | Normalizar px del theme como composición de la escala o token de componente con default |
| Elevation | Existe `--ui-shadow-0..5` (`theme.css:71-76`) | Pasa al core tal cual |
| Border | NO hay tokens de anchura/estilo | `--ui-border-width-*`, `--ui-border-style-*` |
| Focus | Existe `--ui-focus-*` + `.ui-focus-ring` | Pasa al core: thickness, offset, `:focus-visible`, forced colors |
| Density | NO existe; alturas fijas (`2.5rem`, `3.5rem`, `48px`, `56px`) | `--ui-density-{comfortable,compact,roomy}` + `--ui-size-{control,item,icon}` |
| Motion | Solo `--ui-motion-short` + `--ui-easing-standard` (`theme.css:199-200`); dialog usa literales 150/500ms | `--ui-motion-{short,medium,long}` + `--ui-easing-{standard,emphasized,decelerate,accelerate}` |
| Z-index / top layer | Solo `--ui-tooltip-z: 50`; toast hardcodea `z-index: 1000` | `--ui-z-{dialog,menu,popover,toast}`; jerarquía dialog > menu/popover > toast |
| Breakpoints | NO hay tokens | `--ui-breakpoint-{sm,md,lg}` |
| Icons | Existe `.ui-icon` 24px + contrato aria (`icon.css:5-10`) | SVG inline trusted con `aria-hidden="true"` + `focusable="false"`; nombre accesible fuera del SVG |
| Reduced motion | Por media query duplicada por archivo + central (`app.css:52-69`) | Estrategia única en el core, excepciones mínimas |
| Forced colors | Mismo patrón duplicado (`app.css:71-213` + por componente) | Estrategia única en el core |
| Dark mode | Duplicado a mano con drift (`theme.css:203-251` vs `253-299`) | Un solo mecanismo (`color-scheme` + un set de tokens por modo), sin duplicación |
| Responsive | Layouts fluidos por grid (`auto-fit/minmax`, `min()/clamp()`); sin breakpoints | Fluido primero; breakpoints solo cuando el layout fluido no alcanza |

## 5. Tokens

### 5.0 Token ownership (Phase A)

Cada token se clasifica en una de estas categorías:

| Tipo | Definición | Quién lo define | Quién lo sobrescribe | Ejemplo |
|---|---|---|---|---|
| **core token** | Contrato estructural del sistema; valor por defecto neutro | Core (`web/styles/` o `core.css`) | Theme (valor) | `--ui-space-4`, `--ui-focus-thickness`, `--ui-state-hover-opacity` |
| **theme token** | Valor de dirección visual | Theme (`themes/<theme>/theme.css`) | Otro theme | `--ui-color-primary`, `--ui-radius-md` |
| **component token** | Anatomía/rol específico de un componente | Component CSS (scoped) | Theme si el contrato lo permite | `--ui-dialog-scrim`, `--ui-card-radius` |
| **pattern token** | Valor compartido por un patrón (composición) | Pattern CSS | Theme si el contrato lo permite | `--ui-queue-item-height` (futuro) |
| **internal implementation token** | Detalle privado; no forma parte del contrato público | Componente | nadie | `--ui-slider-fill` (knob por instancia) |
| **deprecated/dead token** | Sin consumidor o reemplazado | — | — | `--ui-radius-xl`, `--ui-state-dragged-opacity` |

**Reglas de propiedad**:

1. El **core** declara roles y escalas con valores neutros; los **themes** pueblan valores de dirección visual. Un token core sin valor de theme usa su default.
2. Los **componentes** consumen tokens; declaran tokens propios solo para anatomía que el theme no debe cambiar (alturas de items, tamaños de iconos, paddings estructurales) o como defaults scoped que el theme PUEDE sobrescribir si el contrato lo declara.
3. Los **patterns** consumen componentes y tokens; no crean tokens nuevos sin consumidor real.
4. **Nada se convierte automáticamente en token público**: solo se tokeniza lo que tiene un consumidor real y un grado de libertad que un theme deba cambiar.
5. **internal implementation tokens** usan el mismo prefijo `--ui-` pero se documentan como privados y no forman parte del contrato; no se prometen en el theme contract.
6. **deprecated/dead tokens** se eliminan o se mantienen por compatibilidad con decisión explícita (ver §5.2).

**Matriz de decisión actual (Phase A)**:

| Token / familia | Owner actual | Owner propuesto | Consumidores | Acción |
|---|---|---|---|---|
| `--ui-color-*` (13 roles light/dark) | theme-material | **core** (roles) + theme (valores) | Todos los componentes | Extraer roles al core; Material queda como valores |
| `--ui-type-*` (10 shorthand) | theme-material | **core** (composición size/weight/line-height/letter-spacing) | Todos los componentes | Descomponer en core; mantener alias `--ui-type-*` |
| `--ui-radius-*` (7) | theme-material | **core** (escala) + theme (valores) | Componentes | Extraer escala; normalizar px propios (dialog 28px, card 12px) |
| `--ui-shadow-0..5` | theme-material | **core** | Elevation, Card, FAB, Toast | Pasar al core tal cual |
| `--ui-focus-*` | theme-material | **core** | Focus ring, todos | Pasar al core |
| `--ui-state-*` (opacidades) | theme-material | **core** (opacidades) | Todos los state layers | Pasar al core; colores del layer vía `color-mix()` |
| `--ui-motion-short`, `--ui-easing-standard` | theme-material | **core** | Dialog, Drawer, Menu | Extraer; añadir medium/long + easings con consumidor |
| `--ui-field-*` (5) | theme-material | **theme token** | Text field, Select | Queda en el theme (dirección visual de input) |
| `--ui-dialog-*` (5) | theme-material | **component token** | Dialog, Drawer modal | Queda como componente; theme puede sobrescribir color |
| `--ui-toast-*` (8) | theme-material | **component token** | Toast | Queda como componente; theme puede sobrescribir color |
| `--ui-card-*` (5) | theme-material | **component token** | Card | Queda como componente |
| `--ui-badge-*` (4) | theme-material | **component token** | Badge | Queda como componente |
| `--ui-checkbox/radio/switch/slider/progress/fab/select-*` | theme-material | **component tokens** | Respectivos | Quedan como componente; theme sobrescribe color/shape |
| `--ui-divider-*` | theme-material | **component token** | Divider | Queda como componente |
| `--ui-data-table-*` (scoped) | data-table.css | **component token** (con defaults en core para anatomía) | Data table | Subir defaults al core; theme puede sobrescribir color |
| `--ui-list-*` (scoped) | list.css | **component token** (ídem) | List | Ídem |
| `--ui-menu-*` (scoped) | menu.css | **component token** (ídem) | Menu | Ídem |
| `--ui-navigation-*` / `--ui-nav-bar-*` / `--ui-nav-tab-*` (scoped) | navigation-*.css | **component token** (ídem) | Navegación | Ídem |
| `--ui-segmented-button-*` (scoped) | segmented-button.css | **component token** (ídem) | Segmented | Ídem |
| `--ui-tooltip-*` (scoped) | tooltip.css | **component token** (ídem) | Tooltip | Ídem |
| `--ui-color-surface-container` | NO definido | **core token** (o alias) | base.css, demo | Definir o alias; cerrar gap |
| `--ui-type-display-lg`, `--ui-type-title-md` | NO definido | **core token** | base.css, demo | Definir o alias; cerrar gap |
| `--ui-font-mono` | NO definido | **core token** | demo | Definir |
| `--ui-color-error` vs `--ui-color-danger` | drift | **core token** `danger` + alias `error` | demo | Unificar a `danger`; alias de compatibilidad |
| `--ui-radius-xl` | theme-material (muerto) | **deprecated** | ninguno | Eliminar |
| `--ui-state-dragged-opacity` | theme-material (muerto) | **deprecated** | ninguno | Eliminar |
| `--ui-select-menu-item-icon` | theme-material (muerto) | **deprecated** | ninguno | Eliminar |
| `--ui-slider-fill` | slider.css | **internal implementation** | Slider | Documentar como privado |

### 5.1 Naming

Convención canónica: `--ui-<familia>-<token>` para familias transversales y `--ui-<componente>-<rol>` para cobertura por componente.

| Familia | Prefijo | Ejemplo |
|---|---|---|
| Color semántico | `--ui-color-*` | `--ui-color-primary`, `--ui-color-primary-fg` |
| Tipografía | `--ui-font-*`, `--ui-type-*` | `--ui-font-sans`, `--ui-type-body-md` |
| Spacing | `--ui-space-*` | `--ui-space-4` |
| Radius | `--ui-radius-*` | `--ui-radius-md` |
| Elevation | `--ui-shadow-*` | `--ui-shadow-2` |
| Border | `--ui-border-*` | `--ui-border-width-1` |
| Focus | `--ui-focus-*` | `--ui-focus-thickness` |
| Density | `--ui-density-*` | `--ui-density-comfortable` |
| Size | `--ui-size-*` | `--ui-size-control` |
| Motion | `--ui-motion-*`, `--ui-easing-*` | `--ui-motion-short`, `--ui-easing-standard` |
| State | `--ui-state-*` | `--ui-state-hover-opacity` |
| Z-index | `--ui-z-*` | `--ui-z-toast` |
| Breakpoint | `--ui-breakpoint-*` | `--ui-breakpoint-md` |
| Componente | `--ui-<componente>-*` | `--ui-dialog-scrim`, `--ui-card-radius` |

### 5.2 Reglas del contrato de tokens

1. **Públicos y exclusivos**: todo valor visual expuesto a themes usa `--ui-*`; nada de valores fijos en componentes para color/geometría de controles.
2. **Referenciados ≡ definidos**: ningún `var(--ui-*)` puede referenciar un token inexistente. Gaps conocidos a cerrar en Phase 1:
   - `--ui-color-surface-container` (referenciado en `base.css:40`, sin definir);
   - `--ui-type-display-lg` (`base.css:37`), `--ui-type-title-md` (`demo-whatsapp.css:477`);
   - `--ui-color-error` vs `--ui-color-danger` (drift de naming; se unifica);
   - `--ui-font-mono` (referenciado, sin definir).
3. **Sin tokens muertos**: eliminar o justificar `--ui-radius-xl`, `--ui-state-dragged-opacity`, `--ui-select-menu-item-icon`.
4. **Un solo mecanismo dark**: sin duplicación clase + media query; sin drift (`--ui-switch-track-unselected` solo en una ruta = bug conocido).
5. **State layers theme-aware**: los overlays de estado se pintan con `color-mix()` sobre el token `-fg` (patrón ya usado en `button.css:27-32`, `fab.css:80-88`), nunca con `rgb()` fijo (`button.css:17-18`, `chips.css:32,35,107,141-142` son acoplamientos a corregir).

### 5.3 Cobertura por componente

- Componentes cuyo token family vive en el theme (button, text-field, dialog, toast, card, badge, checkbox, radio, switch, slider, progress, fab, select, divider, icon, elevation, focus-ring).
- Componentes cuyo token family vive scoped en su CSS y el theme NO puede sobrescribir (list, menu, data-table, navigation-bar, navigation-tab, navigation-drawer, segmented-button, tooltip). **La propiedad de cada familia se decide en Phase A** (ver §5.0); el core define el esquema fijo (anatomía) y el theme solo los grados de libertad (color, shape, type, motion).

## 5.4 Theme identity (Phase A)

- **Eliminar el acoplamiento a `class="theme-material"`** hardcodeado en `web/templates/layout.html:2`.
- La identidad del theme debe ser configurable desde el layout o el servidor:

```html
<body class="theme-material">
<body class="theme-basecoat">
```

- La selección puede provenir del servidor (template data-driven) o del documento raíz.
- **No agregar todavía un runtime JavaScript de selección de themes.**
- Sin rebuild por theme: el bundle contiene todos los themes y la clase selecciona.

## 5.5 Semantic colors (Phase A)

Vocabulario semántico a nivel core (los themes definen valores; los componentes consumen tokens):

| Rol | Uso semántico |
|---|---|
| `--ui-color-primary` / `-fg` | Acción principal, elemento enfatizado |
| `--ui-color-secondary` / `-fg` | Acción/énfasis secundario |
| `--ui-color-success` / `-fg` | Éxito (no debe confundirse con primary) |
| `--ui-color-warning` / `-fg` | Advertencia |
| `--ui-color-danger` / `-fg` | Error/destructivo (canónico; `error` pasa a alias) |
| `--ui-color-info` / `-fg` | Informativo |
| `--ui-color-outline` / `-strong` | Bordes/contornos |
| `--ui-color-surface` / `-muted` / `-container` | Superficies |
| `--ui-color-content` / `-muted` | Texto/contenido |
| `--ui-color-border` | Separadores y bordes |
| `--ui-color-scrim` | Fondo de overlays |
| `--ui-color-focus-ring` | Anillo de foco |

Regla: los componentes referencian el **rol** (`var(--ui-color-success)`), nunca un valor; el theme define el valor por rol.

## 5.6 Dark mode (Phase A)

- **Una única rutina** para dark mode; sin duplicación, sin drift, sin selectores contradictorios, sin valores Material hardcodeados fuera del theme.
- Mecanismo acordado: cada theme autocontenido con su bloque dark (clase única o `light-dark()` según decisión del theme contract), eliminando la duplicación `theme.css:203-251` vs `253-299` (bug: `--ui-switch-track-unselected` solo en una ruta).
- `color-scheme` con una sola fuente de verdad.

## 5.7 Referencias rotas (Phase A)

Auditar y resolver; decisión por token:

| Referencia | Decisión | Tipo |
|---|---|---|
| `--ui-color-surface-container` | Definir en el core como rol de superficie | nuevo token |
| `--ui-type-display-lg`, `--ui-type-title-md` | Definir en la escala tipográfica del core | nuevo token |
| `--ui-font-mono` | Definir en el core | nuevo token |
| `--ui-color-error` | Renombrar a `danger` + alias de compatibilidad | rename + alias |
| `--ui-radius-xl`, `--ui-state-dragged-opacity`, `--ui-select-menu-item-icon` | Eliminar (sin consumidores) | eliminación |

## 5.8 State layers (Phase A)

- **Eliminar colores fijos `rgb(...)`** cuando representen hover/focus/pressed/dragged.
- Implementar state layers **theme-aware** usando tokens y `color-mix()` cuando corresponda (patrón ya usado en `button.css:27-32`, `fab.css:80-88`).
- Contrato de estados:

| Estado | Mecanismo |
|---|---|
| hover | `color-mix(in srgb, <fg-token> calc(var(--ui-state-hover-opacity)*100%), transparent)` |
| focus | `:focus-visible` + `--ui-focus-*` + ring |
| pressed | `color-mix(... var(--ui-state-pressed-opacity) ...)` |
| dragged | solo si existe consumidor (hoy no); si no, no tokenizar |
| disabled | `var(--ui-state-disabled-opacity)` sobre el control |
| selected | `:checked`/`aria-current`/`aria-sort` + color de selección del theme |

- Acoplamientos a corregir: `button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`.

## 6. HTML-first

Reglas del markup (derivadas de la implementación, `web/templates/*.html`):

1. Elementos nativos antes que ARIA: `button`, `a`, `ul/ol/li`, `table/caption/th/thead/tbody`, `dialog`, `progress`, `input[type=checkbox|radio|range]`, `select`, `fieldset/legend`, `label`, `hr`.
2. Sin `div`/`span` como controles. Raíz semántica elegida por la acción: `<article>/<a>/<button>` para Card, `<nav>` para navegación, `<ul>` para listas y menús.
3. SVG decorativo con `aria-hidden="true"` + `focusable="false"`; el texto visible (Label) aporta el accessible name.
4. Estado nunca color-only: `:checked`, `aria-sort`, `aria-current`, `role`, `aria-invalid`, `aria-busy` son portadores de estado junto con el color.
5. Foco en el control nativo, con focus ring visible vía `:focus-visible` + forced colors.
6. Los templates NO usan utilidades Tailwind (0 clases `bg-*`/`text-*`): el theming es 100% cascade-time por `var(--ui-*)`.

## 7. Server-first — contratos server-driven

El core define y reutiliza estos contratos (evidencia en `composition-audit.md` §3); las screen recipes NO inventan otros:

1. **HTTP 422 + `X-Loom-Validation: true`** — validación server-side de campos/valores.
   - Sin HX: re-render de página completa con 422, preservando valor + `aria-invalid` + mensaje.
   - Con HX: `htmx:beforeSwap` (`web/static/app.js:1-9`) swapea solo 422 con header.
   - Regla: **la validación nunca dispara toast** (`toast.go:129-133`).
2. **`HX-Trigger: {"loom:toast":{"type":"success","message":"…"}}`** — feedback transitorio server-driven.
   - Vocabulario cerrado `info|success|warning|error`; `error` → `role="alert"`.
   - Región `#loom-toast-region` con `aria-live="polite"`.
   - Sin JS: toast inline persistente; con JS: auto-dismiss 4s/8s pausable.
3. **GET con parámetros estables** (`?q=&sort=&dir=&page=&selection=`) — estado de listados (Data table).
   - La URL es el estado: no-JS = reload completo; HTMX = fragmento `outerHTML`.
   - Vocabularios cerrados sanitizados; orden estable de params + escape.
   - `HX-Request: true` bifurca en el handler: fragmento vs página completa.
4. **POST + 303 SeeOther redirect** — "mover estados" de workflow simple (demo WhatsApp): sin fragmentos, sin JS.

## 8. Progressive enhancement

- HTMX local (`web/static/htmx.min.js`, copiado desde npm `htmx.org`, sin CDN) y opcional.
- El hook local `app.js` solo intercambia fragments HTTP 422 con `X-Loom-Validation: true`.
- Primitivas declarativas con fallback: `<dialog>` + Invoker Commands (`command`/`commandfor`) + `closedby="any"` (navegadores compatibles); Popover API (`popover`/`popovertarget`) para menús; Interest Invokers evaluados y rechazados para tooltip por no-Baseline.
- Toda mejora JS tiene un fallback server-rendered real; nunca se rompe el flujo principal sin JS.

## 9. Accesibilidad

- Contraste AA por diseño en light y dark.
- Reduced motion: `prefers-reduced-motion` central en el core.
- Forced colors: `forced-colors` central en el core.
- Nombres accesibles obligatorios en controles icon-only (Icon button, FAB).
- `aria-live` para feedback transitorio (toast); `role="alert"` para errores persistentes ligados al contexto.
- Foco visible, orden lógico del documento, sin roving focus falso (links reales con `aria-current` en Tabs).

## 10. Qué NO pertenece al core

| Fuera del core | Por qué |
|---|---|
| Valores de paleta (hex Material/Basecoat) | Dirección visual del theme |
| Personalidad de shape (dialog 28px, pill buttons) | Dirección visual del theme |
| Typescale Material y densidad por defecto | Dirección visual del theme |
| Variantes de componente y su CSS | Capa components |
| Criterios de elección de patrones | Capa composition rules |
| Recetas de pantallas | Capa screen recipes |
| JavaScript de componente | Solo como enhancement justificado y auditado platform-first |
| Clases `m3-*` y prefijos no-`ui-*` | Contaminación del contrato (renombrar `m3-select-trigger` → `ui-select-trigger`) |
| `class="theme-material"` hardcodeado en `layout.html` | La selección de theme es decisión del theme contract / runtime |

---

**Definición de done (Phase A/B)**: core documentado con token ownership; tokens `--ui-*` inventariados y gaps cerrados; state layers theme-aware; dark unificado; theme identity data-driven; sin acoplamiento Material en el core; `npm run build` + `go test ./...` + `go vet ./...` verdes; smoke light/dark, narrow/wide, reduced motion, forced colors, no-JS, HTMX; aceptación visual manual.
