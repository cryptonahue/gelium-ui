# Theme Architecture Audit — Gelium UI

> Auditoría read-only de la arquitectura de themes (Phase 1/5 del roadmap). No modifica código, templates, CSS ni tests. Única salida: este handoff.
> Cómo se relaciona con el roadmap: el tema vive en `themes/theme-material/theme.css`; el objetivo es un THEME CONTRACT (Phase 5) que permita `theme-basecoat` (Phase 6) sin duplicar el sistema. Este documento inventaría lo que existe y fija qué congelar antes de extraer.

---

## 1. Resumen ejecutivo

- El sistema tiene **un único theme** (`themes/theme-material/theme.css`, 157 tokens `--ui-*`) que se integra **en build time** por `@import` desde `web/styles/app.css:2` y se empaqueta minificado en `web/static/app.css` (131 KB), que es el **único** asset CSS embebido (`web/assets.go:8`).
- El theming es **token-only por cascade CSS**: los componentes leen `var(--ui-*)` en runtime; no hay utilidades Tailwind involucradas (los templates no usan ni una clase utilitaria `bg-*`/`text-*`). Tailwind actúa solo como bundler/preflight/`@theme inline`.
- El **contrato token es estable en la forma** (`--ui-<familia>-<token>`, `--ui-<componente>-<rol>`) pero **incompleto e inconsistente en cobertura**: hay 5+ tokens referenciados sin definir (`--ui-color-surface-container`, `--ui-type-display-lg`, `--ui-type-title-md`, `--ui-color-error`, `--ui-font-mono`), 3 tokens definidos sin uso, y **8 familias de componentes declaran sus tokens locales y el theme NO puede sobreescribirlos** (list, menu, data-table, nav-bar, nav-tab, navigation-drawer, segmented-button, tooltip).
- El **dark mode está duplicado a mano** en dos rutas (clase `.dark`/`[data-theme]` y `@media prefers-color-scheme`) con **drift real** entre ambas (`--ui-switch-track-unselected` solo existe en la ruta por clase). No se usa `light-dark()`.
- Los **acoplamientos Material residuales** en componentes son pocos pero duros: state-layer con colores hardcodeados (`button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`), geometría en px (`fab.css:29-42,54`, `dialog.css:3-6,26`, `text-field.css:10`), y estilos de demo dentro de `web/styles/`.
- Los **tests ya verifican el contrato CSS** en 3 capas (source de componente, cobertura de tokens del theme, asset compilado), pero **13 tests hardcodean la ruta `theme-material/theme.css` y 3 hardcodean valores hex Material** (`styles_fab_test.go:85-87`, `styles_dialog_test.go:15-29`, `styles_toast_test.go:37-42`). No son theme-agnósticos todavía.
- **Propuesta**: mecanismo multi-theme por **bundle de todos los themes + selección en runtime por clase/attr en `<html>`**, cada theme autocontenido (light+dark), componentes intactos. Requiere antes: formalizar naming, fijar tokens faltantes, tokenizar state-layers, parametrizar tests, y hacer data-driven la clase del theme (`layout.html:2`).

---

## 2. Mecanismo actual de themes — flujo completo theme.css → web/static/app.css

**Fuentes primarias:**
- `package.json:6` — `"build": "tailwindcss -i ./web/styles/app.css -o ./web/static/app.css --minify && node ./scripts/copy-htmx.mjs"`
- `web/styles/app.css` es la **entrada única** (`-i`) y orquesta todo el bundle.
- `web/static/app.css` es el **output**; se embebe con `//go:embed ... static/*` (`web/assets.go:8`).

**Flujo trazado:**

