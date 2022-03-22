package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	log "github.com/sirupsen/logrus"
	"time"
)

func (i *Instance) syncHistory() (bool, error) {
	log.Info("History sync has is been called...")
	fnClones, err := i.db.GetFlowNetworkClones()
	if err != nil {
		return false, err
	}
	var histories []*model.History
	for _, fnc := range fnClones {
		hisLog, err := i.db.GetHistoryLogByFlowNetworkCloneUUID(fnc.UUID)
		if err != nil {
			return false, err
		}
		cli := client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
		pHistories, err := cli.GetProducerHistoriesPoints(hisLog.LastSyncID)
		if err != nil {
			return false, err
		}
		for k, h := range *pHistories {
			h := h // more: https://medium.com/swlh/use-pointer-of-for-range-loop-variable-in-go-3d3481f7ffc9
			histories = append(histories, &h)
			if k == len(*pHistories)-1 { // Update History Log
				hisLog.FlowNetworkCloneUUID = fnc.UUID
				hisLog.LastSyncID = h.ID
				hisLog.Timestamp = time.Now()
				_, err = i.db.UpdateHistoryLog(hisLog)
				if err != nil {
					return false, err
				}
			}
		}
	}
	if len(histories) > 0 {
		_, err = i.db.CreateBulkHistory(histories)
		if err != nil {
			return false, err
		}
	}
	log.Info(fmt.Sprintf("Stored %v rows on %v", len(histories), path))
	return true, nil
}
