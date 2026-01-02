package data

import (
	"context"
	"database/sql"
	"time"
)

type PostsModel struct {
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

func (m PostsModel) Insert(post *Post, feedUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		post.Title,
		post.Url,
		post.Description,
		post.PublishedAt,
		feedUrl,
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

const getPostQry = `
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id
FROM posts
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
LIMIT $2;`

func (m PostsModel) GetAll(user *User, limit int) ([]*Post, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, getPostQry, user.Id, limit)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.Id,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Title,
			&post.Url,
			&post.Description,
			&post.PublishedAt,
			&post.FeedId,
		)

		if err != nil{
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil{
		return nil, err
	}

	return posts, nil
}