package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"producer/entities"
	"time"

	"github.com/IBM/sarama"
)

var (
	brokerAddress string
	topicName     string
	eventTypes    = []string{
		"VEHICLE_SPEED",
		"ENGINE_RPM",
		"ENGINE_TEMPERATURE",
		"FUEL_LEVEL",
		"DISTANCE_TRAVELED",
		"GPS_LOCATION",
		"LIGHT_STATUS",
	}
)

func generateRandomEvent() entities.Event {
	return entities.Event{
		Type:      eventTypes[rand.Intn(len(eventTypes))],
		Timestamp: time.Now().Unix(),
		Value:     rand.Float64() * 100,
	}
}

func createProducer(brokers []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func sendEvent(producer sarama.SyncProducer, topic string, event entities.Event) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Published event:\nType: %s\nTimestamp: %d\nValue: %.2f \n\n", event.Type, event.Timestamp, event.Value)

	return nil
}

func init() {
	brokerAddress = os.Getenv("BROKER_ADDRESS")
	topicName = os.Getenv("TOPIC_NAME")

	if brokerAddress == "" || topicName == "" {
		log.Fatal("BROKER_ADDRESS or TOPIC_NAME not set in .env")
	}
}

func main() {
	producer, err := createProducer([]string{brokerAddress})
	if err != nil {
		log.Fatalf("Failed to producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing producer: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Starting event producer...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			log.Println("Shutting down...")
			return
		case <-ticker.C:
			event := generateRandomEvent()
			if err := sendEvent(producer, topicName, event); err != nil {
				log.Printf("Failed to send event: %v", err)
			}
		}
	}
}
