package slackboard

import (
	"github.com/juju/ratelimit"
	"time"
)

// QPSPerSlackEndpoint controls rate limiting.
type QPSPerSlackEndpoint struct {
	bucket  *ratelimit.Bucket
	maxWait *time.Duration
}

// NewQPSPerSlackEndpoint initializes QPSPerSlackEndpoint.
func NewQPSPerSlackEndpoint(conf ConfToml) *QPSPerSlackEndpoint {
	qps := conf.Core.QPS
	if qps <= 0 {
		return nil
	}

	var maxWait *time.Duration
	duration := conf.Core.MaxDelayDuration
	if 0 <= duration {
		sec := time.Duration(duration) * time.Second
		maxWait = &sec
	}

	return &QPSPerSlackEndpoint{
		ratelimit.NewBucketWithRate(float64(qps), int64(qps)),
		maxWait,
	}
}

// Available takes count from the bucket.
// If it is not available immediately, do nothing and return false.
func (qpsend QPSPerSlackEndpoint) Available() bool {
	if qpsend.bucket.TakeAvailable(1) == 0 {
		return false
	}

	return true
}

// WaitAndAvailable waits until the bucket becomes available
func (qpsend QPSPerSlackEndpoint) WaitAndAvailable() bool {
	maxWait := qpsend.maxWait
	if maxWait == nil {
		// disable rate limiting
		return true
	}
	return qpsend.bucket.WaitMaxDuration(1, *maxWait)
}
