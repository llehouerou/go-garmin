// ratelimit.go
package garmin

import (
	"context"
	"sync"
	"time"
)

type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 15,
		BurstSize:         5,
	}
}

type rateLimiter struct {
	interval time.Duration
	burst    int
	tokens   int
	lastTime time.Time
	mu       sync.Mutex
}

func newRateLimiter(cfg RateLimitConfig) *rateLimiter {
	interval := time.Minute / time.Duration(cfg.RequestsPerMinute)
	return &rateLimiter{
		interval: interval,
		burst:    cfg.BurstSize,
		tokens:   cfg.BurstSize,
		lastTime: time.Now(),
	}
}

func (r *rateLimiter) Wait(ctx context.Context) error {
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.lastTime)
		r.lastTime = now

		newTokens := int(elapsed / r.interval)
		r.tokens += newTokens
		if r.tokens > r.burst {
			r.tokens = r.burst
		}

		if r.tokens > 0 {
			r.tokens--
			r.mu.Unlock()
			return nil
		}

		waitTime := r.interval - (elapsed % r.interval)
		r.mu.Unlock()

		select {
		case <-time.After(waitTime):
			continue // Loop back to try acquiring a token
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
