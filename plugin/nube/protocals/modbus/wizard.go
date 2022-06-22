package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"strconv"
	"time"
)

// wizard make a network/dev/pnt
func (inst *Instance) wizardTCP(body wizard) (string, error) {
	ip := "192.168.15.202"
	if body.IP != "" {
		ip = body.IP
	}
	p := 502
	if body.Port != 0 {
		p = body.Port
	}
	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}
	wizVersion := 0
	if body.WizardVersion != 0 {
		wizVersion = int(body.WizardVersion)
	}

	switch wizVersion {
	case 0:
		var net model.Network
		net.Name = "modbus"
		net.TransportType = model.TransType.IP
		net.PluginPath = "modbus"

		var dev model.Device
		dev.Name = "modbus"
		dev.CommonIP.Host = ip
		dev.CommonIP.Port = p
		dev.AddressId = da
		dev.ZeroMode = boolean.NewTrue()
		dev.PollDelayPointsMS = 5000

		var pnt model.Point
		pnt.Name = "modbus"
		pnt.Description = "modbus"
		pnt.AddressID = integer.New(1) // TODO check conversion
		pnt.ObjectType = string(model.ObjTypeWriteFloat32)

		_, err = inst.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
		if err != nil {
			return "modbus wizard 0 error: on flow-framework add modbus TCP network wizard", err
		}

		return "modbus wizard 0: added network, device, and point", err

	case 1:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.IP
		net.PluginPath = "modbus"

		net.PluginConfId = inst.pluginUUID
		_, err := inst.db.CreateNetwork(&net, false)
		if err != nil {
			inst.modbusErrorMsg(fmt.Sprintf("network creation failure: %s", err))
		}
		inst.modbusDebugMsg("Created a Network")

		for j := 1; j < 4; j++ {
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = boolean.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			_, err := inst.db.CreateDevice(&dev)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("device creation failure: %s", err))
			}
			inst.modbusDebugMsg("Created a Device: ", dev)

			var pnt model.Point
			pnt.Name = "Modbus Pnt " + strconv.Itoa(j)
			pnt.Description = "modbus"
			pnt.AddressID = integer.New(j) // TODO check conversion
			pnt.ObjectType = string(model.ObjTypeWriteFloat32)
			pnt.DeviceUUID = dev.UUID
			pnt.PollPriority = model.PRIORITY_NORMAL
			pnt.PollRate = model.RATE_NORMAL
			pnt.WriteMode = model.ReadOnly
			_, err = inst.db.CreatePoint(&pnt, false, true)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("consumer point creation failure: %s", err))
			}
			inst.modbusDebugMsg("Created a Point for Consumer", pnt)

		}
		return "modbus wizard 1: added networks, devices, and points", err

	case 2:
		for j := 1; j < 4; j++ {
			var net model.Network
			net.Name = "Modbus Net " + strconv.Itoa(j)
			net.TransportType = model.TransType.IP
			net.PluginPath = "modbus"
			time.Sleep(2 * time.Second)

			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = boolean.NewTrue()
			dev.PollDelayPointsMS = 5000

			var pnt model.Point
			pnt.Name = "Modbus Pnt " + strconv.Itoa(j)
			pnt.Description = "modbus"
			pnt.AddressID = integer.New(j) // TODO check conversion
			pnt.ObjectType = string(model.ObjTypeWriteFloat32)

			_, err = inst.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
			if err != nil {
				return "modbus wizard 1: on flow-framework add modbus TCP network wizard", err
			}
		}
		return "modbus wizard 1: added networks, devices, and points", err

	case 3:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.Serial
		net.PluginPath = "modbus"

		net.PluginConfId = inst.pluginUUID
		_, err := inst.db.CreateNetwork(&net, false)
		if err != nil {
			inst.modbusErrorMsg(fmt.Sprintf("network creation failure: %s", err))
		}
		inst.modbusDebugMsg("Created a Network")

		for j := 1; j < 4; j++ {
			time.Sleep(2 * time.Second)
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = boolean.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			dev.FastPollRate = float.New(5.0)
			dev.NormalPollRate = float.New(30.0)
			dev.SlowPollRate = float.New(120.0)
			_, err = inst.db.CreateDevice(&dev)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("device creation failure: %s", err))
			}
			inst.modbusDebugMsg("Created a Device: ", dev)
			for l := 1; l < 6; l++ {
				var pnt model.Point
				pnt.Name = "Modbus Pnt " + strconv.Itoa(l)
				pnt.Description = "modbus"
				pnt.AddressID = integer.New(l) // TODO check conversion
				pnt.ObjectType = string(model.ObjTypeWriteFloat32)
				pnt.DeviceUUID = dev.UUID
				pnt.PollPriority = model.PRIORITY_NORMAL
				pnt.PollRate = model.RATE_NORMAL
				if l == 1 {
					pnt.PollPriority = model.PRIORITY_LOW
				} else if l == 3 {
					pnt.PollPriority = model.PRIORITY_HIGH
				}
				pnt.WriteMode = model.ReadOnly
				_, err = inst.db.CreatePoint(&pnt, false, true)
				if err != nil {
					inst.modbusErrorMsg(fmt.Sprintf("consumer point creation failure: %s", err))
				}
				inst.modbusDebugMsg("Created a Point for Consumer", pnt)
			}
		}
		return "modbus wizard 3: added a network, 3 devices, and 3 points per device", err

	case 4:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.Serial
		net.PluginPath = "modbus"
		net.MaxPollRate = float.New(0.1)

		net.PluginConfId = inst.pluginUUID
		_, err := inst.db.CreateNetwork(&net, false)
		if err != nil {
			inst.modbusErrorMsg(fmt.Sprintf("network creation failure: %s", err))
		}
		inst.modbusDebugMsg("Created a Network")

		for j := 1; j < 2; j++ {
			time.Sleep(2 * time.Second)
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = boolean.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			dev.FastPollRate = float.New(5.0)
			dev.NormalPollRate = float.New(30.0)
			dev.SlowPollRate = float.New(120.0)
			_, err = inst.db.CreateDevice(&dev)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("device creation failure: %s", err))
			}
			inst.modbusDebugMsg("Created a Device: ", dev)
			// pointsArray := [4]int{401, 403, 405, 407}
			pointsArray := [1]int{401}
			for _, l := range pointsArray {
				var pnt model.Point
				pnt.Name = "Modbus Pnt " + strconv.Itoa(l)
				pnt.Description = "modbus"
				pnt.AddressID = integer.New(l) // TODO check conversion
				pnt.ObjectType = string(model.ObjTypeWriteHolding)
				pnt.DataType = string(model.TypeFloat32)
				pnt.DeviceUUID = dev.UUID
				pnt.PollPriority = model.PRIORITY_NORMAL
				pnt.PollRate = model.RATE_NORMAL
				pnt.WriteMode = model.ReadOnly
				pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
				_, err = inst.db.CreatePoint(&pnt, false, true)
				if err != nil {
					inst.modbusErrorMsg(fmt.Sprintf("consumer point creation failure: %s", err))
				}
				inst.modbusDebugMsg("Created a Point for Consumer", pnt)
			}
		}
		return "modbus wizard 4: added a network, 1 device, and 4 points", err

	case 5:
		if body.NameArg != "" && body.AddArg > 0 {
			networkName := "CliniMix-TMV"
			net, err := inst.db.GetNetworkByName(networkName, api.Args{})
			if err != nil || net == nil {
				if net == nil {
					net = &model.Network{}
				}
				net.Name = "CliniMix-TMV"
				net.TransportType = model.TransType.Serial
				net.PluginPath = "modbus"
				net.SerialPort = nstring.New("/dev/ttyXBEE-2")
				net.MaxPollRate = float.New(5)
				net.PluginConfId = inst.pluginUUID
				net, err = inst.addNetwork(net)
				if err != nil {
					inst.modbusErrorMsg(fmt.Sprintf("network creation failure: %s", err))
					inst.modbusDebugMsg("Created a Network")
				}
			} else {
				inst.modbusDebugMsg("Network already exists")
			}

			dev, exists := inst.db.DeviceNameExistsInNetwork(body.NameArg, net.UUID)
			if err != nil || dev == nil || !exists {
				if dev == nil {
					dev = &model.Device{}
				}
				dev.Name = body.NameArg
				dev.CommonIP.Host = "0.0.0.0"
				dev.CommonIP.Port = p
				dev.AddressId = int(body.AddArg)
				dev.ZeroMode = boolean.NewTrue()
				dev.PollDelayPointsMS = 1000
				dev.NetworkUUID = net.UUID
				dev.FastPollRate = float.New(5.0)
				dev.NormalPollRate = float.New(30.0)
				dev.SlowPollRate = float.New(120.0)
				_, err = inst.addDevice(dev)
				if err != nil {
					inst.modbusErrorMsg(fmt.Sprintf("device creation failure: %s", err))
				}
				inst.modbusDebugMsg("Created a Device: ", dev)

			} else {
				inst.modbusDebugMsg("Device already exists")
			}

			type tmvPoint struct {
				AddressID              int
				Name                   string
				Description            string
				ObjectType             model.ObjectType
				DataType               model.DataType
				WriteMode              model.WriteMode
				PollPriority           model.PollPriority
				PollRate               model.PollRate
				PointPriorityArrayMode model.PointPriorityArrayMode
				Fallback               float64
			}

			pointsArray := [26]tmvPoint{
				{
					AddressID:              1001,
					Name:                   "Enable",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               1,
				},
				{
					AddressID:              1002,
					Name:                   "Reset",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               0,
				},
				{
					AddressID:              1003,
					Name:                   "Solenoid-Allow",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               0,
				},
				{
					AddressID:              1004,
					Name:                   "App-Fault",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_LOW,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1006,
					Name:                   "Over-Temperature-Warn",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1007,
					Name:                   "Over-Temperature-Alert",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1008,
					Name:                   "One-Day-Low-Flow-Alert",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1009,
					Name:                   "Five-Day-Low-Flow-Alert",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1010,
					Name:                   "Monthly-Hot-Flush-Status",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1011,
					Name:                   "Solenoid-Status",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadCoil,
					DataType:               model.TypeDigital,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_NORMAL,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1001,
					Name:                   "Temperature-Setpoint",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               40,
				},
				{
					AddressID:              1003,
					Name:                   "Over-Temperature-Offset",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               4,
				},
				{
					AddressID:              1005,
					Name:                   "Low-Flow-Threshold-Setpoint",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               120,
				},
				{
					AddressID:              1007,
					Name:                   "Hot-Flush-Setpoint",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               60,
				},
				{
					AddressID:              1009,
					Name:                   "Hot-Flush-Delay",
					Description:            "modbus",
					ObjectType:             model.ObjTypeWriteHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               300,
				},
				{
					AddressID:              1017,
					Name:                   "Daily-Temperature-Test-1",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               -1,
				},
				{
					AddressID:              1019,
					Name:                   "Daily-Temperature-Test-2",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               -1,
				},
				{
					AddressID:              1021,
					Name:                   "Daily-Temperature-Test-3",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               -1,
				},
				{
					AddressID:              1023,
					Name:                   "Monthly-Mean-Temperature-Test",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               -1,
				},
				{
					AddressID:              1025,
					Name:                   "Total-Flow-Accumulation",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1027,
					Name:                   "One-Day-Flow-Accumulation",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_HIGH,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1029,
					Name:                   "Days-Of-Low-Flow",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1031,
					Name:                   "Over-Temperature-Warn-Count",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              1033,
					Name:                   "Over-Temperature-Alert-Count",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeFloat32,
					WriteMode:              model.ReadOnly,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
					Fallback:               0,
				},
				{
					AddressID:              10008,
					Name:                   "Unix-Timestamp",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeInt32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               0,
				},
				{
					AddressID:              10010,
					Name:                   "Timezone-Offset-Seconds",
					Description:            "modbus",
					ObjectType:             model.ObjTypeReadHolding,
					DataType:               model.TypeInt32,
					WriteMode:              model.WriteAndMaintain,
					PollPriority:           model.PRIORITY_NORMAL,
					PollRate:               model.RATE_SLOW,
					PointPriorityArrayMode: model.PriorityArrayToWriteValue,
					Fallback:               39600,
				},
			}

			for _, point := range pointsArray {
				time.Sleep(1 * time.Second)
				pnt := &model.Point{}
				pnt.Name = point.Name
				pnt.Description = point.Description
				pnt.AddressID = integer.New(point.AddressID)
				pnt.ObjectType = string(point.ObjectType)
				pnt.DataType = string(point.DataType)
				pnt.DeviceUUID = dev.UUID
				pnt.PollPriority = point.PollPriority
				pnt.PollRate = point.PollRate
				pnt.WriteMode = point.WriteMode
				pnt.Fallback = float.New(point.Fallback)
				pnt.PointPriorityArrayMode = point.PointPriorityArrayMode
				_, err = inst.addPoint(pnt)
				if err != nil {
					inst.modbusErrorMsg(fmt.Sprintf("point creation failure: %s", err))
				}
				inst.modbusDebugMsg("Created a Point for Consumer", pnt)
			}
			return "modbus wizard 5: added TMV Points", err
		}
		return "modbus wizard 5: no device name specified in 'arg1'", err
	}
	return "modbus wizard error: unknown wizard version", err
}

// wizard make a network/dev/pnt
func (inst *Instance) wizardSerial(body wizard) (string, error) {

	sp := "/dev/ttyUSB0"
	if body.SerialPort != "" {
		sp = body.SerialPort
	}
	br := 9600
	if body.BaudRate != 0 {
		br = int(body.BaudRate)
	}

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.Serial
	net.PluginPath = "modbus"
	net.SerialPort = &sp
	net.SerialBaudRate = integer.NewUint(uint(br))
	net.SerialParity = nstring.New("none")
	net.SerialDataBits = integer.NewUint(1)
	net.SerialStopBits = integer.NewUint(1)

	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}

	var dev model.Device
	dev.Name = "modbus"
	dev.AddressId = da
	dev.ZeroMode = boolean.NewTrue()
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressID = integer.New(1) // TODO check conversion
	pnt.ObjectType = string(model.ObjTypeWriteCoil)

	pntRet, err := inst.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus serial network wizard", err
	}

	inst.modbusErrorMsg(pntRet, err)
	return "pass: added network and points", err
}
