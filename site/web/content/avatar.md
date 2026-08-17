# Avatar

Avatar is a circular surface that shows initials or an image at one of three sizes. Use it to mark a person or entity next to a name in a feed, list, or comment thread. As a recipe primitive it is decorative by default: pair it with visible text so the identity never depends on the circle alone.

## Specimen

This page renders the real `avatar` template markup live. The decorative initials avatar is `aria-hidden` and must sit next to the visible name that supplies the meaning:

<div class="specimen-block">
<p>Signed in as <span class="ui-avatar ui-avatar--md" aria-hidden="true"><span class="ui-avatar-initials">AR</span></span> Alicia R.</p>
</div>

The three sizes use the closed `ui-avatar--{sm,md,lg}` modifier set:

<div class="specimen-block">
<p><span class="ui-avatar ui-avatar--sm" aria-hidden="true"><span class="ui-avatar-initials">AR</span></span> <span class="ui-avatar ui-avatar--md" aria-hidden="true"><span class="ui-avatar-initials">AR</span></span> <span class="ui-avatar ui-avatar--lg" aria-hidden="true"><span class="ui-avatar-initials">AR</span></span></p>
</div>

An image avatar that carries meaning on its own keeps its `alt` text and is not hidden:

<div class="specimen-block">
<p><span class="ui-avatar ui-avatar--lg"><img class="ui-avatar-image" src="/static/rich-article-image.svg" alt="Alicia R."></span> Alicia R.</p>
</div>

## Guidance

### When to use

Use an avatar to mark identity next to a name — in a feed row, a list item, a comment header, or a member list. Use initials when there is no image; use the image variant when a portrait exists and is current.

### When not to use

Never use an avatar as the only identity signal: a decorative avatar is `aria-hidden`, so the identity must always survive in adjacent visible text. Do not make an avatar a control on its own — wrap it in a real link or button with its own accessible name when it is clickable. Do not invent sizes outside the closed `sm`/`md`/`lg` set, and do not stretch the circle with non-square source images.

### Usability

- The image variant wins over initials whenever `ImageSrc` is set; initials are the fallback when it is empty.
- Sizes map onto the core scale: `sm` = `--ui-size-control`, `md` = `--ui-size-item`, `lg` = `--ui-size-item-lg` (`--ui-avatar-size-*`), with matching font steps for the initials.
- Keep source images square: the circle crops with `object-fit: cover`, and a portrait crop is the expected look.

### Accessibility

- A decorative avatar paired with a visible name renders `aria-hidden="true"`; the surrounding text supplies the name.
- A meaningful image avatar keeps `aria-hidden` off and sets `alt` to the person or entity name.
- A decorative image avatar renders an empty `alt=""` — never a missing alt or the filename.
- In forced-colors mode the circle keeps a `CanvasText` boundary so the identity marker remains visible when color is removed.

## Anatomy

- **`ui-avatar`** — the circular surface: `--ui-avatar-container` (`--ui-color-surface-container`) with `--ui-avatar-fg` ink, a full radius (`--ui-radius-full`), and `flex: none` so it never shrinks inside a row.
- **Size modifiers** — `ui-avatar--sm`, `ui-avatar--md` (default when empty), `ui-avatar--lg`; each sets the `--ui-avatar-size-*` geometry and a font step for the initials.
- **`ui-avatar-initials`** — the text fallback, shown when no image is set.
- **`ui-avatar-image`** — the image layer, `100%` of the circle with `object-fit: cover`.

All paint comes from the scoped `--ui-avatar-*` tokens so the primitive works standalone and a theme may override it globally.

## Variants

- **Sizes** — the closed `sm`/`md`/`lg` modifier set; an unknown size sanitizes to `md`.
- **Content** — initials (decorative by nature, paired with a name) or image (meaningful with `alt`, or decorative with `alt=""`); the image wins over initials whenever it is set.

## Sources

- Registry: `docs/gelium-ui-component-registry.md` (Avatar row, recipe primitive RP) — `.ui-avatar`, `--ui-avatar-*`, `--ui-color-surface-container`, `--ui-radius-full`; no states (decorative, `aria-hidden`); variants `ui-avatar--{sm,md,lg}` and image/initials.
- Vocabulary: `docs/gelium-ui-vocabulary.md` §4 — Feed composes List three-line or Card + avatar + Badge; the avatar is the identity marker.
- Implementation: `lib/templates/avatar.html`, `lib/styles/avatar.css`; view model `avatarView` in `internal/app/avatar.go`.

See also: [List](/components/list) for the two-line identity rows avatars sit in, [Screens](/docs/screens) for identity patterns across screen types.
