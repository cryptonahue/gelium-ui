# Gelium UI — Onboarding de UIs externas (clasificación)

> **Paso 0 del flow de integración.** Este documento formaliza cómo clasificar una
> UI externa (Base UI, Material, Basecoat, shadcn, Protocol, lo que sea) ANTES de
> tocar cualquier archivo. Complementa `gelium-ui-theme-implementation-guide.md`
> (procedimiento de theme, 12 pasos) y `gelium-ui-theme-contract.md` (el contrato).
>
> Regla de oro: **primero se clasifica, después se audita, recién después se
> implementa.** Ninguna UI externa entra al repositorio sin pasar por esta
> clasificación.

---

## 1. La pregunta central

Toda UI externa que quieras traer cae en UNA de estas capas del sistema:

| Capa | ¿Qué es? | ¿Entra a Gelium como? |
|---|---|---|
| **Theme** | Dirección visual (paleta, shape, type, motion) | `themes/<theme>/theme.css` — valores `--ui-*` |
| **Vocabulario** | Nombres canónicos, anatomía, estados de patrones | Doc de vocabulario + propuesta de contrato |
| **Patrón/composición** | Reglas de elección y combinación | `composition-rules.md` / pattern registry |
| **Componente** | Implementación concreta de un patrón | Capa components + propuesta de contrato |
| **Mecanismo** | Comportamiento transversal (dark, focus, state layers) | Core — con consumidor real |
| **Runtime** | JS/React/Web Components como requisito | **DESCARTADO** — viola HTML-first |

**Clasificar es decidir en qué capa cae cada pieza de la UI externa, NO copiarla
entera.** Una UI puede aportar a varias capas a la vez (Basecoat aportó theme;
Mozilla Protocol aportó vocabulario/patrones; Base UI aportaría solo vocabulario).

---

## 2. Árbol de decisión (Paso 0)

```text
¿La UI externa trae dirección visual (paleta, shape, type)?
├── SÍ → ¿Es CSS/estilos que se pueden expresar como valores --ui-*?
│        ├── SÍ → THEME → seguir gelium-ui-theme-implementation-guide.md
│        └── NO → ¿Es un runtime (React/JS obligatorio)?
│                 ├── SÍ → DESCARTAR como theme; evaluar como vocabulario
│                 └── NO → AUDITAR: puede ser vocabulario/patrón
└── NO → ¿Trae patrones, anatomía, estados, comportamiento?
         ├── SÍ → VOCABULARIO/PATRÓN → auditar contra el vocabulario actual
         └── NO → NADA QUE APORTAR → documentar y descartar
```

Reglas duras:

1. **Runtime incompatible (React, Lit, Web Components obligatorios) → nunca como
   runtime.** El core de Gelium es HTML-first, server-rendered, 0 JS requerido.
   La UI externa puede ser *referencia de comportamiento* (focus management,
   estados, keyboard nav) pero su implementación NO se porta.
2. **Prefijos ajenos (`m3-*`, `mzp-*`, clases de framework) → nunca entran.** Se
   traduce al vocabulario `--ui-*` / `ui-*` o se descarta.
3. **Licencia y atribución se auditan en el Paso 0** (ver §5).
4. **Si la clasificación no es obvia → AUDIT read-only primero**
   (`docs/handoffs/<referencia>-audit.md`), sin tocar código.

---

## 3. ¿Qué puede convertirse en core?

El core es el contrato estructural. Algo de una UI externa entra al core SOLO si
cumple TODAS estas condiciones (del core doc §3 y §5.0):

| Criterio | Pregunta | Ejemplo que entra | Ejemplo que NO entra |
|---|---|---|---|
| **Rol o escala estructural** | ¿Es un rol semántico o una escala, no un valor estético? | `--ui-state-hover-opacity`, `--ui-focus-thickness` | `#141218` (color de canvas Material) |
| **Múltiples consumidores** | ¿Lo consumen 2+ componentes? | `--ui-space-*`, `--ui-shadow-0..5` | Token usado por un solo archivo (→ internal) |
| **Consumidor real hoy** | ¿Alguien lo usa AHORA? | `--ui-color-danger` (todos los state layers) | `--ui-radius-xl` (sin consumidor → dead) |
| **Mecanismo, no personalidad** | ¿Es comportamiento transversal? | Dark unificado, state layers `color-mix()`, theme identity | Radio de esquina 28px (personalidad del theme) |

