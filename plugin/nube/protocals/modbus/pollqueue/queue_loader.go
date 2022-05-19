package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

func (pm *NetworkPollManager) RebuildPollingQueue() error {
	//TODO: STOP ANY OTHER QUEUE LOADERS
	pm.pollQueueDebugMsg("RebuildPollingQueue()")
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
		if dev.NetworkUUID == pm.FFNetworkUUID && boolean.IsTrue(dev.Enable) {
			for _, pnt := range dev.Points { //POINTS
				if pnt.DeviceUUID == dev.UUID && boolean.IsTrue(pnt.Enable) {
					pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
					pp.PollPriority = pnt.PollPriority
					pm.pollQueueDebugMsg(fmt.Sprintf("RebuildPollingQueue() pp: %+v", pp))
					pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
					pm.PollQueue.AddPollingPoint(pp)
				} else {
					pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Point (%s) is not enabled./n", pnt.UUID))
				}
			}
		} else {
			pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Device (%s) is not enabled./n", dev.UUID))
		}
	}
	heap.Init(pm.PollQueue.PriorityQueue)
	if wasRunning {
		pm.StartQueueUnloader()
	}
	//TODO: START ANY OTHER REQUIRED QUEUE LOADERS/OPTIMIZERS
	//pm.PrintPollQueuePointUUIDs()
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

func (pm *NetworkPollManager) printPointDebugInfo(pnt *model.Point) {
	printString := "\n\n"
	if pnt != nil {
		printString += fmt.Sprint("Point: ", pnt.UUID, " ", pnt.Name, "\n")
		printString += fmt.Sprint("WriteMode: ", pnt.WriteMode, "\n")
		if pnt.WritePollRequired != nil {
			printString += fmt.Sprint("WritePollRequired: ", *pnt.WritePollRequired, "\n")
		}
		if pnt.ReadPollRequired != nil {
			printString += fmt.Sprint("ReadPollRequired: ", *pnt.ReadPollRequired, "\n")
		}
		if pnt.WriteValue == nil {
			printString += fmt.Sprint("WriteValue: nil", "\n")
		} else {
			printString += fmt.Sprint("WriteValue: ", *pnt.WriteValue, "\n")
		}
		if pnt.OriginalValue == nil {
			printString += fmt.Sprint("OriginalValue: nil", "\n")
		} else {
			printString += fmt.Sprint("OriginalValue: ", *pnt.OriginalValue, "\n")
		}
		if pnt.PresentValue == nil {
			printString += fmt.Sprint("PresentValue: nil", "\n")
		} else {
			printString += fmt.Sprint("PresentValue: ", *pnt.PresentValue, "\n")
		}
		if pnt.CurrentPriority == nil {
			printString += fmt.Sprint("CurrentPriority: nil", "\n")
		} else {
			printString += fmt.Sprint("CurrentPriority: ", *pnt.CurrentPriority, "\n")
		}
		if pnt.Priority != nil {
			if pnt.Priority.P1 != nil {
				printString += fmt.Sprint("_1: ", *pnt.Priority.P1, "\n")
			}
			if pnt.Priority.P2 != nil {
				printString += fmt.Sprint("_2: ", *pnt.Priority.P2, "\n")
			}
			if pnt.Priority.P3 != nil {
				printString += fmt.Sprint("_3: ", *pnt.Priority.P3, "\n")
			}
			if pnt.Priority.P4 != nil {
				printString += fmt.Sprint("_4: ", *pnt.Priority.P4, "\n")
			}
			if pnt.Priority.P5 != nil {
				printString += fmt.Sprint("_5: ", *pnt.Priority.P5, "\n")
			}
			if pnt.Priority.P6 != nil {
				printString += fmt.Sprint("_6: ", *pnt.Priority.P6, "\n")
			}
			if pnt.Priority.P7 != nil {
				printString += fmt.Sprint("_7: ", *pnt.Priority.P7, "\n")
			}
			if pnt.Priority.P8 != nil {
				printString += fmt.Sprint("_8: ", *pnt.Priority.P8, "\n")
			}
			if pnt.Priority.P9 != nil {
				printString += fmt.Sprint("_9: ", *pnt.Priority.P9, "\n")
			}
			if pnt.Priority.P10 != nil {
				printString += fmt.Sprint("_10: ", *pnt.Priority.P10, "\n")
			}
			if pnt.Priority.P11 != nil {
				printString += fmt.Sprint("_11: ", *pnt.Priority.P11, "\n")
			}
			if pnt.Priority.P12 != nil {
				printString += fmt.Sprint("_12: ", *pnt.Priority.P12, "\n")
			}
			if pnt.Priority.P13 != nil {
				printString += fmt.Sprint("_13: ", *pnt.Priority.P13, "\n")
			}
			if pnt.Priority.P14 != nil {
				printString += fmt.Sprint("_14: ", *pnt.Priority.P14, "\n")
			}
			if pnt.Priority.P15 != nil {
				printString += fmt.Sprint("_15: ", *pnt.Priority.P15, "\n")
			}
			if pnt.Priority.P16 != nil {
				printString += fmt.Sprint("_16: ", *pnt.Priority.P16, "\n")
			}
		}
		pm.pollQueueDebugMsg(printString)
		return
	}
	pm.pollQueueDebugMsg("ERROR: INVALID POINT")
}

