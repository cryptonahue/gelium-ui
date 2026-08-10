# Mozilla Protocol Audit — Gelium UI (ex Gelium UI)

> Handoff read-only de investigación (Public content patterns). No modifica código, templates, CSS ni tests. Única salida: este documento.
> Alcance: auditar Mozilla Protocol (protocol.mozilla.org, v22.0.0) como referencia para la fase de *public/content patterns* del sistema Gelium UI, con MAPPING a Gelium y cumplimiento de las restricciones del sistema (server-rendered Go `html/template`, sin CDN, sin JS obligatorio, HTML semántico, tokens `--ui-*`, HTMX opcional, no-JS end-to-end).
> Baseline leído: `README.md`, `docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-vocabulary.md`, `docs/handoffs/vocabulary-audit.md`, `web/templates/layout.html`, `web/templates/card.html`, `web/content/index.md`.
> Evidencia externa (todas accesibles el 10/08/2026): `https://protocol.mozilla.org/`, `https://protocol.mozilla.org/components/detail/{article,billboard,breadcrumb,callout--default,card--overview,cta-link--default,feature-card--default,footer,language-switcher,newsletter--default,notification-bar--default,section-heading,split--default,video}`, `https://protocol.mozilla.org/docs/usage/framework`, `https://github.com/mozilla/protocol` (LICENSE).

---

## 1. Resumen ejecutivo

- **Mozilla Protocol** es el design system de los sitios web de Mozilla/Firefox (fuente: portada de `protocol.mozilla.org`). Hoy lo usa principalmente el equipo de Marketing Websites como front-end de `www.mozilla.org`. Sigue "evolving"; versión actual **22.0.0** (npm `@mozilla-protocol/core`).
- **Stack**: es una **librería de CSS (Sass) + módulos JS** entregada como paquete npm; el HTML se escribe a mano por el consumidor. No es un framework de templates ni de componentes JS. Usa Sass mixins, design tokens y CSS Custom Properties para themes (brand themes Mozilla/Firefox). **Licencia MPL 2.0** (verificado en `LICENSE` del repo). No impone framework: soporta `import` ES/CJS o `<script>` global.
- **De los 14 patrones objetivo, ninguno exige JS para su esencia**: Article, Billboard, Breadcrumb, Callout, Card, CTA Link, Feature Card, Section Heading, Split y Video son 100% estáticos en Protocol. Footer, Language Switcher y Newsletter **dependen de JS en Protocol** (accordion mobile, auto-submit, POST AJAX a Basket) — pero tienen **fallback no-JS nativo** (el `<button>Go</button>` del language switcher y el form submit normal del newsletter se muestran sin JS), lo que valida la tesis no-JS de Gelium.
- **Gap real en Gelium**: de los 14, solo **Card ✅**, **CTA ≈ Button link ✅** y el contenedor `.prose` (≈ Article) existen. Los otros 11 son **gaps**: Breadcrumb ✖ (ya marcado en `vocabulary.md:226-230`), Footer ✖ (no hay `<footer>` en `layout.html`, solo `site-header`), Callout/Banner/Notification bar ✖, Section Heading ✖, Billboard/Hero ✖, Split ✖, Video ✖, Feature Card ✖, Language Switcher ✖, Newsletter ✖.
- **Colisión de naming crítica a resolver**: el "Callout" de Protocol (sección prominente full-width centrada, `mzp-c-callout`) NO es el "Callout" del vocabulario Gelium (nota informativa ignorable, `vocabulary.md:131-138`). Misma palabra, anatomías opuestas. Antes de implementar hay que decidir el canónico Gelium.
- **Clasificación propuesta**: 6 componentes nuevos con contrato (Breadcrumb, Footer, Callout/Banner público, Newsletter, Language Switcher, Video), 4 composiciones (Billboard/Hero, Feature Card, Split, Section Heading), 2 extensiones de existentes (Card público, CTA Link), y diferir variantes JS de Protocol (sticky/scripted, auto-submit, AJAX-only).
- **Anti-patterns a no copiar**: branding Mozilla/Firefox (logo, wordmark, zap, themes de marca), el grid propietario `mzp-l-*` (columns, card-layout, content-container), la familia de JS (`footer.js`, `details.js`, `lang-switcher.js`, `newsletter.js`), la clase de naming `mzp-*`, y las variantes horizontales de Feature Card (ya **deprecadas upstream** a favor de Split).
- **Encaje en roadmap**: el roadmap (`gelium-ui-system-roadmap.md`) tiene fases 0–8 con números, **no existen fases E/F**. La capa "public/content patterns" no está prevista; se propone como extensión de Phase 2 (vocabulario) o fase propia paralela a Phase 3, usando los contratos server-driven ya canónicos (GET params, POST + 303, HTTP 422 + `X-Gelium-Validation`, `loom:toast`).