```
themes/theme-material/theme.css          (source, tokens --ui-*; NO embebido directo)
        │  @import "../../themes/theme-material/theme.css"   (web/styles/app.css:2)
        ▼
web/styles/app.css                       (entrada Tailwind CLI)
  ├─ @import "tailwindcss"                (app.css:1 — preflight + framework)
  ├─ @import theme-material/theme.css     (app.css:2)
  ├─ @import ./base.css + ./button.css + … (app.css:3-33 — 32 archivos de componentes)
  ├─ @source "../templates" ../content    (app.css:35-36 — escaneo de utilidades)
  └─ @theme inline { --color-canvas: var(--ui-color-canvas); … }  (app.css:38-48)
        │  tailwindcss CLI --minify  (package.json:6)  →  lightningcss resuelve imports
        ▼
web/static/app.css                       (131.075 bytes, una línea minificada, SIN @import residuales)
        │  //go:embed static/*  (web/assets.go:8)
        ▼
binario Go  →  servido en /static/app.css  (README.md:59)
```

**Verificado en el output compilado (`web/static/app.css`):**
- El theme se **inlinea**: `.theme-material` aparece 5×, `--ui-color-primary:` 3× (luz + clase dark + media query), el bloque `@media (prefers-color-scheme: dark)` está presente (1×).
- **No** queda `@import` en el output; `light-dark()` **no** se usa (0×). `prefers-reduced-motion` y `forced-colors` (bloques de `app.css:52-69` y `app.css:71-212`) pasan al bundle.
- Los 9 mapeos de `@theme inline` (app.css:38-48) resuelven a `var(--ui-*)` → runtime. Pero **los templates no usan ninguna utilidad Tailwind tokenizada** (0 coincidencias de `bg-primary|text-fg|…` en `web/templates/*.html`); el `@source` es en la práctica inerte para el theming.

**Activación en el documento:**
- `web/templates/layout.html:2` — `<html lang="en" class="theme-material">`. **La clase del theme está hardcodeada** en el template; no es data-driven.

**Conclusión:** el theme es **100% cascade-time CSS**. La selección de theme puede hacerse en runtime (clase/attr) sin recompilar, siempre que el bundle contenga los themes.

---

## 3. Inventario de tokens `--ui-*` en theme-material

**Fuente:** `themes/theme-material/theme.css` (300 líneas).

### 3.1 Convención de naming
- **Global:** `--ui-<familia>-<token>` → ej. `--ui-color-primary`, `--ui-radius-md`, `--ui-shadow-1`, `--ui-state-hover-opacity`, `--ui-motion-short`, `--ui-easing-standard`.
- **Por componente:** `--ui-<componente>-<rol>` → ej. `--ui-field-container`, `--ui-dialog-scrim`, `--ui-toast-action`, `--ui-card-radius`, `--ui-checkbox-size`, `--ui-fab-primary-container`, `--ui-select-menu-item-height`.
- Familias transversales (colors/shape/type/shadow/state/motion) = foundation del core; familias por componente = cobertura del tema. Esta separación es la que el contrato debe formalizar (Phase 5).

### 3.2 Familias (157 tokens únicos definidos)

| Familia | Tokens | Líneas (theme.css) |
|---|---|---|
| color (semántico) | 13 | 10–22 |
| field (text field) | 5 | 25–29 |
| dialog | 5 | 32–35 |
| toast | 8 | 38–45 |
| type (typescale) | 10 | 48–58 |
| radius (shape) | 7 | 61–68 |
| shadow (elevation) | 6 | 71–76 |
| state (opacidades) | 5 | 192–196 |
| focus | 2 | 197–198 |
| motion + easing + font | 3 | 48, 199–200 |
| divider / card / badge | 2+5+4 | 79–93 |
| checkbox / radio / switch | 9+7+15 | 96–130 |
| slider / progress | 10+4 | 133–148 |
| fab | 14 | 151–164 |
| select | 23 | 167–189 |

