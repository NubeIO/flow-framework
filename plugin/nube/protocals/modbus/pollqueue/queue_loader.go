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
	// TODO: STOP ANY OTHER QUEUE LOADERS
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
	for _, dev := range devs { // DEVICES
		if dev.NetworkUUID == pm.FFNetworkUUID && boolean.IsTrue(dev.Enable) {
			for _, pnt := range dev.Points { // POINTS
				if pnt.DeviceUUID == dev.UUID && boolean.IsTrue(pnt.Enable) {
					pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
					pp.PollPriority = pnt.PollPriority
					pm.pollQueueDebugMsg(fmt.Sprintf("RebuildPollingQueue() pp: %+v", pp))
					pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
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
	// TODO: START ANY OTHER REQUIRED QUEUE LOADERS/OPTIMIZERS
	// pm.PrintPollQueuePointUUIDs()
	return nil
}

func (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate bool) {
	pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification Point UUID: %s, writeSuccess: %t, readSuccess: %t, pollTime: %f", pp.FFPointUUID, writeSuccess, readSuccess, pollTimeSecs))

	_, success := pm.PollQueue.OutstandingPollingPoints.RemovePollingPointByPointUUID(pp.FFPointUUID)
	if !success {
		pm.pollQueueErrorMsg("NetworkPollManager.PollingPointCompleteNotification(): couldn't find polling point in OutstandingPollingPoints.  %s /n", pp.FFPointUUID)
	}

	if !pointUpdate {
		pm.PollCompleteStatsUpdate(pp, pollTimeSecs) // This will update the relevant PollManager statistics.
	}

	point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
	if point == nil || err != nil {
		pm.pollQueueErrorMsg("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s /n", pp.FFPointUUID)
		return
	}
	// TODO: potentially only required on writeSuccess (but possibility of lockup on a bad point)
	// Reset poll priority to set value (in cases where pp has been escalated to ASAP).
	pp.PollPriority = point.PollPriority

	val, ok := pm.PollQueue.PointsUpdatedWhilePolling[point.UUID]
	if ok {
		delete(pm.PollQueue.PointsUpdatedWhilePolling, point.UUID)
		if val == true { // point needs an ASAP write
			pp.PollPriority = model.PRIORITY_ASAP
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}
	}

	//pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification: point %+v", point))
	//pm.PrintPointDebugInfo(point)

	// If the device was deleted while this point was being polled, discard the PollingPoint
	if !pointUpdate && !pm.PollQueue.CheckIfActiveDevicesListIncludes(point.DeviceUUID) {
		return
	}

	switch point.WriteMode {
	case model.ReadOnce: // ReadOnce          If read_successful then don't re-add.
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)
		}

	case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.WritePollRequired = boolean.NewFalse()
		if readSuccess {
			point.ReadPollRequired = boolean.NewFalse()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			//pm.pollQueueDebugMsg("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			pm.pollQueueDebugMsg("PollingPointCompleteNotification: NOT readSuccess")
			point.ReadPollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.  TODO: This might conflict with pausing polling on PortUnavailable
			pm.pollQueueDebugMsg("PollingPointCompleteNotification: ABOUT TO ADD POINT")
			pm.PollQueue.AddPollingPoint(pp) // re-add to poll queue immediately
		}

	case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()
		if writeSuccess {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			//pm.pollQueueDebugMsg("duration: ", duration)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			break
		}
		if readSuccess {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// log.Info("duration: ", duration)
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
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		//pm.pollQueueDebugMsg(fmt.Sprintf("WriteAndMaintain point %+v\n", point))
		if point.WriteValue != nil {
			//pm.pollQueueDebugMsg(fmt.Sprintf("WriteAndMaintain WriteValue %+v\n", float.NonNil(point.WriteValue)))
			noPV := true
			var presentValue float64
			if point.PresentValue != nil {
				noPV = false
				presentValue = *point.PresentValue
			}
			if noPV || presentValue != *point.WriteValue {
				if pp.RepollTimer != nil {
					pp.RepollTimer.Stop()
					pp.RepollTimer = nil
				}
				point.WritePollRequired = boolean.NewTrue()
				pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
				pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			} else {
				point.WritePollRequired = boolean.NewFalse()
				duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
				//pm.pollQueueDebugMsg("duration: ", duration)
				// This line sets a timer to re-add the point to the poll queue after the PollRate time.
				pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
				addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
				if !addSuccess {
					pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
				}
			}
		} else {
			//If WriteValue is nil we still need to re-add the point to perform a read
			point.WritePollRequired = boolean.NewFalse()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}
	}

	//pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification (ABOUT TO DB UPDATE): point  %+v", point))
	// point.PrintPointValues()
	// TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
	point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)
	// printPointDebugInfo(point)

}

func (pm *NetworkPollManager) MakePollingPointRepollCallback(pp *PollingPoint, writeMode model.WriteMode) func() {
	// log.Info("MakePollingPointRepollCallback()")
	f := func() {
		//pm.pollQueueDebugMsg(fmt.Sprintf("CALL PollingPointRepollCallback func() pp: %+v", pp))
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

		case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()

		case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.
			return

		case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.
			return

		case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
			point.ReadPollRequired = boolean.NewFalse()
			point.WritePollRequired = boolean.NewTrue()

		case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()

		case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
			point.ReadPollRequired = boolean.NewTrue()
			point.WritePollRequired = boolean.NewFalse()
		}

		// Now add the polling point back to the polling queue
		pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
		pm.PollQueue.AddPollingPoint(pp)

		// TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
		//pm.pollQueueDebugMsg(fmt.Sprintf("pm.DBHandlerRef: %+v", pm.DBHandlerRef))
		point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)
		if err != nil || point == nil {
			pm.pollQueueErrorMsg(fmt.Sprintf("point DB UPDATE FAILED Err: %+v", err))
			return
		}
		//pm.pollQueueDebugMsg(fmt.Sprintf("point after DB UPDATE: %+v", point))
		// printPointDebugInfo(point)
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
		plugin, err := pm.DBHandlerRef.DB.GetPlugin(pm.FFPluginUUID)
		switch priority {
		case model.PRIORITY_ASAP:
			pm.ASAPPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: ASAP Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: ASAP Priority Poll Queue LOCKUP")
			}
			//TODO: update network error to show queue lockup

		case model.PRIORITY_HIGH:
			pm.HighPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: HIGH Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: HIGH Priority Poll Queue LOCKUP")
			}
			//TODO: update network error to show queue lockup

		case model.PRIORITY_NORMAL:
			pm.NormalPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: NORMAL Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: NORMAL Priority Poll Queue LOCKUP")
			}
			//TODO: update network error to show queue lockup

		case model.PRIORITY_LOW:
			pm.LowPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: LOW Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: LOW Priority Poll Queue LOCKUP")
			}
			//TODO: update network error to show queue lockup

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

	case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewFalse()

	case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.
		return

	case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		return

	case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()
	}

	pm.pollQueueDebugMsg("SetPointPollRequiredFlagsBasedOnWriteMode AFTER: point")
	pm.pollQueueDebugMsg("%+v\n", point)
	pm.pollQueueDebugMsg("MODBUS SetPointPollRequiredFlagsBasedOnWriteMode(): PRIORITY")
	pm.pollQueueDebugMsg("%+v\n", point.Priority)

	pm.DBHandlerRef.UpdatePoint(point.UUID, point, true)

}
