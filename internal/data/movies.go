package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.alexedwards.net/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Always hide
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // Hide if empty
	Runtime   Runtime   `json:"runtime,omitempty"` // Hide if empty
	Genres    []string  `json:"genres,omitempty"`  // Hide if empty
	Version   int32     `json:"version"`
}

// Define a MovieModel struct which wraps a sql.DB connection pool
type MovieModel struct {
	DB *sql.DB
}

// The Insert() method accepts a pointer to a Movie struct, which contains the data for the new record.
func (m *MovieModel) Insert(movie *Movie) error {
	// SQL query for inserting and returning system-generated values
	query := `
    INSERT INTO movies (title, year, runtime, genres)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, version`
	// values for the placeholder taken from movie struct
	// Stored in a slice to make it clear which values match which placeholders
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	// execute the query and stored the returned value in the same movie struct
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) Get(id int64) (*Movie, error) {
	// If the ID is less than 1, we skip the database call and return ErrRecordNotFound immediately
	if id < 1 {
		return nil, ErrorRecordNotFound
	}

	// SQL query to retrieve the movie
	query := `
    SELECT id, created_at, title, year, runtime, genres, version
    FROM movies
    WHERE id = $1`

	// create a movie struct to store the result
	var movie Movie

	// run the query with QueryRow() and scan the result into the Movie struct
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		// pq.Array() acts as an adapter that helps Go read PostgreSQL arrays into Go slices.
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	// handles error
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorRecordNotFound
		default:
			return nil, err
		}
	}
	// if movie found return pointer to the movie struct
	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	// SQL query to update the movie and return the new version number
	query := `
    UPDATE movies
    SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
    WHERE id = $5
    RETURNING version`

	// values for the placeholder
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	// Execute the query, then scan the new version into movie.Version
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m *MovieModel) Delete(id int64) error {
	// Return ErrRecordNotFound if the ID is invalid.
	if id < 1 {
		return ErrorRecordNotFound
	}

	// SQL query to delete the record
	query := `DELETE FROM movies WHERE id = $1`

	// Execute the query ( Exec() gives a sql.Result object, which tells us how many rows were affected )
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the movie didn't exist
	if rowsAffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	return nil
}

// collect the movie validation rules in ValidateMovie() function for reusing
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
