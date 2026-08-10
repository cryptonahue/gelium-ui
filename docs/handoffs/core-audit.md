# Core Audit — Gelium UI (ex Gelium UI)

> Handoff read-only de investigación. No modifica código, templates, CSS ni tests.
> Alcance: tokens, acoplamientos Material, foundations, contrato de tokens y build CSS del CORE.
> Fecha: 2026-08-10. Fuentes leídas: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `themes/theme-material/theme.css`, `web/styles/*.css` (31 archivos), `package.json`, `scripts/copy-htmx.mjs`, `web/templates/layout.html`, `web/styles_contract_test.go` y una muestra de `web/styles_*_test.go`.

---

## 1. Resumen ejecutivo

El sistema tiene un único theme (Material) con tokens públicos `--ui-*` que crecieron orgánicamente y **no hay core agnóstico todavía** (confirmado por el propio roadmap, `docs/gelium-ui-system-roadmap.md:33`). El theme define ~240 declaraciones de tokens en 3 esquemas (light, dark explícito, dark por media query); 8 componentes definen además ~129 tokens *scoped a su raíz* que el theme NO sobrescribe.

Los tokens están bien encaminados (semánticos, con light/dark, reduced motion y forced colors casi universales), pero hay cuatro problemas estructurales para un core agnóstico:

1. **Familias de foundation faltantes por completo**: spacing, density, breakpoints, escala de motion, z-index/top-layer, border (anchura/estilo) y colores semánticos complementarios (success/warning/info/outline). Sin ellas, Basecoat no puede re-skinear sin tocar componentes.
2. **Geometría Material incrustada en los componentes** como literales px (`48px`, `56px`, `3.5rem`, `24px`) y **state-layer colors hardcodeados** (`rgb(255 255 255 …)` en Button, `rgb(25 28 28 …)` en Chips) en vez de derivarse de tokens.
3. **Tokens inexistentes referenciados** con fallback silencioso: `--ui-color-surface-container` (`web/styles/base.css:40`), `--ui-type-display-lg` (`base.css:37`), `--ui-type-title-md` (`demo-whatsapp.css:477`), `--ui-color-error` (`demo-whatsapp.css:403`), `--ui-font-mono` (`demo-whatsapp.css:520,528`). Y tokens muertos en el theme: `--ui-radius-xl`, `--ui-state-dragged-opacity`, `--ui-select-menu-item-icon`.
4. **Acoplamiento a la identidad Material en markup y clases**: `class="theme-material"` hardcodeado en `web/templates/layout.html:2`, y la clase `m3-select-trigger` con prefijo m3- en `select-menu.css:23,75` y `app.css:63,174,180`.

El build (Tailwind CSS 4) compila `web/styles/app.css` → `web/static/app.css` y hardcodea la ruta al theme Material (`app.css:2`); para multi-theme hay que parametrizar esa integración. Los tests de CSS ya acoplan valores literales del theme (`web/styles_toast_test.go:33-47`), por lo que la extracción del core romperá tests y hay que planificarlos como tests de contrato.

---

## 2. Inventario de tokens actuales

### 2.1 Familias definidas en `themes/theme-material/theme.css` (scoped a `.theme-material`)

