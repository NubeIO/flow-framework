package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ProducerHistoryDatabase interface for encapsulating database access.
type ProducerHistoryDatabase interface {
	GetProducerHistory(uuid string) (*model.ProducerHistory, error)
	HistoryLatestByProducerUUID(uuid string) (*model.ProducerHistory, error)
	HistoriesAllByProducerUUID(uuid string, order string) ([]*model.ProducerHistory, int64, error)
	GetProducerHistories(args Args) ([]*model.ProducerHistory, error)
	CreateProducerHistory(history *model.ProducerHistory) (*model.ProducerHistory, error)
	DeleteProducerHistory(uuid string) (bool, error)
	DropProducerHistories() (bool, error)
	CreateBulkProducerHistory(history []*model.ProducerHistory) (bool, error)
}
type HistoriesAPI struct {
	DB ProducerHistoryDatabase
}

func (a *HistoriesAPI) GetProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistories(args)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistory(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetProducerHistory(uuid)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) HistoriesAllByProducerUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	order, _ := queryFields(ctx)
	q, _, err := a.DB.HistoriesAllByProducerUUID(uuid, order)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) HistoryLatestByProducerUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.HistoryLatestByProducerUUID(uuid)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateProducerHistory(ctx *gin.Context) {
	body, _ := getBODYHistory(ctx)
	q, err := a.DB.CreateProducerHistory(body)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateBulkProducerHistory(ctx *gin.Context) {
	body, _ := getBODYBulkHistory(ctx)
	q, err := a.DB.CreateBulkProducerHistory(body)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) DeleteProducerHistory(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteProducerHistory(uuid)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) DropProducerHistories(ctx *gin.Context) {
	q, err := a.DB.DropProducerHistories()
	reposeHandler(q, err, ctx)

}
