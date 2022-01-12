package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) SyncStream(body *model.SyncStream) (*model.StreamClone, error) {
	var flowNetworkClone *model.FlowNetworkClone
	if err := d.DB.Where("global_uuid = ?", body.GlobalUUID).Find(&flowNetworkClone).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("we don't have flow_network_clone with global_uuid=%s", body.GlobalUUID))
	}
	mStream, err := json.Marshal(body.Stream)
	if err != nil {
		return nil, err
	}
	streamClone := model.StreamClone{}
	if err = json.Unmarshal(mStream, &streamClone); err != nil {
		return nil, err
	}
	streamClone.FlowNetworkCloneUUID = flowNetworkClone.UUID
	streamClone.SourceUUID = body.Stream.UUID
	var streamClonesModel []*model.StreamClone
	if err = d.DB.Where("source_uuid = ? ", streamClone.SourceUUID).Find(&streamClonesModel).Error; err != nil {
		return nil, err
	}
	if len(streamClonesModel) == 0 {
		streamClone.UUID = utils.MakeTopicUUID(model.CommonNaming.StreamClone)
		if err = d.DB.Create(streamClone).Error; err != nil {
			return nil, err
		}
	} else {
		streamClone.UUID = streamClonesModel[0].UUID
		if err = d.DB.Model(&streamClonesModel[0]).Updates(streamClone).Error; err != nil {
			return nil, err
		}
	}
	return &streamClone, nil
}
