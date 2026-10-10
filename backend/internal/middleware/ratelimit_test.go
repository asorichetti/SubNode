package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(5, 1*time.Minute)

	// Should allow 5 requests
	for i := 0; i < 5; i++ {
		if !rl.Allow("test-key") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	if rl.Allow("test-key") {
		t.Error("6th request should be blocked")
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := NewRateLimiter(2, 1*time.Minute)

	// First key should allow 2 requests
	if !rl.Allow("key1") {
		t.Error("First request for key1 should be allowed")
	}
	if !rl.Allow("key1") {
		t.Error("Second request for key1 should be allowed")
	}
	if rl.Allow("key1") {
		t.Error("Third request for key1 should be blocked")
	}

	// Second key should have its own limit
	if !rl.Allow("key2") {
		t.Error("First request for key2 should be allowed")
	}
	if !rl.Allow("key2") {
		t.Error("Second request for key2 should be allowed")
	}
	if rl.Allow("key2") {
		t.Error("Third request for key2 should be blocked")
	}
}

func TestRateLimiter_WindowExpiration(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)

	// Use up the limit
	if !rl.Allow("test-key") {
		t.Error("First request should be allowed")
	}
	if !rl.Allow("test-key") {
		t.Error("Second request should be allowed")
	}
	if rl.Allow("test-key") {
		t.Error("Third request should be blocked")
	}

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	if !rl.Allow("test-key") {
		t.Error("Request after window expiry should be allowed")
	}
}
