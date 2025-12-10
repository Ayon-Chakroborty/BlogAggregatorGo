package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FeedsModel struct {
	DB *sql.DB
}

type Feed struct {
	Id         int
	Created_at time.Time
	Updated_at time.Time
	Name       string
	Url        string
	User_id    uuid.UUID
}

const insertFeedQry = `
INSERT INTO feeds (name, url, user_id)
VALUES(
	$1,
	$2,
	$3
)

RETURNING id, created_at, updated_at;`

func (m FeedsModel) Insert(feed *Feed) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		feed.Name,
		feed.Url,
		feed.User_id,
	}

	err := m.DB.QueryRowContext(ctx, insertFeedQry, args...).Scan(
		&feed.Id,
		&feed.Created_at,
		&feed.Updated_at,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "feeds_url_key"`:
			return ErrDuplicateUrl
		default:
			return err
		}
	}

	return nil
}

const getFeedQry = `
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds
WHERE url = $1;`

func (m FeedsModel) Get(url string) (*Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Feed{}

	err := m.DB.QueryRowContext(ctx, getFeedQry, url).Scan(
		&f.Id,
		&f.Created_at,
		&f.Updated_at,
		&f.Name,
		&f.Url,
		&f.User_id,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return f, nil
}

const getAllFeedsQry = `SELECT name, url, user_id FROM feeds;`

func (m FeedsModel) GetAll() ([]*Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feeds := []*Feed{}

	rows, err := m.DB.QueryContext(ctx, getAllFeedsQry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feed Feed

		err := rows.Scan(
			&feed.Name,
			&feed.Url,
			&feed.User_id,
		)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, &feed)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
