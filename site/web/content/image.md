# Image

Image is the semantic figure pattern for content images: a real `<img>` with descriptive alt text, intrinsic dimensions, and responsive sources, wrapped in a `<figure>` with an optional caption. Use it when a page embeds a content image that must stay accessible, sized, and responsive without JavaScript. The markup styles through `media.css` — the image templates share the media system and have no dedicated styles file.

## Examples

A single content image with intrinsic dimensions and lazy loading.

<div class="specimen-block">
<figure class="ui-media ui-media-image">
  <img src="https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=60" alt="A mountain lake at dusk" width="1200" height="800" loading="lazy">
  <figcaption>Alpine lake at dusk — the day the trail opened.</figcaption>
</figure>
</div>

An image with responsive sources and a fixed aspect fallback.

<div class="specimen-block">
<figure class="ui-media ui-media-image ui-media--aspect" style="--ui-media-aspect:16 / 9">
  <img src="https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=60" alt="A mountain lake at dusk" width="1200" height="800" loading="lazy" srcset="https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=800&q=60 800w, https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=60 1200w" sizes="(min-width: 60rem) 50vw, 100vw">
  <figcaption>Responsive sources let the browser pick the right file for the viewport.</figcaption>
</figure>
</div>

A `<picture>` element serves modern formats with a fallback.

<div class="specimen-block">
<figure class="ui-media ui-media-picture">
  <picture>
    <source srcset="https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=60&fm=webp" type="image/webp">
    <img src="https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=60" alt="A mountain lake at dusk" width="1200" height="800" loading="lazy">
  </picture>
  <figcaption>Picture sources use the image type marker so browsers skip what they cannot decode.</figcaption>
</figure>
</div>

The specimens above follow the contracts the templates `image.html` defines (`image` and `picture`).

## Guidance

### When to use

Use the image pattern for content images the page actually communicates with — a photograph, a diagram, an illustration — whenever the image needs alt text and a size contract. It earns its place when the image is content, not decoration.

### When not to use

Do not use it for purely decorative images: give those an empty `alt` or move them out of the content flow entirely. Do not use JavaScript-based lazy loaders — the native `loading="lazy"` attribute covers it. If the image is a background or a UI surface, a plain `img` or a CSS background is the right tool, not the figure pattern.

### Usability

- Always set intrinsic `width` and `height` so the layout reserves space and avoids layout shift while the image loads.
- Default to `loading="lazy"` for below-the-fold images; keep eager loading only for the first content image above the fold.
- Use `srcset` with `sizes` when the same image has to resize across viewports; use `<picture>` when you need format variants (for example WebP) with a universal fallback.

### Accessibility

- The `alt` text says what the image means in context; a mountain lake photograph gets a descriptive caption-like alt, a chart gets its data story, and decoration gets `alt=""`.
- The caption is a `<figcaption>` inside the figure, so the image and its caption stay one unit for assistive technology.
- Images and pictures render fluid and `display: block` with `max-width: 100%`, so they never overflow their container.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

The image pattern relies on the media system (`media.css`), not a dedicated image stylesheet — `image.css` does not exist by design.

- **`ui-media`** — the base figure: block flow with `--ui-space-4` vertical margin.
- **`ui-media-image` / `ui-media-picture`** — the figure kinds the templates emit.
- **`ui-media--aspect`** — the aspect fallback: the figure reserves a ratio before paint (`--ui-media-aspect`, default `16 / 9`) with the surface-container background behind it.
- **`figcaption`** — the caption, capped at 65ch in the muted variant color.

## Template contract

The `image` template renders a `<figure>` with an `<img>` carrying `src`, `alt`, `width`, `height`, `loading` (default `lazy`), and optional `srcset` and `sizes`. The `picture` template renders the same figure shell with a `<picture>` whose `<source>` elements carry `srcset`, `sizes`, and `type`. The `alt`, `width`, and `height` attributes are always emitted — they are not optional in the contract.

## Anti-patterns

- Do not omit `width` and `height`: the layout shifts and the reserve-space contract breaks.
- Do not write generic alt text ("image", "photo") — the alt must carry the meaning the image has in context.
- Do not invent image-specific class names: the figure uses the media system classes above.

## Checklist

- The figure has a real `<img>` with descriptive alt text.
- `width` and `height` are set from the intrinsic dimensions.
- `loading="lazy"` is set unless the image is above the fold.
- Responsive images use `srcset`/`sizes` or `<picture>` with a fallback.
- The caption, when present, is a `<figcaption>` inside the figure.

## Accessibility

The figure pattern keeps images and captions in one unit, reserves layout space through intrinsic dimensions, and relies on the platform's native lazy loading. The aspect fallback paints the surface-container color behind the image until it loads, so the shape of the content is visible while the file arrives. In the media system, images render block and fluid so they never overflow narrow containers.
