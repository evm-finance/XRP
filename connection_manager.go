package xrp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	xrpl "github.com/xrpscan/xrpl-go"
	"go.uber.org/zap"
)

// ConnectionManager manages XRPL WebSocket connections with proper pooling
type ConnectionManager struct {
	mu             sync.RWMutex
	clients        map[string]*xrpl.Client
	activeClients  []*xrpl.Client
	maxConnections int
	connectionURL  string
	logger         *zap.Logger
	reconnectDelay time.Duration
	pingInterval   time.Duration
	connectionTTL  time.Duration
	lastUsed       map[string]time.Time
	ctx            context.Context
	cancel         context.CancelFunc
	// Track connections that are currently in use
	inUseConnections map[*xrpl.Client]bool
	// Track connection creation time for TTL
	connectionCreated map[*xrpl.Client]time.Time
	// Circuit breaker for rate limiting
	rateLimitedUntil time.Time
	rateLimitMutex   sync.RWMutex
	// Per-request throttling removed
}

// ConnectionPool represents a pool of XRPL connections
type ConnectionPool struct {
	manager *ConnectionManager
}

// XRPLListener manages the persistent XRPL listener
type XRPLListener struct {
	connectionManager *ConnectionManager
	logger            *zap.Logger
	ctx               context.Context
	cancel            context.CancelFunc
	ledgerChan        chan []byte
	isRunning         bool
	mu                sync.RWMutex
}

// NewConnectionManager creates a new XRPL connection manager
func NewConnectionManager(connectionURL string, maxConnections int, logger *zap.Logger) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = zap.NewNop()
	}

	manager := &ConnectionManager{
		clients:           make(map[string]*xrpl.Client),
		activeClients:     make([]*xrpl.Client, 0, maxConnections),
		maxConnections:    maxConnections,
		connectionURL:     connectionURL,
		logger:            logger.With(zap.String("component", "connection_manager")),
		reconnectDelay:    time.Second * 5,
		pingInterval:      time.Second * 30,
		connectionTTL:     time.Minute * 10,
		lastUsed:          make(map[string]time.Time),
		inUseConnections:  make(map[*xrpl.Client]bool),
		connectionCreated: make(map[*xrpl.Client]time.Time),
		ctx:               ctx,
		cancel:            cancel,
	}

	// Start background cleanup routine
	go manager.cleanupRoutine()

	return manager
}

// GetConnection returns a healthy connection from the pool
func (cm *ConnectionManager) GetConnection() (*xrpl.Client, error) {
	// Check if we're rate limited BEFORE acquiring the lock
	cm.rateLimitMutex.RLock()
	if time.Now().Before(cm.rateLimitedUntil) {
		backoffRemaining := time.Until(cm.rateLimitedUntil)
		cm.rateLimitMutex.RUnlock()
		return nil, fmt.Errorf("XRPL rate limited, waiting %v before retry", backoffRemaining.Round(time.Second))
	}
	cm.rateLimitMutex.RUnlock()

	// Per-request throttling removed to improve responsiveness

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Log current pool status for debugging
	availableConnections := 0
	for _, client := range cm.activeClients {
		if !cm.inUseConnections[client] {
			availableConnections++
		}
	}

	cm.logger.Debug("Connection pool status",
		zap.Int("total_connections", len(cm.activeClients)),
		zap.Int("available_connections", availableConnections),
		zap.Int("in_use_connections", len(cm.activeClients)-availableConnections),
		zap.Int("max_connections", cm.maxConnections))

	// Try to get an existing healthy connection that's not in use
	for _, client := range cm.activeClients {
		if !cm.inUseConnections[client] && cm.isConnectionHealthy(client) {
			cm.inUseConnections[client] = true
			cm.updateLastUsed(client)
			cm.logger.Debug("Reusing existing connection",
				zap.String("client", fmt.Sprintf("%p", client)),
				zap.Int("total_connections", len(cm.activeClients)),
				zap.Int("available_connections", availableConnections-1))
			return client, nil
		}
	}

	// Create new connection if we haven't reached the limit
	if len(cm.activeClients) < cm.maxConnections {
		client, err := cm.createNewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to create new connection: %w", err)
		}

		cm.activeClients = append(cm.activeClients, client)
		cm.inUseConnections[client] = true
		cm.connectionCreated[client] = time.Now()
		cm.updateLastUsed(client)

		cm.logger.Info("Created new XRPL connection",
			zap.String("url", cm.connectionURL),
			zap.String("client", fmt.Sprintf("%p", client)),
			zap.Int("total_connections", len(cm.activeClients)))

		return client, nil
	}

	// If we've reached the limit, log detailed pool status
	cm.logger.Warn("Connection pool full, waiting for available connection",
		zap.Int("max_connections", cm.maxConnections),
		zap.Int("total_connections", len(cm.activeClients)),
		zap.Int("available_connections", availableConnections),
		zap.Int("in_use_connections", len(cm.activeClients)-availableConnections))

	return nil, fmt.Errorf("connection pool full, try again later")
}

