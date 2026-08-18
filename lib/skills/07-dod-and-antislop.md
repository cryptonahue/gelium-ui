# Skill: Definition of Done and anti-slop

Run this checklist before calling a surface done. If any item is NO, fix it — do
not pass.

## UI Definition of Done (product/structure/states/contracts)

**Product**
- [ ] One primary action per screen.
- [ ] Each screen serves the user's job, not the org chart.
- [ ] Empty/error/loading/success all handled.

**Structure**
- [ ] Semantic HTML (native elements), logical heading order (H1 → H2 → H3).
- [ ] Correct list semantics for lists; tables only for tabular data.
- [ ] Links are real `<a href>`; actions are real buttons/forms.

**States**
- [ ] Empty state explains what/why/what-to-do, with a CTA.
- [ ] Error copy names the fix (action pattern), never blames the user.
- [ ] Loading is readable; content doesn't jump on load.

**Contracts**
- [ ] Read = GET + URL state.
- [ ] Mutate = POST + 303 (or 422 + `X-Gelium-Validation` re-render).
- [ ] No-JS flow completes end-to-end.
- [ ] Chat/feedback uses the FEED matrix, persistent ≠ toast.

## Anti-slop (Gelium-aware) — far side of the coin

Avoid:
- nested cards for simple forms;
- stock purple-blue hero gradients;
- an icon tile above every H1;
- gray on saturated fills (low contrast);
- bounce/babel motions;
- "bolder for its own sake" on Operate tables;
- a new font per feature;
- layout tables for non-data; "lorem" filler that hints at content.

OK:
- Material/Basecoat token stacks from the theme;
- brand fonts only where `DESIGN.md` says so;
- spacing from the core `--ui-space-*` scale (no feature-specific scale).

## Verification

Where the consumer repo mirrors this package's monorepo gates, run the same
checks (styles_*_test, copy/contrast contracts). In any consumer repo, at minimum:
grep for `overflow-x: hidden` on body (forbidden), confirm `theme-*` class on
`<html>`, confirm a validation-summary + inline errors on every form, and confirm
no one-off color literals in shipped markup.

## Do not

Skip to pretty CSS. Claim pass without the checks. Redesign under "polish". End
in an endless loop — one polish pass, then stop.
