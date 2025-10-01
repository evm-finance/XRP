package xrp

import (
	"database/sql"
	"fmt"
	"time"

	"qc-defi-graphql-server/internal/database"

	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

// AMMLiquidityService handles liquidity calculations for AMM pools
type AMMLiquidityService struct {
	dbService    database.DatabaseCreate
	priceService PriceServiceInterface
	db           *gorm.DB // Shared database connection
}

// LiquidityCalculation represents the calculated liquidity for a pool
type LiquidityCalculation struct {
	PoolID            string          `json:"pool_id"`
	Asset1Currency    string          `json:"asset1_currency"`
	Asset1Issuer      string          `json:"asset1_issuer"`
	Asset1ValueUSD    decimal.Decimal `json:"asset1_value_usd"`
	Asset2Currency    string          `json:"asset2_currency"`
	Asset2Issuer      string          `json:"asset2_issuer"`
	Asset2ValueUSD    decimal.Decimal `json:"asset2_value_usd"`
	TotalLiquidityUSD decimal.Decimal `json:"total_liquidity_usd"`
	Asset1Percentage  decimal.Decimal `json:"asset1_percentage"`
	Asset2Percentage  decimal.Decimal `json:"asset2_percentage"`
	LastUpdated       time.Time       `json:"last_updated"`
}

// NewAMMLiquidityService creates a new AMM liquidity service
func NewAMMLiquidityService(dbService database.DatabaseCreate, priceService PriceServiceInterface) *AMMLiquidityService {
	return &AMMLiquidityService{
		dbService:    dbService,
		priceService: priceService,
	}
}

// NewAMMLiquidityServiceWithDB creates a new AMM liquidity service with a shared database connection
func NewAMMLiquidityServiceWithDB(db *gorm.DB) *AMMLiquidityService {
	priceService := NewPriceServiceWithDB(db)
	return &AMMLiquidityService{
		dbService:    nil, // Not needed when we have direct DB access
		priceService: priceService,
		db:           db, // Store the shared connection
	}
}

// getDB returns the database connection, using shared connection if available
func (s *AMMLiquidityService) getDB() (*gorm.DB, error) {
	if s.db != nil {
		return s.db, nil
	}
	if s.dbService != nil {
		return s.dbService.GetConnectionXRPDB()
	}
	return nil, fmt.Errorf("no database connection available")
}

// CalculatePoolLiquidity calculates the USD liquidity for a specific AMM pool
func (s *AMMLiquidityService) CalculatePoolLiquidity(account string) (*LiquidityCalculation, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %v", err)
	}

	// Get ALL pool data rows for this account (not just one)
	rows, err := db.Raw(`
		SELECT amount, amount2currency, amount2issuer, amount2value,
		       asset2frozen, lptokencurrency, lptokenissuer, lptokenvalue,
		       tradingfee, liquidity_usd, asset1_value_usd, asset2_value_usd,
		       last_liquidity_update, created_at, last_updated
		FROM xrpAmm
		WHERE account = ?
	`, account).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query pool data: %v", err)
	}
	defer rows.Close()

	// Use nullable types for all database columns
	var poolRecords []struct {
		Amount              sql.NullString
		Amount2Currency     sql.NullString
		Amount2Issuer       sql.NullString
		Amount2Value        sql.NullString
		Asset2Frozen        sql.NullInt64
		LPTokenCurrency     sql.NullString
		LPTokenIssuer       sql.NullString
		LPTokenValue        sql.NullString
		TradingFee          sql.NullInt64
		LiquidityUSD        sql.NullFloat64
		Asset1ValueUSD      sql.NullFloat64
		Asset2ValueUSD      sql.NullFloat64
		LastLiquidityUpdate sql.NullTime
		CreatedAt           sql.NullInt64
		LastUpdated         sql.NullInt64
	}

	// Scan all rows for this pool
	for rows.Next() {
		var record struct {
			Amount              sql.NullString
			Amount2Currency     sql.NullString
			Amount2Issuer       sql.NullString
			Amount2Value        sql.NullString
			Asset2Frozen        sql.NullInt64
			LPTokenCurrency     sql.NullString
			LPTokenIssuer       sql.NullString
			LPTokenValue        sql.NullString
			TradingFee          sql.NullInt64
			LiquidityUSD        sql.NullFloat64
			Asset1ValueUSD      sql.NullFloat64
			Asset2ValueUSD      sql.NullFloat64
			LastLiquidityUpdate sql.NullTime
			CreatedAt           sql.NullInt64
			LastUpdated         sql.NullInt64
		}

		if err := rows.Scan(
			&record.Amount, &record.Amount2Currency, &record.Amount2Issuer, &record.Amount2Value,
			&record.Asset2Frozen, &record.LPTokenCurrency, &record.LPTokenIssuer, &record.LPTokenValue,
			&record.TradingFee, &record.LiquidityUSD, &record.Asset1ValueUSD, &record.Asset2ValueUSD,
			&record.LastLiquidityUpdate, &record.CreatedAt, &record.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan pool record: %v", err)
		}

		poolRecords = append(poolRecords, record)
	}

	if len(poolRecords) == 0 {
		return nil, fmt.Errorf("no pool records found for account: %s", account)
	}

	// For now, use the first record that has valid amount data
	// TODO: Consider aggregating data from multiple records if needed
	var selectedRecord struct {
		Amount              sql.NullString
		Amount2Currency     sql.NullString
		Amount2Issuer       sql.NullString
		Amount2Value        sql.NullString
		Asset2Frozen        sql.NullInt64
		LPTokenCurrency     sql.NullString
		LPTokenIssuer       sql.NullString
		LPTokenValue        sql.NullString
		TradingFee          sql.NullInt64
		LiquidityUSD        sql.NullFloat64
		Asset1ValueUSD      sql.NullFloat64
		Asset2ValueUSD      sql.NullFloat64
		LastLiquidityUpdate sql.NullTime
		CreatedAt           sql.NullInt64
		LastUpdated         sql.NullInt64
	}

	// Find first record with valid amount data
	for _, record := range poolRecords {
		if record.Amount.Valid && record.Amount.String != "" {
			selectedRecord = record
			break
		}
	}

	// If no valid record found, use the first one
	if !selectedRecord.Amount.Valid {
		selectedRecord = poolRecords[0]
	}

	// Extract values safely, handling NULLs
	amount := ""
	if selectedRecord.Amount.Valid {
		amount = selectedRecord.Amount.String
	}

	amount2currency := ""
	if selectedRecord.Amount2Currency.Valid {
		amount2currency = selectedRecord.Amount2Currency.String
	}

	amount2issuer := ""
	if selectedRecord.Amount2Issuer.Valid {
		amount2issuer = selectedRecord.Amount2Issuer.String
	}

	amount2value := ""
	if selectedRecord.Amount2Value.Valid {
		amount2value = selectedRecord.Amount2Value.String
	}

	// Parse balances safely
	amount1, _ := decimal.NewFromString(amount)
	amount2, _ := decimal.NewFromString(amount2value)

	// Determine which asset is which
	var asset1Currency, asset1Issuer, asset2Currency, asset2Issuer string
	var asset1Balance, asset2Balance decimal.Decimal

	if amount2currency == "XRP" {
		// amount2 is XRP, amount is the other token
		asset1Currency = "XRP"
		asset1Issuer = ""
		asset1Balance = amount2.Div(decimal.NewFromFloat(1000000.0)) // Convert drops to XRP from database stored value
		asset2Currency = amount2currency
		asset2Issuer = amount2issuer
		asset2Balance = amount1
	} else {
		// amount is XRP, amount2 is the other token
		asset1Currency = "XRP"
		asset1Issuer = ""
		asset1Balance = amount1.Div(decimal.NewFromFloat(1000000.0)) // Convert drops to XRP from database stored value
		asset2Currency = amount2currency
		asset2Issuer = amount2issuer
		asset2Balance = amount2
	}

	// CRITICAL FIX: Use batch pricing to avoid redundant XRP price calls
	assets := []XRPAsset{
		{Currency: asset1Currency, Issuer: asset1Issuer},
		{Currency: asset2Currency, Issuer: asset2Issuer},
	}

	prices, err := s.priceService.GetMultipleAssetPrices(assets)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch prices: %v", err)
	}

	// Extract individual prices from batch result
	asset1Key := fmt.Sprintf("%s_%s", asset1Currency, asset1Issuer)
	asset1Price, exists := prices[asset1Key]
	if !exists || asset1Price.IsZero() {
		return nil, fmt.Errorf("failed to get price for %s/%s", asset1Currency, asset1Issuer)
	}

	asset2Key := fmt.Sprintf("%s_%s", asset2Currency, asset2Issuer)
	asset2Price, exists := prices[asset2Key]
	if !exists || asset2Price.IsZero() {
		return nil, fmt.Errorf("failed to get price for %s/%s", asset2Currency, asset2Issuer)
	}

	// Calculate USD values
	asset1ValueUSD := asset1Balance.Mul(asset1Price)
	asset2ValueUSD := asset2Balance.Mul(asset2Price)

	// Determine which asset is legitimate (XRP or RLUSD) and calculate liquidity based on that
	var legitimateAssetValueUSD decimal.Decimal

	// Check if asset1 is legitimate (XRP or RLUSD)
	if asset1Currency == "XRP" || (asset1Currency == RLUSDCurrency && asset1Issuer == RLUSDIssuer) {
		legitimateAssetValueUSD = asset1ValueUSD
	} else if asset2Currency == "XRP" || (asset2Currency == RLUSDCurrency && asset2Issuer == RLUSDIssuer) {
		// Check if asset2 is legitimate (XRP or RLUSD)
		legitimateAssetValueUSD = asset2ValueUSD
	} else {
		// Neither asset is legitimate, use the lower value to avoid scam tokens
		if asset1ValueUSD.LessThan(asset2ValueUSD) {
			legitimateAssetValueUSD = asset1ValueUSD
		} else {
			legitimateAssetValueUSD = asset2ValueUSD
		}
	}

	// Calculate total liquidity as 2x the legitimate asset value (or lower value for non-legitimate pairs)
	// This filters out scam tokens with artificially inflated prices
	totalLiquidityUSD := legitimateAssetValueUSD.Mul(decimal.NewFromInt(2))

	// Calculate percentages based on the new total
	var asset1Percentage, asset2Percentage decimal.Decimal
	if totalLiquidityUSD.GreaterThan(decimal.Zero) {
		asset1Percentage = asset1ValueUSD.Div(totalLiquidityUSD).Mul(decimal.NewFromInt(100))
		asset2Percentage = asset2ValueUSD.Div(totalLiquidityUSD).Mul(decimal.NewFromInt(100))
	}

	return &LiquidityCalculation{
		PoolID:            account,
		Asset1Currency:    asset1Currency,
		Asset1Issuer:      asset1Issuer,
		Asset1ValueUSD:    asset1ValueUSD,
		Asset2Currency:    asset2Currency,
		Asset2Issuer:      asset2Issuer,
		Asset2ValueUSD:    asset2ValueUSD,
		TotalLiquidityUSD: totalLiquidityUSD,
		Asset1Percentage:  asset1Percentage,
		Asset2Percentage:  asset2Percentage,
		LastUpdated:       time.Now(),
	}, nil
}

