# Gelium UI

> Themeable, open-code UI components for Tailwind CSS and HTMX.

Current release: **v0.5.3** — monorepo split (publishable `lib/` package + `site/` docs consumer), HTMX 4 (namespaced events, innerMorph, optimistic theme/scheme toggling with server authority), on-this-page rail and GOV.UK-style prev/next pagination, dark mode across standalone recipes, unified 65ch reading column, mobile foundations (touch targets, safe areas, reduced motion, overflow containment contracts).

Gelium UI is a server-rendered component library: semantic HTML, copyable components, token-driven `--ui-*` themes, Tailwind CSS 4 and HTMX served locally. No CDN, no client framework, no hydration. JavaScript is progressive enhancement — everything works without it.

## Dos audiencias

- **Instalar y usar** (desarrollador o LLM): `npm install gelium-ui` (pendiente de publish) — CSS por componente, themes, templates y un bundle prebuilt `dist/gelium.css` + consumer JS `js/gelium.js`. El handler Go embebido es el modo avanzado.
- **Contribuir** (este repo): monorepo npm-workspaces — `lib/` es el paquete publicable, `site/` es el consumidor que lo importa por nombre de paquete (dogfooding).

## Requisitos

- Go 1.24 o posterior.
- Node.js 20 o posterior con npm.

## Setup y build reproducible (contribuidor)

Desde la raíz del proyecto:

```bash
npm install
go mod download
npm run build
```

`npm run build` (raíz) encadena: build del site (`tailwindcss` sobre `site/web/styles/app.css` → `site/web/static/app.css`, copia de `htmx.min.js` e iconos) + build del dist de la librería (`lib/styles/dist-entry.css` → `lib/dist/gelium.css`) + copia de `lib/js/gelium.js` → `site/web/static/gelium.js`. Los artefactos finales dentro de `site/web/static/` y `lib/dist/` se mantienen como artefactos de build y son embebidos en el binario Go mediante `embed`.

## Tests y checks (gates)

```bash
go test ./internal/... ./site/... ./lib/...
go vet ./internal/... ./site/... ./lib/...
npm run build
git diff --check
gofmt -l internal site lib
```

NUNCA uses `go test ./...`: el árbol incluye paquetes con imports pre-split que los gates excluyen a propósito. Los tests usan `httptest` y verifican handlers, forms URL-encoded, HTTP 422/200, render Markdown, contratos HTML/ARIA de componentes, integración HTMX, CSS/tokens y assets embebidos sin snapshots completos.

## Estructura real

```text
gelium-ui/
├── lib/                        ← paquete npm publicable (gelium-ui)
│   ├── assets.go               (//go:embed templates + styles → geliumui/lib)
│   ├── version.go              (AssetsVersion: cache-bust centralizado)
│   ├── package.json            (exports map, files, sideEffects)
│   ├── dist/gelium.css         (bundle prebuilt: preflight + themes + componentes)
│   ├── js/gelium.js            (consumer JS: toast, 422 contract, view transitions, slider)
│   ├── styles/                 (tokens.css, base.css foundation, 47 componentes, index.css manifest)
│   ├── templates/              (45 partials de componentes)
│   ├── themes/                 (theme-material.css, theme-basecoat.css — flat)
│   └── *_test.go               (suites de contrato de los componentes)
├── site/                       ← consumidor de demostración (docs, landing, recipes)
│   ├── package.json            (depende de gelium-ui@0.5.3 por nombre de paquete)
│   └── web/
│       ├── assets.go           (//go:embed templates + content + static → geliumui/site/web)
│       ├── content/            (44 markdown docs)
│       ├── static/             (app.css compilado, gelium.js copiado, app.js chrome, search.js)
│       ├── styles/             (app.css entry, docs-shell.css, docs-chrome.css, +8 site css)
│       └── templates/          (15 shell/landing/recipes templates)
├── internal/app/               (handler Go completo: docs server)
├── scripts/                    (copy-htmx, copy-icons, copy-lib-js, build-lib-dist)
├── go.mod                      (module geliumui — queda en la raíz)
└── package.json                (workspaces: ["lib", "site"])
```

### Separación librería vs docs

- `lib/` = el paquete reutilizable: componentes (`ui-*`), tokens, themes, consumer JS y dist prebuilt. Un consumidor instala esto.
- `site/` = el primer consumidor: shell de docs, landing, recipes, blog. Importa la librería POR NOMBRE DE PAQUETE (`@import "gelium-ui/styles/index.css"`, `gelium-ui/themes/*.css`) — el mismo contrato que un consumidor externo (dogfooding).
- `internal/app` = el handler Go que sirve el site (embed dual: `geliumui/site/web` + `geliumui/lib`).

### Nota sobre `embed`

Go no permite que una directiva `//go:embed` lea rutas padre con `..`. Por eso hay DOS embeds: `lib/assets.go` (templates + styles de la librería) y `site/web/assets.go` (templates + content + static del site). `buildTemplates()` mergea ambos con `template.ParseFS` — los nombres son disjuntos por construcción (collision-guard test).

### Cache-bust

La versión de cache-bust de assets es UN SOLO valor: `lib.AssetsVersion` (0.5.3), renderizado como `?v={{.AssetsVersion}}` en los templates. El test de coherencia pincha que la versión npm == constant y prohíbe `?v=` hardcodeado. Bumpéalo cuando cambie un asset estático.

## Ejecutar

```bash
go run ./cmd/gelium        # sirve en :8787 (PORT=3000 para otro puerto)
```

El sitio de documentación también se sirve desde tu propia aplicación usando `internal/app.New()`:

```go
package main

import (
	"log"
	"net/http"

	"geliumui/internal/app"
)

func main() {
	log.Fatal(http.ListenAndServe(":8787", app.New()))
}
```

## Wire contract (`gelium:*` / `X-Gelium-*`)

Los contratos wire server-driven usan el prefijo del producto: `gelium:toast`
(evento HX-Trigger + `#gelium-toast-region`) y el header `X-Gelium-Validation: true`.
El hook HTMX servido (`lib/js/gelium.js`) y los tests los fijan como canónicos.
Referencia canónica: [`docs/gelium-ui-wire-compatibility.md`](docs/gelium-ui-wire-compatibility.md).

## Alcance

El alcance actual incluye 47 componentes open-code (Button, Text field, Dialog, Toast, Data table, Chips, Navigation, …) con demos server-rendered, dos themes (`theme-material`, `theme-basecoat`) y un form HTMX de validación server-side. No incorpora Lit, Shadow DOM, Astro, `templ`, SQLite, SSE, CLI ni un registry de terceros.

## Licencia

Gelium UI se distribuye bajo la [licencia MIT](LICENSE).
