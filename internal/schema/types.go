package schema

// Manifest is the top-level schema file describing all products and their tools.
type Manifest struct {
	Version  string    `json:"version"`
	Products []Product `json:"products"`
}

// Product represents a business module (approval, staff, attendance, etc).
type Product struct {
	Name        string `json:"name"`        // CLI command name: "approval"
	Description string `json:"description"` // Shown in --help
	Tools       []Tool `json:"tools"`       // Top-level subcommand groups
}

// Tool is a command group (e.g. "list", "detail").
// If Subtools is non-empty, this is a parent command.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subtools    []Tool `json:"subtools,omitempty"` // Leaf commands or nested groups

	// Leaf command fields (only populated for leaves — no Subtools):
	Path     string      `json:"path,omitempty"`     // API endpoint path
	Method   string      `json:"method,omitempty"`   // HTTP method
	Encoding string      `json:"encoding,omitempty"` // "form", "json", "form-nested"
	Output   OutputSpec  `json:"output,omitempty"`
	Params   []Param     `json:"params,omitempty"`
	Body     interface{} `json:"body,omitempty"` // Static body template (map[string]any or nil)
}

// OutputSpec describes how to process the API response.
type OutputSpec struct {
	Unwrap string            `json:"unwrap,omitempty"` // "data.result.data" or "" for direct output
	Labels map[string]string `json:"labels,omitempty"` // Field name → display name (e.g. "employeeName": "发起人")
}

// Param describes a CLI flag that maps to an API parameter.
// Name is the API parameter name in camelCase (e.g. "flowStepId").
// The CLI flag name is auto-derived from Name via toKebabCase conversion in the builder.
type Param struct {
	Name        string      `json:"name"`                  // API parameter name (camelCase)
	Short       string      `json:"short,omitempty"`       // Short flag: "k"
	Type        string      `json:"type"`                  // "string", "int", "bool"
	Default     interface{} `json:"default,omitempty"`     // Default value
	Description string      `json:"description,omitempty"` // Help text
	Required    bool        `json:"required,omitempty"`
	InBody      bool        `json:"inBody,omitempty"` // If true, goes into request body (POSTFormJSON)
	QueryName   string      `json:"queryName,omitempty"` // API param name override, only needed when semantically different from Name
}

// APIName returns the parameter name to use in API requests.
// If QueryName is set, use it directly. Otherwise, Name is already the correct camelCase API name.
func (p Param) APIName() string {
	if p.QueryName != "" {
		return p.QueryName
	}
	return p.Name
}

// IsLeaf returns true if this tool is a leaf command (executes an API call).
func (t Tool) IsLeaf() bool {
	return len(t.Subtools) == 0 && t.Path != ""
}