// ReturnConnection returns a connection to the pool
func (cm *ConnectionManager) ReturnConnection(client *xrpl.Client) {
	if client == nil {
		return
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Mark connection as available
	cm.inUseConnections[client] = false
	cm.updateLastUsed(client)

	cm.logger.Debug("Returned connection to pool",
		zap.String("client", fmt.Sprintf("%p", client)))
}

// createNewConnection creates a new XRPL WebSocket connection
func (cm *ConnectionManager) createNewConnection() (*xrpl.Client, error) {
	client := xrpl.NewClient(xrpl.ClientConfig{
		URL: cm.connectionURL,
	})

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- client.Ping([]byte("PING"))
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("connection ping failed: %w", err)
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("connection timeout: %w", ctx.Err())
	}

	return client, nil
}

// isConnectionHealthy checks if a connection is still healthy
func (cm *ConnectionManager) isConnectionHealthy(client *xrpl.Client) bool {
	if client == nil {
		return false
	}

	// Check if connection is too old
	if created, exists := cm.connectionCreated[client]; exists {
		if time.Since(created) > cm.connectionTTL {
			cm.logger.Debug("Connection expired",
				zap.String("client", fmt.Sprintf("%p", client)),
				zap.Duration("age", time.Since(created)))
			return false
		}
	}

	// Quick ping test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- client.Ping([]byte("PING"))
	}()

	select {
	case err := <-done:
		return err == nil
	case <-ctx.Done():
		cm.logger.Debug("Connection health check timeout",
			zap.String("client", fmt.Sprintf("%p", client)))
		return false
	}
}

// updateLastUsed updates the last used timestamp for a connection
func (cm *ConnectionManager) updateLastUsed(client *xrpl.Client) {
	clientID := fmt.Sprintf("%p", client)
	cm.lastUsed[clientID] = time.Now()
}

// getOldestConnection returns the oldest connection that's not in use
func (cm *ConnectionManager) getOldestConnection() *xrpl.Client {
	if len(cm.activeClients) == 0 {
		return nil
	}

	var oldestClient *xrpl.Client
	var oldestTime time.Time
	first := true

	for _, client := range cm.activeClients {
		// Skip connections that are currently in use
		if cm.inUseConnections[client] {
			continue
		}

		clientID := fmt.Sprintf("%p", client)
		if lastUsed, exists := cm.lastUsed[clientID]; exists {
			if first || lastUsed.Before(oldestTime) {
				oldestClient = client
				oldestTime = lastUsed
				first = false
			}
		}
	}

	return oldestClient
}

