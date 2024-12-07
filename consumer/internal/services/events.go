package services

import (
	"consumer/internal/entities"
	"consumer/internal/repositories"

	"github.com/gin-gonic/gin"
)

type EventsService struct {
	eventsRepository *repositories.EventsRepository
}

func NewEventsService(eventsRepository *repositories.EventsRepository) *EventsService {
	return &EventsService{
		eventsRepository: eventsRepository,
	}
}

func (eventsService *EventsService) Fetch(ctx *gin.Context, page, limit int) ([]entities.Event, int64, error) {
	events, err := eventsService.eventsRepository.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := eventsService.eventsRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return events, totalItems, nil
}
