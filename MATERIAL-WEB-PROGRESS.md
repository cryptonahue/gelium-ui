# Auditoría Material Web → Loom UI

> Snapshot técnico y funcional para decidir qué convertir a Loom UI. Esta auditoría no implementa componentes.

## 1. Snapshot de fuentes

Auditoría ejecutada el `2026-08-08T02:27:27-03:00`.

| Fuente | Ruta | Remote / procedencia | Versión | Commit / fecha | Tracked | Estado observado |
|---|---|---|---:|---|---:|---|
| Material Web upstream | `D:\repos\material-web-upstream` | `https://github.com/material-components/material-web.git` | `@material/web 2.5.0` | `23b638c1d2338b9fcaccef81cf54298cce9c27ed` / 2026-08-04 | 1.493 | `HEAD == origin/main`; limpio después de `git fetch --prune origin` |
| Fork local | `D:\repos\material-web-tailwind` | origin `cryptonahue/material-web-tailwind`; upstream `material-components/material-web` | `@material/web 2.3.0` | `70a1d8e59a127df886819197f24d76bf96524196` / 2025-08-13 | 1.029 | Tiene archivos/árboles untracked; el remote-tracking `upstream/main` local está obsoleto |
| Port CSS local | `D:\repos\material-web-tailwind\material-tailwind` | Árbol untracked dentro del fork; sin historial propio | `@material/tailwind-css 0.1.0` | `UNKNOWN` | 0 | Prototipo local, no repo independiente |

Evidencia principal: `material-web-upstream/package.json`, `material-web-tailwind/package.json`, `material-web-tailwind/material-tailwind/package.json` y comandos Git read-only.

### Método de comparación

Se compararon los índices producidos por `git ls-files -s`, usando ruta relativa y blob Git:

- 1.018 rutas comunes: 733 idénticas y 285 diferentes.
- 475 rutas sólo upstream.
- 11 rutas sólo fork.
- Upstream agrega principalmente `tokens/` (291 archivos), `labs/` (160), Sass (15), scripts/types y `tsconfig.base.json`.
- El port Tailwind no aparece en los 1.029 archivos del fork: todo `material-tailwind/` está untracked.

No se contaron `.git`, `node_modules`, `dist` ni otros artefactos generados como fuente. No se ejecutaron installs ni builds que pudieran modificar las fuentes.

## 2. Inspiraciones y atribución

- **Material Web** (`material-components/material-web`, Apache-2.0): fuente de inventario, tokens, variantes, comportamiento, accesibilidad, demos y tests. Loom reescribirá la implementación para HTML semántico, Tailwind y HTMX.
- **Basecoat UI**: referencia de componentes open-code, Tailwind, JavaScript mínimo y portabilidad entre backends. Loom no es un fork oficial.
- **shadcn/ui**: referencia para registry, componentes copiables y ownership del código por el proyecto consumidor.
- **shadcn-templ**: referencia para distribución/presets, sin atar Loom a Go ni a `templ`.
- **HTMX**: referencia para interacción server-driven, history, validación, fragmentos, `HX-Trigger`, SSE y WebSockets.

Antes de copiar texto, CSS o valores del port local debe resolverse procedencia/licencia: upstream declara Apache-2.0, mientras el árbol sin historial `material-tailwind/package.json` declara MIT.

## 3. Inventario de foundations y tokens

### Sets observados

| Set | SCSS | Component | Reference | System | Estado | Evidencia |
|---|---:|---:|---:|---:|---|---|
| Upstream `tokens/` | 60 | 49 | 2 | 6 | API estable consumida por Material Web 2.5.0 | `material-web-upstream/tokens/_index.scss` |
| Upstream `tokens/v0_192/` | 95 | 84 | 2 | 6 | Shim deprecado | `material-web-upstream/tokens/v0_192/_index.scss` |
| Upstream `tokens/versions/v0_192/` | 95 | 84 | 2 | 6 | Snapshot versionado para trazabilidad; **no** API estable upstream | `material-web-upstream/tokens/versions/v0_192/_index.scss` |
| Upstream `tokens/versions/latest/sass/` | 195 | 177 | 2 | 13 | Generado, Material 3 34.0.21; experimental | `material-web-upstream/tokens/versions/latest/sass/_index.scss` |
| Fork `tokens/` | 60 | 49 | 2 | 6 | API estable antigua | `material-web-tailwind/tokens/_index.scss` |
| Fork `tokens/v0_192/` | 95 | 84 | 2 | 6 | Contenido igual a upstream `versions/v0_192` | `material-web-tailwind/tokens/v0_192/_index.scss` |

`tokens/versions/README.md` advierte explícitamente que todo el árbol versionado puede introducir breaking changes en releases minor/patch y recomienda consumir `@material/web/tokens`. Loom usará la API pública root de este commit como fuente; `versions/v0_192` queda sólo como snapshot fijado por Loom para trazabilidad y `latest` no se seguirá automáticamente.

### Familias foundation

| Familia | Qué existe | Decisión Loom | Evidencia |
|---|---|---|---|
| Color | Palette/reference, roles system, light/dark, surfaces y contraste; `latest` agrega fixed roles y medium/high contrast | Exponer contrato semántico `--ui-color-*`; el theme Material mapea roles Material; mantener pares `on-*` atómicos | `material-web-upstream/tokens/_md-ref-palette.scss`; `tokens/_md-sys-color.scss`; `tokens/versions/latest/sass/_md-sys-color*.scss`; `docs/theming/color.md` |
| Typography | Typeface brand/plain y 15 roles: display/headline/title/body/label × small/medium/large | Mapear family, size, line-height, weight y tracking; diferir ejes de variable fonts | `tokens/_md-ref-typeface.scss`; `tokens/_md-sys-typescale.scss`; `typography/_typescale.scss`; `docs/theming/typography.md` |
| Shape | Escala none→full; `latest` agrega `extra-extra-large` y variantes increased | `--ui-radius-{none,xs,sm,md,lg,xl,full}`; aliases direccionales quedan privados del theme Material | `tokens/_md-sys-shape.scss`; `tokens/versions/latest/sass/_md-sys-shape.scss`; `docs/theming/shape.md` |
| Elevation | Niveles y sombras compuestas; además existe un Custom Element Lit | Conservar seis sombras CSS 0–5; no portar `<md-elevation>` ni DOM auxiliar | `elevation/internal/_elevation.scss`; `elevation/elevation.ts`; `elevation/elevation_test.ts` |
| Motion | Duraciones/easing; `latest` agrega springs, algunos representados como `null` en Sass | Exponer sólo easing/duraciones usados; diferir springs hasta tener representación válida | `tokens/_md-sys-motion.scss`; `tokens/versions/latest/sass/_md-sys-motion.scss` |
| State layers | Hover/focus/pressed/dragged/disabled; `latest` cambia focus/pressed 0.12→0.10 y agrega disabled 0.38 | Variables `--ui-state-*`, fuera de `@theme`; CSS con `prefers-reduced-motion` y `forced-colors` | `tokens/versions/v0_192/_md-sys-state.scss`; `tokens/versions/latest/sass/_md-sys-state.scss` |
| Focus | Focus ring Lit + Sass con tests; `latest` define 3px y offset 2px | Preferir `:focus-visible` y outline nativo; nunca reset global sin fallback | `focus/internal/focus-ring.ts`; `focus/internal/focus-ring_test.ts`; `tokens/versions/latest/sass/_md-sys-state-focus-indicator.scss` |
| Ripple | State machine Lit/Web Animations con touch delay, cancelación y origen del pointer | `defer`; state layer CSS en core y ripple posicional como JS opcional | `ripple/internal/ripple.ts`; `ripple/internal/ripple_test.ts` |
| Common/internal | Imports de conveniencia, controllers y event dispatch hooks de Web Components | No portar como foundations; rescatar sólo comportamiento verificable | `common.ts`; `internal/README.md`; `internal/events/dispatch-hooks.ts` |