// MarkConnectionUnhealthy marks a connection as unhealthy for replacement
func (cm *ConnectionManager) MarkConnectionUnhealthy(client *xrpl.Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Find and remove the unhealthy connection
	for i, activeClient := range cm.activeClients {
		if activeClient == client {
			// Remove from active clients
			cm.activeClients = append(cm.activeClients[:i], cm.activeClients[i+1:]...)

			// Clean up tracking maps
			clientID := fmt.Sprintf("%p", client)
			delete(cm.lastUsed, clientID)
			delete(cm.inUseConnections, client)
			delete(cm.connectionCreated, client)

			cm.logger.Info("Removed unhealthy connection",
				zap.String("client", fmt.Sprintf("%p", client)),
				zap.Int("remaining_connections", len(cm.activeClients)))

			// Try to close the connection gracefully
			if err := client.Close(); err != nil {
				cm.logger.Warn("Failed to close unhealthy connection",
					zap.String("client", fmt.Sprintf("%p", client)),
					zap.Error(err))
			}
			break
		}
	}
}

// replaceConnection replaces an old connection with a new one
func (cm *ConnectionManager) replaceConnection(oldClient, newClient *xrpl.Client) {
	// Find and replace the old connection
	for i, client := range cm.activeClients {
		if client == oldClient {
			cm.activeClients[i] = newClient

			// Transfer tracking data
			oldClientID := fmt.Sprintf("%p", oldClient)
			newClientID := fmt.Sprintf("%p", newClient)

			if lastUsed, exists := cm.lastUsed[oldClientID]; exists {
				cm.lastUsed[newClientID] = lastUsed
				delete(cm.lastUsed, oldClientID)
			}

			if inUse, exists := cm.inUseConnections[oldClient]; exists {
				cm.inUseConnections[newClient] = inUse
				delete(cm.inUseConnections, oldClient)
			}

			cm.connectionCreated[newClient] = time.Now()
			delete(cm.connectionCreated, oldClient)

			// Close old connection
			if err := oldClient.Close(); err != nil {
				cm.logger.Warn("Failed to close old connection during replacement",
					zap.String("client", fmt.Sprintf("%p", oldClient)),
					zap.Error(err))
			}

			cm.logger.Info("Replaced connection",
				zap.String("old_client", fmt.Sprintf("%p", oldClient)),
				zap.String("new_client", fmt.Sprintf("%p", newClient)))
			break
		}
	}
}

// GetConnectionPoolStatus returns the current status of the connection pool
func (cm *ConnectionManager) GetConnectionPoolStatus() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	availableConnections := 0
	inUseConnections := 0
	expiredConnections := 0

	for _, client := range cm.activeClients {
		if !cm.inUseConnections[client] {
			availableConnections++
		} else {
			inUseConnections++
		}

		// Check for expired connections
		if created, exists := cm.connectionCreated[client]; exists {
			if time.Since(created) > cm.connectionTTL {
				expiredConnections++
			}
		}
	}

	return map[string]interface{}{
		"total_connections":     len(cm.activeClients),
		"available_connections": availableConnections,
		"in_use_connections":    inUseConnections,
		"expired_connections":   expiredConnections,
		"max_connections":       cm.maxConnections,
		"connection_url":        cm.connectionURL,
	}
}

// cleanupRoutine periodically cleans up expired and unhealthy connections
func (cm *ConnectionManager) cleanupRoutine() {
	ticker := time.NewTicker(time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			cm.cleanupConnections()
		}
	}
}

