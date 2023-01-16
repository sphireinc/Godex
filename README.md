# Godex

Godex creates a simple interface for creating and story instances of reusable queries

### Example 

```go
q := Codex{
    Table: "posts",
    Queries: Queries{
        SelectById: "SELECT * FROM posts WHERE id = :id",
        SelectOne:  "SELECT * FROM posts WHERE id = :id and post_title = :post_title",
        Select:     "SELECT * FROM posts WHERE created_at > '2020-05-10 12:23:43'",
        Insert:     "INSERT INTO posts(post_id) VALUES (:post_id)",
        Update:     "UPDATE posts SET is_active=:is_active WHERE id=:id",
        Delete:     "DELETE FROM posts WHERE id=:id",
        SoftDelete: "UPDATE posts SET deleted_at=CURRENT_TIMESTAMP() WHERE id=:id",
    },
}

res, err := q.SelectOne()
if err != nil {
	panic(err)
}
fmt.Println(res.postId)
```