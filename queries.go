package Godex

// DefaultQueries stores specific queries to be used by Codex
type DefaultQueries struct {
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
	return c.DB.SelectOne(c.result, c.DefaultQueries.SelectById, args...)
}

// SelectOne item
func (c *Codex) SelectOne(args ...any) (any, error) {
	return c.DB.SelectOne(c, c.DefaultQueries.SelectOne, args...)
}

// Select one
func (c *Codex) Select(args ...any) ([]any, error) {
	var res []any
	_, err := c.DB.Select(res, c.DefaultQueries.Select, args...)
	return res, err
}

// Insert a new record
func (c *Codex) Insert() (int64, error) {
	return c.DB.InsertOne(c.DefaultQueries.Insert, c)
}

// Update a specific record
func (c *Codex) Update() (int64, error) {
	return c.DB.UpdateOne(c.DefaultQueries.Update, c)
}

// Delete a specific record
func (c *Codex) Delete() error {
	return c.DB.DeleteOne(c.DefaultQueries.Delete, c)
}

// SoftDelete a specific record
func (c *Codex) SoftDelete() (int64, error) {
	return c.DB.UpdateOne(c.DefaultQueries.SoftDelete, c)
}
