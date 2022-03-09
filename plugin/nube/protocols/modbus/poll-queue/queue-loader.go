package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
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
	fmt.Println("PrintPollQueuePointUUIDs TOTAL COUNT = ", hasNextPollPoint+pm.PollQueue.PriorityQueue.Len()+pm.PollQueue.PointsOnHold.Len())
	fmt.Print("NextPollPoint: ")
	fmt.Printf("%+v\n", pm.PluginQueueUnloader.NextPollPoint)
	fmt.Print("PollQueue: COUNT = ", pm.PollQueue.PriorityQueue.Len(), ": ")
	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		fmt.Print(pp.FFPointUUID, ", ", pp.PollPriority, "; ")
	}
	fmt.Println("")
	fmt.Print("PointsOnHold COUNT = ", pm.PollQueue.PointsOnHold.Len(), ": ")
	for _, pp := range pm.PollQueue.PointsOnHold.PriorityQueue {
		fmt.Print(pp.FFPointUUID, ", ", pp.PollPriority, ", repoll timer:", pp.RepollTimer != nil, "; ")
	}
	fmt.Println("\n \n")
}

func (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool) {
	log.Infof("modbus-poll: PollingPointCompleteNotification Point UUID: %s, writeSuccess: %t, readSuccess: %t", pp.FFPointUUID, writeSuccess, readSuccess)

	point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID)
	if err != nil {
		fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s /n", pp.FFPointUUID)
	}
	//fmt.Printf("NetworkPollManager.PollingPointCompleteNotification(): writeMode: %s", point.WriteMode)
	//fmt.Println("")

	switch point.WriteMode {
	case poller.ReadOnce: //ReadOnce          If read_successful then don't re-add.
		point.WritePollRequired = utils.NewFalse()
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
		} else {
			point.ReadPollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp)
		}
	case poller.ReadOnly: //ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.WritePollRequired = utils.NewFalse()
		//fmt.Println("ReadOnly: point")
		//fmt.Printf("%+v\n", point)
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			//TODO: point.PollTimer PROPERTY CAUSES FF TO CRASH ON START REMOVED FOR TESTING
			//point.PollTimer = time.AfterFunc(pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID), pm.MakePollingPointRepollCallback(pp))
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			//log.Info("duration: ", duration)
			//time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp))
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp))
			fmt.Println("Modbus PollingPointCompleteNotification(): pp")
			fmt.Printf("%+v\n", pp)
			addSuccess := pm.PollQueue.PointsOnHold.AddPollingPoint(pp)
			if !addSuccess {
				log.Error(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to PointsOnHold slice.  (%s)", pp.FFPointUUID))
			}
		} else {
			point.ReadPollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
	case poller.WriteOnce: //WriteOnce         If write_successful then don't re-add.
		point.ReadPollRequired = utils.NewFalse()
		if writeSuccess {
			point.WritePollRequired = utils.NewFalse()
		} else {
			point.WritePollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
	case poller.WriteOnceReadOnce: //WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		if writeSuccess {
			point.WritePollRequired = utils.NewFalse()
		} else {
			point.WritePollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
		if readSuccess {
			point.ReadPollRequired = utils.NewFalse()
		} else {
			point.ReadPollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
	case poller.WriteAlways: //WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = utils.NewFalse()
		point.WritePollRequired = utils.NewTrue()
		if writeSuccess {
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			//point.PollTimer = time.AfterFunc(pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID), pm.MakePollingPointRepollCallback(pp))
		} else {
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
	case poller.WriteOnceThenRead: //WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = utils.NewTrue()
		if writeSuccess {
			point.WritePollRequired = utils.NewFalse()
		} else {
			point.WritePollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
		if readSuccess {
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			//point.PollTimer = time.AfterFunc(pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID), pm.MakePollingPointRepollCallback(pp))
		} else {
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		}
	case poller.WriteAndMaintain: //WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = utils.NewTrue()
		writeValue := *point.Priority.GetHighestPriorityValue()
		presentValue := *point.PresentValue
		if presentValue != writeValue {
			point.WritePollRequired = utils.NewTrue()
			pm.PollQueue.AddPollingPoint(pp) //re-add to poll queue immediately
		} else {
			point.WritePollRequired = utils.NewFalse()
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			//point.PollTimer = time.AfterFunc(pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID), pm.MakePollingPointRepollCallback(pp))
		}
	}

}

func (pm *NetworkPollManager) MakePollingPointRepollCallback(pp *PollingPoint) func() {
	//log.Info("MakePollingPointRepollCallback()")
	f := func() {
		log.Info("CALL PollingPointRepollCallback func() pp:")
		fmt.Printf("%+v\n", pp)
		pp.RepollTimer = nil
		removeSuccess := pm.PollQueue.PointsOnHold.RemovePollingPointByPointUUID(pp.FFPointUUID)
		if !removeSuccess {
			log.Error(fmt.Sprintf("Modbus MakePollingPointRepollCallback(): polling point could not be found in PointsOnHold.  (%s)", pp.FFPointUUID))
		}
		pm.PollQueue.AddPollingPoint(pp)
	}
	return f
}
