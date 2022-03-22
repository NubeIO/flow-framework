package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type DeviceDatabase interface {
	GetDevices(args Args) ([]*model.Device, error)
	GetDevice(uuid string, args Args) (*model.Device, error)
	GetOneDeviceByArgs(args Args) (*model.Device, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	CreateDevicePlugin(body *model.Device) (*model.Device, error)
	UpdateDevice(uuid string, body *model.Device, fromPlugin bool) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
	DropDevices() (bool, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevices(args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevice(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetOneDeviceByArgs(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetOneDeviceByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateDevice(uuid, body, false)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	q, err := a.DB.CreateDevicePlugin(body)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteDevice(uuid)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) DropDevices(ctx *gin.Context) {
	q, err := a.DB.DropDevices()
	responseHandler(q, err, ctx)
}
