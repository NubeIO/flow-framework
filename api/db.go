package api

import (
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (string, error) //delete all networks, gateways and children
	SyncTopics()                  //sync all the topics into the event bus
	WizardLocalPointMapping(body *WizardLocalMapping) (bool, error)
	WizardRemotePointMapping() (bool, error)
	Wizard2ndFlowNetwork(body *AddNewFlowNetwork) (bool, error)
}
type DatabaseAPI struct {
	DB DBDatabase
}

func (a *DatabaseAPI) DropAllFlow(ctx *gin.Context) {
	q, err := a.DB.DropAllFlow()
	reposeHandler(q, err, ctx)
}

func (a *DatabaseAPI) SyncTopics() {
	a.DB.SyncTopics()
}

type WizardLocalMapping struct {
	PluginName         string `json:"plugin_name"`
	ProducerStreamUUID string `json:"producer_stream_uuid"`
	ConsumerStreamUUID string `json:"consumer_stream_uuid"`
	DeviceUUID         string `json:"device_uuid"`
	ThingUUID          string `json:"thing_uuid"`
}

func getBODYLocal(ctx *gin.Context) (dto *WizardLocalMapping, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func (a *DatabaseAPI) WizardLocalPointMapping(ctx *gin.Context) {
	body, _ := getBODYLocal(ctx)
	mapping, err := a.DB.WizardLocalPointMapping(body)
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardRemotePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardRemotePointMapping()
	reposeHandler(mapping, err, ctx)
}

type AddNewFlowNetwork struct {
	StreamUUID         string `json:"stream_uuid"`
	ProducerUUID       string `json:"producer_uuid"`
	ProducerThingUUID  string `json:"producer_thing_uuid"` // this is the remote point UUID
	ProducerThingClass string `json:"producer_thing_class"`
	ProducerThingType  string `json:"producer_thing_type"`
	FlowToken          string `json:"flow_token"`
}

func getBODYWizard(ctx *gin.Context) (dto *AddNewFlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func (a *DatabaseAPI) Wizard2ndFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYWizard(ctx)
	q, err := a.DB.Wizard2ndFlowNetwork(body)
	reposeHandler(q, err, ctx)
}
