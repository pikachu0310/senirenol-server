package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// InputType is the input device enum
// 0 = keyboard, 1 = button
type InputType uint8

const (
	InputKeyboard InputType = 0
	InputButton   InputType = 1
)

// Difficulty enum: 0=Past,1=Present,2=Future,3=Lycoris,4=Parallel
type Difficulty uint8

const (
	DiffPast Difficulty = iota
	DiffPresent
	DiffFuture
	DiffLycoris
	DiffParallel
)
