package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	SyncTopics() // sync all the topics into the event bus
	WizardP2PMapping(body *model.P2PBody) (bool, error)
	WizardMasterSlavePointMapping() (bool, error)
	WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
	WizardP2PMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
}
type DatabaseAPI struct {
	DB DBDatabase
}

func (a *DatabaseAPI) SyncTopics() {
	a.DB.SyncTopics()
}

func (a *DatabaseAPI) WizardP2PMapping(ctx *gin.Context) {
	body, _ := getP2PBody(ctx)
	mapping, err := a.DB.WizardP2PMapping(body)
	ResponseHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardMasterSlavePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardMasterSlavePointMapping()
	ResponseHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID)
	ResponseHandler(sch, err, ctx)
}

func (a *DatabaseAPI) WizardP2PMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardP2PMappingOnConsumerSideByProducerSide(globalUUID)
	ResponseHandler(sch, err, ctx)
}
