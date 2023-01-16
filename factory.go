package Godex

import (
	mantisDb "github.com/sphireinc/mantis/database"
)

// CreateGodex creates a new instance of Codex
func CreateGodex(
	db mantisDb.MySQL,
	table string,
	SelectById string,
	SelectOne string,
	Select string,
	Insert string,
	Update string,
	Delete string,
	SoftDelete string,
) Codex {
	return Codex{
		Table: table,
		DefaultQueries: DefaultQueries{
			SelectById: SelectById,
			SelectOne:  SelectOne,
			Select:     Select,
			Insert:     Insert,
			Update:     Update,
			Delete:     Delete,
			SoftDelete: SoftDelete,
		},
		Queries: map[string]string{},
		DB:      db,
	}
}
