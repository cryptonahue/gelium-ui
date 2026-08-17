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

### States and feedback

10. [ ] Loading / empty / error defined for every collection (`FEED-LOAD-*`, `FEED-EMPTY`, `FEED-LOAD-FAIL`).  
11. [ ] Forms use summary + inline errors on validation (`FEED-VAL`), not toast-only.  
12. [ ] Success landing is correct (`FEED-OK-PAGE` vs `FEED-OK-TOAST`).  
13. [ ] Destructive flows use confirm (`FEED-CONFIRM`).

### Implementation contracts

14. [ ] Controls follow [Forms](/docs/forms) + [Choose the right control](/docs/choose-the-right-control).  
15. [ ] Server: GET read, POST+303 mutate, 422 + `X-Gelium-Validation` when needed.  
16. [ ] Theme via `theme-*` classes; no one-off hex for core chrome.  
17. [ ] Narrow viewport: stack, `min-width: 0`, no body `overflow-x: hidden` ([Responsive](/docs/responsive)).  
18. [ ] Motion is NONE or justified; `prefers-reduced-motion` honored ([Motion](/docs/motion)).  
19. [ ] Copy follows [Content style](/docs/content-style) (errors say the fix).

### Agent extras

20. [ ] Workflow passes: shape → build → audit → polish (no endless redesign loops).  
21. [ ] `/llms-ux.txt` rules followed (FEED/DATA/JOURNEY/WF ids in notes).  
22. [ ] Domain skeleton used if applicable ([Patterns](/docs/patterns)).  
23. [ ] `bash scripts/ux-detect.sh` clean when working in this monorepo (or consumer equivalent).  
24. [ ] Anti-slop list reviewed ([Agent workflow](/docs/agent-workflow)) without breaking the theme.

## Failure handling

If any item fails: **fix before shipping**. Do not “pass with toast” on validation or skip empty states.

## See also

- [`/llms-ux.txt`](/llms-ux.txt) · [Screens](/docs/screens) · [Feedback](/docs/feedback)
