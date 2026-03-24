package godex

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type testPost struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}

func TestSelectOneIntoSupportsNamedArgs(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(7, "hello")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(7).
		WillReturnRows(rows)

	var post testPost
	err := codex.SelectOneInto(&post, Args{"id": 7})
	if err != nil {
		t.Fatalf("SelectOneInto returned error: %v", err)
	}
	if post.ID != 7 || post.Title != "hello" {
		t.Fatalf("unexpected post: %+v", post)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestRawQueryPassesVariadicArgs(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM posts WHERE id = ? AND title = ?")).
		WithArgs(1, "hello").
		WillReturnRows(rows)

	result, err := codex.RawQuery("SELECT id FROM posts WHERE id = ? AND title = ?", 1, "hello")
	if err != nil {
		t.Fatalf("RawQuery returned error: %v", err)
	}
	result.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestLegacySelectOneUsesBoundPrototype(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	codex.Bind(&testPost{})
	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(9, "legacy")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(9).
		WillReturnRows(rows)

	result, err := codex.SelectOne(Args{"id": 9})
	if err != nil {
		t.Fatalf("SelectOne returned error: %v", err)
	}

	post, ok := result.(*testPost)
	if !ok {
		t.Fatalf("expected *testPost, got %T", result)
	}
	if post.ID != 9 || post.Title != "legacy" {
		t.Fatalf("unexpected post: %+v", post)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestInsertWithUsesNamedExec(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO posts (title) VALUES (?)")).
		WithArgs("new title").
		WillReturnResult(sqlmock.NewResult(42, 1))

	id, err := codex.InsertWith(struct {
		Title string `db:"title"`
	}{Title: "new title"})
	if err != nil {
		t.Fatalf("InsertWith returned error: %v", err)
	}
	if id != 42 {
		t.Fatalf("expected insert id 42, got %d", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestTypedModelSelectOneReturnsConcreteValue(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(13, "typed")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(13).
		WillReturnRows(rows)

	post, err := For[testPost](codex).SelectOne(Args{"id": 13})
	if err != nil {
		t.Fatalf("typed SelectOne returned error: %v", err)
	}
	if post.ID != 13 || post.Title != "typed" {
		t.Fatalf("unexpected post: %+v", post)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestQueryIntoSupportsNamedArgs(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "first").
		AddRow(2, "second")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("hello").
		WillReturnRows(rows)

	var posts []testPost
	err := codex.QueryInto(&posts, "SELECT id, title FROM posts WHERE title = :title", map[string]any{"title": "hello"})
	if err != nil {
		t.Fatalf("QueryInto returned error: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestQueryOneIntoSupportsPositionalArgsSlice(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(3, "positional")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ? AND title = ?")).
		WithArgs(3, "positional").
		WillReturnRows(rows)

	var post testPost
	err := codex.QueryOneInto(&post, "SELECT id, title FROM posts WHERE id = ? AND title = ?", []any{3, "positional"})
	if err != nil {
		t.Fatalf("QueryOneInto returned error: %v", err)
	}
	if post.ID != 3 {
		t.Fatalf("unexpected post: %+v", post)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestQueryOneContextIntoRequiresDestination(t *testing.T) {
	_, _, codex := newMockCodex(t)
	err := codex.QueryOneContextInto(context.Background(), nil, "SELECT 1", nil)
	if !errors.Is(err, ErrDestinationRequired) {
		t.Fatalf("expected ErrDestinationRequired, got %v", err)
	}
}

func TestQueryNamedHelpers(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	codex.RegisterQuery("by_title", "SELECT id, title FROM posts WHERE title = :title")

	oneRow := sqlmock.NewRows([]string{"id", "title"}).AddRow(4, "named")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("named").
		WillReturnRows(oneRow)

	var post testPost
	if err := codex.QueryNamedOneInto("by_title", &post, Args{"title": "named"}); err != nil {
		t.Fatalf("QueryNamedOneInto returned error: %v", err)
	}
	if post.Title != "named" {
		t.Fatalf("unexpected post: %+v", post)
	}

	manyRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(5, "named").
		AddRow(6, "named")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("named").
		WillReturnRows(manyRows)

	var posts []testPost
	if err := codex.QueryNamedContextInto(context.Background(), "by_title", &posts, Args{"title": "named"}); err != nil {
		t.Fatalf("QueryNamedContextInto returned error: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}

	manyRowsWrapper := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(7, "named").
		AddRow(8, "named")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("named").
		WillReturnRows(manyRowsWrapper)

	var postsViaWrapper []testPost
	if err := codex.QueryNamedInto("by_title", &postsViaWrapper, Args{"title": "named"}); err != nil {
		t.Fatalf("QueryNamedInto returned error: %v", err)
	}
	if len(postsViaWrapper) != 2 {
		t.Fatalf("expected 2 posts from wrapper, got %d", len(postsViaWrapper))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestSelectDefaultHelpers(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	byIDRows := sqlmock.NewRows([]string{"id", "title"}).AddRow(8, "by-id")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(8).
		WillReturnRows(byIDRows)

	var byID testPost
	if err := codex.SelectByIDInto(&byID, Args{"id": 8}); err != nil {
		t.Fatalf("SelectByIDInto returned error: %v", err)
	}

	listRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "one").
		AddRow(2, "two")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(listRows)

	var posts []testPost
	if err := codex.SelectContextInto(context.Background(), &posts, nil); err != nil {
		t.Fatalf("SelectContextInto returned error: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestWriteHelpers(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET title = ? WHERE id = ?")).
		WithArgs("updated", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updated, err := codex.UpdateWith(Args{"id": 10, "title": "updated"})
	if err != nil {
		t.Fatalf("UpdateWith returned error: %v", err)
	}
	if updated != 1 {
		t.Fatalf("expected 1 updated row, got %d", updated)
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM posts WHERE id = ?")).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := codex.DeleteWithContext(context.Background(), Args{"id": 10}); err != nil {
		t.Fatalf("DeleteWithContext returned error: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?")).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleted, err := codex.SoftDeleteWith(Args{"id": 10})
	if err != nil {
		t.Fatalf("SoftDeleteWith returned error: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 soft-deleted row, got %d", deleted)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestLegacyHelpers(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	codex.Bind(&testPost{})

	byIDRows := sqlmock.NewRows([]string{"id", "title"}).AddRow(11, "legacy-id")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(11).
		WillReturnRows(byIDRows)

	result, err := codex.SelectById(Args{"id": 11})
	if err != nil {
		t.Fatalf("SelectById returned error: %v", err)
	}
	if result.(*testPost).Title != "legacy-id" {
		t.Fatalf("unexpected legacy result: %+v", result)
	}

	listRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "one").
		AddRow(2, "two")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(listRows)

	values, err := codex.Select()
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	insertPayload := struct {
		ID    int    `db:"id"`
		Title string `db:"title"`
	}{ID: 12, Title: "legacy-write"}
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO posts (title) VALUES (?)")).
		WithArgs("legacy-write").
		WillReturnResult(sqlmock.NewResult(12, 1))
	if _, err := codex.Insert(insertPayload); err != nil {
		t.Fatalf("Insert returned error: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET title = ? WHERE id = ?")).
		WithArgs("legacy-write", 12).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if _, err := codex.Update(insertPayload); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM posts WHERE id = ?")).
		WithArgs(12).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := codex.Delete(Args{"id": 12}); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?")).
		WithArgs(12).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if _, err := codex.SoftDelete(Args{"id": 12}); err != nil {
		t.Fatalf("SoftDelete returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestLegacyHelpersRequireBind(t *testing.T) {
	_, _, codex := newMockCodex(t)

	if _, err := codex.SelectOne(Args{"id": 1}); !errors.Is(err, ErrResultNotBound) {
		t.Fatalf("expected ErrResultNotBound, got %v", err)
	}
	if _, err := codex.Select(); !errors.Is(err, ErrResultNotBound) {
		t.Fatalf("expected ErrResultNotBound, got %v", err)
	}
}

func TestTypedModelHelpers(t *testing.T) {
	db, mock, codex := newMockCodex(t)
	defer db.Close()

	model := For[testPost](codex)

	rawOne := sqlmock.NewRows([]string{"id", "title"}).AddRow(21, "query-one")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(21).
		WillReturnRows(rawOne)
	if _, err := model.QueryOne("SELECT id, title FROM posts WHERE id = :id", Args{"id": 21}); err != nil {
		t.Fatalf("QueryOne returned error: %v", err)
	}

	rawMany := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(22, "query-many").
		AddRow(23, "query-many")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("query-many").
		WillReturnRows(rawMany)
	values, err := model.QueryManyContext(context.Background(), "SELECT id, title FROM posts WHERE title = :title", Args{"title": "query-many"})
	if err != nil {
		t.Fatalf("QueryManyContext returned error: %v", err)
	}
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	rawManyWrapper := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(30, "query-many-wrapper").
		AddRow(31, "query-many-wrapper")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE title = ?")).
		WithArgs("query-many-wrapper").
		WillReturnRows(rawManyWrapper)
	values, err = model.QueryMany("SELECT id, title FROM posts WHERE title = :title", Args{"title": "query-many-wrapper"})
	if err != nil {
		t.Fatalf("QueryMany returned error: %v", err)
	}
	if len(values) != 2 {
		t.Fatalf("expected 2 values from QueryMany, got %d", len(values))
	}

	codex.RegisterQuery("typed_named", "SELECT id, title FROM posts WHERE id = :id")
	namedOne := sqlmock.NewRows([]string{"id", "title"}).AddRow(24, "typed-named")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(24).
		WillReturnRows(namedOne)
	if _, err := model.QueryNamedOne("typed_named", Args{"id": 24}); err != nil {
		t.Fatalf("QueryNamedOne returned error: %v", err)
	}

	namedMany := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(25, "typed-many").
		AddRow(26, "typed-many")
	codex.RegisterQuery("typed_named_many", "SELECT id, title FROM posts")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(namedMany)
	many, err := model.QueryNamedManyContext(context.Background(), "typed_named_many", nil)
	if err != nil {
		t.Fatalf("QueryNamedManyContext returned error: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("expected 2 values, got %d", len(many))
	}

	namedManyWrapper := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(32, "typed-many-wrapper").
		AddRow(33, "typed-many-wrapper")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(namedManyWrapper)
	many, err = model.QueryNamedMany("typed_named_many", nil)
	if err != nil {
		t.Fatalf("QueryNamedMany returned error: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("expected 2 values from QueryNamedMany, got %d", len(many))
	}

	byID := sqlmock.NewRows([]string{"id", "title"}).AddRow(27, "typed-id")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts WHERE id = ?")).
		WithArgs(27).
		WillReturnRows(byID)
	if _, err := model.SelectByID(Args{"id": 27}); err != nil {
		t.Fatalf("SelectByID returned error: %v", err)
	}

	selectRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(28, "typed-select").
		AddRow(29, "typed-select")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(selectRows)
	selected, err := model.SelectContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("SelectContext returned error: %v", err)
	}
	if len(selected) != 2 {
		t.Fatalf("expected 2 selected rows, got %d", len(selected))
	}

	selectRowsWrapper := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(34, "typed-select-wrapper").
		AddRow(35, "typed-select-wrapper")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM posts")).
		WillReturnRows(selectRowsWrapper)
	selected, err = model.Select(nil)
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if len(selected) != 2 {
		t.Fatalf("expected 2 selected rows from wrapper, got %d", len(selected))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func newMockCodex(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *Codex) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New returned error: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	codex := New(sqlxDB, "posts", DefaultQueries{
		SelectById: "SELECT id, title FROM posts WHERE id = :id",
		SelectOne:  "SELECT id, title FROM posts WHERE id = :id",
		Select:     "SELECT id, title FROM posts",
		Insert:     "INSERT INTO posts (title) VALUES (:title)",
		Update:     "UPDATE posts SET title = :title WHERE id = :id",
		Delete:     "DELETE FROM posts WHERE id = :id",
		SoftDelete: "UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id = :id",
	})

	return sqlxDB, mock, codex
}
