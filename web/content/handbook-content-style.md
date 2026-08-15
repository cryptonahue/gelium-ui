# Content style

Copy is part of the component contract, not decoration. Every string a Gelium UI screen renders — errors, toasts, empty states, banners, validation summaries, and the docs themselves — follows one voice: plain English, active voice, short sentences, no jargon, no "please" (GOV.UK content design), and language clear to anyone (Material 3 content design). Editorial mechanics — numbers, dates, capitalization — follow AP Style.

## Error messages

An error message does two jobs: **say what happened** and **say how to fix it** (NNG error-message guidelines, heuristic 9). Lead with the fix in the user's words. Never blame the user, and never leave a consequence alone.

- BAD "Name is required." — states the consequence; the fix is implied.
- GOOD "Enter the project name." — the fix is the message.
- BAD "Invalid input." — names nothing and blames the input.
- GOOD "Choose a status from the list." — names the field and the fix.

Field errors render next to the field that caused them (`aria-invalid` + `aria-describedby`), and the validation summary links each one back to its field. A field error is never reported only in a toast — validation uses the 422 contract (see [Server contracts](/docs/server-contracts)).

## Toasts and transient feedback

A toast states the **verb and the result** of the action that just completed, and carries one action. A static label like "Success" or "Error" is not a message — it says nothing about what happened.

- GOOD "Saved", "Project created", "Could not save — try again".
- BAD "Success", "Error", "Operation completed".

## Empty states

An empty state answers three questions in order: **what is here**, **why it is empty**, and **what to do next**. When a clear next action exists, make it a real control, not prose.

- "No projects yet — Create your first project to get started." (create CTA)
- "No projects match your search. Try adjusting the filters." (second-step action)
- "No queue items match the selected filters. Clear them to see the full queue." (clear-filters CTA)

## Banners and validation summary

Banners and validation summaries use the same voice as errors: what happened, then the fix. The validation summary lists each field error as a link to the field, never a paragraph of prose.

## Docs content

Docs titles are **tasks, not topics**: "Choose the right control", not "Controls". Guidance comes before reference, per the [Information architecture](/docs/information-architecture) rule. Keep sentences to 20 words or fewer, and avoid filler words: "just", "very", "obviously".

## How it is enforced

The copy contract tests in `internal/app/copy_contract_test.go` pin the action pattern on recipe errors and the action language in recipe empty states. The Handbook tests pin this page in the sidebar, the /docs hub, and the sitemap. The [Screen recipes](/recipes/admin-resource) dogfood the voice on real screens.
