package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}

// or takes one or more "done" channels and returns a single channel that will be
// closed as soon as ANY of the input channels is closed.
func or(channels ...<-chan any) <-chan any {
	length := len(channels)

	if length == 0 {
		return nil
	}
	if length == 1 {
		return channels[0]
	}

	doneChan := make(chan any)

	go func() {
		defer close(doneChan)

		switch length {
		// If we only have 2 channels, listen to them directly.
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			// For more than 2 channels, we can't write a select with arbitrary length,
			// so we use recursion:
			// - Listen to the first two channels directly
			// - For the remaining channels, recursively call or(...)
			// - We also append `doneChan` into the recursive call,
			//   so that if a deeper recursive level detects closure, it can
			//   propagate the signal all the way back up to this goroutine.
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or(append(channels[2:], doneChan)...):
			}
		}
	}()

	return doneChan
}
