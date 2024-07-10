package engine

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	portBase = 8000
)

// Dispatcher handles incoming and outgoing connections
type Dispatcher struct {
	conn net.Conn
	port string
	mu   sync.Mutex
}

func NewDispatcher(port string) *Dispatcher {
	return &Dispatcher{
		conn: nil,
		port: port,
	}

}

// StartServer starts the dispatcher as a server
func (d *Dispatcher) StartServer() {
	listener, err := net.Listen("tcp", ":"+d.port)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", d.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		d.mu.Lock()
		d.conn = conn
		d.mu.Unlock()

		fmt.Println("Client connected from", conn.RemoteAddr())

		// Handle incoming messages from client
		go d.handleIncoming()
	}
}

// StartClient starts the dispatcher as a client
func (d *Dispatcher) StartClient(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	d.mu.Lock()
	d.conn = conn
	d.mu.Unlock()

	fmt.Println("Connected to server at", serverAddr)

	// Handle incoming messages from server
	go d.handleIncoming()
}

// handleIncoming handles incoming messages
func (d *Dispatcher) handleIncoming() {
	for {
		buffer := make([]byte, 1024)
		n, err := d.conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		receivedMsg := string(buffer[:n])
		fmt.Printf("Received: %s\n", receivedMsg)

		// Parse the message to extract the timestamp and frequency
		parts := strings.Split(receivedMsg, ",")
		if len(parts) != 3 {
			fmt.Println("Error parsing message: incorrect format")
			return
		}

		messageID, err := strconv.Atoi(strings.TrimPrefix(parts[0], "Message "))
		if err != nil {
			fmt.Println("Error parsing message ID:", err)
			return
		}

		receivedTimestamp, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return
		}

		frequency, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			fmt.Println("Error parsing frequency:", err)
			return
		}

		// Calculate latency
		currentTimestamp := time.Now().UnixNano()
		latency := float64(currentTimestamp-receivedTimestamp) / 1e6 // in milliseconds
		fmt.Printf("Message %d, Latency: %.2f ms, Frequency: %.2f\n", messageID, latency, frequency)

		// Send acknowledgment
		d.sendMessage(fmt.Sprintf("ACK,%d", currentTimestamp))
	}
}

// sendMessage sends a message
func (d *Dispatcher) sendMessage(message string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending:", err)
	}
}

// calculateFrequency simulates frequency calculation based on the message index
func calculateFrequency(index int) float64 {
	// Placeholder calculation based on message index
	return float64(index) * 1.5
}
