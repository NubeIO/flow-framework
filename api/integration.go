package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The IntegrationDatabase interface for encapsulating database access.
type IntegrationDatabase interface {
	GetIntegration(uuid string) (*model.Integration, error)
	GetIntegrations() ([]*model.Integration, error)
	CreateIntegration(body *model.Integration) (*model.Integration, error)
	UpdateIntegration(uuid string, body *model.Integration) (*model.Integration, error)
	DeleteIntegration(uuid string) (bool, error)
}

type IntegrationAPI struct {
	DB IntegrationDatabase
}

func (j *IntegrationAPI) GetIntegration(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetIntegration(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *IntegrationAPI) GetIntegrations(ctx *gin.Context) {
	q, err := j.DB.GetIntegrations()
	ResponseHandler(q, err, ctx)

}

func (j *IntegrationAPI) CreateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateIntegration(body)
	ResponseHandler(q, err, ctx)
}

func (j *IntegrationAPI) UpdateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateIntegration(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (j *IntegrationAPI) DeleteIntegration(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteIntegration(uuid)
	ResponseHandler(q, err, ctx)
}