---

## 2. Qué es Mozilla Protocol

| Atributo | Dato | Evidencia |
|---|---|---|
| Propósito | "A design system for Mozilla and Firefox websites. Establishes a common design language, provides reusable coded components, and outlines high level guidelines for content and accessibility." | https://protocol.mozilla.org/ |
| Uso real | Front-end de `www.mozilla.org` por Mozilla Marketing Websites; objetivo declarado: design system unificado para cualquier sitio Mozilla | https://protocol.mozilla.org/ |
| Estado | "Protocol is still an evolving project"; v22.0.0; npm `@mozilla-protocol/core` | https://protocol.mozilla.org/ (badge npm), https://www.npmjs.com/package/@mozilla-protocol/core |
| Stack | **CSS Sass** (mixins, design tokens, Custom Properties por brand theme) + **módulos JS** (`MzpFooter`, `MzpLangSwitcher`, `MzpNewsletter`, `MzpDetails`) importables vía Webpack/ESM/CJS o `<script>` global. **Sin framework de templates**: markup HTML a mano | https://protocol.mozilla.org/docs/usage/framework |
| Licencia | **MPL 2.0** (Mozilla Public License v2.0) | https://github.com/mozilla/protocol/blob/main/LICENSE |
| No-JS | Los componentes con JS (Footer, Language Switcher, Newsletter, Notification bar scripted) declaran **fallbacks nativos**: el form del language switcher muestra un `<button type="submit">Go</button>` "when JavaScript is not enabled" | https://protocol.mozilla.org/components/detail/language-switcher, /newsletter--default |

Observación para Gelium: Protocol resuelve *themes* con Custom Properties en `:root` por brand theme (`_themes.scss`) — mismo mecanismo conceptual que el theme contract Gelium (Phase 5), aunque Gelium lo formaliza con tokens `--ui-*` en `<html class="theme-*">` (`layout.html:2`).

---

## 3. Distinción: "content/public pattern" vs "componente de aplicación"

**Criterio del sistema (basado en `gelium-ui-system-roadmap.md` capas, `vocabulary.md` §1 y `composition-audit.md`):**

- **Componente de aplicación**: participa en un **workflow con estado** — entrada de datos densa, validación, transiciones de estado, comparación de datos (Form, Data table, Select menu, Steps, Queue, Board, Resource editor). Su valor está en la operación, no en la presentación.
- **Content/public pattern**: **comunica y presenta contenido a una audiencia** — lectura, navegación, suscripción, promoción, contexto editorial. Predominantemente **estático**, sin workflow transaccional ni estado que transicionar; si tiene form, es de propósito único (opt-in / locale) y encaja en los contratos server-driven existentes.

**Justificación por patrón (todos public, ninguno app):**

