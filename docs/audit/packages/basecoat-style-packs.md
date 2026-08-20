# Basecoat style packs in Gelium

Basecoat separates:

```text
base/          shared tokens
components/    structure
styles/*.css   visual packs: vega nova maia lyra mira luma sera rhea
```

Gelium maps that family as product skins:

| Gelium skin        | Basecoat pack | Notes                          |
|--------------------|---------------|--------------------------------|
| `basecoat`         | Vega          | default / docs default         |
| `basecoat-nova`    | Nova          | denser, rounded-lg             |
| `basecoat-maia`    | Maia          | pill / 4xl                     |
| `basecoat-lyra`    | Lyra          | square, smaller type           |
| `basecoat-mira`    | Mira          | compact                        |
| `basecoat-luma`    | Luma          | soft pill                      |
| `basecoat-sera`    | Sera          | uppercase tracking, square     |
| `basecoat-rhea`    | Rhea          | soft 2xl                       |

Color foundations still use `.theme-basecoat`. Pack differences are anatomy
tokens under `html[data-gelium-skin="basecoat-…"]`.

Evidence: `docs/audit/packages/basecoat-css-1.0.2.tgz`.
