package xrp

import (
	"time"

	"gorm.io/gorm"
)

type AMMLiquidity struct {
	ID                uint    `gorm:"primaryKey"`
	PoolID            string  `gorm:"index;not null"`
	Asset1Currency    string  `gorm:"not null"`
	Asset1Issuer      string  `gorm:"not null"`
	Asset1Balance     float64 `gorm:"not null"`
	Asset2Currency    string  `gorm:"not null"`
	Asset2Issuer      string  `gorm:"not null"`
	Asset2Balance     float64 `gorm:"not null"`
	Asset1ValueUsd    float64 `gorm:"not null"`
	Asset2ValueUsd    float64 `gorm:"not null"`
	TotalLiquidityUsd float64 `gorm:"not null"`
	Fee               float64 `gorm:"not null"`
	CreatedAt         int64   `gorm:"not null"`
	UpdatedAt         time.Time
}

func AutoMigrateAMMLiquidity(db *gorm.DB) error {
	return db.AutoMigrate(&AMMLiquidity{})
}
