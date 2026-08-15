# Information architecture

Information architecture (IA) is the rule that decides where a page lives in the docs shell. Gelium UI follows one ordering rule: **concept before reference** — onboarding pages come before lookup pages, so a reader learns how the system thinks before they hunt for a component's API.

## The hierarchy rule

The sidebar is a single ordered list of groups, and every group belongs to one of two kinds:

1. **Concept** — pages that explain how Gelium UI works: the Handbook (Themes, Tokens, Server contracts, Accessibility, Design principles, and this page).
2. **Reference** — pages that document what exists: component sections (Foundation, Actions, Input, Feedback & status, Navigation, Data), then Patterns and Recipes.

The canonical order is:

1. Getting started
2. Handbook (concept)
3. Component sections (reference)
4. Patterns and Recipes (composition)

Rationale: **onboarding before lookup**. A new reader lands on the [Documentation](/docs) hub, reads the Handbook, and only then looks up a component. Reference pages assume the concepts; concepts never assume a component exists. If a Handbook page links a component, that link is a shortcut, not a dependency — the component page never links back as its source of truth.

## When to add a group or page

Before adding a nav group or page, answer three questions:

1. **Who navigates?** Name the reader: a new user, a designer, an integrator, or a maintainer. If you cannot name one, the destination is probably not ready for the nav.
2. **What task?** State the task the reader is trying to finish on that page. A page earns its place by completing one task — not by existing.
3. **Concept or reference?** Run the test: does the page explain *how Gelium UI works* (concept → Handbook) or does it document a specific primitive (reference → component section)? A page that mixes both splits: the concept stays in the Handbook, the API detail stays on the component page.

Rules of thumb:

- A new component page goes in its `docsSections` category (Foundation, Actions, Input, Feedback & status, Navigation, Data), never in the Handbook.
- A new Handbook page explains a cross-cutting concern (themes, tokens, server contracts, accessibility, principles, IA) and goes in the Handbook group, before the component sections.
- Patterns and Recipes stay below components: they compose primitives, so they are read after the primitives are known.

## Agent prompt

The following prompt lets an LLM evaluate or improve the docs IA using Gelium's own nav model (`docsNavFor` in `internal/app/docs.go`). Paste it into any agent:

```text
You are auditing the Gelium UI docs information architecture. The nav model is
docsNavFor in internal/app/docs.go; the canonical rule is "concept before
reference": groups render as Getting started → Handbook → component sections
(Foundation, Actions, Input, Feedback & status, Navigation, Data) → Patterns →
Recipes.

Evaluate the current nav (sidebar groups, /docs hub sections, footer sections,
sitemap) against these criteria:

1. ORDER — every concept group (Handbook) must appear before every reference
   group (components, Patterns, Recipes). Flag any page whose group sits below
   a group that depends on it.
2. PLACEMENT — for each nav destination, run the concept-or-reference test:
   does the page explain how Gelium UI works (concept → Handbook) or document
   a specific primitive (reference → component section)?
3. TASK — each destination must complete one reader task. Flag pages whose
   title is not a task a reader would search for.
4. DRIFT — the nav, /docs hub, footer, and sitemap must be derived from the
   same model (handbookNavLinks, docsSections, componentRoutes). Flag any
   destination that exists in only one of them.

Report: (a) a verdict per criterion, (b) a list of concrete reorder or
re-home proposals with the exact group name each destination should move to,
(c) anything you would add to the Handbook as a concept page. Never propose
moving a component page into the Handbook, and never invent a destination
that has no route.
```

## How it is enforced

Tests pin the rule: `TestHandbookGroupPrecedesComponentSections` asserts the Handbook group sits at index 1 in the nav model, in the rendered sidebar, and in the /docs hub. The nav model, hub, footer, and sitemap share one source of truth (`handbookNavLinks`), so a new Handbook page appears everywhere or nowhere.
