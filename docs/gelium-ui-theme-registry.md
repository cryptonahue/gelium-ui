# Gelium UI — Theme Registry

> Registry de themes del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Fuentes de autoridad: `docs/gelium-ui-theme-contract.md` (contrato), `docs/gelium-ui-theme-implementation-guide.md` (procedimiento), `docs/handoffs/theme-architecture-audit.md` (auditoría), `docs/handoffs/basecoat-audit.md` (referencia Basecoat).

---

## 1. Cómo funciona el sistema de themes

Un theme es una **dirección visual codificada como tokens `--ui-*`**. NO es markup, NO es JS, NO es una copia de otro sistema (`theme-contract.md` §1).

- Un `theme.css` por theme en `themes/<theme>/`, autocontenido (luz + dark), con selector raíz propio `.theme-<name>`.
- `web/styles/app.css` importa CADA theme explícitamente (CSS no globbea); el build produce UN solo `web/static/app.css` embebido con todos los themes adentro (`//go:embed static/*` en `web/assets.go:8`).
- Selección en runtime por clase en el documento raíz — hoy `layout.html:2` hardcodea `<html lang="en" class="theme-material">`; el contrato exige data-driven (Phase A pendiente).
- Dark: UNA rutina por theme (eliminar la duplicación clase + media query con drift actual, `theme.css:203-299`).

## 2. Themes disponibles

| Theme | Directorio | Tokens `--ui-*` | Light | Dark | Estado | Selector raíz |
|---|---|---|---|---|---|---|
| **theme-material** | `themes/theme-material/theme.css` | 269 definiciones (157 únicos) | `.theme-material` | `.theme-material.theme-dark` / `.dark` / `[data-theme="dark"]` + `@media (prefers-color-scheme: dark)` (duplicación con drift conocido: `--ui-switch-track-unselected` falta en la media query) | **Implementado, default** | `.theme-material` |
| theme-basecoat | `themes/theme-basecoat/` | — | — | — | **Pendiente (Phase I)** | `.theme-basecoat` |

### Familias de tokens de theme-material (por tamaño)

| Familia | # tokens | Nota |
|---|---|---|
| color (semántico) | 69 | roles + `-fg` + containers (danger/success/warning/info); `--ui-color-error` es alias de `danger` (core `tokens.css:38`) |
| fab | 26 | `--ui-fab-{primary,surface,secondary}-{container,fg}` + geometría |
| toast | 22 | `--ui-toast-*` |
| switch | 17 | `--ui-switch-*` |
| select | 16 | `--ui-select-*` |
| field (text field) | 15 | `--ui-field-*` |
| dialog | 13 | `--ui-dialog-*` |
| type | 12 | typescale `--ui-type-*` |
| slider | 12 | `--ui-slider-*` |
| size / checkbox / card / radio / progress / shadow / radius / state / badge / motion / font / focus / divider / easing | resto | — |

> Nota del audit: **8 familias de componentes viven scoped en el CSS, NO en el theme**: List, Menu, Data table, Navigation bar, Navigation tab, Navigation drawer, Segmented button, Tooltip (`theme-architecture-audit.md` §3.4). Un theme nuevo solo puede tocarlas declarando los tokens globalmente en `.theme-<name>`.

## 3. Contrato que debe cumplir un theme

El contrato completo está en `docs/gelium-ui-theme-contract.md`. Resumen de lo NO negociable:

1. **Mismo markup, misma API, mismo server contract** — solo cambia la apariencia (`theme-contract.md` §2).
2. **Tokens obligatorios**: familias transversales (color, type, radius, elevation, border, focus, motion, state, z-index, breakpoints, spacing, density/size) + la familia por componente de cada componente del scope (§3).
3. **Naming canónico `--ui-*`**, `danger` es el canónico (no `error`), sin tokens muertos (§5).
4. **Light + dark autocontenido** con UNA rutina, sin duplicación textual (§6).
5. **Cero overrides por theme en `web/styles/`** — los componentes no cambian por theme (§7).
6. **Tests theme-agnósticos** (presencia, no valor) + matriz theme × component × variant × state (§9).
7. `themes/<theme>/README.md` con dirección visual, token mapping, light/dark, divergencias y matriz de cobertura (§10).

## 4. Procedimiento para agregar un theme

