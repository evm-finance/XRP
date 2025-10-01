package xrp

import (
	"fmt"
	"log"
	"time"

	"qc-defi-graphql-server/internal/database"
)

// XRPDataCollectorService orchestrates the collection and storage of XRP data
type XRPDataCollectorService struct {
	tokenService XRPTokenServiceInterface
	ammService   AMMServiceInterface
	dbService    database.DatabaseCreate
}

type XRPDataCollectorInterface interface {
	// Token data collection
	CollectAndStoreTokens() error
	UpdateTokenMetrics() error

	// AMM data collection
	CollectAndStoreAMMPools() error
	UpdateAMMPoolMetrics() error

	// Combined operations
	CollectAllData() error
	GetCollectionStatus() (*CollectionStatus, error)
}

type CollectionStatus struct {
	LastTokenUpdate time.Time `json:"last_token_update"`
	LastAMMUpdate   time.Time `json:"last_amm_update"`
	TokenCount      int       `json:"token_count"`
	AMMPoolCount    int       `json:"amm_pool_count"`
	LastError       string    `json:"last_error,omitempty"`
	IsCollecting    bool      `json:"is_collecting"`
}

func NewXRPDataCollectorService(
	tokenService XRPTokenServiceInterface,
	ammService AMMServiceInterface,
	dbService database.DatabaseCreate,
) XRPDataCollectorInterface {
	return &XRPDataCollectorService{
		tokenService: tokenService,
		ammService:   ammService,
		dbService:    dbService,
	}
}

// CollectAndStoreTokens fetches and stores token data from XRPL Meta API
func (s *XRPDataCollectorService) CollectAndStoreTokens() error {
	log.Println("Starting XRP token data collection...")

	// Fetch tokens from API using the corrected single-call method
	tokens, err := s.tokenService.FetchAndStoreAllTokens()
	if err != nil {
		return fmt.Errorf("failed to fetch and store tokens: %w", err)
	}

	log.Printf("Finished token processing. Fetched and processed %d tokens.", len(tokens))

	return nil
}

// UpdateTokenMetrics updates existing token metrics
func (s *XRPDataCollectorService) UpdateTokenMetrics() error {
	log.Println("Updating XRP token metrics...")

	if err := s.tokenService.UpdateTokenMetrics(); err != nil {
		return fmt.Errorf("failed to update token metrics: %w", err)
	}

	log.Println("Successfully updated XRP token metrics")
	return nil
}

// CollectAndStoreAMMPools fetches and stores AMM pool data
func (s *XRPDataCollectorService) CollectAndStoreAMMPools() error {
	log.Println("Starting XRP AMM pool data collection...")

	// First, get tokens that might have pools
	tokens, err := s.tokenService.GetTokensWithPools()
	if err != nil {
		return fmt.Errorf("failed to get tokens for AMM collection: %w", err)
	}

	log.Printf("Found %d tokens to check for AMM pools", len(tokens))

	// Discover all AMM pools using the discovery method
	log.Println("Discovering AMM pools from XRPL...")
	allPools, err := s.ammService.DiscoverAllAMMPools()
	if err != nil {
		return fmt.Errorf("failed to discover AMM pools: %w", err)
	}

	log.Printf("Found %d AMM pools", len(allPools))

	// Store AMM pools in database using normalized table
	if err := s.ammService.StoreAMMPoolsNormalized(allPools); err != nil {
		return fmt.Errorf("failed to store AMM pools: %w", err)
	}

	log.Printf("Successfully stored %d AMM pools in XRPDB", len(allPools))
	return nil
}

// UpdateAMMPoolMetrics updates existing AMM pool metrics
func (s *XRPDataCollectorService) UpdateAMMPoolMetrics() error {
	log.Println("Updating XRP AMM pool metrics...")

	if err := s.ammService.UpdatePoolMetrics(); err != nil {
		return fmt.Errorf("failed to update AMM pool metrics: %w", err)
	}

	log.Println("Successfully updated XRP AMM pool metrics")
	return nil
}

// CollectAllData performs a complete data collection run
func (s *XRPDataCollectorService) CollectAllData() error {
	log.Println("Starting complete XRP data collection...")

	// Step 1: Collect and store tokens
	if err := s.CollectAndStoreTokens(); err != nil {
		return fmt.Errorf("token collection failed: %w", err)
	}

	// Step 2: Collect and store AMM pools
	if err := s.CollectAndStoreAMMPools(); err != nil {
		return fmt.Errorf("AMM collection failed: %w", err)
	}

	log.Println("Complete XRP data collection finished successfully")
	return nil
}

// GetCollectionStatus returns the current status of data collection
func (s *XRPDataCollectorService) GetCollectionStatus() (*CollectionStatus, error) {
	db, err := s.dbService.GetConnectionXRPDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	status := &CollectionStatus{
		IsCollecting: false, // This would be set by a background job
	}

	// Get token count
	var tokenCount int
	if err := db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&tokenCount).Error; err != nil {
		log.Printf("Warning: failed to get token count: %v", err)
	} else {
		status.TokenCount = tokenCount
	}

	// Get AMM pool count
	var ammPoolCount int
	if err := db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized").Scan(&ammPoolCount).Error; err != nil {
		log.Printf("Warning: failed to get AMM pool count: %v", err)
	} else {
		status.AMMPoolCount = ammPoolCount
	}

	// Get last update times (you might want to add timestamp columns to track this)
	// For now, we'll use current time as placeholder
	status.LastTokenUpdate = time.Now()
	status.LastAMMUpdate = time.Now()

	return status, nil
}