| Familia | Prefijo | Representativos | Líneas |
|---|---|---|---|
| Color semántico (light) | `--ui-color-*` | canvas `#fff7ff`, surface `#f7f2fa`, fg `#1d1b20`, fg-muted `#49454f`, primary `#6750a4`, primary-fg `#ffffff`, secondary `#e8def8`, secondary-fg `#1d192b`, border `#cac4d0`, border-strong `#79747e`, danger `#b3261e`, danger-fg `#ffffff`, focus-ring `#4f378b` | `theme.css:10-22` |
| Color semántico (dark explícito `.theme-dark`) | ídem | canvas `#141218`, primary `#d0bcff`, primary-fg `#381e72`, danger `#f2b8b5`, focus-ring `#e8def8` | `theme.css:203-251` |
| Color semántico (dark auto, media query) | ídem | mismo set + override por componente | `theme.css:253-299` |
| Text field | `--ui-field-*` | container `#f3edf7`, border `#79747e`, border-hover `#1d1b20`, label, error `#b3261e` | `theme.css:25-29` |
| Dialog | `--ui-dialog-*` | container `#ece6f0`, fg, body, scrim `rgb(0 0 0/.32)` | `theme.css:32-35` |
| Toast | `--ui-toast-*` | container `#322f35`, fg `#f3edf7`, radius `4px`, action `#d0bcff`, icon-{info,success,warning,error} | `theme.css:38-45` |
| Tipografía | `--ui-font-*`, `--ui-type-*` | font-sans; display-sm, headline-sm, title-lg, body-lg/md/sm, label-lg/sm, dialog-headline/body (shorthand font única) | `theme.css:48-58` |
| Shape/radius | `--ui-radius-*` | none 0, xs .25rem, sm .5rem, md .75rem, lg 1rem, xl 1.75rem, full 9999px | `theme.css:61-67` |
| Elevation | `--ui-shadow-0..5` | shadow-1 `0 1px 2px rgb(0 0 0/.14), 0 1px 3px rgb(0 0 0/.10)` | `theme.css:71-76` |
| Divider | `--ui-divider-*` | color (alias border), thickness 1px | `theme.css:79-80` |
| Card | `--ui-card-*` | radius 12px, container-elevated/filled/outlined, outline-color | `theme.css:83-87` |
| Badge | `--ui-badge-*` | size 6px, large-size 16px, container (danger), fg | `theme.css:90-93` |
| Checkbox | `--ui-checkbox-*` | size 18px, radius 2px, outline-width 2px, container (primary), icon | `theme.css:96-104` |
| Radio | `--ui-radio-*` | size 20px, radius 50%, outline-width 2px, checked (primary) | `theme.css:107-113` |
| Switch | `--ui-switch-*` | width 52px, height 32px, handle sizes 16/24/28px, track opacities .12/.38 | `theme.css:116-130` |
| Slider | `--ui-slider-*` | track-height 4px, handle 20/24px, active/inactive, disabled-opacity | `theme.css:133-142` |
| Progress | `--ui-progress-*` | track-height 4px, track `#e6e0e9`, indicator (primary) | `theme.css:145-148` |
| FAB | `--ui-fab-*` | container-shape 16/12/28/20px, icon 24/36/28px, extension-gap 12px, containers primary/surface/secondary | `theme.css:151-164` |
| Select | `--ui-select-*` | height 3.5rem, radius 4px, container-filled `#e6e0e9`, menu-container `#ece6f0`, menu-item-height 48px | `theme.css:167-189` |
| State (overlay) | `--ui-state-*` | hover .08, focus .10, pressed .10, dragged .16 (sin uso), disabled .38 | `theme.css:192-196` |
| Focus | `--ui-focus-*` | thickness 3px, offset 2px | `theme.css:197-198` |
| Motion | `--ui-motion-*`, `--ui-easing-*` | motion-short 150ms; easing-standard `cubic-bezier(.2,0,0,1)` | `theme.css:199-200` |

### 2.2 Tokens scoped a la raíz del componente (el theme NO los define ni sobrescribe)

| Componente | Tokens | Líneas |
|---|---|---|
| Data table | container/outline/header-height 56px/row-height 52px/cell-padding 16px/hover-focus-pressed/checkbox 18px/sort-icon 18px | `data-table.css:12-25` |
| List | leading/trailing 16px, icon 24px, heights 56/72/88px, opacities | `list.css:11-24` |
| Menu | container-color/radius/elevation, item-height 48px, spaces 12px, icon 24px, selected | `menu.css:12-29` |
| Tooltip | container/fg/radius/padding/supporting/max-width/offset/z `50` | `tooltip.css:21-28` |
| Navigation bar | height 80px, icon 24px, icon-container 64×32px, indicator 32px, opacities | `navigation-bar.css:16-35` |
| Navigation tab | icon 24px, icon-container 64×32px, indicator 32px, label-height 16px | `navigation-tab.css:17-33` |
| Navigation drawer | width 360px, item-height 56px, padding 12px, icon 24px, opacities | `navigation-drawer.css:17-34` |
| Segmented button | container-height 40px, radius-full, icon 18px, gaps 8/12px | `segmented-button.css:15-26` |

### 2.3 Familias que faltan para un core completo

