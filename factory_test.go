package godex

import (
	"testing"
)

func TestCreateGodex(t *testing.T) {
	q := CreateGodex(
		nil,
		"someTable",
		"SelectById",
		"SelectOne",
		"Select",
		"Insert",
		"Update",
		"Delete",
		"SoftDelete",
	)
	if q.Table != "someTable" {
		t.Fatalf("expected table someTable, got %s", q.Table)
	}
	if q.DefaultQueries.Insert != "Insert" {
		t.Fatalf("expected insert query Insert, got %s", q.DefaultQueries.Insert)
	}
}
