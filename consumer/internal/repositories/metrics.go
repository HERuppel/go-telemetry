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
