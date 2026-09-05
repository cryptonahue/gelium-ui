# Extensibility contract

Gelium UI is extended through server-rendered recipes, registered components, and consumer-owned data adapters. A recipe is the default extension point. Runtime plugins are intentionally out of scope until a real consumer establishes a stable need.

## Contract map

Use the contract that matches the extension before writing implementation code:

- [Application integration and tenant isolation](gelium-ui-application-integration.md)
- [Dashboard metrics](gelium-ui-dashboard-metrics.md)
- [Feedback states](gelium-ui-feedback-recipe.md)
- [Recipe testing](gelium-ui-recipe-testing.md)
- [Export](gelium-ui-export-recipe.md)
- [Import](gelium-ui-import-recipe.md)
- [Relationships and nested resources](gelium-ui-relationships-recipe.md)
- [Screen recipes](gelium-ui-screen-recipes.md)

## Extension order

1. Compose an existing recipe from existing components.
2. Add a consumer-owned data adapter or view model.
3. Add a new component only when composition cannot express the required anatomy or state.
4. Register and document the component after its contract is complete.
5. Consider runtime extension only after repeated consumer evidence and an explicit architecture decision.

## Recipe contract

A recipe should document:

- audience and job to be done;
- route and navigable GET state;
- server-owned view model and data source;
- authorized actions and consumer-owned permission checks;
- semantic HTML composition;
- loading, empty, error, validation, and success states;
- POST + 303 mutations and 422 validation behavior;
- no-JS baseline and optional HTMX enhancement;
- accessibility, responsive, RTL, reduced-motion, and forced-colors behavior;
- recovery paths and non-goals.

Recipes compose components; they do not silently add authentication, authorization, audit persistence, or business rules to Gelium.

## Data source boundary

The consuming application owns queries, joins, sorting semantics, tenant scope, authorization, freshness, and error classification. It supplies a concrete view model with closed values. Gelium renders that model and must not infer sensitive relationships or fetch undeclared fields.

For navigable state, use stable GET parameters. For mutations, use forms with POST and redirect with 303 when navigation follows success. Keep filtering, pagination, and selection in the URL when they affect the rendered result.

## Actions

An action declaration should provide:

```text
label
method
action-url
confirmation-required
availability
success-feedback
failure-feedback
```

The consumer decides `availability` and performs the authorization check at the mutation boundary. Gelium may hide or disable an action for presentation, but that is never the security boundary.

Destructive or consequential actions use a server-rendered confirmation page or dialog, then POST + 303. The action endpoint must re-check authorization and current state.

## State composition

Use the integrated feedback recipe rather than inventing local variants:

- transient result → Toast with a server-rendered fallback;
- persistent success → Banner or inline alert after POST + 303;
- contextual failure → Inline alert;
- multi-field validation → 422 + validation summary;
- unavailable resource → Error state;
- successful empty query → Empty state.

A recipe must not use a Toast as the only channel for validation, authorization denial, critical failure, or a result that must survive navigation.

## Adding a component

A new component requires a complete contract before implementation:

- canonical name and category;
- semantic anatomy and accessible naming;
- variants and relevant states;
- token ownership and theme behavior;
- responsive, RTL, reduced-motion, and forced-colors behavior where applicable;
- server/HTMX contract, if any;
- complete no-JS path;
- alternatives rejected and registry impact.

Follow `lib/skills/14-component-implementation.md` and the component registry checklist. The normal package shape is a template, CSS, source style import, focused contract test, documentation route when needed, and registry entry. Generated assets are rebuilt by the package build and never hand-edited.

## Theme extension

Themes override declared `--ui-*` tokens and preserve component anatomy, state semantics, contrast, focus, reduced-motion, forced-colors, and RTL behavior. A theme must not override scoped component tokens unless the theme contract explicitly permits it.

## Ownership matrix

| Concern | Gelium UI | Consumer application |
|---|---|---|
| HTML/component presentation | Owns | Consumes |
| Recipe composition contract | Documents shared patterns | Chooses and supplies view data |
| Data source and queries | Does not own | Owns |
| Authentication and authorization | Does not implement | Owns and re-checks at mutation |
| Audit persistence | Does not implement | Owns if required |
| Business rules and outcomes | Does not infer | Owns |
| Feedback presentation primitives | Provides | Supplies meaning and content |
| Runtime plugin system | Not provided | Requires explicit future decision |

## Extension checklist

- [ ] Existing recipe or component composition was considered first.
- [ ] The consumer-owned boundary is explicit.
- [ ] Navigable state uses GET and mutations use POST + 303.
- [ ] 422 validation preserves values and field associations.
- [ ] Empty, error, loading, and feedback states are declared.
- [ ] The no-JS path is complete.
- [ ] Accessibility and responsive behavior are documented.
- [ ] New components follow the registry and implementation contract.
- [ ] Generated assets are rebuilt through the existing build.
- [ ] No runtime plugin mechanism was added without a real consumer and architecture decision.

## Related

- [Component registry](gelium-ui-component-registry.md)
- [Pattern registry](gelium-ui-pattern-registry.md)
- [Feedback recipe](gelium-ui-feedback-recipe.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
- [Component implementation skill](../lib/skills/14-component-implementation.md)
