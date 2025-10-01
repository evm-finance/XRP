package xrp

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

// Enhanced AMM Metadata Extractor
type EnhancedAMMMetadataExtractor struct {
	priceService PriceServiceInterface
}

// AMMTransactionMetadata represents comprehensive metadata extracted from AMM transactions
type AMMTransactionMetadata struct {
	// Basic transaction info
	TransactionHash string          `json:"transaction_hash"`
	TransactionType string          `json:"transaction_type"`
	Account         string          `json:"account"`
	LedgerIndex     int             `json:"ledger_index"`
	Timestamp       time.Time       `json:"timestamp"`
	Fee             decimal.Decimal `json:"fee"`

	// AMM-specific data
	AMMAccount string          `json:"amm_account,omitempty"`
	Asset1     XRPAsset        `json:"asset1,omitempty"`
	Asset2     XRPAsset        `json:"asset2,omitempty"`
	Amount1    decimal.Decimal `json:"amount1,omitempty"`
	Amount2    decimal.Decimal `json:"amount2,omitempty"`
	LPTokens   decimal.Decimal `json:"lp_tokens,omitempty"`
	TradingFee int             `json:"trading_fee,omitempty"`

	// Enhanced metadata for different transaction types
	CreateData   *AMMCreateMetadata   `json:"create_data,omitempty"`
	DepositData  *AMMDepositMetadata  `json:"deposit_data,omitempty"`
	WithdrawData *AMMWithdrawMetadata `json:"withdraw_data,omitempty"`
	VoteData     *AMMVoteMetadata     `json:"vote_data,omitempty"`
	BidData      *AMMBidMetadata      `json:"bid_data,omitempty"`

	// Price and value calculations
	Asset1PriceUSD decimal.Decimal `json:"asset1_price_usd"`
	Asset2PriceUSD decimal.Decimal `json:"asset2_price_usd"`
	TotalValueUSD  decimal.Decimal `json:"total_value_usd"`

	// Additional metadata
	RawMetadata map[string]interface{} `json:"raw_metadata,omitempty"`
}

// Specific metadata for each AMM transaction type
type AMMCreateMetadata struct {
	InitialAsset1Amount decimal.Decimal `json:"initial_asset1_amount"`
	InitialAsset2Amount decimal.Decimal `json:"initial_asset2_amount"`
	InitialLPTokens     decimal.Decimal `json:"initial_lp_tokens"`
	TradingFee          int             `json:"trading_fee"`
	InitialPrice        decimal.Decimal `json:"initial_price"`
	CreatedAt           time.Time       `json:"created_at"`
}

type AMMDepositMetadata struct {
	SingleAssetDeposit bool            `json:"single_asset_deposit"`
	DepositAsset1      decimal.Decimal `json:"deposit_asset1"`
	DepositAsset2      decimal.Decimal `json:"deposit_asset2"`
	LPTokensReceived   decimal.Decimal `json:"lp_tokens_received"`
	EffectivePrice     decimal.Decimal `json:"effective_price"`
	PriceImpact        decimal.Decimal `json:"price_impact"`
	SlippageTolerance  decimal.Decimal `json:"slippage_tolerance,omitempty"`
}

type AMMWithdrawMetadata struct {
	SingleAssetWithdraw bool            `json:"single_asset_withdraw"`
	WithdrawAsset1      decimal.Decimal `json:"withdraw_asset1"`
	WithdrawAsset2      decimal.Decimal `json:"withdraw_asset2"`
	LPTokensBurned      decimal.Decimal `json:"lp_tokens_burned"`
	EffectivePrice      decimal.Decimal `json:"effective_price"`
	PriceImpact         decimal.Decimal `json:"price_impact"`
	WithdrawPercentage  decimal.Decimal `json:"withdraw_percentage"`
}

type AMMVoteMetadata struct {
	VotedTradingFee int             `json:"voted_trading_fee"`
	VoteWeight      int             `json:"vote_weight"`
	TotalVotes      int             `json:"total_votes,omitempty"`
	VotePercentage  decimal.Decimal `json:"vote_percentage,omitempty"`
}

type AMMBidMetadata struct {
	BidAmount      decimal.Decimal `json:"bid_amount"`
	BidAsset       XRPAsset        `json:"bid_asset"`
	AuctionSlotNum int             `json:"auction_slot_num"`
	TimeInterval   int             `json:"time_interval"`
	DiscountedFee  int             `json:"discounted_fee"`
	Expiration     time.Time       `json:"expiration,omitempty"`
}

