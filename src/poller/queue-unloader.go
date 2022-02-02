package poller

import (
	"time"
)

// LOOK AT USING:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

// Polling Manager Summary:
//  - Diagram Summary: https://docs.google.com/drawings/d/1priwsaQ6EryRBx1kLQd91REJvHzFyxz7cOHYYXyBNFE/edit?usp=sharing
//  - The QueueUnloader is the only way to get the next PollPoint from a Queue
//  - When a QueueUnloader is stopped, the Worker go routine is closed and the reference to the QueueUnloader is set to nil.

//Questions:
// -

//There should be a function in Modbus(or other protocols) that submits the polling point to the protocol client, then when the poll is completed, it starts a timeout to add the polling point to the queue again.
// NEXT FETCH THE FF POINT AND use time.AfterFunc(DURATION, )
//dbhandler.GormDatabase.GetPoint(pp.FFPointUUID)

type QueueUnloader struct {
	NextPollPoint   *PollingPoint
	NextUnloadTimer *time.Timer
}

func (pm *NetworkPollManager) StartQueueUnloader() {
	ql := &QueueUnloader{nil, nil}
	pm.NetworkQueueUnloader = ql
	if pm.NetworkQueueUnloader.NextPollPoint == nil {
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.NetworkQueueUnloader.NextPollPoint = pp
		}
	}
}

func (pm *NetworkPollManager) StopQueueUnloader() {
	pm.NetworkQueueUnloader.NextUnloadTimer.Stop()
	pm.NetworkQueueUnloader = nil
}

func (pm *NetworkPollManager) postNextPointCallback() {
	if pm.NetworkQueueUnloader != nil && pm.NetworkQueueUnloader.NextPollPoint != nil {
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.NetworkQueueUnloader.NextPollPoint = pp
		}
	}
}

func (pm *NetworkPollManager) GetNextPollingPoint() (pp *PollingPoint, callback func(pp *PollingPoint, writeSuccess, readSuccess bool)) {
	if pm.NetworkQueueUnloader != nil && pm.NetworkQueueUnloader.NextPollPoint != nil {
		pp := pm.NetworkQueueUnloader.NextPollPoint
		pm.NetworkQueueUnloader.NextPollPoint = nil
		pm.NetworkQueueUnloader.NextUnloadTimer = time.AfterFunc(*pm.MaxPollRate, pm.postNextPointCallback)
		return pp, pm.PollingPointCompleteNotification
	}
	return nil, nil
}
