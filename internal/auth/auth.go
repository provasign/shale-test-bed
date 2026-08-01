// Package auth holds the (toy) credential logic for the test bed.
package auth

import (
	"crypto/subtle"
	"sync"
	"unicode"
)

// users is an in-memory credential table. A real service would hash these;
// this repo exists only to generate plausible diffs for Shale demos.
var users = map[string]string{
	"alice": "Wonderland!1",
	"bob":   "BuilderTool!1",
}

// maxFailures is how many consecutive failed attempts lock an account.
const maxFailures = 5

const (
	minPasswordLength = 12
	maxPasswordLength = 128
)

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
	if failures[user] >= maxFailures {
		return false
	}
	want, ok := users[user]
	if ok && ValidPassword(password) && subtle.ConstantTimeCompare([]byte(want), []byte(password)) == 1 {
		failures[user] = 0
		return true
	}
	failures[user]++
	return false
}

// ValidPassword reports whether password satisfies the demo password policy.
func ValidPassword(password string) bool {
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// Locked reports whether the account is currently locked out.
func Locked(user string) bool {
	mu.Lock()
	defer mu.Unlock()
	return failures[user] >= maxFailures
}
