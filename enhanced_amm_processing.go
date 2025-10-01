package xrp

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

// Enhanced AMM Transaction Processing
type EnhancedAMMProcessor struct {
	priceService PriceServiceInterface
}

// NewEnhancedAMMProcessor creates a new enhanced AMM processor
func NewEnhancedAMMProcessor(priceService PriceServiceInterface) *EnhancedAMMProcessor {
	return &EnhancedAMMProcessor{
		priceService: priceService,
	}
}

// ProcessAMMTransaction provides enhanced processing of AMM transactions
func (e *EnhancedAMMProcessor) ProcessAMMTransaction(tx *XRPTransaction) (*AMMTransaction, error) {
	// Start with basic AMM transaction
	ammTx := &AMMTransaction{
		TransactionHash: tx.Hash,
		TransactionType: tx.TransactionType,
		Account:         tx.Account,
		LedgerIndex:     tx.LedgerIndex,
		Timestamp:       int64(tx.Date),
		Metadata:        make(map[string]interface{}),
	}

	// Enhanced metadata extraction
	if err := e.extractEnhancedMetadata(tx, ammTx); err != nil {
		log.Printf("Warning: Error extracting enhanced metadata: %v", err)
		// Continue with basic processing
	}

	// Calculate USD values
	e.calculateUSDValues(ammTx)

	return ammTx, nil
}

// Extract enhanced metadata from transaction
func (e *EnhancedAMMProcessor) extractEnhancedMetadata(tx *XRPTransaction, ammTx *AMMTransaction) error {
	// Parse transaction fee
	if tx.Fee != "" {
		if fee, err := decimal.NewFromString(tx.Fee); err == nil {
			feeDecimal := fee.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
			ammTx.Fee, _ = feeDecimal.Float64()
			ammTx.Metadata["fee_xrp"] = ammTx.Fee
		}
	}

	// Extract assets and amounts based on transaction type
	switch tx.TransactionType {
	case "AMMCreate":
		return e.extractAMMCreateData(tx, ammTx)
	case "AMMDeposit":
		return e.extractAMMDepositData(tx, ammTx)
	case "AMMWithdraw":
		return e.extractAMMWithdrawData(tx, ammTx)
	case "AMMVote":
		return e.extractAMMVoteData(tx, ammTx)
	case "AMMBid":
		return e.extractAMMBidData(tx, ammTx)
	case "AMMDelete":
		return e.extractAMMDeleteData(tx, ammTx)
	}

	return nil
}

// Extract AMMCreate transaction data
func (e *EnhancedAMMProcessor) extractAMMCreateData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	// Extract initial assets and amounts
	if tx.Amount != nil {
		asset1, amount1 := e.parseAssetAmount(tx.Amount)
		ammTx.Asset1 = asset1
		ammTx.Amount1, _ = amount1.Float64()
		ammTx.Metadata["initial_asset1_amount"] = amount1
	}

	// Extract from transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract Amount2 (second asset)
		if amount2Field, exists := meta["Amount2"]; exists {
			asset2, amount2 := e.parseAssetAmount(amount2Field)
			ammTx.Asset2 = asset2
			ammTx.Amount2, _ = amount2.Float64()
			ammTx.Metadata["initial_asset2_amount"] = amount2
		}

		// Extract trading fee
		if tradingFee, exists := meta["TradingFee"]; exists {
			if fee, ok := tradingFee.(int); ok {
				ammTx.Metadata["trading_fee_bps"] = fee // basis points
				ammTx.Metadata["trading_fee_percent"] = decimal.NewFromInt(int64(fee)).Div(decimal.NewFromInt(10000))
			}
		}

		// Extract LP token information
		if lpToken, exists := meta["LPToken"]; exists {
			if lpMap, ok := lpToken.(map[string]interface{}); ok {
				if value, ok := lpMap["value"].(string); ok {
					if lpAmount, err := decimal.NewFromString(value); err == nil {
						ammTx.LPAmount, _ = lpAmount.Float64()
						ammTx.Metadata["initial_lp_tokens"] = lpAmount
					}
				}
				if currency, ok := lpMap["currency"].(string); ok {
					ammTx.Metadata["lp_token_currency"] = currency
				}
				if issuer, ok := lpMap["issuer"].(string); ok {
					ammTx.Metadata["lp_token_issuer"] = issuer
				}
			}
		}

		// Calculate initial price
		if ammTx.Amount1 != 0 && ammTx.Amount2 != 0 {
			initialPrice := decimal.NewFromFloat(ammTx.Amount1).Div(decimal.NewFromFloat(ammTx.Amount2))
			ammTx.Metadata["initial_price"] = initialPrice
			ammTx.Metadata["initial_price_direction"] = fmt.Sprintf("%s per %s",
				ammTx.Asset1.Currency, ammTx.Asset2.Currency)
		}

		// Store pool creation timestamp
		ammTx.Metadata["pool_created_at"] = time.Unix(ammTx.Timestamp, 0).Format(time.RFC3339)
	}

	return nil
}

