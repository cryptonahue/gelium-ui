# Skill: State and feedback

Every surface needs the four states: **empty**, **error**, **loading**, **success**.
Pick the right FEED channel from the decision matrix.

## FEED decision matrix

| Channel | Use for | Not for |
|---|---|---|
| toast | transient, auto-dismissable, one result | persistent info or policy |
| banner (`ui-banner`) | page/site-level durable signal (consent, expiry, maintenance) | field errors |
| inline alert (`ui-inline-alert`) | inline field/group error next to content | global status |
| validation summary | form submit errors (grouped, with links) | single field already inline |
| empty state (`ui-empty-state`) | collection with zero results | a static "no items" text |
| error state (`ui-error-state`) | page-level load failure + retry | inline field error |
| skeleton (`ui-skeleton`) | initial load placeholder | updated nested data |
| callout | distinct explanatory note | persistent page banner |

Persistent success → banner/result page. Transient → toast.

## States checklist per surface

- **Empty** — explain what is here / why empty / what to do (CTA). Never blank.
- **Error** — recognize + recover in the action pattern ("Enter the project
  name." not "Name is required."), never blame the user.
- **Loading** — readable placeholder (skeleton) for initial load; keep the page
  scannable, don't block the whole UI.
- **Success** — confirm the outcome and give the next step.

## Copy rules

Sentences ≤ 25 words. Active voice, plain English, no "please", AP style.
Error copy names the fix, not just the failure.
