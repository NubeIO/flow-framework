package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	mWriter, err := json.Marshal(body.Writer)
	if err != nil {
		return nil, err
	}
	writerClone := model.WriterClone{}
	if err = json.Unmarshal(mWriter, &writerClone); err != nil {
		return nil, err
	}
	_, err = d.GetProducer(body.ProducerUUID, api.Args{})
	if err != nil {
		return nil, errors.New("producer does not exist")
	}
	writerClone.ProducerUUID = body.ProducerUUID
	writerClone.SourceUUID = body.Writer.UUID
	var writerCloneModel []*model.WriterClone
	if err = d.DB.Where("source_uuid = ? ", writerClone.SourceUUID).Find(&writerCloneModel).Error; err != nil {
		return nil, err
	}
	if len(writerCloneModel) == 0 {
		writerClone.UUID = utils.MakeTopicUUID(model.CommonNaming.WriterClone)
		if err = d.DB.Create(writerClone).Error; err != nil {
			return nil, err
		}
	} else {
		writerClone.UUID = writerCloneModel[0].UUID
		if err = d.DB.Model(&writerCloneModel[0]).Updates(writerClone).Error; err != nil {
			return nil, err
		}
	}
	return &writerClone, nil
}
