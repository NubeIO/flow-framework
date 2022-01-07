package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ProducerHistoryDatabase interface for encapsulating database access.
type ProducerHistoryDatabase interface {
	GetProducerHistories(args Args) ([]*model.ProducerHistory, error)
	GetProducerHistoriesByProducerUUID(pUuid string, args Args) ([]*model.ProducerHistory, int64, error)
	GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error)
	GetLatestProducerHistoryByProducerUUID(pUuid string) (*model.ProducerHistory, error)
	GetProducerHistoriesPoints(args Args) ([]*model.History, error)
	CreateBulkProducerHistory(histories []*model.ProducerHistory) (bool, error)
	CreateProducerHistory(history *model.ProducerHistory) (bool, error)
	DeleteAllProducerHistories(args Args) (bool, error)
	DeleteProducerHistoriesByProducerUUID(pUuid string, args Args) (bool, error)
}
type HistoriesAPI struct {
	DB ProducerHistoryDatabase
}

func (a *HistoriesAPI) GetProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistories(args)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, _, err := a.DB.GetProducerHistoriesByProducerUUID(pUuid, args)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetLatestProducerHistoryByProducerName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, err := a.DB.GetLatestProducerHistoryByProducerName(name)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetLatestProducerHistoryByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	q, err := a.DB.GetLatestProducerHistoryByProducerUUID(pUuid)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesPoints(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistoriesPoints(args)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateProducerHistory(ctx *gin.Context) {
	body, _ := getBodyHistory(ctx)
	q, err := a.DB.CreateProducerHistory(body)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateBulkProducerHistory(ctx *gin.Context) {
	body, _ := getBodyBulkHistory(ctx)
	q, err := a.DB.CreateBulkProducerHistory(body)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) DeleteAllProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.DeleteAllProducerHistories(args)
	responseHandler(q, err, ctx)

}

func (a *HistoriesAPI) DeleteProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.DeleteProducerHistoriesByProducerUUID(pUuid, args)
	responseHandler(q, err, ctx)
}
