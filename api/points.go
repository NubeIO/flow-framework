package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error)
	CreatePointPlugin(body *model.Point) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	GetOnePointByArgs(args Args) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	DropPoints() (bool, error)
	GetPointByName(networkName, deviceName, pointName string) (*model.Point, error)
	PointWriteByName(networkName, deviceName, pointName string, body *model.Point, fromPlugin bool) (*model.Point, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoints(args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoint(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) UpdatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdatePoint(uuid, body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) PointWrite(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.PointWrite(uuid, body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetOnePointByArgs(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetOnePointByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	q, err := a.DB.CreatePointPlugin(body)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) DeletePoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeletePoint(uuid)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) DropPoints(ctx *gin.Context) {
	q, err := a.DB.DropPoints()
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointByName(ctx *gin.Context) {
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.GetPointByName(networkName, deviceName, pointName)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) PointWriteByName(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.PointWriteByName(networkName, deviceName, pointName, body, false)
	responseHandler(q, err, ctx)
}

func networkDevicePointNames(ctx *gin.Context) (networkName, deviceName, pointName string) {
	var aType = ArgsType
	var aDefault = ArgsDefault
	networkName = ctx.DefaultQuery(aType.NetworkName, aDefault.NetworkName)
	deviceName = ctx.DefaultQuery(aType.DeviceName, aDefault.DeviceName)
	pointName = ctx.DefaultQuery(aType.PointName, aDefault.PointName)
	return networkName, deviceName, pointName
}
