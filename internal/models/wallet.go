package models

import (
	"time"
)

type Wallet struct {
	ID      int64      `json:"id"`
	UserID  string     `json:"user_id"`
	Balance string    `json:"balance"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	WalletID  int64     `json:"wallet_id"`
	Amount    string   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	MaxBalanceUnidentified = 1000000
	MaxBalanceIdentified   = 10000000
)