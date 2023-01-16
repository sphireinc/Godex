package Godex

import (
	"encoding/json"
	mantisDb "github.com/sphireinc/mantis/database"
)

// Codex provides an easy query store
type Codex struct {
	Table   string         `json:"table,omitempty"`
	Queries Queries        `json:"queries,omitempty"`
	DB      mantisDb.MySQL `json:"db,omitempty"`
}

// String returns a string representation of Codex
func (c *Codex) String() string {
	output, err := json.Marshal(&c)
	if err != nil {
		return ""
	}
	return string(output)
}
