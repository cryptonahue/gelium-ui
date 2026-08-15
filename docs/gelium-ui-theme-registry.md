# Gelium UI — Theme Registry

> Registry de themes del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Fuentes de autoridad: `docs/gelium-ui-theme-contract.md` (contrato), `docs/gelium-ui-theme-implementation-guide.md` (procedimiento), `docs/handoffs/theme-architecture-audit.md` (auditoría), `docs/handoffs/basecoat-audit.md` (referencia Basecoat).

---

## 1. Cómo funciona el sistema de themes

Un theme es una **dirección visual codificada como tokens `--ui-*`**. NO es markup, NO es JS, NO es una copia de otro sistema (`theme-contract.md` §1).

- Un `theme.css` por theme en `themes/<theme>/`, autocontenido (luz + dark), con selector raíz propio `.theme-<name>`.
- `web/styles/app.css` importa CADA theme explícitamente (CSS no globbea); el build produce UN solo `web/static/app.css` embebido con todos los themes adentro (`//go:embed static/*` en `web/assets.go:8`).
- Selección en runtime por clase en el documento raíz — `layout.html:2` es data-driven (`class="{{.ThemeClass}}"`, allowlist server-side en `themeClass()`, `internal/app/server.go`). Cambiar la clase cambia la dirección visual sin rebuild ni JS.
- Dark: cada theme declara UNA rutina por esquema (luz + dark autocontenidos). La cobertura dark usa la **ruta de clase única** (`.theme-dark`/`.dark`/`[data-theme="dark"]`) — sin `@media (prefers-color-scheme: dark)` (unificación Phase A, sin drift).

## 2. Themes disponibles

| Theme | Directorio | Tokens `--ui-*` | Light | Dark | Estado | Selector raíz |
|---|---|---|---|---|---|---|
| **theme-material** | `themes/theme-material/theme.css` | 269 definiciones (169 únicos) | `.theme-material` | `.theme-material.theme-dark` / `.dark` / `[data-theme="dark"]` (clase única, sin media) | **Implementado, default** | `.theme-material` |
| **theme-basecoat** | `themes/theme-basecoat/theme.css` | 269 definiciones (167 únicos) | `.theme-basecoat` | `.theme-basecoat.theme-dark` / `.dark` / `[data-theme="dark"]` (clase única, sin media) | **Implementado (Phase I)** | `.theme-basecoat` |

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

### Familias de tokens de theme-basecoat (por tamaño)