### Cambios relevantes de `latest`

- 84 → 177 módulos component-token y 6 → 13 system-token.
- Nuevos tokens para app bars, button groups/tamaños, banners, FAB menu, nav rail, progress expresivo, split button, toolbar y loading indicator.
- Nuevos contextos light/dark de contraste medio/alto.
- Focus indicator, typescale emphasized, radii nuevos y ejes de variable font.
- Un token generado **no prueba** que exista un componente implementado, documentado o testeado.

Evidencia: `material-web-upstream/tokens/versions/latest/sass/_index.scss` y archivos `_md-comp-*.scss` de ese directorio.

### Contrato recomendado

```css
.theme-material {
  --ui-color-canvas: var(--material-surface);
  --ui-color-surface: var(--material-surface-container);
  --ui-color-fg: var(--material-on-surface);
  --ui-color-fg-muted: var(--material-on-surface-variant);
  --ui-color-primary: var(--material-primary);
  --ui-color-primary-fg: var(--material-on-primary);
  --ui-color-danger: var(--material-error);
  --ui-color-danger-fg: var(--material-on-error);
  --ui-color-border: var(--material-outline-variant);
  --ui-color-focus-ring: var(--material-secondary);
}

@theme inline {
  --color-canvas: var(--ui-color-canvas);
  --color-surface: var(--ui-color-surface);
  --color-fg: var(--ui-color-fg);
  --color-primary: var(--ui-color-primary);
  --color-primary-fg: var(--ui-color-primary-fg);
  --color-danger: var(--ui-color-danger);
  --color-border: var(--ui-color-border);
  --color-focus-ring: var(--ui-color-focus-ring);
  --radius-md: var(--ui-radius-md);
  --shadow-elevation-1: var(--ui-shadow-1);
}
```

En `@theme` deben entrar colores, fonts, tipos estáticos, radii, sombras y easing. State opacity, focus offsets, scrim, component tokens, springs, ripple y ejes variables deben quedar como variables semánticas/privadas fuera de `@theme`.

## 4. Inventario completo de componentes

Leyenda de fuentes:

- `Sí`: implementación/directorio observado.
- `CSS`: CSS local observado, sin prueba reproducible de build ni equivalencia funcional.
- `Parcial`: material relacionado, pero no componente equivalente completo.
- `No`: no se encontró implementación correspondiente.
- Todos comienzan en `NOT_STARTED` para Loom.

