# Basecoat UI — Audit de compatibilidad como theme Gelium UI

> Auditoría **read-only** de investigación de Basecoat UI como candidato a próximo theme (Phase 6 del sistema). No modifica código, templates, CSS ni tests. Única salida: este handoff.
> Relación con el roadmap: `docs/gelium-ui-system-roadmap.md` define el theme contract en Phase 5 (gate duro) y `theme-basecoat` en Phase 6, con scope Button · Text field · Card · Badge · Dialog · Toast · Data table, y el resultado esperado `Button Gelium { theme-material, theme-basecoat }` — NUNCA dos componentes separados. Este documento verifica si Basecoat puede ser una "dirección visual" codificada en tokens `--ui-*` sin imponer markup ni JS propios.

---

## 1. Resumen ejecutivo

- **Basecoat UI v1.0.2** (Ronan Berder, [basecoatui.com](https://basecoatui.com/), [github.com/hunvreus/basecoat](https://github.com/hunvreus/basecoat)) es "toda la magia de shadcn/ui sin React": una librería de componentes **CSS + tokens** para Tailwind CSS 4, con JS vanilla mínimo solo para los componentes interactivos. MIT, framework-agnóstico (HTML plano, Flask, Rails, Laravel, Django, Go), ~4.2k stars, 137 forks, 263 commits, npm `basecoat-css`.
- **No es** un framework CSS standalone ni Web Components: no hay Shadow DOM ni Custom Elements. La mayoría de componentes son CSS-only; de los 11 que piden JS, solo **Toast** está en el scope Gelium (y Gelium ya tiene su propio contrato Toast server-driven no-JS).
- **Modelo de tokens**: set **shadcn/ui-compatible** (~20 tokens semánticos: `--background, --foreground, --card, --popover, --primary, --secondary, --muted, --accent, --destructive, --border, --input, --ring, --radius` + fonts) en **oklch**, dark por **clase `.dark`** única. NO tokeniza elevation, spacing, motion ni typescale: viven como utilidades Tailwind hardcodeadas en el style pack.
- **El mapeo a Gelium es factible solo como traducción de dirección visual**: los ~20 tokens Basecoat deben EXPANDIRSE a los ~157 tokens `--ui-*` que los componentes Gelium consumen (shadow-0..5, motion-short, state-*, type-*, field-*, toast-*, dialog-*…). No hay valores "oficiales" que copiar para esas familias; hay que derivarlos del CSS del style pack elegido.
- **Veredicto**: compatible como **dirección visual** (theme que redefina `--ui-*`), **NO** como dependencia runtime ni como reemplazo de markup. No debe instalarse `basecoat-css`; el theme basecoat es `themes/theme-basecoat/theme.css`. Divergencias duras: anatomía de **Text field** (floating label vs label estático), modelo de **Toast** (toaster JS vs aria-live server-driven) y set de **variantes de Button/Badge** (falta destructive/outline pill en Gelium).
- **Alcance recomendado**: 5 de 7 componentes mapean limpio a tokens (Button, Card, Badge, Dialog, Data table); 2 requieren decisión explícita del theme contract (Text field, Toast). Estimar: la mayor parte del trabajo de Phase 6 no es copiar Basecoat, sino **cerrar el contrato token** que hoy está incompleto (ver `core-audit.md` y `theme-architecture-audit.md`).

---

## 2. Qué es Basecoat UI

### 2.1 Orígenes y propósito
- Posicionamiento explícito: **"All of the shadcn/ui magic, none of the React"** — lleva la ergonomía de shadcn/ui (clases cortas, design system, themes, componentes accesibles) a stacks sin React: "Flask, Rails, Laravel, Django, plain HTML".
- Idea rectora: usar Tailwind sin "class soup"; componentes con clases cortas (`btn`, `input`, `card`) y **HTML semántico**.
- Naturaleza: **librería de componentes CSS** (capas `@layer components` con `@apply` sobre Tailwind v4) + **token set** shadcn-compatible + **JS vanilla mínimo** para interactividad. NO es solo tokens ni solo componentes: es las dos cosas acopladas a Tailwind.

### 2.2 Licencia, mantenimiento, dependencias
- **Licencia**: MIT (`github.com/hunvreus/basecoat/blob/main/LICENSE.md`).
- **Mantenimiento**: proyecto activo de autor único con sponsors ("Built by Ronan Berder — Sponsor me"); release **v1.0.2** (sitio y npm), 4.2k stars, 19 issues / 15 PRs abiertos. Factor de bus alto (autor único), pero licencia MIT y código público → riesgo de adopción manejable para referencia visual.
- **Dependencias runtime**: ninguna obligatoria. Authoring requiere **Tailwind CSS 4**. Iconos **no bundled** (usa Lucide inline, igual que el contrato de iconos Gelium). Fonts: preferencias Geist Sans/Mono con fallback de sistema, no bundled.
- **Package**: npm `basecoat-css`, con CDN jsDelivr, imports npm y templates **Nunjucks/Jinja** (afín a `html/template` de Go).

### 2.3 Arquitectura CSS (evidencia directa del source)
Fuente: `src/css/` en GitHub. Separa **estructura** de **estilo** — conceptualmente la misma separación core/theme que Gelium busca:
- `src/css/base/base.css` → tokens (shadcn set) + `@theme` + reset/base.
- `src/css/components/*.css` → estructura/layout/accesibilidad de cada componente.
- `src/css/styles/*.css` → **style packs** Vega, Nova, Maia, Lyra, Mira, Luma, Sera, Rhea (los 8 estilos del registry shadcn actual). El **style pack ES la dirección visual** (color, radius, shadow, typography, spacing, variantes, estados).
- Entrypoints generados: `basecoat.css` (default = Vega), `basecoat-base.css` (sin style pack), `basecoat-<style>.css`, `basecoat-<style>.cdn.css`.

**Lectura para Gelium**: "Basecoat theme" no es un archivo, es elegir un style pack (p. ej. Vega default) y traducir sus valores al vocabulario `--ui-*`. El style pack "limpio" (sin pack) que ofrece Basecoat equivale al core Gelium agnóstico: tokens + estructura, sin dirección visual.

---

## 3. Markup y semántica

### 3.1 Patrón de clases y variantes (API v1.0)
- Clases cortas por componente: `btn`, `input`, `field`, `label`, `card`, `badge`, `table` / `table-container`, `dialog`, `toast` / `toaster`, `progress`, `alert`.
- Variantes por **atributo** `data-variant` (omisión = primary): `btn` → `primary|secondary|outline|ghost|link|destructive`; `badge` → `primary|secondary|outline|destructive|ghost|link`.
- Tamaños por **atributo** `data-size`: `xs|sm|default|lg|icon|icon-xs|icon-sm|icon-lg`.
- Hay capa `compat` para aliases pre-1.0 (clases de variante): solo migración.

### 3.2 Semántica HTML (coincide con Gelium)
- `btn` → `<button>` para acción, `<a>` para navegación (igual que Gelium).
- `dialog` → elemento nativo `<dialog>` con `aria-labelledby`/`aria-describedby`, estructura `header > h2 + p`, `section`, `footer`. Basecoat envuelve el contenido en un `<div>` interno; Gelium usa hijos directos `.ui-dialog-headline/.ui-dialog-content/.ui-dialog-actions`.
- `card` → `<div class="card">` con `header > h2/p`, `section`, `footer` (+ `card-action`, `data-size="sm"`). Gelium: `<article class="ui-card">` con `ui-card-title`/`ui-card-body` (artículo accionable como `<a>`/`<button>`).
- `field`/`input` → `<div role="group" class="field">` con `<label>`, control (auto-estilado por tipo), `<p>` helper; error vía `aria-invalid="true"`. **Sin floating label** (label estática encima).
- `table` → `<table class="table">` con `caption`, `thead`, `tbody`, `tfoot`; **puramente presentacional** (sin sort/pagination/selection).
- `badge` → `<span class="badge">` pill de texto (variantes, icon inline, link).
- ARIA: `aria-invalid`, `aria-expanded`, `aria-selected`, `aria-current="page"`, `role="status"` (toast), `sr-only`, `role="group"` en fields. Misma filosofía ARIA que Gelium.

### 3.3 Compatibilidad con HTML server-rendered
- **Sí, es el caso de uso primario** de Basecoat (stacks server-rendered sin SPA). Los templates Nunjucks/Jinja que shippea son la prueba de que el markup está pensado para server-render.
- Compatible con el markup semántico Gelium a nivel de **intención** (mismo HTML, misma ARIA). Pero los **nombres de clase son distintos** (`ui-button` vs `btn`): para Gelium esto es irrelevante porque el contrato exige markup Gelium único y themes solo definen tokens — el theme NO renombra clases.

---

## 4. Sistema de tokens

### 4.1 El set de tokens (extraído de `src/css/base/base.css`)
```css
:root {
  --radius: 0.625rem;
  --background, --foreground, --card, --card-foreground, --popover, --popover-foreground,
  --primary, --primary-foreground, --secondary, --secondary-foreground,
  --muted, --muted-foreground, --accent, --accent-foreground,
  --destructive, --border, --input, --ring;
  --chart-1..5, --sidebar*, --scrollbar-track/thumb/width/radius;
  --chevron-down-icon, --chevron-down-icon-50, --check-icon;  /* data-URI SVG */
}
.dark { /* mismo set, valores oscuros */ }
@theme {
  --font-sans: "Geist Sans", …;  --font-mono: "Geist Mono", …;
  --radius-sm/md/lg/xl: calc(var(--radius) ± …);   /* escala derivada */
  --color-*: var(--*) ... /* mapeo a utilidades Tailwind */
}
```
- **Colores**: **oklch** (p. ej. `--primary: oklch(0.205 0 0)`), no hex. Dark mode: **un solo bloque `.dark`** con el set completo (vía `@custom-variant dark (&:is(html.dark *))`). No usa media query automática ni `light-dark()`.
- **Radius**: escala derivada de un `--radius` único (`sm=calc-4px, md=calc-2px, lg=var, xl=calc+4px`).
- **Lo que NO tokeniza** (crítico): elevation, spacing, motion y typescale no existen como variables; son utilidades Tailwind hardcodeadas en los style packs (`shadow-xs/md/lg/xl`, `h-9`, `px-2.5`, `text-sm`, `duration-100`, `transition-[color,box-shadow]`). Tampoco hay opacidades de estado (`hover:bg-primary/80`, `dark:hover:bg-muted/50`).

### 4.2 Mapeo propuesto a `--ui-*` (viable a nivel de color semántico)

| Basecoat / shadcn | Gelium `--ui-*` | Nota |
|---|---|---|
| `--background` | `--ui-color-canvas` | |
| `--foreground` | `--ui-color-fg` | |
| `--card` / `--card-foreground` | `--ui-card-container-*` (+ faltaría `--ui-card-fg`) | |
| `--primary` / `--primary-foreground` | `--ui-color-primary` / `--ui-color-primary-fg` | 1:1 |
| `--secondary` / `--secondary-foreground` | `--ui-color-secondary` / `--ui-color-secondary-fg` | 1:1 |
| `--muted` / `--muted-foreground` | sin equivalente directo (≈ `--ui-color-fg-muted`, `surface-container`) | nuevo token o reuso |
| `--accent` / `--accent-foreground` | sin equivalente | no usado por scope inicial |
| `--destructive` / `--destructive-foreground` | `--ui-color-danger` / `--ui-color-danger-fg` | renaming `error`/`danger` ya pendiente |
| `--border` | `--ui-color-border` | |
| `--input` | `--ui-field-border` / `--ui-select-outline` | |
| `--ring` | `--ui-color-focus-ring` | 1:1 |
| `--radius` (escala) | `--ui-radius-*` | derivar 4 niveles |
| `--font-sans` / `--font-mono` | `--ui-font-sans` / `--ui-font-mono` | **`--ui-font-mono` hoy indefinido en Gelium** (core-audit) |

### 4.3 Gaps del mapeo (qué inventar, no copiar)
- Familias Gelium **sin origen en Basecoat**: `--ui-shadow-0..5` (elevation), `--ui-state-*` (opacidades), `--ui-motion-short`/`--ui-easing-standard`, `--ui-type-*` (typescale Material), `--ui-field-container`, `--ui-dialog-scrim`, `--ui-toast-*`. Sus valores deben derivarse del style pack (sombras/duraciones visibles en el CSS `@apply`), no de tokens públicos.
- **Conversión de color**: Gelium es hex; Basecoat es oklch → convertir valores del style pack elegido a hex (o decidir soportar oklch en `--ui-*`).
- Dirección contraria: tokens Basecoat sin destino en el scope (`--popover-*`, `--accent-*`, `--chart-*`, `--sidebar-*`, icon data-URIs) no se usan en Phase 6.

---

## 5. Estados y variantes

### 5.1 Modelo de estados
- **Pseudo-clases nativas**: `:hover`, `:active`, `:focus-visible`, `:disabled`, `:checked`, `:focus-within`, `:dir(rtl)`. Más selectores avanzados Tailwind: `has-[[data-icon='inline-start']]`, `not-aria-[haspopup]`, `not-last:`, `aria-expanded:`.
- **Atributos ARIA como selectores de estado**: `aria-invalid` (error), `aria-expanded` (menús), `aria-selected` (tabs), `aria-current` (breadcrumb), `data-variant`/`data-size` (variantes).
- **Patrón visual de estados** (en `vega.css`): `hover:bg-primary/80`, `hover:bg-muted`, `focus-visible:ring-ring/50 focus-visible:ring-[3px] focus-visible:border-ring`, `aria-invalid:border-destructive aria-invalid:ring-destructive/20`, `disabled` vía `:disabled` + `opacity-50` (en `in-data-[disabled=true]:opacity-50`), `active:not-aria-[haspopup]:translate-y-px`.
- **NO hay state-layer MD3** (overlay inset con opacidad de token `--ui-state-*`). El feedback es cambio de color/background + ring + micro-movimiento.

### 5.2 Lectura para Gelium
- Gelium **modela estados en el CSS de componentes** (no en el theme), igual que Basecoat → no bloquea el multi-theme.
- Divergencia de mecanismo: los state-layers de Gelium (`--ui-state-hover-opacity` etc.) requieren que el theme basecoat provea opacidades que Basecoat no define (¿0.08? ¿otro?). La estética Basecoat de hover es cambio de color del propio token (`bg-primary/80`), no overlay → traducible a los tokens de color, no a los de opacity.

---

## 6. JavaScript y no-JS end-to-end

### 6.1 Requerimiento de JS por componente
- **CSS-only**: la gran mayoría (Button, Input/Field, Textarea, Native Select, Card, Badge, Table, Pagination, Progress, Skeleton, Alert, Kbd, Label, …).
- **JS vanilla necesario** (11): Accordion, Combobox, Command, Drawer, Dropdown Menu, Popover, Select, Sidebar, Slider, Tabs, Toast.
- **De los 7 del scope Gelium, SOLO Toast requiere JS en Basecoat.**

### 6.2 Tecnología JS
- Vanilla JS auto-init (`window.basecoat.initAll()`), con `force: true` para restore de HTMX history — **diseñado para convivir con HTMX**.
- **No** Web Components, **no** Shadow DOM, **no** Custom Elements obligatorios, **no** React/Radix.
- Dialog usa el `<dialog>` nativo + `showModal()`: el componente no lleva "JS de componente", pero el trigger de los ejemplos es `onclick="...showModal()"`. Gelium ya lo hace **mejor** (más no-JS): `command`/`commandfor` + `closedby="any"` (README.md Dialog). Basecoat no tiene equivalente declarativo documentado.
- Toast: el `toaster` necesita JS para la API (`toaster.toast(config)`, `toast.close()`, auto-dismiss) aunque soporta markup server-rendered dentro del toaster.

### 6.3 Compatibilidad con no-JS Gelium
- **6 de 7 componentes del scope** funcionan sin JS en ambos (Button, Card, Badge, Dialog vía comandos declarativos Gelium, Table, Input).
- **Toast: NO compatible con el modelo Basecoat sin JS.** Gelium conserva su propio contrato: `aria-live` + `HX-Trigger: {"loom:toast":…}` + fallback inline sin JS (README.md, `web/templates/toast.html`). Basecoat aporta solo la **estética** (superficie popover, `data-category`, icon, title/description/action).

---

## 7. Compatibilidad con la arquitectura Gelium

### 7.1 Veredicto
**COMPATIBLE como dirección visual** (theme que redefina `--ui-*`), **INCOMPATIBLE como adopción directa** (importar `basecoat-css`, copiar su CSS de componentes o su markup). El roadmap ya lo presupone: "No es una dependencia de CDN ni un port mecánico de Material Web, Basecoat ni shadcn" y "No reemplazar el markup Gelium por markup Basecoat completamente distinto". Basecoat valida el modelo (tokens + dark class + separación estructura/estilo) y aporta una dirección visual concreta y moderna; el theme basecoat es una **traducción del style pack elegido al vocabulario `--ui-*`**.

### 7.2 Incompatibilidades concretas (top)
1. **Vocabulario de tokens incompleto vs el contrato Gelium**: Basecoat no define elevation, spacing, motion, typescale ni opacidades de estado → `--ui-shadow-*`, `--ui-motion-*`, `--ui-type-*`, `--ui-state-*` no tienen fuente; hay que derivarlos del CSS del style pack (valores hardcodeados en `@apply`) o inventar la escala. Además Gelium referencia tokens aún indefinidos (`--ui-color-surface-container`, `--ui-type-display-lg`, `--ui-type-title-md`, `--ui-color-error`, `--ui-font-mono` — `core-audit.md`): **el contrato debe cerrar esos gaps ANTES de Phase 6**, porque Basecoat no los provee.
2. **Anatomía de Text field**: Gelium usa floating label (filled/outlined, `ui-text-field-control` + `:placeholder-shown`), Basecoat usa label estática + input plano (`field`/`input`). No es re-skin: es un patrón distinto. El markup Gelium no cambia → decisión del theme contract (mantener floating label con estética Basecoat, o documentar divergencia).
3. **Modelo de Toast**: Basecoat requiere JS (`toaster` + auto-dismiss); Gelium exige no-JS end-to-end (aria-live + fallback inline). El theme no puede hacer que el Toast Gelium dependa del runtime Basecoat.
4. **Variantes de Button/Badge**: Gelium tiene primary/secondary/outline/text (Button) y dot/count/label (Badge); Basecoat añade destructive/ghost/link (btn) y pill con data-variant (badge). Destructive NO existe en el Button Gelium (`--ui-color-danger` solo alimenta badge/toast). Faltan tokens de superficie para pills (badge) y para destructive.
5. **Dark mode**: Basecoat usa `.dark` clase única (oklch); Gelium duplica dark en clase + media query con drift (`theme-architecture-audit.md:93`). Basecoat no resuelve la media query automática; el theme contract debe unificar y Basecoat valida el patrón "clase única".
6. **Espacio de color**: oklch vs hex (conversión o decisión de soporte).
7. **Card**: Gelium 3 variantes de elevación (elevated/filled/outlined); Basecoat 1 base + `data-size=sm`. Mapeo: basecoat ≈ outlined/elevated; filled sin equivalente directo (→ `bg-muted`).
8. **Sin dependencia runtime**: instalar `basecoat-css` para "temar" violaría la regla "sin dependencias runtime innecesarias" y el contrato token-only. Es solo referencia.

### 7.3 Lo que SÍ se alinea
- Tokens semánticos por CSS variables (shadcn-compatible) sobre `:root` → mismo mecanismo cascade-time que Gelium.
- Dark por clase en `<html>` → compatible con la clase `theme-*` de `layout.html` (`theme-basecoat.dark` o `[data-theme=dark]`).
- Separación estructura/estilo (base vs style pack) → análoga a la separación core/theme que exige Phase 1.
- Framework-agnostic, HTML server-rendered first, sin Shadow DOM/Custom Elements → cumple todas las restricciones duras.
- Nunjucks/Jinja templates → confirma el modelo de templates server-side (Go `html/template`).

---

## 8. Alcance recomendado (Phase 6 — 7 componentes)

### 8.1 Mapeo limpio (solo tokens; markup Gelium intacto)

| Componente | Veredicto | Notas y tokens a cubrir |
|---|---|---|
| **Button** | ✅ Alto | `--ui-color-*`, `--ui-radius-*`, `--ui-shadow-*`, `--ui-focus-*`. Alinear variantes: primary/secondary/outline mapean; ghost/link/destructive de Basecoat requieren decisión (¿extender variantes Gelium o ignorarlas?). State-layer: Gelium ya usa `color-mix` theme-aware (patrón a completar). |
| **Dialog** | ✅ Alto | Mismo `<dialog>`; `--ui-dialog-container/fg/body/scrim`, `--ui-radius-*`, focus, motion. Conservar trigger declarativo Gelium (`command`/`commandfor`), superior al `onclick` de Basecoat. |
| **Card** | 🟡 Medio-alto | `--ui-card-*` (container-elevated/filled/outlined, outline-color, radius). Basecoat ≈ elevated/outlined con `ring-1`+`shadow-xs`; filled → surface alternativa. `data-size=sm` de Basecoat no aplica (Gelium no tiene size en card). |
| **Badge** | 🟡 Medio | Dot/count es **Gelium-only** (Material); pill de texto de Basecoat mapea a `--ui-badge-container/fg` + variante "large". Decidir si el pill pill-identidad Basecoat (outline/ghost) se agrega como variantes nuevas del badge Gelium. |
| **Data table** | ✅ Alto | Basecoat pinta la tabla presentacional (caption/thead/tbody/tfoot, hover rows, borders) → `--ui-data-table-*` (header-height, row-height, radius, borders, checkbox). **Sort/pagination/selection siguen siendo contratos server-side Gelium** (`data-table.html`), no vienen de Basecoat. |

### 8.2 Requieren decisión del theme contract

| Componente | Divergencia | Decisión |
|---|---|---|
| **Text field** | floating label (Gelium) vs label estático (Basecoat) | Mantener floating label con estética Basecoat (recomendado: menos ruptura de contrato/ARIA), o aprobar variante "flat label" como divergencia documentada (Phase 7 los detectaría). |
| **Toast** | toaster JS (Basecoat) vs aria-live no-JS (Gelium) | Conservar el contrato `loom:toast` + fallback no-JS; tomar solo la estética (superficie popover, `data-category` info/success/warning/error, icon, title/desc/action). |

### 8.3 Orden recomendado dentro de Phase 6
1. Cerrar el **theme contract** (Phase 5) y los gaps de tokens (`--ui-color-surface-container`, `--ui-type-*`, `--ui-font-mono`, unificar `danger`/`error`).
2. Elegir el **style pack** de referencia (default = **Vega**; alternativas Nova/Maia/Lyra/Mira/Luma/Sera/Rhea) y convertir sus valores a hex para `--ui-*`.
3. Implementar por orden de fidelidad: Button → Card → Dialog → Data table → Badge → Text field → Toast.
4. Correr la matriz `theme × component × variant × state` (Phase 7) para detectar desvíos de la dirección visual Basecoat.

---

## 9. Implicaciones para el theme contract (Phase 5)

- El contrato debe exigir por theme: definición de **todas** las familias `--ui-*` referenciadas por los componentes (no solo las de color), dark por clase única sin duplicación, y tests theme-agnósticos (presencia de token, no valor hex — `theme-architecture-audit.md:184`). Basecoat es un caso de prueba perfecto: su set corto obliga a decidir qué es core vs theme.
- **No** crear archivos `themes/theme-basecoat/<componente>.css`: el contrato es token-only; los componentes son únicos (`theme-architecture-audit.md:219`).
- Mecanismo: bundle-all + selección runtime por clase `theme-basecoat` en `<html>`; no `THEME=` build var como base (`theme-architecture-audit.md:220`).
- Basecoat resuelve dos deudas conocidas de Gelium por diseño: dark single-block y tokens semánticos limpios → sirve de guía para la unificación de dark en Phase 1/5.

---

## 10. Evidencia (URLs consultadas)

- Sitio y docs: https://basecoatui.com/ · https://basecoatui.com/introduction/ · https://basecoatui.com/installation/ · https://basecoatui.com/customization/ · https://basecoatui.com/templates/
- Componentes: https://basecoatui.com/components/button/ · /input/ · /field/ · /card/ · /badge/ · /dialog/ · /table/ · /toast/ · /alert-dialog/
- Repo: https://github.com/hunvreus/basecoat · MIT https://github.com/hunvreus/basecoat/blob/main/LICENSE.md
- Source tokens (extraído): https://raw.githubusercontent.com/hunvreus/basecoat/main/src/css/base/base.css
- Source style pack default (Vega): https://raw.githubusercontent.com/hunvreus/basecoat/main/src/css/styles/vega.css
- shadcn/ui theming (token set de referencia): https://ui.shadcn.com/docs/theming
- Repo Gelium (evidencia local): `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-system-roadmap.md`, `docs/handoffs/core-audit.md`, `docs/handoffs/theme-architecture-audit.md`, `themes/theme-material/theme.css`, `web/templates/{button,card,badge,dialog,toast,data-table,text-field}.html`, `web/styles/app.css`.

---

## Anexo — Mapeo de tokens rápido (Basecoat/shadcn → Gelium `--ui-*`)

```text
--background        → --ui-color-canvas
--foreground        → --ui-color-fg
--muted-foreground  → --ui-color-fg-muted
--primary           → --ui-color-primary
--primary-foreground→ --ui-color-primary-fg
--secondary         → --ui-color-secondary
--secondary-foreground → --ui-color-secondary-fg
--destructive       → --ui-color-danger
--destructive-foreground → --ui-color-danger-fg
--border            → --ui-color-border
--input             → --ui-field-border, --ui-select-outline
--ring              → --ui-color-focus-ring
--card / --card-foreground → --ui-card-container-* (+ --ui-card-fg nuevo)
--radius            → escala --ui-radius-sm/md/lg/xl
--font-sans         → --ui-font-sans
--font-mono         → --ui-font-mono (hoy indefinido en Gelium)
--popover / --accent / --chart-* / --sidebar-* → fuera de scope Phase 6
Elevation/Spacing/Motion/Typescale/State-opacity → SIN token público en Basecoat;
  derivar valores del style pack elegido para --ui-shadow-*, --ui-state-*,
  --ui-motion-*, --ui-type-* y geometría de componentes.
```
