package slackboard

import (
	"sync/atomic"
	"time"
)

var (
	WaitClients  int32
	LimitClients int32
)

func init() {
	WaitClients = 0
	LimitClients = 1
}

func rateLimitStart() {
	for atomic.LoadInt32(&WaitClients) >= LimitClients {
		// Slack does not allow to send
		// no more than one message per second.
		// But what bursts over the limit for short periods is allowed.
		// refs: https://api.slack.com/docs/rate-limits
		time.Sleep(1)
	}
	atomic.AddInt32(&WaitClients, 1)
}

func rateLimitEnd() {
	atomic.AddInt32(&WaitClients, -1)
}
