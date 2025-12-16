package ui

import "tfviz/internal/model"

type NodeType int

const (
	NodeTypeGroup NodeType = iota
	NodeTypeResource
	NodeTypeInstance
)

type TreeNode struct {
	Label    string
	Children []*TreeNode
	Parent   *TreeNode
	Expanded bool
	Type     NodeType
	Resource *model.Resource
	Instance *model.Instance
}
