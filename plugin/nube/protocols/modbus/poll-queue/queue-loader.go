package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
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
//  - The QueueLoader puts PollPoints into the Queue

//Questions:
// -

//There should be a function in Modbus(or other protocols) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

func (pm *NetworkPollManager) RebuildPollingQueue() error {
	//TODO: STOP ANY OTHER QUEUE LOADERS
	fmt.Println("RebuildPollingQueue()")
	wasRunning := pm.PluginQueueUnloader != nil
	pm.EmptyQueue()
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, arg)
	if err != nil || len(net.Devices) == 0 {
		return errors.New(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: couldn't find any devices for the network %s/n", pm.FFNetworkUUID))
	}
	devs := net.Devices
	for _, dev := range devs { //DEVICES
		if dev.NetworkUUID == pm.FFNetworkUUID && utils.BoolIsNil(dev.Enable) {
			for _, pnt := range dev.Points { //POINTS
				if pnt.DeviceUUID == dev.UUID && utils.BoolIsNil(pnt.Enable) {
					pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
					pp.PollPriority = pnt.PollPriority
					//fmt.Println("RebuildPollingQueue() pp:")
					//fmt.Printf("%+v\n", pp)
					pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
					pm.PollQueue.AddPollingPoint(pp)
				} else {
					log.Info(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Point (%s) is not enabled./n", pnt.UUID))
				}
			}
		} else {
			log.Info(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Device (%s) is not enabled./n", dev.UUID))
		}
	}
	heap.Init(pm.PollQueue.PriorityQueue)
	if wasRunning {
		pm.StartQueueUnloader()
	}
	//TODO: START ANY OTHER REQUIRED QUEUE LOADERS/OPTIMIZERS
	pm.PrintPollQueuePointUUIDs()
	return nil
}

func (pm *NetworkPollManager) PrintPollQueuePointUUIDs() {
	fmt.Println("")
	hasNextPollPoint := 0
	if pm.PluginQueueUnloader.NextPollPoint != nil {
		hasNextPollPoint = 1
	}
	fmt.Println("PrintPollQueuePointUUIDs TOTAL COUNT = ", hasNextPollPoint+pm.PollQueue.PriorityQueue.Len()+pm.PollQueue.StandbyPollingPoints.Len())
	fmt.Print("NextPollPoint: ")
	fmt.Printf("%+v\n", pm.PluginQueueUnloader.NextPollPoint)
	fmt.Print("PollQueue: COUNT = ", pm.PollQueue.PriorityQueue.Len(), ": ")
	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		fmt.Print(pp.FFPointUUID, " - ", pp.PollPriority, "; ")
	}
	fmt.Println("")
	fmt.Print("StandbyPollingPoints COUNT = ", pm.PollQueue.StandbyPollingPoints.Len(), ": ")
	for _, pp := range pm.PollQueue.StandbyPollingPoints.PriorityQueue {
		fmt.Print(pp.FFPointUUID, " - ", pp.PollPriority, ", repoll timer:", pp.RepollTimer != nil, "; ")
	}
	fmt.Println("\n \n")
}

func (pm *NetworkPollManager) PollCompleteStatsUpdate(pp *PollingPoint, pollTimeSecs float64) {
	fmt.Println("PollCompleteStatsUpdate()")

	pm.AveragePollExecuteTimeSecs = ((pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) + pollTimeSecs) / (float64(pm.TotalPollCount) + 1)
	pm.TotalPollCount++
	enabledTime := time.Since(time.Unix(pm.PollingStartTimeUnix, 0)) * time.Second
	pm.BusyTime = (pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) / enabledTime.Seconds()

	switch pp.PollPriority {
	case poller.PRIORITY_ASAP:
		pm.ASAPPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.ASAPPriorityAveragePollTime = ((pm.ASAPPriorityAveragePollTime * float64(pm.ASAPPriorityPollCountForAvg)) + pollTime) / (float64(pm.ASAPPriorityPollCountForAvg) + 1)
		pm.ASAPPriorityPollCountForAvg++

	case poller.PRIORITY_HIGH:
		pm.HighPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.HighPriorityAveragePollTime = ((pm.HighPriorityAveragePollTime * float64(pm.HighPriorityPollCountForAvg)) + pollTime) / (float64(pm.HighPriorityPollCountForAvg) + 1)
		pm.HighPriorityPollCountForAvg++

	case poller.PRIORITY_NORMAL:
		pm.NormalPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.NormalPriorityAveragePollTime = ((pm.NormalPriorityAveragePollTime * float64(pm.NormalPriorityPollCountForAvg)) + pollTime) / (float64(pm.NormalPriorityPollCountForAvg) + 1)
		pm.NormalPriorityPollCountForAvg++

	case poller.PRIORITY_LOW:
		pm.LowPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.LowPriorityAveragePollTime = ((pm.LowPriorityAveragePollTime * float64(pm.LowPriorityPollCountForAvg)) + pollTime) / (float64(pm.LowPriorityPollCountForAvg) + 1)
		pm.LowPriorityPollCountForAvg++

	}

}

func (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate bool) {
	log.Infof("modbus-poll: PollingPointCompleteNotification Point UUID: %s, writeSuccess: %t, readSuccess: %t, pollTime: %f", pp.FFPointUUID, writeSuccess, readSuccess, pollTimeSecs)

	if !pointUpdate {
		pm.PollCompleteStatsUpdate(pp, pollTimeSecs) // This will update the relevant PollManager statistics.
	}

	point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
	if point == nil || err != nil {
		fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s /n", pp.FFPointUUID)
		return
	}
	log.Infof("modbus-poll: PollingPointCompleteNotification WriteMode: %s", point.WriteMode)

	fmt.Println("PollingPointCompleteNotification: point")
	fmt.Printf("%+v\n", point)
	point.PrintPointValues()

	//If the device was deleted while this point was being polled, discard the PollingPoint
	if !pointUpdate && !pm.PollQueue.CheckIfActiveDevicesListIncludes(point.DeviceUUID) {
		return
	}

	fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): writeMode: %s", point.WriteMode)
	fmt.Println("")
	switch point.WriteMode {
	case poller.ReadOnce: //ReadOnce          If read_successful then don't re-add.
		point.WritePollRequired = utils.NewFalse()
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			point.ReadPollRequired = utils.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)
		}

	case poller.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.WritePollRequired = utils.NewFalse()
		fmt.Println("PollingPointCompleteNotification: ReadOnly")
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			log.Info("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			fmt.Println("PollingPointCompleteNotification: NOT readSuccess")
			point.ReadPollRequired = utils.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			fmt.Println("PollingPointCompleteNotification: ABOUT TO ADD POINT")
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}

	case poller.WriteOnce: //WriteOnce         If write_successful then don't re-add.
		point.ReadPollRequired = utils.NewFalse()
		if writeSuccess {
			point.WritePollRequired = utils.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			point.WritePollRequired = utils.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case poller.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		if utils.BoolIsNil(point.WritePollRequired) && writeSuccess {
			point.WritePollRequired = utils.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (utils.BoolIsNil(point.WritePollRequired) && !writeSuccess) {
			point.WritePollRequired = utils.NewTrue()
			if pointUpdate {
				point.ReadPollRequired = utils.NewTrue()
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
			break
		}
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if utils.BoolIsNil(point.ReadPollRequired) && !readSuccess {
			point.ReadPollRequired = utils.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case poller.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = utils.NewFalse()
		point.WritePollRequired = utils.NewTrue()
		if writeSuccess {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			//log.Info("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case poller.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = utils.NewTrue()
		if utils.BoolIsNil(point.WritePollRequired) && writeSuccess {
			point.WritePollRequired = utils.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (utils.BoolIsNil(point.WritePollRequired) && !writeSuccess) {
			point.WritePollRequired = utils.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
			break
		}
		if readSuccess {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			//log.Info("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case poller.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = utils.NewTrue()
		writeValuePointer := point.Priority.GetHighestPriorityValue()
		if writeValuePointer != nil {
			writeValue := *writeValuePointer
			presentValue := *point.PresentValue
			if presentValue != writeValue {
				if pp.RepollTimer != nil {
					pp.RepollTimer.Stop()
					pp.RepollTimer = nil
				}
				point.WritePollRequired = utils.NewTrue()
				pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
				pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
			} else {
				point.WritePollRequired = utils.NewFalse()
				duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
				log.Info("duration: ", duration)
				// This line sets a timer to re-add the point to the poll queue after the PollRate time.
				pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
				addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
				if !addSuccess {
					log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
				}
			}
		}
	}

	fmt.Println("PollingPointCompleteNotification AFTER (ABOUT TO DB UPDATE): point")
	fmt.Printf("%+v\n", point)
	point.PrintPointValues()
	//TODO: DONT KNOW WHY THIS WAS HERE, REMOVED TO TEST
	pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)
}

func (pm *NetworkPollManager) MakePollingPointRepollCallback(pp *PollingPoint, writeMode poller.WriteMode) func() {
	//log.Info("MakePollingPointRepollCallback()")
	f := func() {
		log.Info("CALL PollingPointRepollCallback func() pp:")
		fmt.Printf("%+v\n", pp)
		pp.RepollTimer = nil
		removeSuccess := pm.PollQueue.StandbyPollingPoints.RemovePollingPointByPointUUID(pp.FFPointUUID)
		if !removeSuccess {
			log.Error(fmt.Sprintf("Modbus MakePollingPointRepollCallback(): polling point could not be found in StandbyPollingPoints.  (%s)", pp.FFPointUUID))
		}

		point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{})
		if point == nil || err != nil {
			fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s /n", pp.FFPointUUID)
			return
		}

		switch writeMode {
		case poller.ReadOnce:
			return

		case poller.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = utils.NewTrue()
			point.WritePollRequired = utils.NewFalse()

		case poller.WriteOnce: //WriteOnce         If write_successful then don't re-add.
			return

		case poller.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
			return

		case poller.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
			point.ReadPollRequired = utils.NewFalse()
			point.WritePollRequired = utils.NewTrue()

		case poller.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = utils.NewTrue()
			point.WritePollRequired = utils.NewFalse()

		case poller.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
			point.ReadPollRequired = utils.NewTrue()
			point.WritePollRequired = utils.NewFalse()
		}

		//Now add the polling point back to the polling queue
		pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
		pm.PollQueue.AddPollingPoint(pp)
	}
	return f
}

func (pm *NetworkPollManager) MakeLockupTimerFunc(priority poller.PollPriority) *time.Timer {
	timeoutDuration := 5 * time.Minute

	switch priority {
	case poller.PRIORITY_ASAP:
		timeoutDuration = pm.ASAPPriorityMaxCycleTime

	case poller.PRIORITY_HIGH:
		timeoutDuration = pm.HighPriorityMaxCycleTime

	case poller.PRIORITY_NORMAL:
		timeoutDuration = pm.NormalPriorityMaxCycleTime

	case poller.PRIORITY_LOW:
		timeoutDuration = pm.LowPriorityMaxCycleTime

	}

	f := func() {
		log.Infof("Polling Lockout Timer Expired! Polling Priority: %d,  Polling Network: %s", priority, pm.FFNetworkUUID)
		switch priority {
		case poller.PRIORITY_ASAP:
			pm.ASAPPriorityLockupAlert = true

		case poller.PRIORITY_HIGH:
			pm.HighPriorityLockupAlert = true

		case poller.PRIORITY_NORMAL:
			pm.NormalPriorityLockupAlert = true

		case poller.PRIORITY_LOW:
			pm.LowPriorityLockupAlert = true

		}
	}

	return time.AfterFunc(timeoutDuration, f)
}

func (pm *NetworkPollManager) SetPointPollRequiredFlagsBasedOnWriteMode(point *model.Point) {
	fmt.Println("SetPointPollRequiredFlagsBasedOnWriteMode BEFORE: point")
	fmt.Printf("%+v\n", point)
	fmt.Println("MODBUS SetPointPollRequiredFlagsBasedOnWriteMode(): PRIORITY")
	fmt.Printf("%+v\n", point.Priority)

	if point == nil {
		fmt.Printf("NetworkPollManager.SetPointPollRequiredFlagsBasedOnWriteMode(): couldn't find point %s /n", point.UUID)
		return
	}

	switch point.WriteMode {
	case poller.ReadOnce:
		return

	case poller.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = utils.NewTrue()
		point.WritePollRequired = utils.NewFalse()

	case poller.WriteOnce: //WriteOnce         If write_successful then don't re-add.
		return

	case poller.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		return

	case poller.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = utils.NewFalse()
		point.WritePollRequired = utils.NewTrue()

	case poller.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = utils.NewTrue()
		point.WritePollRequired = utils.NewTrue()

	case poller.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = utils.NewTrue()
		point.WritePollRequired = utils.NewTrue()
	}

	fmt.Println("SetPointPollRequiredFlagsBasedOnWriteMode AFTER: point")
	fmt.Printf("%+v\n", point)
	fmt.Println("MODBUS SetPointPollRequiredFlagsBasedOnWriteMode(): PRIORITY")
	fmt.Printf("%+v\n", point.Priority)

	pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)

}
