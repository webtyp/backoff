| Test | Rule |
|---|---|
| [TestPolicy_FirstFailureWaitsFirstDelay](../tests/backoff_test.go) | `New(1s, 60s)` → `DelayAfter(1)` → `1s` |
| [TestPolicy_DoublesAfterEachFailure](../tests/backoff_test.go) | `New(1s, 60s)` → `DelayAfter(1..6)` → `1s, 2s, 4s, 8s, 16s, 32s` |
| [TestPolicy_NeverExceedsLongest](../tests/backoff_test.go) | `New(1s, 60s)` → `DelayAfter(7)`, `DelayAfter(8)`, `DelayAfter(100)` → `60s` each |
| [TestPolicy_HugeFailureCountDoesNotOverflow](../tests/backoff_test.go) | `New(1s, 60s)` → `DelayAfter(63)`, `(64)`, `(65)`, `(9223372036854775807)` → `60s` each |
| [TestPolicy_ZeroOrNegativeFailuresWaitFirstDelay](../tests/backoff_test.go) | `New(1s, 60s)` → `DelayAfter(0)`, `DelayAfter(-5)` → `1s` each |
| [TestPolicy_FirstEqualsLongestNeverGrows](../tests/backoff_test.go) | `New(5s, 5s)` → `DelayAfter(1)`, `(2)`, `(10)` → `5s` each |
| [TestPolicy_LongestNotMultipleOfFirst](../tests/backoff_test.go) | `New(3s, 10s)` → `DelayAfter(1..4)` → `3s, 6s, 10s, 10s` |
| [TestNew_RejectsNonPositiveFirst](../tests/backoff_test.go) | `New(0, 60s)` and `New(-1, 60s)` → `err == backoff.ErrFirstDelayNotPositive` and a zero `Policy` |
| [TestNew_RejectsLongestBelowFirst](../tests/backoff_test.go) | `New(2s, 1s)` → `err == backoff.ErrLongestBelowFirst` |
| [TestPolicy_ZeroValueIsRejected](../tests/backoff_test.go) | `backoff.Policy{}` → `DelayAfter(1)` → `0` and `err == backoff.ErrPolicyNotCreated` |
