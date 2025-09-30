package or

// Or takes one Or more "done" channels and returns a single channel that will be
// closed as soon as ANY of the input channels is closed.
func Or(channels ...<-chan any) <-chan any {
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
			// - For the remaining channels, recursively call Or(...)
			// - We also append `doneChan` into the recursive call,
			//   so that if a deeper recursive level detects closure, it can
			//   propagate the signal all the way back up to this goroutine.
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-Or(append(channels[2:], doneChan)...):
			}
		}
	}()

	return doneChan
}