// Extract AMMDeposit transaction data
func (e *EnhancedAMMProcessor) extractAMMDepositData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	// Extract deposit amounts
	if tx.Amount != nil {
		asset1, amount1 := e.parseAssetAmount(tx.Amount)
		ammTx.Asset1 = asset1
		ammTx.Amount1, _ = amount1.Float64()
		ammTx.Metadata["deposit_asset1_amount"] = amount1
	}

	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract second asset deposit
		if amount2Field, exists := meta["Amount2"]; exists {
			asset2, amount2 := e.parseAssetAmount(amount2Field)
			ammTx.Asset2 = asset2
			ammTx.Amount2, _ = amount2.Float64()
			ammTx.Metadata["deposit_asset2_amount"] = amount2
		}

		// Extract LP tokens received
		if lpToken, exists := meta["LPTokenOut"]; exists {
			if lpMap, ok := lpToken.(map[string]interface{}); ok {
				if value, ok := lpMap["value"].(string); ok {
					if lpAmount, err := decimal.NewFromString(value); err == nil {
						ammTx.LPAmount, _ = lpAmount.Float64()
						ammTx.Metadata["lp_tokens_received"] = lpAmount
					}
				}
			}
		}

		// Determine deposit type
		singleAssetDeposit := (ammTx.Amount1 == 0) != (ammTx.Amount2 == 0)
		ammTx.Metadata["single_asset_deposit"] = singleAssetDeposit

		if singleAssetDeposit {
			if ammTx.Amount1 > 0 {
				ammTx.Metadata["deposit_type"] = "single_asset_1"
				ammTx.Metadata["deposited_asset"] = ammTx.Asset1.Currency
			} else {
				ammTx.Metadata["deposit_type"] = "single_asset_2"
				ammTx.Metadata["deposited_asset"] = ammTx.Asset2.Currency
			}
		} else {
			ammTx.Metadata["deposit_type"] = "dual_asset"
		}

		// Calculate effective exchange rate
		if ammTx.Amount1 != 0 && ammTx.Amount2 != 0 {
			effectiveRate := decimal.NewFromFloat(ammTx.Amount1).Div(decimal.NewFromFloat(ammTx.Amount2))
			ammTx.Metadata["effective_exchange_rate"] = effectiveRate
		}

		// Estimate price impact (simplified)
		totalValue := decimal.NewFromFloat(ammTx.Amount1).Add(decimal.NewFromFloat(ammTx.Amount2))
		var priceImpact decimal.Decimal
		if totalValue.GreaterThan(decimal.NewFromInt(10000)) {
			priceImpact = decimal.NewFromFloat(0.005) // 0.5% for large deposits
		} else if totalValue.GreaterThan(decimal.NewFromInt(1000)) {
			priceImpact = decimal.NewFromFloat(0.001) // 0.1% for medium deposits
		} else {
			priceImpact = decimal.NewFromFloat(0.0001) // 0.01% for small deposits
		}
		ammTx.Metadata["estimated_price_impact"] = priceImpact
	}

	return nil
}