- **Spacing scale** (`--ui-space-*`): no existe. Todo padding/gap es rem suelto (`.25rem`, `.5rem`, `1rem`, `1.5rem`).
- **Density** (`--ui-density-*` / `--ui-size-*`): no existe. Las alturas de controles son literales (`2.5rem` button.css:5, `3.5rem` text-field.css:10 y select.css:167 del theme).
- **Motion scale**: solo existe `--ui-motion-short` (150ms) y `--ui-easing-standard`. Dialog usa literales `150ms`/`500ms` (`dialog.css:19,37,39`). No hay duraciones medium/long ni easings emphasized/decelerate/accelerate.
- **Z-index / top layer**: solo existe `--ui-tooltip-z: 50` (scoped, `tooltip.css:28`). Toast hardcodea `z-index: 1000` (`toast.css:22`). No hay escala de overlays (dialog > menu/popover > toast).
- **Border**: no hay tokens de anchura/estilo (solo `--ui-divider-thickness` y `--ui-color-border-*`). Bordes de 1px hardcodeados en componentes.
- **Breakpoints**: no hay tokens. Media queries hardcodeadas (`max-width: 48rem` en `text-field.css:160`, `base.css` no tiene escala).
- **Colores semánticos complementarios**: no hay `--ui-color-success/warning/info/outline/scrim` a nivel core; solo `--ui-toast-icon-*` (por estado) y `--ui-dialog-scrim`.
- **Reduced motion / forced colors**: no son tokens; son bloques `@media` repetidos por archivo + central en `app.css:52-69,71-213`.
- **Dark mode**: existe, pero duplicado manualmente (bloque `.theme-dark` + bloque media query) y con drift (ver 3.6).
- **Tipografía composable**: los `--ui-type-*` son shorthands bundleadas (familia+peso+tamaño+alto); no se puede overridear solo el peso o el tamaño (ver 4.2).
- **Font mono**: `--ui-font-mono` referenciado pero sin definir.

---

## 3. Acoplamientos Material (dónde los estilos asumen decisiones visuales Material)

### 3.1 State-layer colors hardcodeados (rompen con cualquier palette que no sea la de Material)

- `web/styles/button.css:17-18` — `box-shadow: inset 0 0 0 999px rgb(255 255 255 / var(--ui-state-hover-opacity))`. Asume on-primary blanco; ignora `--ui-color-primary-fg`. (El mismo botón ya usa `color-mix(... var(--ui-color-primary) ...)` en `:27-31`, patrón correcto a replicar.)
- `web/styles/icon-button.css:22,25` — mismo `rgb(255 255 255 / …)`.
- `web/styles/chips.css:32,35,141,142` — `rgb(25 28 28 / …)` = on-surface Material `#191c1c` hardcodeado.
- `web/styles/chips.css:107` — `rgb(29 25 43 / …)` = on-secondary-container Material `#1d192b` hardcodeado.
- `web/styles/menu.css:87` y `data-table.css:79` — comentarios y `::before` pintan el state layer con `currentColor`; el color base es correcto, pero depende de que los colores de texto del tema sigan el contraste Material.

### 3.2 Geometría Material como literales en componentes

- `web/styles/tabs.css:42-50` — `height: 48px`, `padding: 0 16px`, `font-size: .875rem`, `font-weight: 500`, `letter-spacing: .00625rem` = label-large Material reimplementado a mano en vez de `var(--ui-type-label-lg)`. Ídem `:85-91` (`24px`) e indicador `3px/2px` (`:98,105-108`).
- `web/styles/data-table.css:14-16` — defaults `56px/52px/16px` "Material anatomy" (comentario `:5-6`).
- `web/styles/fab.css:29-64` — `56px/40px/96px/48px` de anatomía Material.
- `web/styles/menu.css:15-18,34-36` — `48px/12px/24px/112px/8px`.
- `web/styles/navigation-drawer.css:17-22,67` — `360px/56px/12px/24px` + `max-width: calc(100vw - 56px)`.
- `web/styles/list.css:12-17`, `navigation-bar.css:18-31`, `navigation-tab.css:17-29`, `segmented-button.css:15-25` — mismas tallas fijas de anatomía.

### 3.3 Shape del theme no derivado de la escala `--ui-radius-*`

- `theme.css:68` `--ui-dialog-radius: 28px`; `:83` `--ui-card-radius: 12px`; `:40` `--ui-toast-radius: 4px`; `:157-160` FAB shapes 16/12/28/20px; `:168` `--ui-select-radius: 4px`. Son decisiones de personalidad Material incrustadas como literales px en vez de componer `--ui-radius-*`.

### 3.4 Identidad Material en markup/classes