// cleanupConnections removes expired and unhealthy connections
func (cm *ConnectionManager) cleanupConnections() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	initialCount := len(cm.activeClients)
	removedCount := 0

	// Create a new slice to hold healthy connections
	var healthyClients []*xrpl.Client
	var healthyInUse map[*xrpl.Client]bool
	var healthyCreated map[*xrpl.Client]time.Time

	healthyInUse = make(map[*xrpl.Client]bool)
	healthyCreated = make(map[*xrpl.Client]time.Time)

	for _, client := range cm.activeClients {
		shouldRemove := false

		// Check if connection is expired
		if created, exists := cm.connectionCreated[client]; exists {
			if time.Since(created) > cm.connectionTTL {
				shouldRemove = true
				cm.logger.Debug("Removing expired connection",
					zap.String("client", fmt.Sprintf("%p", client)),
					zap.Duration("age", time.Since(created)))
			}
		}

		// Check if connection is healthy
		if !shouldRemove && !cm.isConnectionHealthy(client) {
			shouldRemove = true
			cm.logger.Debug("Removing unhealthy connection",
				zap.String("client", fmt.Sprintf("%p", client)))
		}

		// Check if connection has been in use for too long (potential leak)
		if !shouldRemove && cm.inUseConnections[client] {
			clientID := fmt.Sprintf("%p", client)
			if lastUsed, exists := cm.lastUsed[clientID]; exists {
				if time.Since(lastUsed) > time.Minute*5 { // 5 minutes timeout
					shouldRemove = true
					cm.logger.Warn("Removing connection that has been in use too long (potential leak)",
						zap.String("client", fmt.Sprintf("%p", client)),
						zap.Duration("in_use_duration", time.Since(lastUsed)))
				}
			}
		}

		if shouldRemove {
			removedCount++
			// Close the client connection
			if client != nil {
				_ = client.Close()
			}
		} else {
			healthyClients = append(healthyClients, client)
			healthyInUse[client] = cm.inUseConnections[client]
			if created, exists := cm.connectionCreated[client]; exists {
				healthyCreated[client] = created
			}
		}
	}

	// Update the connection manager with healthy connections
	cm.activeClients = healthyClients
	cm.inUseConnections = healthyInUse
	cm.connectionCreated = healthyCreated

	if removedCount > 0 {
		cm.logger.Info("Connection cleanup completed",
			zap.Int("initial_connections", initialCount),
			zap.Int("removed_connections", removedCount),
			zap.Int("remaining_connections", len(cm.activeClients)))
	}
}

// Close closes all connections and stops the manager
func (cm *ConnectionManager) Close() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cancel()

	// Close all connections
	for _, client := range cm.activeClients {
		if err := client.Close(); err != nil {
			cm.logger.Warn("Failed to close connection during shutdown",
				zap.String("client", fmt.Sprintf("%p", client)),
				zap.Error(err))
		}
	}

	// Clear all maps
	cm.activeClients = nil
	cm.lastUsed = make(map[string]time.Time)
	cm.inUseConnections = make(map[*xrpl.Client]bool)
	cm.connectionCreated = make(map[*xrpl.Client]time.Time)

	cm.logger.Info("Connection manager closed")
}

// GetStats returns connection manager statistics
func (cm *ConnectionManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	activeCount := len(cm.activeClients)
	inUseCount := 0
	for _, inUse := range cm.inUseConnections {
		if inUse {
			inUseCount++
		}
	}

	return map[string]interface{}{
		"active_connections":    activeCount,
		"max_connections":       cm.maxConnections,
		"in_use_connections":    inUseCount,
		"available_connections": activeCount - inUseCount,
		"connection_url":        cm.connectionURL,
		"reconnect_delay":       cm.reconnectDelay.String(),
		"ping_interval":         cm.pingInterval.String(),
		"connection_ttl":        cm.connectionTTL.String(),
	}
}

// GetPoolStatus returns detailed connection pool status for health checks
func (cm *ConnectionManager) GetPoolStatus() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	activeCount := len(cm.activeClients)
	inUseCount := 0
	healthyCount := 0

	for _, client := range cm.activeClients {
		if cm.inUseConnections[client] {
			inUseCount++
		}
		if cm.isConnectionHealthy(client) {
			healthyCount++
		}
	}

	availableCount := activeCount - inUseCount
	poolHealth := "healthy"
	if availableCount == 0 {
		poolHealth = "exhausted"
	} else if availableCount < 2 {
		poolHealth = "low"
	}

	return map[string]interface{}{
		"active_connections":    activeCount,
		"max_connections":       cm.maxConnections,
		"in_use_connections":    inUseCount,
		"available_connections": availableCount,
		"healthy_connections":   healthyCount,
		"pool_health":           poolHealth,
		"connection_url":        cm.connectionURL,
		"utilization_percent":   float64(inUseCount) / float64(cm.maxConnections) * 100,
	}
}

