# Content style

Copy is part of the component contract, not decoration. Every string a Gelium UI screen renders — errors, toasts, empty states, banners, validation summaries, and the docs themselves — follows one voice: plain English, active voice, short sentences, no jargon, no "please" (GOV.UK content design), and language clear to anyone (Material 3 content design). Editorial mechanics — numbers, dates, capitalization — follow AP Style.

## Reading on screen

Most users **scan** pages instead of reading them, and that has been true since eye-tracking studies in the 1990s (NNG); only a minority reads linearly. Scanning follows the **F-pattern**: two horizontal bands — the title and the first lines — then a vertical band on the left as the reader drops down looking for anchors. So the important content goes top-left, never buried mid-paragraph, and every paragraph leads with its point. Web copy is shorter than print: people look for something specific, not pleasure reading (NNG).

## Paragraphs and sentences

- **Paragraphs** — 2-4 sentences, roughly 40-70 words, one topic each.
- **Sentences** — at most 25 words, ideally 20; split any sentence that can be split.
- **Lead paragraph** — inverted pyramid: the conclusion comes first.
- **Lists** — prefer bullets over prose for three or more parallel items.
- **Voice** — active voice, always; drop filler such as "please" and "you should".

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

## Content structure (headings and blocks)

How to choose **H1–H3**, lists, tables, and quotes. Same rules for handbook pages and for long-form Read surfaces. Product **Operate** screens still follow [Screens](/docs/screens) (one H1, few decorative headings).

### Sources

| Topic | Source |
|---|---|
| Scan / F-pattern | [NNG: F-shaped pattern](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/) (see Reading on screen above) |
| Clear structure | GOV.UK content design — short pages, meaningful headings |
| Heading order | WCAG: headings describe topic or purpose; do not skip levels for looks |

### Headings

| Level | When to use | When not |
|---|---|---|
| **H1** | Exactly **one** per URL: the page job or title | Mid-body H1; multiple H1s |
| **H2** | A section you could put on “On this page”; a real topic change | Decorating a single paragraph; skipping from H1 to H3 |
| **H3** | A labeled part **inside** an H2 (sub-rules, “Toast rules”, steps) | Fake H2 because the outline looks empty |
| **H4+** | Almost never in Gelium docs | Deep trees |

**Test:** if the section would not appear in an on-this-page list, it is probably not an H2.

### Canonical handbook outline

```text
H1  Task-shaped title
    Lead paragraph (conclusion first)
H2  Sources (when criteria are sourced)
H2  Main rules or decision matrix
  H3  Nested rule groups only when needed
H2  When not / anti-patterns
H2  Checklist (agents)
H2  See also
```

### Lists, tables, prose, quotes

| Form | When | When not |
|---|---|---|
| **Bullet / numbered list** | Three or more **parallel** items | Two items that need a full sentence of nuance each |
| **Prosa** | Cause/effect, one idea with lead sentence | Long walls of parallel facts (use a list) |
| **Table** | Compare dimensions (when / when not / ID → use) | Layout or “pretty” columns of prose |
| **Blockquote** | Short **external** quote or cited rule | Emphasis, slogans, or “design flair” |
| **Callout / banner (UI)** | Page or system status ([Feedback](/docs/feedback)) | Replacing a heading |
| **Code** | Commands, markup, contracts | Theory and narrative |

### Product UI (Operate) vs docs (Read)

| Surface | Headings | Blocks |
|---|---|---|
| **Operate** (forms, tables, queues) | One H1; section labels via structure, not essay H2 stacks | Lists in empty/error copy; FEED matrices live in docs, not on every screen |
| **Read** (handbook) | Full H1–H3 outline above | Tables for criteria; lists for parallels |
| **Persuade** (landing) | One H1; section H2 sparingly | One primary CTA ([Screens](/docs/screens)) |

### Anti-patterns

- H1 → H3 with no H2 (skipped level for style).
- H2 for every sentence.
- Blockquote used as a highlight box without a source.
- Table used to position non-tabular marketing copy.
- Operate form with a marketing H2 stack above the fields.

## How it is enforced

The copy contract tests in `internal/app/copy_contract_test.go` pin the action pattern on recipe errors and the action language in recipe empty states. The same file carries the sentence-length contract: every sentence in every component page is at most 25 words (code and table rows are stripped before counting), per §Paragraphs and sentences. Structure rules above are pinned by handbook/content-style tests and the agent pack (`DOC-H1`, `DOC-H2`, `DOC-LIST`). The Handbook tests pin this page in the sidebar. Recipes link criteria via the on-page “Maps to Gelium criteria” bridge. The [Screen recipes](/recipes/admin-resource) dogfood the voice on real screens.