**Lo que NUNCA entra al core:**

- Valores de paleta (hex/oklch) de ninguna UI — es dirección visual del theme
- Typescale, densidad, shape por defecto — personalidad del theme
- Implementaciones de componentes de una UI específica — capa components
- Criterios de elección de patrones — capa composition rules
- JavaScript de componente — solo enhancement justificado y auditado platform-first
- Prefijos ajenos (`m3-*`, `mzp-*`) — contaminación del contrato

**Cómo proponer un token core nuevo:** si el mapeo de una UI necesita un token
que el core no tiene (`--ui-card-fg` fue el caso Basecoat), se propone al core CON
contrato y consumidor real — nunca se inventa en el theme.

---

## 4. Qué puede convertirse en theme

Un theme es dirección visual. De una UI externa entran al theme:

1. **Valores de paleta** traducidos a los roles canónicos (`--ui-color-*`)
2. **Shape** (radios, borders) sobre la escala core
3. **Typography** sobre la composición core (size/weight/line-height/letter-spacing)
4. **Elevation / shadows** sobre `--ui-shadow-*`
5. **Motion** (duraciones, easings) sobre `--ui-motion-*` / `--ui-easing-*`
6. **State opacities** sobre `--ui-state-*`

Reglas (de `theme-contract.md`):

- El theme define **valores** para tokens del contrato; NO markup, NO JS.
- El theme solo toca los **grados de libertad** que el contrato declara
  (color, shape, type, motion) — nunca anatomía de componente.
- Un theme nuevo debe pasar la **misma suite theme-agnóstica** que Material
  (presencia de familias, cobertura light/dark) — ver
  `gelium-ui-theme-verification.md`.
- Divergencias de anatomía (ej. floating label de Text field) se documentan en
  el README del theme, nunca como CSS de theme sobre markup distinto.

---

## 5. Checklist del Paso 0 (una UI nueva)

- [ ] Clasificada la UI en al menos una capa (theme / vocabulario / patrón /
      componente / mecanismo / descartada)
- [ ] Si es runtime obligatorio → descartada como runtime; documentado qué
      aportaría como referencia
- [ ] Licencia y atribución auditadas y compatibles
- [ ] Colisiones de naming detectadas contra el vocabulario Gelium
      (el caso "Callout" de Protocol: resolvidas ANTES de implementar)
- [ ] Si propone tokens core → listados con consumidor real propuesto
- [ ] Si es theme → precondiciones de `theme-implementation-guide.md` §0
      verificadas (STOP/BLOCKED si fallan)
- [ ] Decisión registrada en `gelium-ui-theme-registry.md` (o descarte
      documentado)

---

## 6. Ejemplos clasificados

| UI | Clasificación | Qué aportó / aportaría | Estado |
|---|---|---|---|
| **Material 3 / Material Web** | Referencia de vocabulario + theme | Roles de color, estados, anatomía; theme `theme-material` | Integrado |
| **Basecoat UI** | Theme | Style pack traducido a `--ui-*` (theme `theme-basecoat`, Phase I) | Integrado |
| **Mozilla Protocol** | Vocabulario + patrones public-content | 14 patrones como referencia, NO runtime; colisión "Callout" resuelta | Audited |
| **Base UI (base-ui.com)** | Vocabulario/referencia (headless primitives) | Focus management, estados, keyboard nav como referencia; **NO runtime** (React) | Pendiente de audit |
| **shadcn/ui** | Futuro preset visual "shadcn-like" | Dirección visual posible; no port directo | Roadmap |
| **Cualquier runtime React/Lit** | Descartada como runtime | Documentar qué referencia aporta, no portar | — |

---

**Definición de done**: toda UI externa que se considere entra primero por este
doc (clasificación registrada), y el resultado alimenta el audit
(`docs/handoffs/<referencia>-audit.md`) y, si aplica, el procedimiento de theme
o la propuesta de contrato al core.