| Componente | Categoría | Upstream | Fork | Material Tailwind | Variantes | Estados / accesibilidad | JS local / HTMX pattern | Prioridad | Loom status | Evidencia |
|---|---|---|---|---|---|---|---|---|---|---|
| Button | Core | Sí | Sí | CSS | elevated, filled, tonal, outlined, text; icon | disabled, focus-visible, loading a definir; button/link semantics | Sin JS para visual; submit/action HTMX | P0 | NOT_STARTED | `material-web-upstream/button/{elevated,filled,filled-tonal,outlined,text}-button.ts`; `button/*_test.ts`; `docs/components/button.md`; `material-tailwind/src/components/button/button.css` |
| Checkbox | Core | Sí | Sí | CSS | checked, unchecked, indeterminate | disabled, focus, form association, label y mixed state | HTML nativo; validación/submit HTMX | P1 | NOT_STARTED | `material-web-upstream/checkbox/checkbox.ts`; `checkbox/checkbox_test.ts`; `docs/components/checkbox.md`; `material-tailwind/src/components/input/checkbox/checkbox.css` |
| Chips | Core | Sí | Sí | CSS | assist, filter, input, suggestion, chip-set | selected, disabled, removable; teclado/roles según uso | JS mínimo para selección local o estado server-side | P2 | NOT_STARTED | `material-web-upstream/chips/{assist,filter,input,suggestion}-chip.ts`; `chips/*_test.ts`; `docs/components/chip.md`; `material-tailwind/src/components/feedback/chip/chip.css` |
| Dialog | Core | Sí | Sí | CSS | standard, full-screen vía tokens | modal, open/closed, responsive; focus trap, Escape, restore focus, `aria-modal` | `hx-get` devuelve contenido; `<dialog>` + JS mínimo de lifecycle | P0 | NOT_STARTED | `material-web-upstream/dialog/dialog.ts`; `dialog/dialog_test.ts`; `docs/components/dialog.md`; `tokens/_md-comp-full-screen-dialog.scss`; `material-tailwind/src/components/feedback/dialog/dialog.css` |
| Divider | Core | Sí | Sí | CSS | horizontal, inset, vertical según markup/theme | decorativo vs semántico; orientación | Sin JS/HTMX | P2 | NOT_STARTED | `material-web-upstream/divider/divider.ts`; `divider/divider_test.ts`; `docs/components/divider.md`; `material-tailwind/src/components/layout/divider/divider.css` |
| Elevation | Core/foundation | Sí | Sí | CSS utility | niveles 0–5 | visual; no debe alterar semántica | Sin JS; utility/theme | P2 | NOT_STARTED | `material-web-upstream/elevation/elevation.ts`; `elevation/elevation_test.ts`; `docs/components/elevation.md`; `material-tailwind/src/{tokens,utilities}/elevation.css` |
| FAB | Core | Sí | Sí | CSS | standard, branded; tamaños/extended aparecen en tokens | disabled, focus; nombre accesible para icon-only | Acción normal/HTMX | P2 | NOT_STARTED | `material-web-upstream/fab/{fab,branded-fab}.ts`; `fab/fab_test.ts`; `docs/components/fab.md`; `material-tailwind/src/components/layout/fab/fab.css` |
| Field | Core primitive | Sí | Sí | Parcial (textfield CSS) | filled, outlined | label, supporting/error text, disabled; no usar sólo color | Referencia interna para Text field; no se distribuye separado | skip | NOT_STARTED | `material-web-upstream/field/{filled,outlined}-field.ts`; `field/*_test.ts`; `material-tailwind/src/components/input/textfield/textfield.css` |
| Focus ring | Core/foundation | Sí | Sí | CSS utility | inward/outward en upstream | `:focus-visible`, forced-colors, no reset global | Sin HTMX; JS no necesario en primera versión | P1 | NOT_STARTED | `material-web-upstream/focus/internal/focus-ring.ts`; `focus/internal/focus-ring_test.ts`; `docs/components/focus-ring.md`; `material-tailwind/src/utilities/focus-ring.css` |
| Icon | Core | Sí | Sí | No (claim sin CSS) | font/SVG strategy a definir | decorativo `aria-hidden` o nombre accesible | Sin HTMX | P2 | NOT_STARTED | `material-web-upstream/icon/icon.ts`; `icon/icon_test.ts`; `docs/components/icon.md`; ausencia de `material-tailwind/src/components/layout/icon/icon.css` |
| Icon button | Core | Sí | Sí | No | standard, filled, tonal, outlined; toggle | selected, disabled, focus; accessible name obligatorio | Acción/toggle local o HTMX | P2 | NOT_STARTED | `material-web-upstream/iconbutton/{icon-button,filled-icon-button,filled-tonal-icon-button,outlined-icon-button}.ts`; `iconbutton/*_test.ts`; `docs/components/icon-button.md` |
| List | Core | Sí | Sí | CSS | list/list-item; one/two/three-line visual | active, disabled, selected; semántica depende del contenido | Navegación o selección; fragmentos/paginación HTMX | P1 | NOT_STARTED | `material-web-upstream/list/{list,list-item}.ts`; `list/*_test.ts`; `docs/components/list.md`; `material-tailwind/src/components/data-display/list/list.css` |
| Menu | Core | Sí | Sí | CSS | menu, item, submenu | open, selected, disabled; roving focus, arrows, Escape | JS mínimo de popup/teclado; acciones HTMX | P1 | NOT_STARTED | `material-web-upstream/menu/{menu,menu-item,sub-menu}.ts`; `menu/menu_test.ts`; `docs/components/menu.md`; `material-tailwind/src/components/navigation/menu/menu.css` |
| Progress | Core | Sí | Sí | CSS | circular, linear; determinate/indeterminate | `role=progressbar`, value/min/max o label de loading | Polling o SSE progress | P1 | NOT_STARTED | `material-web-upstream/progress/{circular,linear}-progress.ts`; `progress/*_test.ts`; `docs/components/progress.md`; `material-tailwind/src/components/feedback/progress/progress.css` |
| Radio | Core | Sí | Sí | CSS | radio group | checked, disabled, focus; name/group semantics | HTML nativo; submit/validation HTMX | P2 | NOT_STARTED | `material-web-upstream/radio/radio.ts`; `radio/radio_test.ts`; `docs/components/radio.md`; `material-tailwind/src/components/input/radio/radio.css` |
| Ripple | Core/foundation | Sí | Sí | CSS aproximado | hover/press; pointer-origin sólo upstream JS | disabled, touch delay/cancel; respetar reduced motion | JS opcional; no es patrón HTMX | defer | NOT_STARTED | `material-web-upstream/ripple/internal/ripple.ts`; `ripple/internal/ripple_test.ts`; `docs/components/ripple.md`; `material-tailwind/src/utilities/ripple.css` |
| Select | Core | Sí | Sí | CSS | filled, outlined, option | open, selected, error, disabled; label/keyboard | Nativo primero; autocomplete remoto HTMX aparte | P1 | NOT_STARTED | `material-web-upstream/select/{filled-select,outlined-select,select-option}.ts`; `select/select_test.ts`; `docs/components/select.md`; `material-tailwind/src/components/input/select/select.css` |
| Slider | Core | Sí | Sí | CSS | single/range; tamaños nuevos sólo tokens | value, disabled, focus; keyboard y ARIA range | HTML range/JS local; HTMX al commit | defer | NOT_STARTED | `material-web-upstream/slider/slider.ts`; `slider/slider_test.ts`; `docs/components/slider.md`; `material-tailwind/src/components/input/slider/slider.css` |
| Switch | Core | Sí | Sí | CSS | selected/unselected; icon opcional | disabled, focus; role/switch naming | HTML checkbox; acción HTMX opcional | P2 | NOT_STARTED | `material-web-upstream/switch/switch.ts`; `switch/*_test.ts`; `docs/components/switch.md`; `material-tailwind/src/components/input/switch/switch.css` |
| Tabs | Core | Sí | Sí | CSS | primary, secondary; icon/text | active, disabled, focus; tablist/tab/tabpanel y arrows | `hx-get` + history + panel swap | P1 | NOT_STARTED | `material-web-upstream/tabs/{tabs,primary-tab,secondary-tab}.ts`; `tabs/tabs_test.ts`; `docs/components/tabs.md`; `material-tailwind/src/components/navigation/tabs/tabs.css` |
| Text field | Core | Sí | Sí | CSS | filled, outlined; textarea/prefix/suffix en API | label, helper, error, required, disabled; HTTP 422 | Form validation server-side y fragment swap | P0 | NOT_STARTED | `material-web-upstream/textfield/{filled,outlined}-text-field.ts`; `textfield/*_test.ts`; `docs/components/text-field.md`; `material-tailwind/src/components/input/textfield/textfield.css` |
| ARIA primitives | Labs | Sí | No | No | command, menu, tabs | patrones de asociación, teclado y roles | No portar como componente; usar como referencia | skip | NOT_STARTED | `material-web-upstream/labs/aria/{command,query-associated}.ts`; `labs/aria/menu/*`; `labs/aria/tabs/*`; tests en esos árboles |
| Badge | Labs | Sí | Sí | CSS | small/large o count/dot a validar | contenido decorativo vs anunciado | Sin HTMX; contador puede actualizarse por swap/SSE | P2 | NOT_STARTED | `material-web-upstream/labs/badge/badge.ts`; `labs/badge/badge_test.ts`; `material-tailwind/src/components/layout/badge/badge.css` |
| Behaviors | Labs | Sí | Sí | No | form-associated, validation, submitter, focusable | ElementInternals/Custom Element lifecycle | No portar runtime; traducir escenarios a HTML/Go | skip | NOT_STARTED | `material-web-upstream/labs/behaviors/{form-associated,form-submitter,constraint-validation}.ts`; respectivos `*_test.ts` |
| Card | Labs | Sí | Sí | CSS | elevated, filled, outlined | focus/click sólo si interactiva; no div clickeable sin semántica | Link/form/action HTMX según contenido | P2 | NOT_STARTED | `material-web-upstream/labs/card/{elevated,filled,outlined}-card.ts`; `labs/card/*_test.ts`; `material-tailwind/src/components/layout/card/card.css` |
| Item | Labs | Sí | Sí | No | primitive de slots | Semántica `UNKNOWN`; sin tests observados | No portar hasta existir contrato reusable | defer | NOT_STARTED | `material-web-upstream/labs/item/item.ts`; `labs/item/demo/*`; no `*_test.ts` en el directorio |
| Navigation bar | Labs | Sí | Sí | CSS | destinations; variantes locales no validadas | active, focus, badges; nav semantics | Navegación HTMX/history opcional | P2 | NOT_STARTED | `material-web-upstream/labs/navigationbar/navigation-bar.ts`; `labs/navigationbar/md-navigation-bar_test.ts`; `material-tailwind/src/components/navigation/navigation-bar/navigation-bar.css` |
| Navigation drawer | Labs | Sí | Sí | No | standard, modal | open/closed; focus trap, Escape, responsive | JS mínimo; carga/links HTMX opcionales | defer | NOT_STARTED | `material-web-upstream/labs/navigationdrawer/{navigation-drawer,navigation-drawer-modal}.ts`; tests/docs específicos `UNKNOWN` |
| Navigation tab | Labs | Sí | Sí | Parcial (tabs CSS) | destination tab | active/focus; nav vs tab semantics no deben mezclarse | History/navigation HTMX | P2 | NOT_STARTED | `material-web-upstream/labs/navigationtab/navigation-tab.ts`; `labs/navigationtab/md-navigation-tab_test.ts`; `material-tailwind/src/components/navigation/tabs/tabs.css` |
| Segmented button | Labs | Sí | Sí | No | outlined | selected, disabled, focus; single/multi semantics | Form/state HTMX opcional | defer | NOT_STARTED | `material-web-upstream/labs/segmentedbutton/outlined-segmented-button.ts`; tests/docs específicos `UNKNOWN` |
| Segmented button set | Labs | Sí | Sí | No | outlined set | selection model y keyboard; tests `UNKNOWN` | Local selection o submit HTMX | defer | NOT_STARTED | `material-web-upstream/labs/segmentedbuttonset/outlined-segmented-button-set.ts`; `labs/segmentedbuttonset/demo/*` |
| GB | Labs/experimental | Sí | No | No | appbar, badge, button, card, checkbox, chip, divider, FAB, focus, iconbutton, list, menu, radio y más | API experimental; tests mayormente shared/styles | No portar en bloque; cosechar tokens/patrones selectivos | skip | NOT_STARTED | `material-web-upstream/labs/gb/components/**`; `labs/gb/components/shared/*_test.ts`; `labs/gb/styles/adopt-styles_test.ts` |

