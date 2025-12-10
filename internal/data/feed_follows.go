package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FeedFollowsModel struct {
	DB *sql.DB
}

type FeedFollow struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uuid.UUID
	FeedId    int
	UserName  string
	FeedName  string
	Url       string
}

const insertFeedFollowsQry = `
WITH new_feed_follow AS(
	INSERT INTO feed_follows (user_id, feed_id) 
	VALUES (
		$1,
		$2
	)
	RETURNING *
)

SELECT 
	new_feed_follow.*,
	users.name as user_name,
	feeds.name as feed_name,
	feeds.url as feed_url
FROM new_feed_follow
INNER JOIN users ON users.id = new_feed_follow.user_id
INNER JOIN feeds on feeds.id = new_feed_follow.feed_id;`

func (m FeedFollowsModel) Insert(f *FeedFollow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		f.UserId,
		f.FeedId,
	}

	err := m.DB.QueryRowContext(ctx, insertFeedFollowsQry, args...).Scan(
		&f.Id,
		&f.CreatedAt,
		&f.UpdatedAt,
		&f.UserId,
		&f.FeedId,
		&f.UserName,
		&f.FeedName,
		&f.Url,
	)

	if err != nil {
		return err
	}

	return nil
}

const getFeedFollowUserQry = `
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_follows.feed_id,
	feeds.name as feed_name,
	feeds.url as feed_url
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;`

func (m FeedFollowsModel) Get(userId uuid.UUID) ([]*FeedFollow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, getFeedFollowUserQry, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*FeedFollow{}

	for rows.Next() {
		feed := &FeedFollow{}

		err := rows.Scan(
			&feed.Id,
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.UserId,
			&feed.FeedId,
			&feed.FeedName,
			&feed.Url,
		)

		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
