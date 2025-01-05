package main

import (
	"fmt"
	"github.com/macurandb/go-eventbus-pattern/eventbus"
	"time"
)

func handleUserCreated(event eventbus.Event) {
	fmt.Printf("Processing event: %v\n", event.Data)
	time.Sleep(500 * time.Millisecond) // Simulate some processing
	fmt.Println("Finished processing")
}

func main() {
	bus := eventbus.NewEventBus()

	// Subscribe a predefined function to the "user_created" event
	bus.Subscribe("user_created", handleUserCreated)

	// Create and publish a "user_created" event
	event := eventbus.Event{
		Type:      "user_created",
		Timestamp: time.Now(),
		Source:    "auth_service",
		Data: map[string]interface{}{
			"user_id": "12345",
			"email":   "user@example.com",
		},
	}

	fmt.Println("Publishing event...")
	bus.Publish("user_created", event)

	fmt.Println("All subscribers have finished processing")
}
