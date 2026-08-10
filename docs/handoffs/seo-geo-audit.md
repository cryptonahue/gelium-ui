# SEO & GEO Audit — Gelium UI (ex Gelium UI)

> Handoff read-only de investigación (SEO + Generative Engine Optimization). No modifica código, templates, CSS ni tests. Única salida: este documento.
> Alcance: auditar el estado SEO/GEO del sitio de documentación server-rendered de Gelium UI (Go `net/http` + `html/template` + `embed`, Markdown interno en `web/content/*.md`), y proponer un contrato SEO mínimo y un contrato GEO mínimo alineados con el system roadmap y las restricciones del sistema (sin CDN, sin JS obligatorio, HTML semántico primero, server-rendered).
> Baseline leído: `README.md`, `web/templates/layout.html`, `web/templates/{button,demo-whatsapp,demo-whatsapp-admin}.html`, `internal/app/{server,docs,routes,button,demo_whatsapp}.go`, `internal/app/{server,docs}_test.go`, `web/content/{index,button,dialog,data-table}.md`, `web/assets.go`, `web/static/`, `package.json`, `go.mod`, `docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-vocabulary.md`, `COMPONENT-ROADMAP.md`, `docs/handoffs/{mozilla-protocol-audit,vocabulary-audit,core-audit}.md`.

---

## 1. Estado SEO actual del sitio de docs

| Dimensión | Estado | Evidencia |
|---|---|---|
| `<title>` | ✅ Presente, server-driven por ruta | `layout.html:6` `{{.Title}} · Gelium UI`; cada handler setea `Title` (`server.go:123` home, `docs.go:97` /docs, `button.go:24`, etc.). Demo usa títulos hardcodeados y **stale**: `LoomChat · WhatsApp demo` (`demo-whatsapp.html:7`), `LoomChat · Modo administrador` (`demo-whatsapp-admin.html:7`) |
| `<meta name="description">` | ❌ Ausente | No existe en `layout.html:3-10` ni en ninguna página |
| `<link rel="canonical">` | ❌ Ausente | No existe; no hay URL canónica absoluta en el `<head>` |
| Open Graph / Twitter Cards | ❌ Ausente | Sin `og:*` ni `twitter:*` en ninguna página |
| `robots` (meta o robots.txt) | ❌ Ausente | Sin `<meta name="robots">` ni `/robots.txt`; `staticAsset` solo sirve `app.css`/`htmx.min.js`/`app.js` (`server.go:100-104`) |
| `sitemap.xml` | ❌ Ausente | No existe archivo ni ruta |
| `lang` | ✅ `en` | `layout.html:2` `<html lang="en" class="theme-material">` (también en los dos demos, `demo-whatsapp.html:3`, `demo-whatsapp-admin.html:3`) |
| H1 único | ✅ Sí (desde Markdown) | Todo `web/content/*.md` arranca con `# ...` (ej. `button.md:1`, `data-table.md:1`, `dialog.md:1`; home `index.md:1`; /docs generado `# Documentation`, `docs.go:85`). El layout **no** aporta su propio `h1` (la marca es `<a class="brand">`, `layout.html:13`) → un solo `h1` por página |
| Jerarquía de headings | ◐ H2 plano | `## ` sin subniveles en la mayoría (ej. `data-table.md:5,34,42,46,52,56,65,84,88`); contenido factual con headings descriptivos, **no interrogativos** |
| HTML semántico | ✅ Parcial | `<header>` + `<nav aria-label="Primary">` (`layout.html:12-15`), `<main class="docs-shell">` (`:16`), `<article class="prose">` (`:17`). **Falta `<footer>`** (layout termina en `toast-region`, `:51`; grep global: cero `<footer>`) |
| URLs limpias y estables | ✅ Sí | Registro único `componentRoutes()` (`routes.go:16-47`): `/`, `/docs`, `/components/{slug}` sin hashes ni `.html`. Query server-driven estable en data-table (`?sort=`, `?q=`, `?page=`, `data_table.go`) |
| Favicon / theme-color | ❌ Ausente | Sin `<link rel="icon">`, sin `theme-color` |

