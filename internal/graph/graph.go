package graph

import "tfviz/internal/model"

type Node struct {
	ID          string
	Instance    *model.Instance
	DependsOn   []*Node
	ReverseDeps []*Node
}

type Graph struct {
	Nodes map[string]*Node
}

func Build(state *model.State) *Graph {
	g := &Graph{
		Nodes: make(map[string]*Node),
	}

	for _, r := range state.Resources {
		for _, i := range r.Instances {
			g.Nodes[i.ID] = &Node{
				ID:       i.ID,
				Instance: i,
			}
		}
	}

	for _, n := range g.Nodes {
		for _, dep := range n.Instance.DependsOn {
			if target, ok := g.Nodes[dep]; ok {
				n.DependsOn = append(n.DependsOn, target)
				target.ReverseDeps = append(target.ReverseDeps, n)
			}
		}
	}

	return g
}
