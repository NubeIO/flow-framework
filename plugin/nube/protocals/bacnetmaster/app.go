package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/writemode"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
// addNetwork add network. Called via API call (or wizard)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body == nil {
		inst.bacnetErrorMsg("addNetwork(): nil network object")
		return nil, errors.New("empty network body, no network created")
	}
	if body.NetworkInterface == "" {
		return nil, errors.New("network interface can not be empty try, eth0")
	}

	inst.bacnetDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if network == nil || err != nil {
		inst.bacnetErrorMsg("addNetwork(): failed to create bacnet network: ", body.Name)
		return nil, errors.New("failed to create bacnet network")
	}
	err = inst.bacnetStoreNetwork(network)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("issue on add bacnet-device to store err:%s", err.Error()))
	}

	if boolean.IsTrue(network.Enable) {
		conf := inst.GetConfig().(*Config)
		pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.LogLevel}
		pollManager := NewPollManager(&pollQueueConfig, &inst.db, network.UUID, inst.pluginUUID, float.NonNil(network.MaxPollRate))
		pollManager.StartPolling()
		inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
	}
	return network, nil
}

// addDevice add device. Called via API call (or wizard)
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.bacnetDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.bacnetDebugMsg("addDevice(): failed to create bacnet device: ", body.Name)
		return nil, errors.New("failed to create bacnet device")
	}

	err = inst.bacnetStoreDevice(device)
	if err != nil {
		return nil, errors.New("issue on add bacnet-device to store")
	}

	inst.bacnetDebugMsg("addDevice(): ", body.UUID)
	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

