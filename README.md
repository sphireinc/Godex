# Godex

[![Go](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/sphireinc/Godex/actions/workflows/go.yml/badge.svg)](https://github.com/sphireinc/Godex/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sphireinc/godex/v2.svg)](https://pkg.go.dev/github.com/sphireinc/godex/v2)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Godex is a small `sqlx`-based query store for reusable table-level queries.

- [API Reference](API.md)
- [Migration Guide](MIGRATING.md)

<div align="center">
    <img src="logo_grayscale.png" width="400px" alt="logo" />
</div>

## Example

### Creating a Codex
```go
import godex "github.com/sphireinc/godex/v2"

db := sqlx.MustConnect("mysql", dsn)

q := godex.NewWithQueries(db, "posts", godex.DefaultQueries{
    SelectById: "SELECT id, title, created_at FROM posts WHERE id = :id",
    SelectOne:  "SELECT id, title, created_at FROM posts WHERE id = :id AND title = :title",
    Select:     "SELECT id, title, created_at FROM posts WHERE created_at >= :created_after",
    Insert:     "INSERT INTO posts (title) VALUES (:title)",
    Update:     "UPDATE posts SET title = :title WHERE id = :id",
    Delete:     "DELETE FROM posts WHERE id = :id",
    SoftDelete: "UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id = :id",
}, map[string]string{
    "SelectUsersByFirstName": "SELECT id, first_name FROM users WHERE first_name = :first_name",
})

type Post struct {
    ID    int    `db:"id"`
    Title string `db:"title"`
}

posts := godex.For[Post](q)
```

### Reading one row

```go
post, err := posts.SelectOne(godex.Args{
    "id":    154,
    "title": "Hello World",
})
if err != nil {
    panic(err)
}

fmt.Println(post.ID)
```

### Reading many rows

```go
recentPosts, err := posts.Select(godex.Args{
    "created_after": time.Now().Add(-24 * time.Hour),
})
if err != nil {
    panic(err)
}
fmt.Println(len(recentPosts))
```

### Running a custom named query

```go
type User struct {
    ID        int    `db:"id"`
    FirstName string `db:"first_name"`
}

var users []User
err := q.QueryNamedInto("SelectUsersByFirstName", &users, godex.Args{
    "first_name": "John",
})
if err != nil {
    panic(err)
}
```

### Writing data

```go
id, err := q.InsertWith(struct {
    Title string `db:"title"`
}{
    Title: "Hello World",
})
if err != nil {
    panic(err)
}
fmt.Println(id)
```

## Notes for V2

- This module is `github.com/sphireinc/godex/v2`.
- Godex no longer depends on `mantis`.
- Queries can use either named parameters like `:id` or positional placeholders like `?`.
- Named parameters accept `godex.Args`, plain `map[string]any`, or structs with `db` tags.
- The preferred read API is the typed facade from `godex.For[T](...)`.
- Legacy helpers such as `SelectOne()` still exist for compatibility, but they are deprecated.

## Migrating From v1

- Update imports from `github.com/sphireinc/godex` to `github.com/sphireinc/godex/v2`.
- Replace Mantis-backed setup with `sqlx.DB`.
- Prefer `godex.For[T](...)` and the `*Into` / `*With` methods over the deprecated legacy helpers.
