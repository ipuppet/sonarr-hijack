package middleware

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// check error
		for _, err := range c.Errors {
			c.JSON(c.Writer.Status(), gin.H{
				"status":  false,
				"message": err.Error(),
			})
			c.Abort()

			return
		}
	}
}