- `web/templates/layout.html:2` — `<html lang="en" class="theme-material">`: el nombre del theme está hardcodeado en el layout; todo el árbol depende de esa clase.
- `web/styles/select-menu.css:23,75,93-94` y `web/styles/app.css:63,174,180` — clase `m3-select-trigger` con prefijo `m3-` (contamina el contrato de clases; debería ser `ui-select-trigger`).

### 3.5 Motion Material hardcodeado

- `web/styles/dialog.css:19` — `transition: … 150ms …` con literal; `:37,39` `transition-duration: 500ms` literal (ignora `--ui-motion-short`).
- Ídem motion literales en `navigation-drawer.css:74,86,90` (`150ms`/`500ms`).

### 3.6 Dark mode duplicado con drift

- `theme.css:225-226,230-231` — indentación rota (declaraciones a columna 0), síntoma de merges manuales.
- `theme.css:253-299` (media query) NO redefine `--ui-switch-track-unselected`, pero el bloque `.theme-dark` (`:225`) sí → el switch track queda claro en auto-dark mientras el resto oscurece.
- El dark explícito y el auto redefinen el mismo set de tokens duplicando ~40 líneas; cualquier cambio requiere tocarlos en dos bloques.

### 3.7 Otros

- `web/styles/base.css:3` — `html { color-scheme: light dark }` vs `theme.css:7` `.theme-material { color-scheme: light }`: dos fuentes de verdad.
- `web/styles/app.css:38-48` — mapping `@theme inline` solo cubre 8 tokens (canvas, surface, fg, primary, primary-fg, border, font-sans, radius-md, shadow-elevation-1); el resto del sistema usa vars crudas. Contrato parcial con Tailwind.

---

## 4. Foundations faltantes

| Primitiva | Estado | Evidencia |
|---|---|---|
| Focus ring | ✅ Existe como tokens (`--ui-focus-*`) y `.ui-focus-ring` | `focus-ring.css:7-14`, `theme.css:197-198` |
| Elevation | ✅ Existe (`--ui-shadow-0..5`) y utilities `.ui-elevation-*` | `elevation.css:4-12` |
| Spacing scale | ❌ No existe ningún token `--ui-space-*` | — |
| Radius scale | ✅ Existe `--ui-radius-*`, pero componentes/shape del theme usan px propios | `theme.css:61-67` vs `:68,:83,:40` |
| Density | ❌ No existe; alturas fijas en componentes | `button.css:5`, `text-field.css:10` |
| Motion scale | ⚠️ Solo `short`+`standard`; dialog usa literales | `theme.css:199-200`, `dialog.css:19,37,39` |
| State layer tokens | ✅ Opacidades existen; ❌ colores del layer hardcodeados | `theme.css:192-196`, `button.css:17`, `chips.css:32` |
| Border (anchura/estilo) | ❌ No hay tokens (solo `--ui-color-border*`) | — |
| Z-index / top layer | ❌ Solo `--ui-tooltip-z` scoped; toast hardcodea `1000` | `tooltip.css:28`, `toast.css:22` |
| Breakpoints | ❌ Sin tokens; media queries hardcodeadas | `text-field.css:160` |
| Reduced motion | ⚠️ Por media query duplicada por archivo + central | `app.css:52-69` + 1 bloque por componente |
| Forced colors | ⚠️ Mismo patrón duplicado | `app.css:71-213` + por componente |
| Dark mode | ⚠️ Existe pero duplicado y con drift | `theme.css:203-251` vs `253-299` |
| Icons | ✅ `.ui-icon` 24px + contrato aria | `icon.css:5-10` |
| Tipografía | ⚠️ Escala existe pero shorthand bundleada y con huecos | `theme.css:48-58` |

---

## 5. Contrato de tokens actual

**No hay contrato formal de naming/cobertura.** Los tokens crecieron orgánicamente por componente. Evidencia:

