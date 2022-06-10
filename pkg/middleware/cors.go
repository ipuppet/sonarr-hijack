package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ultagic.com/pkg/config"
)

var (
	corsConfig *config.Config
)

func Cors(app string) gin.HandlerFunc {
	if corsConfig == nil {
		corsConfig = &config.Config{
			Filename: "cors.json",
		}
		corsConfig.Init()
		corsConfig.AddNotifyer(config.LoggerNotifyer())
	}

	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		cors, err := corsConfig.Get(app)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		for _, allowHost := range cors.([]interface{}) {
			if origin == allowHost.(string) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
				c.Header("Access-Control-Allow-Credentials", "true")

				break
			}
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
