package tests

import (
	"testing"

	"webtyp.com/backoff"
	"webtyp.com/time"
)

const second time.Nano = 1000000000

func newPolicy(t *testing.T, first, longest time.Nano) backoff.Policy {
	t.Helper()
	p, err := backoff.New(first, longest)
	if err != nil {
		t.Fatalf("New(%d, %d): %v", first, longest, err)
	}
	return p
}

func wantDelay(t *testing.T, p backoff.Policy, failures int64, want time.Nano) {
	t.Helper()
	got, err := p.DelayAfter(failures)
	if err != nil {
		t.Fatalf("DelayAfter(%d): %v", failures, err)
	}
	if got != want {
		t.Fatalf("DelayAfter(%d) = %d, want %d", failures, got, want)
	}
}

func TestPolicy_FirstFailureWaitsFirstDelay(t *testing.T) {
	p := newPolicy(t, second, 60*second)
	wantDelay(t, p, 1, second)
}

func TestPolicy_DoublesAfterEachFailure(t *testing.T) {
	p := newPolicy(t, second, 60*second)
	wantDelay(t, p, 1, second)
	wantDelay(t, p, 2, 2*second)
	wantDelay(t, p, 3, 4*second)
	wantDelay(t, p, 4, 8*second)
	wantDelay(t, p, 5, 16*second)
	wantDelay(t, p, 6, 32*second)
}

func TestPolicy_NeverExceedsLongest(t *testing.T) {
	p := newPolicy(t, second, 60*second)
	wantDelay(t, p, 7, 60*second)
	wantDelay(t, p, 8, 60*second)
	wantDelay(t, p, 100, 60*second)
}

func TestPolicy_HugeFailureCountDoesNotOverflow(t *testing.T) {
	p := newPolicy(t, second, 60*second)
	wantDelay(t, p, 63, 60*second)
	wantDelay(t, p, 64, 60*second)
	wantDelay(t, p, 65, 60*second)
	wantDelay(t, p, 9223372036854775807, 60*second)
}

func TestPolicy_ZeroOrNegativeFailuresWaitFirstDelay(t *testing.T) {
	p := newPolicy(t, second, 60*second)
	wantDelay(t, p, 0, second)
	wantDelay(t, p, -5, second)
}

func TestPolicy_FirstEqualsLongestNeverGrows(t *testing.T) {
	p := newPolicy(t, 5*second, 5*second)
	wantDelay(t, p, 1, 5*second)
	wantDelay(t, p, 2, 5*second)
	wantDelay(t, p, 10, 5*second)
}

func TestPolicy_LongestNotMultipleOfFirst(t *testing.T) {
	p := newPolicy(t, 3*second, 10*second)
	wantDelay(t, p, 1, 3*second)
	wantDelay(t, p, 2, 6*second)
	wantDelay(t, p, 3, 10*second)
	wantDelay(t, p, 4, 10*second)
}

func TestNew_RejectsNonPositiveFirst(t *testing.T) {
	p1, err1 := backoff.New(0, 60*second)
	if err1 != backoff.ErrFirstDelayNotPositive {
		t.Fatalf("expected ErrFirstDelayNotPositive, got %v", err1)
	}
	if p1 != (backoff.Policy{}) {
		t.Fatalf("expected zero policy, got %v", p1)
	}

	p2, err2 := backoff.New(-1, 60*second)
	if err2 != backoff.ErrFirstDelayNotPositive {
		t.Fatalf("expected ErrFirstDelayNotPositive, got %v", err2)
	}
	if p2 != (backoff.Policy{}) {
		t.Fatalf("expected zero policy, got %v", p2)
	}
}

func TestNew_RejectsLongestBelowFirst(t *testing.T) {
	_, err := backoff.New(2*second, second)
	if err != backoff.ErrLongestBelowFirst {
		t.Fatalf("expected ErrLongestBelowFirst, got %v", err)
	}
}

func TestPolicy_ZeroValueIsRejected(t *testing.T) {
	p := backoff.Policy{}
	got, err := p.DelayAfter(1)
	if err != backoff.ErrPolicyNotCreated {
		t.Fatalf("expected ErrPolicyNotCreated, got %v", err)
	}
	if got != 0 {
		t.Fatalf("expected 0 delay, got %v", got)
	}
}
