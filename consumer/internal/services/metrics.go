package services

import (
	"consumer/internal/entities"
	"consumer/internal/repositories"
	"context"
)

type MetricsService struct {
	metricsRepository *repositories.MetricsRepository
}

func NewMetricsService(metricsRepository *repositories.MetricsRepository) *MetricsService {
	return &MetricsService{
		metricsRepository: metricsRepository,
	}
}

func (metricsService *MetricsService) Fetch(ctx context.Context) ([]entities.Metrics, error) {
	return metricsService.metricsRepository.FindAll(ctx)
}
