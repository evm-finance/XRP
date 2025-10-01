package xrp

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// AMMValidationError represents validation errors with detailed context
type AMMValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *AMMValidationError) Error() string {
	return fmt.Sprintf("AMM validation error [%s]: %s (field: %s, value: %s)", 
		e.Code, e.Message, e.Field, e.Value)
}

// AMMValidationResult contains validation results and any errors found
type AMMValidationResult struct {
	IsValid       bool                 `json:"is_valid"`
	Errors        []*AMMValidationError `json:"errors"`
	Warnings      []*AMMValidationError `json:"warnings"`
	CleanedData   map[string]interface{} `json:"cleaned_data,omitempty"`
}

// ValidatedAMMData represents cleaned and validated AMM pool data
type ValidatedAMMData struct {
	Account           string    `json:"account"`
	Amount            float64   `json:"amount"`
	AmountCurrency    string    `json:"amount_currency"`
	AmountIssuer      string    `json:"amount_issuer"`
	Amount2Value      float64   `json:"amount2_value"`
	Amount2Currency   string    `json:"amount2_currency"`
	Amount2Issuer     string    `json:"amount2_issuer"`
	Asset2Frozen      bool      `json:"asset2_frozen"`
	LPTokenValue      float64   `json:"lp_token_value"`
	LPTokenCurrency   string    `json:"lp_token_currency"`
	LPTokenIssuer     string    `json:"lp_token_issuer"`
	TradingFeeBPS     int       `json:"trading_fee_bps"`
	LedgerIndex       int       `json:"ledger_index"`
	Validated         bool      `json:"validated"`
	LastAPIFetch      time.Time `json:"last_api_fetch"`
}

// AMMValidator provides comprehensive validation for AMM pool data
type AMMValidator struct {
	MaxTradingFeeBPS    int
	MinPoolValue        float64
	RequireValidation   bool
}

// NewAMMValidator creates a new validator with default settings
func NewAMMValidator() *AMMValidator {
	return &AMMValidator{
		MaxTradingFeeBPS:  1000, // 10% max trading fee
		MinPoolValue:      0.01, // Minimum 0.01 for amounts
		RequireValidation: true,
	}
}

// ValidateAMMInfoResponse validates a complete amm_info API response
func (v *AMMValidator) ValidateAMMInfoResponse(response interface{}) (*AMMValidationResult, error) {
	result := &AMMValidationResult{
		IsValid:     false,
		Errors:      []*AMMValidationError{},
		Warnings:    []*AMMValidationError{},
		CleanedData: make(map[string]interface{}),
	}

	// Parse response structure
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "response",
			Value:   "invalid_type",
			Message: "Response is not a valid map structure",
			Code:    "INVALID_RESPONSE_TYPE",
		})
		return result, nil
	}

	// Validate response status
	if err := v.validateResponseStatus(responseMap, result); err != nil {
		return result, err
	}

	// Extract and validate result structure
	resultData, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "result",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid result field in response",
			Code:    "MISSING_RESULT",
		})
		return result, nil
	}

	// Extract and validate AMM data
	ammData, ok := resultData["amm"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amm",
			Value:   "missing_or_invalid", 
			Message: "Missing or invalid amm field in result",
			Code:    "MISSING_AMM_DATA",
		})
		return result, nil
	}

	// Validate individual AMM fields
	validatedData, err := v.validateAMMData(ammData, resultData, result)
	if err != nil {
		return result, err
	}

	// Set validation result
	result.IsValid = len(result.Errors) == 0
	if result.IsValid {
		// Convert validated data to map for compatibility
		result.CleanedData = v.validatedDataToMap(validatedData)
	}

	return result, nil
}

// validateResponseStatus checks the basic response structure and status
func (v *AMMValidator) validateResponseStatus(response map[string]interface{}, result *AMMValidationResult) error {
	status, ok := response["status"].(string)
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "status",
			Value:   "missing",
			Message: "Missing status field in response",
			Code:    "MISSING_STATUS",
		})
		return nil
	}

	if status != "success" {
		errorMsg := "Unknown error"
		if errorCode, exists := response["error"].(string); exists {
			errorMsg = errorCode
		}
		
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "status",
			Value:   status,
			Message: fmt.Sprintf("API returned error status: %s", errorMsg),
			Code:    "API_ERROR_STATUS",
		})
		return nil
	}

	return nil
}

