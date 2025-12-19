package data

import (
	"context"
	"database/sql"
	"time"
)

type PostsModels struct {
	DB *sql.DB
}

type Post struct {
	Id          int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedId      int
}

const insertPostQry = `
INSERT INTO posts (title, url, description, published_at, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	(SELECT feeds.id FROM feeds WHERE feeds.url = $5)
)
	
RETURNING id, created_at, updated_at;`

func (m PostsModels) Insert(post *Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		post.Title,
		post.Url,
		post.Description,
		post.PublishedAt,
	}

	err := m.DB.QueryRowContext(ctx, insertPostQry, args...).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "posts_url_key"`:
			return ErrDuplicateUrl
		default:
			return err
		}
	}

	return nil
}
