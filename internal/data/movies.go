package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Always hide
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // Hide if empty
	Runtime   Runtime   `json:"runtime,omitempty"` // Hide if empty
	Genres    []string  `json:"genres,omitempty"`  // Hide if empty
	Version   int32     `json:"version"`
}
