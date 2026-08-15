package app

import "strings"

// docDate is the ISO-8601 (2006-01-02) publication and modification dates for
// one component docs page, sourced from the git history of its content file
// (first commit = published, last commit = modified; equal when a file was
// never touched after creation). The values feed both the visible provenance
// line and the TechArticle JSON-LD so the article and the structured data
// always agree (GEO §7).
type docDate struct {
	Published string
	Modified  string
}

// componentDates maps every registered component slug to its published and
// modified dates. index.md is not a component (it renders /) and is excluded.
// The table is the single source of dates: no goldmark/frontmatter change.
var componentDates = map[string]docDate{
	"badge":             {Published: "2026-08-09", Modified: "2026-08-09"},
	"button":            {Published: "2026-08-09", Modified: "2026-08-09"},
	"card":              {Published: "2026-08-09", Modified: "2026-08-09"},
	"checkbox":          {Published: "2026-08-14", Modified: "2026-08-14"},
	"chips":             {Published: "2026-08-14", Modified: "2026-08-14"},
	"data-table":        {Published: "2026-08-10", Modified: "2026-08-14"},
	"dialog":            {Published: "2026-08-09", Modified: "2026-08-11"},
	"divider":           {Published: "2026-08-09", Modified: "2026-08-09"},
	"elevation":         {Published: "2026-08-14", Modified: "2026-08-14"},
	"fab":               {Published: "2026-08-09", Modified: "2026-08-09"},
	"focus-ring":        {Published: "2026-08-14", Modified: "2026-08-14"},
	"icon":              {Published: "2026-08-14", Modified: "2026-08-14"},
	"icon-button":       {Published: "2026-08-09", Modified: "2026-08-09"},
	"list":              {Published: "2026-08-14", Modified: "2026-08-14"},
	"menu":              {Published: "2026-08-14", Modified: "2026-08-14"},
	"navigation-bar":    {Published: "2026-08-14", Modified: "2026-08-14"},
	"navigation-drawer": {Published: "2026-08-14", Modified: "2026-08-14"},
	"navigation-tab":    {Published: "2026-08-14", Modified: "2026-08-14"},
	"progress":          {Published: "2026-08-14", Modified: "2026-08-14"},
	"radio":             {Published: "2026-08-14", Modified: "2026-08-14"},
	"segmented-button":  {Published: "2026-08-14", Modified: "2026-08-14"},
	"select":            {Published: "2026-08-14", Modified: "2026-08-14"},
	"slider":            {Published: "2026-08-14", Modified: "2026-08-14"},
	"switch":            {Published: "2026-08-14", Modified: "2026-08-14"},
	"tabs":              {Published: "2026-08-14", Modified: "2026-08-14"},
	"text-field":        {Published: "2026-08-09", Modified: "2026-08-09"},
	"toast":             {Published: "2026-08-14", Modified: "2026-08-14"},
	"tooltip":           {Published: "2026-08-10", Modified: "2026-08-14"},
}

// docDatesFor returns the published/modified dates for a content slug and
// whether the slug has a table entry. Slugs absent from the table (including
// unregistered /components/* paths) report not-found so callers emit no dates.
func docDatesFor(slug string) (docDate, bool) {
	d, ok := componentDates[strings.TrimPrefix(slug, "/components/")]
	return d, ok
}