**Resumen**: el sitio tiene lo estructuralmente valioso (server-rendered, títulos por ruta, semántica HTML, URLs limpias, un h1) y le falta TODO el head SEO de marketing: description, canonical, OG/Twitter, robots, sitemap, structured data. Es un estado "SEO técnico de base, cero contrato de metadata".

---

## 2. Metadata actual: qué hay y qué falta

**Qué hay** — `layout.html:3-10` (head completo del layout):
```html
<meta charset="utf-8">            <!-- :4 -->
<meta name="viewport" ...>        <!-- :5 -->
<title>{{.Title}} · Gelium UI</title>  <!-- :6 -->
<link rel="stylesheet" href="/static/app.css?v=0.4.0">  <!-- :7 -->
<script defer src="/static/htmx.min.js?v=0.4.0"></script>  <!-- :8 -->
<script defer src="/static/app.js?v=0.4.0"></script>       <!-- :9 -->
```

**Qué falta** para un contrato SEO mínimo (server-driven por ruta):

1. **Fuente única de metadata por ruta** — hoy `pageView` (`server.go:14-50`) solo tiene `Title`, `Content`, `Nav` y datos de previews. No hay `Description`, `Canonical`, `Robots`, `OG`, `JSONLD`. La metadata se resolvería donde ya se resuelve todo lo demás: en el **handler** (`server.go:132-170` `renderMarkdownStatus`) poblada desde el **registro de rutas** (`routes.go:16-47`) + datos del documento.
2. **`meta description`** por ruta (derivable del intro de cada `*.md`).
3. **`canonical`** con URL absoluta (requiere base URL configurable — hoy no existe `BASE_URL`/host, `server.go` no usa `r.Host`).
4. **Open Graph + Twitter** (`og:title`, `og:description`, `og:url`, `og:type`, `twitter:card`).
5. **`robots`** por ruta (permite `noindex` en demos/POST como `POST /examples/*`, `server.go:77-81`, y en `/demo/whatsapp/admin`).
6. **JSON-LD** (ver §3).
7. **Favicon + `theme-color`** triviales.
8. **Datos del documento para GEO**: versión, fecha, autoría (ver §7).

---

## 3. Structured data (JSON-LD)

**Estado: cero.** Grep global: sin `application/ld+json`, sin `schema.org`, sin `itemprop`. No hay ningún JSON-LD en Go, templates ni contenido.

**Qué sería necesario para documentación de un sistema UI** (tipos schema.org):

| Tipo | Dónde | Por qué |
|---|---|---|
| `WebSite` | `/` (home) | Entidad raíz: `name` Gelium UI, `url`, `inLanguage`, `publisher` |
| `Organization` / `SoftwareSourceCode` | `/`, `/docs` | Autoría y repo (provenance) |
| `BreadcrumbList` | `/components/*` | Jerarquía Home → Docs → Componente (hoy no existe breadcrumb visual ni estructural, `vocabulary.md:226-230`) |
| `SoftwareApplication` / `WebPage` | cada `/components/*` | `name`, `applicationCategory` (DeveloperApplication), `softwareVersion` (0.4.0, `package.json:3`), `operatingSystem`, `license` (MIT, `README.md:126-128`) |
| `TechArticle` (opcional) | páginas de componentes | Envuelve el contenido dogfoodeado; `about`, `keywords` |

Implementación encaja con el stack: `pageView.JSONLD template.HTML` poblado en el handler y emitido en `layout.html` antes de `</head>`. Sin JS obligatorio (JSON-LD es declarativo).

---

## 4. Performance SEO

