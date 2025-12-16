package group

import (
	"tfviz/internal/model"
)

type AWSALBRule struct {
	Index *Index
}

func (r AWSALBRule) Match(res *model.Resource) bool {
	return isALBRelated(res)
}

func (r AWSALBRule) GroupKey(res *model.Resource) string {
	// If this *is* the LB, it is its own root
	if isALB(res) {
		return res.Address
	}

	// Walk dependencies of all instances
	for _, inst := range res.Instances {
		for _, dep := range inst.DependsOn {
			if parent, ok := r.Index.InstanceToResource[dep]; ok {
				if isALB(parent) {
					return parent.Address
				}
			}
		}
	}

	// Fallback (should be rare)
	return res.Address
}

func (r AWSALBRule) GroupName(res *model.Resource) string {
	return "Application Load Balancer"
}
