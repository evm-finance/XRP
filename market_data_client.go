package xrp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"qc-defi-graphql-server/internal/models"
)

// MarketDataClient defines a shared interface for fetching market/token data
// used by both the Screener and Terminal features. Implementations should
// aggregate data from QuantifyCrypto and XRPL Meta (or other sources).
type MarketDataClient interface {
	// GetTokenList returns up to N tokens including XRP as the first entry
	// when QuantifyCrypto data is available. Results should be mapped to
	// models.XRPTokenFields and generally sorted by marketcap desc.
	GetTokenList(limit int) ([]*models.XRPTokenFields, error)
	GetXRPOverview() (*XRPOverview, error)
}

// QuantifyCryptoXRPLClient is the default implementation that combines
// QuantifyCrypto (for XRP) with XRPL Meta tokens.
type QuantifyCryptoXRPLClient struct {
	httpClient  *http.Client
	qcAccessKey string
	qcSecretKey string
}

// NewMarketDataClient constructs a MarketDataClient. QC keys are optional; if
// not provided, the client will skip the QuantifyCrypto XRP augmentation.
func NewMarketDataClient(qcAccessKey, qcSecretKey string) MarketDataClient {
	return &QuantifyCryptoXRPLClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		qcAccessKey: qcAccessKey,
		qcSecretKey: qcSecretKey,
	}
}

// QuantifyCryptoResponse represents the response from QuantifyCrypto API
type QuantifyCryptoResponse struct {
	Status struct {
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
	} `json:"status"`
	Data struct {
		ID                string  `json:"id"`
		Rank              int     `json:"rank"`
		CoinSymbol        string  `json:"coin_symbol"`
		CoinName          string  `json:"coin_name"`
		Marketcap         float64 `json:"marketCap"`
		CirculatingSupply float64 `json:"circulating_supply"`
		CoinPrice         float64 `json:"coin_price"`
	} `json:"data"`
	Currency struct {
		ConvertRate float64 `json:"convert_rate"`
		Code        string  `json:"code"`
		Locale      string  `json:"locale"`
	} `json:"currency"`
}

func (c *QuantifyCryptoXRPLClient) GetTokenList(limit int) ([]*models.XRPTokenFields, error) {
	if limit <= 0 {
		limit = 100
	}

	// Try to augment with XRP first (optional)
	var result []*models.XRPTokenFields
	if c.qcAccessKey != "" && c.qcSecretKey != "" {
		if xrpToken, err := c.getXRPFromQuantifyCrypto(); err != nil {
			log.Printf("⚠️ [MARKET DATA] Failed to fetch XRP from QuantifyCrypto: %v", err)
		} else if xrpToken != nil {
			result = append(result, xrpToken)
		}
	}

	// Fetch XRPL tokens from XRPL Meta API
	tokens, err := c.getXRPLMetaTokens(limit - len(result))
	if err != nil {
		return nil, err
	}
	result = append(result, tokens...)

	// Sort by marketcap desc
	sort.Slice(result, func(i, j int) bool {
		return result[i].Marketcap > result[j].Marketcap
	})

	return result, nil
}

