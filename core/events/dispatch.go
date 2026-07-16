package events

import (
	"context"
	"errors"
)

// Dispatch delivers event to all listeners registered for topic.
// Returns immediately — listeners run concurrently in their own goroutines.
// Listener errors are routed to the bus error handler (default: zerolog).
//
// Use this for fire-and-forget side effects (audit logs, notifications, etc.).
func Dispatch[T any](ctx context.Context, bus *EventBus, topic string, event T) {
	if bus == nil {
		return
	}
	// ctx is stored on the envelope so listeners receive the caller's context.
	bus.mu.RLock()
	if bus.closed {
		bus.mu.RUnlock()
		return
	}
	subs := bus.subs[topic]
	bus.mu.RUnlock()

	for _, sub := range subs {
		env := envelope{
			ctx:     ctx,
			payload: event,
		}
		bus.wg.Add(1)
		select {
		case sub.ch <- env:
		default:
			bus.wg.Done()
			bus.errorHandler(topic, errors.New("events: subscriber channel full, event dropped"))
		}
	}
}

// DispatchSync delivers event and blocks until all listeners have returned.
// Returns a combined error if any listener failed.
//
// Use this when side-effects must complete before the caller continues —
// e.g. inside a DB transaction, or in tests.
func DispatchSync[T any](ctx context.Context, bus *EventBus, topic string, event T) error {
	if bus == nil {
		return nil
	}

	bus.mu.RLock()
	if bus.closed {
		bus.mu.RUnlock()
		return errors.New("events: bus is shut down")
	}
	subs := bus.subs[topic]
	bus.mu.RUnlock()

	if len(subs) == 0 {
		return nil
	}

	doneChans := make([]chan error, 0, len(subs))

	for _, sub := range subs {
		done := make(chan error, 1)
		env := envelope{
			ctx:     ctx,
			payload: event,
			done:    done,
		}
		bus.wg.Add(1)
		select {
		case sub.ch <- env:
			doneChans = append(doneChans, done)
		default:
			bus.wg.Done()
			err := errors.New("events: subscriber channel full, event dropped")
			bus.errorHandler(topic, err)
			done <- err
			doneChans = append(doneChans, done)
		}
	}

	// Collect all errors.
	var errs []error
	for _, ch := range doneChans {
		if err := <-ch; err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
