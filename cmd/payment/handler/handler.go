package handler

import (
	"net/http"
	"paymentfc/cmd/payment/usecase"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler interface {
	HandleXenditWebhook(c *gin.Context)
	HandleCreateInvoice(c *gin.Context)
	HandlerDownloadPDFInvoice(c *gin.Context)
}

type paymentHandler struct {
	usecase            usecase.PaymentUseCase
	xenditWebhookToken string
}

func NewPaymentHandler(usecase usecase.PaymentUseCase, webhookToken string) PaymentHandler {
	return &paymentHandler{
		usecase:            usecase,
		xenditWebhookToken: webhookToken,
	}
}

func (h *paymentHandler) HandleXenditWebhook(c *gin.Context) {
	var payload models.XenditWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"payload": payload,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "invalid payload",
			"error message": err.Error(),
		})
		return

	}

	// validate web hook token
	headerWebhookToken := c.GetHeader("x-callback-token")
	if h.xenditWebhookToken != headerWebhookToken {
		logger.Logger.WithFields(logrus.Fields{
			"callbackToken": headerWebhookToken,
		}).Errorf("Invalid Webhook Token: %s", headerWebhookToken)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid webhook token!"})
		return
	}

	err := h.usecase.ProcessPaymentWebhook(c.Request.Context(), payload)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"payload": payload,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	return
}

func (h *paymentHandler) HandleCreateInvoice(c *gin.Context) {
	var payload models.OrderCreatedEvent

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error message": "bad request",
		})
		return
	}

	return
}
func (h *paymentHandler) HandlerDownloadPDFInvoice(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)

	filePath, err := h.usecase.DownloadPDFInvoice(c.Request.Context(), orderID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).WithError(err).Errorf("h.usecase.DownloadPdfInvoice() got error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err.Error,
		})
		return

	}
	c.FileAttachment(filePath, filePath)
}
