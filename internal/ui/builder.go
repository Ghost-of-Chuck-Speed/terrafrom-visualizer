package ui

import "tfviz/internal/group"

func BuildTree(groups map[string]*group.Group) *TreeNode {
	root := &TreeNode{
		Label:    "Terraform Resources",
		Children: []*TreeNode{},
		Expanded: true,
		Type:     NodeTypeGroup,
	}

	for _, g := range groups {
		groupNode := &TreeNode{
			Label:    g.Name,
			Children: []*TreeNode{},
			Parent:   root,
			Expanded: true,
			Type:     NodeTypeGroup,
		}

		for _, res := range g.Resources {
			label := res.Name
			if label == "" {
				label = res.Address
			}

			resNode := &TreeNode{
				Label:    label,
				Parent:   groupNode,
				Expanded: false,
				Type:     NodeTypeResource,
				Resource: res,
			}

			groupNode.Children = append(groupNode.Children, resNode)
		}

		root.Children = append(root.Children, groupNode)
	}

	return root
}
