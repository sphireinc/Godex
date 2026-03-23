package godex

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSoftDeleteWhereClause(t *testing.T) {
	if got := softDeleteWhereClause(""); got != " deleted_at IS NULL" {
		t.Fatalf("unexpected clause: %q", got)
	}
	if got := softDeleteWhereClause("id = :id"); got != " deleted_at IS NULL AND id = :id" {
		t.Fatalf("unexpected clause with where: %q", got)
	}
}

func TestPretty(t *testing.T) {
	if got := Pretty(`{"id":1}`); !strings.Contains(got, "\n") {
		t.Fatalf("expected pretty JSON, got %q", got)
	}
	if got := Pretty("not-json"); got != "" {
		t.Fatalf("expected empty string for invalid json, got %q", got)
	}
}

func TestString(t *testing.T) {
	codex := New(nil, "posts", DefaultQueries{SelectOne: "SELECT 1"})
	value := codex.String()
	if !strings.Contains(value, `"table":"posts"`) {
		t.Fatalf("expected table in json string, got %s", value)
	}
	if strings.Contains(value, `"DB"`) {
		t.Fatalf("expected DB field to be omitted, got %s", value)
	}
}

func TestOpenError(t *testing.T) {
	_, err := Open("definitely-invalid-driver", "dsn", "posts", DefaultQueries{}, nil)
	if err == nil {
		t.Fatal("expected Open to fail for invalid driver")
	}
}

func TestOpenSuccess(t *testing.T) {
	db, _, err := sqlmock.NewWithDSN("godex-open")
	if err != nil {
		t.Fatalf("sqlmock.NewWithDSN returned error: %v", err)
	}
	defer db.Close()

	codex, err := Open("sqlmock", "godex-open", "posts", DefaultQueries{SelectOne: "SELECT 1"}, nil)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	if codex.DB == nil {
		t.Fatal("expected DB to be configured")
	}
}

func TestOpenDB(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New returned error: %v", err)
	}
	defer db.Close()

	codex := OpenDB("sqlmock", db, "posts", DefaultQueries{SelectOne: "SELECT 1"}, nil)
	if codex.DB == nil {
		t.Fatal("expected DB to be configured")
	}
	if codex.Table != "posts" {
		t.Fatalf("expected table posts, got %s", codex.Table)
	}
}

func TestRegisterAndLookupQuery(t *testing.T) {
	codex := New(nil, "posts", DefaultQueries{})
	codex.RegisterQuery("by_title", "SELECT * FROM posts WHERE title = :title")

	query, err := codex.LookupQuery("by_title")
	if err != nil {
		t.Fatalf("LookupQuery returned error: %v", err)
	}
	if query != "SELECT * FROM posts WHERE title = :title" {
		t.Fatalf("unexpected query: %s", query)
	}
}

func TestRegisterQueryInitializesMap(t *testing.T) {
	codex := &Codex{}
	codex.RegisterQuery("by_id", "SELECT * FROM posts WHERE id = :id")
	if codex.Queries["by_id"] == "" {
		t.Fatal("expected query map to be initialized")
	}
}

func TestLookupQueryNotFound(t *testing.T) {
	codex := New(nil, "posts", DefaultQueries{})
	_, err := codex.LookupQuery("missing")
	if !errors.Is(err, ErrQueryNotFound) {
		t.Fatalf("expected ErrQueryNotFound, got %v", err)
	}
}

func TestEnsureDBAndValidationHelpers(t *testing.T) {
	var nilCodex *Codex
	if err := nilCodex.ensureDB(); !errors.Is(err, ErrDBNotConfigured) {
		t.Fatalf("expected ErrDBNotConfigured for nil codex, got %v", err)
	}

	codex := New(nil, "posts", DefaultQueries{})
	if err := codex.ensureDB(); !errors.Is(err, ErrDBNotConfigured) {
		t.Fatalf("expected ErrDBNotConfigured, got %v", err)
	}

	if err := validateDestination(nil); !errors.Is(err, ErrDestinationRequired) {
		t.Fatalf("expected destination error for nil, got %v", err)
	}

	var post struct{}
	if err := validateDestination(post); !errors.Is(err, ErrDestinationRequired) {
		t.Fatalf("expected destination error for non-pointer, got %v", err)
	}

	if err := validateDestination(&post); err != nil {
		t.Fatalf("expected pointer destination to succeed, got %v", err)
	}
}

func TestArgumentHelpers(t *testing.T) {
	if got := normalizeArgs(nil); got != nil {
		t.Fatalf("expected nil args, got %#v", got)
	}
	if got := normalizeArgs([]any{1}); got != 1 {
		t.Fatalf("expected single arg, got %#v", got)
	}
	if got := normalizeArgs([]any{1, 2}); len(got.([]any)) != 2 {
		t.Fatalf("expected variadic slice, got %#v", got)
	}

	fallback := "fallback"
	if got := normalizeLegacyExecArg(fallback, nil); got != fallback {
		t.Fatalf("expected fallback, got %#v", got)
	}
	if got := normalizeLegacyExecArg(fallback, []any{1}); got != 1 {
		t.Fatalf("expected single legacy arg, got %#v", got)
	}
	if got := normalizeLegacyExecArg(fallback, []any{1, 2}); len(got.([]any)) != 2 {
		t.Fatalf("expected legacy variadic slice, got %#v", got)
	}
}

func TestNamedParameterDetection(t *testing.T) {
	if !hasNamedParameters("SELECT * FROM posts WHERE id = :id") {
		t.Fatal("expected named parameter detection to succeed")
	}
	if hasNamedParameters("SELECT * FROM posts WHERE id = ?") {
		t.Fatal("did not expect named parameter detection for positional query")
	}
}

func TestLegacyBindingHelpers(t *testing.T) {
	codex := New(nil, "posts", DefaultQueries{})

	if _, err := codex.legacyResultType(); !errors.Is(err, ErrResultNotBound) {
		t.Fatalf("expected ErrResultNotBound, got %v", err)
	}

	codex.result = struct{}{}
	if _, err := codex.legacyResultType(); !errors.Is(err, ErrResultNotBound) {
		t.Fatalf("expected ErrResultNotBound for non-pointer prototype, got %v", err)
	}

	codex.Bind(&struct{ ID int }{})
	resultType, err := codex.legacyResultType()
	if err != nil {
		t.Fatalf("legacyResultType returned error: %v", err)
	}
	if resultType.Kind() != reflect.Struct {
		t.Fatalf("expected struct result type, got %s", resultType.Kind())
	}

	instance, err := codex.newLegacyResult()
	if err != nil {
		t.Fatalf("newLegacyResult returned error: %v", err)
	}
	if _, ok := instance.(*struct{ ID int }); !ok {
		t.Fatalf("unexpected instance type: %T", instance)
	}
}

func TestExecNamedNilPayload(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New returned error: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	codex := New(sqlxDB, "posts", DefaultQueries{})
	_, err = codex.execNamed(t.Context(), "INSERT INTO posts(title) VALUES (:title)", nil)
	if err == nil || !strings.Contains(err.Error(), "named exec requires a payload") {
		t.Fatalf("expected nil payload error, got %v", err)
	}
}