### 3.3 Light vs dark
- **Base (luz):** bloque `.theme-material { color-scheme: light; … }` (theme.css:6–201) define los 157 tokens.
- **Ruta dark A (por clase/attr):** `.theme-material.theme-dark, .theme-material.dark, .theme-material[data-theme="dark"]` (theme.css:203–251) — ~40 overrides.
- **Ruta dark B (por sistema):** `@media (prefers-color-scheme: dark)` + `.theme-material:not([data-theme="light"])` (theme.css:253–299) — ~40 overrides **duplicados a mano**.
- **NO usa `light-dark()`** (0 en output compilado). No hay pares de tokens por scheme; todo es duplicación textual.
- **Drift entre rutas dark:** `--ui-switch-track-unselected` se redefine en la ruta A (theme.css:225) pero **NO** en la media query (253–299) → con `prefers-color-scheme: dark` el switch no seleccionado conserva el valor de luz `#e6e0e9`. Además las líneas 225–226 tienen indentación rota (edit manual). Es un ejemplo concreto de por qué duplicar dark a mano no escala.
- 41 tokens aparecen ≥3× (light + ruta A + ruta B); 115 solo en luz.

### 3.4 Cobertura por componente (quién tiene tokens en el theme)

| Componente | ¿Tokens en theme.css? | Fuente de tokens |
|---|---|---|
| Button, Text field, Dialog, Toast, Card, Badge, Checkbox, Radio, Switch, Slider, Progress, FAB, Select, Select menu, Divider, Icon, Elevation, Focus ring | **Sí** | `theme.css` |
| List, Menu, Data table, Navigation bar, Navigation tab, Navigation drawer, Segmented button, Tooltip | **No** | tokens locales scoped en el CSS del componente (defaults en el root) |

Los 8 componentes de la segunda fila declaran sus propios defaults scoped (p. ej. `.ui-data-table { --ui-data-table-*: … }`, `data-table.css:12–25`; lista completa en §6). El theme **podría** sobreescribirlos globalmente declarándolos en `.theme-material`, pero **hoy no lo hace** → para estos 8, un nuevo theme no tiene palanca más allá de los `--ui-color-*` base.

### 3.5 Gaps de definición (referenciados sin definir en el theme)

| Token | Dónde se referencia | Impacto |
|---|---|---|
| `--ui-color-surface-container` | `base.css:40` + 14× en `demo-whatsapp.css` (63, 102, 145, …) | **Sin fallback** → `var()` cae a inherited/transparent. Bug visual real en luz y dark. |
| `--ui-type-display-lg` | `base.css:37` (`.landing-hero h1`) | Sin fallback → font inválida, hereda. |
| `--ui-type-title-md` | `demo-whatsapp.css:477` | Sin fallback. |
| `--ui-color-error` | `demo-whatsapp.css:403` (`var(--ui-color-error, #b3261e)`) | **Drift de naming**: el theme define `--ui-color-danger`, no `error`. |
| `--ui-font-mono` | `demo-whatsapp.css:520,528` | Con fallback `, monospace` (degrada bien) pero el theme no lo define. |
| `--ui-slider-fill` | `slider.css:32` | No es token de theme: es knob por instancia (porcentaje del slider). OK. |

### 3.6 Tokens definidos sin uso real en componentes (4)

`--ui-color-danger-fg` (usado indirectamente por `--ui-badge-fg: var(--ui-color-danger-fg)` en theme.css:93 → OK), `--ui-radius-xl` (sin uso), `--ui-select-menu-item-icon` (sin uso), `--ui-state-dragged-opacity` (sin uso; sin soporte de drag en componentes).

---

## 4. Estilos de componentes: tokens vs valores fijos

**Métricas:** 1.181 usos de `var(--ui-*)` en `web/styles/*.css`; **0 hex Material hardcodeado en componentes reales** (solo `demo-whatsapp.css`, que es demo). La higiene token es alta, pero hay acoplamientos de otra naturaleza:

### 4.1 State-layer con colores base hardcodeados (lo más importante a tokenizar)

Los overlays de estado (hover/pressed) se pintan con `inset box-shadow` de un color **RGB fijo**, no derivado de tokens:

- `web/styles/button.css:17-18` — `box-shadow: inset 0 0 0 999px rgb(255 255 255 / var(--ui-state-hover-opacity))` (y `pressed`). Blanco hardcodeado para variantes rellenas.
- `web/styles/icon-button.css:22,25` — `rgb(255 255 255 / …)` idem.
- `web/styles/chips.css:32,35,107,141-142` — `rgb(25 28 28 / …)` y `rgb(29 25 43 / …)` (on-surface / on-secondary-container MD3).

