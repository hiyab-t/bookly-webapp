package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Books BookModel
}

func NewModels(dbpool *pgxpool.Pool) Models {
	return Models{
		Books: BookModel{DBpool: dbpool},
	}
}