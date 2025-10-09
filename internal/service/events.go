package service

import "github.com/lavatee/liceum_backend/internal/repository"

type EventsService struct {
	repo *repository.Repository
}

func NewEventsService(repo *repository.Repository) *EventsService {
	return &EventsService{
		repo: repo,
	}
}