Notas:

- La documentación estable vive en `material-web-upstream/docs/components/`; labs tiene cobertura desigual.
- `field` es una primitive interna de text field/select, no un control de formulario completo por sí solo.
- El toast del vertical slice no tiene componente core equivalente: usar snackbar como referencia de tokens/feedback y `material-tailwind/src/components/feedback/snackbar/snackbar.css`, pero diseñar contrato Loom propio con `HX-Trigger`.

### Campos obligatorios complementarios y evidencia literal

Esta tabla completa explícitamente tokens, complejidad, documentación/tests y rutas literales. `Baja`, `Media` y `Alta` describen la complejidad estimada de traducir a HTML + Tailwind, incluyendo comportamiento local; no son una afirmación de implementación terminada.

| Componente | Tokens relevantes | Complejidad de port | Documentación / tests observados | Evidencia literal upstream / fork / port CSS |
|---|---|---|---|---|
| Button | filled, outlined, text, tonal, elevation/state | Baja | `material-web-upstream/docs/components/button.md`; `material-web-upstream/button/filled-button_test.ts` | `material-web-upstream/button/filled-button.ts`; `material-web-tailwind/button/filled-button.ts`; `material-web-tailwind/material-tailwind/src/components/button/button.css`; `material-web-upstream/tokens/_md-comp-filled-button.scss` |
| Checkbox | checkbox + state/color | Media | `material-web-upstream/docs/components/checkbox.md`; `material-web-upstream/checkbox/checkbox_test.ts` | `material-web-upstream/checkbox/checkbox.ts`; `material-web-tailwind/checkbox/checkbox.ts`; `material-web-tailwind/material-tailwind/src/components/input/checkbox/checkbox.css`; `material-web-upstream/tokens/_md-comp-checkbox.scss` |
| Chips | assist/filter/input/suggestion | Media | `material-web-upstream/docs/components/chip.md`; `material-web-upstream/chips/filter-chip_test.ts` | `material-web-upstream/chips/filter-chip.ts`; `material-web-tailwind/chips/filter-chip.ts`; `material-web-tailwind/material-tailwind/src/components/feedback/chip/chip.css`; `material-web-upstream/tokens/_md-comp-filter-chip.scss` |
| Dialog | dialog/full-screen, scrim, elevation | Alta | `material-web-upstream/docs/components/dialog.md`; `material-web-upstream/dialog/dialog_test.ts` | `material-web-upstream/dialog/dialog.ts`; `material-web-tailwind/dialog/dialog.ts`; `material-web-tailwind/material-tailwind/src/components/feedback/dialog/dialog.css`; `material-web-upstream/tokens/_md-comp-dialog.scss` |
| Divider | divider color/thickness | Baja | `material-web-upstream/docs/components/divider.md`; `material-web-upstream/divider/divider_test.ts` | `material-web-upstream/divider/divider.ts`; `material-web-tailwind/divider/divider.ts`; `material-web-tailwind/material-tailwind/src/components/layout/divider/divider.css`; `material-web-upstream/tokens/_md-comp-divider.scss` |
| Elevation | system/component elevation | Baja | `material-web-upstream/docs/components/elevation.md`; `material-web-upstream/elevation/elevation_test.ts` | `material-web-upstream/elevation/elevation.ts`; `material-web-tailwind/elevation/elevation.ts`; `material-web-tailwind/material-tailwind/src/utilities/elevation.css`; `material-web-upstream/tokens/_md-sys-elevation.scss` |
| FAB | standard/branded FAB | Baja | `material-web-upstream/docs/components/fab.md`; `material-web-upstream/fab/fab_test.ts` | `material-web-upstream/fab/fab.ts`; `material-web-tailwind/fab/fab.ts`; `material-web-tailwind/material-tailwind/src/components/layout/fab/fab.css`; `material-web-upstream/tokens/_md-comp-fab.scss` |
| Field | filled/outlined field | Media; sólo referencia, no componente Loom independiente | `material-web-upstream/docs/components/text-field.md`; `material-web-upstream/field/filled-field_test.ts` | `material-web-upstream/field/filled-field.ts`; `material-web-tailwind/field/filled-field.ts`; `material-web-tailwind/material-tailwind/src/components/input/textfield/textfield.css`; `material-web-upstream/tokens/_md-comp-filled-field.scss` |
| Focus ring | focus ring + focus indicator | Media | `material-web-upstream/docs/components/focus-ring.md`; `material-web-upstream/focus/internal/focus-ring_test.ts` | `material-web-upstream/focus/internal/focus-ring.ts`; `material-web-tailwind/focus/internal/focus-ring.ts`; `material-web-tailwind/material-tailwind/src/utilities/focus-ring.css`; `material-web-upstream/tokens/_md-comp-focus-ring.scss` |
| Icon | icon size/color | Baja | `material-web-upstream/docs/components/icon.md`; `material-web-upstream/icon/icon_test.ts` | `material-web-upstream/icon/icon.ts`; `material-web-tailwind/icon/icon.ts`; port CSS exacto: `UNKNOWN`; `material-web-upstream/tokens/_md-comp-icon.scss` |
| Icon button | standard/filled/tonal/outlined | Media | `material-web-upstream/docs/components/icon-button.md`; `material-web-upstream/iconbutton/icon-button_test.ts` | `material-web-upstream/iconbutton/icon-button.ts`; `material-web-tailwind/iconbutton/icon-button.ts`; port CSS exacto: `UNKNOWN`; `material-web-upstream/tokens/_md-comp-icon-button.scss` |
| List | list/list-item | Media | `material-web-upstream/docs/components/list.md`; `material-web-upstream/list/list_test.ts` | `material-web-upstream/list/list.ts`; `material-web-tailwind/list/list.ts`; `material-web-tailwind/material-tailwind/src/components/data-display/list/list.css`; `material-web-upstream/tokens/_md-comp-list.scss` |
| Menu | menu/menu-item, elevation/state | Alta | `material-web-upstream/docs/components/menu.md`; `material-web-upstream/menu/menu_test.ts` | `material-web-upstream/menu/menu.ts`; `material-web-tailwind/menu/menu.ts`; `material-web-tailwind/material-tailwind/src/components/navigation/menu/menu.css`; `material-web-upstream/tokens/_md-comp-menu.scss` |
| Progress | circular/linear progress | Media | `material-web-upstream/docs/components/progress.md`; `material-web-upstream/progress/linear-progress_test.ts` | `material-web-upstream/progress/linear-progress.ts`; `material-web-tailwind/progress/linear-progress.ts`; `material-web-tailwind/material-tailwind/src/components/feedback/progress/progress.css`; `material-web-upstream/tokens/_md-comp-linear-progress.scss` |
| Radio | radio + state/color | Baja | `material-web-upstream/docs/components/radio.md`; `material-web-upstream/radio/radio_test.ts` | `material-web-upstream/radio/radio.ts`; `material-web-tailwind/radio/radio.ts`; `material-web-tailwind/material-tailwind/src/components/input/radio/radio.css`; `material-web-upstream/tokens/_md-comp-radio.scss` |
| Ripple | ripple/state/motion | Alta | `material-web-upstream/docs/components/ripple.md`; `material-web-upstream/ripple/internal/ripple_test.ts` | `material-web-upstream/ripple/internal/ripple.ts`; `material-web-tailwind/ripple/internal/ripple.ts`; `material-web-tailwind/material-tailwind/src/utilities/ripple.css`; `material-web-upstream/tokens/_md-comp-ripple.scss` |
| Select | filled/outlined select | Alta | `material-web-upstream/docs/components/select.md`; `material-web-upstream/select/select_test.ts` | `material-web-upstream/select/filled-select.ts`; `material-web-tailwind/select/filled-select.ts`; `material-web-tailwind/material-tailwind/src/components/input/select/select.css`; `material-web-upstream/tokens/_md-comp-filled-select.scss` |
| Slider | slider + state/motion | Alta | `material-web-upstream/docs/components/slider.md`; `material-web-upstream/slider/slider_test.ts` | `material-web-upstream/slider/slider.ts`; `material-web-tailwind/slider/slider.ts`; `material-web-tailwind/material-tailwind/src/components/input/slider/slider.css`; `material-web-upstream/tokens/_md-comp-slider.scss` |
| Switch | switch + state/color | Media | `material-web-upstream/docs/components/switch.md`; `material-web-upstream/switch/switch_test.ts` | `material-web-upstream/switch/switch.ts`; `material-web-tailwind/switch/switch.ts`; `material-web-tailwind/material-tailwind/src/components/input/switch/switch.css`; `material-web-upstream/tokens/_md-comp-switch.scss` |
| Tabs | primary/secondary tab | Alta | `material-web-upstream/docs/components/tabs.md`; `material-web-upstream/tabs/tabs_test.ts` | `material-web-upstream/tabs/tabs.ts`; `material-web-tailwind/tabs/tabs.ts`; `material-web-tailwind/material-tailwind/src/components/navigation/tabs/tabs.css`; `material-web-upstream/tokens/_md-comp-primary-tab.scss` |
| Text field | filled/outlined text-field | Alta | `material-web-upstream/docs/components/text-field.md`; `material-web-upstream/textfield/filled-text-field_test.ts` | `material-web-upstream/textfield/filled-text-field.ts`; `material-web-tailwind/textfield/filled-text-field.ts`; `material-web-tailwind/material-tailwind/src/components/input/textfield/textfield.css`; `material-web-upstream/tokens/_md-comp-filled-text-field.scss` |
| ARIA primitives | system state/type; sin component-token propio confirmado | Alta si se portara; recomendado skip | `material-web-upstream/labs/aria/menu/menuitem_test.ts`; `material-web-upstream/labs/aria/tabs/tab_test.ts` | `material-web-upstream/labs/aria/command.ts`; fork: `NO`; port CSS: `NO`; `material-web-upstream/tokens/_md-sys-state.scss` |
| Badge | badge | Baja | docs específicas: `UNKNOWN`; `material-web-upstream/labs/badge/badge_test.ts` | `material-web-upstream/labs/badge/badge.ts`; `material-web-tailwind/labs/badge/badge.ts`; `material-web-tailwind/material-tailwind/src/components/layout/badge/badge.css`; `material-web-upstream/tokens/_md-comp-badge.scss` |
| Behaviors | system state/form; sin component-token | Alta si se portara; recomendado skip | docs específicas: `UNKNOWN`; `material-web-upstream/labs/behaviors/form-associated_test.ts` | `material-web-upstream/labs/behaviors/form-associated.ts`; `material-web-tailwind/labs/behaviors/form-associated.ts`; port CSS: `NO`; `material-web-upstream/tokens/_md-sys-state.scss` |
| Card | elevated/filled/outlined card | Baja | `material-web-upstream/labs/card/demo/stories.ts`; `material-web-upstream/labs/card/elevated-card_test.ts` | `material-web-upstream/labs/card/elevated-card.ts`; `material-web-tailwind/labs/card/elevated-card.ts`; `material-web-tailwind/material-tailwind/src/components/layout/card/card.css`; `material-web-upstream/tokens/_md-comp-elevated-card.scss` |
| Item | generic item | Media | `material-web-upstream/labs/item/demo/stories.ts`; tests: `UNKNOWN` | `material-web-upstream/labs/item/item.ts`; `material-web-tailwind/labs/item/item.ts`; port CSS: `NO`; `material-web-upstream/tokens/_md-comp-item.scss` |
| Navigation bar | navigation bar | Media | `material-web-upstream/labs/navigationbar/demo/stories.ts`; `material-web-upstream/labs/navigationbar/md-navigation-bar_test.ts` | `material-web-upstream/labs/navigationbar/navigation-bar.ts`; `material-web-tailwind/labs/navigationbar/navigation-bar.ts`; `material-web-tailwind/material-tailwind/src/components/navigation/navigation-bar/navigation-bar.css`; `material-web-upstream/tokens/_md-comp-navigation-bar.scss` |
| Navigation drawer | navigation drawer | Alta | docs/tests específicos: `UNKNOWN` | `material-web-upstream/labs/navigationdrawer/navigation-drawer.ts`; `material-web-tailwind/labs/navigationdrawer/navigation-drawer.ts`; port CSS: `NO`; `material-web-upstream/tokens/_md-comp-navigation-drawer.scss` |
| Navigation tab | navigation/state tokens; component-token específico `UNKNOWN` | Media | docs específicas: `UNKNOWN`; `material-web-upstream/labs/navigationtab/md-navigation-tab_test.ts` | `material-web-upstream/labs/navigationtab/navigation-tab.ts`; `material-web-tailwind/labs/navigationtab/navigation-tab.ts`; `material-web-tailwind/material-tailwind/src/components/navigation/tabs/tabs.css`; `material-web-upstream/tokens/_md-comp-navigation-bar.scss` |
| Segmented button | outlined segmented button | Media | docs/tests específicos: `UNKNOWN` | `material-web-upstream/labs/segmentedbutton/outlined-segmented-button.ts`; `material-web-tailwind/labs/segmentedbutton/outlined-segmented-button.ts`; port CSS: `NO`; `material-web-upstream/tokens/_md-comp-outlined-segmented-button.scss` |
| Segmented button set | outlined segmented button | Alta | `material-web-upstream/labs/segmentedbuttonset/demo/stories.ts`; tests: `UNKNOWN` | `material-web-upstream/labs/segmentedbuttonset/outlined-segmented-button-set.ts`; `material-web-tailwind/labs/segmentedbuttonset/outlined-segmented-button-set.ts`; port CSS: `NO`; `material-web-upstream/tokens/_md-comp-outlined-segmented-button.scss` |
| GB | tokens propios por componente | Alta; recomendado skip como familia | `material-web-upstream/labs/gb/components/button/demo/stories.ts`; `material-web-upstream/labs/gb/components/shared/directives_test.ts` | `material-web-upstream/labs/gb/components/button/button.ts`; fork: `NO`; port CSS: `NO`; `material-web-upstream/labs/gb/components/button/_button-tokens.scss` |