| Patrón | Por qué es public/content, no app |
|---|---|
| Article | Lectura editorial long-form; el contenido ES el producto; sin interacción. |
| Billboard/Hero | Landing/marketing; primera impresión; CTA a navegación. |
| Breadcrumb | Navegación de contenido por jerarquía de sitio; 100% estático; sin workflow. |
| Callout | Contenido promocional/informativo prominente; no resuelve tarea operativa. |
| Card | Exploración/scan de contenido editorial (media + tag + meta + CTA); el card Gelium actual también sirve KPIs de app, pero el patrón *público* es el de contenido. |
| CTA Link | Acción promocional de navegación (enlace destacado), no transacción. |
| Feature Card | Variante editorial grande de Card; promocional. |
| Footer | Chrome del sitio público: navegación secundaria, legal, locale; sin estado. |
| Language Switcher | Localización de contenido; form de propósito único (GET locale), no workflow. |
| Newsletter | Suscripción de contenido (opt-in marketing); form de propósito único, no workflow de datos. |
| Notification Bar | Feedback contextual persistente sobre contenido público (≈ Banner Gelium); distinto del Toast de app (transitorio, `loom:toast`). |
| Section Heading | Organización de contenido en secciones; tipografía, no interacción. |
| Split | Layout promocional media/texto (dos columnas, RTL-aware); presentación pura. |
| Video | Presentación de media con controles nativos; sin workflow. |

Regla derivada: **ningún patrón de esta lista transiciona estado ni recolecta datos de negocio**; los dos que usan forms (Newsletter, Language Switcher) son conversión de audiencia y se resuelven con los contratos server-driven canónicos (GET locale / POST subscribe + 422 + success view), sin JS obligatorio.

---

## 4. Los 14 componentes: anatomía, semántica, JS y MAPPING a Gelium

Leyenda: **✅ existe en Gelium** · **◐ parcial (patrón ad-hoc / composición posible)** · **✖ gap (crear)** · **⏸ diferir variantes JS**.

### 4.1 Resumen mapping

| # | Patrón Protocol | Root semántico Protocol | JS en Protocol | JS en Gelium | Mapping a Gelium | Estado |
|---|---|---|---|---|---|---|
| 1 | Article | `<article class="mzp-c-article">` + h1 + `.mzp-c-article-intro` | 0 | 0 | `.prose` de `layout.html:17` (contenido Markdown) — formalizar tipografía | ✅ |
| 2 | Billboard (Hero) | `<div class="mzp-c-billboard mzp-l-billboard-left\|right">` + img + `.mzp-c-billboard-title/desc` + cta-link | 0 | 0 | `hero-action` (`layout.html:18`) = ad-hoc home; crear composición Hero | ◐ |
| 3 | Breadcrumb | `<nav aria-label="breadcrumbs" class="mzp-c-breadcrumb">` + `<ol>` + `<li aria-current="page">` | 0 | 0 | **Gap** (`vocabulary.md:226-230`); la semántica Gelium coincide 1:1 | ✖ |
| 4 | Callout | `<section class="mzp-c-callout">` + `.mzp-l-content` + title + desc; full-width centrado | 0 | 0 | **Colisión de naming** con Callout Gelium (`vocabulary.md:131-138`); decidir canónico | ✖ |
| 5 | Card | `<section class="mzp-c-card mzp-has-aspect-3-2">` + `<a class="mzp-c-card-block-link">` + media/tag/title/desc/cta/meta | 0 | 0 | Card Gelium (`card.html`); **añadir slots público** (media, tag, aspect, CTA) | ✅+ext |
| 6 | CTA Link | `<a class="mzp-c-cta-link">` | 0 | 0 | Button link Gelium (`button.html`, README:112-116) | ✅ |
| 7 | Feature Card | `<section class="mzp-c-card-feature mzp-has-aspect-16-9">` + media + content + cta | 0 | 0 | Composición Card + CTA Link; horizontales **deprecadas upstream** → no implementar | ◐ |
| 8 | Footer | `<footer class="mzp-c-footer">` + nav primario (logo + sections `h5`+`ul`) + nav secundario (lang switcher + social + legal) | **Sí** (`MzpFooter` + `MzpDetails` accordion mobile) | 0 (`<details>/<summary>` nativos) | **Gap**: `layout.html` solo tiene `site-header` (líneas 12-15), sin `<footer>` | ✖ |
| 9 | Language Switcher | `<form class="mzp-c-language-switcher" method="get">` + label + select + submit | **Sí** (`MzpLangSwitcher.init()` auto-submit) | 0 (GET form con submit siempre visible) | **Gap**; requiere `lang`/RTL | ✖ |
| 10 | Newsletter | `<aside class="mzp-c-newsletter">` + `<form method="post">` email + country + lang + checkboxes + privacy | **Sí** (`MzpNewsletter.init()` POST AJAX a Basket + validación) | 0 (POST server; HTMX opcional) | **Gap**; contrato 422 + success view | ✖ |
| 11 | Notification Bar | `<aside class="mzp-c-notification-bar">` + `<p>` + variantes tone + dismiss + cta | 0 base; variantes *scripted/sticky* sí | 0 (dismiss = POST + 303) | **≈ Banner Gelium** (`vocabulary.md:122-129`); alias Notification Bar | ✖ |
| 12 | Section Heading | `<h2 class="mzp-c-section-heading">` (centered) | 0 | 0 | **Gap**; = utilidad tipográfica, no componente | ✖ |
| 13 | Split | `<section class="mzp-c-split">` + `.mzp-c-split-body` + `.mzp-c-split-media`; bidi RTL | 0 | 0 | **Gap**; composición media + tipografía + CTA | ✖ |
| 14 | Video | `<div class="mzp-c-video">` + `<video controls>` o `<iframe>`; aspect 16:9/4:3 | 0 (controles nativos) | 0 | **Gap**; contenedor media responsive | ✖ |

