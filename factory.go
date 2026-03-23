package godex

import (
	"github.com/jmoiron/sqlx"
)

// CreateGodex creates a new instance of Codex.
//
// Deprecated: use New or NewWithQueries instead.
func CreateGodex(
	db *sqlx.DB,
	table string,
	SelectById string,
	SelectOne string,
	Select string,
	Insert string,
	Update string,
	Delete string,
	SoftDelete string,
) Codex {
	return *New(db, table, DefaultQueries{
		SelectById: SelectById,
		SelectOne:  SelectOne,
		Select:     Select,
		Insert:     Insert,
		Update:     Update,
		Delete:     Delete,
		SoftDelete: SoftDelete,
	})
}
