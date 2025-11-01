package websocket

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a test message
	testMessage := Message{
		Type:      "test",
		WidgetID:  "widget-1",
		Data:      "test data",
		Timestamp: time.Now().Unix(),
	}

	// Broadcast message
	hub.Broadcast(testMessage.Type, testMessage.WidgetID, testMessage.Data)

	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHubClientCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients initially, got %d", hub.ClientCount())
	}
}

func TestMessageSerialization(t *testing.T) {
	msg := Message{
		Type:      "update",
		WidgetID:  "widget-1",
		Data:      map[string]interface{}{"key": "value"},
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var unmarshaled Message
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if unmarshaled.Type != msg.Type {
		t.Errorf("Type mismatch: %s != %s", unmarshaled.Type, msg.Type)
	}

	if unmarshaled.WidgetID != msg.WidgetID {
		t.Errorf("WidgetID mismatch: %s != %s", unmarshaled.WidgetID, msg.WidgetID)
	}
}
