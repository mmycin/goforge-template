// Package events provides a production-ready, in-process event bus for GoForge.
//
// Usage:
//
//	// register once at startup
//	events.On(bus, todo.EventTodoCreated, &todo.AuditOnCreated{DB: db})
//
//	// fire from a service
//	events.Dispatch(ctx, bus, todo.EventTodoCreated, TodoCreated{ID: t.ID, Title: t.Title})
package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// subscription holds the channel and goroutine for a single registered listener.
// It is type-erased at this layer; the generic On/Dispatch functions handle typing.
type subscription struct {
	topic  string
	ch     chan envelope
	cancel context.CancelFunc
	wg     *sync.WaitGroup // points to the bus-level WaitGroup
}

// envelope carries a single event delivery to a subscriber goroutine.
type envelope struct {
	ctx     context.Context
	payload any
	// done is non-nil only for DispatchSync deliveries.
	// The subscriber closes it after the listener returns.
	done chan error
}

// EventBus is the central event dispatcher.
// Create one with New() and inject it wherever events need to be fired.
// Register listeners with On() before calling Dispatch().
type EventBus struct {
	mu           sync.RWMutex
	subs         map[string][]*subscription
	wg           sync.WaitGroup // tracks all in-flight handler calls
	bufferSize   int
	errorHandler func(topic string, err error)
	closed       bool
}

// New creates a ready-to-use EventBus.
func New(opts ...Option) *EventBus {
	b := &EventBus{
		subs:       make(map[string][]*subscription),
		bufferSize: 100,
		errorHandler: func(topic string, err error) {
			log.Error().Str("topic", topic).Err(err).Msg("event listener error")
		},
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

// on is the unexported core of On[T]. It registers a type-erased handler function
// for the given topic and starts its dedicated goroutine.
func (b *EventBus) on(topic string, handle func(context.Context, any) error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	sub := &subscription{
		topic:  topic,
		ch:     make(chan envelope, b.bufferSize),
		cancel: cancel,
		wg:     &b.wg,
	}

	b.subs[topic] = append(b.subs[topic], sub)

	// Each subscriber has its own goroutine for the lifetime of the bus.
	go func() {
		for {
			select {
			case env, ok := <-sub.ch:
				if !ok {
					// channel closed by Shutdown
					return
				}
				err := handle(env.ctx, env.payload)
				if err != nil {
					b.errorHandler(topic, err)
				}
				// Signal DispatchSync waiters if this was a sync delivery.
				if env.done != nil {
					env.done <- err
				}
				b.wg.Done()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Shutdown stops accepting new events and waits for all in-flight listener
// goroutines to finish, or until ctx is cancelled.
//
// Called by boot/kernel.go during graceful shutdown. The user never calls this.
func (b *EventBus) Shutdown(ctx context.Context) error {
	b.mu.Lock()
	b.closed = true
	// Close all subscriber channels — their goroutines will drain and exit.
	for _, subs := range b.subs {
		for _, sub := range subs {
			sub.cancel()
			close(sub.ch)
		}
	}
	b.mu.Unlock()

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("events: shutdown timed out: %w", ctx.Err())
	}
}

// On registers listener for the given topic.
// Call this once at startup (in internal/services/kernel.go).
// T must match the event type passed to Dispatch for this topic.
func On[T any](bus *EventBus, topic string, listener Listener[T]) {
	bus.on(topic, func(ctx context.Context, payload any) error {
		return listener.Handle(ctx, payload.(T))
	})
}
