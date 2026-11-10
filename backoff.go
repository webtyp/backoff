// Package backoff says how long to wait before the next retry: the delay
// doubles after each failure, from a first delay up to a longest delay.
//
// It never sleeps, never retries and never reads a clock; the caller does.
//
//	policy, err := backoff.New(time.Second, 60*time.Second) // 1 s first, 60 s longest
//	if err != nil { return err }
//	wait, err := policy.DelayAfter(3) // 4 s: the delay after 3 failures
//
// Flow: https://github.com/webtyp/backoff/blob/main/docs/diagrams/DELAY_FLOW.md
package backoff

import "webtyp.com/time"

// backoffError is the concrete type of this package's errors. Compare with
// ==, never errors.Is: == between two error values is fine on the stdlib
// toolchain and the constants are the only values this package returns.
type backoffError string

func (e backoffError) Error() string { return string(e) }

const (
	// ErrFirstDelayNotPositive: New received a first delay of zero or less.
	ErrFirstDelayNotPositive backoffError = "backoff: the first delay must be greater than zero, so no policy was created and nothing was waited or retried. Pass a first delay such as time.Second to backoff.New."

	// ErrLongestBelowFirst: New received a longest delay smaller than the first.
	ErrLongestBelowFirst backoffError = "backoff: the longest delay is smaller than the first delay, so no policy was created and nothing was waited or retried. Pass a longest delay equal to or greater than the first delay to backoff.New."

	// ErrPolicyNotCreated: DelayAfter was called on a Policy that did not come from New.
	ErrPolicyNotCreated backoffError = "backoff: this policy was not created with backoff.New, so it has no delays and nothing was waited. Create it with backoff.New(first, longest) and call DelayAfter on the result."
)

// Policy holds the first and the longest delay. Only New creates a valid one;
// the zero Policy answers ErrPolicyNotCreated.
type Policy struct {
	first   time.Duration
	longest time.Duration
}

// New creates a Policy. first must be greater than zero and longest must be at
// least first.
//
// Example: policy, err := backoff.New(time.Second, 60*time.Second)
func New(first, longest time.Duration) (Policy, error) {
	if first <= 0 {
		return Policy{}, ErrFirstDelayNotPositive
	}
	if longest < first {
		return Policy{}, ErrLongestBelowFirst
	}
	return Policy{first: first, longest: longest}, nil
}

// DelayAfter answers how long to wait after the given number of failed
// attempts: the first delay after 1 failure (and after 0 or a negative count),
// then doubling after each further failure, never more than the longest delay.
//
// Example: wait, err := policy.DelayAfter(3) // first delay of 1 s → 4 s
func (p Policy) DelayAfter(failures int64) (time.Duration, error) {
	if p.first <= 0 {
		return 0, ErrPolicyNotCreated
	}
	delay := p.first
	for step := int64(1); step < failures; step++ {
		// delay*2 >= longest, written without multiplying so it cannot overflow
		if delay >= p.longest-delay {
			return p.longest, nil
		}
		delay *= 2
	}
	return delay, nil
}
