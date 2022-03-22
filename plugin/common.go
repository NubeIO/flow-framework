package plugin

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

const (
	NetworksURL = "/networks"
	DevicesURL  = "/devices"
	PointsURL   = "/points"
)

func GetBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}
