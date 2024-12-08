package main

import (
	"consumer/internal/controllers"
	"consumer/internal/repositories"
	"consumer/internal/services"
	"consumer/internal/store"
	"context"
	"log"
	"net/http"
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
	brokerAddress            string
	brokerPort               string
	topicName                string
	mongoURI                 string
	mongoDBName              string
	mongoDBCollection        string
	mongoDBMetricsCollection string
)

func init() {
	brokerAddress = os.Getenv("BROKER_ADDRESS")
	brokerPort = os.Getenv("BROKER_PORT")
	topicName = os.Getenv("TOPIC_NAME")
	mongoURI = os.Getenv("MONGO_URI")
	mongoDBName = os.Getenv("MONGO_DB_NAME")
	mongoDBCollection = os.Getenv("MONGO_DB_COLLECTION")
	mongoDBMetricsCollection = os.Getenv("MONGO_DB_METRICS_COLLECTION")

	if brokerAddress == "" || brokerPort == "" || topicName == "" || mongoURI == "" || mongoDBName == "" || mongoDBCollection == "" || mongoDBMetricsCollection == "" {
		log.Fatal("BROKER_ADDRESS | BROKER_PORT | TOPIC_NAME | MONGO vars not set in .env")
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
	metricsCollection := mongoStore.Database.Collection(mongoDBMetricsCollection)

	eventsRepository := repositories.NewEventsRepository(collection)
	eventsService := services.NewEventsService(eventsRepository)
	eventsController := controllers.NewEventsController(eventsService)

	metricsRepository := repositories.NewMetricsRepository(metricsCollection)
	metricsService := services.NewMetricsService(metricsRepository)
	metricsController := controllers.NewMetricsController(metricsService)

	brokerConnection := brokerAddress + ":" + brokerPort

	consumerGroup, err := services.SetupKafkaConsumer([]string{brokerConnection}, "events-consumer-group")
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	consumerHandler := services.NewConsumerGroupHandler(eventsRepository, metricsRepository)

	r := gin.Default()

	r.GET("/events", eventsController.Fetch)
	r.GET("/events/metrics-by-day", eventsController.GetEventMetricsByDay)
	r.GET("/metrics-since-day-one", metricsController.FetchSinceDayOne)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

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
