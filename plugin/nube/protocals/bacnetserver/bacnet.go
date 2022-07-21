package main

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/segmentation"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPort = 47809
)

func (inst *Instance) bacnetNetworkInit() {

	networks, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{WithDevices: true})
	if err != nil {
		return
	}
	for _, net := range networks {
		err := inst.bacnetNetwork(net)
		if err != nil {
			log.Errorln("bacnet-server init network error:", err)
			continue
		}
		for _, dev := range net.Devices {
			err := inst.bacnetDevice(dev)
			if err != nil {
				log.Errorln("bacnet-server init device error:", err)
				continue
			}
		}
	}

}

func (inst *Instance) initBacStore() {
	inst.BacStore = network.NewStore()
}

// bacnetNetwork add or update an instance a bacnet network that is cached in bacnet lib
func (inst *Instance) bacnetNetwork(net *model.Network) error {
	bacnetNet := &network.Network{
		Interface: net.NetworkInterface,
		Port:      integer.NonNil(net.Port),
		StoreID:   net.UUID,
	}
	return inst.BacStore.UpdateNetwork(net.UUID, bacnetNet)
}

// getBacnetNetwork get an instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) getBacnetNetwork(networkUUID string) (*network.Network, error) {
	return inst.BacStore.GetNetwork(networkUUID)
}

// closeBacnetNetwork delete the instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) closeBacnetNetwork(networkUUID string) (bool, error) {
	net, err := inst.BacStore.GetNetwork(networkUUID)
	if err != nil {
		return false, err
	}
	net.NetworkClose()
	return true, nil
}

// getBacnetDevice get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) getBacnetDevice(deviceUUID string) (*network.Device, error) {
	return inst.BacStore.GetDevice(deviceUUID)
}

// bacnetDevice add or update an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) bacnetDevice(dev *model.Device) error {
	max := intToUint32(integer.NonNil(dev.MaxADPU))
	seg := uint32(setSegmentation(dev.Segmentation))
	d := &network.Device{
		Ip:            dev.CommonIP.Host,
		Port:          dev.CommonIP.Port,
		DeviceID:      integer.NonNil(dev.DeviceObjectId),
		StoreID:       dev.UUID,
		NetworkNumber: integer.NonNil(dev.NetworkNumber),
		MacMSTP:       integer.NonNil(dev.DeviceMac),
		MaxApdu:       max,
		Segmentation:  seg,
	}

	net, _ := inst.getBacnetNetwork(dev.NetworkUUID)
	return inst.BacStore.UpdateDevice(dev.UUID, net, d)
}

// getDev get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) doReadValue(pnt *model.Point, networkUUID, deviceUUID string) (float64, error) {
	object, _, isBool := setObjectType(pnt.ObjectType)
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteValue:       nil,
		WriteNull:        false,
		WritePriority:    0,
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	// get network
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return 0, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return 0, err
	}
	var outValue float64
	if isBool {
		readBool, err := dev.PointReadBool(bp)
		if err != nil {
			log.Errorln("bacnet-server-read-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = Unit32ToFloat64(readBool)

	} else {
		readFloat32, err := dev.PointReadFloat32(bp)
		if err != nil {
			log.Errorln("bacnet-server-read-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = float32ToFloat64(readFloat32)
	}
	log.Infoln("bacnet-server-POINT-READ:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", outValue)
	return outValue, nil
}

func (inst *Instance) doWrite(pnt *model.Point, networkUUID, deviceUUID string) error {
	val := float.NonNil(pnt.WriteValue)
	object, isWrite, isBool := setObjectType(pnt.ObjectType)
	writePriority := integer.NonNil(pnt.WritePriority)
	if writePriority == 0 {
		writePriority = 16
	}
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteNull:        false,
		WritePriority:    uint8(writePriority),
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return err
	}
	if isWrite {
		if isBool {
			err = dev.PointWriteBool(bp, float64ToUint32(val))
			if err != nil {
				log.Errorln("bacnet-server-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		} else {
			err = dev.PointWriteAnalogue(bp, float64ToFloat32(val))
			if err != nil {
				log.Errorln("bacnet-server-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		}
	}
	log.Infoln("bacnet-server-POINT-WRITE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
	return nil

}

func setObjectType(object string) (obj btypes.ObjectType, isWritable, isBool bool) {
	switch object {
	case "analog_input":
		return btypes.AnalogInput, false, false
	case "analog_output":
		return btypes.AnalogOutput, true, false
	case "analog_value":
		return btypes.AnalogValue, true, false
	case "binary_input":
		return btypes.BinaryInput, false, true
	case "binary_output":
		return btypes.BinaryOutput, true, true
	case "binary_value":
		return btypes.BinaryValue, true, true
	case "multi_state_input":
		return btypes.MultiStateInput, false, false
	case "multi_state_output":
		return btypes.MultiStateOutput, true, false
	case "multi_state_value":
		return btypes.MultiStateValue, true, false
	default:
		return btypes.AnalogInput, false, false
	}
}

type SegmentedType string

const (
	SegmentedBoth     SegmentedType = "segmentation_both"
	SegmentedTransmit SegmentedType = "segmentation_transmit"
	SegmentedReceive  SegmentedType = "segmentation_receive"
	NoSegmentation    SegmentedType = "no_segmentation"
)

func setSegmentation(SegmentedType string) (out segmentation.SegmentedType) {
	switch SegmentedType {
	case string(SegmentedBoth):
		return segmentation.SegmentedBoth
	case string(SegmentedTransmit):
		return segmentation.SegmentedTransmit
	case string(SegmentedReceive):
		return segmentation.SegmentedReceive
	case string(NoSegmentation):
		return segmentation.NoSegmentation
	default:
		return segmentation.NoSegmentation
	}
}

func (inst *Instance) doWriteBool(networkUUID, deviceUUID string, pnt *network.Point, value uint32) error {
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return err
	}

	err = dev.PointWriteBool(pnt, value)
	if err != nil {
		return err
	}
	return nil

}

func (inst *Instance) whoIs(networkUUID string, opts *bacnet.WhoIsOpts) (resp []btypes.Device, err error) {
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	devices, err := net.Whois(opts)
	if err != nil {
		return nil, err
	}
	return devices, err
}

func (inst *Instance) devicePoints(deviceUUID string) (resp []*network.PointDetails, err error) {
	getNetwork, err := inst.db.GetNetworkByDeviceUUID(deviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	net, err := inst.getBacnetNetwork(getNetwork.UUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return nil, err
	}

	resp, err = dev.GetDevicePoints(btypes.ObjectInstance(dev.DeviceID))
	if err != nil {
		return nil, err
	}
	return

}

func intToUint32(value int) uint32 {
	var y = uint32(value)
	return y
}

func Unit32ToFloat64(value uint32) float64 {
	var y = float64(value)
	return y
}
func float32ToFloat64(value float32) float64 {
	var y = float64(value)
	return y
}

func float64ToUint32(value float64) uint32 {
	var y = uint32(value)
	return y
}

func float64ToFloat32(value float64) float32 {
	var y = float32(value)
	return y
}
