package repositories

import (
	"consumer/internal/entities"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MetricsRepository struct {
	collection *mongo.Collection
}

func NewMetricsRepository(collection *mongo.Collection) *MetricsRepository {
	return &MetricsRepository{
		collection: collection,
	}
}

func (metricsRepository *MetricsRepository) Upsert(ctx context.Context, event entities.Event) error {
	update := bson.M{
		"$inc": bson.M{
			"count": 1,
			"sum":   event.Value,
		},
	}

	_, err := metricsRepository.collection.UpdateOne(
		ctx,
		bson.M{"eventType": event.Type},
		update,
		options.Update().SetUpsert(true),
	)

	return err
}
func (metricsRepository *MetricsRepository) FindAll(ctx context.Context) ([]entities.Metrics, error) {
	var metricsToReturn []entities.Metrics

	cursor, err := metricsRepository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var metric entities.MetricsSinceDayOne
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}

		if metric.Count > 0 {
			metric.AverageValue = metric.Sum / float64(metric.Count)
		} else {
			metric.AverageValue = 0
		}

		metricsToReturn = append(metricsToReturn, entities.Metrics{
			EventType:    metric.EventType,
			AverageValue: metric.AverageValue,
			Count:        metric.Count,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return metricsToReturn, nil
}
