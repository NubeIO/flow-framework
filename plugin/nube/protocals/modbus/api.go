package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/mbmodel"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	help          = "/help"
	listSerial    = "/list/serial"
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

type Scan struct {
	Start  uint32 `json:"start"`
	Count  uint32 `json:"count"`
	IsCoil bool   `json:"is_coil"`
}

type Body struct {
	Network       *model.Network
	Device        *model.Device
	Point         *model.Point
	Client        `json:"client"`
	Operation     `json:"request_body"`
	Scan          `json:"scan"`
	ReturnArray   bool  `json:"return_array"`
	IsSerial      bool  `json:"is_serial"`
	DeviceAddress uint8 `json:"device_address"`
}

func bodyClient(ctx *gin.Context) (dto Body, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// serialWizard
type wizard struct {
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	SerialPort    string `json:"serial_port"`
	BaudRate      uint   `json:"baud_rate"`
	DeviceAddr    uint   `json:"device_addr"`
	WizardVersion uint   `json:"wizard_version"`
}

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	out.Add(model.ObjTypeAnalogInput)
	out.Add(model.ObjTypeAnalogOutput)
	out.Add(model.ObjTypeAnalogValue)
	out.Add(model.ObjTypeBinaryInput)
	out.Add(model.ObjTypeBinaryOutput)
	out.Add(model.ObjTypeBinaryValue)
	return out
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := i.addNetwork(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, network)
		}
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := i.addDevice(body)
		fmt.Println(err)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, device)
		}
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := i.addPoint(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, point)
		}
	})

	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, serial)
			return
		}
	})
	mux.POST("/modbus/point/operation", func(ctx *gin.Context) {
		body, err := bodyClient(ctx)
		netType := body.Network.TransportType
		mbClient, err := i.setClient(body.Network, body.Device, false)
		if err != nil {
			modbusErrorMsg(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if netType == model.TransType.Serial || netType == model.TransType.LoRa {
			if body.Device.AddressId >= 1 {
				mbClient.RTUClientHandler.SlaveID = byte(body.Device.AddressId)
			}
		} else if netType == model.TransType.IP {
			url, err := uurl.JoinIpPort(body.Device.Host, body.Device.Port)
			if err != nil {
				modbusErrorMsg(fmt.Sprintf("failed to validate device IP %s\n", url))
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			mbClient.TCPClientHandler.Address = url
			mbClient.TCPClientHandler.SlaveID = byte(body.Device.AddressId)
		}
		_, responseValue, err := networkRequest(mbClient, body.Point, false)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		modbusDebugMsg("responseValue", responseValue)
		ctx.JSON(http.StatusOK, responseValue)
		return
	})
	mux.POST("/modbus/wizard/tcp", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		n, err := i.wizardTCP(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, n)
			return
		}
	})
	mux.POST("/modbus/wizard/serial", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		serial, err := i.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, serial)
			return
		}
	})
	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetPointSchema())
	})
}
