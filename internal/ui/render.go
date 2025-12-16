package ui

func RenderNodeDetails(node *TreeNode) string {
	if node == nil {
		return ""
	}
	return RenderDetails(node)
}
