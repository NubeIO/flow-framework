package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/src/poller"
	log "github.com/sirupsen/logrus"
	"time"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

//Priority Polling Summary:
//  - Diagram Summary: https://docs.google.com/drawings/d/1priwsaQ6EryRBx1kLQd91REJvHzFyxz7cOHYYXyBNFE/edit?usp=sharing
//  - Protocol client runs as a worker go routine, pulls jobs from the ProtocolPriorityPollQueue.  Rate is dictated by the availability of the protocol client.
//  - ProtocolPriorityPollQueue is fed by the (multiple) NetworkPriorityPollQueue. One feeder queue for each network, should respect the network polling delays (etc).
//  - Device priority queues are fed by points using `time.Ticker` triggers on each point (configured based on push rate setting).
//  - Device priority queues check that the device priority queues don't already have that point in them.  Or they have a flag that is reset when they are polled.
//  - In all priority queues the most significant (lowest int) priority is selected first.

//Questions:
// - at what level should we specify the fast, normal, and slow poll rates?  Plugin? Network? Device?  I'm thinking Device level
// - Are there device poll rate limitations? Rather than setting at the network level.
// - Should write values should be given a higher priority in the poll queue.  I think probably.  High Priority Writes -> High Priority Reads -> Normal Priority Writes -> Normal Priority Reads -> etc
// - How do I get FF Points by UUID?
// - Are FF Points shared by multiple plugins?
//     - Can I store a Timer as a new property in FF Points?
//     - Can I store a PollRate and PollPriority in FF Points?

// TODO: Add in special PollPoints that are for bundled operations.  Should support multiple protocols (maybe the bundle properties are dependent on the plugin?)
//There should be a function in Modbus(or other protocols) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

type NetworkPriorityPollQueue struct {
	PriorityQueue        *PriorityPollQueue //This is the queue that is polling points are drawn from
	StandbyPollingPoints *PriorityPollQueue //This is a slice that contains polling points that are not in the active polling queue, it is mostly a reference so that we can periodically find out if any points have been dropped from polling.
	FFPluginUUID         string
	FFNetworkUUID        string
	ActiveDevicesList    []string //UUIDs of devices that have points in the queue
}

func (nq *NetworkPriorityPollQueue) AddPollingPoint(pp *PollingPoint) bool {
	if pp.FFNetworkUUID != nq.FFNetworkUUID {
		log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: PollingPoint FFNetworkUUID does not match the queue FFNetworkUUID. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID)
		return false
	}
	pp.QueueEntryTime = time.Now().Unix()
	success := nq.PriorityQueue.AddPollingPoint(pp)
	if !success {
		//log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: point already exists in poll queue. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID)
		return false
	}
	nq.AddDeviceToActiveDevicesList(pp.FFDeviceUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) bool {
	fmt.Println("RemovePollingPointByPointUUID(): ", pointUUID)
	nq.PriorityQueue.RemovePollingPointByPointUUID(pointUUID)
	nq.StandbyPollingPoints.RemovePollingPointByPointUUID(pointUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) RemovePollingPointByDeviceUUID(deviceUUID string) bool {
	fmt.Println("RemovePollingPointByDeviceUUID(): ", deviceUUID)
	nq.PriorityQueue.RemovePollingPointByDeviceUUID(deviceUUID)
	nq.StandbyPollingPoints.RemovePollingPointByDeviceUUID(deviceUUID)
	nq.RemoveDeviceFromActiveDevicesList(deviceUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority poller.PollPriority) bool {
	nq.PriorityQueue.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	nq.StandbyPollingPoints.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	return true
}
func (nq *NetworkPriorityPollQueue) GetNextPollingPoint() (*PollingPoint, error) {
	//if nq.PriorityQueue.Enable {
	pp, err := nq.PriorityQueue.GetNextPollingPoint()
	if err != nil {
		log.Errorf("NetworkPriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue. FFNetworkUUID: %s \n", nq.FFNetworkUUID)
		return nil, errors.New(fmt.Sprintf("NetworkPriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue"))
	}
	return pp, nil
	//}
	//fmt.Println("(nq *NetworkPriorityPollQueue) GetNextPollingPoint(): nq.PriorityQueue.Enable is false ")
	//return nil, errors.New(fmt.Sprintf("NetworkPriorityPollQueue %s is not enabled.", nq.FFNetworkUUID))
}
func (nq *NetworkPriorityPollQueue) Start() {
	//nq.PriorityQueue.Start()
}
func (nq *NetworkPriorityPollQueue) Stop() {
	//nq.PriorityQueue.Stop()
	nq.EmptyQueue()
}
func (nq *NetworkPriorityPollQueue) EmptyQueue() {
	nq.PriorityQueue.EmptyQueue()
	refQueue := make([]*PollingPoint, 0)
	rq := &PriorityPollQueue{refQueue}
	nq.StandbyPollingPoints = rq
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
			//remove the devUUID from ActiveDevicesList
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
	return false
}

// THIS IS THE BASE PriorityPollQueue Type and defines the base methods used to implement the `heap` library.  https://pkg.go.dev/container/heap
type PriorityPollQueue struct {
	//Enable        bool
	PriorityQueue []*PollingPoint
}

func (q *PriorityPollQueue) Len() int { return len(q.PriorityQueue) }
func (q *PriorityPollQueue) Less(i, j int) bool {
	return q.PriorityQueue[i].PollPriority < q.PriorityQueue[j].PollPriority
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
func (q *PriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) bool {
	index := q.GetPollingPointIndexByPointUUID(pointUUID)
	if index >= 0 {
		heap.Remove(q, index)
		return true
	}
	return false
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
			heap.Remove(q, index)
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
func (q *PriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority poller.PollPriority) bool {
	index := q.GetPollingPointIndexByPointUUID(pointUUID)
	if index >= 0 {
		q.PriorityQueue[index].PollPriority = newPriority
		heap.Fix(q, index)
		return true
	}
	return false
}

//func (q *PriorityPollQueue) Start() { q.Enable = true }  //TODO: add queue startup code
//func (q *PriorityPollQueue) Stop()  { q.Enable = false } //TODO: add queue stop code
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
	PollPriority   poller.PollPriority
	FFPointUUID    string
	FFDeviceUUID   string
	FFNetworkUUID  string
	FFPluginUUID   string
	RepollTimer    *time.Timer
	QueueEntryTime int64
}

func NewPollingPoint(FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID string) *PollingPoint {
	pp := &PollingPoint{poller.PRIORITY_NORMAL, FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID, nil, 0}
	//WHATEVER FUNCTION CALLS NewPollingPoint NEEDS TO SET THE PRIORITY
	return pp
}

func NewPollingPointWithPriority(FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID string, priority poller.PollPriority) *PollingPoint {
	pp := &PollingPoint{priority, FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID, nil, 0}
	return pp
}
