# Tarea OpenCode — Toast Loom UI

Leé completo:

1. `AI-COMPONENT-IMPLEMENTER-PROMPT.md`
2. `COMPONENT-ROADMAP.md`

Implementá exactamente el componente Toast (Wave 0 / P0) siguiendo el contrato del prompt implementer.

## Asignación

```text
COMPONENTE: Toast
SLUG: toast
CATEGORÍA: LOOM_ONLY
VARIANTES REQUERIDAS: success, error, info (y neutral si surge de auditoría)
ESTADOS REQUERIDOS: visible, dismissible, stacked/cola si la auditoría lo exige, reduced-motion
COMPORTAMIENTOS REQUERIDOS: región aria-live; dismiss visible; timeout/pausa hover-focus si la plataforma no alcanza sin JS; feedback no-JS inline o full-page
FLUJO SERVER/HTMX, SI APLICA: demo dogfoodeada; wire contract exacto tipo HX-Trigger {"loom:toast":{"type":"success","message":"..."}} si aplica; sin JS debe quedar feedback server-rendered
CRITERIOS ESPECÍFICOS DE ACEPTACIÓN: platform-first + Material Snackbar como referencia visual (no componente Snackbar separado); cero React/Lit/Shadow DOM; JS mínimo solo si gap demostrado; TDD RED→GREEN; docs dogfood /components/toast; no tocar :8787 ni loom.exe; smoke solo :8788
MATRIZ DE NAVEGADORES SOPORTADA: Chromium/Edge/Firefox/Safari actuales del entorno local (Baseline 2025+ donde aplique); documentar fallbacks
FECHA/SNAPSHOT BASELINE WEB: 2026-08-08
UPSTREAM SNAPSHOT ID PROVISTO POR COORDINADOR: material-web-upstream local en D:\repos\material-web-upstream (árbol local de referencia; no uses Git)
UPSTREAM MANIFEST/HASH PROVISTO: usar MATERIAL-WEB-PROGRESS.md + inspección de rutas literales locales
MODO ASIGNADO: EXCLUSIVE_INTEGRATION
WORKSPACE AUTORIZADO: D:\repos\loom-ui
WORKER ID: opencode-toast-01
RESERVA ID: toast-exclusive-2026-08-08
PATHS POSEÍDOS: web/templates/toast.html (nuevo), web/content/toast.md (nuevo), y con reserva exclusiva los shared necesarios: internal/app/server.go, internal/app/server_test.go, web/templates/layout.html, web/styles/app.css, themes/theme-material/theme.css, web/styles_contract_test.go, web/static/app.css, web/static/htmx.min.js, web/static/app.js (solo si gap JS demostrado), README.md, package.json (version bump si hace falta)
PATHS PROHIBIDOS ADICIONALES: MATERIAL-WEB-PROGRESS.md, PROMPT-MATERIAL-WEB-INVENTORY.md, COMPONENT-ROADMAP.md, AI-COMPONENT-IMPLEMENTER-PROMPT.md, loom.exe, repos material-*
BASELINE SHA-256: capturá SHA-256 de cada archivo shared antes de editar; abortá ABORTED_ON_DRIFT si cambia sin tu escritura
NUEVA VERSIÓN DE ASSETS PROPUESTA: 0.4.0
```

## Reglas clave

- No Git.
- No React/Lit/Shadow DOM/Astro/templ/CDN.
- Auditar primero HTML/CSS moderno (aria-live, popover, etc.) y Snackbar Material como referencia visual.
- TDD estricto: test RED observado antes de producción.
- Spec review PASS luego quality APPROVED.
- Smoke solo PORT=8788; no tocar :8787 ni loom.exe.
- Entregá el formato del prompt con estado final permitido y checklist de aceptación del usuario.
- Windows Bash/MSYS: sintaxis POSIX.

Empezá leyendo los dos documentos de contrato y después ejecutá el slice Toast de punta a punta.
