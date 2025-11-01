package test

import (
	"testing"
	"time"

	"github.com/glance-project/glance/internal/websocket"
)

func TestHubRegister(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	client := &websocket.Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	if len(hub.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(hub.Clients))
	}

	hub.Stop()
}

func TestHubUnregister(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	client := &websocket.Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client
	time.Sleep(100 * time.Millisecond)

	hub.Unregister <- client
	time.Sleep(100 * time.Millisecond)

	if len(hub.Clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(hub.Clients))
	}

	hub.Stop()
}

func TestHubBroadcast(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	client := &websocket.Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client
	time.Sleep(100 * time.Millisecond)

	message := []byte("test message")
	hub.Broadcast <- message

	time.Sleep(100 * time.Millisecond)

	select {
	case msg := <-client.Send:
		if string(msg) != string(message) {
			t.Errorf("Expected %s, got %s", message, msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("Expected message on client send channel")
	}

	hub.Stop()
}

func TestMultipleClients(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	clients := make([]*websocket.Client, 3)
	for i := 0; i < 3; i++ {
		clients[i] = &websocket.Client{
			Hub:  hub,
			Send: make(chan []byte, 256),
		}
		hub.Register <- clients[i]
	}

	time.Sleep(100 * time.Millisecond)

	if len(hub.Clients) != 3 {
		t.Errorf("Expected 3 clients, got %d", len(hub.Clients))
	}

	message := []byte("broadcast to all")
	hub.Broadcast <- message

	time.Sleep(100 * time.Millisecond)

	for _, client := range clients {
		select {
		case msg := <-client.Send:
			if string(msg) != string(message) {
				t.Errorf("Expected %s, got %s", message, msg)
			}
		case <-time.After(1 * time.Second):
			t.Error("Expected message on client send channel")
		}
	}

	hub.Stop()
}
