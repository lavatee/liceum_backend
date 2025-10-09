package service

import "github.com/lavatee/liceum_backend/internal/model"

type Events interface {
	CreateEvent(event model.Event) (int, error)
	CreateEventBlocks(blocks []model.EventBlock) error
	GetCurrentEvents() ([]model.Event, error)
	GetAllTasks() ([]model.Event, error)
}