**Técnica theme-aware ya existente (a emular):** `color-mix(in srgb, var(--ui-color-primary) calc(var(--ui-state-*) * 100%), transparent)` en `button.css:27-32`, y `color-mix(... currentColor …)` en `fab.css:80-88`. Solo 4 archivos usan `color-mix` (text-field ×4, fab ×3, button ×2, select-menu ×1). Un theme que cambie primary no afecta los overlays blancos de button/icon-button/chips.

### 4.2 Geometría Material en px/rem fijos

- `button.css:5` — `min-height: 2.5rem` (48px target MD3).
- `fab.css:29-30` — `56px`; `fab.css:35-36` — `40px`; `fab.css:41-42` — `96px`; `fab.css:54` — `48px` (4 tamaños Material). (Los tests los exigen: `styles_fab_test.go:28-38`.)
- `dialog.css:3-6` — `280px/140px/560px/48px`; `dialog.css:26` — padding `16px 24px 24px`.
- `text-field.css:10` — `height: 3.5rem`; `switch.css:40` — aritmética `calc(var(--ui-switch-width) - … - 12px)`; `switch.css:66` — `inset-inline-start: 6px`.
- `data-table.css:14-16` — `56px/52px/16px` pero **ya tokenizado** (`--ui-data-table-header-height` etc., scoped en `data-table.css:12-25`). Patrón a replicar.

Esto es **anatomía** (definición del vocabulario visual) vs **dirección visual** (paleta/shape/motion). La pregunta de contrato es: ¿la geometría es parte del theme o del core? Recomendación: **core** (esquema fijo), tokens solo para los grados de libertad que un theme quiere cambiar (radius, surface, color).

### 4.3 Typography / shape / elevation: bien tokenizados
- `font: var(--ui-type-*)` en todos los componentes (p. ej. `button.css:12`, `card.css:20`, `dialog.css:22-23`).
- `border-radius: var(--ui-radius-*)` y `var(--ui-*-radius)` en casi todos (excepción notable: `text-field.css:12` usa `var(--ui-radius-xs)`, OK; `select.css` usa `var(--ui-select-radius)`).
- Elevation 100% vía `--ui-shadow-*` (`fab.css:32`, `toast.css:39`, `card.css:22`).
- Los únicos valores fijos de layout/demo (`component-preview`, `card-demo-grid`, `data-table-demo-*`) son spacing de documentación, no parte del contrato.

### 4.4 El puente `@theme inline` (app.css:38-48)
Expone solo 9 tokens a Tailwind (`canvas, surface, fg, primary, primary-fg, border, font-sans, radius-md, shadow-elevation-1`). Con 0 utilidades usadas en templates, hoy es **dormante**: el theming real vive en las reglas `.ui-*` que leen `--ui-*` directo. No bloquea el multi-theme pero conviene decidir si se expande o se elimina.

---

## 5. Estado/variantes por componente — qué debe cubrir un theme

**Mecanismos de modelado (todos en el CSS de componentes, no en el theme):**

