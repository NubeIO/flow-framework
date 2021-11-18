package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncWriterDatabase interface {
	SyncWriter(fn *model.SyncWriter) (*model.WriterClone, error)
}

type SyncWriterAPI struct {
	DB SyncWriterDatabase
}

func (a *SyncWriterAPI) SyncWriter(ctx *gin.Context) {
	body, _ := getBodySyncWriter(ctx)
	q, err := a.DB.SyncWriter(body)
	responseHandler(q, err, ctx)
}