func (c *QuantifyCryptoXRPLClient) getXRPFromQuantifyCrypto() (*models.XRPTokenFields, error) {
	req, err := http.NewRequest("GET", "https://quantifycrypto.com/api/v1/coins/XRP?currency=USD&include_signals=true&signal_type=trend", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create QC request: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("QC-Access-Key", c.qcAccessKey)
	req.Header.Set("QC-Secret-Key", c.qcSecretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch XRP from QC: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("QuantifyCrypto API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read QC response body: %w", err)
	}

	var qcResp QuantifyCryptoResponse
	if err := json.Unmarshal(body, &qcResp); err != nil {
		return nil, fmt.Errorf("failed to parse QC XRP JSON: %w", err)
	}

	return &models.XRPTokenFields{
		Currency:      "XRP",
		IssuerAddress: "",
		TokenName:     "XRP",
		IssuerName:    "XRP Ledger",
		Icon:          "https://xumm.app/assets/icons/currencies/XRP.png",
		Marketcap:     qcResp.Data.Marketcap,
		Price:         qcResp.Data.CoinPrice,
		Supply:        qcResp.Data.CirculatingSupply,
		Liquidity:     0,
		Volume24H:     0,
	}, nil
}

func (c *QuantifyCryptoXRPLClient) getXRPLMetaTokens(limit int) ([]*models.XRPTokenFields, error) {
	if limit <= 0 {
		limit = 100
	}
	resp, err := c.httpClient.Get("https://s1.xrplmeta.org/tokens")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch XRPL Meta tokens: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("XRPL Meta API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XRPL Meta response: %w", err)
	}

	var tokenList models.XRPTokenList
	if err := json.Unmarshal(body, &tokenList); err != nil {
		return nil, fmt.Errorf("failed to parse XRPL Meta JSON: %w", err)
	}

	var out []*models.XRPTokenFields
	count := 0
	for _, t := range tokenList.Tokens {
		if t.Meta.Token.Name == "" {
			continue
		}
		if count >= limit {
			break
		}
		out = append(out, &models.XRPTokenFields{
			Currency:      t.Currency,
			IssuerAddress: t.Issuer,
			TokenName:     t.Meta.Token.Name,
			IssuerName:    t.Meta.Issuer.Name,
			Icon:          t.Meta.Token.Icon,
			Marketcap:     t.Metrics.Marketcap,
			Price:         t.Metrics.Price,
			Supply:        0,
			Liquidity:     t.LiquidityUSD,
			Volume24H:     t.Metrics.Volume24H,
		})
		count++
	}
	return out, nil
}

// GetXRPOverview queries QuantifyCrypto v1 coins/XRP with signals to map terminal fields
func (c *QuantifyCryptoXRPLClient) GetXRPOverview() (*XRPOverview, error) {
	if c.qcAccessKey == "" || c.qcSecretKey == "" {
		return nil, fmt.Errorf("QuantifyCrypto credentials not configured")
	}

	req, err := http.NewRequest("GET", "https://quantifycrypto.com/api/v1/coins/XRP?currency=USD&include_signals=true&signal_type=trend", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create QC request: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("QC-Access-Key", c.qcAccessKey)
	req.Header.Set("QC-Secret-Key", c.qcSecretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch QC XRP overview: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("QuantifyCrypto API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read QC response: %w", err)
	}

	var qc struct {
		Data struct {
			CoinPrice         float64            `json:"coin_price"`
			PriceBtc          float64            `json:"price_btc"`
			Change24hPct      float64            `json:"percent_change_24h"`
			High24h           float64            `json:"high_24h"`
			Low24h            float64            `json:"low_24h"`
			Marketcap         float64            `json:"marketCap"`
			CirculatingSupply float64            `json:"circulating_supply"`
			Volume24h         float64            `json:"volume_24h"`
			Support           float64            `json:"support"`
			Resistance        float64            `json:"resistance"`
			SafetyScore       float64            `json:"safety_score"`
			TechnicalScore    float64            `json:"technical_score"`
			Signals           map[string]float64 `json:"signals"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &qc); err != nil {
		return nil, fmt.Errorf("failed to parse QC XRP overview JSON: %w", err)
	}

	overview := &XRPOverview{
		PriceUSD:          qc.Data.CoinPrice,
		PriceBTC:          qc.Data.PriceBtc,
		Change24hPct:      qc.Data.Change24hPct,
		High24hUSD:        qc.Data.High24h,
		Low24hUSD:         qc.Data.Low24h,
		MarketCapUSD:      qc.Data.Marketcap,
		Volume24hUSD:      qc.Data.Volume24h,
		CirculatingSupply: qc.Data.CirculatingSupply,
		SupportUSD:        qc.Data.Support,
		ResistanceUSD:     qc.Data.Resistance,
		SafetyScore:       qc.Data.SafetyScore,
		TechnicalScore:    qc.Data.TechnicalScore,
		TrendScores:       qc.Data.Signals,
	}
	return overview, nil
}
