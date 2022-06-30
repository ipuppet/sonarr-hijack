package main

import (
	"log"

	"sonarr-hijack/app/proxy"

	"github.com/gin-gonic/gin"
	_ "github.com/ipuppet/gtools/flags"
	"golang.org/x/sync/errgroup"
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
