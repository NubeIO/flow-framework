package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"time"
)

// The ProducerHistoryDatabase interface for encapsulating database access.
type ProducerHistoryDatabase interface {
	GetProducerHistories(args Args) ([]*model.ProducerHistory, error)
	GetProducerHistoriesByProducerUUID(pUuid string, args Args) ([]*model.ProducerHistory, int64, error)
	GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error)
	GetLatestProducerHistoryByProducerUUID(pUuid string) (*model.ProducerHistory, error)
	GetProducerHistoriesPoints(args Args) ([]*model.History, error)
	GetProducerHistoriesPointsForSync(id string, timeStamp string) ([]*model.History, error)
	CreateProducerHistory(history *model.ProducerHistory) (*model.ProducerHistory, error)
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

func (a *HistoriesAPI) GetProducerHistoriesPointsForSync(ctx *gin.Context) {
	id, timeStamp := buildProducerHistoryPointsSyncArgs(ctx)
	q, err := a.DB.GetProducerHistoriesPointsForSync(id, timeStamp)
	responseHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateProducerHistory(ctx *gin.Context) {
	body, _ := getBodyHistory(ctx)
	body.Timestamp = time.Now().UTC()
	q, err := a.DB.CreateProducerHistory(body)
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
