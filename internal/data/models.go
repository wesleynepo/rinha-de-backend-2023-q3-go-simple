package data

import (
    "database/sql"
    "errors"
)


var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
    People PeopleModel
}

func NewModels(db *sql.DB) Models {
    return Models{
        People: PeopleModel{DB: db},
    }
}
