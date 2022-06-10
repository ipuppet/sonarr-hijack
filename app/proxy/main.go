package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ultagic.com/app/proxy/handler"
	"ultagic.com/pkg/server"
)

func GetServer(addr string) *http.Server {
	return server.GetServer(addr, func(engine *gin.Engine) {
		handler.LoadRouters(engine)
	})
}
