package limiter

import (
	"github.com/juju/ratelimit"
)

var Rate rateLimiter

type rateLimiter struct {
	limiter *ratelimit.Bucket
}

func (r *rateLimiter) SetLimiter(limit int) {
	if limit > 0 {
		bps := float64(limit * (1 << 20) / 8)
		r.limiter = ratelimit.NewBucketWithRate(bps, int64(bps)*3)
	}
}

func (r *rateLimiter) Wait(size int64) {
	if r.limiter != nil {
		r.limiter.Wait(size)
	}
}

func (r *rateLimiter) IsSet() bool {
	return r.limiter != nil
}
