package parser

import "errors"

// Validate rejects input the parser must never see.
func Validate(s string) error {
	if s == "" {
		return errors.New("empty input")
	}
	if len(s) > 1<<16 {
		return errors.New("input too large")
	}
	return nil
}
