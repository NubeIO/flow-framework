package pollqueue

import (
	"container/heap"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/config"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"math"
	"time"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

type NetworkPollManager struct {
	config       *config.Config
	DBHandlerRef *dbhandler.Handler

	Enable                    bool
	MaxPollRate               time.Duration
	PollQueue                 *NetworkPriorityPollQueue
	PluginQueueUnloader       *QueueUnloader
	StatsCalcTimer            *time.Ticker
	PortUnavailableTimeout    *time.Timer
	QueueCheckerTimer         *time.Ticker
	QueueCheckerCancelChannel chan bool

	// References
	FFNetworkUUID string
	FFPluginUUID  string

	// Statistics
	MaxPollExecuteTimeSecs        float64       // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
	AveragePollExecuteTimeSecs    float64       // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
	MinPollExecuteTimeSecs        float64       // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
	TotalPollQueueLength          int64         // number of polling points in the current queue.
	TotalStandbyPointsLength      int64         // number of polling points in the standby list.
	TotalPointsOutForPolling      int64         // number of points currently out for polling (currently being handled by the protocol plugin).
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
	EnabledTime                   float64       // time in seconds that the statistics have been running for.
	PortUnavailableTime           float64       // time in seconds that the serial port has been unavailable.
	PortUnavailableStartTime      int64         // unix time (seconds) when port became unavailable.  Used for calculating downtime.
}

func (pm *NetworkPollManager) StartPolling() {
	if !pm.Enable {
		pm.Enable = true
		pm.RebuildPollingQueue()
	} else if pm.PluginQueueUnloader == nil {
		pm.StartQueueUnloader()
	}
	//Also start the Queue Checker
	pm.StartQueueCheckerAndStats()
	pm.StartPollingStatistics()
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
	if pm.QueueCheckerTimer != nil && pm.QueueCheckerCancelChannel != nil {
		pm.StopQueueCheckerAndStats()
	}
}

func (pm *NetworkPollManager) PausePolling() { // POLLING SHOULD NOT BE PAUSED FOR LONG OR THE QUEUE WILL BECOME TOO LONG
	pm.pollQueueDebugMsg("PausePolling()")
	pm.Enable = false
	var nextPP *PollingPoint = nil
	if pm.PluginQueueUnloader.NextPollPoint != nil {
		nextPP = pm.PluginQueueUnloader.NextPollPoint
		pm.PollQueue.AddPollingPoint(nextPP) // add the next point back into the queue
	}
	pm.StopQueueUnloader()
}

func (pm *NetworkPollManager) UnpausePolling() {
	pm.pollQueueDebugMsg("UnpausePolling()")
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

func (pm *NetworkPollManager) PollQueueErrorChecking() {
	pm.pollQueueDebugMsg("NetworkPollManager.PollQueueErrorChecking")
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, api.Args{WithDevices: true, WithPoints: true})
	if net == nil || err != nil {
		pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Network Not Found/n")
	}
	if boolean.IsFalse(net.Enable) { //If network isn't enabled, there should be no points in the polling queues
		if pm.PollQueue.PriorityQueue.Len() > 0 {
			pm.PollQueue.PriorityQueue.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Found PollingPoints in PriorityQueue of a disabled network/n")
		}
		if pm.PollQueue.StandbyPollingPoints.Len() > 0 {
			pm.PollQueue.StandbyPollingPoints.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Found PollingPoints in StandbyPollingPoints of a disabled network./n")
		}
		if pm.PollQueue.OutstandingPollingPoints.Len() > 0 {
			pm.PollQueue.OutstandingPollingPoints.EmptyQueue()
			pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Found PollingPoints in OutstandingPollingPoints of a disabled network./n")
		}
	}
	for _, dev := range net.Devices {
		deviceExistsInQueue := pm.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID)
		if boolean.IsFalse(net.Enable) || boolean.IsFalse(dev.Enable) {
			if deviceExistsInQueue {
				pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Device UUID exist in Poll Queues for a disabled device./n")
				pm.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)
				continue
			}
		}
		if boolean.IsTrue(net.Enable) && boolean.IsTrue(dev.Enable) && !deviceExistsInQueue {
			pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Device UUID doesn't exist in active devices list./n")
		}
		if dev.Points != nil {
			for _, pnt := range dev.Points {
				if pnt != nil {
					if boolean.IsFalse(net.Enable) || boolean.IsFalse(dev.Enable) || boolean.IsFalse(pnt.Enable) {
						pp, _ := pm.PollQueue.GetPollingPointByPointUUID(pnt.UUID)
						if pp != nil {
							pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Found disabled point in poll queue./n")
							pm.PollQueue.RemovePollingPointByPointUUID(pnt.UUID)
						}
						continue
					}
					if boolean.IsTrue(net.Enable) && boolean.IsTrue(dev.Enable) && boolean.IsTrue(pnt.Enable) {
						pp, err := pm.PollQueue.GetPollingPointByPointUUID(pnt.UUID)
						if pp == nil || err != nil {
							pm.pollQueueErrorMsg("NetworkPollManager.PollQueueErrorChecking: Polling point doesn't exist for point ", pnt.Name, "/n")
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

func (pm *NetworkPollManager) StartQueueCheckerAndStats() {
	period := 30 * time.Second
	pm.MaxPollRate = period
	if pm.QueueCheckerTimer != nil {
		pm.QueueCheckerTimer.Stop()
	}
	pm.QueueCheckerTimer = nil
	if pm.QueueCheckerCancelChannel != nil {
		pm.QueueCheckerCancelChannel <- true
	}
	pm.QueueCheckerCancelChannel = nil
	ticker := time.NewTicker(period)
	pm.QueueCheckerTimer = ticker
	done := make(chan bool)
	pm.QueueCheckerCancelChannel = done
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				//pm.pollQueueDebugMsg("RELOAD QUEUE TICKER")
				pm.PollQueueErrorChecking()
				pm.PrintPollQueueStatistics()
			}
		}
	}()
}

func (pm *NetworkPollManager) StopQueueCheckerAndStats() {
	pm.QueueCheckerTimer.Stop()
	pm.QueueCheckerTimer = nil
	pm.QueueCheckerCancelChannel <- true
	pm.QueueCheckerCancelChannel = nil
}

func (pm *NetworkPollManager) StartPollingStatistics() {
	pm.pollQueueDebugMsg("StartPollingStatistics()")
	pm.PollingStartTimeUnix = time.Now().Unix()
	pm.AveragePollExecuteTimeSecs = 0
	pm.MaxPollExecuteTimeSecs = 0
	pm.MinPollExecuteTimeSecs = 0
	pm.ASAPPriorityAveragePollTime = 0
	pm.HighPriorityAveragePollTime = 0
	pm.NormalPriorityAveragePollTime = 0
	pm.LowPriorityAveragePollTime = 0
	pm.TotalPollCount = 0
	pm.ASAPPriorityPollCount = 0
	pm.HighPriorityPollCount = 0
	pm.NormalPriorityPollCount = 0
	pm.LowPriorityPollCount = 0
	pm.ASAPPriorityPollCountForAvg = 0
	pm.HighPriorityPollCountForAvg = 0
	pm.NormalPriorityPollCountForAvg = 0
	pm.LowPriorityPollCountForAvg = 0
	pm.ASAPPriorityLockupAlert = false
	pm.HighPriorityLockupAlert = false
	pm.NormalPriorityLockupAlert = false
	pm.LowPriorityLockupAlert = false
	pm.PortUnavailableTime = 0
	pm.PortUnavailableStartTime = 0
}

func (pm *NetworkPollManager) PollCompleteStatsUpdate(pp *PollingPoint, pollTimeSecs float64) {
	pm.pollQueueDebugMsg("PollCompleteStatsUpdate()")

	if pm.MaxPollExecuteTimeSecs == 0 || pollTimeSecs > pm.MaxPollExecuteTimeSecs {
		pm.MaxPollExecuteTimeSecs = pollTimeSecs
	}
	if pm.MinPollExecuteTimeSecs == 0 || pollTimeSecs < pm.MinPollExecuteTimeSecs {
		pm.MinPollExecuteTimeSecs = pollTimeSecs
	}
	pm.AveragePollExecuteTimeSecs = ((pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) + pollTimeSecs) / (float64(pm.TotalPollCount) + 1)
	pm.TotalPollCount++
	pm.EnabledTime = time.Since(time.Unix(pm.PollingStartTimeUnix, 0)).Seconds()
	pm.BusyTime = math.Round((((pm.AveragePollExecuteTimeSecs*float64(pm.TotalPollCount))/pm.EnabledTime)*100)*1000) / 1000 // percentage rounded to 3 decimal places

	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
	pm.TotalStandbyPointsLength = int64(pm.PollQueue.StandbyPollingPoints.Len())
	pm.TotalPointsOutForPolling = int64(pm.PollQueue.OutstandingPollingPoints.Len())

	pm.ASAPPriorityPollQueueLength = 0
	pm.HighPriorityPollQueueLength = 0
	pm.NormalPriorityPollQueueLength = 0
	pm.LowPriorityPollQueueLength = 0

	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		if pp != nil {
			switch pp.PollPriority {
			case model.PRIORITY_ASAP:
				pm.ASAPPriorityPollQueueLength++
			case model.PRIORITY_HIGH:
				pm.HighPriorityPollQueueLength++
			case model.PRIORITY_NORMAL:
				pm.NormalPriorityPollQueueLength++
			case model.PRIORITY_LOW:
				pm.LowPriorityPollQueueLength++
			}
		}
	}
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())

	switch pp.PollPriority {
	case model.PRIORITY_ASAP:
		pm.ASAPPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.ASAPPriorityAveragePollTime = ((pm.ASAPPriorityAveragePollTime * float64(pm.ASAPPriorityPollCountForAvg)) + pollTime) / (float64(pm.ASAPPriorityPollCountForAvg) + 1)
		pm.ASAPPriorityPollCountForAvg++

	case model.PRIORITY_HIGH:
		pm.HighPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.HighPriorityAveragePollTime = ((pm.HighPriorityAveragePollTime * float64(pm.HighPriorityPollCountForAvg)) + pollTime) / (float64(pm.HighPriorityPollCountForAvg) + 1)
		pm.HighPriorityPollCountForAvg++

	case model.PRIORITY_NORMAL:
		pm.NormalPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.NormalPriorityAveragePollTime = ((pm.NormalPriorityAveragePollTime * float64(pm.NormalPriorityPollCountForAvg)) + pollTime) / (float64(pm.NormalPriorityPollCountForAvg) + 1)
		pm.NormalPriorityPollCountForAvg++

	case model.PRIORITY_LOW:
		pm.LowPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.LowPriorityAveragePollTime = ((pm.LowPriorityAveragePollTime * float64(pm.LowPriorityPollCountForAvg)) + pollTime) / (float64(pm.LowPriorityPollCountForAvg) + 1)
		pm.LowPriorityPollCountForAvg++

	}

}

func (pm *NetworkPollManager) PartialPollStatsUpdate() {
	pm.pollQueueDebugMsg("PartialPollStatsUpdate()")
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
	pm.TotalStandbyPointsLength = int64(pm.PollQueue.StandbyPollingPoints.Len())
	pm.TotalPointsOutForPolling = int64(pm.PollQueue.OutstandingPollingPoints.Len())

	pm.EnabledTime = time.Since(time.Unix(pm.PollingStartTimeUnix, 0)).Seconds()

	if pm.PortUnavailableTimeout != nil {
		pm.PortUnavailableTime += time.Since(time.Unix(pm.PortUnavailableStartTime, 0)).Seconds()
		pm.PortUnavailableStartTime = time.Now().Unix()
	}

	pm.ASAPPriorityPollQueueLength = 0
	pm.HighPriorityPollQueueLength = 0
	pm.NormalPriorityPollQueueLength = 0
	pm.LowPriorityPollQueueLength = 0

	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		if pp != nil {
			switch pp.PollPriority {
			case model.PRIORITY_ASAP:
				pm.ASAPPriorityPollQueueLength++
			case model.PRIORITY_HIGH:
				pm.HighPriorityPollQueueLength++
			case model.PRIORITY_NORMAL:
				pm.NormalPriorityPollQueueLength++
			case model.PRIORITY_LOW:
				pm.LowPriorityPollQueueLength++
			}
		}
	}
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
}

func (pm *NetworkPollManager) PortUnavailable() {
	pm.PortUnavailableStartTime = time.Now().Unix()
	pm.PausePolling()
}

func (pm *NetworkPollManager) PortAvailable() {
	pm.PartialPollStatsUpdate()
	pm.PrintPollQueueStatistics()
	pm.UnpausePolling()
}
