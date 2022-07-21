package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

type NetworkPriorityPollQueue struct {
	Config                    *Config
	PriorityQueue             *PriorityPollQueue // This is the queue that is polling points are drawn from
	StandbyPollingPoints      *PriorityPollQueue // This is a slice that contains polling points that are not in the active polling queue, it is mostly a reference so that we can periodically find out if any points have been dropped from polling.
	OutstandingPollingPoints  *PriorityPollQueue // this is a slice that contains polling points that are currently out for polling.
	PointsUpdatedWhilePolling map[string]bool    // UUIDs of points that have been updated while they were out for polling.  bool is true if the point needs to be written ASAP
	QueueUnloader             *QueueUnloader
	FFPluginUUID              string
	FFNetworkUUID             string
	ActiveDevicesList         []string // UUIDs of devices that have points in the queue
}

func (nq *NetworkPriorityPollQueue) AddPollingPoint(pp *PollingPoint) bool {
	nq.pollQueueDebugMsg("NetworkPriorityPollQueue AddPollingPoint(): ", pp.FFPointUUID)
	if pp.FFNetworkUUID != nq.FFNetworkUUID {
		nq.pollQueueErrorMsg(fmt.Sprintf("NetworkPriorityPollQueue.AddPollingPoint: PollingPoint FFNetworkUUID does not match the queue FFNetworkUUID. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID))
		if pp.LockupAlertTimer != nil {
			pp.LockupAlertTimer.Stop()
		}
		return false
	}
	if nq.PriorityQueue.GetPollingPointIndexByPointUUID(pp.FFPointUUID) != -1 {
		log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: PollingPoint %s already exists in polling queue. \n", pp.FFPointUUID)
		if pp.LockupAlertTimer != nil {
			pp.LockupAlertTimer.Stop()
		}
		return false
	}
	if nq.StandbyPollingPoints.GetPollingPointIndexByPointUUID(pp.FFPointUUID) != -1 {
		// point exists in the StandbyPollingPoints list, remove it and add immediately.
		nq.RemovePollingPointByPointUUID(pp.FFPointUUID)
	}
	if nq.OutstandingPollingPoints.GetPollingPointIndexByPointUUID(pp.FFPointUUID) != -1 {
		_, ok := nq.PointsUpdatedWhilePolling[pp.FFPointUUID]
		if !ok {
			nq.PointsUpdatedWhilePolling[pp.FFPointUUID] = false
		}
		return true
	}

	pp.QueueEntryTime = time.Now().Unix()
	success := nq.PriorityQueue.AddPollingPoint(pp)
	if !success {
		// log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: point already exists in poll queue. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID)
		return false
	}
	nq.AddDeviceToActiveDevicesList(pp.FFDeviceUUID)
	return true
}

func (nq *NetworkPriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) (pp *PollingPoint, success bool) {
	nq.pollQueueDebugMsg("RemovePollingPointByPointUUID(): ", pointUUID)
	pp = nil
	if nq.QueueUnloader != nil && nq.QueueUnloader.NextPollPoint != nil && nq.QueueUnloader.NextPollPoint.FFPointUUID == pointUUID {
		pp = nq.QueueUnloader.NextPollPoint
		if nq.QueueUnloader.NextPollPoint.LockupAlertTimer != nil {
			nq.QueueUnloader.NextPollPoint.LockupAlertTimer.Stop()
		}
		nq.QueueUnloader.NextPollPoint = nil

	}
	pp, _ = nq.PriorityQueue.RemovePollingPointByPointUUID(pointUUID)
	if pp != nil {
		nq.pollQueueDebugMsg("RemovePollingPointByPointUUID(): Point is in the PriorityQueue")
	}
	pp, _ = nq.StandbyPollingPoints.RemovePollingPointByPointUUID(pointUUID)
	if pp != nil {
		nq.pollQueueDebugMsg("RemovePollingPointByPointUUID(): Point is in the StandbyPollingPoints Queue")
	}
	return pp, true
}

func (nq *NetworkPriorityPollQueue) RemovePollingPointByDeviceUUID(deviceUUID string) bool {
	nq.pollQueueDebugMsg("RemovePollingPointByDeviceUUID(): ", deviceUUID)
	nq.PriorityQueue.RemovePollingPointByDeviceUUID(deviceUUID)
	nq.StandbyPollingPoints.RemovePollingPointByDeviceUUID(deviceUUID)
	nq.OutstandingPollingPoints.RemovePollingPointByDeviceUUID(deviceUUID)
	nq.RemoveDeviceFromActiveDevicesList(deviceUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority model.PollPriority) bool {
	nq.PriorityQueue.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	nq.StandbyPollingPoints.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	return true
}
func (nq *NetworkPriorityPollQueue) GetPollingPointByPointUUID(pointUUID string) (pp *PollingPoint, err error) {
	nq.pollQueueDebugMsg("NetworkPriorityPollQueue GetPollingPointByPointUUID(): ", pointUUID)
	pp = nil
	if nq.QueueUnloader != nil && nq.QueueUnloader.NextPollPoint != nil && nq.QueueUnloader.NextPollPoint.FFPointUUID == pointUUID {
		pp = nq.QueueUnloader.NextPollPoint
		return pp, nil
	}
	pollQueueIndex := nq.PriorityQueue.GetPollingPointIndexByPointUUID(pointUUID)
	if pollQueueIndex != -1 {
		return nq.PriorityQueue.PriorityQueue[pollQueueIndex], nil
	}
	standbyIndex := nq.StandbyPollingPoints.GetPollingPointIndexByPointUUID(pointUUID)
	if standbyIndex != -1 {
		return nq.StandbyPollingPoints.PriorityQueue[standbyIndex], nil
	}
	outstandingIndex := nq.OutstandingPollingPoints.GetPollingPointIndexByPointUUID(pointUUID)
	if outstandingIndex != -1 {
		return nq.OutstandingPollingPoints.PriorityQueue[outstandingIndex], nil
	}

	return nil, errors.New(fmt.Sprint("couldn't find point: ", pointUUID))
}

func (nq *NetworkPriorityPollQueue) GetNextPollingPoint() (*PollingPoint, error) {
	pp, err := nq.PriorityQueue.GetNextPollingPoint()
	if err != nil {
		// nq.pollQueueDebugMsg(fmt.Sprintf("NetworkPriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue. FFNetworkUUID: %s \n", nq.FFNetworkUUID))
		return nil, errors.New(fmt.Sprintf("NetworkPriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue"))
	}
	return pp, nil
}
func (nq *NetworkPriorityPollQueue) Start() {
	// nq.PriorityQueue.Start()
}
func (nq *NetworkPriorityPollQueue) Stop() {
	// nq.PriorityQueue.Stop()
	nq.EmptyQueue()
}
func (nq *NetworkPriorityPollQueue) EmptyQueue() {
	nq.PriorityQueue.EmptyQueue()
	refQueue := make([]*PollingPoint, 0)
	rq := &PriorityPollQueue{refQueue}
	nq.StandbyPollingPoints = rq
	outstandingQueue := make([]*PollingPoint, 0)
	opq := &PriorityPollQueue{outstandingQueue}
	nq.StandbyPollingPoints = opq
}
func (nq *NetworkPriorityPollQueue) CheckIfActiveDevicesListIncludes(devUUID string) bool {
	for _, dev := range nq.ActiveDevicesList {
		if dev == devUUID {
			return true
		}
	}
	return false
}
func (nq *NetworkPriorityPollQueue) AddDeviceToActiveDevicesList(devUUID string) bool {
	for _, dev := range nq.ActiveDevicesList {
		if dev == devUUID {
			return false
		}
	}
	nq.ActiveDevicesList = append(nq.ActiveDevicesList, devUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) RemoveDeviceFromActiveDevicesList(devUUID string) bool {
	for index, dev := range nq.ActiveDevicesList {
		if dev == devUUID {
			// remove the devUUID from ActiveDevicesList
			nq.ActiveDevicesList[index] = nq.ActiveDevicesList[len(nq.ActiveDevicesList)-1]
			nq.ActiveDevicesList = nq.ActiveDevicesList[:len(nq.ActiveDevicesList)-1]
			return true
		}
	}
	return false
}
func (nq *NetworkPriorityPollQueue) CheckPollingQueueForDevUUID(devUUID string) bool {
	for _, pp := range nq.PriorityQueue.PriorityQueue {
		if pp.FFDeviceUUID == devUUID {
			return true
		}
	}
	for _, pp := range nq.StandbyPollingPoints.PriorityQueue {
		if pp.FFDeviceUUID == devUUID {
			return true
		}
	}
	for _, pp := range nq.OutstandingPollingPoints.PriorityQueue {
		if pp.FFDeviceUUID == devUUID {
			return true
		}
	}
	return false
}

// THIS IS THE BASE PriorityPollQueue Type and defines the base methods used to implement the `heap` library.  https://pkg.go.dev/container/heap
type PriorityPollQueue struct {
	// Enable        bool
	PriorityQueue []*PollingPoint
}

func (q *PriorityPollQueue) Len() int { return len(q.PriorityQueue) }
func (q *PriorityPollQueue) Less(i, j int) bool {
	if len(q.PriorityQueue) <= i && len(q.PriorityQueue) <= j {
		return false
	}
	iPriority := q.PriorityQueue[i].PollPriority
	iPriorityNum := 0
	switch iPriority {
	case model.PRIORITY_ASAP:
		iPriorityNum = 0
	case model.PRIORITY_HIGH:
		iPriorityNum = 1
	case model.PRIORITY_NORMAL:
		iPriorityNum = 2
	case model.PRIORITY_LOW:
		iPriorityNum = 3
	}
	jPriority := q.PriorityQueue[j].PollPriority
	jPriorityNum := 0
	switch jPriority {
	case model.PRIORITY_ASAP:
		jPriorityNum = 0
	case model.PRIORITY_HIGH:
		jPriorityNum = 1
	case model.PRIORITY_NORMAL:
		jPriorityNum = 2
	case model.PRIORITY_LOW:
		jPriorityNum = 3
	}

	return iPriorityNum < jPriorityNum
}
func (q *PriorityPollQueue) Swap(i, j int) {
	q.PriorityQueue[i], q.PriorityQueue[j] = q.PriorityQueue[j], q.PriorityQueue[i]
}
func (q *PriorityPollQueue) Push(x interface{}) {
	item := x.(*PollingPoint)
	q.PriorityQueue = append(q.PriorityQueue, item)
}
func (q *PriorityPollQueue) Pop() interface{} {
	old := q.PriorityQueue
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	q.PriorityQueue = old[0 : n-1]
	return item
}
func (q *PriorityPollQueue) GetPollingPointIndexByPointUUID(pointUUID string) int {
	for index, pp := range q.PriorityQueue {
		if pp.FFPointUUID == pointUUID {
			return index
		}
	}
	return -1
}
func (q *PriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) (pp *PollingPoint, success bool) {
	pp = nil
	index := q.GetPollingPointIndexByPointUUID(pointUUID)
	if index >= 0 {
		pp = heap.Remove(q, index).(*PollingPoint)
		if pp != nil {
			// pollQueueDebugMsg("RemovePollingPointByPointUUID() pp: %+v\n", pp)
		}
		if pp.RepollTimer != nil {
			pp.RepollTimer.Stop()
		}
		if pp.LockupAlertTimer != nil {
			pp.LockupAlertTimer.Stop()
		}
		return pp, true
	}
	return pp, false
}
func (q *PriorityPollQueue) RemovePollingPointByDeviceUUID(deviceUUID string) bool {
	index := 0
	for index < q.Len() {
		pp := q.PriorityQueue[index]
		if pp.FFDeviceUUID == deviceUUID {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
			}
			heap.Remove(q, index)
			if pp.LockupAlertTimer != nil {
				pp.LockupAlertTimer.Stop()
			}
		} else {
			index++
		}
	}
	return true
}
func (q *PriorityPollQueue) RemovePollingPointByNetworkUUID(networkUUID string) bool {
	index := 0
	for index < q.Len() {
		pp := q.PriorityQueue[index]
		if pp.FFNetworkUUID == networkUUID {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
			}
			heap.Remove(q, index)
			if pp.LockupAlertTimer != nil {
				pp.LockupAlertTimer.Stop()
			}
		} else {
			index++
		}
	}
	return true
}
func (q *PriorityPollQueue) AddPollingPoint(pp *PollingPoint) bool {
	index := q.GetPollingPointIndexByPointUUID(pp.FFPointUUID)
	if index == -1 {
		heap.Push(q, pp)
		return true
	}
	return false
}
func (q *PriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority model.PollPriority) bool {
	index := q.GetPollingPointIndexByPointUUID(pointUUID)
	if index >= 0 {
		q.PriorityQueue[index].PollPriority = newPriority
		heap.Fix(q, index)
		return true
	}
	return false
}

// func (q *PriorityPollQueue) Start() { q.Enable = true }  //TODO: add queue startup code
// func (q *PriorityPollQueue) Stop()  { q.Enable = false } //TODO: add queue stop code
func (q *PriorityPollQueue) EmptyQueue() {
	for q.Len() > 0 {
		heap.Pop(q)
	}
}
func (q *PriorityPollQueue) GetNextPollingPoint() (*PollingPoint, error) {
	if q.Len() > 0 {
		pp := heap.Pop(q).(*PollingPoint)
		return pp, nil
	}
	return nil, errors.New("PriorityPollQueue is not enabled")
}

type PollingPoint struct {
	PollPriority     model.PollPriority
	FFPointUUID      string
	FFDeviceUUID     string
	FFNetworkUUID    string
	FFPluginUUID     string
	RepollTimer      *time.Timer
	QueueEntryTime   int64
	LockupAlertTimer *time.Timer
}

func NewPollingPoint(ffPointUUID, ffDeviceUUID, ffNetworkUUID, ffPluginUUID string) *PollingPoint {
	pp := &PollingPoint{model.PRIORITY_NORMAL, ffPointUUID, ffDeviceUUID, ffNetworkUUID, ffPluginUUID, nil, 0, nil}
	// WHATEVER FUNCTION CALLS NewPollingPoint NEEDS TO SET THE PRIORITY
	return pp
}

func NewPollingPointWithPriority(ffPointUUID, ffDeviceUUID, ffNetworkUUID, ffPluginUUID string, priority model.PollPriority) *PollingPoint {
	pp := &PollingPoint{priority, ffPointUUID, ffDeviceUUID, ffNetworkUUID, ffPluginUUID, nil, 0, nil}
	return pp
}

func (nq *NetworkPriorityPollQueue) pollQueueDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(nq.Config.LogLevel, "DEBUG") {
		prefix := "Modbus Poll Queue: "
		log.Info(prefix, args)
	}
}

func (nq *NetworkPriorityPollQueue) pollQueueErrorMsg(args ...interface{}) {
	prefix := "Modbus Poll Queue: "
	log.Error(prefix, args)
}
