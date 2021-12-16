package main

import (
	"net/http"

	baseModel "github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	SerialPort string `json:"serial_port"`
	BaudRate   uint   `json:"baud_rate"`
	DeviceAddr uint   `json:"device_addr"`
}

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

type Urls struct {
	Help       []string `json:"help"`
	ListSerial string   `json:"list_serial"`
}

type T2 struct {
	Capabilities struct {
		Help []string `json:"help,omitempty"`
	} `json:"capabilities,omitempty"`
	ObjectType struct {
		Options  interface{} `json:"options"`
		Type     string      `json:"type"`
		Required bool        `json:"required"`
	} `json:"object_type"`
}

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	out.Add(baseModel.ObjTypeAnalogInput)
	out.Add(baseModel.ObjTypeAnalogOutput)
	out.Add(baseModel.ObjTypeAnalogValue)
	out.Add(baseModel.ObjTypeBinaryInput)
	out.Add(baseModel.ObjTypeBinaryOutput)
	out.Add(baseModel.ObjTypeBinaryValue)
	return out
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(help, func(ctx *gin.Context) {

		var h T2
		//h.Enabled = model.CommonNaming
		//a1 := []string{fmt.Sprintf("http://0.0.0.0:1660/api/plugins/api/%s%s" ,name, help), "GET", "POST", "PATCH"}
		//a1 := []string{fmt.Sprintf("http://0.0.0.0:1660/api/plugins/api/%s%s" ,name, help), "GET", "POST", "PATCH"}
		//h.Capabilities.Help = a1
		h.ObjectType.Options = supportedObjects()
		h.ObjectType.Type = "array"
		h.ObjectType.Required = true

		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, h)
		}
	})
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/modbus/point/operation", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := i.setClient(body.Client, "", false, body.IsSerial)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			err := cli.SetUnitId(body.DeviceAddress)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "request was invalid, failed to SetUnitId")
				return
			}
			request, err := parseRequest(body.Operation)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "request was invalid, try readCoil or writeCoil")
				return
			}
			rArray, responseValue, err := networkRequest(cli, request)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				if body.ReturnArray {
					log.Info(responseValue)
					ctx.JSON(http.StatusOK, rArray)
				} else {
					ctx.JSON(http.StatusOK, responseValue)
				}
			}
		}
	})
	mux.POST("/modbus/wizard/tcp", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		n, err := i.wizardTCP(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, n)
		}
	})
	mux.POST("/modbus/wizard/serial", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		serial, err := i.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/scan/bool", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := i.setClient(body.Client, "", false, false)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			found, _ := performBoolScan(cli, body.Scan.IsCoil, body.Scan.Start, body.Scan.Count)
			ctx.JSON(http.StatusOK, found)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.GetPointSchema())
	})
}
