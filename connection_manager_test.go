package xrp

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestConnectionManager(t *testing.T) {
	logger := zap.NewNop()

	// Test creating connection manager
	cm := NewConnectionManager("wss://xrplcluster.com", 3, logger)
	defer cm.Close()

	if cm == nil {
		t.Fatal("Connection manager should not be nil")
	}

	if cm.maxConnections != 3 {
		t.Fatalf("Expected max connections 3, got %d", cm.maxConnections)
	}

	if cm.connectionURL != "wss://xrplcluster.com" {
		t.Fatalf("Expected URL wss://xrplcluster.com, got %s", cm.connectionURL)
	}

	// Test getting connection (this may fail if XRPL is not available, but we can test the logic)
	client, err := cm.GetConnection()
	if err != nil {
		t.Logf("Expected error when XRPL is not available: %v", err)
		// This is expected in test environment
		return
	}

	if client == nil {
		t.Fatal("Client should not be nil when connection succeeds")
	}

	// Test connection reuse
	client2, err := cm.GetConnection()
	if err != nil {
		t.Fatalf("Second connection should reuse existing: %v", err)
	}

	if client2 != client {
		t.Log("Second connection may be different (pool behavior)")
	}
}

func TestXRPLListener(t *testing.T) {
	logger := zap.NewNop()
	cm := NewConnectionManager("wss://xrplcluster.com", 2, logger)
	defer cm.Close()

	listener := NewXRPLListener(cm, logger)

	if listener == nil {
		t.Fatal("Listener should not be nil")
	}

	if listener.IsRunning() {
		t.Fatal("Listener should not be running initially")
	}

	// Test starting and stopping listener
	go listener.StartListener()
	time.Sleep(time.Millisecond * 100) // Give it time to start

	if !listener.IsRunning() {
		t.Log("Listener may not be running due to connection issues (expected in test)")
	}

	listener.Stop()
	time.Sleep(time.Millisecond * 100) // Give it time to stop

	if listener.IsRunning() {
		t.Fatal("Listener should be stopped")
	}
}

func TestNewXRPService(t *testing.T) {
	logger := zap.NewNop()

	service := NewXRPService("wss://xrplcluster.com", 5, logger)

	if service == nil {
		t.Fatal("XRP service should not be nil")
	}

	if service.connectionManager == nil {
		t.Fatal("Connection manager should be initialized")
	}

	// Test getting client (may fail in test environment)
	_, err := service.GetClient()
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
		// This is expected when XRPL is not available
	}
}