// Extract AMMWithdraw transaction data
func (e *EnhancedAMMProcessor) extractAMMWithdrawData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	// Extract withdrawal amounts
	if tx.Amount != nil {
		asset1, amount1 := e.parseAssetAmount(tx.Amount)
		ammTx.Asset1 = asset1
		ammTx.Amount1, _ = amount1.Float64()
		ammTx.Metadata["withdraw_asset1_amount"] = amount1
	}

	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract second asset withdrawal
		if amount2Field, exists := meta["Amount2"]; exists {
			asset2, amount2 := e.parseAssetAmount(amount2Field)
			ammTx.Asset2 = asset2
			ammTx.Amount2, _ = amount2.Float64()
			ammTx.Metadata["withdraw_asset2_amount"] = amount2
		}

		// Extract LP tokens burned
		if lpToken, exists := meta["LPTokenIn"]; exists {
			if lpMap, ok := lpToken.(map[string]interface{}); ok {
				if value, ok := lpMap["value"].(string); ok {
					if lpAmount, err := decimal.NewFromString(value); err == nil {
						ammTx.LPAmount, _ = lpAmount.Float64()
						ammTx.Metadata["lp_tokens_burned"] = lpAmount
					}
				}
			}
		}

		// Determine withdrawal type
		singleAssetWithdraw := (ammTx.Amount1 == 0) != (ammTx.Amount2 == 0)
		ammTx.Metadata["single_asset_withdraw"] = singleAssetWithdraw

		if singleAssetWithdraw {
			if ammTx.Amount1 > 0 {
				ammTx.Metadata["withdraw_type"] = "single_asset_1"
				ammTx.Metadata["withdrawn_asset"] = ammTx.Asset1.Currency
			} else {
				ammTx.Metadata["withdraw_type"] = "single_asset_2"
				ammTx.Metadata["withdrawn_asset"] = ammTx.Asset2.Currency
			}
		} else {
			ammTx.Metadata["withdraw_type"] = "dual_asset"
		}

		// Calculate withdrawal percentage (requires pool state)
		// This is a simplified estimate
		if ammTx.LPAmount > 0 {
			// Estimate total LP supply (would need actual pool data)
			estimatedTotalLP := decimal.NewFromFloat(ammTx.LPAmount).Mul(decimal.NewFromInt(10)) // Rough estimate
			withdrawPercentage := decimal.NewFromFloat(ammTx.LPAmount).Div(estimatedTotalLP).Mul(decimal.NewFromInt(100))
			ammTx.Metadata["estimated_withdraw_percentage"] = withdrawPercentage
		}
	}

	return nil
}

// Extract AMMVote transaction data
func (e *EnhancedAMMProcessor) extractAMMVoteData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract trading fee vote
		if tradingFee, exists := meta["TradingFee"]; exists {
			if fee, ok := tradingFee.(int); ok {
				ammTx.Metadata["voted_trading_fee_bps"] = fee
				ammTx.Metadata["voted_trading_fee_percent"] = decimal.NewFromInt(int64(fee)).Div(decimal.NewFromInt(10000))
			}
		}

		// Extract vote weight
		if voteSlots, exists := meta["VoteSlots"]; exists {
			if slotsArray, ok := voteSlots.([]interface{}); ok {
				totalWeight := 0
				for _, slot := range slotsArray {
					if slotMap, ok := slot.(map[string]interface{}); ok {
						if weight, ok := slotMap["VoteWeight"].(int); ok {
							if account, ok := slotMap["Account"].(string); ok && account == tx.Account {
								ammTx.Metadata["vote_weight"] = weight
							}
							totalWeight += weight
						}
					}
				}
				ammTx.Metadata["total_vote_weight"] = totalWeight

				if voteWeight, exists := ammTx.Metadata["vote_weight"]; exists {
					if weight, ok := voteWeight.(int); ok && totalWeight > 0 {
						votePercentage := decimal.NewFromInt(int64(weight)).Div(decimal.NewFromInt(int64(totalWeight))).Mul(decimal.NewFromInt(100))
						ammTx.Metadata["vote_percentage"] = votePercentage
					}
				}
			}
		}
	}

	return nil
}

