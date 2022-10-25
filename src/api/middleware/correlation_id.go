package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := uuid.New()
		c.Set("correlation_id", correlationId.String())
		c.Next()
	}
}
