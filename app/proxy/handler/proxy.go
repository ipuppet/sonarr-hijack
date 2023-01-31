package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sonarr-hijack/app/proxy/logic"
)

func LoadRouters(e *gin.Engine) {
	e.Any("/jackett/api/:apiVer/indexers/:indexer/results/:feedType/api", func(c *gin.Context) {
		// 以下参数暂时没用
		type UriParam struct {
			ApiVer   string `uri:"apiVer" binding:"required"`
			Indexer  string `uri:"indexer" binding:"required"`
			FeedType string `uri:"feedType" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp, err := logic.Proxy(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer resp.Body.Close()

		// header
		for key, values := range resp.Header {
			if len(values) == 1 {
				c.Writer.Header().Set(key, values[0])
			} else {
				c.Writer.Header().Set(key, values[0])
				for _, value := range values[1:] {
					c.Writer.Header().Add(key, value)
				}
			}
		}

		c.DataFromReader(
			resp.StatusCode,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
			resp.Body,
			nil,
		)
	})
}
