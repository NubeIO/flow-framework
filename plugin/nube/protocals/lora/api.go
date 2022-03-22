package main

import (
	"github.com/NubeIO/flow-framework/plugin"
	"net/http"

	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/lora_model"
	"github.com/gin-gonic/gin"
)

const (
	networks = "/networks"
	devices  = "/devices"
	points   = "/points"
)

const (
	help          = "/help"
	restartSerial = "/serial/restart"
	listSerial    = "/serial/list"
	wizardSerial  = "/wizard/serial"
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

// wizard
type wizard struct {
	SerialPort string `json:"serial_port"`
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
}

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(networks, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, network)
		}
	})
	mux.POST(devices, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, device)
		}
	})
	mux.POST(points, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		err := inst.addPoint(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, body)
		}
	})
	mux.POST(restartSerial, func(ctx *gin.Context) {
		err := inst.Disable()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			err := inst.Enable()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, "ok")
			}
		}
	})
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := inst.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST(wizardSerial, func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		serial, err := inst.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lora_model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lora_model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lora_model.GetPointSchema())
	})
}
