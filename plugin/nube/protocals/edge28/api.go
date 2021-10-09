package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET("/edge/ping", func(ctx *gin.Context) {
		p, err := i.rest.PingServer()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.GET("/edge/read/ui/all", func(ctx *gin.Context) {
		p, err := i.rest.GetUIs()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.GET("/edge/read/di/all", func(ctx *gin.Context) {
		p, err := i.rest.GetDIs()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/edge/wizard", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		p, err := i.wizard(body)
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
}