// validateAMMData validates the core AMM pool data
func (v *AMMValidator) validateAMMData(ammData map[string]interface{}, resultData map[string]interface{}, result *AMMValidationResult) (*ValidatedAMMData, error) {
	validated := &ValidatedAMMData{
		LastAPIFetch: time.Now(),
	}

	// Validate account (pool ID)
	if err := v.validateAccount(ammData, validated, result); err != nil {
		return nil, err
	}

	// Validate amounts
	if err := v.validateAmounts(ammData, validated, result); err != nil {
		return nil, err
	}

	// Validate trading fee
	if err := v.validateTradingFee(ammData, validated, result); err != nil {
		return nil, err
	}

	// Validate LP token data
	if err := v.validateLPToken(ammData, validated, result); err != nil {
		return nil, err
	}

	// Validate ledger metadata
	if err := v.validateLedgerMetadata(resultData, validated, result); err != nil {
		return nil, err
	}

	return validated, nil
}

// validateAccount validates the AMM pool account field
func (v *AMMValidator) validateAccount(ammData map[string]interface{}, validated *ValidatedAMMData, result *AMMValidationResult) error {
	account, ok := ammData["account"].(string)
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "account",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid account field",
			Code:    "INVALID_ACCOUNT",
		})
		return nil
	}

	// Validate XRPL address format
	xrplAddressRegex := regexp.MustCompile(`^r[a-zA-Z0-9]{24,34}$`)
	if !xrplAddressRegex.MatchString(account) {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "account",
			Value:   account,
			Message: "Invalid XRPL address format",
			Code:    "INVALID_ADDRESS_FORMAT",
		})
		return nil
	}

	validated.Account = account
	return nil
}

// validateAmounts validates pool asset amounts
func (v *AMMValidator) validateAmounts(ammData map[string]interface{}, validated *ValidatedAMMData, result *AMMValidationResult) error {
	// Validate Asset 1 (amount)
	amount, ok := ammData["amount"].(string)
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid amount field",
			Code:    "INVALID_AMOUNT",
		})
		return nil
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount",
			Value:   amount,
			Message: "Cannot parse amount as number",
			Code:    "INVALID_AMOUNT_FORMAT",
		})
		return nil
	}

	// Convert XRP drops to XRP if needed
	// Note: Conversion disabled pending verification of data source format
	// Large amounts (>1M) may indicate drops format requiring conversion
	// if amountFloat > 1000000 {
	//	amountFloat = amountFloat / 1000000.0
	// }

	if amountFloat < v.MinPoolValue {
		result.Warnings = append(result.Warnings, &AMMValidationError{
			Field:   "amount",
			Value:   fmt.Sprintf("%.8f", amountFloat),
			Message: "Amount is very small, may indicate data quality issues",
			Code:    "SMALL_AMOUNT_WARNING",
		})
	}

	validated.Amount = amountFloat
	validated.AmountCurrency = "XRP"
	validated.AmountIssuer = ""

	// Validate Asset 2 (amount2)
	amount2, ok := ammData["amount2"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount2",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid amount2 field",
			Code:    "INVALID_AMOUNT2",
		})
		return nil
	}

	currency2, ok := amount2["currency"].(string)
	if !ok || currency2 == "" {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount2.currency",
			Value:   "missing_or_empty",
			Message: "Missing or empty amount2 currency",
			Code:    "INVALID_CURRENCY",
		})
		return nil
	}

	issuer2, ok := amount2["issuer"].(string)
	if !ok || issuer2 == "" {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount2.issuer",
			Value:   "missing_or_empty",
			Message: "Missing or empty amount2 issuer",
			Code:    "INVALID_ISSUER",
		})
		return nil
	}

	value2, ok := amount2["value"].(string)
	if !ok {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount2.value",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid amount2 value",
			Code:    "INVALID_AMOUNT2_VALUE",
		})
		return nil
	}

	value2Float, err := strconv.ParseFloat(value2, 64)
	if err != nil {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "amount2.value",
			Value:   value2,
			Message: "Cannot parse amount2 value as number",
			Code:    "INVALID_AMOUNT2_FORMAT",
		})
		return nil
	}

	validated.Amount2Value = value2Float
	validated.Amount2Currency = currency2
	validated.Amount2Issuer = issuer2

	if frozen, ok := ammData["asset2_frozen"].(bool); ok {
		validated.Asset2Frozen = frozen
	}

	return nil
}