// addPoint add point. Called via API call (or wizard)
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.bacnetDebugMsg("addPoint(): ", body.Name)

	if writemode.IsWriteable(body.WriteMode) {
		body.WritePollRequired = boolean.NewTrue()
	} else {
		body.WritePollRequired = boolean.NewFalse()
	}
	body.ReadPollRequired = boolean.NewTrue()

	isTypeBool := checkForBooleanType(body.ObjectType, body.DataType)
	body.IsTypeBool = nils.NewBool(isTypeBool)

	isOutput := checkForOutputType(body.ObjectType)
	body.IsOutput = nils.NewBool(isOutput)

	// point, err = inst.db.CreatePoint(body, true, false)
	point, err = inst.db.CreatePoint(body, true, true)
	if point == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): failed to create bacnet point: ", body.Name)
		return nil, errors.New("failed to create bacnet point")
	}
	inst.bacnetDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	// net, err := inst.db.DB.GetNetworkByDeviceUUID(point.DeviceUUID, api.Args{})
	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("addPoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}

	if boolean.IsTrue(point.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = point.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		// netPollMan.PollQueue.AddPollingPoint(pp)
		// netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	}
	return point, nil

}

// updateNetwork update network. Called via API call.
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.bacnetDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updateNetwork():  nil network object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = errors.New("network not enabled").Error()
		body.CommonFault.LastFail = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil || network == nil {
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(network.UUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updateNetwork(): cannot find NetworkPollManager for network: ", network.UUID)
		return
	}

	if boolean.IsTrue(network.Enable) == false && netPollMan.Enable == true {
		// DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()
	} else if boolean.IsTrue(network.Enable) == true && netPollMan.Enable == false {
		// DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
	}

	return network, nil
}

// updateDevice update device. Called via API call.
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.bacnetDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = errors.New("device not enabled").Error()
		body.CommonFault.LastFail = time.Now().UTC()
	}

	dev, err := inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil || dev == nil {
		return nil, err
	}

	err = inst.bacnetStoreDevice(device)
	if err != nil {
		return nil, err
	}

	if boolean.IsTrue(dev.Enable) == true { // If Enabled we need to GetDevice so we get Points
		dev, err = inst.db.GetDevice(dev.UUID, api.Args{WithPoints: true})
		if err != nil || dev == nil {
			return nil, err
		}
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updateDevice(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}

	if boolean.IsFalse(dev.Enable) && netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)

	} else if boolean.IsTrue(dev.Enable) && !netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		// DO POLLING ENABLE ACTIONS FOR DEVICE
		for _, pnt := range dev.Points {
			if boolean.IsTrue(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				//netPollMan.PollQueue.AddPollingPoint(pp)  //This is the original tested way, above is new so that on device update, it will re-poll write-once points
			}
		}

	} else if boolean.IsTrue(dev.Enable) {
		// TODO: Currently on every device update, all device points are removed, and re-added.
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)
		for _, pnt := range dev.Points {
			if boolean.IsTrue(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				//netPollMan.PollQueue.AddPollingPoint(pp)  //This is the original tested way, above is new so that on device update, it will re-poll write-once points
			}
		}
	}
	// TODO: NEED TO ACCOUNT FOR OTHER CHANGES ON DEVICE.  It would be useful to have a way to know if the device polling rates were changed.

	device, err = inst.db.UpdateDevice(dev.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// updatePoint update point. Called via API call.
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.bacnetDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updatePoint(): nil point object")
		return
	}

	/*
		pnt, err := inst.db.GetPoint(body.UUID, api.Args{WithPriority: true})
		if pnt == nil || err != nil {
			inst.bacnetErrorMsg("could not find pointID: ", pp.FFPointUUID)
			netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
			continue
		}

	*/

	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = errors.New("point not enabled").Error()
		body.CommonFault.LastFail = time.Now().UTC()
	}

	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updatePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(point, err)
		return
	}

	if boolean.IsTrue(point.Enable) && boolean.IsTrue(dev.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		// TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = point.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		// netPollMan.PollQueue.AddPollingPoint(pp)
		// netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
		// DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
	}

	return point, nil
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {

	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	inst.bacnetDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.bacnetDebugMsg("writePoint(): nil point object")
		return
	}

	inst.bacnetDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.bacnetDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	/* TODO: ONLY NEEDED IF THE WRITE VALUE IS WRITTEN ON COV (CURRENTLY IT IS WRITTEN ANYTIME THERE IS A WRITE COMMAND).
	point, err = inst.db.GetPoint(pntUUID, apinst.Args{})
	if err != nil || point == nil {
		inst.bacnetErrorMsg("writePoint(): bad response from GetPoint(), ", err)
		return nil, err
	}

	previousWriteVal := -1.11
	if isWriteable(point.WriteMode) {
		previousWriteVal = utils.Float64IsNil(point.WriteValue)
	}
	*/

	// body.WritePollRequired = utils.NewTrue() // TODO: commented out this section, seems like useless

	point, _, isWriteValueChange, _, err := inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	//  TODO: THIS SECTION MIGHT BE USEFUL IF WE ADD ASAP PRIORITY FOR IMMEDIATE POINT WRITES
	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("writePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(point, err)
		return
	}

	if boolean.IsTrue(point.Enable) {
		if isWriteValueChange { //if the write value has changed, we need to re-add the point so that it is polled asap (if required)
			pp, _ := netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
			if pp == nil {
				if netPollMan.PollQueue.OutstandingPollingPoints.GetPollingPointIndexByPointUUID(point.UUID) > -1 {
					if writemode.IsWriteable(point.WriteMode) {
						netPollMan.PollQueue.PointsUpdatedWhilePolling[point.UUID] = true // this triggers a write post at ASAP priority (for writeable points).
						point.WritePollRequired = boolean.NewTrue()
						if point.WriteMode != model.WriteAlways && point.WriteMode != model.WriteOnce {
							point.ReadPollRequired = boolean.NewTrue()
						} else {
							point.ReadPollRequired = boolean.NewFalse()
						}
					} else {
						netPollMan.PollQueue.PointsUpdatedWhilePolling[point.UUID] = false //
						point.WritePollRequired = boolean.NewFalse()
					}
					return point, nil
				} else {
					inst.bacnetDebugMsg("writePoint(): cannot find PollingPoint for point (could be out for polling: ", point.UUID)
					inst.pointUpdateErr(point, errors.New(fmt.Sprint("writePoint(): cannot find PollingPoint for point: ", point.UUID)))
					return point, err
				}
			}
			pp.PollPriority = model.PRIORITY_ASAP
			netPollMan.PollQueue.AddPollingPoint(pp)
			// netPollMan.PollQueue.UpdatePollingPointByPointUUID(point.UUID, model.PRIORITY_ASAP)

			/*
				netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
				//DO POLLING ENABLE ACTIONS FOR POINT
				//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
				pp := pollqueue.NewPollingPoint(body.UUID, body.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = body.PollPriority
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				//netPollMan.PollQueue.AddPollingPoint(pp)
				//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
			*/
		}
	} else {
		// DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(pntUUID)
	}
	return point, nil
}

