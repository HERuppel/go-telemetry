package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"producer/entities"
	"time"

	"github.com/IBM/sarama"
)

var (
	brokerAddress string
	brokerPort    string
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

const (
	dispatchEventInterval     = 5
	pingConsumerRetryInterval = 5
	pingConsumerMaxRetries    = 10
)

func generateRandomEvent() entities.Event {
	return entities.Event{
		Type:      eventTypes[rand.Intn(len(eventTypes))],
		Timestamp: time.Now().Add(-3 * time.Hour).Unix(),
		Value:     rand.Float64() * 100,
	}
}

func createProducer(brokers []string) (sarama.SyncProducer, error) {
	var producer sarama.SyncProducer
	var err error
	const (
		maxRetries    = 10
		retryInterval = 5
	)

	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer(brokers, nil)
		if err == nil {
			log.Println("Producer created successfully!")
			return producer, nil
		}

		log.Printf("Error connecting to Kafka: %v. Trying again  in %d seconds...", err, retryInterval)
		time.Sleep(retryInterval * time.Second)
	}

	return nil, fmt.Errorf("não foi possível criar o producer após %d tentativas: %v", maxRetries, err)
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

	log.Printf("Published event:\nType: %s\nTimestamp: %d\nValue: %f \n\n", event.Type, event.Timestamp, event.Value)

	return nil
}

func pingConsumer() bool {
	resp, err := http.Get("http://consumer:3333/ping")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func waitForConsumer() bool {
	for i := 0; i < pingConsumerMaxRetries; i++ {
		if pingConsumer() {
			log.Println("Consumer is up!")
			return true
		}

		log.Printf("Consumer not available. Retrying in %d seconds...\n", pingConsumerRetryInterval)
		time.Sleep(pingConsumerRetryInterval * time.Second)
	}

	log.Println("Failed to reach consumer, max number of attempts.")
	return false
}

func init() {
	brokerAddress = os.Getenv("BROKER_ADDRESS")
	brokerPort = os.Getenv("BROKER_PORT")
	topicName = os.Getenv("TOPIC_NAME")

	if brokerAddress == "" || brokerPort == "" || topicName == "" {
		log.Fatal("BROKER_ADDRESS | BROKER_PORT | TOPIC_NAME not set in .env")
	}
}

func main() {
	if !waitForConsumer() {
		log.Fatalf("Couldn't reach consumer...")
	}

	brokerConnection := brokerAddress + ":" + brokerPort

	producer, err := createProducer([]string{brokerConnection})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing producer: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Starting event producer...")
	ticker := time.NewTicker(dispatchEventInterval * time.Second)
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
