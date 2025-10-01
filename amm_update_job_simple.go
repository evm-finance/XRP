package xrp

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// SimpleAMMUpdateJob manages periodic AMM pool updates with discovery
type SimpleAMMUpdateJob struct {
	db                *gorm.DB
	ammService        *AMMService
	priceCollector    *AMMPriceCollector
	updateInterval    time.Duration
	discoveryInterval int // Run discovery every N update cycles
	ticker            *time.Ticker
	stopChan          chan struct{}
	running           bool
	mu                sync.RWMutex

	// Monitoring
	updateCount        int64
	errorCount         int64
	lastUpdateTime     time.Time
	lastUpdateSuccess  bool
	lastDiscoveryTime  time.Time
	lastDiscoveryCount int
}

// AMMUpdateStats provides statistics about the update job
type AMMUpdateStats struct {
	Running           bool          `json:"running"`
	LastUpdateTime    time.Time     `json:"last_update_time"`
	LastUpdateSuccess bool          `json:"last_update_success"`
	UpdateCount       int64         `json:"update_count"`
	ErrorCount        int64         `json:"error_count"`
	UpdateInterval    time.Duration `json:"update_interval"`
	NextUpdateIn      time.Duration `json:"next_update_in"`
}

// NewSimpleAMMUpdateJob creates a new AMM update job with discovery
func NewSimpleAMMUpdateJob(db *gorm.DB, ammService *AMMService, priceCollector *AMMPriceCollector, updateInterval time.Duration) *SimpleAMMUpdateJob {
	return &SimpleAMMUpdateJob{
		db:                db,
		ammService:        ammService,
		priceCollector:    priceCollector,
		updateInterval:    updateInterval,
		discoveryInterval: 10, // Run discovery every 10 update cycles (30 minutes with 3-minute intervals)
		stopChan:          make(chan struct{}),
		running:           false,
	}
}

// Start begins the continuous update process
func (job *SimpleAMMUpdateJob) Start() error {
	job.mu.Lock()
	defer job.mu.Unlock()

	if job.running {
		return fmt.Errorf("AMM update job is already running")
	}

	job.running = true
	go job.runUpdateLoop()

	log.Printf("✅ [AMM UPDATE JOB] Started with %v update interval", job.updateInterval)
	return nil
}

// Stop halts the update process
func (job *SimpleAMMUpdateJob) Stop() error {
	job.mu.Lock()
	defer job.mu.Unlock()

	if !job.running {
		return fmt.Errorf("AMM update job is not running")
	}

	job.stopChan <- struct{}{}
	job.running = false

	log.Printf("🛑 [AMM UPDATE JOB] Stopped after %d updates (%d errors)", job.updateCount, job.errorCount)
	return nil
}

// GetStats returns current job statistics
func (job *SimpleAMMUpdateJob) GetStats() AMMUpdateStats {
	job.mu.Lock()
	defer job.mu.Unlock()

	stats := AMMUpdateStats{
		Running:           job.running,
		LastUpdateTime:    job.lastUpdateTime,
		LastUpdateSuccess: job.lastUpdateSuccess,
		UpdateCount:       job.updateCount,
		ErrorCount:        job.errorCount,
		UpdateInterval:    job.updateInterval,
	}

	if job.running && !job.lastUpdateTime.IsZero() {
		nextUpdate := job.lastUpdateTime.Add(job.updateInterval)
		stats.NextUpdateIn = time.Until(nextUpdate)
	}

	return stats
}

// runUpdateLoop is the main update loop
func (job *SimpleAMMUpdateJob) runUpdateLoop() {
	log.Printf("🔄 [AMM UPDATE JOB] Background loop started - running initial update immediately")

	// Run initial update immediately
	job.performUpdate()

	ticker := time.NewTicker(job.updateInterval)
	defer ticker.Stop()

	// Add a heartbeat ticker for monitoring (every 30 seconds)
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	log.Printf("⏰ [AMM UPDATE JOB] Timer set for %v intervals - waiting for next update cycle", job.updateInterval)

	cycleCount := 1
	for {
		select {
		case <-ticker.C:
			cycleCount++
			log.Printf("⏰ [AMM UPDATE JOB] Timer triggered - starting scheduled update cycle #%d", cycleCount)
			job.performUpdate()
		case <-heartbeatTicker.C:
			// Periodic heartbeat for monitoring
			job.mu.Lock()
			if job.running {
				nextUpdate := job.lastUpdateTime.Add(job.updateInterval)
				timeUntilNext := time.Until(nextUpdate)
				log.Printf("💓 [AMM UPDATE JOB] Heartbeat - Status: Running | Next update in: %v | Total cycles: %d", timeUntilNext, job.updateCount)
			}
			job.mu.Unlock()
		case <-job.stopChan:
			log.Printf("🛑 [AMM UPDATE JOB] Stop signal received - shutting down update loop")
			return
		}
	}
}

