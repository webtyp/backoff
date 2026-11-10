## 1. Overview in One Minute
This repository contains a simple backoff package which calculates how long a caller should wait before retrying an action. It guarantees that the delay doubles on each attempt up to a strict cap, starting from a base duration. It never handles retries, sleeping, or random number generation directly.

## 2. Repository Pieces & On-Disk Layout
| Piece | Responsibility |
|---|---|
| `backoff.go` | Provides the `Policy` struct, `New` constructor, and `DelayAfter` method. |
| `tests/backoff_test.go` | Externally verifies edge cases and rule conformance. |
| `docs/TESTS.md` | Maps each test in the suite to the specific business rule it proves. |
| `docs/diagrams/DELAY_FLOW.md` | Flowchart outlining the logic of `DelayAfter`. |
| `docs/diagrams/PLACEMENT.md` | Diagram outlining where the package fits into the larger ecosystem. |

## 3. Step-by-Step Flows

### delay after a failure
The happy path allows the caller to determine a progressively longer wait time without surpassing a defined threshold. It calls `DelayAfter` on a valid `Policy` created with `New`.

The error paths are triggered when:
1. `New` is passed a zero or negative first delay.
2. `New` is passed a longest delay that is smaller than the first delay.
3. `DelayAfter` is called on a zero `Policy`.

The full visual flowchart of this procedure is found in [DELAY_FLOW.md](diagrams/DELAY_FLOW.md).

## 4. Security Boundaries & Contracts
* A zero `Policy` is strictly rejected and yields an error instead of operating with unsafe defaults.
* There is no jitter because a reliable random source is unavailable in the WASM target; avoiding jitter also keeps the code fully deterministic.
* There is no `Sleep` since the callers may operate concurrently and handles event looping differently across consumers.
* Only `webtyp.com/time` is used, and the standard library is never imported to comply with TinyGo limits.
