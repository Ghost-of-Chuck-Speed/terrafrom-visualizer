package state

import (
	"fmt"

	"tfviz/internal/model"
)

func Normalize(raw *RawState) (*model.State, error) {
	out := &model.State{}

	for _, rr := range raw.Resources {
		res := &model.Resource{
			Type:     rr.Type,
			Name:     rr.Name,
			Provider: rr.Provider,
		}

		for _, inst := range rr.Instances {
			index := normalizeIndex(inst.IndexKey)

			addr := rr.Type + "." + rr.Name
			if index != "" {
				addr = fmt.Sprintf("%s[%s]", addr, index)
			}

			i := &model.Instance{
				ID:         addr,
				IndexKey:   index,
				Attributes: inst.Attributes,
				DependsOn:  inst.DependsOn,
			}

			res.Instances = append(res.Instances, i)
		}

		res.Address = rr.Type + "." + rr.Name
		out.Resources = append(out.Resources, res)
	}

	return out, nil
}

func normalizeIndex(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