1. **Clases de variante:** `.ui-button-primary/secondary/outline/text` (`button.css:23-26`), `.ui-fab-primary/surface/secondary/lowered` (`fab.css:68-75`), `.ui-text-field-outlined/filled` (`text-field.css:72,93`), `.ui-select-filled/outlined` (`select.css:48,55`), `.ui-card-elevated/filled/outlined` (`card.css:22-24`), `.ui-data-table-cell--label/--checkbox/--sortable` (`data-table.css:65,111,136`).
2. **Pseudo-estados nativos:** `:hover:not(:disabled):not([aria-disabled="true"])` (`button.css:17`), `:active`, `:focus-visible` (`button.css:19-22`, `fab.css:85-89`), `:checked` (`switch.css:34-41`, `checkbox/radio`), `:disabled` (`button.css:33`, `switch.css:46-49`), `:focus-within` (`text-field.css:32,98`), `:placeholder-shown` (`text-field.css:33-34`), `:has(input:checked)` (`data-table.css:104-107`), `:has(select:disabled)` (`select.css:88`).
3. **Estados ARIA/atributos:** `[aria-disabled="true"]` (`button.css:17,33`, `fab.css:79-90`), `[aria-invalid="true"]` (`select.css:45-47`), `aria-selected` (menu/select-menu), `aria-current="page"` (paginación data table + clase `--current`, `data-table.css:204-207`).
4. **Estado disabled genérico:** `opacity: var(--ui-state-disabled-opacity)` en todos los componentes (`button.css:34`, `fab.css:91`, `text-field.css:137`, `data-table.css:210`, `switch.css:47-49`).
5. **State-layer MD3:** inset overlay con `--ui-state-hover/focus/pressed/dragged/disabled-opacity` (`theme.css:192-196`) sobre `currentColor`/`color-mix`/RGB fijo (§4.1).
6. **Reduced motion / forced colors:** centralizado en `app.css:52-69` (reduced) y `app.css:71-212` (forced-colors), con refuerzos por componente en algunos archivos (`data-table.css:334-359`).

**Lectura para el contrato:** el theme **no** modela estados; modela tokens. Un theme nuevo debe cubrir, por componente, los tokens de color/surface/shape/type/motion que esos estados consumen. La superficie de estados es estable y queda documentada en los tests (ver §6).

---

## 6. Pruebas de estilos — qué cubren y cómo verificarían un theme nuevo

**Arquitectura de tests (3 capas), todas en `web/styles_*_test.go`:**

1. **Contrato del source del componente** — `sourceComponentCSS(t, "fab.css")` lee del embed `//go:embed styles/*.css` (`styles_contract_test.go:12-13`) y asevera selectores + referencias token exactas del archivo fuente:
   - `styles_fab_test.go:20-57` (selectores `.ui-fab-*`, `var(--ui-fab-*)`, estado hover/active con exclusión de `aria-disabled`).
   - `styles_data_table_test.go:9-50` (anatomía + `:has(input:checked)`).
   - `styles_button_test.go:9-33` (orden disabled-después-de-interactivos; estados).
2. **Cobertura de tokens del theme** — `repositoryFile(t, "themes", "theme-material", "theme.css")` lee el theme **por ruta** y asevera que defina los tokens del componente (13 tests). **⚠️ 3 tests hardcodean valores hex Material:**
   - `styles_fab_test.go:85-87` — cuenta `--ui-fab-primary-container: #4f378b` + `#36343b` + `#4a4458` (≥3).
   - `styles_dialog_test.go:15-29` — asevera `#ece6f0/#1d1b20/#49454f` luz y `#2b2930/#e6e0e9/#cac4d0` dark; `strings.Count(… "#2b2930;") != 2` falla si cambia el mecanismo dark.
   - `styles_toast_test.go:29-42` — asevera `#322f35/#f3edf7/#d0bcff/#ece6f0/#1d1b20`.
   - **Un theme con paleta distinta (Basecoat) rompería estos tests. El contrato debe volverlos theme-agnósticos (assert de presencia de token, no de valor).**
3. **Asset compilado embebido** — `Assets.ReadFile("static/app.css")` asevera que el bundle contenga selectores/tokens (`styles_fab_test.go:110-131`). Depende de que `npm run build` se haya corrido (el `web/static/app.css` embebido es un artefacto de build).
4. **Entrada app.css** — `sourceAppCSS(t)` concatena source en orden fijo (`styles_contract_test.go:15-52`) para assert de reduced-motion/forced-colors (`styles_fab_test.go:93-108`, `styles_button_test.go:50-65`). **⚠️ La lista de `sourceAppCSS` está desactualizada** (`styles_contract_test.go:21-42`): no incluye icon-button, fab, chips, tabs, navigation-*, segmented-button, menu, navigation-drawer, data-table, tooltip ni demo-whatsapp. Los componentes nuevos se cubren vía `sourceComponentCSS` propio, pero el helper compartido mintió la cobertura.