- El propio roadmap lo reconoce: `docs/gelium-ui-system-roadmap.md:32` — *"Tokens públicos `--ui-*` | Parcial; no formalizado como contrato"*. La Phase 1 (`:94-119`) pide "contrato de naming y cobertura" y DoD "tokens `--ui-*` inventariados" — es trabajo pendiente.
- Naming de facto: familias por prefijo `--ui-color-*`, `--ui-type-*`, `--ui-radius-*`, `--ui-shadow-*`, `--ui-state-*`, `--ui-focus-*`, `--ui-motion-*`/`--ui-easing-*` y `--ui-<componente>-*`. No hay convención publicada (p.ej. `--ui-color-fg` vs `--ui-color-on-surface`; `--ui-color-surface-container` referenciado pero inexistente).
- La **documentación** es descentralizada: tablas de tokens en los `.md` de cada componente (`web/content/data-table.md:75`, `web/content/slider.md:18`, etc.) + comentarios en los `.css` + los tests.
- El **único "contrato" ejecutable** son los tests Go que asertan strings del CSS y valores literales del theme: `web/styles_contract_test.go` (concatenación de `app.css`), y p.ej. `web/styles_toast_test.go:28-48` aserta `--ui-toast-container: #322f35;` y que aparece **3 veces** (light + dark explícito + media). Consecuencia: mover/renombrar tokens a un core rompe tests; hay que migrar estos asertos a tests de contrato (existencia/cobertura) en vez de valores literales.
- 3 tokens definidos y muertos: `--ui-radius-xl` (`theme.css:66`), `--ui-state-dragged-opacity` (`theme.css:195`), `--ui-select-menu-item-icon` (`theme.css:187`).
- 5 tokens referenciados y no definidos: `--ui-color-surface-container` (`base.css:40`), `--ui-type-display-lg` (`base.css:37`), `--ui-type-title-md` (`demo-whatsapp.css:477`), `--ui-color-error` (`demo-whatsapp.css:403`), `--ui-font-mono` (`demo-whatsapp.css:520,528`).

---

## 6. Estructura de build CSS

- **Entrypoint**: `web/styles/app.css` (Tailwind 4). Importa `tailwindcss` → **theme** (`../../themes/theme-material/theme.css`, línea 2) → `base.css` → 29 archivos por componente → `demo-whatsapp.css` (`app.css:1-33`). Declara `@source` para `../templates` y `../content` (`:35-36`) y un bloque `@theme inline` (`:38-48`) que mapea 8 tokens a utilities Tailwind. Termina con `@keyframes ui-spin`, el bloque central `prefers-reduced-motion` (`:52-69`) y el bloque central `forced-colors` (`:71-213`).
- **Integración del theme**: por `@import` con ruta hardcodeada a Material (`app.css:2`). Para multi-theme hay que parametrizar esta línea (variable de entorno/build, p.ej. `THEME=basecoat`), porque hoy el único punto de selección es la ruta de import.
- **Compilación**: `npm run build` → `tailwindcss -i ./web/styles/app.css -o ./web/static/app.css --minify && node ./scripts/copy-htmx.mjs` (`package.json:6`). `copy-htmx.mjs` copia `htmx.min.js` desde `node_modules/htmx.org` a `web/static/` (`scripts/copy-htmx.mjs:4-6`).
- **Assets**: `web/static/app.css`, `app.js` y `htmx.min.js` se mantienen como outputs exclusivos del build y se embeben vía `go:embed` en `web/assets.go` (ver `README.md:24-27,103-105`).
- **Layering**: cada archivo envuelve sus reglas en `@layer components` (o `base`); el theme no usa `@layer` (es un set de custom properties). Los bloques `reduced-motion`/`forced-colors` se repiten por componente + central.
- **Tests de estructura**: `web/styles_contract_test.go:21-42` mantiene la lista canónica de imports en orden y concatena los fuentes para asertar contratos; hay que mantenerla sincronizada con `app.css`.

---

## 7. Propuesta de core agnóstico

### 7.1 Token families que el core debe definir (una sola vez, sin dirección visual)