// Extract AMMBid transaction data
func (e *EnhancedAMMProcessor) extractAMMBidData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	// Extract bid amount
	if tx.Amount != nil {
		asset, amount := e.parseAssetAmount(tx.Amount)
		ammTx.Asset1 = asset
		ammTx.Amount1, _ = amount.Float64()
		ammTx.Metadata["bid_amount"] = amount
		ammTx.Metadata["bid_asset"] = asset.Currency
	}

	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract auction slot information
		if auctionSlot, exists := meta["AuctionSlot"]; exists {
			if slotMap, ok := auctionSlot.(map[string]interface{}); ok {
				if timeInterval, ok := slotMap["TimeInterval"].(int); ok {
					ammTx.Metadata["auction_time_interval"] = timeInterval
				}
				if discountedFee, ok := slotMap["DiscountedFee"].(int); ok {
					ammTx.Metadata["discounted_fee_bps"] = discountedFee
					ammTx.Metadata["discounted_fee_percent"] = decimal.NewFromInt(int64(discountedFee)).Div(decimal.NewFromInt(10000))
				}
				if expiration, ok := slotMap["Expiration"].(string); ok {
					ammTx.Metadata["auction_expiration"] = expiration
				}
			}
		}

		// Calculate bid value in USD
		if ammTx.Amount1 > 0 && e.priceService != nil {
			if price, err := e.priceService.GetAssetPrice(ammTx.Asset1); err == nil {
				bidValueUSD := decimal.NewFromFloat(ammTx.Amount1).Mul(price)
				ammTx.Metadata["bid_value_usd"] = bidValueUSD
			}
		}
	}

	return nil
}

// Extract AMMDelete transaction data
func (e *EnhancedAMMProcessor) extractAMMDeleteData(tx *XRPTransaction, ammTx *AMMTransaction) error {
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		// Extract final pool state before deletion
		if deletedNode, exists := meta["DeletedNode"]; exists {
			if nodeMap, ok := deletedNode.(map[string]interface{}); ok {
				if finalObject, exists := nodeMap["FinalFields"]; exists {
					if finalMap, ok := finalObject.(map[string]interface{}); ok {
						// Extract final balances
						if amount, ok := finalMap["Amount"].(string); ok {
							if finalAmount, err := decimal.NewFromString(amount); err == nil {
								ammTx.Amount1, _ = finalAmount.Div(decimal.NewFromInt(1000000)).Float64() // Convert drops to XRP
								ammTx.Metadata["final_asset1_balance"] = ammTx.Amount1
							}
						}

						if amount2, ok := finalMap["Amount2"]; ok {
							if amount2Map, ok := amount2.(map[string]interface{}); ok {
								if value, ok := amount2Map["value"].(string); ok {
									if finalAmount2, err := decimal.NewFromString(value); err == nil {
										ammTx.Amount2, _ = finalAmount2.Float64()
										ammTx.Metadata["final_asset2_balance"] = ammTx.Amount2
									}
								}
							}
						}

						// Calculate final liquidity value
						totalFinalValue := decimal.NewFromFloat(ammTx.Amount1).Add(decimal.NewFromFloat(ammTx.Amount2)) // Simplified
						ammTx.Metadata["final_liquidity_estimate"] = totalFinalValue
					}
				}
			}
		}

		// Record deletion timestamp
		ammTx.Metadata["pool_deleted_at"] = time.Unix(ammTx.Timestamp, 0).Format(time.RFC3339)
	}

	return nil
}

// Parse asset and amount from transaction field
func (e *EnhancedAMMProcessor) parseAssetAmount(field interface{}) (XRPAsset, decimal.Decimal) {
	var asset XRPAsset
	var amount decimal.Decimal

	switch v := field.(type) {
	case string:
		// XRP amount in drops
		if amt, err := decimal.NewFromString(v); err == nil {
			asset = XRPAsset{Currency: "XRP"}
			amount = amt.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
		}
	case map[string]interface{}:
		// Token amount
		if currency, ok := v["currency"].(string); ok {
			asset.Currency = currency
		}
		if issuer, ok := v["issuer"].(string); ok {
			asset.Issuer = issuer
		}
		if value, ok := v["value"].(string); ok {
			if amt, err := decimal.NewFromString(value); err == nil {
				amount = amt
			}
		}
	}

	return asset, amount
}

