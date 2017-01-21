package slackboard

import (
	"github.com/juju/ratelimit"
)

type QPSPerSlackEndpoint struct {
	bucket *ratelimit.Bucket
}

func NewQPSPerSlackEndpoint(conf ConfToml) *QPSPerSlackEndpoint {
	qps := conf.Core.QPS
	if qps <= 0 {
		return nil
	}

	return &QPSPerSlackEndpoint{
		ratelimit.NewBucketWithRate(float64(qps), int64(qps)),
	}
}

func (qpsend QPSPerSlackEndpoint) Available() bool {
	if qpsend.bucket.TakeAvailable(1) == 0 {
		return false
	}

	return true
}
