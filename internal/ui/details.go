package ui

import (
	"fmt"
)

// RenderDetails prints details for the selected TreeNode
func RenderDetails(n *TreeNode) string {
	if n == nil {
		return ""
	}

	switch n.Type {
	case NodeTypeGroup:
		return fmt.Sprintf("Group: %s\nResources: %d", n.Label, len(n.Children))
	case NodeTypeResource:
		if n.Resource == nil {
			return fmt.Sprintf("Resource: %s (no details)", n.Label)
		}

		// Extract Name tag if exists
		nameTag := ""
		if n.Resource.Attributes != nil {
			if tags, ok := n.Resource.Attributes["tags"].(map[string]any); ok {
				if v, ok := tags["Name"]; ok {
					nameTag = fmt.Sprintf("Name Tag: %v\n", v)
				}
			}
		}

		// Extract ARN if exists
		arn := ""
		if n.Resource.Attributes != nil {
			if a, ok := n.Resource.Attributes["arn"].(string); ok {
				arn = fmt.Sprintf("ARN: %s\n", a)
			}
		}

		return fmt.Sprintf(
			"Resource: %s\nType: %s\n%s%sInstances: %d",
			n.Label,
			n.Resource.Type,
			arn,
			nameTag,
			len(n.Resource.Instances),
		)

	case NodeTypeInstance:
		if n.Instance == nil {
			return fmt.Sprintf("Instance: %s (no details)", n.Label)
		}
		id := n.Instance.ID
		index := n.Instance.IndexKey
		return fmt.Sprintf("Instance: %s\nIndex: %s", id, index)

	default:
		return n.Label
	}
}
