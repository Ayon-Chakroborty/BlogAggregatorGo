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
	Id            int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserId        uuid.UUID
	LastFetchedAt time.Time
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
		feed.UserId,
	}

	err := m.DB.QueryRowContext(ctx, insertFeedQry, args...).Scan(
		&feed.Id,
		&feed.CreatedAt,
		&feed.UpdatedAt,
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
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
WHERE url = $1;`

func (m FeedsModel) Get(url string) (*Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Feed{}

	err := m.DB.QueryRowContext(ctx, getFeedQry, url).Scan(
		&f.Id,
		&f.CreatedAt,
		&f.UpdatedAt,
		&f.Name,
		&f.Url,
		&f.UserId,
		&f.LastFetchedAt,
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
			&feed.UserId,
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

const updateLastFetchedAtFeedQry = `
UPDATE feeds
SET last_fetched_at = $1
WHERE url = $2;`

func (m FeedsModel) UpdatedLastFetched(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, updateLastFetchedAtFeedQry, time.Now(), url)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

const getNextFeedToFetchQry = `
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;`

func (m FeedsModel) GetOrderedFeeds() ([]*Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, getNextFeedToFetchQry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds :=[]*Feed{}

	for rows.Next(){
		var feed Feed
		
		err := rows.Scan(
			&feed.Id,
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.Name,
			&feed.Url,
			&feed.UserId,
		)
		if err != nil{
			return nil, err
		}

		feed.LastFetchedAt = time.Now()

		feeds = append(feeds, &feed)
	}

	return feeds, nil
}