// deleteNetwork delete network. Called via API call.
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.bacnetDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteNetwork(): nil network object")
		return
	}
	found := false
	for index, netPollMan := range inst.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == body.UUID {
			netPollMan.StopPolling()
			// Next remove the NetworkPollManager from the slice in polling instance
			inst.NetworkPollManagers[index] = inst.NetworkPollManagers[len(inst.NetworkPollManagers)-1]
			inst.NetworkPollManagers = inst.NetworkPollManagers[:len(inst.NetworkPollManagers)-1]
			found = true
		}
	}
	if !found {
		inst.bacnetDebugMsg("deleteNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	ok, err = inst.closeBacnetStoreNetwork(body.UUID)
	return ok, nil
}

// deleteNetwork delete device. Called via API call.
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.bacnetDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteDevice(): nil device object")
		return
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("deleteDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
		return
	}
	netPollMan.PollQueue.RemovePollingPointByDeviceUUID(body.UUID)
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// deletePoint delete point. Called via API call.
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.bacnetDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deletePoint(): nil point object")
		return
	}

	dev, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("addPoint(): bad response from GetDevice()")
		return false, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)

	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(body, err)
		return
	}

	netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
	otherPointsOnSameDeviceExist := netPollMan.PollQueue.CheckPollingQueueForDevUUID(body.DeviceUUID)
	if !otherPointsOnSameDeviceExist {
		netPollMan.PollQueue.RemoveDeviceFromActiveDevicesList(body.DeviceUUID)
	}
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// pointUpdate update point. Called from within plugin.
func (inst *Instance) pointUpdate(point *model.Point, value float64, writeSuccess, readSuccess, clearFaults bool) (*model.Point, error) {
	if clearFaults {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Info
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
		point.CommonFault.LastOk = time.Now().UTC()
	}

	if readSuccess {
		if value != float.NonNil(point.OriginalValue) {
			point.ValueUpdatedFlag = boolean.NewTrue() // Flag so that UpdatePointValue() will broadcast new value to producers. TODO: MAY NOT BE NEEDED.
		}
		point.OriginalValue = float.New(value)
	}
	point.InSync = boolean.NewTrue() // TODO: MAY NOT BE NEEDED.

	_, err := inst.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		inst.bacnetDebugMsg("MODBUS UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

// pointUpdateErr update point with errors. Called from within plugin.
func (inst *Instance) pointUpdateErr(point *model.Point, err error) (*model.Point, error) {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		inst.bacnetDebugMsg(" pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdateErr update point with errors. Called from within plugin.
func (inst *Instance) deviceUpdateErr(device *model.Device, err error) (*model.Device, error) {
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = model.MessageLevel.Fail
	device.CommonFault.MessageCode = model.CommonFaultCode.PointError
	device.CommonFault.Message = err.Error()
	device.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdateDevice(device.UUID, device, true)
	if err != nil {
		inst.bacnetDebugMsg(" deviceUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdateErr update point with errors. Called from within plugin.
func (inst *Instance) networkUpdateErr(network *model.Network, err error) (*model.Network, error) {
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = model.MessageLevel.Fail
	network.CommonFault.MessageCode = model.CommonFaultCode.PointError
	network.CommonFault.Message = err.Error()
	network.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdateNetwork(network.UUID, network, true)
	if err != nil {
		inst.bacnetDebugMsg(" networkUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

func (inst *Instance) getNetworks() ([]*model.Network, error) {
	return inst.db.GetNetworks(api.Args{})
}
