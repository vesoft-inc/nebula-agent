package limiter

import (
	"github.com/juju/ratelimit"
)

var limiter *ratelimit.Bucket

func SetAgentRateLimiter(limit int) {
	if limit > 0 {
		bps := float64(limit * (1 << 20) / 8)
		limiter = ratelimit.NewBucketWithRate(bps, int64(bps)*3)
	}
}

func Wait(size int64) {
	if limiter != nil {
		limiter.Wait(size)
	}
}

func IsSet() bool {
	return limiter != nil
}
