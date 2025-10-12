package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/liceum_backend/internal/model"
)

type Events interface {
	CreateEvent(event model.Event) (int, error)
	DeleteEvent(eventId int) error
	CreateEventBlocks(blocks []model.EventBlock) error
	DeleteEventBlock(blockId int) error
	EditEventInfo(event model.Event) error
	EditBlockInfo(block model.EventBlock) error
	GetCurrentEvents() ([]model.Event, error)
	GetAllEvents() ([]model.Event, error)
	GetOneEvent(eventId int) (model.Event, error)
	GetOneBlock(blockId int) (model.EventBlock, error)
}

type Repository struct {
	Events
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Events: NewEventsPostgres(db),
	}
}
