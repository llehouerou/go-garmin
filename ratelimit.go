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

		newTokens := int(elapsed / r.interval)
		if newTokens > 0 {
			// Only advance lastTime by the time that produced tokens,
			// preserving fractional remainder for the next check.
			r.lastTime = r.lastTime.Add(time.Duration(newTokens) * r.interval)
			r.tokens += newTokens
			if r.tokens > r.burst {
				r.tokens = r.burst
			}
		}

		if r.tokens > 0 {
			r.tokens--
			r.mu.Unlock()
			return nil
		}

		// Wait until the next token is due
		nextToken := r.lastTime.Add(r.interval)
		waitTime := nextToken.Sub(now)
		if waitTime <= 0 {
			waitTime = r.interval
		}
		r.mu.Unlock()

		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
