# UI definition of done

Use this checklist before accepting UI from a human or an LLM. It does **not** replace user research; it raises the floor so deliveries match Gelium’s sourced criteria.

Pass = every applicable item is true. Cite page IDs (`FEED-*`, `DATA-*`, `JOURNEY-*`) in the PR or agent summary.

## DoD checklist

### Product framing

1. [ ] **User job** stated in one sentence.  
2. [ ] **Screen type** chosen ([Screens](/docs/screens)).  
3. [ ] **Journey shape** named if >1 step ([Journeys](/docs/journeys)).  
4. [ ] **One primary action** visible per view.

### Structure

5. [ ] Hierarchy: H1 → context → primary → main → secondary.  
6. [ ] Nav pattern justified ([Screens](/docs/screens) nav table).  
7. [ ] Collection uses a `DATA-*` pattern ([Data display](/docs/data-display)).  
8. [ ] Density mode named for the surface ([Density](/docs/density)).

### States and feedback

9. [ ] Loading / empty / error defined for every collection (`FEED-LOAD-*`, `FEED-EMPTY`, `FEED-LOAD-FAIL`).  
10. [ ] Forms use summary + inline errors on validation (`FEED-VAL`), not toast-only.  
11. [ ] Success landing is correct (`FEED-OK-PAGE` vs `FEED-OK-TOAST`).  
12. [ ] Destructive flows use confirm (`FEED-CONFIRM`).

### Implementation contracts

13. [ ] Controls follow [Forms](/docs/forms) + [Choose the right control](/docs/choose-the-right-control).  
14. [ ] Server: GET read, POST+303 mutate, 422 + `X-Gelium-Validation` when needed.  
15. [ ] Theme via `theme-*` classes; no one-off hex for core chrome.  
16. [ ] Narrow viewport: stack, `min-width: 0`, no body `overflow-x: hidden` ([Responsive](/docs/responsive)).  
17. [ ] Motion is NONE or justified; `prefers-reduced-motion` honored ([Motion](/docs/motion)).  
18. [ ] Copy follows [Content style](/docs/content-style) (errors say the fix).

### Agent extras

19. [ ] `/llms-ux.txt` rules followed (FEED/DATA/JOURNEY ids in notes).  
20. [ ] Domain skeleton used if applicable ([Patterns](/docs/patterns)) instead of inventing a new IA.

## Failure handling

If any item fails: **fix before shipping**. Do not “pass with toast” on validation or skip empty states.

## See also

- [`/llms-ux.txt`](/llms-ux.txt) · [Screens](/docs/screens) · [Feedback](/docs/feedback)
