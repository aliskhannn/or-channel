# or-channel

A small Go utility package that merges multiple **done channels** into one.  
The resulting channel is closed as soon as **any** of the input channels is closed.  

This pattern is useful for coordinating goroutines and handling cancellation in concurrent code.

---

## Installation

```bash
go get github.com/aliskhannn/or-channel
````

---

## Usage

```go
package main

import (
	"fmt"
	"time"

	or "github.com/aliskhannn/or-channel"
)

func sig(after time.Duration) <-chan any {
	c := make(chan any)
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()

	<-or.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}
```

---

## Development

### Run example

```bash
make run
```

### Run tests

```bash
make test
```

---

## Project structure

```
.
├── or/
│   ├── or.go
│   └── or_test.go
└── examples/
│   └── main.go
├── Makefile
├── README.md
```