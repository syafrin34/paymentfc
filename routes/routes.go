package routes

import (
	"paymentfc/cmd/payment/handler"
	"paymentfc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler) {
	router.Use(middleware.RequestLogger())

	router.POST("/v1/paymnet/webhook", paymentHandler.HandleXenditWebhook)
	router.GET("v1/invoice/:order_id/pdf", paymentHandler.HandlerDownloadPDFInvoice)

	// internal useg
	router.GET("v1/failed_payments", paymentHandler.HandleFailedPayments)

}
