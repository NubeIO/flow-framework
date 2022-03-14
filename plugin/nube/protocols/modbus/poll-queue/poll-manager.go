package pollqueue

import (
	"container/heap"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
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
	StatsCalcTimer      time.Ticker

	//References
	FFNetworkUUID string
	FFPluginUUID  string

	//Statistics
	AveragePollExecuteTimeSecs    float64       //time in seconds for polling to complete (poll response time, doesn't include the time in queue).
	TotalPollQueueLength          int64         //number of polling points in the current queue.
	TotalStandbyPointsLength      int64         //number of polling points in the standby list.
	ASAPPriorityPollQueueLength   int64         //number of ASAP priority polling points in the current queue.
	HighPriorityPollQueueLength   int64         //number of High priority polling points in the current queue.
	NormalPriorityPollQueueLength int64         //number of Normal priority polling points in the current queue.
	LowPriorityPollQueueLength    int64         //number of Low priority polling points in the current queue.
	ASAPPriorityAveragePollTime   float64       //average time in seconds between ASAP priority polling point added to current queue, and polling complete.
	HighPriorityAveragePollTime   float64       //average time in seconds between High priority polling point added to current queue, and polling complete.
	NormalPriorityAveragePollTime float64       //average time in seconds between Normal priority polling point added to current queue, and polling complete.
	LowPriorityAveragePollTime    float64       //average time in seconds between Low priority polling point added to current queue, and polling complete.
	TotalPollCount                int64         //total number of polls completed.
	ASAPPriorityPollCount         int64         //total number of ASAP priority polls completed.
	HighPriorityPollCount         int64         //total number of High priority polls completed.
	NormalPriorityPollCount       int64         //total number of Normal priority polls completed.
	LowPriorityPollCount          int64         //total number of Low priority polls completed.
	ASAPPriorityPollCountForAvg   int64         //number of poll times included in avg polling time for ASAP priority (some are excluded because they have been modified while in the queue).
	HighPriorityPollCountForAvg   int64         //number of poll times included in avg polling time for High priority (some are excluded because they have been modified while in the queue).
	NormalPriorityPollCountForAvg int64         //number of poll times included in avg polling time for Normal priority (some are excluded because they have been modified while in the queue).
	LowPriorityPollCountForAvg    int64         //number of poll times included in avg polling time for Low priority (some are excluded because they have been modified while in the queue).
	ASAPPriorityMaxCycleTime      time.Duration //threshold setting for triggering a lockup alert for ASAP priority.
	HighPriorityMaxCycleTime      time.Duration //threshold setting for triggering a lockup alert for High priority.
	NormalPriorityMaxCycleTime    time.Duration //threshold setting for triggering a lockup alert for Normal priority.
	LowPriorityMaxCycleTime       time.Duration //threshold setting for triggering a lockup alert for Low priority.
	ASAPPriorityLockupAlert       bool          //alert if poll time has exceeded the ASAPPriorityMaxCycleTime
	HighPriorityLockupAlert       bool          //alert if poll time has exceeded the HighPriorityMaxCycleTime
	NormalPriorityLockupAlert     bool          //alert if poll time has exceeded the NormalPriorityMaxCycleTime
	LowPriorityLockupAlert        bool          //alert if poll time has exceeded the LowPriorityMaxCycleTime
	PollingStartTimeUnix          int64         //unix time (seconds) at polling start time.  Used for calculating Busy Time.
	BusyTime                      float64       //percent of the time that the plugin is actively polling.

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

func (pm *NetworkPollManager) ReAddDevicePoints(devUUID string) { //This is triggered by a user who wants to update the device poll times for standby points
	var arg api.Args
	arg.WithPoints = true
	dev, err := pm.DBHandlerRef.GetDevice(devUUID, arg)
	if dev == nil || err != nil {
		log.Error("Modbus: ReAddDevicePoints(): cannot find device ", devUUID)
		return
	}
	pm.PollQueue.RemovePollingPointByDeviceUUID(devUUID)
	for _, pnt := range dev.Points {
		if utils.BoolIsNil(pnt.Enable) {
			pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
			pp.PollPriority = pnt.PollPriority
			pm.PollQueue.AddPollingPoint(pp)
		}
	}
}

func NewPollManager(dbHandler *dbhandler.Handler, FFNetworkUUID, FFPluginUUID string) *NetworkPollManager {
	// Make the main priority polling queue
	queue := make([]*PollingPoint, 0)
	pq := &PriorityPollQueue{queue}
	heap.Init(pq) //Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	// Make the reference slice that contains points that are not in the current polling queue.
	refQueue := make([]*PollingPoint, 0)
	rq := &PriorityPollQueue{refQueue}
	adl := make([]string, 0)
	npq := &NetworkPriorityPollQueue{pq, rq, FFPluginUUID, FFNetworkUUID, adl}
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
