package state

import (
	"encoding/json"
	"fmt"
	"os"
	"tfviz/internal/model"
)

// ParseFile reads a Terraform state file from disk, parses it, and normalizes it.
func ParseFile(filename string) (*model.State, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rawState RawState
	if err := json.Unmarshal(data, &rawState); err != nil {
		return nil, err
	}

	normalizedState, err := Normalize(&rawState)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize state: %w", err)
	}
	return normalizedState, nil
}
