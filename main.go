package godex

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

// Codex provides an easy query store
type Codex struct {
	Table          string            `json:"table,omitempty"`
	DefaultQueries DefaultQueries    `json:"default_queries,omitempty"`
	Queries        map[string]string `json:"queries"`
	DB             *sqlx.DB          `json:"-"`
	result         any
}

// Args provides an alias for map[string]any for arguments passed to named queries.
type Args map[string]any

// CxArgs is a deprecated compatibility alias for Args.
//
// Deprecated: use Args instead.
type CxArgs = Args

// New creates a new Codex backed by a sqlx database handle.
func New(db *sqlx.DB, table string, defaults DefaultQueries) *Codex {
	return NewWithQueries(db, table, defaults, nil)
}

// NewWithQueries creates a new Codex and preloads custom named queries.
func NewWithQueries(db *sqlx.DB, table string, defaults DefaultQueries, queries map[string]string) *Codex {
	if queries == nil {
		queries = map[string]string{}
	}
	return &Codex{
		Table:          table,
		DefaultQueries: defaults,
		Queries:        queries,
		DB:             db,
	}
}

// Open opens a sqlx database connection and returns a new Codex.
func Open(driverName, dataSourceName, table string, defaults DefaultQueries, queries map[string]string) (*Codex, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewWithQueries(db, table, defaults, queries), nil
}

// OpenDB opens a standard library database handle and wraps it in sqlx.
func OpenDB(driverName string, db *sql.DB, table string, defaults DefaultQueries, queries map[string]string) *Codex {
	return NewWithQueries(sqlx.NewDb(db, driverName), table, defaults, queries)
}

// String returns a string representation of Codex
func (c *Codex) String() string {
	output, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(output)
}
