# backoff
<img src="docs/img/badges.svg">

How long to wait before the next retry: the delay doubles after each failure, from a first delay up to a longest delay.

## Why This Exists

A retry loop that waits a fixed 1 s hammers a service that is down; one that doubles without a cap waits an hour. This package provides a fixed multiplier of 2, a first delay, and a capped longest delay. It operates on `webtyp.com/time`. The entire process is detailed in the [DELAY_FLOW.md](docs/diagrams/DELAY_FLOW.md).

**Without this vs. With this**
* Without this: Every developer writes a private loop (or includes `math/rand` and `time.Sleep`).
* With this: One consistent policy creation and a single call per failure, without internal sleeping or external stdlib imports.

## Before You Start

```bash
go get webtyp.com/backoff@v0.2.0
```

## Step-by-Step Guide

### Step 1: Create the policy
**Goal:** Define the delays.
```go
package main

import (
	"webtyp.com/backoff"
	"webtyp.com/fmt"
	"webtyp.com/time"
)

func main() {
	policy, err := backoff.New(time.Second, 60*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(policy)
}
```
**Expected output:** `{1000000000 60000000000}`

### Step 2: Compute the delay after each failure
**Goal:** Find how long to wait after `failures` attempts.
```go
package main

import (
	"webtyp.com/backoff"
	"webtyp.com/fmt"
	"webtyp.com/time"
)

func main() {
	policy, _ := backoff.New(time.Second, 60*time.Second)

	d1, _ := policy.DelayAfter(1)
	fmt.Println(d1) // 1000000000

	d2, _ := policy.DelayAfter(2)
	fmt.Println(d2) // 2000000000

	d3, _ := policy.DelayAfter(3)
	fmt.Println(d3) // 4000000000
}
```
**Expected output:**
```
1000000000
2000000000
4000000000
```

### Step 3: What the errors mean
**Goal:** Handle creation or usage errors properly. See [backoff.go](backoff.go) for details.
```go
package main

import (
	"webtyp.com/backoff"
	"webtyp.com/fmt"
	"webtyp.com/time"
)

func main() {
	_, err1 := backoff.New(-1, time.Second)
	fmt.Println(err1)

	_, err2 := backoff.New(2*time.Second, time.Second)
	fmt.Println(err2)

	var empty backoff.Policy
	_, err3 := empty.DelayAfter(1)
	fmt.Println(err3)
}
```
**Expected output:** (The respective error texts)

## Documentation Map

| Order | Document | What question it answers & when to read it |
|---|---|---|
| 1 | `docs/ARCHITECTURE.md` | What are the boundaries, flows, and structural decisions of this package? |
| 2 | `docs/diagrams/DELAY_FLOW.md` | How is the delay computed under different failure scenarios? |
| 3 | `docs/diagrams/PLACEMENT.md` | Where does this package sit in the ecosystem? |
| 4 | `docs/TESTS.md` | What specific edge cases and rules does this package formally verify? |
