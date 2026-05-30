package schema

import "embed"

// Embedded manifests baked into the binary
//
//go:embed schemas/*.json
var EmbeddedFS embed.FS
