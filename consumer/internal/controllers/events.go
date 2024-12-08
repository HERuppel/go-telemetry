package controllers

import (
	"consumer/internal/entities"
	"consumer/internal/services"
	"consumer/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EventsController struct {
	eventsService *services.EventsService
}

func NewEventsController(eventsService *services.EventsService) *EventsController {
	return &EventsController{
		eventsService: eventsService,
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

	events, totalItems, err := eventsController.eventsService.Fetch(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	var response entities.EventsResponse

	if len(events) == 0 {
		response = entities.EventsResponse{
			Page:   int64(pageInt),
			Limit:  int64(limitInt),
			Count:  totalItems,
			Events: []entities.Event{},
		}

	} else {
		response = entities.EventsResponse{
			Page:   int64(pageInt),
			Limit:  int64(limitInt),
			Count:  totalItems,
			Events: events,
		}

	}

	ctx.JSON(http.StatusOK, response)
}

// GetMetricsByDay returns some event metrics by a given date
// @Summary Get metrics by day
// @Description Retrieves metrics reading the db by a given date
// @Tags Metrics
// @Accept json
// @Produce json
// @Param date query string false "Date to filter metrics"
// @Success 200 {array} entities.Metrics
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/metrics-by-day [get]
func (eventsController *EventsController) GetEventMetricsByDay(ctx *gin.Context) {
	date := ctx.DefaultQuery("date", time.Now().Add(-3*time.Hour).Format("2006-01-02"))

	start, end, err := utils.GetUnixStartAndEndOfDay(date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	metrics, err := eventsController.eventsService.GetEventMetrics(ctx, start, end)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	var response []entities.Metrics

	if len(metrics) == 0 {
		response = []entities.Metrics{}
	} else {
		response = metrics
	}

	ctx.JSON(http.StatusOK, response)
}
