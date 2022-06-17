package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwrest"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyDevice(ctx *gin.Context) (dto lwmodel.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

const chirpName = "admin"
const chirpPass = "Helensburgh2508"

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	cli := lwrest.NewChirp(chirpName, chirpPass, ip, port)

	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.addPoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.updateNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.updateDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.updatePoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBodyPointWriter(ctx)
		uuid := plugin.ResolveID(ctx)
		point, err := inst.writePoint(uuid, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.deleteNetwork(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		ok, err := inst.deleteDevice(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.deletePoint(body)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.GET("/lorawan/organizations", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/gateways", func(ctx *gin.Context) {
		p, err := cli.GetGateways()
		if err != nil {
			log.Info(err, "ERROR ON gateways")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/applications", func(ctx *gin.Context) {
		p, err := cli.GetApplications()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/devices", func(ctx *gin.Context) {
		p, err := cli.GetDevices()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.GetDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON GetDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.PUT("/lorawan/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyDevice(ctx)
		if err != nil {
			return
		}
		_, err = cli.EditDevice(eui, device)
		if err != nil {
			log.Info(err, "ERROR ON GetDevice")
			ctx.JSON(http.StatusBadRequest, "fail")
		} else {
			ctx.JSON(http.StatusOK, "ok")
		}
	})

	mux.DELETE("/lorawan/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.DeleteDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON DeleteDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}

	})

	mux.GET("/lorawan/device-profiles", func(ctx *gin.Context) {
		p, err := cli.GetDeviceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON device-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/service-profiles", func(ctx *gin.Context) {
		p, err := cli.GetServiceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON service-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/gateway-profiles", func(ctx *gin.Context) {
		p, err := cli.GetGatewayProfiles()
		if err != nil {
			log.Info(err, "ERROR ON gateway-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lwmodel.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lwmodel.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lwmodel.GetPointSchema())
	})

}