### 4.2 Detalle por patrón (anatomía + variantes + evidencia)

1. **Article** — `<article class="mzp-c-article">` con `h1.mzp-c-article-title` y `p.mzp-c-article-intro` (intro con font-size mayor). "Should be the primary content on the page." Variantes: ninguna (intro opcional). **Gelium**: el `.prose` de `layout.html:17` ya renderiza el contenido Markdown como `<article>`; el patrón existe de facto. Formalizar como contrato tipográfico (`--ui-type-display-lg`/`title-md` — gaps ya listados en core-audit) + opcional intro. **Tokens**: `--ui-type-*`, `--ui-color-*`, `--ui-space-*`.

2. **Billboard (Hero)** — `<div class="mzp-c-billboard mzp-l-billboard-left|right">` con `.mzp-c-billboard-image-container` + `.mzp-c-billboard-content` + `h2.mzp-c-billboard-title` + `p.mzp-c-billboard-desc` + `a.mzp-c-cta-link`. "The image and copy stack vertically on small screens, and run full-width horizontally on larger screens." Imagen 346×346; **título ≤35 chars, desc ≤100 chars** (copywriting, no contrato). RTL: layout NO se invierte. **Gelium**: home usa `hero-action` (`layout.html:18`); proponer composición Hero = Section Heading + CTA Link + media opcional. **Tokens**: `--ui-type-display-*`, `--ui-space-*`, `--ui-color-*`.

3. **Breadcrumb** — `<nav aria-label="breadcrumbs" class="mzp-c-breadcrumb">` + `<ol class="mzp-c-breadcrumb-list">` + `<li class="mzp-c-breadcrumb-item">` con `<a>`; actual = `<li aria-current="page">` sin link. Variante dark `mzp-t-dark`. **Semántica idéntica a la propuesta Gelium** (`vocabulary.md:229`). Componente nuevo, 100% estático. **Tokens**: `--ui-type-body-*`, `--ui-color-primary/outline`, `--ui-space-*`.

4. **Callout** — `<section class="mzp-c-callout">` + `.mzp-l-content` + `.mzp-c-callout-content` (title + desc). Full-width, texto centrado, no para long-form. Variantes: fondo secondary/tertiary/dark, ancho sm/md/lg/xl, `mzp-l-compact`. Refactor reciente: reemplazó a `mzp-c-call-out` (con guion). **No-nos**: no anidar en otro `mzp-l-content`. **Colisión con Gelium**: ver §7. **Gelium**: elegir canónico (recomendado: "Promo Section" o reusar "Callout" Gelium para la nota). Estático.

