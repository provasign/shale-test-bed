// Package auth holds the (toy) credential logic for the test bed.
package auth

import (
	"crypto/subtle"
	"sync"
)

// users is an in-memory credential table. A real service would hash these;
// this repo exists only to generate plausible diffs for Shale demos.
var users = map[string]string{
	"alice": "wonderland",
	"bob":   "builder",
}

// maxFailures is how many consecutive failed attempts lock an account.
const maxFailures = 5

var (
	mu       sync.Mutex
	failures = map[string]int{}
)

// Check reports whether the user/password pair is valid. After maxFailures
// consecutive failures the account is locked and Check always returns false
// until the process restarts.
func Check(user, password string) bool {
	mu.Lock()
	defer mu.Unlock()
	if failures[user] > maxFailures {
		return false
	}
	want, ok := users[user]
	if ok && subtle.ConstantTimeCompare([]byte(want), []byte(password)) == 1 {
		failures[user] = 0
		return true
	}
	failures[user]++
	return false
}

// Locked reports whether the account is currently locked out.
func Locked(user string) bool {
	mu.Lock()
	defer mu.Unlock()
	return failures[user] > maxFailures
}
