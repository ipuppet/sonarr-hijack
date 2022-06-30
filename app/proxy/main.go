package proxy

import (
	"net/http"

	"sonarr-hijack/app/proxy/handler"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/server"
)

func GetServer(addr string) *http.Server {
	return server.GetServer(addr, func(engine *gin.Engine) {
		handler.LoadRouters(engine)
	})
}
