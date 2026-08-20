# lib/fonts — self-hosted theme webfonts

This directory holds the **self-hosted, subsetted WOFF2 webfonts** that Gelium
themes ship. They travel with the library (`lib.Assets` embeds `fonts/*`) and
are served by the site at `/static/fonts/<file>` — never from a CDN.

## Layout

Flat, globally idempotent filenames following the `<theme-class>-<file>`
convention so preloaded URLs never collide across themes:

```
theme-material-alden.woff2     # example — the theme font file a theme preloads
```

The font set a theme ships is declared in `themeDirection.Fonts`
(`internal/app/server.go`); only fonts referenced by an allowlisted theme are
servable (the `/static/fonts/` namespace is closed).

## Why self-host + WOFF2

- **No third-party DNS / connection**: fonts resolve on your own origin, so no
  separate connection setup delays First Paint / LCP.
- **WOFF2** already carries Brotli compression.
- **Subsetting** (Unicode ranges only for the glyphs a theme needs) keeps each
  weight small (~tens of KB).

See `docs/gelium-ui-font-contract.md` for the full quality contract every font
here must satisfy (licensing, subsetting, `font-display: swap`, metric
overrides, contrast, language coverage).
