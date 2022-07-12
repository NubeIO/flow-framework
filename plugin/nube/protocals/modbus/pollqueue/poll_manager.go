package pollqueue

import (
	"container/heap"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/config"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

type NetworkPollManager struct {
	config       *config.Config
	DBHandlerRef *dbhandler.Handler

	Enable                 bool
	MaxPollRate            time.Duration
	PollQueue              *NetworkPriorityPollQueue
	PluginQueueUnloader    *QueueUnloader
	StatsCalcTimer         *time.Ticker
	PortUnavailableTimeout *time.Timer

	// References
	FFNetworkUUID string
	FFPluginUUID  string

	// Statistics
	AveragePollExecuteTimeSecs    float64       // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
	TotalPollQueueLength          int64         // number of polling points in the current queue.
	TotalStandbyPointsLength      int64         // number of polling points in the standby list.
	ASAPPriorityPollQueueLength   int64         // number of ASAP priority polling points in the current queue.
	HighPriorityPollQueueLength   int64         // number of High priority polling points in the current queue.
	NormalPriorityPollQueueLength int64         // number of Normal priority polling points in the current queue.
	LowPriorityPollQueueLength    int64         // number of Low priority polling points in the current queue.
	ASAPPriorityAveragePollTime   float64       // average time in seconds between ASAP priority polling point added to current queue, and polling complete.
	HighPriorityAveragePollTime   float64       // average time in seconds between High priority polling point added to current queue, and polling complete.
	NormalPriorityAveragePollTime float64       // average time in seconds between Normal priority polling point added to current queue, and polling complete.
	LowPriorityAveragePollTime    float64       // average time in seconds between Low priority polling point added to current queue, and polling complete.
	TotalPollCount                int64         // total number of polls completed.
	ASAPPriorityPollCount         int64         // total number of ASAP priority polls completed.
	HighPriorityPollCount         int64         // total number of High priority polls completed.
	NormalPriorityPollCount       int64         // total number of Normal priority polls completed.
	LowPriorityPollCount          int64         // total number of Low priority polls completed.
	ASAPPriorityPollCountForAvg   int64         // number of poll times included in avg polling time for ASAP priority (some are excluded because they have been modified while in the queue).
	HighPriorityPollCountForAvg   int64         // number of poll times included in avg polling time for High priority (some are excluded because they have been modified while in the queue).
	NormalPriorityPollCountForAvg int64         // number of poll times included in avg polling time for Normal priority (some are excluded because they have been modified while in the queue).
	LowPriorityPollCountForAvg    int64         // number of poll times included in avg polling time for Low priority (some are excluded because they have been modified while in the queue).
	ASAPPriorityMaxCycleTime      time.Duration // threshold setting for triggering a lockup alert for ASAP priority.
	HighPriorityMaxCycleTime      time.Duration // threshold setting for triggering a lockup alert for High priority.
	NormalPriorityMaxCycleTime    time.Duration // threshold setting for triggering a lockup alert for Normal priority.
	LowPriorityMaxCycleTime       time.Duration // threshold setting for triggering a lockup alert for Low priority.
	ASAPPriorityLockupAlert       bool          // alert if poll time has exceeded the ASAPPriorityMaxCycleTime
	HighPriorityLockupAlert       bool          // alert if poll time has exceeded the HighPriorityMaxCycleTime
	NormalPriorityLockupAlert     bool          // alert if poll time has exceeded the NormalPriorityMaxCycleTime
	LowPriorityLockupAlert        bool          // alert if poll time has exceeded the LowPriorityMaxCycleTime
	PollingStartTimeUnix          int64         // unix time (seconds) at polling start time.  Used for calculating Busy Time.
	BusyTime                      float64       // percent of the time that the plugin is actively polling.

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
	// pm.EmptyQueue()

	pm.StopQueueUnloader()
	// TODO: STOP ANY QUEUE LOADERS
	for _, pp := range pm.PollQueue.StandbyPollingPoints.PriorityQueue {
		if pp != nil && pp.RepollTimer != nil {
			pp.RepollTimer.Stop()
		}
	}
	pm.PollQueue.EmptyQueue()
}

func (pm *NetworkPollManager) PausePolling() { // POLLING SHOULD NOT BE PAUSED FOR LONG OR THE QUEUE WILL BECOME TOO LONG
	pm.Enable = false
	var nextPP *PollingPoint = nil
	if pm.PluginQueueUnloader.NextPollPoint != nil {
		nextPP = pm.PluginQueueUnloader.NextPollPoint
		pm.PollQueue.AddPollingPoint(nextPP) // add the next point back into the queue
	}
	pm.StopQueueUnloader()
}

func (pm *NetworkPollManager) UnpausePolling() {
	pm.Enable = true
	pm.PortUnavailableTimeout = nil
	pm.StartQueueUnloader()
}

func (pm *NetworkPollManager) EmptyQueue() {
	pm.StopQueueUnloader()
	pm.PollQueue.EmptyQueue()
}

func (pm *NetworkPollManager) ReAddDevicePoints(devUUID string) { // This is triggered by a user who wants to update the device poll times for standby points
	var arg api.Args
	arg.WithPoints = true
	dev, err := pm.DBHandlerRef.GetDevice(devUUID, arg)
	if dev == nil || err != nil {
		pm.pollQueueErrorMsg("ReAddDevicePoints(): cannot find device ", devUUID)
		return
	}
	pm.PollQueue.RemovePollingPointByDeviceUUID(devUUID)
	for _, pnt := range dev.Points {
		if boolean.IsTrue(pnt.Enable) {
			pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
			pp.PollPriority = pnt.PollPriority
			pm.PollQueue.AddPollingPoint(pp)
		}
	}
}

func NewPollManager(conf *config.Config, dbHandler *dbhandler.Handler, ffNetworkUUID, ffPluginUUID string) *NetworkPollManager {
	// Make the main priority polling queue
	queue := make([]*PollingPoint, 0)
	pq := &PriorityPollQueue{queue}
	heap.Init(pq)                        // Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	refQueue := make([]*PollingPoint, 0) // Make the reference slice that contains points that are not in the current polling queue.
	rq := &PriorityPollQueue{refQueue}
	heap.Init(rq)                                // Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	outstandingQueue := make([]*PollingPoint, 0) // Make the reference slice that contains points that are not in the current polling queue.
	opq := &PriorityPollQueue{outstandingQueue}
	heap.Init(opq)
	adl := make([]string, 0)
	pqu := &QueueUnloader{nil, nil, nil}
	puwp := make(map[string]bool)
	npq := &NetworkPriorityPollQueue{conf, pq, rq, opq, puwp, pqu, ffPluginUUID, ffNetworkUUID, adl}
	pm := new(NetworkPollManager)
	pm.config = conf
	pm.PollQueue = npq
	pm.PluginQueueUnloader = pqu
	pm.DBHandlerRef = dbHandler
	maxPollRate := 1000 * time.Millisecond
	pm.MaxPollRate = maxPollRate // TODO: MaxPollRate should come from a network property,but I can't find it. Also kinda implemented in StartQueueUnloader().
	pm.FFNetworkUUID = ffNetworkUUID
	pm.FFPluginUUID = ffPluginUUID
	pm.ASAPPriorityMaxCycleTime, _ = time.ParseDuration("2m")
	pm.HighPriorityMaxCycleTime, _ = time.ParseDuration("5m")
	pm.NormalPriorityMaxCycleTime, _ = time.ParseDuration("15m")
	pm.LowPriorityMaxCycleTime, _ = time.ParseDuration("60m")
	return pm
}

func (pm *NetworkPollManager) GetPollRateDuration(rate model.PollRate, deviceUUID string) time.Duration {
	pm.pollQueueDebugMsg("GetPollRateDuration(): ", rate)
	var arg api.Args
	device, err := pm.DBHandlerRef.GetDevice(deviceUUID, arg)
	if err != nil {
		pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.GetPollRateDuration(): couldn't find device %s/n", deviceUUID))
	}
	pm.pollQueueDebugMsg("GetPollRateDuration() device poll times: ", device.FastPollRate, device.NormalPollRate, device.SlowPollRate)

	var duration time.Duration
	switch rate {
	case model.RATE_FAST:
		pm.pollQueueDebugMsg("GetPollRateDuration(): FAST")
		fastRateDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", *device.FastPollRate))
		if fastRateDuration <= 100*time.Millisecond {
			duration = 10 * time.Second
		} else {
			duration = fastRateDuration
		}

	case model.RATE_NORMAL:
		pm.pollQueueDebugMsg("GetPollRateDuration(): NORMAL")
		normalRateDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", *device.NormalPollRate))
		if normalRateDuration <= 500*time.Millisecond {
			duration = 30 * time.Second
		} else {
			duration = normalRateDuration
		}

	case model.RATE_SLOW:
		pm.pollQueueDebugMsg("GetPollRateDuration(): SLOW")
		slowRateDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", *device.SlowPollRate))
		if slowRateDuration <= 1*time.Second {
			duration = 120 * time.Second
		} else {
			duration = slowRateDuration
		}

	default:
		pm.pollQueueDebugMsg("GetPollRateDuration(): UNKNOWN")
		normalRateDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", *device.NormalPollRate))
		if normalRateDuration <= 500*time.Millisecond {
			duration = 30 * time.Second
		} else {
			duration = normalRateDuration
		}
	}

	if duration.Milliseconds() <= 100 {
		duration = 30 * time.Second
		pm.pollQueueErrorMsg("NetworkPollManager.GetPollRateDuration: invalid PollRate duration. Set to 30 seconds/n")
	}
	return duration
}

func (pm *NetworkPollManager) PollingFinished(pp *PollingPoint, pollStartTime time.Time, writeSuccess, readSuccess bool, callback func(pp *PollingPoint, writeSuccess bool, readSuccess bool, pollTimeSecs float64, pointUpdate bool)) {
	pollEndTime := time.Now()
	pollDuration := pollEndTime.Sub(pollStartTime)
	pollTimeSecs := pollDuration.Seconds()
	callback(pp, writeSuccess, readSuccess, pollTimeSecs, false) // (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool)
}

func (pm *NetworkPollManager) CheckAllPointsExistInQueues() {
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, api.Args{WithDevices: true, WithPoints: true})
	if net == nil || err != nil {
		pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Network Not Found/n")
	}
	if boolean.IsFalse(net.Enable) { //If network isn't enabled, there should be no points in the polling queues
		if pm.PollQueue.PriorityQueue.Len() > 0 {
			pm.PollQueue.PriorityQueue.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Found PollingPoints in PriorityQueue of a disabled network/n")
		}
		if pm.PollQueue.StandbyPollingPoints.Len() > 0 {
			pm.PollQueue.StandbyPollingPoints.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Found PollingPoints in StandbyPollingPoints of a disabled network./n")
		}
		if pm.PollQueue.OutstandingPollingPoints.Len() > 0 {
			pm.PollQueue.OutstandingPollingPoints.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Found PollingPoints in OutstandingPollingPoints of a disabled network./n")
		}
	}
	for _, dev := range net.Devices {
		if dev.Points != nil {
			deviceExistsInQueue := pm.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID)
			if boolean.IsFalse(dev.Enable) {
				if deviceExistsInQueue {
					pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Device UUID exist in Poll Queues for a disabled device./n")
					pm.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)
					continue
				}
			}
			if boolean.IsTrue(dev.Enable) && !deviceExistsInQueue {
				pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Device UUID doesn't exist in active devices list./n")
			}
			for _, pnt := range dev.Points {
				if pnt != nil {
					if boolean.IsFalse(pnt.Enable) {
						pp, _ := pm.PollQueue.GetPollingPointByPointUUID(pnt.UUID)
						if pp != nil {
							pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Found disabled point in poll queue./n")
							pm.PollQueue.RemovePollingPointByPointUUID(pnt.UUID)
						}
						continue
					}
					if boolean.IsTrue(pnt.Enable) {
						pp, err := pm.PollQueue.GetPollingPointByPointUUID(pnt.UUID)
						if pp == nil || err != nil {
							pm.pollQueueErrorMsg("NetworkPollManager.CheckAllPointsExistInQueues: Polling point doesn't exist for point ", pnt.Name, "/n")
							pp = NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
							pp.PollPriority = pnt.PollPriority
							pm.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode.
						}
						continue
					}
				}
			}
		}
	}
}
