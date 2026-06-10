// Package ratelimit provides a simple in-memory token-bucket limiter keyed
// by client IP. Good enough for a single instance; a real deployment would
// use a shared store.
package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Limiter allows burst requests per key, refilling at rate per minute.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64 // tokens per second
	burst   float64
	now     func() time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

// New returns a Limiter allowing perMinute requests per key with the same burst.
func New(perMinute int) *Limiter {
	return &Limiter{
		buckets: make(map[string]*bucket),
		rate:    float64(perMinute) / 60.0,
		burst:   float64(perMinute),
		now:     time.Now,
	}
}

// Allow reports whether the key may proceed, consuming a token if so.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	now := l.now()
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Middleware wraps h, returning 429 when the client IP exceeds the limit.
func (l *Limiter) Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !l.Allow(ip) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		h(w, r)
	}
}
