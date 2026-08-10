# Gelium UI — Theme Implementation Guide

> Guía operativa para implementar un theme en Gelium UI (antes Gelium UI).
> Fase 5 del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Complementa `docs/gelium-ui-theme-contract.md` (el contrato) con el procedimiento exacto paso a paso.
> Base: `docs/handoffs/theme-architecture-audit.md`, `docs/handoffs/core-audit.md`, `docs/handoffs/basecoat-audit.md`.

---

## 0. Precondiciones

Antes de tocar cualquier archivo de theme:

- [ ] Phase 1 (core) cerrada: gaps de tokens cerrados, naming unificado (`danger`), state layers theme-aware, dark unificado, `--ui-font-mono` definido.
- [ ] Phase 5 (theme contract) aprobada.
- [ ] Suite de tests refactorizada a theme-agnóstica (presencia, no valor).
- [ ] Clase del theme en `layout.html` data-driven.
- [ ] Audit read-only de la referencia terminada (`docs/handoffs/<referencia>-audit.md`).

Si alguna precondición falla: STOP, reportar `BLOCKED`, no implementar.

---

## 1. Inspección

1. Leer la referencia visual (docs oficiales, source CSS, style packs). Identificar:
   - token set y valores (anotar espacio de color: hex/oklch);
   - variantes y estados por componente;
   - requisitos JS vs CSS-only;
   - anatomía que difiere del markup Gelium;
   - licencia y atribución.
2. Inventariar los componentes Gelium del scope del theme y sus tokens (`var(--ui-*)` consumidos).
3. Registrar el baseline: `git status`, commit actual, tests verdes ANTES de empezar.

**Resultado**: `docs/handoffs/<referencia>-audit.md` + lista de tokens a cubrir por componente.

---

## 2. Token mapping

1. Construir la tabla de mapeo referencia → `--ui-*` (template: `basecoat-audit.md` Anexo).
2. Decidir el espacio de color del theme: **hex por defecto** (convención Gelium actual); aprobar oklch explícitamente si se adopta.
3. Para cada familia Gelium sin origen directo en la referencia (elevation, motion, typescale, spacing, state-opacity), **derivar valores del CSS de la referencia** (sombras, duraciones, tipos, espaciados visibles en `@apply`) — documentar de dónde salió cada uno.
4. Marcar tokens sin destino (p. ej. `--popover-*`, `--accent-*`, `--chart-*` de Basecoat) como fuera de scope.
5. Si el mapeo necesita tokens nuevos (`--ui-card-fg`), proponerlos al core con contrato — no inventar en el theme.

**Resultado**: tabla de mapeo completa y aprobada.

---

## 3. Implementación

### 3.1 Estructura

```text
themes/<theme>/
├── theme.css       # tokens --ui-* autocontenido (luz + dark)
└── README.md       # dirección visual, mapping, divergencias, matriz
```

### 3.2 theme.css

1. Selector raíz: `.theme-<name>` (canónico) + aliases (`[data-theme="<name>"]` si aplica).
2. Bloque luz: definir TODAS las familias transversales + por componente del scope.
3. Bloque dark: UNA rutina (clase única o `light-dark()` según decisión del contrato); sin duplicación textual.
4. `color-scheme` consistente con el core (una sola fuente de verdad).
5. Respetar naming canónico `--ui-*`; sin prefijos de terceros.

### 3.3 Integración

1. Agregar el import explícito en `web/styles/app.css` junto a los otros themes.
2. Activar la clase en runtime (template/server data-driven): `class="theme-<name>"` o `data-theme="<name>"`.
3. Verificar que `npm run build` produce UN solo `web/static/app.css` con ambos themes inlineados.

---

## 4. Cobertura de componentes

1. Por cada componente del scope, verificar que TODOS sus `var(--ui-*)` están definidos en el theme.
2. Ejecutar la matriz: componente × variante × estado (rest/hover/focus/pressed/disabled/selected/error/loading) en browser.
3. Componentes con anatomía divergente de la referencia:
   - si es solo estética → cubrir con tokens;
   - si es patrón distinto (p. ej. Text field floating label vs label estático) → **decisión del theme contract**: mantener patrón Gelium con estética del theme, o aprobar variante como divergencia documentada. Nunca cambiar el markup Gelium para acomodar la referencia.

---

