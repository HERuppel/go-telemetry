package services

import (
	"consumer/internal/entities"
	"consumer/internal/repositories"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	eventsRepository  *repositories.EventsRepository
	metricsRepository *repositories.MetricsRepository
}

func NewConsumerGroupHandler(eventsRepository *repositories.EventsRepository, metricsRepository *repositories.MetricsRepository) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		eventsRepository:  eventsRepository,
		metricsRepository: metricsRepository,
	}
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		var event entities.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		log.Printf("Received event:\n Type: %s\n Timestamp: %d\n Value: %.2f\n\n", event.Type, event.Timestamp, event.Value)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if err := consumerGroupHandler.eventsRepository.Insert(ctx, event); err != nil {
			log.Printf("Error inserting into MongoDB: %v", err)
		}

		if err := consumerGroupHandler.metricsRepository.Upsert(ctx, event); err != nil {
			log.Printf("Error upserting into MongoDB: %v", err)
		}

		cancel()

		sess.MarkMessage(msg, "")
	}
	return nil
}

func SetupKafkaConsumer(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	const (
		maxRetries    = 10
		retryInterval = 5
	)
	var consumerGroup sarama.ConsumerGroup
	var err error

	for i := 0; i < maxRetries; i++ {
		consumerGroup, err = sarama.NewConsumerGroup(brokers, groupID, config)
		if err == nil {
			log.Println("Kafka consumer connected successfully!")
			return consumerGroup, nil
		}

		log.Printf("Error connecting to Kafka: %v. Retrying in %d seconds...", err, retryInterval)
		time.Sleep(retryInterval * time.Second)
	}

	return nil, err
}
