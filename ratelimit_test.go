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

func TestRateLimiterNoInfiniteLoop(t *testing.T) {
	// Regression test: the rate limiter previously entered an infinite loop
	// when tokens were depleted and Wait was called at a non-aligned time.
	// The bug: lastTime was reset to now on every iteration, splitting elapsed
	// time across iterations so integer division never produced a token.
	rl := newRateLimiter(RateLimitConfig{
		RequestsPerMinute: 15, // 4 second interval (same as production)
		BurstSize:         1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Exhaust burst
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	// Sleep to create a non-aligned elapsed time.
	// This triggers the bug: elapsed splits across loop iterations,
	// and integer division never produces a token.
	time.Sleep(1 * time.Second)

	// This would hang forever with the old code
	start := time.Now()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("second wait failed (may have hung): %v", err)
	}
	elapsed := time.Since(start)

	// Should complete in ~3 seconds (4s interval minus 1s already elapsed)
	if elapsed > 5*time.Second {
		t.Errorf("wait took %v, expected ~3s", elapsed)
	}
	if elapsed < 2*time.Second {
		t.Errorf("wait was too fast (%v), expected ~3s", elapsed)
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
