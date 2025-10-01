package xrp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"qc-defi-graphql-server/internal/database"
	"qc-defi-graphql-server/internal/models"
)

type XRPTokenService struct {
	dbService database.DatabaseCreate
	client    *http.Client
	apiURL    string
}

type XRPTokenServiceInterface interface {
	FetchTokensFromAPI() ([]models.XRPTokenData, error)
	FetchAndStoreAllTokens() ([]models.XRPTokenData, error)
	StoreTokens(tokens []models.XRPTokenData) error
	GetTokensWithPools() ([]models.XRPTokenData, error)
	UpdateTokenMetrics() error
}

func NewXRPTokenService(dbService database.DatabaseCreate) XRPTokenServiceInterface {
	return &XRPTokenService{
		dbService: dbService,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://s1.xrplmeta.org",
	}
}

// FetchTokensFromAPI fetches token data from the XRPL Meta API
func (s *XRPTokenService) FetchTokensFromAPI() ([]models.XRPTokenData, error) {
	url := fmt.Sprintf("%s/tokens", s.apiURL)
	fmt.Printf("Fetching all tokens from %s...\n", url)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenList models.XRPTokenList
	if err := json.Unmarshal(body, &tokenList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token list: %w", err)
	}

	fmt.Printf("Successfully fetched %d tokens.\n", len(tokenList.Tokens))
	return tokenList.Tokens, nil
}

// FetchAndStoreAllTokens fetches all available tokens from the single endpoint and stores them.
// This replaces the incorrect pagination logic.
func (s *XRPTokenService) FetchAndStoreAllTokens() ([]models.XRPTokenData, error) {
	// 1. Fetch all tokens from the API in a single call.
	allTokens, err := s.FetchTokensFromAPI()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all tokens: %w", err)
	}

	totalTokens := len(allTokens)
	fmt.Printf("Fetched a total of %d tokens. Now preparing to store them in batches.\n", totalTokens)

	// 2. Store the tokens in smaller batches to avoid a single large transaction.
	batchSize := 100
	for i := 0; i < totalTokens; i += batchSize {
		end := i + batchSize
		if end > totalTokens {
			end = totalTokens
		}
		batch := allTokens[i:end]

		fmt.Printf("Storing batch %d/%d (%d tokens)...\n", (i/batchSize)+1, (totalTokens/batchSize)+1, len(batch))
		if err := s.StoreTokens(batch); err != nil {
			// Log a warning but continue processing other batches.
			fmt.Printf("Warning: failed to store token batch: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully stored batch.\n")
		}
	}

	fmt.Printf("Total tokens processed for storage: %d\n", totalTokens)
	return allTokens, nil
}

// StoreTokens stores token data in the xrpTokens table
func (s *XRPTokenService) StoreTokens(tokens []models.XRPTokenData) error {
	if len(tokens) == 0 {
		return nil // Nothing to store
	}
	db, err := s.dbService.GetConnectionXRPDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Prepare batch insert
	for _, token := range tokens {
		// Convert token data to database format
		tokenRecord := map[string]interface{}{
			"issuer":      token.Issuer,
			"currency":    token.Currency,
			"name":        token.Meta.Token.Name,
			"icon":        token.Meta.Token.Icon,
			"description": token.Meta.Token.Description,
			"trustlines":  token.Metrics.Trustlines,
			"holders":     token.Metrics.Holders,
			"supply_xrpl": token.Metrics.Supply,
			"marketcap":   token.Metrics.Marketcap,
			"price":       token.Metrics.Price,
			"volume24h":   token.Metrics.Volume24H,
			"volume7d":    token.Metrics.Volume7D,
		}

		// Convert weblinks to JSON string
		if len(token.Meta.Token.Weblinks) > 0 {
			weblinksJSON, _ := json.Marshal(token.Meta.Token.Weblinks)
			tokenRecord["weblinks"] = string(weblinksJSON)
		}

		// Upsert token data
		sql := `INSERT INTO xrpTokens 
				(issuer, currency, token_name, icon, description, trustlines, holders, supply_xrpl, marketcap, price, volume_24h, weblinks) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE 
				token_name = VALUES(token_name), icon = VALUES(icon), description = VALUES(description),
				trustlines = VALUES(trustlines), holders = VALUES(holders), supply_xrpl = VALUES(supply_xrpl),
				marketcap = VALUES(marketcap), price = VALUES(price), volume_24h = VALUES(volume_24h),
				updated_at = NOW(), weblinks = VALUES(weblinks)`

		if err := tx.Exec(sql,
			tokenRecord["issuer"], tokenRecord["currency"], tokenRecord["name"],
			tokenRecord["icon"], tokenRecord["description"], tokenRecord["trustlines"],
			tokenRecord["holders"], tokenRecord["supply_xrpl"], tokenRecord["marketcap"],
			tokenRecord["price"], tokenRecord["volume24h"],
			tokenRecord["weblinks"]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert token: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTokensWithPools returns tokens that have AMM pools
func (s *XRPTokenService) GetTokensWithPools() ([]models.XRPTokenData, error) {
	db, err := s.dbService.GetConnectionXRPDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	query := `
		SELECT DISTINCT t.issuer, t.currency, t.token_name, t.icon, t.description, 
		       t.trustlines, t.holders, t.supply, t.marketcap, t.price, t.volume24h, t.volume7d
		FROM xrpTokens t
		INNER JOIN xrpAmm_normalized a ON (t.currency = a.amount2currency AND t.issuer = a.amount2issuer)
		WHERE a.liquidity_usd > 100
		ORDER BY t.marketcap DESC
		LIMIT 100
	`

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	var tokens []models.XRPTokenData
	for rows.Next() {
		var token models.XRPTokenData
		var name, icon, description, supply, volume7d string
		var trustlines, holders int
		var marketcap, price, volume24h float64

		if err := rows.Scan(&token.Issuer, &token.Currency, &name, &icon, &description,
			&trustlines, &holders, &supply, &marketcap, &price, &volume24h, &volume7d); err != nil {
			continue // Skip invalid rows
		}

		// Populate the token struct
		token.Meta.Token.Name = name
		token.Meta.Token.Icon = icon
		token.Meta.Token.Description = description
		token.Metrics.Trustlines = trustlines
		token.Metrics.Holders = holders
		token.Metrics.Supply = supply
		token.Metrics.Marketcap = marketcap
		token.Metrics.Price = price
		token.Metrics.Volume24H = volume24h
		token.Metrics.Volume7D = volume7d

		// Query for the XRP/token pool liquidity using the correct table and column names
		poolSQL := `SELECT liquidity_usd FROM xrpAmm_normalized WHERE amount2currency = ? AND amount2issuer = ? AND (amount = 'XRP' OR amount2currency = 'XRP') LIMIT 1`
		var liquidity float64
		db.Raw(poolSQL, token.Currency, token.Issuer).Scan(&liquidity)
		token.LiquidityUSD = liquidity

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// UpdateTokenMetrics updates token metrics from the API
func (s *XRPTokenService) UpdateTokenMetrics() error {
	tokens, err := s.FetchTokensFromAPI()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens for metrics update: %w", err)
	}

	return s.StoreTokens(tokens)
}