5. **Card** — `<section class="mzp-c-card mzp-has-aspect-{1-1,3-2,16-9}">` + `<a class="mzp-c-card-block-link">` + media wrapper + content (tag opcional con icon `mzp-has-video/audio`, title ≤50 chars, desc ≤150, cta, meta). Variantes de tamaño: extra-small / small (default) / medium / large; dark theme. **Card Gelium** (`card.html`): `<article>`/`<a>`/`<button>` con title/body y elevación elevated/filled/outlined — **sin media, sin tag, sin aspect ratio, sin CTA**. La raíz semántica de Protocol (`<a>` block-link sobre `<section>`) choca con la regla Gelium "el control interno es el foco" (`vocabulary.md:35`) — conservar el patrón Gelium (a/button como control) y añadir slots. **Tokens**: `--ui-card-*`, `--ui-radius-*`, `--ui-elevation-*`, `--ui-type-*`.

6. **CTA Link** — `<a class="mzp-c-cta-link">`. "A prominent link that stands apart... usually part of another component such as a Card or Picto." Variantes: sizes, con iconos. **Gelium**: Button con `Href` ya renderiza `<a>` con states disabled/loading y `aria-disabled`/`aria-busy` (README:112-116); cubre la necesidad. Opcional: clase pública de "cta link" inline con icono para composiciones editoriales. **Tokens**: `--ui-color-primary`, `--ui-type-*`, `--ui-focus-ring`.

7. **Feature Card** — `<section class="mzp-c-card-feature mzp-has-aspect-{16-9,3-2}">` + media + content (title + desc + cta). **"The horizontal variants for this component are deprecated. Use the Split component instead."** → señal fuerte: en Gelium implementar como **composición Card + CTA Link**, no componente propio; las variantes horizontales no se portan. Estático.

8. **Footer** — `<footer class="mzp-c-footer">` + `.mzp-c-footer-primary` (logo + `.mzp-c-footer-sections` de `section.mzp-c-footer-section` con `h5.mzp-c-footer-heading` + `ul.mzp-c-footer-list`) + `.mzp-c-footer-secondary` (language switcher, social, legal/terms). **JS en Protocol**: `MzpFooter.init()` + `MzpDetails` para secciones expandibles en móvil. **Gelium**: no existe (`layout.html` termina en `main` + `toast-region`, sin `<footer>`). Implementar con `<details>/<summary>` nativos para el accordion móvil → **zero JS**. **Tokens**: `--ui-color-surface`, `--ui-type-*`, `--ui-space-*`.

9. **Language Switcher** — `<form class="mzp-c-language-switcher" method="get" action="#">` + `<label>` + `<select class="mzp-js-language-switcher-select">` + `<button type="submit">Go</button>` ("shown when JavaScript is not enabled"). JS: auto-submit on change + callback analytics. Asume primer path de la URL = locale (`/en-US/...`). **Gelium**: form GET con `select` nativo + botón submit siempre visible (no hace falta JS para nada; HTMX opcional para swap del `<html lang>`). Requiere i18n (etiqueta localizada) y manejo de `lang`/RTL. **Tokens**: `--ui-color-surface`, `--ui-select-*` reuso, `--ui-space-*`.

10. **Newsletter** — `<aside class="mzp-c-newsletter">` + `<form method="post" action="https://basket.mozilla.org/news/subscribe/">` + `source_url` hidden + header (title + tagline) + `fieldset.mzp-c-newsletter-content` (email required, country select, language select, newsletter checkboxes, privacy checkbox required) + submit + `mzp-c-newsletter-thanks` success block. JS: `MzpNewsletter.init()` = POST AJAX + validación client + callbacks analytics; **fallback no-JS = form POST normal**. **Gelium**: form POST server-driven con contrato 422 + `X-Gelium-Validation` + vista success (reusa el patrón de validación existente, README:35); HTMX opcional. Requiere i18n (copy, selects) y privacy (accesible). **Tokens**: `--ui-color-*`, `--ui-text-field-*`, `--ui-button-*`, `--ui-space-*`.

11. **Notification Bar** — `<aside class="mzp-c-notification-bar">` + `<p>` (mensaje corto, 1 frase). Variantes tone: `mzp-t-success|warning|error|click`; botón dismiss opcional; `mzp-c-notification-bar-cta` para acción prominente; variantes *scripted* y *sticky* requieren JS. **Gelium**: mapea al gap **Banner** (`vocabulary.md:122-129`, persistente a nivel página/sitio) — Protocol lo cataloga como "contextual feedback", semánticamente el hermano público del Toast Gelium. Base estática; dismiss = form POST + 303; **diferir** sticky/scripted (JS). **Tokens**: `--ui-color-success/warning/error/info`, `--ui-color-surface`, `--ui-radius-*`.

