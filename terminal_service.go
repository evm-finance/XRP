package xrp

import (
	"fmt"
	"qc-defi-graphql-server/internal/models"
)

// TerminalService provides data aggregation for the Terminal page.
// It relies on the shared MarketDataClient so it stays consistent with Screener.
type TerminalService struct {
	marketData MarketDataClient
}

func NewTerminalService(client MarketDataClient) *TerminalService {
	return &TerminalService{marketData: client}
}

// GetTerminalOverview returns a concise dataset for the Terminal page,
// using the same external API sources as the screener.
func (s *TerminalService) GetTerminalOverview(limit int) ([]*models.XRPTokenFields, error) {
	if s.marketData == nil {
		return nil, fmt.Errorf("market data client not configured")
	}
	return s.marketData.GetTokenList(limit)
}

// GetXRPTerminalOverview fetches detailed XRP overview for Terminal page widgets
func (s *TerminalService) GetXRPTerminalOverview() (*XRPOverview, error) {
	if s.marketData == nil {
		return nil, fmt.Errorf("market data client not configured")
	}
	return s.marketData.GetXRPOverview()
}
