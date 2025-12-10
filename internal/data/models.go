package data

import "database/sql"

// parent model to hold all models
type Models struct {
	UsersModel       UsersModel
	FeedsModel       FeedsModel
	FeedFollowsModel FeedFollowsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UsersModel:       UsersModel{DB: db},
		FeedsModel:       FeedsModel{DB: db},
		FeedFollowsModel: FeedFollowsModel{DB: db},
	}
}
