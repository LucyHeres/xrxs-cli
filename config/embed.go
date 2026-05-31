package config

import (
	"embed"
	"io/fs"
)

//go:embed schemas/*
var schemaData embed.FS

// SchemaFS returns an fs.FS rooted at the embedded schemas directory.
func SchemaFS() fs.FS {
	sub, _ := fs.Sub(schemaData, "schemas")
	return sub
}
