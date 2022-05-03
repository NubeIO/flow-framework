package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
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
			continue
		}
		cli := client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
		pHistories, err := cli.GetProducerHistoriesPointsForSync(hisLog.LastSyncID, hisLog.Timestamp)
		if err != nil {
			continue
		}
		for k, h := range *pHistories {
			h := h // more: https://medium.com/swlh/use-pointer-of-for-range-loop-variable-in-go-3d3481f7ffc9
			histories = append(histories, &h)
			if k == len(*pHistories)-1 { // Update History Log
				hisLog.FlowNetworkCloneUUID = fnc.UUID
				hisLog.LastSyncID = h.ID
				hisLog.Timestamp = h.Timestamp
				_, err = i.db.UpdateHistoryLog(hisLog)
				if err != nil {
					continue
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
