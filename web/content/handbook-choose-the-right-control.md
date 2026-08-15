# Choose the right control

Every input decision starts from one question: **how does the user's choice behave?** Answer that first, then pick the component that matches the behavior — never the other way around.

## The decision table

| Situation | Component | Why |
|---|---|---|
| User picks **one** option from 2–5, all options should stay visible and comparable | [Radio](/components/radio) | All choices on screen; native keyboard arrow selection |
| User picks **one** option from a long list, or options carry a lot of text | [Select](/components/select) | Collapsed; the list is browsed, not scanned |
| User picks **many** independent options | [Checkbox](/components/checkbox) | Each choice is independent and toggles alone |
| A **single on/off setting** (one switch, not a group) | [Switch](/components/switch) | Immediate state change for one setting |
| A **continuous or approximate value** (volume, range, percentage) | [Slider](/components/slider) | Direct manipulation of a magnitude |
| A **precise value** (price, date, exact number) | [Text field](/components/text-field) | Keyboard entry is more accurate than dragging |
| A **short free-form answer** | [Text field](/components/text-field) | Text input |
| A **long free-form answer** | Text area (multi-line text field) | Paragraph input |
| A **command or action** on the page | [Button](/components/button) | Primary action control |
| A **choice between views or modes** | [Tabs](/components/tabs), [Segmented buttons](/components/segmented-button) | Persistent view switcher vs local mode toggle |

## Rules of thumb

- **Radio vs Select**: 5 or fewer options → Radio (visible, comparable). 6 or more, or options with long labels → Select.
- **Radio vs Checkbox**: one answer from the group → Radio. Multiple independent answers → Checkbox. If you catch yourself saying "select all that apply", it is Checkbox, not Radio.
- **Switch vs Checkbox**: Switch for ONE independent setting that takes effect immediately. Checkbox when the choice is part of a group of related options or needs a submit to apply.
- **Slider vs Text field**: Slider for approximate values where precision does not matter. Never use a Slider for exact critical values (price, date, account number) — the user cannot aim precisely.
- **Radio vs Select vs Menu**: Radio and Select are *data entry* (a form answer). Menu is *command selection* (choose an action). If the value is submitted with the form, it is Radio/Select; if choosing immediately performs an action, it is Menu.
- **Prefill, don't default**: a Radio group should start with a deliberate selection or explicit "None", not with the first option silently chosen.

## What this page is not

This page decides *which control* for *which situation*. It does not document anatomy, variants, or server contracts — those live on each component page. If a component's usage section contradicts this table, the component page wins for its own scope; tell us about the contradiction in the [Information architecture](/docs/information-architecture) audit instead of silently splitting.
