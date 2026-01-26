// ratelimit_test.go
package garmin

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := newRateLimiter(RateLimitConfig{
		RequestsPerMinute: 60, // 1 per second
		BurstSize:         2,
	})

	ctx := context.Background()

	// First two should be immediate (burst)
	start := time.Now()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("second wait failed: %v", err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Error("burst requests should be immediate")
	}

	// Third should wait
	start = time.Now()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("third wait failed: %v", err)
	}
	if time.Since(start) < 900*time.Millisecond {
		t.Error("third request should have waited")
	}
}

func TestRateLimiterContextCancellation(t *testing.T) {
	rl := newRateLimiter(RateLimitConfig{
		RequestsPerMinute: 60, // 1 per second
		BurstSize:         1,
	})

	ctx, cancel := context.WithCancel(context.Background())

	// Exhaust the limiter
	_ = rl.Wait(ctx)

	// Cancel context before next wait
	cancel()

	// Next wait should fail immediately due to cancelled context
	if err := rl.Wait(ctx); err == nil {
		t.Error("expected context cancellation error")
	}
}