### Mecanismo Phase H (2 pasos mínimos)

1. **Importar el theme en el bundle**: agregar `@import "../../themes/<theme>/theme.css";` en `web/styles/app.css` (lista explícita, junto a los otros themes).
2. **Seleccionarlo en runtime por clase**: `<html class="theme-<name>">` (o `data-theme="<name>"`) en `layout.html`, data-driven desde el server cuando Phase A lo haga.

Con esto, cambiar la clase cambia la dirección visual **sin rebuild ni JS** (`theme-contract.md` §2, `theme-architecture-audit.md` §7).

### Procedimiento completo (12 pasos)

El flujo operativo completo (inspección → token mapping → implementación → cobertura → variantes/estados → dark/light → responsive → reduced motion/forced colors → documentación → tests → smoke → aceptación) vive en **`docs/gelium-ui-theme-implementation-guide.md`** — NO se duplica aquí. Precondiciones de gate (STOP/BLOCKED si fallan): Phase 1 core cerrada, contrato aprobado, tests theme-agnósticos, clase data-driven, audit read-only de la referencia (`implementation-guide.md` §0).

## 5. Matriz theme × component × variant × state (para Basecoat)

> Estado: **pendiente** — se ejecuta en Phase I (Basecoat) y en Phase 7 del roadmap (2+ themes). Esta matriz es el contrato de qué debe cubrir Basecoat, derivado del scope de Phase I (`roadmap.md:274-282`): Button, Text field, Card, Badge, Dialog, Toast, Data table.

| Componente | Variantes | Estados a cubrir | Tokens en theme-material | Basecoat pendiente |
|---|---|---|---|---|
| Button | primary/secondary/outline/text, link, loading | rest/hover/focus/pressed/disabled/loading | `--ui-button-*` + color semántico | ✅ |
| Text field | filled/outlined | rest/hover/focus/disabled/error | `--ui-field-*` | ✅ (decisión: floating label como variante de theme, `theme-contract.md:103`) |
| Card | elevated/filled/outlined | rest/hover/focus/pressed | `--ui-card-*` (+ `--ui-card-fg` opcional) | ✅ |
| Badge | dot/count/large + tones error/success/warning/info | nunca color-only | `--ui-badge-*` | ✅ |
| Dialog | confirm/page | open/closed, light dismiss | `--ui-dialog-*` | ✅ |
| Toast | info/success/warning/error | transitorio auto-dismiss | `--ui-toast-*` | ✅ |
| Data table | sortable/selection/pagination | rest/hover/focus/pressed/selected/empty/error | `--ui-data-table-*` (scoped, NO en theme) | ✅ (requiere decidir propiedad de familia scoped, `theme-architecture-audit.md` §3.4) |

**Gaps conocidos a cerrar antes de Basecoat** (del audit, §3.5/§3.6):

- `--ui-color-surface-container`, `--ui-type-display-lg`, `--ui-type-title-md`, `--ui-font-mono` ya cerrados en core (`tokens.css:23,145,146,138`).
- `--ui-color-error` vs `danger`: alias temporal en `tokens.css:38`, canónico `danger`.
- State layers hardcodeados (`button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`) → tokenizar con `color-mix`.
- Dark duplicado con drift (`theme.css:203-299`) → una rutina.
- 13 tests hardcodean la ruta `theme-material/theme.css`; 3 aseveran hex Material → parametrizar.

## 6. Divergencias documentadas de theme-material

- Toast conserva `loom:toast` + `aria-live` + fallback no-JS (`theme-contract.md:88`, decisión Phase I).
- Floating label de Text field es patrón Gelium; Basecoat decide si es variante de theme o divergencia documentada (`implementation-guide.md` §4.2).
- Variantes que un theme nuevo aporte y Gelium no tenga (p. ej. `destructive` Button, pills de Badge) se resuelven como decisión del contrato: extender el core o documentar divergencia — nunca CSS de theme sobre markup distinto (`theme-contract.md` §7).

---

**Definición de done (Phase J)**: registry refleja el estado real (un theme implementado, Basecoat pendiente con matriz de scope), el procedimiento referencia los 2 pasos del mecanismo Phase H y el guía completo sin duplicarlo, y los gaps pre-Basecoat están listados con evidencia.
