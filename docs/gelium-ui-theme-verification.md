# Gelium UI — Theme Verification

> Guía de verificación theme-agnóstica (Phase C del system roadmap).
> Complementa `docs/gelium-ui-theme-contract.md` (el contrato) y `docs/gelium-ui-theme-implementation-guide.md` (el procedimiento).
> Base: `docs/handoffs/theme-architecture-audit.md`.
> Gate: **un theme nuevo solo entra cuando la suite entera es theme-agnóstica y pasa sin asumir valores Material.**

---

## 1. Qué se verifica y por qué

La suite verifica **contratos de presencia**, nunca valores. El objetivo de Phase C es que un theme nuevo
(theme-basecoat, Phase I) pase la suite **sin editar tests**: la suite descubre themes por glob
(`themes/*/theme.css`) y asevera que cada uno defina las familias de tokens que el sistema consume.

El principio es: el theme modela **tokens**, no estados ni variantes. Los estados (hover/focus/pressed/
disabled/selected/error/loading/empty) son parte del CSS de componentes y se verifican con tokens, no con
literales. Un theme nuevo cubre el contraste light/dark redefiniendo los tokens de color; los componentes
no cambian.

## 2. Cómo correr la matriz

```bash
# 1. Build del asset compilado (los tests de asset compilado lo requieren)
npm run build

# 2. Suite completa
go test ./...

# 3. Matriz formal theme × component × scheme + component × state
go test ./web/ -run 'TestThemeMatrix' -v

# 4. Contratos CSS (tokens, light/dark, forced-colors, sourceAppCSS sync)
go test ./web/ -run 'TestMaterialTheme|TestSourceAppCSS|TestNoColorLiterals' -v

# 5. Estática
go vet ./...
```

La matriz vive en `web/styles_theme_matrix_test.go`:

- `TestThemeMatrixCoversEveryAvailableTheme` — itera `availableThemes(t)` (glob de `themes/*/theme.css`,
  nunca una ruta hardcodeada). Por cada theme y componente verifica, por presencia (no valor):
  - la familia `--ui-<componente>-*` está definida en el bloque **light**;
  - cobertura dark en las **dos rutas** (clase `.theme-dark` + `@media (prefers-color-scheme: dark)`):
    - **directa** (field, dialog, toast, card, switch, slider, progress, fab, select y la familia semántica
      de color): al menos un token de la familia redefinido en cada ruta dark;
    - **derivada** (badge, checkbox, radio, divider): la familia vive en light y referencia `--ui-color-*`
      semánticos; la legibilidad dark viene de que ambos bloques dark redefinen esos colores.
  - el componente cubre sus estados documentados con `var(--ui-*)`, nunca literales
    (`TestThemeMatrixStateCoverageIsTokenDriven`).
- `TestSourceAppCSSKeepsCoreBeforeThemeCascade` — `sourceAppCSS` respeta el orden del build:
  `tokens.css` (core) antes que el theme (overrides).
- `TestThemeMatrixDefaultThemeIsPartOfDiscovery` — el theme que el repo entrega hoy (theme-material)
  es parte del discovery; la matriz no puede quedar vacía silenciosamente.

## 3. Checklist de un theme nuevo (del theme contract, §12)

Un theme pasa la suite si y solo si cumple **todo** lo siguiente:

- [ ] Define **todas** las familias obligatorias del contrato en su bloque light
      (`--ui-<componente>-*` para button-vía-color-semántico, text-field `--ui-field-*`, dialog, toast,
      card, badge, checkbox, radio, switch, slider, progress, fab, select, divider).
- [ ] Light + dark autocontenido en un solo `themes/<name>/theme.css`, con selector raíz `.theme-<name>`.
- [ ] Dark cubierto en **las dos rutas**: clase (`.<name>.theme-dark`/`.dark`/`[data-theme="dark"]`)
      **y** `@media (prefers-color-scheme: dark)`. O bien redefine la familia directamente, o bien
      garantiza que los `--ui-color-*` semánticos que la familia referencia se redefinen en ambas rutas.
- [ ] Naming canónico `--ui-*`; `danger` es el canónico (no `error`).
- [ ] Cero tokens muertos, cero gaps: todo `var(--ui-*)` referenciado por los componentes tiene definición.
- [ ] `npm run build` + `go test ./...` + `go vet ./...` verdes **sin aseverar valores Material**.
- [ ] Matriz `TestThemeMatrix*` pasando.
- [ ] Smoke: light/dark, narrow/wide, reduced motion, forced colors, no-JS, HTMX, estados
      empty/loading/error, teclado.
- [ ] Aceptación visual manual + `themes/<name>/README.md` con divergencias documentadas.
- [ ] Los componentes y contratos del core NO se modificaron para acomodar el theme.

El flujo operativo paso a paso está en `docs/gelium-ui-theme-implementation-guide.md`.

## 4. Cómo se extiende la matriz cuando exista theme-basecoat

1. Crear `themes/theme-basecoat/theme.css` con `.theme-basecoat { … }` (light + dark autocontenido).
2. **No se edita ningún test**: `availableThemes` descubre `theme-basecoat` por glob y la matriz lo
   recorre automáticamente.
3. `app.css` añade un import explícito (`@import "../../themes/theme-basecoat/theme.css";`) y la clase
   del theme se activa en runtime desde el template (mecanismo Phase H).
4. Si Basecoat necesita una familia nueva o una divergencia de variante, se resuelve como **decisión del
   theme contract** (extender el core con contrato, o documentar divergencia), nunca CSS de theme sobre
   markup distinto.
5. Si el test `TestThemeMatrixCoversEveryAvailableTheme` falla para basecoat, el error nombra
   exactamente la familia o el color semántico faltante por componente y por ruta dark — es la checklist
   automatizada.

## 5. Desincronizaciones conocidas y su guarda

| Área | Guarda | Estado (Phase C) |
|---|---|---|
| Orden `tokens.css` → theme en el build | `TestSourceAppCSSKeepsCoreBeforeThemeCascade` | ✅ sincronizado |
| Lista de imports de `sourceAppCSS` vs `app.css` | verificación de orden en el helper | ✅ sincronizado |
| Literales de color en componentes | `TestNoColorLiteralsInComponents` (excluye tokens/demo/app + bloques forced-colors) | ✅ cubre 50 archivos |
| Valor hex Material en tests | refactor a presencia (fab/dialog/toast ya migrados) | ✅ sin aserciones de valor |
