package vite

// Manifest type.
type Manifest map[string]ManifestChunk

// ManifestChunk type.
type ManifestChunk struct {
	Src            string   `json:"src,omitempty"`
	File           string   `json:"file"`
	CSS            []string `json:"css,omitempty"`
	Assets         []string `json:"assets,omitempty"`
	IsEntry        bool     `json:"isEntry,omitempty"`
	Name           string   `json:"name,omitempty"`
	Names          []string `json:"names,omitempty"`
	IsDynamicEntry bool     `json:"isDynamicEntry,omitempty"`
	Imports        []string `json:"imports,omitempty"`
	DynamicImports []string `json:"dynamicImports,omitempty"`
}
