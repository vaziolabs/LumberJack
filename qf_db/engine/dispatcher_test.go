package engine

import (
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	// Handle OS interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\nExiting...")
			os.Exit(0)
		}
	}()

	// Create 4 instances of Dispatcher
	var handlers []*Dispatcher
	for i := 0; i < 4; i++ {
		port := fmt.Sprintf("%d", portBase+i)
		d := NewDispatcher(port)
		handlers = append(handlers, d)

		go d.StartServer()
	}

	// Wait for servers to start
	time.Sleep(time.Second)

	// Connect each handler to others as clients
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i != j {
				serverAddr := "localhost:" + handlers[j].port
				go handlers[i].StartClient(serverAddr)
			}
		}
	}

	// Example sending messages (for demonstration purposes)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		timestamp := time.Now().UnixNano()
		frequency := calculateFrequency(i + 1)
		message := fmt.Sprintf("Message %d (%d): %f", i+1, timestamp, frequency)
		handlers[i%4].sendMessage(message)
	}

	// Keep the main goroutine running indefinitely to handle signal
	select {}
}