**Cómo verificaría un theme nuevo:** hoy NO puede verificarse sin editar tests, porque (a) la ruta `themes/theme-material/theme.css` está hardcodeada en 13 tests y (b) los valores hex Material están aseverados. La parametrización necesaria es la del theme contract: un test que (1) lea `themes/<theme>/theme.css`, (2) verifique que defina **todo** token referenciado por los componentes (cerrar los gaps de §3.5), (3) verifique que cada componente tenga su familia de tokens, (4) verifique dark en ambas rutas sin asumir valores, (5) ejecute la matriz `theme × componente × variante × estado` sobre el asset compilado.

---

## 7. Propuesta de mecanismo multi-theme

**Restricciones de partida:**
- `web/static/app.css` es el **único** asset CSS embebido (`web/assets.go:8`); `go:embed` no admite `..` (README.md:103-105).
- El theme se inyecta hoy por `@import` en build (`app.css:2`) y la clase se hardcodea en `layout.html:2`.
- El contrato es **token-only**; los componentes leen `var(--ui-*)` en cascade → **la selección de theme no necesita rebuild**.

**Mecanismo propuesto — bundle de todos los themes + selección runtime por clase:**

```text
themes/
├── theme-material/theme.css   →  .theme-material { --ui-* … }  (luz + dark, autocontenido)
└── theme-basecoat/theme.css   →  .theme-basecoat  { --ui-* … }  (luz + dark, autocontenido)
        │
        │  web/styles/app.css importa CADA theme (listado explícito; CSS no globbea)
        ▼
web/static/app.css             (UN solo asset, contiene todos los themes)
        │  embed
        ▼
runtime: <html class="theme-material">  ←  data-driven desde template/server
        (o <html data-theme="basecoat">; dark interno de cada theme)
```

1. **Un `theme.css` por theme** en `themes/<theme>/`, cada uno con su selector raíz propio `.theme-<name>` (renombrar `.theme-material` → `.theme-material` se mantiene como selector canónico, pero el contrato exige uno por theme). Autocontenido: luz + dark (sin depender de la duplicación por media query manual; ver nota abajo).
2. **app.css importa todos**: `@import "../../themes/theme-material/theme.css"; @import "../../themes/theme-basecoat/theme.css"; …`. Sin glob (CSS no soporta globs); la lista se mantiene explícita en `app.css` (ya es la convención de import de componentes).
3. **Selección en runtime** por clase/attr en `<html>` (`layout.html:2` pasaría a template-driven: `{{.ThemeClass}}` o `data-theme`). Un binary, todos los themes, sin rebuild por theme.
4. **Dark:** el contrato debe decidir entre (a) pares `light-dark()` si hay soporte objetivo, o (b) una **única** rutina por theme (`[data-theme="dark"]`/`.dark` + media query generada o combinada) eliminando la duplicación textual actual. **La duplicación a mano (theme.css:203-299) es el primer candidato a borrar.**
5. **NO**: archivos por componente por theme (`themes/<theme>/<componente>.css`) — duplica estructura y build surface sin necesidad, ya que el contrato es token-only y los componentes son únicos.
6. **Variable de build (`THEME=basecoat npm run build`)**: válida como modo "solo un theme por binario" (selección en `app.css` condicionada por env en el script), pero **insuficiente** para el objetivo del roadmap (misma doc, misma suite, switch de dirección visual). El bundle-all es el mecanismo base; la variable de build es una optimización opcional de footprint.

**Lo que NO cambia:** componentes (`web/styles/*.css`, `web/templates/*.html`), contratos HTMX, `@theme inline` (es theme-agnóstico), tests (una vez parametrizados).

---

## 8. Riesgos de extracción del core de theme-material

**Qué se rompe si se extrae el core sin preparación:**

