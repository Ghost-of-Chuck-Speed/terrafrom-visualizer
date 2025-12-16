package state

type RawState struct {
	Version          int           `json:"version"`
	TerraformVersion string        `json:"terraform_version"`
	Serial           int           `json:"serial"`
	Lineage          string        `json:"lineage"`
	Resources        []RawResource `json:"resources"`
}

type RawResource struct {
	Mode      string        `json:"mode"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	Instances []RawInstance `json:"instances"`
}

type RawInstance struct {
	SchemaVersion int            `json:"schema_version"`
	Attributes    map[string]any `json:"attributes"`
	DependsOn     []string       `json:"depends_on"`
	IndexKey      any            `json:"index_key,omitempty"`
}