// Calculate USD values for the transaction
func (e *EnhancedAMMProcessor) calculateUSDValues(ammTx *AMMTransaction) {
	if e.priceService == nil {
		return
	}

	var totalValueUSD decimal.Decimal

	// CRITICAL FIX: Use batch pricing to avoid redundant XRP price calls
	assets := []XRPAsset{}
	if ammTx.Asset1.Currency != "" {
		assets = append(assets, ammTx.Asset1)
	}
	if ammTx.Asset2.Currency != "" {
		assets = append(assets, ammTx.Asset2)
	}

	if len(assets) == 0 {
		return
	}

	// Get all prices with a single XRP price fetch
	prices, err := e.priceService.GetMultipleAssetPrices(assets)
	if err != nil {
		log.Printf("Failed to get batch prices: %v", err)
		return
	}

	// Process asset1 price
	if ammTx.Asset1.Currency != "" {
		key := fmt.Sprintf("%s_%s", ammTx.Asset1.Currency, ammTx.Asset1.Issuer)
		if price, exists := prices[key]; exists && price.GreaterThan(decimal.Zero) {
			asset1ValueUSD := decimal.NewFromFloat(ammTx.Amount1).Mul(price)
			ammTx.Metadata["asset1_price_usd"] = price
			ammTx.Metadata["asset1_value_usd"] = asset1ValueUSD
			totalValueUSD = totalValueUSD.Add(asset1ValueUSD)
		}
	}

	// Process asset2 price
	if ammTx.Asset2.Currency != "" {
		key := fmt.Sprintf("%s_%s", ammTx.Asset2.Currency, ammTx.Asset2.Issuer)
		if price, exists := prices[key]; exists && price.GreaterThan(decimal.Zero) {
			asset2ValueUSD := decimal.NewFromFloat(ammTx.Amount2).Mul(price)
			ammTx.Metadata["asset2_price_usd"] = price
			ammTx.Metadata["asset2_value_usd"] = asset2ValueUSD
			totalValueUSD = totalValueUSD.Add(asset2ValueUSD)
		}
	}

	// Store total USD value
	if totalValueUSD.GreaterThan(decimal.Zero) {
		ammTx.Metadata["total_value_usd"] = totalValueUSD

		// Categorize transaction size
		var sizeCategory string
		switch {
		case totalValueUSD.GreaterThanOrEqual(decimal.NewFromInt(100000)):
			sizeCategory = "whale"
		case totalValueUSD.GreaterThanOrEqual(decimal.NewFromInt(10000)):
			sizeCategory = "large"
		case totalValueUSD.GreaterThanOrEqual(decimal.NewFromInt(1000)):
			sizeCategory = "medium"
		default:
			sizeCategory = "small"
		}
		ammTx.Metadata["transaction_size_category"] = sizeCategory
	}
}

// ProcessLedgerAMMTransactions processes all AMM transactions in a ledger with enhanced metadata
func (e *EnhancedAMMProcessor) ProcessLedgerAMMTransactions(ledger *XRPLedger) ([]*AMMTransaction, error) {
	var enhancedAMMTransactions []*AMMTransaction

	if ledger.Ledger == nil || ledger.Ledger.Transactions == nil {
		return enhancedAMMTransactions, nil
	}

	for _, tx := range ledger.Ledger.Transactions {
		if e.isAMMTransaction(tx) {
			enhancedTx, err := e.ProcessAMMTransaction(tx)
			if err != nil {
				log.Printf("Error processing AMM transaction %s: %v", tx.Hash, err)
				continue
			}
			enhancedAMMTransactions = append(enhancedAMMTransactions, enhancedTx)
		}
	}

	// Add ledger-level statistics
	if len(enhancedAMMTransactions) > 0 {
		e.addLedgerStatistics(ledger, enhancedAMMTransactions)
	}

	return enhancedAMMTransactions, nil
}

// Add statistics about AMM activity in the ledger
func (e *EnhancedAMMProcessor) addLedgerStatistics(ledger *XRPLedger, ammTxs []*AMMTransaction) {
	stats := map[string]interface{}{
		"ledger_index":      ledger.LedgerIndex,
		"total_amm_txs":     len(ammTxs),
		"ledger_close_time": ledger.Ledger.CloseTime,
		"transaction_types": make(map[string]int),
		"total_volume_usd":  decimal.Zero,
		"unique_accounts":   make(map[string]bool),
	}

	var totalVolumeUSD decimal.Decimal
	uniqueAccounts := make(map[string]bool)
	txTypes := make(map[string]int)

	for _, tx := range ammTxs {
		// Count transaction types
		txTypes[tx.TransactionType]++

		// Track unique accounts
		uniqueAccounts[tx.Account] = true

		// Sum total volume
		if totalValue, exists := tx.Metadata["total_value_usd"]; exists {
			if value, ok := totalValue.(decimal.Decimal); ok {
				totalVolumeUSD = totalVolumeUSD.Add(value)
			}
		}
	}

	stats["transaction_types"] = txTypes
	stats["total_volume_usd"] = totalVolumeUSD
	stats["unique_accounts_count"] = len(uniqueAccounts)

	// Calculate average transaction size
	if len(ammTxs) > 0 && totalVolumeUSD.GreaterThan(decimal.Zero) {
		stats["avg_transaction_size_usd"] = totalVolumeUSD.Div(decimal.NewFromInt(int64(len(ammTxs))))
	}

	// Store statistics in each transaction for easy access
	for _, tx := range ammTxs {
		tx.Metadata["ledger_statistics"] = stats
	}
}

