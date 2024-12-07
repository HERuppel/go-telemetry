package main

import (
	"consumer/internal/controllers"
	"consumer/internal/repositories"
	"consumer/internal/services"
	"consumer/internal/store"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "consumer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	brokerAddress     string
	topicName         string
	mongoURI          string
	mongoDBName       string
	mongoDBCollection string
)

func init() {
	brokerAddress = os.Getenv("BROKER_ADDRESS")
	topicName = os.Getenv("TOPIC_NAME")
	mongoURI = os.Getenv("MONGO_URI")
	mongoDBName = os.Getenv("MONGO_DB_NAME")
	mongoDBCollection = os.Getenv("MONGO_DB_COLLECTION")

	if brokerAddress == "" || topicName == "" || mongoURI == "" || mongoDBName == "" || mongoDBCollection == "" {
		log.Fatal("BROKER_ADDRESS | TOPIC_NAME | MONGO vars not set in .env")
	}
}

// @title           Go Telemtry
// @version         1.0
// @description     API that consumes, stores and returns events received with Kafka

// @host      localhost:3333
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	mongoStore := store.NewMongoStore(mongoURI, mongoDBName)
	defer mongoStore.Close()

	collection := mongoStore.Database.Collection(mongoDBCollection)

	eventsRepository := repositories.NewEventsRepository(collection)
	eventsService := services.NewEventsService(eventsRepository)
	eventsController := controllers.NewEventsController(eventsService)

	consumerGroup, err := services.SetupKafkaConsumer([]string{brokerAddress}, "events-consumer-group")
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	consumerHandler := services.NewConsumerGroupHandler(eventsRepository)

	r := gin.Default()

	r.GET("/events", eventsController.Fetch)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go r.Run(":3333")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topicName}, consumerHandler); err != nil {
				log.Printf("Error during Kafka consumption: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan

	log.Println("Shutting down...")
	cancel()
}
