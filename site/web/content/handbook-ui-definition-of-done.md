# UI definition of done

Use this checklist before accepting UI from a human or an LLM. It does **not** replace user research; it raises the floor so deliveries match Gelium’s sourced criteria.

Pass = every applicable item is true. Cite page IDs (`FEED-*`, `DATA-*`, `JOURNEY-*`) in the PR or agent summary.

## DoD checklist

### Product framing

1. [ ] **User job** stated in one sentence.  
2. [ ] **Surface mode** Operate | Read | Persuade | Experience ([Agent workflow](/docs/agent-workflow)).  
3. [ ] **Screen type** chosen ([Screens](/docs/screens)).  
4. [ ] **Journey shape** named if >1 step ([Journeys](/docs/journeys)).  
5. [ ] **One primary action** visible per view.

### Structure

6. [ ] Hierarchy: H1 → context → primary → main → secondary.  
7. [ ] Nav pattern justified ([Screens](/docs/screens) nav table).  
8. [ ] Collection uses a `DATA-*` pattern ([Data display](/docs/data-display)).  
9. [ ] Density mode named for the surface ([Density](/docs/density)).

### Page and section architecture

10. [ ] Every major region has a distinct purpose and a `SECTION-CONTRACT` before components or tokens are selected ([Page + section architecture](/docs/page-section-architecture)).
11. [ ] DOM order is entry → primary → supporting → recovery; repeated card anatomy is not mistaken for page architecture.
12. [ ] Section actions are subordinate unless the section owns the one declared page-primary action.
13. [ ] Each applicable data/action section has recoverable loading, empty, error, and success behavior with a no-JS path.

### States and feedback

14. [ ] Loading / empty / error defined for every collection (`FEED-LOAD-*`, `FEED-EMPTY`, `FEED-LOAD-FAIL`).
15. [ ] Forms use summary + inline errors on validation (`FEED-VAL`), not toast-only.
16. [ ] Success landing is correct (`FEED-OK-PAGE` vs `FEED-OK-TOAST`).
17. [ ] Destructive flows use confirm (`FEED-CONFIRM`).

### Implementation contracts

18. [ ] Controls follow [Forms](/docs/forms) + [Choose the right control](/docs/choose-the-right-control).
19. [ ] Server: GET read, POST+303 mutate, 422 + `X-Gelium-Validation` when needed.
20. [ ] Theme via `theme-*` classes; no one-off hex for core chrome.
21. [ ] Narrow viewport: stack, `min-width: 0`, no body `overflow-x: hidden` ([Responsive](/docs/responsive)).
22. [ ] Motion is NONE or justified; `prefers-reduced-motion` honored ([Motion](/docs/motion)).
23. [ ] Copy follows [Content style](/docs/content-style) (errors say the fix).

### Agent extras

24. [ ] Workflow passes: shape → architecture → build → audit → polish (no endless redesign loops).
25. [ ] `/llms-ux.txt` rules followed (FEED/DATA/JOURNEY/WF/ARCH ids in notes).
26. [ ] Domain skeleton used if applicable ([Patterns](/docs/patterns)).
27. [ ] `bash scripts/ux-detect.sh` clean when working in this monorepo (or consumer equivalent).
28. [ ] Anti-slop list reviewed ([Agent workflow](/docs/agent-workflow)) without breaking the theme.

## Failure handling

If any item fails: **fix before shipping**. Do not “pass with toast” on validation or skip empty states.

## See also

- [`/llms-ux.txt`](/llms-ux.txt) · [Screens](/docs/screens) · [Feedback](/docs/feedback)