| Aspecto | Estado | Evidencia |
|---|---|---|
| Assets embebidos | ✅ | `web/assets.go:8` `//go:embed templates/*.html content/*.md static/*`; sin CDN (`README.md:24` "No se usa CDN") |
| CSS minificado | ✅ | Build lane Tailwind `--minify` (`package.json:6`) → `web/static/app.css` (131 075 bytes, 1 línea) |
| JS | ◐ | 2 scripts `defer` (`layout.html:8-9`); HTMX es enhancement, nunca requisito (`README.md:24`, `index.md:11`) |
| Orden de render | ✅ | CSS en `<head>` (`:7`), JS diferido, contenido antes que scripts |
| Cache | ❌ Débil | `staticAsset` fuerza `Cache-Control: no-cache` (`server.go:116`) + cache-busting por query `?v=0.4.0` (`layout.html:7-9`, testeado en `server_test.go:85-99`). **Sin `max-age`/ETag/immutable**: cada visita revalida los assets |
| Compresión | ❌ Ausente | Sin gzip/brotli middleware; las páginas HTML se escriben planas (`server.go:167-169`) |
| Imágenes | N/A | No hay imágenes en docs; demos usan emoji/SVG inline (`demo-whatsapp.html:20-21`) |
| Peso total assets | ⚠ | ~186 KB (`app.css` 131 KB + `htmx.min.js` 51 KB + `app.js` 3.7 KB) |

Oportunidades de bajo costo: `Cache-Control: public, max-age=31536000, immutable` para `/static/*` versionado (el `?v=` ya desambigua), gzip, favicon.

---

## 5. GEO (Generative Engine Optimization)

Evaluación de si el sitio es **citable** por motores generativos:

| Criterio GEO | Estado | Evidencia / Análisis |
|---|---|---|
| **Entidades claras** | ⚠ Inconsistente | Nombre del sistema aparece como "Gelium UI" (`index.md:1`, `layout.html:13`) pero **el README todavía dice `# Gelium UI`** (`README.md:1`) y los demos usan **"LoomChat"** (`demo-whatsapp.html:7,15`; `demo_whatsapp.go:222-223`). Tres identidades coexisten → ambigüedad de entidad. Versión solo implícita en query strings (`?v=0.4.0`), sin entidad visible |
| **Licencia** | ◐ En repo, no en web | MIT en `README.md:126-128` + `LICENSE`; **no expuesto en el sitio** ni en metadata machine-readable |
| **Autoría / ownership** | ❌ No hay | Sin autor/publisher/organization en ninguna página |
| **Provenance** | ❌ No hay | Sin fechas de publicación/actualización en `web/content/*.md` (ni frontmatter, goldmark default sin YAML) ni en el HTML |
| **Citaciones / URLs estables** | ✅ Fuerte | URLs limpias y estables (`routes.go:16-47`), contenido factual y específico por componente, dogfooding real (el preview es el componente real, `docs.go:86`) |
| **Headings interrogativos** | ❌ No | Headings descriptivos ("Anatomy", "States", "When to use it", `data-table.md:5,46,84`) — no formato pregunta ("What is X?", "How do I use X?"). Menos extraíble como respuesta a pregunta |
| **Machine-readable metadata** | ❌ Cero | Sin JSON-LD, sin meta description, sin Open Graph (§2-3) |
| **Ausencia de ambigüedad** | ⚠ Riesgo | Colisión de naming Gelium/Gelidium/LoomChat; docs de componentes en inglés sobre tema técnico estable ("Button", "Dialog") → contenido de bajo riesgo de desambiguación |

**Conclusión**: la fundación GEO es sólida (server-rendered, URLs estables, contenido factual dogfoodeado), pero el sitio es **invisible para citación enriquecida**: sin metadata estructurada, sin entidad de sistema/versión/licencia expuesta, sin fechas, sin autoría, y con identidad de marca dividida entre tres nombres.

---