| Familia | # tokens | Nota |
|---|---|---|
| color (semántico) | 23 | roles + `-fg` + containers; dirección Vega (oklch → hex); `danger` canónico, sin alias `error` |
| switch | 15 | 32×18.4px track (Basecoat), pressed sin crecimiento |
| select / fab | 14 | select 2.25rem→`--ui-size-field`; fab derivado de primary/secondary surfaces |
| type | 12 | typescale derivada de los steps Tailwind de Vega (`text-sm` etc.) |
| slider / radius / field / dialog / toast | 10 / 6 / 5 / 5 / 8 | anatomía Basecoat (h-1.5 tracks, 0.625rem base radius, toast popover) |
| state | 5 | hover/pressed .20 (≈ bg-primary/80), disabled .50 (opacity-50) |
| shadow / border / card / badge / focus / motion / size / font / radio / checkbox / progress / divider / easing | resto | — |

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
2. **Seleccionarlo en runtime por clase**: `<html class="theme-<name>">` (o `data-theme="<name>")` en `layout.html` (línea 2), data-driven desde el server (`{{.ThemeClass}}`, allowlist en `themeClass()`, `internal/app/server.go`). **Placement: `<html>` es el documento raíz** — el literal `<body>` del roadmap no se adopta (decisión Phase H, `theme-contract.md` §2).

Con esto, cambiar la clase cambia la dirección visual **sin rebuild ni JS** (`theme-contract.md` §2, `theme-architecture-audit.md` §7).

### Procedimiento completo (12 pasos)

El flujo operativo completo (inspección → token mapping → implementación → cobertura → variantes/estados → dark/light → responsive → reduced motion/forced colors → documentación → tests → smoke → aceptación) vive en **`docs/gelium-ui-theme-implementation-guide.md`** — NO se duplica aquí. Precondiciones de gate (STOP/BLOCKED si fallan): Phase 1 core cerrada, contrato aprobado, tests theme-agnósticos, clase data-driven, audit read-only de la referencia (`implementation-guide.md` §0).

## 5. Matriz theme × component × variant × state

> Estado: **ejecutada** — `TestThemeMatrixCoversEveryAvailableTheme` (web) recorre por glob todos los themes en disco (theme-material + theme-basecoat) y verifica por componente: familia de tokens en light, cobertura dark en la **ruta de clase única** (`.theme-dark`/`.dark`/`[data-theme="dark"]`, sin media query), y estados cubiertos por `var(--ui-*)`. `TestBasecoatTheme*` (web) pin las familias mandatorias del scope Phase I. La tabla conserva el scope del contrato (Phase I): Button, Text field, Card, Badge, Dialog, Toast, Data table.

| Componente | Variantes | Estados a cubrir | Tokens en theme-material | theme-basecoat |
|---|---|---|---|---|
| Button | primary/secondary/outline/text, link, loading | rest/hover/focus/pressed/disabled/loading | `--ui-button-*` + color semántico | ✅ |
| Text field | filled/outlined | rest/hover/focus/disabled/error | `--ui-field-*` | ✅ (floating label conservado; `--ui-size-field` 3rem, divergencia documentada) |
| Card | elevated/filled/outlined | rest/hover/focus/pressed | `--ui-card-*` (+ `--ui-card-fg` opcional) | ✅ |
| Badge | dot/count/large + tones error/success/warning/info | nunca color-only | `--ui-badge-*` | ✅ |
| Dialog | confirm/page | open/closed, light dismiss | `--ui-dialog-*` | ✅ |
| Toast | info/success/warning/error | transitorio auto-dismiss | `--ui-toast-*` | ✅ |
| Data table | sortable/selection/pagination | rest/hover/focus/pressed/selected/empty/error | `--ui-data-table-*` (scoped, NO en theme) | ✅ (anatomía scoped en `data-table.css`; theme pinta los colores semánticos) |

**Gaps cerrados antes de Phase I** (del audit, §3.5/§3.6):

- `--ui-color-surface-container`, `--ui-type-display-lg`, `--ui-type-title-md`, `--ui-font-mono` cerrados en core (`tokens.css:23,145,146,138`).
- `--ui-color-error` vs `danger`: alias temporal en `tokens.css:38`, canónico `danger`.
- State layers hardcodeados (`button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`) → tokenizados con `color-mix`.
- Dark: ambos themes declaran rutas clase + media consistentes (verificación §2).
- Tests theme-agnósticos: 13 tests parametrizados; matriz recorre themes por glob; sin aserciones de valor hex.

## 6. Divergencias documentadas

- Toast conserva `gelium:toast` + `aria-live` + fallback no-JS (`theme-contract.md:88`, decisión Phase I).
- Floating label de Text field es patrón Gelium; Basecoat lo conserva con `--ui-size-field: 3rem` (divergencia documentada en `themes/theme-basecoat/README.md` §5, no es variante de theme).
- Variantes que un theme nuevo aporte y Gelium no tenga (p. ej. `destructive` Button, pills de Badge) se resuelven como decisión del contrato: extender el core o documentar divergencia — nunca CSS de theme sobre markup distinto (`theme-contract.md` §7). En Basecoat quedan fuera de scope Phase I y documentadas.

---

**Definición de done (Phase J)**: registry refleja el estado real (dos themes implementados, cada uno con su README de divergencias), el procedimiento referencia los 2 pasos del mecanismo Phase H y el guía completo sin duplicarlo, y la matriz theme × component × variant × state se ejecuta por glob sobre todos los themes.
