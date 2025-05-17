package main

import (
	"context"
	"paymentfc/cmd/payment/handler"
	"paymentfc/cmd/payment/repository"
	"paymentfc/cmd/payment/resource"
	"paymentfc/cmd/payment/service"
	"paymentfc/cmd/payment/usecase"
	"paymentfc/config"
	"paymentfc/infrastructure/constant"
	"paymentfc/infrastructure/logger"
	"paymentfc/kafka"
	"paymentfc/models"
	"paymentfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	//redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)
	kafkaWriter := kafka.NewWriter(cfg.Kafka.Broker, cfg.Kafka.KafkaTopics[constant.KafkaTopicPaymentSuccess])
	logger.SetupLogger()

	databaseRepository := repository.NewPaymentDatabase(db)
	publisherRepository := repository.NewKafkaPublisher(kafkaWriter)

	paymentService := service.NewPaymentService(databaseRepository, publisherRepository)
	paymentUsecase := usecase.NewPaymentUseCase(paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)
	xenditRepository := repository.NewXenditClient(cfg.Xendit.XenditAPIKey)
	xenditService := service.NewXenditService(databaseRepository, xenditRepository)
	xenditUsecase := usecase.NewXenditUseCase(xenditService)

	// kafka consumer
	kafka.StartOrderConsumer(cfg.Kafka.Broker, cfg.Kafka.KafkaTopics[constant.KafkaTopicOrderCreated], func(event models.OrderCreatedEvent) {
		if err := xenditUsecase.CreateInvoice(context.Background(), event); err != nil {
			logger.Logger.Printf("Failed Handling order created event: ", err.Error())
		}
	})

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, paymentHandler)
	router.Run(":" + port)
	logger.Logger.Printf("Server running on port: %s", port)

}
