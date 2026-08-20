# Official package evidence (local)

Downloaded for token translation. Do not import these packages into the production
Gelium CSS bundle as-is.

## Basecoat

- Package: `basecoat-css@1.0.2`
- File: `basecoat-css-1.0.2.tgz`
- Button rules extracted from the published CDN bundle inside the tarball:
  `package/dist/basecoat.cdn.css` (`.btn` selectors)
- Snapshot: `basecoat-btn-extracted.css`

Basecoat component CSS in the npm package is mostly Tailwind `@apply` stubs;
the compiled CDN bundle is the authoritative published anatomy for translation.

## Base UI

- Package: `@base-ui/react@1.7.0`
- File: `base-ui-react-1.7.0.tgz`
- The package is headless and ships **no component visual CSS**.
- Docs-inspired Button demo CSS is shipped inside the package docs markdown:
  `package/docs/react/components/button.md` (CSS Modules demo)
- Snapshot: `baseui-button-hero.module.css`

Always label Gelium Base UI visuals as **docs-inspired / Gelium-authored**, never
as official Base UI package CSS.
