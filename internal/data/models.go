package data

import (
	"database/sql"
	"errors"
)

var ErrorRecordNotFound = errors.New("record not found")

// Define Models struct which wraps all the database models
type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
}

// NewModels returns a Models struct with the real MovieModel
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: &MovieModel{DB: db}, // use pointer to match method receivers
	}
}

// newMockModels returns a Models struct with the mock MovieModel
func newMockModels() Models {
	return Models{
		Movies: MockMovieModel{},
	}
}
