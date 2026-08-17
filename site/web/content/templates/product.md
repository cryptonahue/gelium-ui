# PRODUCT.md template (consumer apps)

Copy to your application repo root as `PRODUCT.md`. Agents should read it before shaping UI. This is **product context**, not Gelium internals.

```markdown
# Product

## Audience
Who uses this and in what situation?

## Primary jobs
1. …
2. …

## Non-goals
What we explicitly will not build in this slice.

## Voice
Tone for UI copy (see also Gelium Content style): plain / formal / …

## Constraints
- Stack: server-rendered HTML + gelium-ui (npm)
- Theme: theme-material | theme-basecoat
- JS: progressive only / HTMX optional
- A11y: keyboard, contrast, reduced-motion

## Anti-references
Looks or patterns we do **not** want (e.g. “generic purple SaaS landing”).

## Success metrics (optional)
What “good” means for this release.
```

Gelium docs: `/docs/agent-workflow`, `/llms-ux.txt`, `/docs/ui-definition-of-done`.
