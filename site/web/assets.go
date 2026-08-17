package web

import "embed"

// Assets contains the documentation source and templates. Keeping the embed
// declaration here allows the public web/ tree to remain the source of truth.
//
//go:embed templates/*.html content/*.md content/templates/*.md static/*
var Assets embed.FS
