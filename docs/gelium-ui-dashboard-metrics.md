# Dashboard metrics contract

Gelium UI can render metric cards, but the consuming application owns the meaning, authorization, freshness, and source of every value. A demo metric is not a production data contract.

## Quick path

Before rendering a production metric, define:

1. Its authorized source.
2. The time period and timezone used for the value.
3. The unit and aggregation.
4. Its freshness or `updated-at` value.
5. The meaning and direction of any comparison.
6. Who may view it.
7. Its empty and error behavior.

## Metric definition

Each metric should have an application-owned definition with at least:

```text
key
label
value
unit
period
timezone
updated-at
source
visibility
state
```

Optional comparison data must add:

```text
previous-value
delta
delta-direction
comparison-period
```

`delta-direction` must describe meaning, not only arithmetic sign. For example, a lower error rate may be positive even when the numeric delta is negative.

## Responsibility boundary

| Concern | Owner | Gelium UI role |
|---|---|---|
| Query and source of truth | Consumer application | Render supplied value |
| Tenant/role/metric visibility | Consumer application | Omit unauthorized metrics |
| Period, timezone, unit, and aggregation | Consumer application | Display the declared context |
| Freshness and `updated-at` | Consumer application | Make staleness visible when relevant |
| Comparison and delta semantics | Consumer application | Render the supplied explanation |
| Loading, empty, and error presentation | Shared contract | Provide readable states and recovery |
| Chart data and accessible equivalent | Consumer application | Render approved representation |

## States

### Value

Show the value, label, unit, period, and any freshness context needed to interpret it. Do not present a number without its unit or time basis when those details affect meaning.

### Empty

Use an empty state when the authorized source has no data for the requested period. Explain whether there is no data, no access, or no completed period. Do not render zero unless zero is the factual value.

### Loading

Use a loading state only when a real metric region is waiting for a server or enhanced request. A first server-rendered response is not a fabricated loading phase.

### Error

Use an error state or inline alert when the source fails. Explain the recovery action, such as retrying or choosing another period. Do not replace an unavailable value with zero or stale data without labeling it.

## Filters and periods

If a dashboard supports filters, the state belongs in GET parameters such as `from`, `to`, `status`, `owner`, or `tenant`, subject to the application's authorization rules. Links must be deep-linkable and preserve the selected period and filters.

The application must define:

- inclusive or exclusive date boundaries;
- timezone used for boundaries;
- maximum period length;
- behavior for invalid or reversed ranges;
- whether data is live, delayed, or closed at period end.

## Accessibility and trust

- Pair every value with visible text naming its metric and unit.
- Do not communicate positive or negative meaning through color alone.
- Provide a textual or tabular equivalent for charts.
- Expose freshness in readable text, not only a tooltip.
- Keep unauthorized metrics out of the response; do not render disabled placeholders that reveal their existence.

## Checklist

- [ ] Source and aggregation are authorized and documented.
- [ ] Period, timezone, and unit are visible and unambiguous.
- [ ] Freshness or `updated-at` is defined where relevant.
- [ ] Delta direction has domain meaning, not only a plus/minus sign.
- [ ] Metric visibility follows the consumer's authorization layer.
- [ ] Empty, loading, and error states are distinct and recoverable.
- [ ] GET filters are stable and deep-linkable.
- [ ] Any chart has a readable equivalent.

## Related

- [Application integration](gelium-ui-application-integration.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
- [Server contracts](../lib/skills/05-server-contracts.md)
