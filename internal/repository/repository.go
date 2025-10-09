package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/liceum_backend/internal/model"
)

type Events interface {
	CreateEvent(event model.Event) (int, error)
	CreateEventBlocks(blocks []model.EventBlock) error
}

type Repository struct {
	Events
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Events: NewEventsPostgres(db),
	}
}