// performUpdate executes a single update cycle
func (job *SimpleAMMUpdateJob) performUpdate() {
	startTime := time.Now()
	nextCycleNum := job.updateCount + 1
	log.Printf("🔄 [AMM UPDATE JOB] ============ STARTING UPDATE CYCLE #%d ============", nextCycleNum)
	log.Printf("🕐 [AMM UPDATE JOB] Cycle start time: %s", startTime.Format("2006-01-02 15:04:05 UTC"))

	var updateError error
	defer func() {
		job.mu.Lock()
		job.lastUpdateTime = time.Now()
		job.lastUpdateSuccess = updateError == nil

		if updateError == nil {
			job.updateCount++
			duration := time.Since(startTime)
			log.Printf("✅ [AMM UPDATE JOB] ============ CYCLE #%d COMPLETED ============", job.updateCount)
			log.Printf("⏱️  [AMM UPDATE JOB] Cycle duration: %v", duration)
			log.Printf("📊 [AMM UPDATE JOB] Total successful cycles: %d", job.updateCount)
			log.Printf("⏰ [AMM UPDATE JOB] Next update in: %v", job.updateInterval)

			// Log periodic summary every 10 cycles
			if job.updateCount%10 == 0 {
				log.Printf("📈 [AMM UPDATE JOB] ===== 10-CYCLE SUMMARY =====")
				log.Printf("📈 [AMM UPDATE JOB] Total successful updates: %d", job.updateCount)
				log.Printf("📈 [AMM UPDATE JOB] Total errors: %d", job.errorCount)
				successRate := float64(job.updateCount) / float64(job.updateCount+job.errorCount) * 100
				log.Printf("📈 [AMM UPDATE JOB] Success rate: %.1f%%", successRate)
				log.Printf("📈 [AMM UPDATE JOB] Runtime: %v", time.Since(startTime))
			}
		} else {
			job.errorCount++
			duration := time.Since(startTime)
			log.Printf("❌ [AMM UPDATE JOB] ============ CYCLE #%d FAILED ============", nextCycleNum)
			log.Printf("⏱️  [AMM UPDATE JOB] Failed after: %v", duration)
			log.Printf("💥 [AMM UPDATE JOB] Error: %v", updateError)
			log.Printf("📊 [AMM UPDATE JOB] Total errors: %d", job.errorCount)
			log.Printf("⏰ [AMM UPDATE JOB] Will retry in: %v", job.updateInterval)
		}
		job.mu.Unlock()
	}()

	// Step 0: Periodic AMM Discovery (every N cycles to find new pools)
	if job.updateCount%int64(job.discoveryInterval) == 0 {
		log.Printf("🔍 [AMM UPDATE JOB] Step 0: Running periodic AMM discovery (every %d cycles)...", job.discoveryInterval)
		job.runDiscovery()
	}

	// Step 1: Update pool metrics with fresh XRPL data (reuses complete discovery logic)
	log.Printf("📊 [AMM UPDATE JOB] Step 1: Updating pool metrics with fresh XRPL data...")
	if err := job.ammService.UpdatePoolMetrics(); err != nil {
		updateError = fmt.Errorf("pool metrics update failed: %w", err)
		return
	}
	log.Printf("✅ [AMM UPDATE JOB] Step 1 complete: Pool metrics updated with ALL database fields")

	// Step 2: Collect prices from AMM pools
	log.Printf("💰 [AMM UPDATE JOB] Step 2: Collecting token prices...")
	priceResult, err := job.priceCollector.CollectPricesFromAMMPools()
	if err != nil {
		updateError = fmt.Errorf("price collection failed: %w", err)
		return
	}

	log.Printf("✅ [AMM UPDATE JOB] Collected prices for %d/%d tokens",
		priceResult.PricedTokens,
		priceResult.TotalTokens)

	// Step 3: Store updated prices to database
	log.Printf("💾 [AMM UPDATE JOB] Step 3: Storing prices to database...")
	if err := job.priceCollector.StorePricesToDatabase(priceResult.Prices); err != nil {
		updateError = fmt.Errorf("failed to store prices: %w", err)
		return
	}

	// Step 4: Compute liquidity by XRP amount for better price accuracy
	log.Printf("📈 [AMM UPDATE JOB] Step 4: Computing liquidity values...")
	if err := job.ammService.computeLiquidityByXRPAmount(); err != nil {
		// Non-critical error, log but don't fail the update
		log.Printf("⚠️  [AMM UPDATE JOB] Liquidity computation failed: %v", err)
	}

	// Log summary
	log.Printf("📊 [AMM UPDATE JOB] Update summary:")
	log.Printf("   - Tokens priced: %d", priceResult.PricedTokens)
	log.Printf("   - Failed prices: %d", priceResult.FailedTokens)
	log.Printf("   - Update duration: %v", time.Since(startTime))
}

// ForceUpdate triggers an immediate update outside the regular schedule
func (job *SimpleAMMUpdateJob) ForceUpdate() error {
	job.mu.Lock()
	if !job.running {
		job.mu.Unlock()
		return fmt.Errorf("AMM update job is not running")
	}
	job.mu.Unlock()

	log.Printf("🔄 [AMM UPDATE JOB] Force update requested")
	go job.performUpdate()
	return nil
}

// runDiscovery executes AMM pool discovery to find new pools
func (job *SimpleAMMUpdateJob) runDiscovery() {
	discoveryStartTime := time.Now()

	job.mu.Lock()
	job.lastDiscoveryTime = discoveryStartTime
	job.mu.Unlock()

	log.Printf("🔍 [AMM DISCOVERY] Starting AMM pool discovery...")

	// Run full pool discovery to find new pools
	newPools, err := job.ammService.DiscoverAllAMMPools()
	if err != nil {
		log.Printf("❌ [AMM DISCOVERY] Discovery failed: %v", err)
		return
	}

	poolCount := len(newPools)

	job.mu.Lock()
	job.lastDiscoveryCount = poolCount
	job.mu.Unlock()

	discoveryDuration := time.Since(discoveryStartTime)
	log.Printf("✅ [AMM DISCOVERY] Discovery completed in %v: found %d pools", discoveryDuration, poolCount)

	if poolCount > 0 {
		log.Printf("📊 [AMM DISCOVERY] New pools found and stored in database")
	} else {
		log.Printf("📊 [AMM DISCOVERY] No new pools discovered")
	}
}