## 5. Diferencias upstream vs fork

1. La API estable de foundations casi no cambió: `color/`, `common.ts`, gran parte de typography/elevation/focus/ripple y docs de theming son equivalentes; upstream agrega `typography/md-typescale-styles.ts`.
2. Upstream reorganiza v0.192 bajo `tokens/versions/v0_192` y agrega `tokens/versions/latest`; el fork conserva `tokens/v0_192`.
3. Upstream agrega `labs/aria` y `labs/gb`; el fork carece de ambos.
4. `form-submitter` pasa de `internal/controller/` en el fork a `labs/behaviors/` upstream.
5. Upstream reescribe `internal/events/dispatch-hooks.ts` y amplía sus tests; es lógica de Custom Elements, no core Loom.
6. Upstream 2.5.0 agrega pipeline separado para `labs/gb`, metadatos Sass y `custom-elements.json`; fork 2.3.0 no tiene `build:manifest` ni el campo `customElements`.
7. El fork tracked no contiene Tailwind. El port CSS es un árbol aparte, untracked y sin procedencia Git verificable.

## 6. Qué rescatar de `material-tailwind`

Rescatar sólo después de normalizar y verificar:

- La separación conceptual `tokens/`, `themes/`, `utilities/`, `components/` de `src/material-tailwind.css`.
- CSS de componentes como spike visual: button, textfield, checkbox, radio, switch, select, slider, tabs, menu, navigation bar, card, divider, FAB, badge, dialog, snackbar, progress, tooltip, chip, list y table.
- HTML de `examples/*.html` como materia prima de previews.
- Nombres de variantes y selectores que coincidan con evidencia upstream.
- Sombras precomputadas y typography como referencia, no copia automática.

