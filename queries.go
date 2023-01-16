package Godex

// Queries stores specific queries to be used by Codex
type Queries struct {
	SelectById string `json:"selectById,omitempty"`
	SelectOne  string `json:"selectOne,omitempty"`
	Select     string `json:"select,omitempty"`
	Insert     string `json:"insert,omitempty"`
	Update     string `json:"update,omitempty"`
	Delete     string `json:"delete,omitempty"`
	SoftDelete string `json:"softDelete,omitempty"`
}

// RawQuery performs a query with given args
func (c *Codex) RawQuery(query string, args ...any) (any, error) {
	return c.DB.RawQuery(query, args)
}

// SelectById one item
func (c *Codex) SelectById(args ...any) (any, error) {
	return c.DB.SelectOne(c, c.Queries.SelectById, args...)
}

// SelectOne item
func (c *Codex) SelectOne(args ...any) (any, error) {
	return c.DB.SelectOne(c, c.Queries.SelectOne, args...)
}

// Select one
func (c *Codex) Select(into []any, args ...any) (any, []any, error) {
	record, err := c.DB.Select(&into, c.Queries.Select, args...)
	return record, into, err
}

// Insert a new record
func (c *Codex) Insert() (int64, error) {
	return c.DB.InsertOne(c.Queries.Insert, c)
}

// Update a specific record
func (c *Codex) Update() (int64, error) {
	return c.DB.UpdateOne(c.Queries.Update, c)
}

// Delete a specific record
func (c *Codex) Delete() error {
	return c.DB.DeleteOne(c.Queries.Delete, c)
}

// SoftDelete a specific record
func (c *Codex) SoftDelete() (int64, error) {
	return c.DB.UpdateOne(c.Queries.SoftDelete, c)
}
