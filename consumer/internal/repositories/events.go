package repositories

import (
	"consumer/internal/entities"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventsRepository struct {
	collection *mongo.Collection
}

func NewEventsRepository(collection *mongo.Collection) *EventsRepository {
	return &EventsRepository{
		collection: collection,
	}
}

func (eventsRepository *EventsRepository) Insert(ctx context.Context, event entities.Event) error {
	_, err := eventsRepository.collection.InsertOne(ctx, event)
	return err
}

func (eventsRepository *EventsRepository) FindAll(ctx context.Context, page, limit int) ([]entities.Event, error) {
	var events []entities.Event

	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cur, err := eventsRepository.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var event entities.Event
		if err := cur.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (eventsRepository *EventsRepository) Count(ctx context.Context) (int64, error) {
	count, err := eventsRepository.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (eventsRepository *EventsRepository) GetEventMetricsByDay(ctx context.Context, start, end int64) ([]entities.Metrics, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: start},
				{Key: "$lte", Value: end},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "averageValue", Value: bson.D{{Key: "$avg", Value: "$value"}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "eventType", Value: "$_id"},
			{Key: "count", Value: 1},
			{Key: "averageValue", Value: bson.D{{Key: "$round", Value: bson.A{"$averageValue", 2}}}},
			{Key: "_id", Value: 0},
		}}},
	}

	cursor, err := eventsRepository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error on aggregate: %v", err)
	}
	defer cursor.Close(ctx)

	var metrics []entities.Metrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, fmt.Errorf("error on parse metrics %v", err)
	}

	return metrics, nil
}