No considerarlo librería lista:

- Usa Tailwind `^3.4.6`, no Tailwind 4 `@theme`.
- Faltan `scripts/build-plugin.js`, `scripts/build-tokens.js`, `scripts/test-output.js`, `plugin.js` y `dist/index.d.ts`, aunque `package.json` los declara.
- Hay referencias de tokens inconsistentes (`--md-sys-elevation-1` vs `--md-sys-elevation-level1`; nombres state-layer divergentes).
- No existe suite de tests de componentes; los checks son principalmente presencia de archivos/demos manuales.
- `PROJECT-STATUS.md` afirma componentes/ejemplos ausentes (`icon`, `surface`, `tree`, feedback/data-display examples), WCAG AA e IE11 sin evidencia suficiente.
- La SPA de docs no usa Markdown ni HTMX y su pipeline no es reproducible con los manifests observados.

Evidencia: `material-tailwind/package.json`, `src/**`, `docs/**`, `PROJECT-STATUS.md` y ausencia comprobada de las rutas declaradas.

## 7. Qué no portar

- LitElement, decorators y `CSSResult`.
- Shadow DOM y `:host` como contrato obligatorio.
- Custom Elements `md-*` como arquitectura.
- `AttachableController`, ElementInternals polyfills y controllers ligados al lifecycle de Web Components.
- Event dispatch hooks internos.
- La state machine completa de ripple dentro del core.
- Sass generators/validators como pipeline principal de Loom.
- Los 177 component tokens de `latest` dentro de `@theme`.
- Claims de compatibilidad, WCAG o completitud sin tests.

## 8. Traducción conceptual a Loom

| Material Web | Loom UI |
|---|---|
| Custom Element | HTML semántico + clases/atributos |
| Lit property | atributo HTML, input nativo o estado server-side |
| Lit event | evento HTML/HTMX; header de respuesta cuando corresponda |
| Shadow DOM CSS | Tailwind + tokens semánticos + selectores acotados |
| Component token Material | variable privada o mapping del `theme-material` |
| Imperative JS state | respuesta server-side/HTMX cuando sea remoto; JS local mínimo cuando el navegador lo exige |
| Web Test Runner | `go test`, `httptest`, tests contractuales HTML/ARIA y Playwright puntual |
| Catálogo Lit/Eleventy | docs Go dogfoodeadas con previews y endpoints reales |

HTMX no sustituye comportamiento estrictamente local como focus trap, medición/posicionamiento de popup, Escape o pointer-origin ripple. Ese comportamiento debe usar primitivas nativas o módulos JS pequeños y explícitos.

## 9. Matriz de progreso Loom

Estados permitidos:

- `NOT_STARTED`: sólo inventariado.
- `IN_REVIEW`: contrato/evidencia en revisión.
- `DESIGN_READY`: markup, tokens, estados, accesibilidad y patrón HTMX acordados.
- `IMPLEMENTED`: implementación presente.
- `DOCS_READY`: documentación y preview dogfoodeados.
- `TESTED`: gates automatizados aprobados.
- `BLOCKED`: dependencia o decisión impide avanzar.

Estado actual:

| Área | Estado | Próximo gate |
|---|---|---|
| Auditoría Material Web | IN_REVIEW | Revisión humana del presente informe y cierre de UNKNOWN críticos |
| Theme Material | NOT_STARTED | Usar API pública `tokens/` del commit fijado y sólo las familias enumeradas abajo |
| Button | NOT_STARTED | Contrato P0; icon slot con SVG inline, sin depender del sistema Icon P2 |
| Form/Input | NOT_STARTED | Text field nativo P0 + HTTP 422; `field/` queda sólo como referencia |
| Dialog | NOT_STARTED | `<dialog>` nativo + contrato P0 de focus/lifecycle definido abajo |
| Toast | NOT_STARTED | Componente Loom-only P0 + `HX-Trigger`; snackbar sólo referencia visual |
| Go/docs app | NOT_STARTED | Estructura y build reproducible |
| SQLite/SSE demo | NOT_STARTED | MVP visual/HTMX validado |
| Registry/CLI/themes extra | NOT_STARTED | Vertical slice y realtime validados |

## 10. Primer vertical slice recomendado

```text
button → form/input → dialog → toast
```

Este slice valida tokens, variantes, estados, labels/errors, requests HTMX, HTTP 422, fragmentos HTML, overlays/focus, `HX-Trigger`, feedback, accesibilidad y documentación.

### Contrato mínimo

- **Button:** primary, secondary, outline, disabled, loading, focus e icon slot. El icono P0 será SVG inline aportado por el consumidor, con `aria-hidden` cuando sea decorativo; el catálogo/sistema Icon sigue P2.
- **Form/Input:** un único componente Loom `Text field` sobre `<input>`/`<textarea>` nativos, con label, helper, error, disabled, validación server-side y respuesta 422. `field/` se usa como referencia Material, no como paquete independiente. Loom instalará un listener global `htmx:beforeSwap`: cuando `event.detail.xhr.status === 422`, asignará `event.detail.shouldSwap = true` y `event.detail.isError = false`; así el fragmento de error se intercambia en el target normal sin convertir 422 en un falso éxito HTTP.
- **Dialog:** `<dialog>` nativo, apertura con `hx-get`, fragmento Go, `showModal()`, Escape/cierre, foco inicial, restore al trigger, responsive y nombre/descripción ARIA. El comportamiento local forma parte del componente Dialog P0; Focus ring P1 no es una dependencia.
- **Toast:** componente propio de Loom; nombre de evento exacto `loom:toast`. El servidor responderá, por ejemplo, `HX-Trigger: {"loom:toast":{"type":"success","message":"Guardado"}}`; un listener global `document.body.addEventListener('loom:toast', handler)` leerá `event.detail.type/message` y renderizará la región `aria-live`. Incluye success/error/warning y auto-dismiss pausable. Snackbar Material sólo inspira tokens/visual.

