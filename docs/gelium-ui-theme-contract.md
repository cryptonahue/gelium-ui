# Gelium UI — Theme Contract

> Contrato de themes del sistema Gelium UI.
> Fases A–C del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/handoffs/theme-architecture-audit.md`, `docs/handoffs/core-audit.md`, `docs/handoffs/basecoat-audit.md`.
> Gate duro: **no se implementa ningún theme nuevo (incluido Basecoat) antes de cerrar las fases A–H.**

---

## 1. Qué es un theme en Gelium UI

Un theme es una **dirección visual** codificada como tokens `--ui-*`. NO es markup, NO es JS, NO es una copia de otro sistema, NO es una dependencia npm.

```text
Button Gelium
  ├── Gelium UI Material    (tokens Material)
  └── Gelium UI Basecoat    (tokens Basecoat)
```

Nunca:

```text
Button Material        ← no: dos componentes separados por theme
Button Basecoat
```

Reglas del theme:

1. Define tokens `--ui-*`.
2. Conserva semántica HTML (los templates son únicos y no cambian).
3. Conserva la API de componentes (clases `ui-*`, variantes, estados).
4. Conserva contratos HTMX y server-driven.
5. Cubre estados y variantes (superficie de estados estable en el CSS de componentes).
6. Soporta light/dark.
7. Documenta divergencias.
8. Funciona con la misma documentación y suite de tests.

## 2. Mecanismo (acordado — Phase H)

- **Bundle de todos los themes + selección en runtime por clase/attr en el documento raíz** (`<html>` o `<body>`).
- Un `theme.css` por theme en `themes/<theme>/`, autocontenido (luz + dark), con selector raíz propio `.theme-<name>`.
- `web/styles/app.css` importa CADA theme explícitamente (CSS no globbea): un solo `web/static/app.css` embebido, todos los themes adentro.
- Selección data-driven desde template/server: `<html class="theme-material">` deja de estar hardcodeado (`layout.html:2`).
- **Mecanismo mínimo (Phase H)**: `<body class="theme-material">` / `<body class="theme-basecoat">`; selección desde el servidor o documento raíz. **Sin runtime JavaScript de selección de themes.**
- NO: archivos `themes/<theme>/<componente>.css` (el contrato es token-only; los componentes son únicos).
- La variable de build `THEME=basecoat npm run build` es una optimización opcional de footprint, no el mecanismo base.
- Dark: **una sola rutina por theme** (eliminar la duplicación clase + media query actual con drift, `theme.css:203-299`); decidir `light-dark()` vs clase única al implementar.

El contrato es: **mismo markup, mismo HTML semántico, misma API, mismo server contract, apariencia distinta.**

## 3. Tokens obligatorios

Un theme DEBE definir todas las familias que los componentes referencian. Contrato de cobertura:

### 3.1 Familias transversales (core)

| Familia | Tokens | Nota |
|---|---|---|
| Color semántico | `--ui-color-{canvas,surface,surface-container,primary,secondary,error,warning,success,info,outline,outline-strong,scrim,focus-ring}` + variantes `-fg` | roles, no marcas |
| Tipografía | `--ui-font-{sans,mono}` + `--ui-type-*` (steps display/headline/title/body/label) | cerrar gaps `display-lg`, `title-md` |
| Radius | `--ui-radius-{none,xs,sm,md,lg,xl,full}` | |
| Elevation | `--ui-shadow-0..5` | |
| Border | `--ui-border-width-*`, `--ui-border-style-*` | |
| Focus | `--ui-focus-{thickness,offset}` | |
| Motion | `--ui-motion-{short,medium,long}` + `--ui-easing-{standard,emphasized,decelerate,accelerate}` | |
| State | `--ui-state-{hover,focus,pressed,dragged,disabled}-opacity` | |
| Z-index | `--ui-z-{dialog,menu,popover,toast}` | |
| Breakpoints | `--ui-breakpoint-{sm,md,lg}` | |
| Spacing | `--ui-space-0..N` | |
| Density/size | `--ui-density-*`, `--ui-size-{control,item,icon}` | |

### 3.2 Familias por componente

Todo componente con tokens en el theme (button, text-field, dialog, toast, card, badge, checkbox, radio, switch, slider, progress, fab, select, divider, icon, elevation, focus-ring) DEBE tener su familia definida por el theme.

### 3.3 Componentes con tokens scoped (propiedad decidida en Phase A)

List, Menu, Data table, Navigation bar, Navigation tab, Navigation drawer, Segmented button, Tooltip declaran tokens scoped a su raíz. **Decisión Phase A (token ownership)**: los valores de anatomía (alturas, tamaños, padding de items) viven en el **core** como defaults scoped; el theme **puede sobrescribir** grados de libertad (color, shape, type, motion, densidad) declarándolos en `.theme-<name>` cuando el contrato lo declara. Ver matriz completa en `docs/gelium-ui-core.md` §5.0.

## 4. Tokens opcionales

- `--ui-card-fg` (nuevo, sugerido por mapeo Basecoat).
- Tokens por componente adicionales documentados como opcionales.
- Todo token opcional DEBE tener fallback o default en el core.

## 5. Naming

- Convención: `--ui-<familia>-<token>` (transversal) y `--ui-<componente>-<rol>` (cobertura).
- Unificar naming pendiente: `--ui-color-error` vs `--ui-color-danger` → **`danger`** es el canónico (ya usado por badge/toast); `--ui-color-surface-container` se define formalmente.
- Sin tokens muertos: eliminar/justificar `--ui-radius-xl`, `--ui-state-dragged-opacity`, `--ui-select-menu-item-icon`.
- Clases públicas solo `ui-*`: renombrar `m3-select-trigger` → `ui-select-trigger`; sin prefijos de terceros.
- `--ui-font-mono` se define formalmente en el core.

## 6. Light/dark

- Cada theme define luz + dark autocontenido.
- Un solo mecanismo dark (sin duplicación con drift). Basecoat valida el patrón de clase única (`.dark`); el contrato puede adoptar `light-dark()` si el soporte objetivo lo permite, o clase única + media query generada.
- `color-scheme` consistente: una sola fuente de verdad (hoy `base.css:3` y `theme.css:7` discrepan).

## 7. Component overrides

- Los componentes NO cambian por theme. Cero overrides por theme en `web/styles/`.
- Si un theme necesita una variante nueva (p. ej. `destructive` de Button, pills de Badge en Basecoat), se resuelve como **decisión del theme contract**: extender el componente en el core (con contrato), o documentar divergencia. Nunca CSS de theme sobre markup distinto.

## 8. Anatomía que NO puede cambiar

- Semántica HTML (artículo/lista/tabla/dialog/controls nativos).
- Clases de componente `ui-*` y variantes existentes.
- Contratos server-driven (422, `loom:toast`, GET params).
- Comportamiento no-JS end-to-end.
- Estructura de estados del componente (hover/focus/pressed/disabled/selected/error).

## 9. Tests obligatorios

La suite DEBE ser theme-agnóstica (hoy 13 tests hardcodean la ruta `theme-material/theme.css` y 3 aseveran valores hex Material — refactor obligatorio):

1. **Cobertura de tokens**: cada `themes/<theme>/theme.css` define TODO token referenciado por los componentes (`var(--ui-*)` sin gaps).
2. **Presencia, no valor**: assert de existencia de familia/cobertura light+dark, NO de valores hex (eliminar `styles_fab_test.go:85-87`, `styles_dialog_test.go:15-29`, `styles_toast_test.go:37-42` como aserciones de valor).
3. **Dark en ambas rutas** sin asumir valores.
4. **Matriz theme × component × variant × state** sobre el asset compilado (Phase 7).
5. **Sin valores hardcodeados**: ningún componente usa color de estado/geometría de control fijo (detectar regresión de `button.css:17-18`, `chips.css:32,35,107,141-142`).
6. `styles_contract_test.go` en sync con el nuevo `app.css` (lista `sourceAppCSS` desactualizada a corregir).

## 10. Documentación obligatoria por theme

Cada theme publica `themes/<theme>/README.md` con:

- Dirección visual y referencia (p. ej. Basecoat style pack Vega).
- Token mapping (origen → `--ui-*`).
- Light/dark soportado.
- Variantes/estados cubiertos.
- Divergencias documentadas (componente × divergencia × decisión).
- Cobertura de componentes (matriz).

## 11. Procedimiento para agregar un theme

1. Inspeccionar la referencia (auditoría read-only; `docs/handoffs/basecoat-audit.md` es el template).
2. Cerrar gaps del core y naming (Phase 1) — obligatorio antes.
3. Elegir la dirección visual concreta (style pack si la referencia tiene varios).
4. Traducir valores al vocabulario `--ui-*` (conversión de color decidida: hex por defecto; oklch si se aprueba).
5. Derivar familias que la referencia no tokeniza (elevation, motion, type, state, spacing) del CSS de la referencia.
6. Escribir `themes/<theme>/theme.css` (luz + dark autocontenido).
7. Integrar el import en `app.css` y activar por clase en runtime.
8. Cubrir la matriz de componentes/variantes/estados del scope.
9. Correr suite completa (tests theme-agnósticos) + build + vet.
10. Smoke: light/dark, narrow/wide, reduced motion, forced colors, no-JS, HTMX, estados empty/loading/error.
11. Aceptación visual manual.
12. Escribir `themes/<theme>/README.md`.

## 12. Definition of done de un theme

- [ ] Todos los tokens obligatorios definidos (cobertura §3) sin gaps.
- [ ] Naming canónico `--ui-*` respetado.
- [ ] Light + dark funcionando sin duplicación manual.
- [ ] Suite de tests theme-agnóstica pasando (presencia, no valor).
- [ ] `npm run build` + `go test ./...` + `go vet ./...` verdes.
- [ ] Matriz theme × component × variant × state ejecutada.
- [ ] Smoke completo en browser (consola sin errores, keyboard, no-JS, HTMX, light/dark, narrow/wide, reduced motion, forced colors).
- [ ] Aceptación visual manual.
- [ ] `themes/<theme>/README.md` con divergencias documentadas.
- [ ] Los componentes y contratos del core NO se modificaron para acomodar el theme.

---

**Gate**: esta fase (A–C) cierra antes de Phase I (Basecoat). Un theme que requiera cambiar markup, clases, contratos o comportamiento no-JS es un cambio de core, no un theme, y se trata como tal.
