package group

import "tfviz/internal/model"

type DefaultRule struct{}

func (r DefaultRule) Match(res *model.Resource) bool {
	return true
}

func (r DefaultRule) GroupKey(res *model.Resource) string {
	return res.Type
}

func (r DefaultRule) GroupName(res *model.Resource) string {
	return "Other Resources"
}
