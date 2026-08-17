// Package lib contains the reusable Gelium UI component library: the
// component templates and their styles, embedded as a single asset bundle
// for consumers.
package lib

import "embed"

// LibAssets is the embed file system carrying the component library sources.
// The templates and styles live next to this declaration so //go:embed works
// for both in-tree and external consumers.
//
//go:embed templates/*.html styles/*.css
var LibAssets embed.FS
