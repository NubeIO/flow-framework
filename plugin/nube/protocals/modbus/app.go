package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/config"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/pollqueue"
	"github.com/NubeIO/flow-framework/utils/array"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"go.bug.st/serial"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
// addNetwork add network. Called via API call (or wizard)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body == nil {
		inst.modbusErrorMsg("addNetwork(): nil network object")
		return nil, errors.New("empty network body, no network created")
	}
	inst.modbusDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if network == nil || err != nil {
		inst.modbusErrorMsg("addNetwork(): failed to create modbus network: ", body.Name)
		return nil, errors.New("failed to create modbus network")
	}

	if boolean.IsTrue(body.Enable) {
		conf := inst.GetConfig().(*config.Config)
		pollManager := pollqueue.NewPollManager(conf, &inst.db, network.UUID, inst.pluginUUID)
		pollManager.StartPolling()
		inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
	}
	return network, nil
}

// addDevice add device. Called via API call (or wizard)
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.modbusDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.modbusDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.modbusDebugMsg("addDevice(): failed to create modbus device: ", body.Name)
		return nil, errors.New("failed to create modbus device")
	}

	inst.modbusDebugMsg("addDevice(): ", body.UUID)

	if boolean.IsFalse(device.Enable) {
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

// addPoint add point. Called via API call (or wizard)
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.modbusDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.modbusDebugMsg("addPoint(): ", body.Name)

	if isWriteable(body.WriteMode) {
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
		inst.modbusDebugMsg("addPoint(): failed to create modbus point: ", body.Name)
		return nil, errors.New("failed to create modbus point")
	}
	inst.modbusDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("addPoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
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
	inst.modbusDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updateNetwork():  nil network object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.NetworkError
		body.CommonFault.Message = errors.New("network disabled").Error()
		body.CommonFault.LastFail = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil || network == nil {
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(network.UUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("updateNetwork(): cannot find NetworkPollManager for network: ", network.UUID)
		return
	}

	if boolean.IsFalse(network.Enable) && netPollMan.Enable == true {
		// DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	} else if boolean.IsTrue(network.Enable) && netPollMan.Enable == false {
		// DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
		network.CommonFault.InFault = false
		network.CommonFault.MessageLevel = model.MessageLevel.Info
		network.CommonFault.MessageCode = model.CommonFaultCode.Ok
		network.CommonFault.Message = errors.New("").Error()
		network.CommonFault.LastOk = time.Now().UTC()
		inst.db.ClearErrorsForAllDevicesOnNetwork(network.UUID, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network, true)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

// updateDevice update device. Called via API call.
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.modbusDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
		body.CommonFault.Message = "device disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}

	dev, err := inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil || dev == nil {
		return nil, err
	}

	if boolean.IsTrue(dev.Enable) { // If Enabled we need to GetDevice so we get Points
		dev, err = inst.db.GetDevice(dev.UUID, api.Args{WithPoints: true})
		if err != nil || dev == nil {
			return nil, err
		}
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("updateDevice(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}
	if boolean.IsFalse(dev.Enable) && netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		inst.db.SetErrorsForAllPointsOnDevice(dev.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)

	} else if boolean.IsTrue(dev.Enable) && !netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		// DO POLLING ENABLE ACTIONS FOR DEVICE
		dev.CommonFault.InFault = false
		dev.CommonFault.MessageLevel = model.MessageLevel.Info
		dev.CommonFault.MessageCode = model.CommonFaultCode.Ok
		dev.CommonFault.Message = ""
		dev.CommonFault.LastOk = time.Now().UTC()
		err = inst.db.ClearErrorsForAllPointsOnDevice(dev.UUID)
		if err != nil {
			inst.modbusDebugMsg("updateDevice(): error on ClearErrorsForAllPointsOnDevice(): ", err)
		}
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
		dev.CommonFault.InFault = false
		dev.CommonFault.MessageLevel = model.MessageLevel.Info
		dev.CommonFault.MessageCode = model.CommonFaultCode.Ok
		dev.CommonFault.Message = ""
		dev.CommonFault.LastOk = time.Now().UTC()
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

	device, err = inst.db.UpdateDevice(dev.UUID, dev, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// updatePoint update point. Called via API call.
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.modbusDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updatePoint(): nil point object")
		return
	}

	/*
		pnt, err := inst.db.GetPoint(body.UUID, api.Args{WithPriority: true})
		if pnt == nil || err != nil {
			inst.modbusErrorMsg("could not find pointID: ", pp.FFPointUUID)
			netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
			continue
		}

	*/

	inst.modbusDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.modbusDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = errors.New("point disabled").Error()
		body.CommonFault.LastFail = time.Now().UTC()
	}

	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil || point == nil {
		inst.modbusDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("updatePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(point, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
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

	inst.modbusDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.modbusDebugMsg("writePoint(): nil point object")
		return
	}

	inst.modbusDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.modbusDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	/* TODO: ONLY NEEDED IF THE WRITE VALUE IS WRITTEN ON COV (CURRENTLY IT IS WRITTEN ANYTIME THERE IS A WRITE COMMAND).
	point, err = inst.db.GetPoint(pntUUID, apinst.Args{})
	if err != nil || point == nil {
		inst.modbusErrorMsg("writePoint(): bad response from GetPoint(), ", err)
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
		inst.modbusDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	//  TODO: THIS SECTION MIGHT BE USEFUL IF WE ADD ASAP PRIORITY FOR IMMEDIATE POINT WRITES
	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("writePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(point, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}

	if boolean.IsTrue(point.Enable) {
		if isWriteValueChange { //if the write value has changed, we need to re-add the point so that it is polled asap (if required)
			pp, _ := netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
			if pp == nil {
				if netPollMan.PollQueue.OutstandingPollingPoints.GetPollingPointIndexByPointUUID(point.UUID) > -1 {
					if isWriteable(point.WriteMode) {
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
					point, err = inst.db.UpdatePoint(point.UUID, point, true)
					if err != nil || point == nil {
						inst.modbusDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
						return nil, err
					}
					return point, nil
				} else {
					inst.modbusDebugMsg("writePoint(): cannot find PollingPoint for point (could be out for polling: ", point.UUID)
					inst.pointUpdateErr(point, "writePoint(): cannot find PollingPoint for point: ", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
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
	inst.modbusDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deleteNetwork(): nil network object")
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
		inst.modbusDebugMsg("deleteNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// deleteNetwork delete device. Called via API call.
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.modbusDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deleteDevice(): nil device object")
		return
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("deleteDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
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
	inst.modbusDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deletePoint(): nil point object")
		return
	}

	dev, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("addPoint(): bad response from GetDevice()")
		return false, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)

	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		inst.pointUpdateErr(body, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
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

	_, err = inst.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		inst.modbusDebugMsg("MODBUS UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

// pointUpdateErr update point with errors. Called from within plugin.
func (inst *Instance) pointUpdateErr(point *model.Point, message string, messageLevel string, messageCode string) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = messageLevel
	point.CommonFault.MessageCode = messageCode
	point.CommonFault.Message = message
	point.CommonFault.LastFail = time.Now().UTC()
	err = inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.modbusDebugMsg(" pointUpdateErr()", err)
		return err
	}
	return nil
}

// deviceUpdateErr update device with errors. Called from within plugin.
func (inst *Instance) deviceUpdateErr(device *model.Device, err error) error {
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = model.MessageLevel.Warning
	device.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
	device.CommonFault.Message = err.Error()
	device.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdateDevice(device.UUID, device, true)
	if err != nil {
		inst.modbusDebugMsg(" deviceUpdateErr()", err)
		return err
	}
	return nil
}

// networkUpdateErr update network with errors. Called from within plugin.
func (inst *Instance) networkUpdateErr(network *model.Network, err error) error {
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = model.MessageLevel.Fail
	network.CommonFault.MessageCode = model.CommonFaultCode.PointError
	network.CommonFault.Message = err.Error()
	network.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdateNetwork(network.UUID, network, true)
	if err != nil {
		inst.modbusDebugMsg(" networkUpdateErr()", err)
		return err
	}
	return nil
}

// listSerialPorts list all serial ports on host
func (inst *Instance) listSerialPorts() (*array.Array, error) {
	ports, err := serial.GetPortsList()
	p := array.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
