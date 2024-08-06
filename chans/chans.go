package chans

import "sync"

// FanIn merges given channels into one and returns it.
// All messages from given channels are directed to the merged channel.
// FanIn does not close provided channels.
// Merged channel will be closed automatically when all of the given channels are closed.
func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, c := range channels {
		go func(c <-chan T) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// NonBlockSend tries to send given value through the given channel in a non-blocking manner.
// If the send would block then this is func aborts the message and do nothing.
func NonBlockSend[T any](c chan<- T, v T) {
	select {
	case c <- v:
	default:
		return
	}
}