### Tokens fijados para este slice

Fuente: API pública root `material-web-upstream/tokens/` del commit `23b638c1d2338b9fcaccef81cf54298cce9c27ed`. `versions/v0_192` se conserva sólo como trazabilidad del snapshot.

- system: color, typescale, shape, elevation y state;
- component: text/filled/outlined button, filled/outlined text field y dialog;
- Loom-only: focus ring semántico y toast (`--ui-toast-*`), sin copiar un componente inexistente;
- adopción selectiva de 34.0.21 para V1: únicamente `thickness: 3px` y `outer-offset: 2px` para un outline exterior nativo; se omite `inner-offset: -3px` porque Loom no dibujará el focus ring inward;
- diferidos: fixed roles, contrast schemes, emphasized type, variable-font axes, radii increased, springs y familias expresivas por tamaño.

Evidencia literal: `material-web-upstream/tokens/_md-sys-color.scss`, `_md-sys-typescale.scss`, `_md-sys-shape.scss`, `_md-sys-elevation.scss`, `_md-sys-state.scss`, `_md-comp-text-button.scss`, `_md-comp-filled-button.scss`, `_md-comp-outlined-button.scss`, `_md-comp-filled-text-field.scss`, `_md-comp-outlined-text-field.scss`, `_md-comp-dialog.scss` y `material-web-upstream/tokens/versions/latest/sass/_md-sys-state-focus-indicator.scss`.

## 11. HTMX/realtime backlog

1. Validación server-side y fragmentos 422.
2. Dialog remoto.
3. Toast/eventos con `HX-Trigger`.
4. Tabs + history.
5. Tables con filtros/paginación.
6. Multi-target responses sólo donde reduzcan complejidad.
7. SSE notifications.
8. SSE progress.
9. WebSockets opcionales cuando exista interacción bidireccional real.

## 12. Go MVP

Estructura objetivo posterior a la aprobación de la auditoría:

```text
cmd/loom/
internal/
web/
├── templates/
├── content/
└── static/
registry/
themes/
└── theme-material/
```

Stack y gates:

- Go `net/http`, `html/template`, `embed`.
- Markdown con renderer integrado a los mismos layouts/componentes.
- Tailwind 4 con un input CSS y un build reproducible.
- HTMX con endpoints reales y templates de fragmentos.
- `go test ./...`, `httptest`, contratos HTML/ARIA y Playwright sólo para keyboard/focus/dialog/swaps.
- Docs y playground construidos con los componentes documentados, sin duplicar markup.

Después del vertical slice:

- SQLite con WAL, `busy_timeout`, migrations y transactions.
- Publicar al broker sólo después de commit exitoso.
- Broker en memoria para una instancia.
- SSE con heartbeat, cancelación por request context y buffers acotados.
- Demo: dos pestañas; una modifica SQLite, Go publica, SSE notifica y HTMX actualiza la otra.

## 13. Próximos pasos priorizados

1. Revisar este informe; la verificación automática ya confirmó que sólo cambió este archivo dentro de Loom UI.
2. Aprobar las decisiones ya explicitadas para icon slot, Text field, dialog/focus, toast y source de tokens. La procedencia del port CSS no bloquea: V1 se reimplementa desde evidencia upstream y no copia ese CSS.
3. Fijar la API pública `tokens/` del commit auditado; usar `versions/v0_192` sólo para trazabilidad y adoptar de 34.0.21 únicamente focus 3px/2px en V1.
4. Definir tokens semánticos y `theme-material`.
5. Crear el MVP Go + Tailwind 4 + HTMX y docs dogfoodeadas.
6. Implementar button → form/input → dialog → toast.
7. Validar HTTP 422, focus, `HX-Trigger` y browser tests.
8. Agregar SQLite/WAL + broker + SSE y demo de dos pestañas.
9. Extraer registry reusable.
10. Recién entonces agregar `theme-basecoat`, `theme-neutral`, CLI y adapters.

## UNKNOWN y riesgos abiertos

- `tokens/versions/latest` es generado/inestable; validación interna exhaustiva de Google: `UNKNOWN`.
- Token generado no equivale a componente implementado.
- Procedencia exacta de cada CSS del port local: `UNKNOWN` por falta de historial.
- Contraste WCAG del port local no fue medido; claims de WCAG: `UNKNOWN`.
- 298 referencias de variables CSS sin definición local reportadas por análisis estático necesitan clasificación entre runtime inputs y errores.
- Variable-font axes con fonts estáticas: `UNKNOWN`.
- Motion springs `null` en Sass requieren otra representación.
- Tests/docs de item, navigation drawer y segmented buttons son incompletos o no observados.
- No se ejecutaron builds por la restricción de sólo lectura; los pipelines rotos se determinan por scripts/artefactos ausentes y manifests inconsistentes.

### Gates de cierre

| Riesgo/UNKNOWN | ¿Bloquea V1? | Método de cierre | Gate |
|---|---|---|---|
| Procedencia del port CSS | No, mientras no se copie | Reimplementar desde contratos/evidencia Apache-2.0 upstream; revisión legal antes de cualquier copia textual | Pre-commit del primer CSS reutilizado |
| Contraste de theme Material | Sí para declarar accesibilidad | Test automatizado de pares semánticos + revisión forced-colors | Antes de `DESIGN_READY` del theme |
| Dialog/focus | Decisión cerrada; falta probar | Playwright: apertura, foco inicial, Tab, Escape, restore al trigger | Antes de `TESTED` de Dialog |
| Iconos del button | No | SVG inline slot + tests de nombre accesible/decorativo; catálogo Icon queda P2 | Antes de `TESTED` de Button |
| Toast/snackbar | Decisión cerrada; falta probar | Test de payload `HX-Trigger`, `aria-live`, pausa y dismiss | Antes de `TESTED` de Toast |
| Variables CSS indefinidas del port | No para V1 | No importar el bundle; clasificar sólo si una pieza concreta se reutiliza | Revisión por componente |
| Latest/variable fonts/springs | No; diferidos | Spike separado y snapshot exacto si se promueven | Post-MVP |

## Operación realizada

- Se ejecutaron análisis Git, inventarios, lecturas de manifests/código/docs/tests y tres líneas paralelas de auditoría.
- Un subagente de componentes encontró un bloqueo de acceso read-only al port local; el agente principal completó y contrastó la matriz mediante inventarios locales read-only y rutas verificadas.
- No se emitieron operaciones de escritura sobre código tracked de `material-web-upstream` o `material-web-tailwind`; ambos estados tracked quedaron sin cambios.
- `material-tailwind/` es un árbol untracked y no se capturó un manifiesto interno pre-run: Git sólo prueba que siguió apareciendo como una entrada untracked, no que cada byte interno permaneciera idéntico. Los tres scripts Node auditados fueron ejecutados como checks, pero esta auditoría no afirma un hash pre/post del árbol interno.
- Documento actualizado: `D:\repos\loom-ui\MATERIAL-WEB-PROGRESS.md`.

### Comandos y comprobaciones realmente ejecutados

