# Godex API

This document describes the public API exposed by `github.com/sphireinc/godex/v2`.

## Import

```go
import godex "github.com/sphireinc/godex/v2"
```

## Core Types

### `type Codex`

`Codex` is the main query store. It holds:

- a table name
- a set of default CRUD-style queries
- a map of custom named queries
- a `*sqlx.DB` handle

Public fields:

```go
type Codex struct {
    Table          string
    DefaultQueries DefaultQueries
    Queries        map[string]string
    DB             *sqlx.DB
}
```

### `type DefaultQueries`

`DefaultQueries` stores the built-in query set used by the default read and write helpers.

```go
type DefaultQueries struct {
    SelectById string
    SelectOne  string
    Select     string
    Insert     string
    Update     string
    Delete     string
    SoftDelete string
}
```

### `type Args`

`Args` is a convenience alias for named query arguments.

```go
type Args map[string]any
```

Use it with queries that contain named placeholders such as `:id` or `:title`.

### `type CxArgs`

Deprecated compatibility alias for `Args`.

## Constructors

### `func New(db *sqlx.DB, table string, defaults DefaultQueries) *Codex`

Creates a new `Codex` with no custom queries.

### `func NewWithQueries(db *sqlx.DB, table string, defaults DefaultQueries, queries map[string]string) *Codex`

Creates a new `Codex` and preloads custom named queries.

### `func Open(driverName, dataSourceName, table string, defaults DefaultQueries, queries map[string]string) (*Codex, error)`

Opens a database using `sqlx.Open` and returns a `Codex`.

### `func OpenDB(driverName string, db *sql.DB, table string, defaults DefaultQueries, queries map[string]string) *Codex`

Wraps an existing `*sql.DB` in `sqlx` and returns a `Codex`.

### `func CreateGodex(...) Codex`

Deprecated compatibility constructor. Prefer `New` or `NewWithQueries`.

## Query Registration

### `func (c *Codex) RegisterQuery(name, query string)`

Adds or replaces a custom named query in `c.Queries`.

### `func (c *Codex) LookupQuery(name string) (string, error)`

Returns a custom named query by name.

Returns `ErrQueryNotFound` when the query does not exist.

## Raw Query APIs

### `func (c *Codex) RawQuery(query string, args ...any) (*sql.Rows, error)`

Runs a raw SQL query with positional arguments using `context.Background()`.

### `func (c *Codex) RawQueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)`

Same as `RawQuery`, but accepts a `context.Context`.

## Untyped Read APIs

These APIs scan directly into caller-provided destinations.

### `func (c *Codex) QueryOneInto(dest any, query string, arg any) error`

Runs any query and scans one row into `dest`.

### `func (c *Codex) QueryOneContextInto(ctx context.Context, dest any, query string, arg any) error`

Context-aware version of `QueryOneInto`.

### `func (c *Codex) QueryInto(dest any, query string, arg any) error`

Runs any query and scans many rows into `dest`.

### `func (c *Codex) QueryContextInto(ctx context.Context, dest any, query string, arg any) error`

Context-aware version of `QueryInto`.

### `func (c *Codex) QueryNamedOneInto(name string, dest any, arg any) error`

Looks up a registered custom query and scans one row into `dest`.

### `func (c *Codex) QueryNamedOneContextInto(ctx context.Context, name string, dest any, arg any) error`

Context-aware version of `QueryNamedOneInto`.

### `func (c *Codex) QueryNamedInto(name string, dest any, arg any) error`

Looks up a registered custom query and scans many rows into `dest`.

### `func (c *Codex) QueryNamedContextInto(ctx context.Context, name string, dest any, arg any) error`

Context-aware version of `QueryNamedInto`.

## Default Read APIs

These methods use the query strings stored in `DefaultQueries`.

### `func (c *Codex) SelectByIDInto(dest any, arg any) error`

Runs `DefaultQueries.SelectById` and scans one row into `dest`.

### `func (c *Codex) SelectByIDContextInto(ctx context.Context, dest any, arg any) error`

Context-aware version of `SelectByIDInto`.

### `func (c *Codex) SelectOneInto(dest any, arg any) error`

Runs `DefaultQueries.SelectOne` and scans one row into `dest`.

### `func (c *Codex) SelectOneContextInto(ctx context.Context, dest any, arg any) error`

Context-aware version of `SelectOneInto`.

### `func (c *Codex) SelectInto(dest any, arg any) error`

Runs `DefaultQueries.Select` and scans many rows into `dest`.

### `func (c *Codex) SelectContextInto(ctx context.Context, dest any, arg any) error`

Context-aware version of `SelectInto`.

## Write APIs

These methods execute the query strings stored in `DefaultQueries`.

### `func (c *Codex) InsertWith(arg any) (int64, error)`

Runs `DefaultQueries.Insert` using named execution and returns the last insert ID.

### `func (c *Codex) InsertWithContext(ctx context.Context, arg any) (int64, error)`

Context-aware version of `InsertWith`.

### `func (c *Codex) UpdateWith(arg any) (int64, error)`

Runs `DefaultQueries.Update` and returns the number of affected rows.

### `func (c *Codex) UpdateWithContext(ctx context.Context, arg any) (int64, error)`

Context-aware version of `UpdateWith`.

### `func (c *Codex) DeleteWith(arg any) error`

Runs `DefaultQueries.Delete`.

### `func (c *Codex) DeleteWithContext(ctx context.Context, arg any) error`

