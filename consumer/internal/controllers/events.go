package controllers

import (
	"consumer/internal/entities"
	"consumer/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventsController struct {
	eventService *services.EventsService
}

func NewEventsController(eventsService *services.EventsService) *EventsController {
	return &EventsController{
		eventService: eventsService,
	}
}

// Fetch returns stored events received with Kafka
// @Summary Get events
// @Description Retrieve a paginated list of events stored in the database.
// @Tags Events
// @Accept json
// @Produce json
// @Param limit query int true "Number of items per page"
// @Param page query int true "Page number"
// @Success 200 {array} entities.EventsResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events [get]
func (eventsController *EventsController) Fetch(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PAGE"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_LIMIT"})
		return
	}

	events, totalItems, err := eventsController.eventService.Fetch(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := entities.EventsResponse{
		Page:   int64(pageInt),
		Limit:  int64(limitInt),
		Count:  totalItems,
		Events: events,
	}

	ctx.JSON(http.StatusOK, response)
}
