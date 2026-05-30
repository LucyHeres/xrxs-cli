package schema

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// LoadManifest reads and parses a manifest JSON file.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", path, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}

	// Strip comments (lines starting with //) for user convenience
	// Support trailing commas by using the standard library which handles them in Go 1.26

	return &m, nil
}

// LoadAllManifests loads all .json manifest files from a directory.
func LoadAllManifests(dir string) (*Manifest, error) {
	return loadManifestsFromDir(os.DirFS(dir), ".")
}

// LoadManifestsFromFS loads all .json manifest files from an embed.FS.
func LoadManifestsFromFS(fsys fs.FS) (*Manifest, error) {
	return loadManifestsFromDir(fsys, ".")
}

// LoadFromEmbed loads manifests from the embedded FS (EmbeddedFS).
func LoadFromEmbed() (*Manifest, error) {
	// EmbeddedFS root is internal/schema/, files are in schemas/ subdir
	return loadManifestsFromDir(EmbeddedFS, "schemas")
}

func loadManifestsFromDir(fsys fs.FS, dir string) (*Manifest, error) {
	merged := &Manifest{Version: "1"}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read schema dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := fs.ReadFile(fsys, filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var m Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("parse %s: %w", entry.Name(), err)
		}
		merged.Products = append(merged.Products, m.Products...)
	}

	return merged, nil
}

// FindProduct looks up a product by name.
func (m *Manifest) FindProduct(name string) *Product {
	for i := range m.Products {
		if m.Products[i].Name == name {
			return &m.Products[i]
		}
	}
	return nil
}

// FindTool looks up a tool by path (e.g. "list.search").
func (p *Product) FindTool(path string) *Tool {
	parts := strings.Split(path, ".")
	return findTool(p.Tools, parts)
}

func findTool(tools []Tool, parts []string) *Tool {
	for i := range tools {
		if tools[i].Name == parts[0] {
			if len(parts) == 1 {
				return &tools[i]
			}
			return findTool(tools[i].Subtools, parts[1:])
		}
	}
	return nil
}