// UpdateAllPoolLiquidity calculates and updates liquidity for all AMM pools
func (s *AMMLiquidityService) UpdateAllPoolLiquidity() error {
	db, err := s.getDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	// Get all unique pool accounts
	rows, err := db.Raw("SELECT DISTINCT account FROM xrpAmm_normalized").Rows()
	if err != nil {
		return fmt.Errorf("failed to query pools: %v", err)
	}
	defer rows.Close()

	var accounts []string
	for rows.Next() {
		var account string
		if err := rows.Scan(&account); err != nil {
			fmt.Printf("⚠️  Failed to scan account row: %v\n", err)
			continue
		}
		accounts = append(accounts, account)
	}

	fmt.Printf("🔄 Updating liquidity for %d unique pools...\n", len(accounts))

	totalProcessed := 0
	totalSucceeded := 0
	totalFailed := 0

	// Calculate and update liquidity for each pool
	for i, account := range accounts {
		totalProcessed++
		calculation, err := s.CalculatePoolLiquidity(account)
		if err != nil {
			totalFailed++
			fmt.Printf("⚠️  Failed to calculate liquidity for pool %s (index %d): %v\n", account, i+1, err)
			continue
		}

		// Update the database - FIXED: Use correct normalized table
		updateSQL := `
			UPDATE xrpAmm_normalized 
			SET liquidity_usd = ?, 
				asset1_value_usd = ?, 
				asset2_value_usd = ?, 
				last_updated = ?
			WHERE account = ?
		`

		result := db.Exec(updateSQL,
			calculation.TotalLiquidityUSD,
			calculation.Asset1ValueUSD,
			calculation.Asset2ValueUSD,
			calculation.LastUpdated,
			account,
		)

		if result.Error != nil {
			totalFailed++
			fmt.Printf("⚠️  Failed to update liquidity for pool %s (index %d): %v\n", account, i+1, result.Error)
		} else {
			totalSucceeded++
			fmt.Printf("✅ Pool %d/%d: %s - $%s USD liquidity\n",
				i+1, len(accounts), account, calculation.TotalLiquidityUSD.StringFixed(2))
		}
	}

	fmt.Printf("\n--- Liquidity Update Summary ---\n")
	fmt.Printf("Total unique pools processed: %d\n", totalProcessed)
	fmt.Printf("Succeeded: %d\n", totalSucceeded)
	fmt.Printf("Failed/skipped: %d\n", totalFailed)
	fmt.Println("🎉 Liquidity update complete!")
	return nil
}

