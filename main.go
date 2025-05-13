package main

import (
	"paymentfc/cmd/payment/handler"
	"paymentfc/cmd/payment/repository"
	"paymentfc/cmd/payment/resource"
	"paymentfc/cmd/payment/service"
	"paymentfc/cmd/payment/usecase"
	"paymentfc/config"
	"paymentfc/infrastructure/constant"
	"paymentfc/infrastructure/logger"
	"paymentfc/kafka"
	"paymentfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	//redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)
	kafkaWriter := kafka.NewWriter(cfg.Kafka.Broker, cfg.Kafka.KafkaTopics[constant.KafkaTopicPaymentSuccess])

	databaseRepository := repository.NewPaymentDatabase(db)
	publisherRepository := repository.NewKafkaPublisher(kafkaWriter)
	paymentService := service.NewPaymentService(databaseRepository, publisherRepository)
	paymentUsecase := usecase.NewPaymentUseCase(paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)
	logger.SetupLogger()

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, paymentHandler)
	router.Run(":" + port)
	logger.Logger.Printf("Server running on port: %s", port)

}
