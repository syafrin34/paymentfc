package middleware

import (
	"context"
	"paymentfc/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		ctx := context.WithValue(timeoutCtx, "request_id", requestId)
		c.Request = c.Request.WithContext(ctx)

		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		requestLog := logger.Fields{
			"request_field": requestId,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status":        c.Writer.Status(),
			"latency":       latency,
		}
		if c.Writer.Status() == 200 || c.Writer.Status() == 201 {
			logger.Logger.WithFields(requestLog).Info("Request success")

		} else {

			logger.Logger.WithFields(requestLog).Info("Request Error")

		}
	}

}
