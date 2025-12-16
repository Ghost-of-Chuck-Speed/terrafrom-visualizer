package group

import "tfviz/internal/model"

type Group struct {
	Key       string
	Name      string
	Resources []*model.Resource
}
