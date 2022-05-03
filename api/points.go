package api

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	GetOnePointByArgs(args Args) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	DropPoints() (bool, error)
	GetPointByName(networkName, deviceName, pointName string) (*model.Point, error)
	PointWriteByName(networkName, deviceName, pointName string, body *model.Point, fromPlugin bool) (*model.Point, error)

	CreatePointPlugin(body *model.Point) (*model.Point, error)
	UpdatePointPlugin(uuid string, body *model.Point) (*model.Point, error)
	DeletePointPlugin(uuid string) (bool, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoints(args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointsBulk(ctx *gin.Context) {
	body, _ := getBODYBulkPoints(ctx)
	q, err := a.DB.GetPointsBulk(body)
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
	q, err := a.DB.UpdatePointPlugin(uuid, body)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) PointWrite(ctx *gin.Context) {
	fmt.Println(fmt.Sprintf("PointWrite ctx %+v", ctx))
	fmt.Println(fmt.Sprintf("PointWrite ctx.Request %+v", ctx.Request))
	fmt.Println(fmt.Sprintf("PointWrite ctx.Request.Body %+v", ctx.Request.Body))
	bodyPrint, err := ioutil.ReadAll(ctx.Request.Body) //THIS GETS THE RAW BYTE STREAM OF BODY
	fmt.Println(fmt.Sprintf("PointWrite ctx.Request.Body %s", bodyPrint))
	fmt.Println(fmt.Sprintf("PointWrite ctx.Request.Body err %+v", err))

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
	q, err := a.DB.DeletePointPlugin(uuid)
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