// NewEnhancedAMMMetadataExtractor creates a new enhanced metadata extractor
func NewEnhancedAMMMetadataExtractor(priceService PriceServiceInterface) *EnhancedAMMMetadataExtractor {
	return &EnhancedAMMMetadataExtractor{
		priceService: priceService,
	}
}

// ExtractAMMMetadata extracts comprehensive metadata from an AMM transaction
func (e *EnhancedAMMMetadataExtractor) ExtractAMMMetadata(tx *XRPTransaction) (*AMMTransactionMetadata, error) {
	if !e.isAMMTransaction(tx) {
		return nil, fmt.Errorf("transaction %s is not an AMM transaction", tx.Hash)
	}

	// Extract basic metadata
	metadata := &AMMTransactionMetadata{
		TransactionHash: tx.Hash,
		TransactionType: tx.TransactionType,
		Account:         tx.Account,
		LedgerIndex:     tx.LedgerIndex,
		Timestamp:       time.Unix(int64(tx.Date), 0),
		RawMetadata:     make(map[string]interface{}),
	}

	// Extract fee
	if fee, err := decimal.NewFromString(tx.Fee); err == nil {
		metadata.Fee = fee.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
	}

	// Extract transaction-specific metadata based on type
	switch tx.TransactionType {
	case "AMMCreate":
		if err := e.extractAMMCreateMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMCreate metadata: %v", err)
		}
	case "AMMDeposit":
		if err := e.extractAMMDepositMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMDeposit metadata: %v", err)
		}
	case "AMMWithdraw":
		if err := e.extractAMMWithdrawMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMWithdraw metadata: %v", err)
		}
	case "AMMVote":
		if err := e.extractAMMVoteMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMVote metadata: %v", err)
		}
	case "AMMBid":
		if err := e.extractAMMBidMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMBid metadata: %v", err)
		}
	case "AMMDelete":
		if err := e.extractAMMDeleteMetadata(tx, metadata); err != nil {
			log.Printf("Error extracting AMMDelete metadata: %v", err)
		}
	}

	// Calculate USD values
	e.calculateUSDValues(metadata)

	return metadata, nil
}

// Extract metadata for AMMCreate transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMCreateMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	createData := &AMMCreateMetadata{
		CreatedAt: metadata.Timestamp,
	}

	// Extract assets and amounts from transaction
	asset1, asset2, amount1, amount2, err := e.extractAssetsAndAmounts(tx)
	if err != nil {
		return fmt.Errorf("failed to extract assets and amounts: %v", err)
	}

	metadata.Asset1 = asset1
	metadata.Asset2 = asset2
	metadata.Amount1 = amount1
	metadata.Amount2 = amount2

	createData.InitialAsset1Amount = amount1
	createData.InitialAsset2Amount = amount2

	// Extract trading fee from transaction data
	if tradingFee := e.extractTradingFee(tx); tradingFee > 0 {
		createData.TradingFee = tradingFee
		metadata.TradingFee = tradingFee
	}

	// Calculate initial price
	if !amount2.IsZero() {
		createData.InitialPrice = amount1.Div(amount2)
	}

	// Extract LP tokens from transaction metadata
	if lpTokens := e.extractLPTokenAmount(tx); lpTokens.GreaterThan(decimal.Zero) {
		createData.InitialLPTokens = lpTokens
		metadata.LPTokens = lpTokens
	}

	metadata.CreateData = createData
	return nil
}

// Extract metadata for AMMDeposit transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMDepositMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	depositData := &AMMDepositMetadata{}

	// Extract assets and amounts
	asset1, asset2, amount1, amount2, err := e.extractAssetsAndAmounts(tx)
	if err != nil {
		return fmt.Errorf("failed to extract assets and amounts: %v", err)
	}

	metadata.Asset1 = asset1
	metadata.Asset2 = asset2
	metadata.Amount1 = amount1
	metadata.Amount2 = amount2

	depositData.DepositAsset1 = amount1
	depositData.DepositAsset2 = amount2

	// Determine if single asset deposit
	depositData.SingleAssetDeposit = (amount1.IsZero()) != (amount2.IsZero())

	// Extract LP tokens received
	if lpTokens := e.extractLPTokenAmount(tx); lpTokens.GreaterThan(decimal.Zero) {
		depositData.LPTokensReceived = lpTokens
		metadata.LPTokens = lpTokens
	}

	// Calculate effective price and price impact
	if !amount2.IsZero() {
		depositData.EffectivePrice = amount1.Div(amount2)
		depositData.PriceImpact = e.calculatePriceImpact(asset1, asset2, amount1, amount2)
	}

	metadata.DepositData = depositData
	return nil
}

