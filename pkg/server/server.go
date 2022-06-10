package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ultagic.com/pkg/middleware"
	"ultagic.com/pkg/utils"
)

var (
	logger *log.Logger
)

func init() {
	logger = utils.Logger("server")
}

func GetServer(addr string, handle func(engine *gin.Engine)) *http.Server {
	engine := gin.Default()

	engine.Use(middleware.ErrorHandler())

	handle(engine)

	logger.Println("server listening on: " + addr)
	fmt.Println("server listening on: " + addr)

	return &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
