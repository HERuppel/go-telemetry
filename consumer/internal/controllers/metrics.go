package controllers

import (
	"consumer/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricsController struct {
	metricsService *services.MetricsService
}

func NewMetricsController(metricsService *services.MetricsService) *MetricsController {
	return &MetricsController{
		metricsService: metricsService,
	}
}

// FetchSinceDayOne returns metrics aggregated by event type since day one
// @Summary Get metrics since application day one
// @Description Retrieves aggregated metrics since application day one
// @Tags Metrics
// @Accept json
// @Produce json
// @Success 200 {array} entities.MetricsResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /metrics-since-day-one [get]
func (metricsController *MetricsController) FetchSinceDayOne(ctx *gin.Context) {
	metrics, err := metricsController.metricsService.Fetch(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}
