# Data display

Choosing **how to show a collection** (table vs list vs cards) is a product decision, not a decoration choice. This page gives when/when-not rules adapted from government design systems and common data-UI practice, mapped to Gelium primitives.

## Sources

| Topic | Source |
|---|---|
| Tables as data | [GOV.UK: table](https://design-system.service.gov.uk/components/table/) (tabular data; not for layout) |
| Side / section structure | [USWDS: side navigation](https://designsystem.digital.gov/components/side-navigation/) (IA, not grids) |
| Lists and process | [USWDS: process list](https://designsystem.digital.gov/components/process-list/), list patterns in USWDS components |
| Scanning | [NNG: F-shaped pattern](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/) |
| Collection states | [Feedback](/docs/feedback) (`FEED-EMPTY`, `FEED-LOAD-LIST`, `FEED-ROW`) |
| Narrow tables | [Responsive](/docs/responsive) — scroll **inside** the table region |

## Collection patterns

| ID | Pattern | When to use | When not | Gelium |
|---|---|---|---|---|
| **DATA-TABLE** | Data table | Compare **many attributes** across rows; sort/scan columns; admin density | Layout of a marketing page; 2 fields that fit a list | `data-table` + local scroll wrapper; recipe [Admin Resource](/recipes/admin-resource) |
| **DATA-LIST** | Linear list | Primary label + few meta lines; clear row action | Wide numeric comparison across 8 columns | `list` (and row actions) |
| **DATA-CARDS** | Card grid/stack | Browse entities as objects (image, title, short pitch); uneven content height | Dense reconciliation of numbers/status across dozens of fields | `card` / `feature-card` + stack or simple grid; stack on narrow |
| **DATA-FEED** | Activity feed | Time-ordered events, social/activity | Master data admin | Recipe [Public feed](/recipes/public-feed) |
| **DATA-QUEUE** | Work queue | Items to process; status + action dominate | Read-only catalog | Recipe [Ops Queue](/recipes/ops-queue) |
| **DATA-DESC** | Description list / detail fields | One record’s attributes | Hundreds of peer records | Detail screen: rows of label/value, `divider` |

GOV.UK: use tables for **tabular data**, not to position non-tabular content ([table](https://design-system.service.gov.uk/components/table/)).

## Filters, sort, pagination

| Concern | Do | Don’t |
|---|---|---|
| Filters | Put above the collection; use [Forms](/docs/forms) controls | Hide all filters in unlabeled icons only |
| Default sort | Document the default (e.g. newest first) | Surprise reorder with no affordance |
| Pagination | For large sets; keep filters in query/GET | Infinite scroll as the only path when users need item N |
| URL state | Filters/sort/page on **GET** query | POST-only filter that breaks share/back |

## Columns and density

| Rule | Guidance |
|---|---|
| Column count | Prefer **fewer columns** on first view; secondary data on detail |
| Actions column | One obvious row primary; rest in menu or detail |
| Admin density | Compact rows OK; still meet `--ui-touch-target` on controls |
| Consumer browse | More whitespace; cards/list over dense tables |
| Numbers | Align consistently; don’t mix units in one column without labels |

## States (mandatory)

Every collection view must define:

1. **Loading** → `FEED-LOAD-LIST`  
2. **Empty** → `FEED-EMPTY`  
3. **Error** → `FEED-LOAD-FAIL`  
4. **Row failure** (if row actions) → `FEED-ROW`  
5. **Partial batch** → `FEED-PARTIAL`

## Anti-patterns

- Table used as a 3-column layout for prose cards.
- Cards for a 12-column financial compare.
- Empty `<table>` with no `empty-state`.
- Horizontal page scroll because the table wasn’t contained ([Responsive](/docs/responsive)).
- Row errors only as a global toast (`FEED-ROW`).

## Checklist (agents)

1. Pick `DATA-*` id before choosing a component.
2. List columns/fields **and** which live only on detail.
3. Wire GET query for filter/sort/page if present.
4. Implement all five collection states.
5. Check ~360px: stack filters; table scrolls inside.

## See also

- [Screens](/docs/screens) · [Journeys](/docs/journeys) · [Patterns](/docs/patterns) · [`/llms-ux.txt`](/llms-ux.txt)
