package ui

import "tfviz/internal/model"

func ResourceLabel(r *model.Resource) string {
	// Try tags.Name first
	for _, inst := range r.Instances {
		if tags, ok := inst.Attributes["tags"].(map[string]any); ok {
			if name, ok := tags["Name"].(string); ok {
				return r.Type + " : " + name
			}
		}
		// fallback to name attribute
		if name, ok := inst.Attributes["name"].(string); ok {
			return r.Type + " : " + name
		}
		// fallback to ARN
		if arn, ok := inst.Attributes["arn"].(string); ok {
			return r.Type + " : " + arn
		}
	}
	// fallback to Terraform address
	return r.Address
}

func InstanceLabel(i *model.Instance) string {
	if arn, ok := i.Attributes["arn"].(string); ok {
		return arn
	}
	if name, ok := i.Attributes["name"].(string); ok {
		return name
	}
	return i.ID
}
