# Migrating to Godex v2

`v2` is the first intentionally breaking cleanup release for Godex.

## Import Path

Change:

```go
import "github.com/sphireinc/godex"
```

To:

```go
import godex "github.com/sphireinc/godex/v2"
```

## Major Changes

- The Mantis dependency is removed.
- `Codex.DB` is now a `*sqlx.DB`.
- The preferred read API is the typed facade returned by `godex.For[T](...)`.
- Legacy helpers such as `SelectOne()`, `Select()`, `Insert()`, `Update()`, `Delete()`, and `SoftDelete()` remain for compatibility but are deprecated.

## Before

```go
q := Godex.CreateGodex(db, "posts", byID, one, many, insert, update, del, softDelete)
q.Bind(&Post{})
result, err := q.SelectOne(Godex.CxArgs{"id": 10})
```

## After

```go
q := godex.New(db, "posts", godex.DefaultQueries{
    SelectById: byID,
    SelectOne:  one,
    Select:     many,
    Insert:     insert,
    Update:     update,
    Delete:     del,
    SoftDelete: softDelete,
})

post, err := godex.For[Post](q).SelectOne(godex.Args{"id": 10})
```

## Write Operations

Move from implicit payloads to explicit payloads:

```go
id, err := q.InsertWith(struct {
    Title string `db:"title"`
}{
    Title: "Hello World",
})
```