func (pm *NetworkPollManager) printPollingPointDebugInfo(pp *PollingPoint) {
	if pp != nil {
		pm.pollQueueDebugMsg(fmt.Sprintf("ModbusPolling() pp %+v", pp))
	}
}

func (pm *NetworkPollManager) PollCompleteStatsUpdate(pp *PollingPoint, pollTimeSecs float64) {
	pm.pollQueueDebugMsg("PollCompleteStatsUpdate()")

	pm.AveragePollExecuteTimeSecs = ((pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) + pollTimeSecs) / (float64(pm.TotalPollCount) + 1)
	pm.TotalPollCount++
	enabledTime := time.Since(time.Unix(pm.PollingStartTimeUnix, 0)) * time.Second
	pm.BusyTime = (pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) / enabledTime.Seconds()

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

func (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate bool) {
	pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification Point UUID: %s, writeSuccess: %t, readSuccess: %t, pollTime: %f", pp.FFPointUUID, writeSuccess, readSuccess, pollTimeSecs))

	if !pointUpdate {
		pm.PollCompleteStatsUpdate(pp, pollTimeSecs) // This will update the relevant PollManager statistics.
	}

	point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
	if point == nil || err != nil {
		fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s /n", pp.FFPointUUID)
		return
	}
	//TODO: potentially only required on writeSuccess (but possibility of lockup on a bad point)
	//Reset poll priority to set value (in cases where pp has been escalated to ASAP).
	pp.PollPriority = point.PollPriority

	pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification: point %+v", point))
	//point.PrintPointValues()

	//If the device was deleted while this point was being polled, discard the PollingPoint
	if !pointUpdate && !pm.PollQueue.CheckIfActiveDevicesListIncludes(point.DeviceUUID) {
		return
	}

	switch point.WriteMode {
	case model.ReadOnce: //ReadOnce          If read_successful then don't re-add.
		point.WritePollRequired = boolean.NewFalse()
		if readSuccess {
			point.ReadPollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			point.ReadPollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)
		}

	case model.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.WritePollRequired = boolean.NewFalse()
		if readSuccess {
			point.ReadPollRequired = boolean.NewFalse()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			pm.pollQueueDebugMsg("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			pm.pollQueueDebugMsg("PollingPointCompleteNotification: NOT readSuccess")
			point.ReadPollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.  TODO: This might conflict with pausing polling on PortUnavailable
			pm.pollQueueDebugMsg("PollingPointCompleteNotification: ABOUT TO ADD POINT")
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}

	case model.WriteOnce: //WriteOnce         If write_successful then don't re-add.
		point.ReadPollRequired = boolean.NewFalse()
		if writeSuccess {
			point.WritePollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			point.WritePollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case model.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		if boolean.IsTrue(point.WritePollRequired) && writeSuccess {
			point.WritePollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (boolean.IsTrue(point.WritePollRequired) && !writeSuccess) {
			point.WritePollRequired = boolean.NewTrue()
			if pointUpdate {
				point.ReadPollRequired = boolean.NewTrue()
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
			break
		}
		if readSuccess {
			point.ReadPollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if boolean.IsTrue(point.ReadPollRequired) && !readSuccess {
			point.ReadPollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case model.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()
		if writeSuccess {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			pm.pollQueueDebugMsg("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case model.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		if boolean.IsTrue(point.WritePollRequired) && writeSuccess {
			point.WritePollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (boolean.IsTrue(point.WritePollRequired) && !writeSuccess) {
			point.WritePollRequired = boolean.NewTrue()
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
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
		}

	case model.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		pm.pollQueueDebugMsg(fmt.Sprintf("WriteAndMaintain point %+v\n", point))
		writeValuePointer := point.Priority.GetHighestPriorityValue()
		if writeValuePointer != nil {
			pm.pollQueueDebugMsg(fmt.Sprintf("WriteAndMaintain writeValuePointer %+v\n", *writeValuePointer))
			writeValue := *writeValuePointer
			noPV := true
			var presentValue float64
			if point.PresentValue != nil {
				noPV = false
				presentValue = *point.PresentValue
			}
			if noPV || presentValue != writeValue {
				if pp.RepollTimer != nil {
					pp.RepollTimer.Stop()
					pp.RepollTimer = nil
				}
				point.WritePollRequired = boolean.NewTrue()
				pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
				pm.PollQueue.AddPollingPoint(pp)                              //re-add to poll queue immediately
			} else {
				point.WritePollRequired = boolean.NewFalse()
				duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
				pm.pollQueueDebugMsg("duration: ", duration)
				// This line sets a timer to re-add the point to the poll queue after the PollRate time.
				pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
				addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
				if !addSuccess {
					pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
				}
			}
		}
	}

	pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification (ABOUT TO DB UPDATE): point  %+v", point))
	//point.PrintPointValues()
	//TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
	point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)
	//printPointDebugInfo(point)

}

func (pm *NetworkPollManager) MakePollingPointRepollCallback(pp *PollingPoint, writeMode model.WriteMode) func() {
	//log.Info("MakePollingPointRepollCallback()")
	f := func() {
		pm.pollQueueDebugMsg(fmt.Sprintf("CALL PollingPointRepollCallback func() pp: %+v", pp))
		pp.RepollTimer = nil
		_, removeSuccess := pm.PollQueue.StandbyPollingPoints.RemovePollingPointByPointUUID(pp.FFPointUUID)
		if !removeSuccess {
			pm.pollQueueErrorMsg(fmt.Sprintf("Modbus MakePollingPointRepollCallback(): polling point could not be found in StandbyPollingPoints.  (%s)", pp.FFPointUUID))
		}

		point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{})
		if point == nil || err != nil {
			pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s", pp.FFPointUUID))
			return
		}

		switch writeMode {
		case model.ReadOnce:
			return

		case model.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()

		case model.WriteOnce: //WriteOnce         If write_successful then don't re-add.
			return

		case model.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
			return

		case model.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
			point.ReadPollRequired = boolean.NewFalse()
			point.WritePollRequired = boolean.NewTrue()

		case model.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()

		case model.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()
		}

		//Now add the polling point back to the polling queue
		pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) //starts a countdown for queue lockup alerts.
		pm.PollQueue.AddPollingPoint(pp)

		//TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
		pm.pollQueueDebugMsg(fmt.Sprintf("pm.DBHandlerRef: %+v", pm.DBHandlerRef))
		point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)
		if err != nil || point == nil {
			pm.pollQueueErrorMsg(fmt.Sprintf("point DB UPDATE FAILED Err: %+v", err))
			return
		}
		pm.pollQueueDebugMsg(fmt.Sprintf("point after DB UPDATE: %+v", point))
		//printPointDebugInfo(point)
	}
	return f
}

func (pm *NetworkPollManager) MakeLockupTimerFunc(priority model.PollPriority) *time.Timer {
	timeoutDuration := 5 * time.Minute

	switch priority {
	case model.PRIORITY_ASAP:
		timeoutDuration = pm.ASAPPriorityMaxCycleTime

	case model.PRIORITY_HIGH:
		timeoutDuration = pm.HighPriorityMaxCycleTime

	case model.PRIORITY_NORMAL:
		timeoutDuration = pm.NormalPriorityMaxCycleTime

	case model.PRIORITY_LOW:
		timeoutDuration = pm.LowPriorityMaxCycleTime

	}

	f := func() {
		pm.pollQueueDebugMsg("Polling Lockout Timer Expired! Polling Priority: %d,  Polling Network: %s", priority, pm.FFNetworkUUID)
		switch priority {
		case model.PRIORITY_ASAP:
			pm.ASAPPriorityLockupAlert = true

		case model.PRIORITY_HIGH:
			pm.HighPriorityLockupAlert = true

		case model.PRIORITY_NORMAL:
			pm.NormalPriorityLockupAlert = true

		case model.PRIORITY_LOW:
			pm.LowPriorityLockupAlert = true

		}
	}

	return time.AfterFunc(timeoutDuration, f)
}

func (pm *NetworkPollManager) SetPointPollRequiredFlagsBasedOnWriteMode(point *model.Point) {
	pm.pollQueueDebugMsg("SetPointPollRequiredFlagsBasedOnWriteMode BEFORE: point")
	pm.pollQueueDebugMsg("%+v\n", point)
	pm.pollQueueDebugMsg("MODBUS SetPointPollRequiredFlagsBasedOnWriteMode(): PRIORITY")
	pm.pollQueueDebugMsg("%+v\n", point.Priority)

	if point == nil {
		pm.pollQueueDebugMsg("NetworkPollManager.SetPointPollRequiredFlagsBasedOnWriteMode(): couldn't find point %s /n", point.UUID)
		return
	}

	switch point.WriteMode {
	case model.ReadOnce:
		return

	case model.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewFalse()

	case model.WriteOnce: //WriteOnce         If write_successful then don't re-add.
		return

	case model.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		return

	case model.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()
	}

	pm.pollQueueDebugMsg("SetPointPollRequiredFlagsBasedOnWriteMode AFTER: point")
	pm.pollQueueDebugMsg("%+v\n", point)
	pm.pollQueueDebugMsg("MODBUS SetPointPollRequiredFlagsBasedOnWriteMode(): PRIORITY")
	pm.pollQueueDebugMsg("%+v\n", point.Priority)

	pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)

}