// GetPoolLiquidity retrieves the current liquidity calculation for a pool
func (s *AMMLiquidityService) GetPoolLiquidity(account string) (*LiquidityCalculation, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %v", err)
	}

	row := db.Raw(`
		SELECT liquidity_usd, asset1_value_usd, asset2_value_usd, last_liquidity_update
		FROM xrpAmm
		WHERE account = ?
		LIMIT 1
	`, account).Row()

	var liquidityUSD, asset1ValueUSD, asset2ValueUSD decimal.Decimal
	var lastUpdate time.Time
	if err := row.Scan(&liquidityUSD, &asset1ValueUSD, &asset2ValueUSD, &lastUpdate); err != nil {
		return nil, fmt.Errorf("pool liquidity not found: %v", err)
	}

	// Get pool details for asset information
	poolRow := db.Raw(`
		SELECT amount2currency, amount2issuer
		FROM xrpAmm
		WHERE account = ?
		LIMIT 1
	`, account).Row()

	var amount2currency, amount2issuer string
	if err := poolRow.Scan(&amount2currency, &amount2issuer); err != nil {
		return nil, fmt.Errorf("pool details not found: %v", err)
	}

	// Calculate percentages
	var asset1Percentage, asset2Percentage decimal.Decimal
	if liquidityUSD.GreaterThan(decimal.Zero) {
		asset1Percentage = asset1ValueUSD.Div(liquidityUSD).Mul(decimal.NewFromInt(100))
		asset2Percentage = asset2ValueUSD.Div(liquidityUSD).Mul(decimal.NewFromInt(100))
	}

	return &LiquidityCalculation{
		PoolID:            account,
		Asset1Currency:    "XRP",
		Asset1Issuer:      "",
		Asset1ValueUSD:    asset1ValueUSD,
		Asset2Currency:    amount2currency,
		Asset2Issuer:      amount2issuer,
		Asset2ValueUSD:    asset2ValueUSD,
		TotalLiquidityUSD: liquidityUSD,
		Asset1Percentage:  asset1Percentage,
		Asset2Percentage:  asset2Percentage,
		LastUpdated:       lastUpdate,
	}, nil
}

