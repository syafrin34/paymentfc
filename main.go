package main

import (
	"context"
	"encoding/json"
	"fmt"
	"paymentfc/cmd/payment/handler"
	"paymentfc/cmd/payment/repository"
	"paymentfc/cmd/payment/resource"
	"paymentfc/cmd/payment/service"
	"paymentfc/cmd/payment/usecase"
	"paymentfc/config"
	"paymentfc/grpc"
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
	kafkaWriter := kafka.NewWriter(cfg.Kafka.Broker, cfg.Kafka.Topics[constant.KafkaTopicPaymentSuccess])
	fmt.Printf("Kafka Broker: %s\n", cfg.Kafka.Broker)
	fmt.Printf("Kafka Topics: %#v\n", cfg.Kafka.Topics)
	logger.SetupLogger()

	grpcUserClient := grpc.NewUserClient()
	userInfo, err := grpcUserClient.GetUserInfoByUserID(context.Background(), 10)
	if err != nil {
		logger.Logger.Fatalf("error get user info using grpc", err.Error())
	}
	tempDebug39, _ := json.Marshal(userInfo)
	fmt.Printf("\n==== DEBUG main.go - line 39 ======\n\n%s\n\n========\n\n\n", string(tempDebug39))
	databaseRepository := repository.NewPaymentDatabase(db)
	publisherRepository := repository.NewKafkaPublisher(kafkaWriter)

	paymentService := service.NewPaymentService(databaseRepository, publisherRepository)
	paymentUsecase := usecase.NewPaymentUseCase(paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase, cfg.Xendit.XenditWebhook)
	xenditRepository := repository.NewXenditClient(cfg.Xendit.XenditAPIKey)
	xenditService := service.NewXenditService(*grpcUserClient, databaseRepository, xenditRepository)
	xenditUsecase := usecase.NewXenditUseCase(xenditService)

	//scheduler
	scheduler := service.SchedulerService{
		Database:       databaseRepository,
		Xendit:         xenditRepository,
		Publisher:      publisherRepository,
		PaymentService: paymentService,
	}
	scheduler.StartCheckPendingInvoices()
	scheduler.StartProcessPendingPaymentRequests()
	scheduler.StartProcessFailedPaymentRequests()
	scheduler.StartSweepingExpiredPendingPayments()
	// REMOVE ME - DEBUG PURPOSE ONLY
	tempDebug69, _ := json.Marshal(cfg.Kafka)
	fmt.Printf("\n==== DEBUG main.go - Line: 69 ===== \n\n%s\n\n=====================\n\n\n", string(tempDebug69))
	// END OF REMOVE ME

	// kafka consumer
	// potential less efficient  --> traffict gede
	kafka.StartOrderConsumer(cfg.Kafka.Broker, cfg.Kafka.Topics[constant.KafkaTopicOrderCreated], func(event models.OrderCreatedEvent) {

		if cfg.Toggle.DisableCreateInvoiceDirectly {
			err := paymentUsecase.ProcessPaymentRequests(context.Background(), event)
			if err != nil {
				logger.Logger.Println("Enabled Handle Order Created Event: ", err.Error())
			}

		} else {
			err := xenditUsecase.CreateInvoice(context.Background(), event)
			if err != nil {
				logger.Logger.Println("Failed Handling order created event: ", err.Error())
			}

		}

	})

	// current condition
	/*
		- user checkout order
		- order execute checkout --> publish event order.created
		- payment service akan memproses create invoice
	*/

	// new condition
	/*
		- user checkout order
		- order execute event order.created
		- payment service akan simpan event yang dari order.created
		- payment akan menyediakan backgroud process utk create invoice per batch
	*/

	// cons: - data tidak execute secara realtime
	// pertimbangan : transactional especially payment process --> harus lebih consistency dan stability
	// pro:
	/*
		sample scenario:
			-xendit team informaed there will be maintenance for 5 minutes (12:00 - 12:05)
			- kita bisa hold execute payment requests sampai xendit stable
			- data dari order service  (order.created) tidak menumpuk
	*/

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, paymentHandler)
	router.Run(":" + port)
	logger.Logger.Printf("Server running on port: %s", port)

}
