# REF-SHELL

Use for header / navbar / footer / account / search decisions on any surface.
Screenshots and external products inform hierarchy; the **consumer’s existing
shared shell is authoritative**.

## Audit lenses

- **Identity cluster:** `Me`, account menu, settings, logout — one coherent place.
  Do not re-implement account actions in the page body when the shell owns them.
- **Notifications:** separate destination only when it is a distinct job.
- **Search:** shell-owned vs page-owned; do not duplicate both without a reason.
- **Responsive:** what collapses, what becomes a menu, what stays visible.
- **No duplicate chrome:** one top nav system; reject second header/rail/bottom bar
  from a reference unless the product explicitly has no shell yet.

## Gelium filter

- Native links and forms; no-JS completion
- Real registered routes only
- CSRF on logout and other POST shell actions
- Accessible names; catalog icons only when they improve scan
- Never invent avatars, brand marks, or a parallel navigation system
- Page body links must not restate shell destinations the header already exposes
  (Architect checks shared header before proposing body chrome)

## Packet note

```text
References: REF-SHELL — preserve existing top nav; account stays in shell.
Rejected: reference bottom nav / left rail (shell already settled).
```
