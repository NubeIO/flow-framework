package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	//restart plugin
	mux.POST("/broker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})

}