1. **8 familias de tokens viven en los componentes, no en el theme** (`--ui-list-*`, `--ui-menu-*`, `--ui-data-table-*`, `--ui-nav-bar-*`, `--ui-nav-tab-*`, `--ui-navigation-drawer-*`, `--ui-segmented-button-*`, `--ui-tooltip-*`, scoped en `data-table.css:12-25`, `list.css`, `menu.css`, `tooltip.css`, etc.). Al extraer el core hay que **decidir la propiedad** de cada familia; si se extrae solo lo que hoy está en theme-material, Basecoat no puede tocar 8 componentes.
2. **Tokens con gaps de definición** (§3.5) se convertirían en "contrato formal" con agujeros: `--ui-color-surface-container` (base.css:40), `--ui-type-display-lg`/`title-md`, `--ui-color-error` vs `danger`, `--ui-font-mono`. **Congelar el naming y cerrar gaps ANTES** de escribir el contrato.
3. **State-layers hardcodeados** (`button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`): un theme que redefina `primary` no repinta los overlays → inconsistencia visual que la matriz cross-theme (Phase 7) detectaría tarde. Tokenizar con `color-mix` (patrón ya usado en `button.css:27-32`, `fab.css:80-88`) antes de la extracción.
4. **Dark duplicado a mano** (theme.css:203-299) con drift demostrado (`--ui-switch-track-unselected` en theme.css:225, ausente en 253-299): cualquier extracción que replique el patrón propaga el bug a cada theme nuevo.
5. **Tests Material-assertivos**: `repositoryFile(t, "themes","theme-material","theme.css")` en 13 tests + valores hex en `styles_fab_test.go:85-87`, `styles_dialog_test.go:15-29`, `styles_toast_test.go:37-42` y aserciones de geometría px (`styles_fab_test.go:28-38`). **Un theme nuevo no pasa la suite sin refactor de tests.**
6. **`layout.html:2` hardcodeado** y `sourceAppCSS` desactualizado (`styles_contract_test.go:21-42`): la extracción que toque app.css debe actualizar la lista del helper y hacer data-driven la clase del theme.
7. **README.md:89-90** — estructura desactualizada (lista solo `web/styles/app.css`, no los split files); menor, pero la extracción debe actualizarlo.

**Orden de congelado sugerido (antes de Phase 1/5):**
1. Inventario cerrado y gaps cerrados (§3.5, §3.6) + naming `danger`/`error` unificado.
2. Tokenizar state-layers (§4.1) y decidir geometría core-vs-theme (§4.2).
3. Unificar dark (eliminar duplicación, una rutina por theme).
4. Parametrizar tests (ruta de theme + valores) → suite theme-agnóstica.
5. Data-driven la clase del theme en `layout.html`.
6. Recién entonces: extraer el core agnóstico y escribir el contrato (Phase 5 gate).

---

## Anexo A — Inventario de archivos relevantes

- `themes/theme-material/theme.css` — tokens `--ui-*` (único theme).
- `web/styles/app.css` — entrada del build (imports, `@theme inline`, reduced-motion, forced-colors).
- `web/styles/*.css` — 32 archivos: base + 31 componentes/demos.
- `web/static/app.css` — output compilado (131 KB), embebido.
- `web/assets.go:8` — `//go:embed templates/*.html content/*.md static/*`.
- `web/templates/layout.html:2` — `<html class="theme-material">` (hardcodeado).
- `package.json:6` — script `build` (tailwindcss + copy-htmx).
- `scripts/copy-htmx.mjs` — copia `htmx.min.js` → `web/static/` (no toca CSS).
- `web/styles_contract_test.go` — helper `sourceAppCSS` (lista desactualizada) + `repositoryFile`.
- `web/styles_*_test.go` — 3 capas de test CSS (source componente / tokens theme / asset compilado).
- `docs/gelium-ui-system-roadmap.md` — roadmap (Phases 0-8; theme contract en Phase 5).
- `README.md` — setup/build, nota de embed, estructura (parcialmente desactualizada).
