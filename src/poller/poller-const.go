package poller

type PollPriority string

const (
	PRIORITY_ASAP   PollPriority = "asap"
	PRIORITY_HIGH   PollPriority = "high"
	PRIORITY_NORMAL PollPriority = "normal"
	PRIORITY_LOW    PollPriority = "low"
)

type PollRate string

const (
	RATE_FAST   PollRate = "fast"
	RATE_NORMAL PollRate = "normal"
	RATE_SLOW   PollRate = "slow"
)
