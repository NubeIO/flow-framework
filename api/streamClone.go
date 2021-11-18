package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type StreamCloneDatabase interface {
	GetStreamClones(args Args) ([]*model.StreamClone, error)
	GetStreamClone(uuid string, args Args) (*model.StreamClone, error)
	DeleteStreamClone(uuid string) error
}

type StreamCloneAPI struct {
	DB StreamCloneDatabase
}

func (j *StreamCloneAPI) GetStreamClones(ctx *gin.Context) {
	args := buildStreamCloneArgs(ctx)
	q, err := j.DB.GetStreamClones(args)
	responseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) GetStreamClone(ctx *gin.Context) {
	args := buildStreamCloneArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStreamClone(uuid, args)
	responseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) DeleteStreamClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	err := j.DB.DeleteStreamClone(uuid)
	responseHandler(nil, err, ctx)
}
