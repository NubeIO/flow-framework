package pollqueue

import (
	"container/heap"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/poller"
	log "github.com/sirupsen/logrus"
	"time"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

// Polling Manager Summary:
//  - Diagram Summary: https://docs.google.com/drawings/d/1priwsaQ6EryRBx1kLQd91REJvHzFyxz7cOHYYXyBNFE/edit?usp=sharing
//  - The PollManager Inserts, Removes, and Updates PollingPoints from the PriorityPollQueue based on the settings of the point.  It considers the Network/Device/Point Enables, Point Write-Modes, etc.
//  - When a ProtocolPollWorker finishes a poll, it will start a timer on the FF Point (based on poll rate) that will re-trigger the PollManager to create a PollingPoint if necessary.
//  -

//Questions:
// -

//There should be a function in Modbus(or other protocols) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

type NetworkPollManager struct {
	DBHandlerRef *dbhandler.Handler

	Enable              bool
	MaxPollRate         time.Duration
	PollQueue           *NetworkPriorityPollQueue
	PluginQueueUnloader *QueueUnloader

	//References
	FFNetworkUUID string
	FFPluginUUID  string

	//Statistics
	AveragePollTime               time.Duration
	TotalPollQueueLength          int
	TotalStandbyPointsLength      int
	HighPriorityPollQueueLength   int
	NormalPriorityPollQueueLength int
	LowPriorityPollQueueLength    int
	HighPriorityAveragePollTime   int
	NormalPriorityAveragePollTime int
	LowPriorityAveragePollTime    int
	BusyTime                      int //percent
	HighPriorityMaxCycleTime      time.Duration
	NormalPriorityMaxCycleTime    time.Duration
	LowPriorityMaxCycleTime       time.Duration
	HighPriorityLockupAlert       bool
	NormalPriorityLockupAlert     bool
	LowPriorityLockupAlert        bool
}

func (pm *NetworkPollManager) StartPolling() {
	if !pm.Enable {
		pm.Enable = true
		pm.RebuildPollingQueue()
	}
	if pm.PluginQueueUnloader == nil {
		pm.StartQueueUnloader()
	}
}

func (pm *NetworkPollManager) StopPolling() {
	pm.Enable = false
	//pm.EmptyQueue()

	pm.StopQueueUnloader()
	//TODO: STOP ANY QUEUE LOADERS
	for _, pp := range pm.PollQueue.StandbyPollingPoints.PriorityQueue {
		pp.RepollTimer.Stop()
	}
	pm.PollQueue.EmptyQueue()
}

func (pm *NetworkPollManager) PausePolling() { //POLLING SHOULD NOT BE PAUSED FOR LONG OR THE QUEUE WILL BECOME TOO LONG
	pm.Enable = false
	var nextPP *PollingPoint = nil
	if pm.PluginQueueUnloader.NextPollPoint != nil {
		nextPP = pm.PluginQueueUnloader.NextPollPoint
	}
	pm.StopQueueUnloader()
	pm.PollQueue.AddPollingPoint(nextPP) //add the next point back into the queue
}

func (pm *NetworkPollManager) UnpausePolling() {
	pm.Enable = true
	pm.StartQueueUnloader()
}

func (pm *NetworkPollManager) EmptyQueue() {
	pm.StopQueueUnloader()
	pm.PollQueue.EmptyQueue()
}

func NewPollManager(dbHandler *dbhandler.Handler, FFNetworkUUID, FFPluginUUID string) *NetworkPollManager {
	// Make the main priority polling queue
	queue := make([]*PollingPoint, 0)
	pq := &PriorityPollQueue{queue}
	heap.Init(pq) //Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	// Make the reference slice that contains points that are not in the current polling queue.
	refQueue := make([]*PollingPoint, 0)
	rq := &PriorityPollQueue{refQueue}
	npq := &NetworkPriorityPollQueue{pq, rq, FFPluginUUID, FFNetworkUUID}
	pqu := &QueueUnloader{nil, nil, nil}
	pm := new(NetworkPollManager)
	pm.PollQueue = npq
	pm.PluginQueueUnloader = pqu
	pm.DBHandlerRef = dbHandler
	maxpollrate := 1000 * time.Millisecond
	pm.MaxPollRate = maxpollrate //TODO: MaxPollRate should come from a network property,but I can't find it
	pm.FFNetworkUUID = FFNetworkUUID
	pm.FFPluginUUID = FFPluginUUID
	return pm
}

func (pm *NetworkPollManager) GetPollRateDuration(rate poller.PollRate, deviceUUID string) time.Duration {
	//fmt.Println("GetPollRateDuration()")
	var arg api.Args
	device, err := pm.DBHandlerRef.GetDevice(deviceUUID, arg)
	if err != nil {
		fmt.Printf("NetworkPollManager.GetPollRateDuration(): couldn't find device %s/n", deviceUUID)
	}

	var duration time.Duration
	switch rate {
	case poller.RATE_FAST:
		if device.FastPollRate == nil {
			duration = 60 * time.Second
		} else {
			duration = *device.FastPollRate
		}
	case poller.RATE_NORMAL:
		if device.NormalPollRate == nil {
			duration = 10 * time.Second
		} else {
			duration = *device.NormalPollRate
		}
	case poller.RATE_SLOW:
		if device.SlowPollRate == nil {
			duration = 60 * time.Second
		} else {
			duration = *device.SlowPollRate
		}
	default:
		if device.NormalPollRate == nil {
			duration = 60 * time.Second
		} else {
			duration = *device.NormalPollRate
		}
	}

	if duration.Milliseconds() <= 100 {
		duration = 60 * time.Second
		log.Info("NetworkPollManager.GetPollRateDuration: invalid PollRate duration. Set to 60 seconds/n")
	}
	return duration
}
