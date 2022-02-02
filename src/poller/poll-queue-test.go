package poller

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	log "github.com/sirupsen/logrus"
)

func PollQueueTest() {
	h := &dbhandler.Handler{}
	dbhandler.Init(h)
	var arg api.Args
	net, err := h.DB.GetNetworkByPluginName("modbus", arg)
	if err != nil {
		fmt.Printf("PollQueueTest: no modbus networks found.\n")
	}
	log.Info(net)
	for {
		pp := net.PollManager.GetNextPollingPoint()
		if pp != nil {
			fmt.Printf("PollQueueTest Polling Point: priority = %d, uuid = %s\n", pp.PollPriority, pp.FFPointUUID)
		}
	}

}