12. **Section Heading** — `<h2 class="mzp-c-section-heading">` (centered, más aire que un heading normal). **Gelium**: no es un componente con contrato — es **utilidad tipográfica** sobre `--ui-type-*` (heading + variante centered). Propuesta: definir en vocabulario/composición, no crear template. Estático.

13. **Split** — `<section class="mzp-c-split">` + `.mzp-c-split-container` + `.mzp-c-split-body` (tipografía NO aplicada por defecto) + `.mzp-c-split-media` (img/video). Variantes: `mzp-l-split-reversed` (se invierte en RTL), body narrow/wide (33/66), h/v alignment, media overflow / constrain-height, pop top/bottom, center-on-sm-md, hide-media-on-sm-md; ancho md/lg/xl (sin sm). **No-nos**: no anidar en `mzp-l-content`. **Gelium**: composición de dos columnas responsive (stack en mobile) con media + tipografía + CTA; puede ser patrón de composición sin componente propio o un template `split.html` mínimo. **RTL automático** (bidi) — alineado con el requisito RTL del sistema. **Tokens**: `--ui-space-*`, breakpoints, `--ui-type-*`, `--ui-radius-*`.

14. **Video** — `<div class="mzp-c-video">` + `<video controls poster>`/`<source>` o `<iframe>` (YouTube). Aspect 16:9 default / `mzp-has-aspect-4-3`. "Best used inside another component" (Split, Card, Callout). **Gelium**: contenedor media responsive mínimo, controles nativos, poster + `loading="lazy"`, `aspect-ratio` CSS. Estático. **Tokens**: `--ui-radius-*`, `--ui-color-*`.

---

## 5. Compatibilidad con las restricciones del sistema

| Patrón | 100% estático | Requiere JS en Protocol | Requiere JS en Gelium | i18n/l10n |
|---|---|---|---|---|
| Article | ✅ | — | 0 | copy length sensible al idioma |
| Billboard/Hero | ✅ | — | 0 | copy length sensible (35/100 chars); RTL no se invierte |
| Breadcrumb | ✅ | — | 0 | etiquetas de nivel localizadas |
| Callout | ✅ | — | 0 | copy corto centrado |
| Card | ✅ | — | 0 | título/desc localizables |
| CTA Link | ✅ | — | 0 | texto de acción localizable |
| Feature Card | ✅ | — | 0 | ídem Card |
| Footer | ✅ (Gelium) | Sí (accordion mobile) | **0** (`<details>/<summary>`) | etiquetas de secciones; `lang` |
| Language Switcher | ✅ | Sí (auto-submit) | **0** (GET form + submit visible) | **crítico**: locales, `<html lang>`, RTL, primer path URL |
| Newsletter | ✅ | Sí (AJAX Basket) | **0** (POST server, H opcional) | **crítico**: copy, country/language selects, privacy |
| Notification Bar | ✅ base | solo variantes scripted/sticky | **0** (dismiss = POST + 303) | copy localizable |
| Section Heading | ✅ | — | 0 | texto localizable |
| Split | ✅ | — | 0 | **RTL automático** (bidi) |
| Video | ✅ | — | 0 | captions/transcript, poster |

**Conclusión**: los 14 son compatibles con "no-JS end-to-end". Los tres que usan JS en Protocol (Footer, Language Switcher, Newsletter) tienen **alternativa nativa Gelium**: `<details>/<summary>`, form GET con submit visible y POST server-side con 422. Solo Language Switcher y Newsletter requieren **regionalización real** (locales + RTL + copy).

---

## 6. Propuesta de clasificación Gelium

### 6.1 Componentes nuevos con contrato (template + handler + test)

