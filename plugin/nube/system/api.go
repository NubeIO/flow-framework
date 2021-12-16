package main

import (
	"net/http"

	"github.com/NubeIO/flow-framework/model"
	system_model "github.com/NubeIO/flow-framework/plugin/nube/system/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

//markdown guide
const helpText = `
# LoRa Help Guide

help

### line 2
*new tab*
<a href="https://stackoverflow.com" target="_blank">New Tab</a>
You will never use anything else than this [website].

- this is some normal texy
- this is ***some*** normal texy
-- aaaaaa
1. First item
2. Second item
3. Third item

| Syntax | Description |
| ----------- | ----------- |
| Header | Title |
| Paragraph | Text |

- [x] Write the press release
- [ ] Update the website
- [ ] Contact the media

this is *some* normal texy`

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

const (
	help     = "/help"
	helpHTML = "/help/guide"
)

var Supports = struct {
	Network     bool `json:"network"`
	NetworkCRUD bool `json:"networkCRUD"`
}{
	Network:     true,
	NetworkCRUD: true,
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, system_model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, system_model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, system_model.GetPointSchema())
	})

	mux.GET(help, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Supports)
	})
	mux.GET(helpHTML, func(ctx *gin.Context) {
		md := []byte(helpText)
		output := markdown.ToHTML(md, nil, nil)
		ctx.Writer.Write(output)
	})
	mux.GET("/system/schedule/store/:name", func(ctx *gin.Context) {
		obj, ok := i.store.Get(resolveName(ctx))
		if ok != true {
			ctx.JSON(http.StatusBadRequest, "no schedule exists")
		} else {
			ctx.JSON(http.StatusOK, obj)
		}
	})
}
