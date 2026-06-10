// Package auth holds the (toy) credential logic for the test bed.
package auth

import "crypto/subtle"

// users is an in-memory credential table. A real service would hash these;
// this repo exists only to generate plausible diffs for Shale demos.
var users = map[string]string{
	"alice": "wonderland",
	"bob":   "builder",
}

// Check reports whether the user/password pair is valid.
func Check(user, password string) bool {
	want, ok := users[user]
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(want), []byte(password)) == 1
}
