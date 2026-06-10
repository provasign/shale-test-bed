package ratelimit

import (
	"testing"
	"time"
)

func TestAllowExhaustsAndRefills(t *testing.T) {
	l := New(10)
	clock := time.Unix(0, 0)
	l.now = func() time.Time { return clock }

	for i := 0; i < 10; i++ {
		if !l.Allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if l.Allow("1.2.3.4") {
		t.Fatal("11th request should be limited")
	}
	if !l.Allow("5.6.7.8") {
		t.Fatal("different key should not be limited")
	}

	clock = clock.Add(6 * time.Second) // refills one token at 10/min
	if !l.Allow("1.2.3.4") {
		t.Fatal("request after refill should be allowed")
	}
}