| Componente | Contrato server | Observaciones |
|---|---|---|
| **Breadcrumb** | Datos: slices de (Label, Href, IsCurrent) | Semántica idéntica a Protocol; cero JS |
| **Footer** | Layout site (`Nav`, `Sections`, `Legal`, opcional `LangSwitcher`) | `<details>/<summary>` nativos para móvil; cero JS |
| **Callout/Banner público** | Variantes tone + título + cuerpo + CTA opcional | **Resolver colisión de naming §7.1 antes de implementar** |
| **Notification Bar** | Variantes tone + mensaje + dismiss opcional | dismiss = POST + 303 + Toast; variantes sticky/scripted **diferidas** |
| **Newsletter** | POST + `X-Gelium-Validation` 422 + success view | reusa text-field/select/button/checkbox Gelium; i18n |
| **Language Switcher** | GET `?lang=` → 303 a URL locale | form nativo + submit visible; i18n/RTL |
| **Video** | `Src`, `Poster`, `Aspect`, `Transparent`/caption | contenedor media responsive |

### 6.2 Composiciones de componentes existentes (patrón en vocabulario, sin template nuevo o template mínimo)

- **Billboard/Hero** = Section Heading + desc + CTA Link + media opcional (hoy `hero-action`, `layout.html:18`).
- **Feature Card** = Card + CTA Link; NO portar variantes horizontales (deprecadas upstream, §7).
- **Split** = dos columnas (body + media) sobre Card/Video + tipografía + CTA; RTL bidi; posible `split.html` mínimo si el patrón se repite.
- **Section Heading** = utilidad tipográfica `--ui-type-*` + variante centered; no componente.

### 6.3 Extensiones de existentes

- **Card** → añadir slots públicos: media (con aspect-ratio), tag, meta, CTA; conservar la raíz `<a>/<button>` como control (regla `vocabulary.md:35`, divergencia con Protocol §7).
- **CTA Link** → variante pública inline con icono (opcional; Button link ya cubre la acción).

### 6.4 Diferir

- Variantes JS de Protocol: Notification bar *scripted/sticky*, Language Switcher *auto-submit*, Newsletter *AJAX-only*, Footer accordion JS.
- Grid/layout propietario (`mzp-l-content`, Columns, Card Layout, Content Container, Main-with-sidebar).
- Branding (Logo, Wordmark, Zap) y themes de marca — fuera de Gelium por definición (sin branding copiado).

---

## 7. Divergencias y anti-patterns (qué NO copiar)

1. **Branding Mozilla/Firefox**: logo, wordmark, zap, brand themes, colores de producto. Gelium es multi-theme con tokens semánticos `--ui-*`, nunca marca.
2. **Grid propietario `mzp-l-*`**: content-container, columns, card-layout, main-with-sidebar. Gelium usa responsive fluido (auto-fit/minmax, `min()`/`clamp()`, `composition-audit.md`).
3. **Familia JS de Protocol**: `footer.js`, `details.js`, `lang-switcher.js`, `newsletter.js`, paquete npm `@mozilla-protocol/core` (Sass/Webpack). Gelium: templates Go + Tailwind + cero JS de componente.
4. **Naming `mzp-*`** y Sass mixins: Gelium tiene su propio naming `ui-*` y tokens; se adopta la *semántica* del patrón, no el CSS.
5. **Feature Card horizontal**: deprecada upstream a favor de Split → no implementar.
6. **Anti-pattern de anidado** (no-nos de Protocol): Callout y Split tienen contenedor interno propio; no anidar en otro contenedor de contenido — incorporar como regla de composición Gelium.
7. **Guías de copy como contrato**: límites de chars (35/100 billboard, 50/150 card) son *copywriting guidelines* de un sitio editorial de marca; adoptar como recomendación, no como validación de componente.
8. **Card block-link**: Protocol envuelve todo el card en un `<a>` sobre `<section>`; Gelium mantiene "el control interno es el foco" (`vocabulary.md:35`) — no portar el block-link completo como tabstop.
9. **Variantes JS-dependentes** (sticky/scripted, auto-submit, AJAX): violan la regla "ningún término exige JS para el flujo principal" (`vocabulary.md:311`) si se portan tal cual.

### 7.1 Colisión de naming pendiente de resolución