## 6. Propuesta de contrato SEO mínimo (server-driven por ruta)

**Campos por ruta** — resolver en el handler (`internal/app/server.go`) y emitir en `web/templates/layout.html`:

| Campo | Ruta de origen | Emisión |
|---|---|---|
| `title` | `pageView.Title` (ya existe) + sufijo de marca | `layout.html:6` |
| `description` | intro del `*.md` (1ª frase) o dato del registro de rutas | `<meta name="description">` |
| `canonical` | base URL configurable (`BASE_URL` env, nueva) + `r.URL.Path` | `<link rel="canonical">` |
| `robots` | default `index,follow`; `noindex` en demos y rutas POST | `<meta name="robots">` |
| `og:*` + `twitter:card` | derivados de `title`/`description`/`canonical`; `og:type` website | `<meta property="og:...">` |
| `JSON-LD` | `pageView.JSONLD` (`template.HTML`) | `<script type="application/ld+json">` pre-`</head>` |
| `lang` | ya `en` (`layout.html:2`) | — |

**Dónde se implementaría:**
- **Handler**: extender `pageView` (`server.go:14-50`) con `Description`, `Canonical`, `Robots`, `JSONLD template.HTML`; un resolver `metaFor(route, pageView)` en el registro (`routes.go:16-47`) o un `map[path]meta` derivado de `componentRoutes()`; poblar en `renderMarkdownStatus` (`server.go:152-169`).
- **Template**: bloque condicional en `layout.html` head para description/canonical/robots/OG/JSON-LD (zero-JS, semántico).
- **Datos del documento**: frontmatter YAML simple por `web/content/*.md` (description, version, updated) parseado en el handler, o tabla Go por slug — la decisión de fuente única debe tomarse en el roadmap (grep confirma que goldmark se usa sin extensiones, `go.mod:5`, `server.go:63`).
- **Tests**: el patrón ya existe (`server_test.go:72-82` assert de contratos en HTML) — se extendería con asserts de `meta description`, canonical y JSON-LD.

---

## 7. Propuesta de contrato GEO mínimo

1. **Entidad única y versionada**: fijar "Gelium UI" como único nombre visible (reemplazar `# Gelium UI` en `README.md:1` y `LoomChat` en `demo-whatsapp*.html`), y exponer versión + licencia + repo en cada página (`WebSite`/`Organization`/`SoftwareApplication` JSON-LD con `softwareVersion` 0.4.0 y `license` MIT).
2. **Provenance**: fecha de publicación/actualización por documento (frontmatter o metadata Go) renderizada en el `article` (`layout.html:17`) y en JSON-LD (`datePublished`, `dateModified`); autoría en `Organization`/`author`.
3. **Citaciones**: mantener URLs estables (`routes.go`), añadir **BreadcrumbList** (JSON-LD + breadcrumb visual — componente ya planeado, `vocabulary.md:226-230`) y enlazar entre docs (ya existe cross-linking, ej. `data-table.md:86`).
4. **Headings citables**: convertir las secciones clave a formato pregunta ("What is Button?", "When to use a data table?") o añadir una sección FAQ por componente — la información ya existe en los `.md`.
5. **Machine-readable**: JSON-LD completo + `meta description` + Open Graph (§6) — prerrequisito para citación enriquecida.
6. **Estructura de docs citables**: cada `web/content/*.md` ya sigue "intención → anatomía → estados → cuándo usar" (ej. `data-table.md`); formalizar el patrón (intro factual de 1-2 oraciones, versión, fecha, `When to use it`, FAQ) como contrato de contenido para que el texto sea extraíble como respuesta verifiable.

---

## 8. Public content patterns relevantes para docs/sitios públicos

El roadmap de sistema **no tiene fases E/F** (tiene fases numeradas 0–8, `gelium-ui-system-roadmap.md:381-393`); la capa "public content patterns" ya fue auditada contra Mozilla Protocol en `docs/handoffs/mozilla-protocol-audit.md` (14 patrones). Estado cruzado con lo que existe hoy:

