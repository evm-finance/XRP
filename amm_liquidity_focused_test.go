package xrp

import (
	"testing"

	"github.com/shopspring/decimal"
	"qc-defi-graphql-server/internal/database"
)

func TestAMMLiquidityCalculationOnly(t *testing.T) {
	// Single database connection for this test
	dbService := database.NewDatabaseCreateService()
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	// Create liquidity service with shared database connection
	liquidityService := NewAMMLiquidityServiceWithDB(db)

	t.Log("🧪 Running AMM liquidity calculation for all pools...")

	// Update all pool liquidity
	err = liquidityService.UpdateAllPoolLiquidity()
	if err != nil {
		t.Fatalf("❌ Failed to update all pool liquidity: %v", err)
	}

	t.Log("✅ AMM liquidity calculation completed successfully!")

	// Verify the results
	calculations, err := liquidityService.GetAllPoolLiquidity()
	if err != nil {
		t.Fatalf("❌ Failed to get updated liquidity data: %v", err)
	}

	t.Logf("📊 Updated liquidity data for %d pools", len(calculations))

	// Show summary
	var totalLiquidity decimal.Decimal
	for _, calc := range calculations {
		totalLiquidity = totalLiquidity.Add(calc.TotalLiquidityUSD)
	}
	t.Logf("💰 Total liquidity across all pools: $%s USD", totalLiquidity.StringFixed(2))

	// Show details for first few pools
	for i, calc := range calculations {
		if i < 5 { // Show first 5 pools
			t.Logf("   Pool %d: %s - $%s USD (%s%% XRP, %s%% %s)",
				i+1, calc.PoolID, calc.TotalLiquidityUSD.StringFixed(2),
				calc.Asset1Percentage.StringFixed(1), calc.Asset2Percentage.StringFixed(1), calc.Asset2Currency)
		}
	}

	t.Log("✅ AMM liquidity calculation test completed!")
}