| Término | Protocol (protocol.mozilla.org) | Vocabulario Gelium (`vocabulary.md:131-138`) |
|---|---|---|
| **Callout** | Sección **prominente full-width centrada** (title + desc + CTA + media), hero-like | Nota **informativa ignorable** (`<aside>` con heading opcional), tipo tip box |
| **Notification Bar** | Feedback contextual con tones success/warning/error/click + dismiss | ≈ **Banner** Gelium (persistente página/sitio) |

Acción requerida en Phase 2/3: decidir si Gelium adopta el término "Callout" con anatomía Gelium (tip box) y nombra el patrón Protocol como "Promo Section/Hero", o viceversa. No implementar ninguno hasta resolver el canónico.

---

## 8. Encaje en el roadmap

El roadmap tiene fases **numeradas 0–8** (`gelium-ui-system-roadmap.md`); **no existen fases E/F**. Propuesta de encaje de la capa "Public content patterns":

- **Extender Phase 2** (vocabulario): añadir los 14 términos a la capa public/content con su clasificación (componente/composición/utilidad).
- **Fase propia o paralela a Phase 3** (composición): definir reglas públicas (Hero = composición, Feature Card = composición, Section Heading = utilidad) y resolver la colisión Callout/Notification Bar antes de Phase 4 (recetas públicas).
- **Dependencias**: tokens de Phase 1 (typography `--ui-type-*` incl. gaps `display-lg`, spacing, colores semánticos success/warning/info — ver core-audit) y componentes base ya existentes (Button, Card, Text field, Select, Checkbox).
- **Contratos**: reusar los 4 contratos server-driven canónicos (`composition-audit.md` §9); Newsletter y Language Switcher no inventan contrato nuevo.
- **Gate Phase 5**: los templates públicos son markup Gelium único sobre tokens; los themes solo aportan tokens (sin duplicar markup por theme).

---

## 9. Top 5 recomendaciones

1. **Clasificar la capa public/content explícitamente** en el roadmap (extensión de Phase 2 + reglas de composición paralelas a Phase 3), separándola de los componentes de aplicación por el criterio del §3 (presentación/conversión vs workflow con estado).
2. **Resolver las dos colisiones de naming antes de implementar** (§7.1): "Callout" (Protocol ≠ Gelium) y "Notification Bar" (≈ Banner Gelium). Fijar canónicos en `gelium-ui-vocabulary.md`.
3. **Implementar primero los componentes 100% estáticos con semántica nativa**: Breadcrumb, Section Heading (utilidad), Video, Card público (slots media/tag/CTA) y Banner/Notification Bar base — todos zero-JS, directos.
4. **Portar solo con alternativa nativa los 3 que usan JS en Protocol**: Footer con `<details>/<summary>`, Language Switcher con GET form + submit visible, Newsletter con POST + 422 + success view (HTMX opcional) — preservando el principio no-JS end-to-end.
5. **No portar JS, grid propietario ni branding de Protocol** (§7): adoptar la semántica de patrones y las guías de copy como referencia, nunca el CSS/Sass `mzp-*`, los módulos JS ni los themes de marca; Feature Card horizontal se descarta por deprecación upstream.

---

## Fuentes

- https://protocol.mozilla.org/ — portada: propósito, estado, versionado.
- https://protocol.mozilla.org/docs/usage/framework — stack Sass/tokens/Custom Properties/themes/mixins.
- https://protocol.mozilla.org/components/detail/{article,billboard,breadcrumb,callout--default,card--overview,cta-link--default,feature-card--default,footer,language-switcher,newsletter--default,notification-bar--default,section-heading,split--default,video} — anatomía, clases, variantes, JS y no-nos por componente.
- https://github.com/mozilla/protocol/blob/main/LICENSE — licencia MPL 2.0.
- https://www.npmjs.com/package/@mozilla-protocol/core — paquete npm v22.0.0.
- Gelium: `README.md`, `docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-vocabulary.md`, `docs/handoffs/{vocabulary,core,composition}-audit.md`, `web/templates/layout.html`, `web/templates/card.html`, `web/content/index.md`.
