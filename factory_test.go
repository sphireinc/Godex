package Godex

import (
	mantisDb "github.com/sphireinc/mantis/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGodex(t *testing.T) {
	db := mantisDb.MySQL{}
	q := CreateGodex(
		db,
		"someTable",
		"SelectById",
		"SelectOne",
		"Select",
		"Insert",
		"Update",
		"Delete",
		"SoftDelete",
	)
	assert.Equal(t, "someTable", q.Table)
	assert.Equal(t, "Insert", q.DefaultQueries.Insert)
}
