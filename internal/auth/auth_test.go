package auth

import "testing"

func TestCheck(t *testing.T) {
	cases := []struct {
		user, pass string
		want       bool
	}{
		{"alice", "wonderland", true},
		{"alice", "wrong", false},
		{"nobody", "anything", false},
	}
	for _, c := range cases {
		if got := Check(c.user, c.pass); got != c.want {
			t.Errorf("Check(%q, %q) = %v, want %v", c.user, c.pass, got, c.want)
		}
	}
}
