package models

// Post represents a blog post.
type Post struct {
	ID      int    `db:"id"`
	Content string `db:"content"`
	Author  string `db:"author"`
}