// GetAllPoolLiquidity retrieves liquidity for all pools
func (s *AMMLiquidityService) GetAllPoolLiquidity() ([]*LiquidityCalculation, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %v", err)
	}

	rows, err := db.Raw(`
		SELECT account, liquidity_usd, asset1_value_usd, asset2_value_usd, 
		       amount2currency, amount2issuer, last_liquidity_update
		FROM xrpAmm
		WHERE liquidity_usd > 0
		ORDER BY liquidity_usd DESC
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query pool liquidity: %v", err)
	}
	defer rows.Close()

	var calculations []*LiquidityCalculation
	for rows.Next() {
		var account, amount2currency, amount2issuer string
		var liquidityUSD, asset1ValueUSD, asset2ValueUSD decimal.Decimal
		var lastUpdate time.Time

		if err := rows.Scan(&account, &liquidityUSD, &asset1ValueUSD, &asset2ValueUSD,
			&amount2currency, &amount2issuer, &lastUpdate); err != nil {
			continue
		}

		// Calculate percentages
		var asset1Percentage, asset2Percentage decimal.Decimal
		if liquidityUSD.GreaterThan(decimal.Zero) {
			asset1Percentage = asset1ValueUSD.Div(liquidityUSD).Mul(decimal.NewFromInt(100))
			asset2Percentage = asset2ValueUSD.Div(liquidityUSD).Mul(decimal.NewFromInt(100))
		}

		calculation := &LiquidityCalculation{
			PoolID:            account,
			Asset1Currency:    "XRP",
			Asset1Issuer:      "",
			Asset1ValueUSD:    asset1ValueUSD,
			Asset2Currency:    amount2currency,
			Asset2Issuer:      amount2issuer,
			Asset2ValueUSD:    asset2ValueUSD,
			TotalLiquidityUSD: liquidityUSD,
			Asset1Percentage:  asset1Percentage,
			Asset2Percentage:  asset2Percentage,
			LastUpdated:       lastUpdate,
		}

		calculations = append(calculations, calculation)
	}

	return calculations, nil
}
