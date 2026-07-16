package events

// Option configures an EventBus at construction time.
type Option func(*EventBus)

// WithBufferSize sets the channel buffer depth for every subscriber.
// A full buffer causes the event to be dropped and the error handler called.
// Default: 100.
func WithBufferSize(n int) Option {
	return func(b *EventBus) {
		if n > 0 {
			b.bufferSize = n
		}
	}
}

// WithErrorHandler replaces the default error handler (zerolog) with fn.
// fn is called when a listener returns an error or a subscriber channel is full.
func WithErrorHandler(fn func(topic string, err error)) Option {
	return func(b *EventBus) {
		if fn != nil {
			b.errorHandler = fn
		}
	}
}
