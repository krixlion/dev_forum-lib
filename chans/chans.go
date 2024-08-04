package chans

// FanIn directs all messages from given channels to the returned channel.
func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)

	for _, c := range channels {
		go func(c <-chan T) {
			for v := range c {
				out <- v
			}
		}(c)
	}

	return out
}

// NonBlockSend tries to send given value through the given channel in a non-blocking manner.
// If the send would block then this is func aborts the message and do nothing.
func NonBlockSend[T any](c chan T, v T) {
	select {
	case c <- v:
	default:
		return
	}
}
