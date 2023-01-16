package Godex

import (
	"encoding/json"
	mantisDb "github.com/sphireinc/mantis/database"
)

// Codex provides an easy query store
type Codex struct {
	Table          string            `json:"table,omitempty"`
	DefaultQueries DefaultQueries    `json:"default_queries,omitempty"`
	Queries        map[string]string `json:"queries"`
	DB             mantisDb.MySQL    `json:"db,omitempty"`
	result         any
}

// CxArgs provides an alias for map[string]any for arguments passed to the named queries
type CxArgs map[string]any

// String returns a string representation of Codex
func (c *Codex) String() string {
	output, err := json.Marshal(&c)
	if err != nil {
		return ""
	}
	return string(output)
}
