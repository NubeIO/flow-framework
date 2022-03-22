package main

import (
	model "github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/edge_model"
	edgerest "github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/restclient"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetPointSchema())
	})

	mux.GET("/edge/ping", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.PingServer()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, p)
			return
		}
	})
	mux.GET("/edge/read/ui/all", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.GetUIs()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, p)
			return
		}
	})
	mux.GET("/edge/read/di/all", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.GetDIs()
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
	mux.POST("/edge/write/uo", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		ioNum := body.IONum
		val := body.Value
		p, err := rest.WriteUO(ioNum, val)
		if err != nil {
			log.Info(err, "ERROR ON write uo")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/edge/write/do", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		ioNum := body.IONum
		val := body.Value
		p, err := rest.WriteDO(ioNum, val)
		if err != nil {
			log.Info(err, "ERROR ON write do")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
}
