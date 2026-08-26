# Skill: Usability checklist (per-screen)

Run this over EVERY surface before DoD (`skills/07`). Binary questions only —
any NO is a defect to fix or escalate to the user. Sources: Gelium UX
principles (13), GOV.UK service patterns, Nielsen heuristics, cognitive
accessibility.

## Comprehension

- [ ] Could a first-time user grasp this screen's purpose in 5 seconds?
- [ ] Does every control look like what it does (affordance matches behavior)?
- [ ] Is navigation visible and persistent — never hidden behind hover,
      gesture, or JS-only reveals?
- [ ] Is current location explicit (`aria-current`, breadcrumb, active tab)?
- [ ] Is state recognizable from the URL/screen — no recall-only flows?
- [ ] Is essential info visible on the page, never tooltip-only?

## Feedback

- [ ] After every action, does something say what happened and what to do next?
- [ ] Is loading visible while work runs (progress/skeleton/`aria-busy`)?
- [ ] Is persistent feedback in a persistent channel (inline alert/banner),
      NOT a toast? Toasts only for transient events.
- [ ] Are status changes announced (`role="status"` polite / `role="alert"`
      assertive — correct one used)?

## Errors

- [ ] Can this step fail without losing user input (422 re-render preserves
      values)?
- [ ] Does each error message say WHAT happened + HOW to fix it — plain
      language, no blame ("Enter a date", not "Invalid input")?
- [ ] Does a form-level validation summary link to each failing field?
- [ ] Are destructive actions confirmed and separated from primary actions?
- [ ] Is there an undo or recovery path where the action is consequential?

## Forms & inputs

- [ ] Every input has a visible, programmatic label (no placeholder-as-label)?
- [ ] Is the right control chosen for the choice's behavior (closed list →
      select/radio, not free text)?
- [ ] Are constraints prevented up front (`required`, type, pattern, helper
      text) rather than corrected after submit?
- [ ] One clear primary submit; secondary actions visually quieter?
- [ ] Errors are per-field (`aria-invalid` + `aria-describedby`) AND
      summarized at top?

## Content & copy

- [ ] Plain English, active voice, short sentences — no jargon, no "please"?
- [ ] Button labels name the action ("Save changes", never "Submit"/"OK")?
- [ ] Empty states explain what/why + give a CTA — never a bare void?
- [ ] Same concept = same word everywhere (internal consistency)?

## Touch & pointer

- [ ] All targets ≥ 44px with spacing between them?
- [ ] No hover-only interaction (touch and keyboard can reach everything)?
- [ ] Works at narrow viewport: fluid layout, contained horizontal scroll,
      no `overflow-x: hidden` on body masking loss?
- [ ] Motion explains change only; honors `prefers-reduced-motion`?

## Trust

- [ ] Does the screen do exactly what it says — no surprise side effects?
- [ ] Are permissions/costs/irreversibility stated BEFORE the commit?
- [ ] Does the flow survive closing the tab mid-way (URL state, resumable)?
- [ ] Consistent with the rest of the app (same component, same token, same
      word for the same job)?

## Compact final pass

Purpose in 5s · affordances honest · status always visible · errors recoverable
without data loss · labels real · copy plain · touch-safe · trustworthy.
Any NO → fix or ask the user before calling the surface done.

Then run `skills/07-dod-and-antislop.md`.