Context-aware version of `DeleteWith`.

### `func (c *Codex) SoftDeleteWith(arg any) (int64, error)`

Runs `DefaultQueries.SoftDelete` and returns the number of affected rows.

### `func (c *Codex) SoftDeleteWithContext(ctx context.Context, arg any) (int64, error)`

Context-aware version of `SoftDeleteWith`.

## Typed Read API

The preferred read API is the typed facade returned by `godex.For[T](...)`.

### `type Model[T any]`

`Model[T]` wraps a `*Codex` and returns concrete typed values instead of requiring destination pointers.

### `func For[T any](c *Codex) Model[T]`

Creates a typed API for model `T`.

### `func (m Model[T]) QueryOne(query string, arg any) (T, error)`

Runs any query and returns one value of type `T`.

### `func (m Model[T]) QueryOneContext(ctx context.Context, query string, arg any) (T, error)`

Context-aware version of `QueryOne`.

### `func (m Model[T]) QueryMany(query string, arg any) ([]T, error)`

Runs any query and returns many values of type `T`.

### `func (m Model[T]) QueryManyContext(ctx context.Context, query string, arg any) ([]T, error)`

Context-aware version of `QueryMany`.

### `func (m Model[T]) QueryNamedOne(name string, arg any) (T, error)`

Runs a registered custom query and returns one value of type `T`.

### `func (m Model[T]) QueryNamedOneContext(ctx context.Context, name string, arg any) (T, error)`

Context-aware version of `QueryNamedOne`.

### `func (m Model[T]) QueryNamedMany(name string, arg any) ([]T, error)`

Runs a registered custom query and returns many values of type `T`.

### `func (m Model[T]) QueryNamedManyContext(ctx context.Context, name string, arg any) ([]T, error)`

Context-aware version of `QueryNamedMany`.

### `func (m Model[T]) SelectByID(arg any) (T, error)`

Runs `DefaultQueries.SelectById` and returns one value of type `T`.

### `func (m Model[T]) SelectByIDContext(ctx context.Context, arg any) (T, error)`

Context-aware version of `SelectByID`.

### `func (m Model[T]) SelectOne(arg any) (T, error)`

Runs `DefaultQueries.SelectOne` and returns one value of type `T`.

### `func (m Model[T]) SelectOneContext(ctx context.Context, arg any) (T, error)`

Context-aware version of `SelectOne`.

### `func (m Model[T]) Select(arg any) ([]T, error)`

Runs `DefaultQueries.Select` and returns many values of type `T`.

### `func (m Model[T]) SelectContext(ctx context.Context, arg any) ([]T, error)`

Context-aware version of `Select`.

## Legacy Compatibility APIs

These methods remain available for older callers, but they are deprecated.

### `func (c *Codex) Bind(result any) *Codex`

Registers a pointer prototype used by the legacy read helpers.

### `func (c *Codex) SelectById(args ...any) (any, error)`

Deprecated. Prefer `SelectByIDInto` or `For[T](...).SelectByID`.

### `func (c *Codex) SelectOne(args ...any) (any, error)`

Deprecated. Prefer `SelectOneInto` or `For[T](...).SelectOne`.

### `func (c *Codex) Select(args ...any) ([]any, error)`

Deprecated. Prefer `SelectInto` or `For[T](...).Select`.

### `func (c *Codex) Insert(args ...any) (int64, error)`

Deprecated. Prefer `InsertWith`.

### `func (c *Codex) Update(args ...any) (int64, error)`

Deprecated. Prefer `UpdateWith`.

### `func (c *Codex) Delete(args ...any) error`

Deprecated. Prefer `DeleteWith`.

### `func (c *Codex) SoftDelete(args ...any) (int64, error)`

Deprecated. Prefer `SoftDeleteWith`.

## Utility Function

### `func Pretty(str string) string`

Formats a JSON string with indentation. Returns an empty string if the input is not valid JSON.

## Errors

### `ErrDBNotConfigured`

Returned when a `Codex` has no configured database handle.

### `ErrResultNotBound`

Returned by legacy read helpers when `Bind(...)` has not been called.

### `ErrQueryNotFound`

Returned when a named query lookup fails.

### `ErrDestinationRequired`

Returned when a destination argument is nil or not a pointer.

## Argument Rules

- If a query uses named parameters like `:id`, pass `godex.Args`, `map[string]any`, or a struct.
- If a query uses positional placeholders like `?`, pass a single value or a slice of arguments where appropriate.
- Read helpers that end in `Into` require a non-nil pointer destination.
- Write helpers ending in `With` expect an explicit payload object.

## Minimal Example

```go
type Post struct {
    ID    int    `db:"id"`
    Title string `db:"title"`
}

db := sqlx.MustConnect("mysql", dsn)

q := godex.New(db, "posts", godex.DefaultQueries{
    SelectById: "SELECT id, title FROM posts WHERE id = :id",
    SelectOne:  "SELECT id, title FROM posts WHERE id = :id",
    Select:     "SELECT id, title FROM posts",
    Insert:     "INSERT INTO posts (title) VALUES (:title)",
    Update:     "UPDATE posts SET title = :title WHERE id = :id",
    Delete:     "DELETE FROM posts WHERE id = :id",
    SoftDelete: "UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id = :id",
})

post, err := godex.For[Post](q).SelectOne(godex.Args{"id": 1})
if err != nil {
    panic(err)
}

fmt.Println(post.Title)
```