// Extract metadata for AMMWithdraw transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMWithdrawMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	withdrawData := &AMMWithdrawMetadata{}

	// Extract assets and amounts
	asset1, asset2, amount1, amount2, err := e.extractAssetsAndAmounts(tx)
	if err != nil {
		return fmt.Errorf("failed to extract assets and amounts: %v", err)
	}

	metadata.Asset1 = asset1
	metadata.Asset2 = asset2
	metadata.Amount1 = amount1
	metadata.Amount2 = amount2

	withdrawData.WithdrawAsset1 = amount1
	withdrawData.WithdrawAsset2 = amount2

	// Determine if single asset withdraw
	withdrawData.SingleAssetWithdraw = (amount1.IsZero()) != (amount2.IsZero())

	// Extract LP tokens burned
	if lpTokens := e.extractLPTokenAmount(tx); lpTokens.GreaterThan(decimal.Zero) {
		withdrawData.LPTokensBurned = lpTokens
		metadata.LPTokens = lpTokens
	}

	// Calculate effective price and price impact
	if !amount2.IsZero() {
		withdrawData.EffectivePrice = amount1.Div(amount2)
		withdrawData.PriceImpact = e.calculatePriceImpact(asset1, asset2, amount1, amount2)
	}

	metadata.WithdrawData = withdrawData
	return nil
}

// Extract metadata for AMMVote transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMVoteMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	voteData := &AMMVoteMetadata{}

	// Extract voting data from transaction
	if voteFee := e.extractVotedTradingFee(tx); voteFee > 0 {
		voteData.VotedTradingFee = voteFee
	}

	if voteWeight := e.extractVoteWeight(tx); voteWeight > 0 {
		voteData.VoteWeight = voteWeight
	}

	metadata.VoteData = voteData
	return nil
}

// Extract metadata for AMMBid transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMBidMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	bidData := &AMMBidMetadata{}

	// Extract bid amount and asset
	asset, amount, err := e.extractBidAssetAndAmount(tx)
	if err != nil {
		return fmt.Errorf("failed to extract bid asset and amount: %v", err)
	}

	bidData.BidAsset = asset
	bidData.BidAmount = amount

	// Extract auction slot details
	bidData.AuctionSlotNum = e.extractAuctionSlotNumber(tx)
	bidData.TimeInterval = e.extractTimeInterval(tx)
	bidData.DiscountedFee = e.extractDiscountedFee(tx)

	metadata.BidData = bidData
	return nil
}

// Extract metadata for AMMDelete transactions
func (e *EnhancedAMMMetadataExtractor) extractAMMDeleteMetadata(tx *XRPTransaction, metadata *AMMTransactionMetadata) error {
	// Extract basic AMM account information
	if ammAccount := e.extractAMMAccount(tx); ammAccount != "" {
		metadata.AMMAccount = ammAccount
	}

	return nil
}

// Helper methods for extracting specific data from transactions

func (e *EnhancedAMMMetadataExtractor) extractAssetsAndAmounts(tx *XRPTransaction) (XRPAsset, XRPAsset, decimal.Decimal, decimal.Decimal, error) {
	var asset1, asset2 XRPAsset
	var amount1, amount2 decimal.Decimal

	// Extract from Amount field
	if tx.Amount != nil {
		asset1, amount1 = e.parseAssetAmount(tx.Amount)
	}

	// Extract from Amount2 field if it exists in metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if amount2Field, exists := meta["Amount2"]; exists {
			asset2, amount2 = e.parseAssetAmount(amount2Field)
		}
	}

	// Also check in the main transaction fields
	if asset1.Currency == "" && tx.TakerGets != nil {
		asset1, amount1 = e.parseAssetAmount(tx.TakerGets)
	}
	if asset2.Currency == "" && tx.TakerPays != nil {
		asset2, amount2 = e.parseAssetAmount(tx.TakerPays)
	}

	return asset1, asset2, amount1, amount2, nil
}

