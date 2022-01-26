package poller

import (
	"container/heap"
	"time"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

//Priority Polling Summary:
//  - Protocol client runs as a worker go routine, pulls jobs from the protocol plugin priority queue.  Rate is dictated by the availability of the protocol client.
//  - Protocol plugin priority queue is fed by the (multiple) network priority queues. One feeder queue for each network, should respect the network polling delays (etc).
//  - Network priority queues are fed by device priority queues.
//  - Device priority queues are fed by points using `time.Ticker` triggers on each point (configured based on push rate setting).
//  - Device priority queues check that the device priority queues don't already have that point in them.  Or they have a flag that is reset when they are polled.
//  - in all priority queues the most significant (lowest int) priority is selected first.

//Questions:
// - at what level should we specify the fast, normal, and slow poll rates?  Plugin? Network? Device?  I'm thinking Device level
// - Should write values should be given a higher priority in the poll queue

type ProtocolPollQueue struct {
	MaxPollRate   time.Duration
	PriorityQueue PriorityPollQueue
}

type DevicePriorityPollQueue struct {
	FastPollRate   time.Duration
	NormalPollRate time.Duration
	SlowPollRate   time.Duration
	PriorityQueue  PriorityPollQueue
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
func (q *PriorityPollQueue) GetPollingPointIndex(uuid string) int {
	for index, pp := range q.PriorityQueue {
		if pp.FFPointUUID == uuid {
			return index
		}
	}
	return -1
}
func (q *PriorityPollQueue) RemovePollingPoint(uuid string) {
	index := q.GetPollingPointIndex(uuid)
	heap.Remove(q, index)
}
func (q *PriorityPollQueue) AddPollingPoint(FFPointUUID string, PollPriority PollPriority) bool {
	if q.GetPollingPointIndex(FFPointUUID) >= 0 {
		pp := NewPollingPoint(FFPointUUID)
		pp.PollPriority = PollPriority
		heap.Push(q, pp)
		return true
	}
	return false
}

//func (q *PriorityPollQueue) UpdatePollingPoint(uuid string, priority PollPriority, pollRate PollRate, pollRateTime time.Duration) {
func (q *PriorityPollQueue) UpdatePollingPoint(uuid string, priority PollPriority) {
	index := q.GetPollingPointIndex(uuid)
	q.PriorityQueue[index].PollPriority = priority
	//q.PriorityQueue[index].PollRate = pollRate
	//q.PriorityQueue[index].PollRateTime = pollRateTime
	heap.Fix(q, index)
}
func (q *PriorityPollQueue) Start() { q.Enable = true }  //TODO: add queue startup code
func (q *PriorityPollQueue) Stop()  { q.Enable = false } //TODO: add queue stop code
func (q *PriorityPollQueue) EmptyQueue() {
	for q.Len() > 0 {
		heap.Pop(q)
	}
}

type PollingPoint struct {
	FFPointUUID  string
	PollPriority PollPriority
	//PollRate     PollRate      //TODO: this might not belong in this type.  Poll timing is probably handled by the routine that adds points to the Queue
	//PollRateTime time.Duration //TODO: this might not belong in this type.  Poll timing is probably handled by the routine that adds points to the Queue
}

func NewPollingPoint(FFPointUUID string) *PollingPoint {
	return &PollingPoint{FFPointUUID, PRIORITY_NORMAL}
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