// validateTradingFee validates the trading fee field
func (v *AMMValidator) validateTradingFee(ammData map[string]interface{}, validated *ValidatedAMMData, result *AMMValidationResult) error {
	tradingFee, ok := ammData["trading_fee"].(float64)
	if !ok {
		if feeInt, ok := ammData["trading_fee"].(int); ok {
			tradingFee = float64(feeInt)
		} else {
			result.Warnings = append(result.Warnings, &AMMValidationError{
				Field:   "trading_fee",
				Value:   "missing_or_invalid",
				Message: "Missing or invalid trading fee, defaulting to 0",
				Code:    "MISSING_TRADING_FEE",
			})
			validated.TradingFeeBPS = 0
			return nil
		}
	}

	tradingFeeBPS := int(tradingFee)
	if tradingFeeBPS < 0 || tradingFeeBPS > v.MaxTradingFeeBPS {
		result.Errors = append(result.Errors, &AMMValidationError{
			Field:   "trading_fee",
			Value:   fmt.Sprintf("%d", tradingFeeBPS),
			Message: fmt.Sprintf("Trading fee out of valid range (0-%d)", v.MaxTradingFeeBPS),
			Code:    "INVALID_TRADING_FEE_RANGE",
		})
		return nil
	}

	validated.TradingFeeBPS = tradingFeeBPS
	return nil
}

// validateLPToken validates LP token data
func (v *AMMValidator) validateLPToken(ammData map[string]interface{}, validated *ValidatedAMMData, result *AMMValidationResult) error {
	lpToken, ok := ammData["lp_token"].(map[string]interface{})
	if !ok {
		result.Warnings = append(result.Warnings, &AMMValidationError{
			Field:   "lp_token",
			Value:   "missing",
			Message: "Missing LP token data",
			Code:    "MISSING_LP_TOKEN",
		})
		return nil
	}

	if value, ok := lpToken["value"].(string); ok {
		if valueFloat, err := strconv.ParseFloat(value, 64); err == nil {
			validated.LPTokenValue = valueFloat
		}
	}

	if currency, ok := lpToken["currency"].(string); ok {
		validated.LPTokenCurrency = currency
	}

	if issuer, ok := lpToken["issuer"].(string); ok {
		validated.LPTokenIssuer = issuer
	}

	return nil
}

// validateLedgerMetadata validates ledger-related metadata
func (v *AMMValidator) validateLedgerMetadata(resultData map[string]interface{}, validated *ValidatedAMMData, result *AMMValidationResult) error {
	if ledgerIndex, ok := resultData["ledger_current_index"].(float64); ok {
		validated.LedgerIndex = int(ledgerIndex)
	} else if ledgerIndex, ok := resultData["ledger_current_index"].(int); ok {
		validated.LedgerIndex = ledgerIndex
	} else {
		result.Warnings = append(result.Warnings, &AMMValidationError{
			Field:   "ledger_current_index",
			Value:   "missing_or_invalid",
			Message: "Missing or invalid ledger index",
			Code:    "MISSING_LEDGER_INDEX",
		})
	}

	if validated_flag, ok := resultData["validated"].(bool); ok {
		validated.Validated = validated_flag
	} else {
		validated.Validated = v.RequireValidation
	}

	return nil
}

// validatedDataToMap converts ValidatedAMMData to a map for compatibility
func (v *AMMValidator) validatedDataToMap(data *ValidatedAMMData) map[string]interface{} {
	return map[string]interface{}{
		"account":            data.Account,
		"amount":             data.Amount,
		"amount_currency":    data.AmountCurrency,
		"amount_issuer":      data.AmountIssuer,
		"amount2_value":      data.Amount2Value,
		"amount2_currency":   data.Amount2Currency,
		"amount2_issuer":     data.Amount2Issuer,
		"asset2_frozen":      data.Asset2Frozen,
		"lp_token_value":     data.LPTokenValue,
		"lp_token_currency":  data.LPTokenCurrency,
		"lp_token_issuer":    data.LPTokenIssuer,
		"trading_fee_bps":    data.TradingFeeBPS,
		"ledger_index":       data.LedgerIndex,
		"validated":          data.Validated,
		"last_api_fetch":     data.LastAPIFetch,
	}
}