Los comandos se ejecutaron desde las rutas indicadas, en modo read-only salvo `git fetch`, que sólo actualizó metadata remota de `.git` en el clone upstream.

```text
# cwd: D:\repos\material-web-upstream
git status --short
git remote -v
git branch --show-current
git log -1 --format=%H%n%cI%n%s
git ls-files | wc -l
git fetch --prune origin
git rev-parse origin/main
git rev-list --count HEAD..origin/main
git ls-tree --name-only HEAD
git ls-tree --name-only HEAD:labs
git ls-tree --name-only HEAD:tokens
git ls-tree -r --name-only HEAD:tokens/versions/latest | wc -l
git ls-tree -r --name-only HEAD:tokens/versions/v0_192 | wc -l
git ls-files *_test.* */test/* | wc -l
git ls-files docs/*.md docs/**/*.md | wc -l

# cwd: D:\repos\material-web-tailwind
git status --short
git remote -v
git branch --show-current
git log -1 --format=%H%n%cI%n%s
git ls-files | wc -l

# cwd: D:\repos\loom-ui
date --iso-8601=seconds

# cwd: D:\repos\material-web-tailwind\material-tailwind
node build-checker.js

# cwd: D:\repos\material-web-tailwind\material-tailwind\docs
node scripts/test-theme-system.js
node scripts/verify-complete-system.js
```

`node build-checker.js` terminó con exit code 0, pero su salida reportó 1 error y 46 warnings, incluido `docs/src/pages/ThemeBuilderPage.js`; por eso se registra como diagnóstico, no como gate aprobado.

También se ejecutó este bloque literal con `execute_code`; contiene rutas, entradas y validaciones concretas (sin placeholders):

```python
from pathlib import Path
import hashlib, json, re, subprocess
repos = Path(r'D:\repos')
loom = repos / 'loom-ui'
upstream = repos / 'material-web-upstream'
fork = repos / 'material-web-tailwind'
report = loom / 'MATERIAL-WEB-PROGRESS.md'
text = report.read_text(encoding='utf-8')
expected_unchanged = {
    'PROMPT-MATERIAL-WEB-INVENTORY.md':
        '696e40db84bcf8d3942cf45e3ff2c7a986b47ba0bca0d8cdd288bdc13bcf227c',
    'README.md':
        '251e59e543a42fa652cf89d25f139a83df9ab0278e4b238cc697b57d1f95dd34',
}
loom_hashes = {}
for name in [
    'MATERIAL-WEB-PROGRESS.md',
    'PROMPT-MATERIAL-WEB-INVENTORY.md',
    'README.md',
]:
    path = loom / name
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    loom_hashes[name] = {
        'size': path.stat().st_size,
        'sha256': digest,
        'unchanged_from_baseline': (
            expected_unchanged.get(name) == digest
            if name in expected_unchanged else None
        ),
    }

def git_index(path):
    raw = subprocess.check_output(
        ['git', '-C', str(path), 'ls-files', '-s'], text=True
    )
    return {
        line.split(None, 3)[3]: line.split()[1]
        for line in raw.splitlines()
    }

upstream_index, fork_index = git_index(upstream), git_index(fork)
common = set(upstream_index) & set(fork_index)
comparison = {
    'upstream': len(upstream_index),
    'fork': len(fork_index),
    'common': len(common),
    'same_common': sum(upstream_index[p] == fork_index[p] for p in common),
    'changed_common': sum(upstream_index[p] != fork_index[p] for p in common),
    'upstream_only': len(set(upstream_index) - set(fork_index)),
    'fork_only': len(set(fork_index) - set(upstream_index)),
}
required = [f'## {number}.' for number in range(1, 14)]
primary = text.split('| Componente | Categoría |', 1)[1].split('Notas:', 1)[0]
supplement = text.split(
    '### Campos obligatorios complementarios y evidencia literal', 1
)[1].split('## 5.', 1)[0]
literal_paths = {
    value for value in re.findall(r'`([^`\n]+)`', supplement)
    if value.startswith('material-web')
    and not any(mark in value for mark in ['*', '{', '}'])
}
result = {
    'loom_files': loom_hashes,
    'git_comparison': comparison,
    'required_sections_missing': [
        heading for heading in required if heading not in text
    ],
    'primary_component_rows': sum(
        1 for line in primary.splitlines()
        if line.startswith('| ') and 'NOT_STARTED' in line
    ),
    'supplement_component_rows': sum(
        1 for line in supplement.splitlines()
        if line.startswith('| ') and not line.startswith('| Componente |')
    ),
    'literal_paths_checked': len(literal_paths),
    'missing_literal_paths': sorted(
        value for value in literal_paths if not (repos / value).exists()
    ),
    'upstream_status': subprocess.check_output(
        ['git', '-C', str(upstream), 'status', '--short'], text=True
    ).splitlines(),
    'fork_status': subprocess.check_output(
        ['git', '-C', str(fork), 'status', '--short'], text=True
    ).splitlines(),
    'declared_node_scripts_exist': {
        'build-checker.js':
            (fork / 'material-tailwind' / 'build-checker.js').exists(),
        'docs/scripts/test-theme-system.js':
            (fork / 'material-tailwind/docs/scripts/test-theme-system.js').exists(),
        'docs/scripts/verify-complete-system.js':
            (fork / 'material-tailwind/docs/scripts/verify-complete-system.js').exists(),
    },
}
print(json.dumps(result, indent=2, ensure_ascii=False))
```

Resultado real de esa ejecución:

```text
required_sections_missing=[]
primary_component_rows=32
supplement_component_rows=32
literal_paths_checked=166
missing_literal_paths=[]
PROMPT-MATERIAL-WEB-INVENTORY.md unchanged_from_baseline=true
README.md unchanged_from_baseline=true
upstream_status=[]
fork_status=[las mismas seis entradas untracked del baseline]
declared_node_scripts_exist={los tres: true}
git_comparison={1493, 1029, 1018, 733, 285, 475, 11}
```

Las lecturas de archivos, búsquedas y patches se hicieron mediante herramientas Hermes (`read_file`, `search_files`, `patch`, `write_file`), no mediante comandos shell. No se ejecutaron `npm install`, `npm build`, `npm test`, generadores, formatters ni builds que escribieran en las fuentes.

### Prueba de aislamiento

`D:\repos\loom-ui` no es un repositorio Git, por lo que se capturó un manifiesto SHA-256 antes de escribir:

| Archivo | SHA-256 inicial | Resultado final |
|---|---|---|
| `PROMPT-MATERIAL-WEB-INVENTORY.md` | `696e40db84bcf8d3942cf45e3ff2c7a986b47ba0bca0d8cdd288bdc13bcf227c` | Sin cambios; hash idéntico |
| `README.md` | `251e59e543a42fa652cf89d25f139a83df9ab0278e4b238cc697b57d1f95dd34` | Sin cambios; hash idéntico |
| `MATERIAL-WEB-PROGRESS.md` | `9610053b526bcf5f2f319de2dcd7c837191c387049a58dc10967bb771071b708` | Actualizado por esta auditoría |

Baseline de fuentes antes de escribir:

- upstream: `git status --short` vacío;
- fork: exactamente `.claude/`, `catalog/pnpm-lock.yaml`, `material-design-icons-master.zip`, `material-design-icons-master/`, `material-tailwind/` y `temp_focus.txt` como untracked.

La comprobación posterior conservó upstream limpio y exactamente las mismas seis entradas untracked en el fork.
