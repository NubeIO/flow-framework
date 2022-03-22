package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
)

func (a *FlowClient) GetProducerHistoriesPoints(lastSyncId int) (*[]model.History, error) {
	req := a.client.R().
		SetResult(&[]model.History{}).SetQueryParam("id_gt", fmt.Sprintf("%v", lastSyncId))
	resp, err := req.Get("/api/histories/producers/points")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetPointsProducerHistories: %s", err)
		} else {
			return nil, fmt.Errorf("GetPointsProducerHistories: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetPointsProducerHistories: %s", resp)
	}
	return resp.Result().(*[]model.History), nil
}
