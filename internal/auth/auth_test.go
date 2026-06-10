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

func TestLockoutAfterFiveFailures(t *testing.T) {
	failures = map[string]int{} // reset state

	for i := 0; i < 5; i++ {
		Check("bob", "wrong")
	}
	if !Locked("bob") {
		t.Fatal("account should be locked after 5 consecutive failures")
	}
	if Check("bob", "builder") {
		t.Fatal("locked account must fail even with the right password")
	}
}

func TestSuccessResetsFailureCount(t *testing.T) {
	failures = map[string]int{} // reset state

	Check("alice", "wrong")
	Check("alice", "wrong")
	if !Check("alice", "wonderland") {
		t.Fatal("correct password should still work before lockout")
	}
	if failures["alice"] != 0 {
		t.Fatalf("success should reset failures, got %d", failures["alice"])
	}
}
