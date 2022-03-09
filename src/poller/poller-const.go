package poller

type PollPriority int

const (
	PRIORITY_ASAP PollPriority = iota //Only Read Point Value.
	PRIORITY_HIGH
	PRIORITY_NORMAL
	PRIORITY_LOW
)

type PollRate int

const (
	RATE_FAST   PollRate = iota //Only Read Point Value.
	RATE_NORMAL                 //Write the value on COV, don't Read.
	RATE_SLOW                   //Write the value on every poll (poll rate defined by setting).
)
