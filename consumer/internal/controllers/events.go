package controllers

import (
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

	events, err := eventsController.eventService.Fetch(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
