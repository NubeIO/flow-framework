package poller

import (
	"container/heap"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

//Priority Polling Summary:
//  - Diagram Summary: https://docs.google.com/drawings/d/1priwsaQ6EryRBx1kLQd91REJvHzFyxz7cOHYYXyBNFE/edit?usp=sharing
//                    https://docs.google.com/drawings/d/1DPltaxg1Pm-xguK36QUbHyKNFauBbuO9G07GtQ9rFYs/edit?usp=sharing
//  - Protocol client runs as a worker go routine, pulls jobs from the ProtocolPriorityPollQueue.  Rate is dictated by the availability of the protocol client.
//  - ProtocolPriorityPollQueue is fed by the (multiple) NetworkPriorityPollQueue. One feeder queue for each network, should respect the network polling delays (etc).
//  - Network priority queues are fed by device priority queues.
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

//There should be a function in Modbus(or other protocols) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

type NetworkPriorityPollQueue struct {
	PriorityQueue PriorityPollQueue
	MaxPollRate   time.Duration
	FFNetworkUUID string
}

func (nq *NetworkPriorityPollQueue) AddPollingPoint(pp *PollingPoint) bool {
	if pp.FFNetworkUUID != nq.FFNetworkUUID {
		log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: PollingPoint FFNetworkUUID does not match the queue FFNetworkUUID. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID)
		return false
	}
	success := nq.PriorityQueue.AddPollingPoint(pp)
	if !success {
		log.Errorf("NetworkPriorityPollQueue.AddPollingPoint: point already exists in poll queue. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pp.FFPointUUID)
		return false
	}
	return true
}
func (nq *NetworkPriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) bool {
	success := nq.PriorityQueue.RemovePollingPointByPointUUID(pointUUID)
	if !success {
		log.Errorf("NetworkPriorityPollQueue.RemovePollingPointByPointUUID: point does not exists in poll queue. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pointUUID)
		return false
	}
	return true
}
func (nq *NetworkPriorityPollQueue) RemovePollingPointByDeviceUUID(deviceUUID string) bool {
	nq.PriorityQueue.RemovePollingPointByDeviceUUID(deviceUUID)
	return true
}
func (nq *NetworkPriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority PollPriority) bool {
	success := nq.PriorityQueue.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	if !success {
		log.Errorf("NetworkPriorityPollQueue.UpdatePollingPointByPointUUID: point does not exists in poll queue. FFNetworkUUID: %s  FFPointUUID: %s \n", nq.FFNetworkUUID, pointUUID)
		return false
	}
	return true
}
func (nq *NetworkPriorityPollQueue) GetNextPollingPoint() (*PollingPoint, error) {
	if nq.PriorityQueue.Enable {
		pp, err := nq.PriorityQueue.GetNextPollingPoint()
		if err != nil {
			log.Errorf("NetworkPriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue. FFNetworkUUID: %s \n", nq.FFNetworkUUID)
		}
		return pp, nil
	}
	return nil, errors.New(fmt.Sprintf("NetworkPriorityPollQueue %s is not enabled.", nq.FFNetworkUUID))
}
func (nq *NetworkPriorityPollQueue) Start() { nq.PriorityQueue.Start() } //TODO: add queue startup code
func (nq *NetworkPriorityPollQueue) Stop()  { nq.PriorityQueue.Stop() }  //TODO: add queue stop code
func (nq *NetworkPriorityPollQueue) EmptyQueue() {
	nq.PriorityQueue.EmptyQueue()
}

type DevicePriorityPollQueue struct {
	PriorityQueue      *PriorityPollQueue
	NetworkQueue       *NetworkPriorityPollQueue
	NetworkQueueLoader *QueueLoader
	FastPollRate       time.Duration
	NormalPollRate     time.Duration
	SlowPollRate       time.Duration
	FFDeviceUUID       string
}

func (dq *DevicePriorityPollQueue) AddPollingPoint(pp *PollingPoint) bool {
	if pp.FFDeviceUUID != dq.FFDeviceUUID {
		log.Errorf("DevicePriorityPollQueue.AddPollingPoint: PollingPoint FFDeviceUUID does not match the queue FFDeviceUUID. FFDeviceUUID: %s  FFPointUUID: %s \n", dq.FFDeviceUUID, pp.FFPointUUID)
		return false
	}
	success := dq.PriorityQueue.AddPollingPoint(pp)
	if !success {
		log.Errorf("DevicePriorityPollQueue.AddPollingPoint: point already exists in poll queue. FFDeviceUUID: %s  FFPointUUID: %s \n", dq.FFDeviceUUID, pp.FFPointUUID)
		return false
	}
	return true
}
func (dq *DevicePriorityPollQueue) RemovePollingPointByPointUUID(pointUUID string) bool {
	success := dq.PriorityQueue.RemovePollingPointByPointUUID(pointUUID)
	if !success {
		log.Errorf("DevicePriorityPollQueue.RemovePollingPointByPointUUID: point does not exists in poll queue. FFDeviceUUID: %s  FFPointUUID: %s \n", dq.FFDeviceUUID, pointUUID)
		return false
	}
	return true
}
func (dq *DevicePriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority PollPriority) bool {
	success := dq.PriorityQueue.UpdatePollingPointByPointUUID(pointUUID, newPriority)
	if !success {
		log.Errorf("DevicePriorityPollQueue.UpdatePollingPointByPointUUID: point does not exists in poll queue. FFDeviceUUID: %s  FFPointUUID: %s \n", dq.FFDeviceUUID, pointUUID)
		return false
	}
	return true
}
func (dq *DevicePriorityPollQueue) GetNextPollingPoint() (*PollingPoint, error) {
	if dq.PriorityQueue.Enable {
		pp, err := dq.PriorityQueue.GetNextPollingPoint()
		if err != nil {
			log.Errorf("DevicePriorityPollQueue.GetNextPollingPoint: no PollingPoints in queue. FFDeviceUUID: %s \n", dq.FFDeviceUUID)
		}
		return pp, nil
	}
	return nil, errors.New(fmt.Sprintf("DevicePriorityPollQueue %s is not enabled.", dq.FFDeviceUUID))
}
func (dq *DevicePriorityPollQueue) Start() { //TODO: find how to get FFNetwork, use it to get the *NetworkPriorityPollQueue
	dq.PriorityQueue.Start()
	// GET FF NETWORK, THEN GET THE *NetworkPriorityPollQueue
	//nq := GET FFNETWORK.PollingQueue
	dq.NetworkQueueLoader = dq.NewQueueLoader()
	go dq.NetworkQueueLoader.Worker(dq, np, dq.NetworkQueueLoader.CancelChan)
}
func (dq *DevicePriorityPollQueue) Stop() { //TODO: find how to get FFNetwork, use it to get the *NetworkPriorityPollQueue
	dq.PriorityQueue.Stop()
	close(dq.NetworkQueueLoader.CancelChan)
	dq.NetworkQueueLoader = nil
	dq.EmptyQueue()
}
func (dq *DevicePriorityPollQueue) EmptyQueue() {
	dq.PriorityQueue.EmptyQueue()
}
func (*DevicePriorityPollQueue) NewQueueLoader() *QueueLoader {
	cancelChan := make(chan struct{})
	f := func(dq *DevicePriorityPollQueue, nq *NetworkPriorityPollQueue, cancelChan <-chan struct{}) {
		for {
			select {
			case <-cancelChan:
				return

			default:
				if dq.PriorityQueue.Enable && dq.PriorityQueue.Len() > 0 {
					pp, err := dq.GetNextPollingPoint()
					if err == nil {
						if nq.PriorityQueue.Enable {
							nq.AddPollingPoint(pp)
						}
					}
				}
			}
		}
	}
	ql := &QueueLoader{f, cancelChan}
	return ql
}

/*
func NewQueueLoader(dq *DevicePriorityPollQueue, nq *NetworkPriorityPollQueue) *QueueLoader {
	pointChan := make(chan *PollingPoint)
	cancelChan := make(chan struct{})
	//f := func(pointChan <-chan *PollingPoint, cancelChan <-chan struct{}) {
	f := func(dq *DevicePriorityPollQueue, nq *NetworkPriorityPollQueue, pointChan <-chan *PollingPoint, cancelChan <-chan struct{}) {
		for {
			select {
			case <-cancelChan:
				return

			case point := <-pointChan:
				if dq.PriorityQueue.Enable && dq.PriorityQueue.Len() > 0 {
					pp, err := dq.GetNextPollingPoint()
					if err != nil {
					} else {
						if nq.PriorityQueue.Enable {
							nq.AddPollingPoint(pp)
						}
					}
				}
			}
		}
	}
	ql := &QueueLoader{f, pointChan, cancelChan}
	return ql
}
*/

type QueueLoader struct {
	Worker func(dq *DevicePriorityPollQueue, nq *NetworkPriorityPollQueue, cancelChan <-chan struct{})
	//PointChan		chan *PollingPoint
	CancelChan chan struct{}
}

type PriorityPollQueue struct {
	Enable        bool
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
	for index, pp := range q.PriorityQueue {
		if pp.FFDeviceUUID == deviceUUID {
			heap.Remove(q, index)
		}
	}
	return true
}
func (q *PriorityPollQueue) RemovePollingPointByNetworkUUID(networkUUID string) bool {
	for index, pp := range q.PriorityQueue {
		if pp.FFNetworkUUID == networkUUID {
			heap.Remove(q, index)
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
func (q *PriorityPollQueue) UpdatePollingPointByPointUUID(pointUUID string, newPriority PollPriority) bool {
	index := q.GetPollingPointIndexByPointUUID(pointUUID)
	if index >= 0 {
		q.PriorityQueue[index].PollPriority = newPriority
		heap.Fix(q, index)
		return true
	}
	return false
}
func (q *PriorityPollQueue) Start() { q.Enable = true }  //TODO: add queue startup code
func (q *PriorityPollQueue) Stop()  { q.Enable = false } //TODO: add queue stop code
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
	PollPriority  PollPriority
	FFPointUUID   string
	FFDeviceUUID  string
	FFNetworkUUID string
	FFPluginUUID  string
}

func NewPollingPoint(FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID string) *PollingPoint {
	pp := &PollingPoint{PRIORITY_NORMAL, FFPointUUID, FFDeviceUUID, FFNetworkUUID, FFPluginUUID}
	// FIND THE FUNCTION TO GET THE FF POINT BY UUID
	// If we are storing the PollPriority and PollRate in the FF Point, get and set those values.
	// SET THE PollRateTime based on Device Poll Rates
	return pp
}

type PollPriority int

const (
	PRIORITY_HIGH   PollPriority = iota //Only Read Point Value.
	PRIORITY_NORMAL                     //Write the value on COV, don't Read.
	PRIORITY_LOW                        //Write the value on every poll (poll rate defined by setting).
)

type PollRate int

const (
	RATE_FAST   PollRate = iota //Only Read Point Value.
	RATE_NORMAL                 //Write the value on COV, don't Read.
	RATE_SLOW                   //Write the value on every poll (poll rate defined by setting).
)