// NewXRPLListener creates a new XRPL listener with persistent connection
func NewXRPLListener(connectionManager *ConnectionManager, logger *zap.Logger) *XRPLListener {
	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = zap.NewNop()
	}

	return &XRPLListener{
		connectionManager: connectionManager,
		logger:            logger.With(zap.String("component", "xrpl_listener")),
		ctx:               ctx,
		cancel:            cancel,
		ledgerChan:        make(chan []byte, 100), // Buffered channel for ledger data
		isRunning:         false,
	}
}

// StartListener starts the persistent XRPL listener based on EVM legacy pattern
func (xl *XRPLListener) StartListener() {
	xl.mu.Lock()
	if xl.isRunning {
		xl.mu.Unlock()
		return
	}
	xl.isRunning = true
	xl.mu.Unlock()

	xl.logger.Info("Starting XRPL listener")

	for {
		select {
		case <-xl.ctx.Done():
			xl.logger.Info("XRPL listener context cancelled")
			return
		default:
			xl.runListener()
		}
	}
}

// runListener runs a single listener session with automatic reconnection
func (xl *XRPLListener) runListener() {
	// Get connection from pool
	client, err := xl.connectionManager.GetConnection()
	if err != nil {
		xl.logger.Error("Failed to get XRPL connection", zap.Error(err))
		time.Sleep(xl.connectionManager.reconnectDelay)
		return
	}

	// Subscribe to ledger stream
	_, err = client.Subscribe([]string{xrpl.StreamTypeLedger})
	if err != nil {
		xl.logger.Error("Failed to subscribe to ledger stream", zap.Error(err))
		time.Sleep(xl.connectionManager.reconnectDelay)
		return
	}

	xl.logger.Info("Successfully subscribed to XRPL ledger stream")

	// Listen for ledger updates
	for {
		select {
		case <-xl.ctx.Done():
			return
		case ledgerData := <-client.StreamLedger:
			xl.processLedgerData(ledgerData)
		case <-time.After(time.Minute * 5):
			// Timeout - reconnect
			xl.logger.Warn("XRPL listener timeout, reconnecting")
			return
		}
	}
}

// processLedgerData processes incoming ledger data
func (xl *XRPLListener) processLedgerData(ledgerData []byte) {
	var ledgerInfo map[string]interface{}
	if err := json.Unmarshal(ledgerData, &ledgerInfo); err != nil {
		xl.logger.Error("Failed to unmarshal ledger data", zap.Error(err))
		return
	}

	if ledgerIndex, ok := ledgerInfo["ledger_index"]; ok {
		xl.logger.Info("Received XRPL ledger",
			zap.Any("ledger_index", ledgerIndex))
	}

	// Send to processing channel (non-blocking)
	select {
	case xl.ledgerChan <- ledgerData:
	default:
		xl.logger.Warn("Ledger processing channel full, dropping data")
	}
}

// GetLedgerChannel returns the channel for ledger data
func (xl *XRPLListener) GetLedgerChannel() <-chan []byte {
	return xl.ledgerChan
}

// Stop stops the XRPL listener
func (xl *XRPLListener) Stop() {
	xl.mu.Lock()
	defer xl.mu.Unlock()

	if !xl.isRunning {
		return
	}

	xl.cancel()
	xl.isRunning = false
	close(xl.ledgerChan)

	xl.logger.Info("XRPL listener stopped")
}

// TriggerRateLimitBackoff sets the rate limit backoff period
func (cm *ConnectionManager) TriggerRateLimitBackoff(duration time.Duration) {
	cm.rateLimitMutex.Lock()
	defer cm.rateLimitMutex.Unlock()

	cm.rateLimitedUntil = time.Now().Add(duration)
	cm.logger.Warn("🛑 XRPL rate limit detected - pausing all requests",
		zap.Duration("backoff_duration", duration),
		zap.Time("resume_at", cm.rateLimitedUntil))
}

// IsRunning returns whether the listener is currently running
func (xl *XRPLListener) IsRunning() bool {
	xl.mu.RLock()
	defer xl.mu.RUnlock()
	return xl.isRunning
}
