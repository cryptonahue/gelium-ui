# Language switcher

The language switcher is a GET navigation form for changing the document language. Choosing a language and pressing the visible submit navigates to `?lang=<code>`; the server answers with a 303 redirect to the localized URL. It is deliberately not a POST mutation and never auto-submits from script: the control works with zero JavaScript, exactly like every other Gelium UI pattern.

## Examples

A labeled switcher with three locales, Spanish currently selected.

<div class="specimen-block">
<form class="ui-language-switcher" method="get" action="/components/language-switcher">
  <label class="ui-language-switcher-label" for="lang-docs">Language</label>
  <div class="ui-language-switcher-control">
    <select class="ui-language-switcher-select" id="lang-docs" name="lang">
      <option value="en">English</option>
      <option value="es" selected>Español</option>
      <option value="pt">Português</option>
    </select>
    <button class="ui-button ui-button-secondary" type="submit"><span>Switch</span></button>
  </div>
</form>
</div>

Without the label slot the control stands alone — keep an accessible name another way if you drop it (for example `aria-label` on the `<select>`).

<div class="specimen-block">
<form class="ui-language-switcher" method="get" action="/components/language-switcher">
  <div class="ui-language-switcher-control">
    <select class="ui-language-switcher-select" id="lang-footer" name="lang" aria-label="Language">
      <option value="en">English</option>
      <option value="es">Español</option>
    </select>
    <button class="ui-button ui-button-secondary" type="submit"><span>Switch</span></button>
  </div>
</form>
</div>

## Guidance

### When to use

Use a language switcher when one document is served in several languages and the user must be able to move between them at any time — typically in the footer or the header chrome. It belongs wherever locale is a property of the whole page, not of one fragment.

### When not to use

Do not use it to change display preferences such as theme or density — those are [theme](/docs/themes) concerns with their own switchers. Do not use it as a general select-and-submit filter: a data filter belongs in [Select](/components/select) inside its own form contract. And do not use it when the page ships in a single language — a one-option switcher is noise.

### Usability

- Keep the submit visible. Auto-submitting on `change` breaks keyboard review, screen-reader expectations, and no-JS users alike.
- Name locales in their own language (`Español`, not `Spanish`) so speakers can find theirs without reading the current language first.
- Mark the current locale with `selected`; the server renders it, the user never guesses where they are.
- Answer the GET with `303 See Other` so the localized URL becomes the address bar URL and refresh/reload stays correct.

### Accessibility

- The native `<select>` carries the semantics; pair it with a visible `<label>`, or give the control an `aria-label` when the label slot is omitted.
- The focus ring uses `--ui-focus-*` tokens with `:focus-visible`, so keyboard focus is visible and mouse clicks stay clean.
- In forced-colors mode the label paints `CanvasText` and the select keeps a `CanvasText` border, so the control survives palette removal.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`.ui-language-switcher`** — the root `<form method="get">`: a wrapping flex row with `--ui-language-switcher-gap`.
- **`.ui-language-switcher-label`** — the optional visible label bound to the select by `for`/`id`; body-sm type in muted foreground.
- **`.ui-language-switcher-control`** — the row holding the select and the submit button together.
- **`.ui-language-switcher-select`** — the native `<select name="lang">`: surface background, small radius, border token, and the shared focus-visible ring.
- **Submit** — a standard [Button](/components/button) (`secondary`) rendered by the `button` template; always visible.

## Variants

- With label — the default pairing of visible label and control.
- Label-free — the control alone, named through `aria-label` on the select.

## Anti-patterns

- Do not POST from the switcher: changing language is navigation, not mutation — a GET keeps it idempotent, cacheable, and bookmarkable.
- Do not hide the submit button or fire the form from `change` events; the visible submit is the no-JS contract.
- Do not hardcode `left`/`right` styling: the flex row flows in document direction, so RTL documents mirror it for free.

## Checklist

- The form is `method="get"` with `name="lang"` on the select.
- The current locale renders `selected`; locales are written in their own language.
- The submit button is visible and reachable by keyboard before the change takes effect.
- The server answers the GET with `303 See Other` to the localized URL.
- A visible label or an `aria-label` names the control.

## Accessibility

The switcher is built entirely from native elements — a labelled `<select>` inside a real form — so assistive technology announces the role, value, and current locale without extra ARIA. Focus follows the natural tab order into the select and then the visible submit, both painted with the shared focus-visible ring, and forced-colors mode keeps label and boundary legible when the palette is removed.
