package or

import (
	"testing"
	"time"
)

func TestOr_NoChannels(t *testing.T) {
	if Or() != nil {
		t.Error("expected nil when no channels passed")
	}
}

func TestOr_SingleChannel(t *testing.T) {
	ch := make(chan interface{})
	if got := Or(ch); got != ch {
		t.Error("expected the same channel for single input")
	}
}

func TestOr_MultipleChannels(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	out := Or(ch1, ch2)

	go func() {
		defer close(ch2)
		time.Sleep(1 * time.Second)
	}()

	select {
	case <-out:
	case <-time.After(2 * time.Second):
		t.Error("expected the channel to be closed after 1s")
	}
}