| Patrón | Hoy en Gelium UI | Brecha |
|---|---|---|
| **Hero/Billboard** | ◐ Ad-hoc: `hero-action` solo en home (`layout.html:18`, CTA "Read the docs" `server.go:124-128`) | Composición formal Hero = Section Heading + CTA + media (mozilla-protocol-audit §6.2) |
| **Breadcrumb** | ✖ No existe | Componente planeado (`vocabulary.md:226-230`); requisito SEO/GEO (BreadcrumbList) |
| **Footer** | ✖ No existe (`layout.html` termina en `main`+`toast-region`) | `<footer>` con nav secundaria, legal, locale (`<details>/<summary>` zero-JS) |
| **Article** | ◐ De facto: `<article class="prose">` (`layout.html:17`) | Formalizar como contrato tipográfico + intro + fecha/autoría |
| **CTA** | ✅ `ui-button` link (Button open-code, `button.html:6`, `README.md:112-116`) | Suficiente; opcional variante pública inline |
| **Card** | ✅ `ui-card` (`card.html`, `<article>/<a>/<button>`) | Faltan slots públicos: media, tag, meta, CTA (mozilla-protocol-audit §6.3) |
| **Language Switcher** | ✖ No existe | GET form + submit visible; requiere i18n/RTL |
| **Newsletter** | ✖ No existe | POST + 422 + success view (contratos server-driven canónicos) |
| **Notification Bar** | ✖ No existe (≈ Banner Gelium, `vocabulary.md:122-129`) | Base estática + dismiss POST+303 |
| **Section Heading / Split / Video** | ✖ No existen | Utilidad tipográfica / composición / contenedor media |

Recomendación: no duplicar el trabajo de `mozilla-protocol-audit.md`; integrar sus 14 patrones como extensión de Phase 2 del roadmap y añadir **SEO/GEO como fase/contrato propio** (este handoff) con dependencia de metadata + JSON-LD (§6-7).

---

## Top 5 gaps GEO (con ruta de evidencia)

1. **Identidad de marca dividida** — "Gelium UI" vs "Gelium UI" (`README.md:1`) vs "LoomChat" (`demo-whatsapp.html:7`, `demo_whatsapp.go:222-223`): unificar nombre + entidad `WebSite`/`Organization`.
2. **Cero machine-readable metadata** — sin JSON-LD, sin meta description, sin OG (grep global vacío; head solo `layout.html:3-10`): implementar contrato §6.
3. **Sin provenance ni autoría** — `web/content/*.md` sin fechas ni autor; sin `datePublished`/`dateModified`: añadir frontmatter + render en `article` (`layout.html:17`).
4. **Headings no interrogativos** — secciones descriptivas ("Anatomy", "States", `data-table.md:5,46`) en vez de formato pregunta/FAQ: reformular para extracción como respuesta.
5. **Sin breadcrumb ni jerarquía visible** — `vocabulary.md:226-230` marca el gap; requerido para `BreadcrumbList` y navegación de contexto: componente + JSON-LD.

---

## Fuentes

- Código: `internal/app/server.go`, `internal/app/docs.go`, `internal/app/routes.go`, `internal/app/button.go`, `internal/app/demo_whatsapp.go`, `internal/app/{server,docs}_test.go`.
- Templates: `web/templates/layout.html`, `web/templates/{button,demo-whatsapp,demo-whatsapp-admin}.html`.
- Contenido: `web/content/{index,button,dialog,data-table}.md` (29 docs en `web/content/`).
- Build/assets: `web/assets.go`, `web/static/`, `package.json`, `go.mod`.
- Docs del sistema: `README.md`, `docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-vocabulary.md`, `COMPONENT-ROADMAP.md`, `docs/handoffs/mozilla-protocol-audit.md`.