// Check if transaction is AMM-related
func (e *EnhancedAMMProcessor) isAMMTransaction(tx *XRPTransaction) bool {
	ammTypes := []string{
		"AMMCreate", "AMMDelete", "AMMDeposit", "AMMWithdraw",
		"AMMBid", "AMMVote",
	}

	for _, ammType := range ammTypes {
		if tx.TransactionType == ammType {
			return true
		}
	}
	return false
}

// GetTransactionSummary provides a human-readable summary of an enhanced AMM transaction
func (e *EnhancedAMMProcessor) GetTransactionSummary(ammTx *AMMTransaction) string {
	summary := fmt.Sprintf("%s transaction by %s", ammTx.TransactionType, ammTx.Account)

	// Add transaction-specific details
	switch ammTx.TransactionType {
	case "AMMCreate":
		if ammTx.Amount1 > 0 && ammTx.Amount2 > 0 {
			summary += fmt.Sprintf(" - Created pool with %.6f %s and %.6f %s",
				ammTx.Amount1, ammTx.Asset1.Currency,
				ammTx.Amount2, ammTx.Asset2.Currency)
		}
	case "AMMDeposit":
		if singleAsset, exists := ammTx.Metadata["single_asset_deposit"]; exists && singleAsset.(bool) {
			if ammTx.Amount1 > 0 {
				summary += fmt.Sprintf(" - Single asset deposit: %.6f %s", ammTx.Amount1, ammTx.Asset1.Currency)
			} else {
				summary += fmt.Sprintf(" - Single asset deposit: %.6f %s", ammTx.Amount2, ammTx.Asset2.Currency)
			}
		} else {
			summary += fmt.Sprintf(" - Dual asset deposit: %.6f %s + %.6f %s",
				ammTx.Amount1, ammTx.Asset1.Currency,
				ammTx.Amount2, ammTx.Asset2.Currency)
		}
	case "AMMWithdraw":
		if singleAsset, exists := ammTx.Metadata["single_asset_withdraw"]; exists && singleAsset.(bool) {
			if ammTx.Amount1 > 0 {
				summary += fmt.Sprintf(" - Single asset withdrawal: %.6f %s", ammTx.Amount1, ammTx.Asset1.Currency)
			} else {
				summary += fmt.Sprintf(" - Single asset withdrawal: %.6f %s", ammTx.Amount2, ammTx.Asset2.Currency)
			}
		} else {
			summary += fmt.Sprintf(" - Dual asset withdrawal: %.6f %s + %.6f %s",
				ammTx.Amount1, ammTx.Asset1.Currency,
				ammTx.Amount2, ammTx.Asset2.Currency)
		}
	case "AMMVote":
		if fee, exists := ammTx.Metadata["voted_trading_fee_percent"]; exists {
			summary += fmt.Sprintf(" - Voted for %s%% trading fee", fee.(decimal.Decimal).StringFixed(4))
		}
	case "AMMBid":
		if ammTx.Amount1 > 0 {
			summary += fmt.Sprintf(" - Bid %.6f %s for auction slot", ammTx.Amount1, ammTx.Asset1.Currency)
		}
	case "AMMDelete":
		summary += " - Deleted AMM pool"
	}

	// Add USD value if available
	if totalValue, exists := ammTx.Metadata["total_value_usd"]; exists {
		if value, ok := totalValue.(decimal.Decimal); ok && value.GreaterThan(decimal.Zero) {
			summary += fmt.Sprintf(" (≈$%s USD)", value.StringFixed(2))
		}
	}

	return summary
}
