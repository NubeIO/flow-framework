package poller

import (
	"container/heap"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/dbhandler"
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
	Enable               bool
	MaxPollRate          *time.Duration
	PollQueue            *NetworkPriorityPollQueue
	NetworkQueueUnloader *QueueUnloader

	//References
	FFNetworkUUID string
	FFPluginUUID  string

	//Statistics
	AveragePollTime               time.Duration
	TotalPollQueueLength          int
	HighPriorityPollQueueLength   int
	NormalPriorityPollQueueLength int
	SlowPriorityPollQueueLength   int
	HighPriorityPollsPerMinute    int
	NormalPriorityPollsPerMinute  int
	SlowPriorityPollsPerMinute    int
	BusyTime                      int //percent
}

func (pm *NetworkPollManager) StartPolling() {
	if !pm.Enable {
		pm.Enable = true
		pm.RebuildPollingQueue()
	}
	if pm.NetworkQueueUnloader == nil {
		pm.StartQueueUnloader()
	}
}

func (pm *NetworkPollManager) StopPolling() {
	pm.Enable = false
	//TODO: STOP ANY QUEUE LOADERS
	pm.EmptyQueue()
	pm.StopQueueUnloader()
}

func (pm *NetworkPollManager) PausePolling() { //POLLING SHOULD NOT BE PAUSED FOR LONG OR THE QUEUE WILL BECOME TOO LONG
	pm.Enable = false
	var nextPP *PollingPoint = nil
	if pm.NetworkQueueUnloader.NextPollPoint != nil {
		nextPP = pm.NetworkQueueUnloader.NextPollPoint
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

func NewPollManager(FFNetworkUUID, FFPluginUUID string) *NetworkPollManager {
	queue := make([]*PollingPoint, 0)
	pq := &PriorityPollQueue{false, queue}
	heap.Init(pq)
	nq := &NetworkPriorityPollQueue{pq, FFNetworkUUID}
	pm := new(NetworkPollManager)
	pm.PollQueue = nq
	maxpollrate := 1000 * time.Millisecond
	pm.MaxPollRate = &maxpollrate //TODO: MaxPollRate should come from a network property,but I can't find it
	pm.FFNetworkUUID = FFNetworkUUID
	pm.FFPluginUUID = FFPluginUUID
	return pm
}

func (pm *NetworkPollManager) GetPollRateDuration(rate PollRate, deviceUUID string) time.Duration {
	h := &dbhandler.Handler{}
	dbhandler.Init(h)
	var arg api.Args
	device, err := h.DB.GetDevice(deviceUUID, arg)
	if err != nil {
		fmt.Printf("NetworkPollManager.GetPollRateDuration(): couldn't find device %s/n", deviceUUID)
	}

	var duration time.Duration
	switch rate {
	case RATE_FAST:
		duration = *device.FastPollRate
	case RATE_NORMAL:
		duration = *device.NormalPollRate
	case RATE_SLOW:
		duration = *device.SlowPollRate
	default:
		duration = *device.NormalPollRate
	}
	if duration.Milliseconds() <= 100 {
		duration = 60 * time.Second
		log.Info("NetworkPollManager.GetPollRateDuration: invalid PollRate duration. Set to 60 seconds/n")
	}
	return duration
}
