package models

import (
	"time"
)

type Wallet struct {
	ID      int64      `json:"id"`
	UserID  int     `json:"user_id"`
	Balance string    `json:"balance"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	WalletID  int64     `json:"wallet_id"`
	Amount    string   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type TopUpRequest struct {
	WalletID int  `json:"wallet_id" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

const (
	MaxBalanceUnidentified = 1000000
	MaxBalanceIdentified   = 10000000
)