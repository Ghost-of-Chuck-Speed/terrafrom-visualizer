package group

import "tfviz/internal/model"

type Index struct {
	InstanceToResource map[string]*model.Resource
}

func BuildIndex(state *model.State) *Index {
	idx := &Index{
		InstanceToResource: make(map[string]*model.Resource),
	}

	for _, r := range state.Resources {
		for _, i := range r.Instances {
			idx.InstanceToResource[i.ID] = r
		}
	}

	return idx
}