## 5. Variantes y estados

1. Verificar que las variantes existentes del componente (p. ej. Button primary/secondary/outline/text) se ven coherentes.
2. Verificar estados: hover/focus/pressed (state layers theme-aware con `color-mix`, sin `rgb()` fijo), disabled (opacidad del token), selected (`:checked`, `aria-current`, `aria-sort`), error (`aria-invalid`).
3. Si la referencia aporta variantes que Gelium no tiene (p. ej. `destructive` en Basecoat), decidir: extender el componente en el core (con contrato y tests) o documentar como fuera de scope del theme.

---

## 6. Dark / light

1. Verificar light completo.
2. Verificar dark completo (ambas rutas si el contrato lo permite: clase y preferencia del sistema).
3. Verificar que no hay drift entre rutas (el bug `--ui-switch-track-unselected` de theme-material es la regresión a evitar).
4. Contrast AA en ambos modos.

---

## 7. Responsive

1. Verificar narrow/wide con layouts fluidos (grid `auto-fit/minmax`, `min()/clamp()`).
2. Sin breakpoints nuevos salvo necesidad real (tokens `--ui-breakpoint-*`).
3. Overlays fluidos (`calc(100vw - n)`) sin regresión.

---

## 8. Reduced motion y forced colors

1. `prefers-reduced-motion`: verificar que las excepciones mínimas del core se mantienen y el theme no reintroduce animaciones no reducibles.
2. `forced-colors`: verificar que ningún estado depende solo del color; controles nativos y `currentColor`/`CanvasText` sobreviven.

---

## 9. Documentación

Escribir `themes/<theme>/README.md`:

- dirección visual y referencia (con URLs y licencia);
- style pack elegido (si la referencia tiene varios);
- tabla de token mapping;
- luz/dark;
- componentes del scope y su cobertura;
- divergencias: componente × divergencia × decisión;
- notas de conversión de color.

---

## 10. Tests

1. Suite theme-agnóstica pasando (cobertura de tokens, presencia no valor, dark ambas rutas, matriz).
2. Si se agregaron tokens nuevos al core: tests de contrato actualizados.
3. `styles_contract_test.go` en sync con `app.css`.

```bash
npm run build
go test ./...
go vet ./...
```

---

## 11. Smoke

En el puerto permitido por el entorno, browser real:

- consola sin errores;
- keyboard (tab, enter, escape, arrows donde aplique);
- no-JS (flujo principal completo sin JS);
- HTMX (enhancement: validación 422, toast `loom:toast`, refresh data table);
- light/dark;
- narrow/wide;
- reduced motion;
- forced colors;
- estados empty/loading/error.

---

## 12. Aceptación visual

- Checklist manual revisada por el mantenedor antes de declarar el theme terminado.
- Matriz theme × component × variant × state ejecutada (Phase 7 si hay 2+ themes).
- Evidencia guardada: capturas o URLs de cada estado.

---

## Checklist final de un theme

- [ ] Precondiciones §0 cumplidas.
- [ ] Audit §1 entregado.
- [ ] Token mapping §2 completo y aprobado.
- [ ] `themes/<theme>/theme.css` autocontenido, luz + dark sin duplicación.
- [ ] Import en `app.css` + selección runtime.
- [ ] Cobertura §4 de todos los componentes del scope.
- [ ] Variantes y estados §5 verificados.
- [ ] Light/dark §6, responsive §7, reduced motion/forced colors §8 verificados.
- [ ] `themes/<theme>/README.md` §9 publicado.
- [ ] Tests §10 verdes (build + go test + go vet).
- [ ] Smoke §11 completo.
- [ ] Aceptación visual §12 aprobada.
- [ ] Sin cambios en markup/clases/contratos del core (o cambios de core con su propio contrato).

---

## Cierre

Reportar con el formato del roadmap:

```text
STATUS: READY_FOR_INTEGRATION
FILES CREATED: themes/<theme>/theme.css, themes/<theme>/README.md
FILES MODIFIED: web/styles/app.css, layout (clase data-driven)
DECISIONS: …
TESTS: npm run build; go test ./...; go vet ./...
BUILD: OK/FAIL
SMOKE: …
EVIDENCE: …
OPEN QUESTIONS: …
RISKS: …
NEXT PHASE: Phase 7 cross-theme validation
```
