package data

import "database/sql"

// parent model to hold all models
type Models struct {
	UsersModel UsersModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UsersModel: UsersModel{DB: db},
	}
}
