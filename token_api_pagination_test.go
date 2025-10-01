package xrp_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type TokenList struct {
	Count  int64       `json:"count"`
	Tokens []TokenData `json:"tokens"`
}

type TokenData struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Meta     struct {
		Token struct {
			Name string `json:"name"`
		} `json:"token"`
	} `json:"meta"`
}

func countUniqueTokens(tokens []TokenData) int {
	unique := make(map[string]struct{})
	for _, tok := range tokens {
		key := tok.Currency + ":" + tok.Issuer
		unique[key] = struct{}{}
	}
	return len(unique)
}

func TestFetchAndCountUniqueTokens(t *testing.T) {
	testCases := []int{1000, 10000}
	for _, limit := range testCases {
		url := fmt.Sprintf("https://s1.xrplmeta.org/tokens?page=1&limit=%d", limit)
		t.Logf("Fetching tokens from API: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Failed to fetch tokens: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		var tokenList TokenList
		if err := json.NewDecoder(resp.Body).Decode(&tokenList); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("Total tokens returned (limit %d): %d", limit, len(tokenList.Tokens))
		uniqueCount := countUniqueTokens(tokenList.Tokens)
		t.Logf("Unique (currency, issuer) pairs (limit %d): %d", limit, uniqueCount)

		// Show first 5 tokens
		t.Log("First 5 tokens:")
		for i := 0; i < 5 && i < len(tokenList.Tokens); i++ {
			tok := tokenList.Tokens[i]
			t.Logf("  %d. %s (%s) - Issuer: %s", i+1, tok.Meta.Token.Name, tok.Currency, tok.Issuer)
		}

		// Show last 5 tokens if we have more than 5
		if len(tokenList.Tokens) > 5 {
			t.Log("Last 5 tokens:")
			start := len(tokenList.Tokens) - 5
			for i := start; i < len(tokenList.Tokens); i++ {
				tok := tokenList.Tokens[i]
				t.Logf("  %d. %s (%s) - Issuer: %s", i+1, tok.Meta.Token.Name, tok.Currency, tok.Issuer)
			}
		}
	}
}
