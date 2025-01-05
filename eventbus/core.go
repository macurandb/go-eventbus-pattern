package eventbus

import (
	"sync"
	"time"
)

type Event struct {
	Type      string      // The type of the event
	Timestamp time.Time   // Timestamp of the event
	Source    string      // Source of the event (e.g., system, service)
	Data      interface{} // Data associated with the event
}

// EventBus represents a generic event bus.
type EventBus struct {
	subscribers map[EventType]map[int]Subscriber
	mutex       sync.Mutex
	counter     int // Unique counter for subscriber IDs
}

// EventType is the type representing different event types.
type EventType string

// Subscriber is a function that handles events.
type Subscriber func(event Event)

// NewEventBus creates a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType]map[int]Subscriber),
		counter:     0,
	}
}

// Subscribe adds a subscriber for a specific event type and returns a unique ID.
func (eb *EventBus) Subscribe(eventType EventType, subscriber Subscriber) int {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make(map[int]Subscriber)
	}

	eb.counter++
	id := eb.counter
	eb.subscribers[eventType][id] = subscriber

	return id
}

// Unsubscribe removes a subscriber for a specific event type using its ID.
func (eb *EventBus) Unsubscribe(eventType EventType, id int) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		delete(subscribers, id)

		// Cleanup the event type map if there are no subscribers left.
		if len(subscribers) == 0 {
			delete(eb.subscribers, eventType)
		}
	}
}

// Publish sends an event to all subscribers of a specific event type.
func (eb *EventBus) Publish(eventType EventType, event Event) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	var wg sync.WaitGroup

	if subscribers, ok := eb.subscribers[eventType]; ok {
		for _, subscriber := range subscribers {
			wg.Add(1)
			go func(sub Subscriber) {
				defer wg.Done()
				sub(event)
			}(subscriber)
		}
	}

	wg.Wait() // Wait for all goroutines to finish
}
