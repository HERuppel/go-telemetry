package repositories

import (
	"consumer/internal/entities"
	"context"

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