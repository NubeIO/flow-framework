package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Scan struct {
	Start  uint32 `json:"start"`
	Count  uint32 `json:"count"`
	IsCoil bool   `json:"is_coil"`
}

type Body struct {
	Client    `json:"client"`
	Operation `json:"request_body"`
	Scan      `json:"scan"`
}

func bodyClient(ctx *gin.Context) (dto Body, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

type Enable struct {
	Type        string
	ReadOnly    string
	WriteValues interface{}
}

type T struct {
	Enabled      bool        `json:"enabled"`
	ProtocolType interface{} `json:"protocol_type"`
	Capabilities []string    `json:"capabilities"`
	Networks     struct {
		Uuid   string      `json:"uuid"`
		Name   string      `json:"name"`
		Enable interface{} `json:"enable"`
	} `json:"networks"`
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET("/modbus/help/network", func(ctx *gin.Context) {

		h := new(T)
		h.ProtocolType = model.TransType

		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, h)
		}
	})
	mux.GET("/modbus/list/serial", func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/modbus/point/tcp/operation", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := setClient(body.Client)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			request, err := parseRequest(body.Operation)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "request was invalid, try readCoil or writeCoil")
				return
			}
			r, err := DoOperations(cli, request)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, r)
			}

		}
	})
	mux.POST("/modbus/wizard/tcp", func(ctx *gin.Context) {
		serial, err := i.wizardTCP()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/modbus/wizard/serial", func(ctx *gin.Context) {
		serial, err := i.wizardSerial()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/scan/bool", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := setClient(body.Client)
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

}