1. **Color semántico**: `--ui-color-{canvas,surface,primary,secondary,error,warning,success,info,outline,outline-strong,scrim,focus-ring}` + variantes `*-fg` (on-). El core declara los *roles* con valores neutros de contraste por defecto; Material provee la palette.
2. **Tipografía composable**: descomponer los `--ui-type-*` shorthand en `--ui-font-family-*`, `--ui-font-size-*`, `--ui-font-weight-*`, `--ui-line-height-*`, `--ui-letter-spacing-*` por step (display/headline/title/body/label) y reconstruir `--ui-type-*` como alias. Cerrar los huecos (`display-lg`, `title-md`) hoy referenciados y ausentes.
3. **Spacing**: escala `--ui-space-0..N` (base 4px o 0.25rem) y usarla en todos los paddings/gaps de componentes.
4. **Radius**: mantener `--ui-radius-*`; normalizar los px del theme (`dialog 28px`, `card 12px`, `fab`, `select 4px`) como composición de la escala o token de componente con default.
5. **Density/size**: `--ui-density-{comfortable,compact,roomy}` y `--ui-size-{control,item,icon}` derivando las alturas hoy literales (`2.5rem`, `3.5rem`, `48px`, `56px`, `24px`).
6. **Elevation**: ya existe (`--ui-shadow-0..5`); pasa al core tal cual.
7. **Border**: `--ui-border-width-*`, `--ui-border-style-*` (reusar `--ui-color-border*`).
8. **Focus**: ya existe (`--ui-focus-*`); pasa al core.
9. **Motion**: `--ui-motion-{short,medium,long}` + `--ui-easing-{standard,emphasized,decelerate,accelerate}`; reemplazar literales de Dialog/NAV-drawer.
10. **State layer**: `--ui-state-*` opacities + **color del layer por `color-mix()` sobre el token `-fg`** (no `rgb()` fijos).
11. **Z-index/top-layer**: escala `--ui-z-{dialog,menu,popover,toast}`; reemplaza `z-index: 1000` de toast y `--ui-tooltip-z`.
12. **Breakpoints**: `--ui-breakpoint-{sm,md,lg}` y media queries por tokens.
13. **Reduced motion / forced colors**: estrategia única en el core (bloques centrales con las excepciones mínimas), eliminando la duplicación por archivo.
14. **Dark mode**: un solo mecanismo (`color-scheme` + un set de tokens por modo), sin duplicar el bloque.

### 7.2 Qué extraer del theme Material hacia el core

- Los tokens **cross-component** que hoy viven en Material pero son estructurales: alturas de controles/items, tamaño de iconos (24px), padding inline de items (12/16px), opacidades de estado, focus, motion, elevation, state layers.
- Los **nombres semánticos** ya existentes (`--ui-color-primary`, `--ui-color-fg`, `--ui-field-*`, `--ui-dialog-scrim`) quedan en el core como contratos; Material solo aporta *valores* de palette, personalidad de shape/tipografía/elevación y densidad elegida.
- Los **alias de compatibilidad** (`--ui-color-error` → `--ui-color-danger`, `--ui-color-surface-container` → token core `surface-container`) para eliminar referencias rotas.
- Quitar del markup/classes la identidad: `class="theme-material"` sale del layout (pasa a ser decisión del theme contract / config), y `m3-select-trigger` se renombra a `ui-select-trigger`.
- El `@theme inline` de `app.css:38-48` se convierte en mapping automático/declarado de los tokens core, no de Material.

### 7.3 Qué NO debe extraerse

- Los **valores** de la palette Material (hex), la personalidad de shape (28px dialogo, pill buttons), el typescale Material y la densidad por defecto: son dirección visual del theme.

### 7.4 Impacto en tests (planificarlo)

`web/styles_toast_test.go:28-48` y similares asertan valores literales y *cantidades* de declaraciones del theme; la extracción romperá estos tests. Migrar a tests de contrato: existencia de familia de tokens, cobertura light/dark, presencia de `var(--ui-*)` en componentes, y un test que verifique que ningún componente usa valores hardcodeados para color de estado/geometría de controles. Mantener `styles_contract_test.go` en sync con el nuevo `app.css`.

---

## 8. Archivos leídos (evidencia)

- `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`
- `themes/theme-material/theme.css` (300 líneas)
- `web/styles/app.css` + `base.css`, `button.css`, `text-field.css`, `dialog.css`, `toast.css`, `focus-ring.css`, `elevation.css`, `icon.css`, `divider.css`, `card.css`, `badge.css`, `checkbox.css`, `radio.css`, `switch.css`, `select.css`, `select-menu.css`, `slider.css`, `progress.css`, `icon-button.css`, `fab.css`, `list.css`, `chips.css`, `tabs.css`, `navigation-bar.css`, `navigation-tab.css`, `segmented-button.css`, `menu.css`, `navigation-drawer.css`, `data-table.css`, `tooltip.css`, `demo-whatsapp.css`
- `package.json`, `scripts/copy-htmx.mjs`, `web/templates/layout.html`
- Tests: `web/styles_contract_test.go`, `web/styles_button_test.go`, `web/styles_toast_test.go` (+ inventario `web/styles_*_test.go`)
