package gates

// RoutePlanResult is the machine-readable handoff for a task-start route.
// It describes the next workflow step; it never authorizes delivery.
type RoutePlanResult struct {
	Route              Route  `json:"route"`
	Status             string `json:"status"`
	Next               string `json:"next"`
	RequiresDesignGate bool   `json:"requires_design_gate"`
}

// PlanRoute maps a validated route to its smallest next action.
func PlanRoute(route Route) (RoutePlanResult, bool) {
	plans := map[Route]RoutePlanResult{
		RouteDirectExempt: {
			Route: route, Status: "pass", Next: "implement-bounded-action",
		},
		RouteDelegatedDirect: {
			Route: route, Status: "pass", Next: "delegate-bounded-action",
		},
		RouteDesignGated: {
			Route: route, Status: "pass", Next: "prepare-visible-architecture-packet", RequiresDesignGate: true,
		},
		RouteEscalate: {
			Route: route, Status: "needs-decision", Next: "request-concrete-decision",
		},
		RouteFullSDD: {
			Route: route, Status: "pass", Next: "start-openspec",
		},
	}
	plan, ok := plans[route]
	return plan, ok
}
