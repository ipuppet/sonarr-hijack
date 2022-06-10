package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"ultagic.com/app/proxy"
	_ "ultagic.com/pkg/flags"
)

var (
	g errgroup.Group
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	proxyServer := proxy.GetServer(":9117")

	g.Go(func() error {
		return proxyServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
