# DESIGN.md template (consumer apps)

Copy to your application repo root as `DESIGN.md`. Records **visual lane** without forking Gelium’s token system.

```markdown
# Design

## Surface mix
List major URLs and mode (Operate | Read | Persuade | Experience):
- / → Persuade
- /app/* → Operate
- /help/* → Read

## Theme
- html class: theme-material | theme-basecoat
- Dark: theme-dark class route only (not media-only as sole authority)

## Density
- Operate surfaces: cozy | compact
- Read/Persuade: comfortable | cozy

## Components
Prefer Gelium `ui-*` partials from `gelium-ui/templates`. New primitives need a token-first rationale.

## Motion
Default MOTION-NONE on Operate. Honor prefers-reduced-motion. See Gelium /docs/motion.

## Anti-slop (project)
Add project-specific bans. Keep Gelium anti-slop from /docs/agent-workflow.

## Brand notes
Logo, any extra brand color **mapped into tokens** — do not hardcode one-off hex in features.
```

Gelium docs: `/docs/themes`, `/docs/tokens`, `/docs/density`, `/docs/agent-workflow`.
