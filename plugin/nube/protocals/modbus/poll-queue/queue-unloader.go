package pollqueue

import (
	"github.com/NubeIO/flow-framework/api"
	"time"
	//log "github.com/sirupsen/logrus"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

// Polling Manager Summary:
//  - Diagram Summary: https://docs.google.com/drawings/d/1priwsaQ6EryRBx1kLQd91REJvHzFyxz7cOHYYXyBNFE/edit?usp=sharing
//  - The QueueUnloader is the only way to get the next PollPoint from a Queue
//  - When a QueueUnloader is stopped, the Worker go routine is closed and the reference to the QueueUnloader is set to nil.

//Questions:
// -

//There should be a function in Modbus(or other protocals) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

type QueueUnloader struct {
	NextPollPoint *PollingPoint
	//NextUnloadTimer 		*time.Timer
	NextUnloadTimer *time.Ticker
	CancelChannel   chan bool
}

func (pm *NetworkPollManager) StartQueueUnloader() {
	pollQueueDebugMsg("StartQueueUnloader() 1")
	pm.StopQueueUnloader()
	pollQueueDebugMsg("StartQueueUnloader() 2")
	ql := &QueueUnloader{nil, nil, nil}
	pm.PluginQueueUnloader = ql
	if pm.PluginQueueUnloader.NextPollPoint == nil {
		pollQueueDebugMsg("StartQueueUnloader() pm.PluginQueueUnloader.NextPollPoint == nil")
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
		}
	}
	var netArg api.Args
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, netArg)
	if err != nil {
		pollQueueDebugMsg("NetworkPollManager.StartQueueUnloader(): couldn't find network %s \n", pm.FFNetworkUUID)
		return
	}
	pollQueueDebugMsg("NetworkPollManager.StartQueueUnloader(): net.MaxPollRate %d \n", net.MaxPollRate)
	refreshRate := net.MaxPollRate
	if pm.MaxPollRate <= 0*time.Second {
		refreshRate = 1 * time.Second
	}
	pm.MaxPollRate = refreshRate
	ticker := time.NewTicker(refreshRate)
	pm.PluginQueueUnloader.NextUnloadTimer = ticker
	done := make(chan bool)
	pm.PluginQueueUnloader.CancelChannel = done

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				pollQueueDebugMsg("RELOAD QUEUE TICKER")
				pm.postNextPointCallback()
			}
		}
	}()
}

func (pm *NetworkPollManager) StopQueueUnloader() {
	pollQueueDebugMsg("StopQueueUnloader()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextUnloadTimer != nil && pm.PluginQueueUnloader.CancelChannel != nil {
		pm.PluginQueueUnloader.NextUnloadTimer.Stop()
		pm.PluginQueueUnloader.CancelChannel <- true
	}
	pm.PluginQueueUnloader = nil
}

//This function should be called from the Polling service. It will start a timer that posts the next polling point.
func (pm *NetworkPollManager) GetNextPollingPoint() (pp *PollingPoint, callback func(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate bool)) {
	pollQueueDebugMsg("GetNextPollingPoint()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint != nil {
		pp := pm.PluginQueueUnloader.NextPollPoint
		pm.PluginQueueUnloader.NextPollPoint = nil
		//Moving the line below to a reoccurring timer instead.
		//pm.PluginQueueUnloader.NextUnloadTimer = time.AfterFunc(pm.MaxPollRate, pm.postNextPointCallback)
		return pp, pm.PollingPointCompleteNotification
	}
	pollQueueDebugMsg("GetNextPollingPoint(): No pollingPoint available")
	return nil, nil
}

//This is the callback function that is called by the timer made in (pm *NetworkPollManager) GetNextPollingPoint().
func (pm *NetworkPollManager) postNextPointCallback() {
	pollQueueDebugMsg("postNextPointCallback()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint == nil {
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
		}
	}
}
