package ui

import "fmt"

func RenderTree(root *TreeNode, selected *TreeNode, width int) []string {
	var lines []string
	renderNode(root, selected, 0, &lines, width)
	return lines
}

func renderNode(n *TreeNode, selected *TreeNode, indent int, lines *[]string, width int) {
	if n.Parent != nil { // skip root label
		prefix := ""
		for i := 0; i < indent-1; i++ {
			prefix += "  "
		}

		cursor := "  "
		if n == selected {
			cursor = "▶ "
		}

		line := fmt.Sprintf("%s%s%s", prefix, cursor, n.Label)
		if len(line) > width {
			line = line[:width-1]
		}

		*lines = append(*lines, line)
	}

	if !n.Expanded {
		return
	}

	for _, c := range n.Children {
		renderNode(c, selected, indent+1, lines, width)
	}
}
