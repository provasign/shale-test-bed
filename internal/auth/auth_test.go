package auth

import "testing"

func TestCheck(t *testing.T) {
	cases := []struct {
		user, pass string
		want       bool
	}{
		{"alice", "Wonderland!1", true},
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
	if Check("bob", "BuilderTool!1") {
		t.Fatal("locked account must fail even with the right password")
	}
}

func TestSuccessResetsFailureCount(t *testing.T) {
	failures = map[string]int{} // reset state

	Check("alice", "wrong")
	Check("alice", "wrong")
	if !Check("alice", "Wonderland!1") {
		t.Fatal("correct password should still work before lockout")
	}
	if failures["alice"] != 0 {
		t.Fatalf("success should reset failures, got %d", failures["alice"])
	}
}

func TestValidPassword(t *testing.T) {
	cases := []struct {
		name, password string
		want           bool
	}{
		{"valid", "CorrectHorse!1", true},
		{"too short", "Short!1", false},
		{"missing lowercase", "PASSWORDONLY!1", false},
		{"missing uppercase", "passwordonly!1", false},
		{"missing digit", "PasswordOnly!", false},
		{"missing special", "PasswordOnly1", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ValidPassword(c.password); got != c.want {
				t.Fatalf("ValidPassword(%q) = %v, want %v", c.password, got, c.want)
			}
		})
	}
}
