# Gelium UI — Font Contract

> Contrato de fuentes (webfonts) del sistema Gelium UI.
> Mecanismo de carga por theme + estándar de calidad. Se apoya en
> `docs/gelium-ui-theme-contract.md` (el theme define tokens `--ui-*`, no markup)
> y en el contrato tipográfico existente (medida 65ch, WCAG 1.4.12, steps
> GOV.UK/USWDS — ver `site/web/content/handbook-typography.md`).

---

## 1. Qué es una fuente en Gelium UI

Una fuente es un **recurso self-hosted, embebido y subseteado** que un theme
declara y pre-carga. NO es un `<link>` a Google Fonts/CDN en el markup, NO es
una dependencia npm, NO es un TTF/OTF servido en producción.

El theme sigue siendo una **dirección visual codificada como tokens `--ui-*`**;
la fuente alimenta los tokens `--ui-font-sans`, `--ui-font-mono` (y opcional
`--ui-font-display`) que los steps del typescale consumen vía
`--ui-type-*-family`. Mismo markup, misma API, mismo server contract — solo
cambia la apariencia.

---

## 2. Mecanismo (implementado)

1. **Los `.woff2` viven en `lib/fonts/`** (flat) y se embeben con
   `//go:embed fonts/*` en `lib.Assets` (`lib/assets.go`). Se sirven
   self-hosted en `/static/fonts/<file>` — **sin CDN**.
2. **El theme declara su set de fuentes** en `themeDirection.Fonts`
   (`internal/app/server.go`), con nombres flat `theme-<class>-<familia>.woff2`.
3. **`layout.html` pre-carga** cada fuente del theme resuelto con
   `{{themePreloadFonts .ThemeClass .AssetsVersion}}` → emite
   `<link rel="preload" as="font" href="/static/fonts/...?v=..." type="font/woff2" crossorigin>`.
4. **`@font-face` + `--ui-font-*`** se declaran en el `theme.css` del theme
   (junto a los demás tokens), apuntando a esas URLs locales.
5. **Namespace cerrado**: `fontAsset` solo sirve un `.woff2` si algún theme
   allowlisted lo declara por nombre exacto. Nada más es legible.

Pará una fuente: agregás el `.woff2` a `lib/fonts/`, la declará en
`themeDirection.Fonts` de ese theme, y el `@font-face` en su `theme.css`. Los
`<link>` de preload se generan solos.

---

## 3. Estándar de calidad (obligatorio por fuente)

### A. Entrega y rendimiento

1. **Self-hosted, `.woff2`** — nunca CDN externo, nunca TTF/OTF en producción.
   WOFF2 ya incluye compresión Brotli.
2. **Subsetting por Unicode range** — solo los glifos que el theme usa
   (`latin`, `+latin-ext` si hay clientes con acentos AR/BR/ES/PT). Payload 50%+
   menor.
3. **Solo los pesos necesarios** — p. ej. 400/500/600 para la sans + 400 para
   la display. No subir la family variable completa.
4. **`font-display: swap`** en todo `@font-face` (texto visible desde el inicio
   → buen LCP) **+ preload** para minimizar el swap (el preload ya lo emite el
   layout).
5. **Peso ≤ ~40 KB por weight** (woff2 subseteado real). Más pesado = falta
   subsetear.

### B. Estabilidad visual / CLS (el paso que la mayoría se salta)

6. **Fallback métricamente compatible** — todo `@font-face` de fuente real debe
   ir acompañado de un `@font-face` de fallback con los font-metric overrides:
   `size-adjust`, `ascent-override`, `descent-override`, `line-gap-override`,
   apuntando a una fuente de sistema (Inter/Georgia/UI). Esto **elimina el
   layout shift** cuando la fuente real reemplaza al fallback. Sin esto, `swap`
   = CLS.

   ```css
   @font-face {
     font-family: "Fallback for Alden Sans";
     src: local("Inter");
     size-adjust: 98.3%;
     ascent-override: 99%;
     descent-override: 27%;
   }
   ```

7. **Un solo `--ui-font-sans`/`--ui-font-mono` por theme** (fuente de verdad),
   que se propaga por todos los steps vía `--ui-type-*-family`. Cero markup.

### C. Legibilidad / accesibilidad (hereda del contrato actual)

8. La tipografía del theme debe cumplir el contrato tipográfico existente:
   body ≥ 16px, medida 65ch, WCAG 1.4.12 text spacing, line-height ≥ 1.5 en
   body, **tracking NO negativo en texto pequeño** (Linear usa `-0.011em` a
   13–16px — validar contra 1.4.12: tracking negativo en cuerpo es riesgo).
9. **Contraste tipográfico**: `--ui-color-fg` / `--ui-color-fg-muted` del theme
   deben pasar WCAG AA contra canvas/surface/surface-container. Fuentes nuevas
   no deben forzar tones que rompan contraste.
10. **Familia display solo para titulares** — el serif/display nunca en
    body/caption (regla que Alden mismo establece: "never use it in UI").

### D. Cobertura de idiomas

11. **Unicode range declarado por font**. Si un cliente es AR/BR/MX, el theme
    debe cubrir `latin-ext` (glifos es/por). El font stack incluye un fallback
    que cubra glifos fuera del subset sin que se vea roto (p. ej.
    `--ui-font-sans: "Alden Sans", ui-sans-serif`).

### E. Licencias (obligatorio)

12. **Licencia verificada y documentada** — la mayoría de Google Fonts son SIL
    OFL u otras libres, pero el contrato exige que cada fuente en `lib/fonts/`
    tenga su archivo de licencia y una nota en el README del theme. No asumir.

---

## 4. Verificación (automatizable)

- **Tests Go:** un test lee cada `theme.css` y verifica que todo `@font-face` es
  `.woff2` + `font-display: swap` + URL local (`/static/fonts/...`, no http) +
  declara Unicode range + los overrides de métricas en su fallback.
- **Mecanismo:** tests existentes (`internal/app/fonts_test.go`) verifican que
  el preload se emite correctamente, que el namespace `/static/fonts/` es
  cerrado, y que Material/Basecoat (sin fuentes) no emiten preload — sin
  regresión.
- **Contraste:** extender los tests AA existentes a los nuevos tokens de color
  que definan las fuentes.
- **Runtime (smoke):** verificar que el `AssetsVersion` aparece en los hrefs de
  preload y que el swap no produce CLS medible en la galería.

---

## 5. Documentación por theme

Cada theme con fuentes publica en su `themes/<theme>/README.md`:

- Fuentes incluidas (familia, pesos, archivos `lib/fonts/`).
- Token mapping: fuente → `--ui-font-*` / `--ui-type-*-family`.
- Subset/idiomas cubiertos y Unicode ranges.
- Fallback elegido y sus overrides de métricas (anti-CLS).
- Licencia de cada fuente.
- Divergencias documentadas.

---

**Gate**: una fuente que requiera cambiar markup, clases, contratos o
comportamiento no-JS es un cambio de core, no de fuente — se trata como tal. El
mecanismo está diseñado para que agregar un theme con fuentes **no toque** los
componentes ni el layout tipográfico.
