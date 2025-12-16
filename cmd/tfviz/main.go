package main

import (
	"fmt"
	"os"
	"tfviz/internal/group"
	"tfviz/internal/state"
	"tfviz/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tfviz <tfstate-file>")
		return
	}

	filename := os.Args[1]
	st, err := state.ParseFile(filename)
	if err != nil {
		fmt.Println("Error parsing state:", err)
		return
	}

	// Apply default grouping rules
	idx := group.BuildIndex(st)
	rules := []group.Rule{
		group.AWSALBRule{Index: idx},
		group.AWSEKSRule{},
		group.AWSRDSRule{},
		group.AWSS3Rule{},
		group.AWSVPCRule{},
		group.DefaultRule{},
	}
	groups := group.Apply(rules, st)

	// Initialize TUI model
	model := ui.NewModel(groups)
	if err := ui.Run(model); err != nil {
		fmt.Println("Error running UI:", err)
	}
}