func (e *EnhancedAMMMetadataExtractor) parseAssetAmount(field interface{}) (XRPAsset, decimal.Decimal) {
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

func (e *EnhancedAMMMetadataExtractor) extractTradingFee(tx *XRPTransaction) int {
	// Look for TradingFee in transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if fee, exists := meta["TradingFee"]; exists {
			if feeInt, ok := fee.(int); ok {
				return feeInt
			}
			if feeFloat, ok := fee.(float64); ok {
				return int(feeFloat)
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractLPTokenAmount(tx *XRPTransaction) decimal.Decimal {
	// Look for LP token amounts in transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if lpToken, exists := meta["LPToken"]; exists {
			if lpMap, ok := lpToken.(map[string]interface{}); ok {
				if value, ok := lpMap["value"].(string); ok {
					if amount, err := decimal.NewFromString(value); err == nil {
						return amount
					}
				}
			}
		}
	}
	return decimal.Zero
}

func (e *EnhancedAMMMetadataExtractor) extractVotedTradingFee(tx *XRPTransaction) int {
	// Extract voted trading fee from transaction
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if fee, exists := meta["VotedTradingFee"]; exists {
			if feeInt, ok := fee.(int); ok {
				return feeInt
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractVoteWeight(tx *XRPTransaction) int {
	// Extract vote weight from transaction
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if weight, exists := meta["VoteWeight"]; exists {
			if weightInt, ok := weight.(int); ok {
				return weightInt
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractBidAssetAndAmount(tx *XRPTransaction) (XRPAsset, decimal.Decimal, error) {
	// Extract bid asset and amount from transaction
	if tx.Amount != nil {
		asset, amount := e.parseAssetAmount(tx.Amount)
		return asset, amount, nil
	}
	return XRPAsset{}, decimal.Zero, fmt.Errorf("no bid amount found")
}

func (e *EnhancedAMMMetadataExtractor) extractAuctionSlotNumber(tx *XRPTransaction) int {
	// Extract auction slot number from transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if slot, exists := meta["AuctionSlot"]; exists {
			if slotInt, ok := slot.(int); ok {
				return slotInt
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractTimeInterval(tx *XRPTransaction) int {
	// Extract time interval from transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if interval, exists := meta["TimeInterval"]; exists {
			if intervalInt, ok := interval.(int); ok {
				return intervalInt
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractDiscountedFee(tx *XRPTransaction) int {
	// Extract discounted fee from transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if fee, exists := meta["DiscountedFee"]; exists {
			if feeInt, ok := fee.(int); ok {
				return feeInt
			}
		}
	}
	return 0
}

func (e *EnhancedAMMMetadataExtractor) extractAMMAccount(tx *XRPTransaction) string {
	// Extract AMM account from transaction metadata
	if meta, ok := tx.Meta.(map[string]interface{}); ok {
		if account, exists := meta["AMMAccount"]; exists {
			if accountStr, ok := account.(string); ok {
				return accountStr
			}
		}
	}
	return ""
}

func (e *EnhancedAMMMetadataExtractor) calculatePriceImpact(asset1, asset2 XRPAsset, amount1, amount2 decimal.Decimal) decimal.Decimal {
	// Simplified price impact calculation
	if amount1.IsZero() || amount2.IsZero() {
		return decimal.Zero
	}

	// Simple heuristic: larger trades relative to pool size have higher impact
	totalValue := amount1.Add(amount2) // Simplified
	if totalValue.GreaterThan(decimal.NewFromInt(10000)) {
		return decimal.NewFromFloat(0.005) // 0.5% impact for large trades
	} else if totalValue.GreaterThan(decimal.NewFromInt(1000)) {
		return decimal.NewFromFloat(0.001) // 0.1% impact for medium trades
	}
	return decimal.NewFromFloat(0.0001) // 0.01% impact for small trades
}

func (e *EnhancedAMMMetadataExtractor) calculateUSDValues(metadata *AMMTransactionMetadata) {
	if e.priceService == nil {
		return
	}

	// CRITICAL FIX: Use batch pricing to avoid redundant XRP price calls
	assets := []XRPAsset{}
	if metadata.Asset1.Currency != "" {
		assets = append(assets, metadata.Asset1)
	}
	if metadata.Asset2.Currency != "" {
		assets = append(assets, metadata.Asset2)
	}

	if len(assets) > 0 {
		prices, err := e.priceService.GetMultipleAssetPrices(assets)
		if err == nil {
			if metadata.Asset1.Currency != "" {
				key := fmt.Sprintf("%s_%s", metadata.Asset1.Currency, metadata.Asset1.Issuer)
				if price, exists := prices[key]; exists {
					metadata.Asset1PriceUSD = price
				}
			}
			if metadata.Asset2.Currency != "" {
				key := fmt.Sprintf("%s_%s", metadata.Asset2.Currency, metadata.Asset2.Issuer)
				if price, exists := prices[key]; exists {
					metadata.Asset2PriceUSD = price
				}
			}
		}
	}

	// Calculate total USD value
	metadata.TotalValueUSD = metadata.Amount1.Mul(metadata.Asset1PriceUSD).Add(metadata.Amount2.Mul(metadata.Asset2PriceUSD))
}

func (e *EnhancedAMMMetadataExtractor) isAMMTransaction(tx *XRPTransaction) bool {
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

// ExtractMetadataFromLedger extracts AMM metadata from all transactions in a ledger
func (e *EnhancedAMMMetadataExtractor) ExtractMetadataFromLedger(ledger *XRPLedger) ([]*AMMTransactionMetadata, error) {
	var ammMetadata []*AMMTransactionMetadata

	if ledger.Ledger == nil || ledger.Ledger.Transactions == nil {
		return ammMetadata, nil
	}

	for _, tx := range ledger.Ledger.Transactions {
		if e.isAMMTransaction(tx) {
			metadata, err := e.ExtractAMMMetadata(tx)
			if err != nil {
				log.Printf("Error extracting metadata for transaction %s: %v", tx.Hash, err)
				continue
			}
			ammMetadata = append(ammMetadata, metadata)
		}
	}

	return ammMetadata, nil
}

// GetTransactionSummary provides a human-readable summary of an AMM transaction
func (e *EnhancedAMMMetadataExtractor) GetTransactionSummary(metadata *AMMTransactionMetadata) string {
	switch metadata.TransactionType {
	case "AMMCreate":
		return fmt.Sprintf("Created AMM pool with %s %s and %s %s",
			formatAmount(metadata.Amount1), metadata.Asset1.Currency,
			formatAmount(metadata.Amount2), metadata.Asset2.Currency)
	case "AMMDeposit":
		if metadata.DepositData.SingleAssetDeposit {
			if metadata.Amount1.GreaterThan(decimal.Zero) {
				return fmt.Sprintf("Deposited %s %s (single asset)",
					formatAmount(metadata.Amount1), metadata.Asset1.Currency)
			} else {
				return fmt.Sprintf("Deposited %s %s (single asset)",
					formatAmount(metadata.Amount2), metadata.Asset2.Currency)
			}
		}
		return fmt.Sprintf("Deposited %s %s + %s %s",
			formatAmount(metadata.Amount1), metadata.Asset1.Currency,
			formatAmount(metadata.Amount2), metadata.Asset2.Currency)
	case "AMMWithdraw":
		if metadata.WithdrawData.SingleAssetWithdraw {
			if metadata.Amount1.GreaterThan(decimal.Zero) {
				return fmt.Sprintf("Withdrew %s %s (single asset)",
					formatAmount(metadata.Amount1), metadata.Asset1.Currency)
			} else {
				return fmt.Sprintf("Withdrew %s %s (single asset)",
					formatAmount(metadata.Amount2), metadata.Asset2.Currency)
			}
		}
		return fmt.Sprintf("Withdrew %s %s + %s %s",
			formatAmount(metadata.Amount1), metadata.Asset1.Currency,
			formatAmount(metadata.Amount2), metadata.Asset2.Currency)
	case "AMMVote":
		return fmt.Sprintf("Voted for %d trading fee", metadata.VoteData.VotedTradingFee)
	case "AMMBid":
		return fmt.Sprintf("Bid %s %s for auction slot",
			formatAmount(metadata.BidData.BidAmount), metadata.BidData.BidAsset.Currency)
	case "AMMDelete":
		return "Deleted AMM pool"
	default:
		return fmt.Sprintf("AMM transaction: %s", metadata.TransactionType)
	}
}

func formatAmount(amount decimal.Decimal) string {
	if amount.GreaterThanOrEqual(decimal.NewFromInt(1000000)) {
		return fmt.Sprintf("%sM", amount.Div(decimal.NewFromInt(1000000)).StringFixed(2))
	} else if amount.GreaterThanOrEqual(decimal.NewFromInt(1000)) {
		return fmt.Sprintf("%sK", amount.Div(decimal.NewFromInt(1000)).StringFixed(2))
	}
	return amount.StringFixed(6)
}
