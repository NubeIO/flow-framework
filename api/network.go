package api

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type NetworkDatabase interface {
	GetNetworkByPlugin(uuid string, withChildren bool, withPoints bool, transport string) (*model.Network, error)
	GetNetworks(args Args) ([]*model.Network, error)
	GetNetwork(uuid string, args Args) (*model.Network, error)
	CreateNetwork(network *model.Network) (*model.Network, error)
	UpdateNetwork(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetwork(uuid string) (bool, error)
	DropNetworks() (bool, error)
}
type NetworksAPI struct {
	DB  NetworkDatabase
	Bus eventbus.BusService
}

func (a *NetworksAPI) GetNetworks(ctx *gin.Context) {
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworks(args)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworkByPlugin(ctx *gin.Context) {
	uuid := resolveID(ctx)
	withChildren, withPoints, _ := withChildrenArgs(ctx)
	q, err := a.DB.GetNetworkByPlugin(uuid, withChildren, withPoints, "") //TODO fix this need to add in like "serial"
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetwork(uuid, args)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) CreateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	q, err := a.DB.CreateNetwork(body)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) UpdateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateNetwork(uuid, body)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteNetwork(uuid)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) DropNetworks(ctx *gin.Context) {
	q, err := a.DB.DropNetworks()
	reposeHandler(q, err, ctx)
}
