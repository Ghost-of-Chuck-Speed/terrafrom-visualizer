package model

type State struct {
	Resources []*Resource
}

type Resource struct {
	Address    string
	Type       string
	Name       string
	Provider   string
	Instances  []*Instance
	Attributes map[string]any
}

type Instance struct {
	ID         string
	IndexKey   string
	Attributes map[string]any
	DependsOn  []string
}
