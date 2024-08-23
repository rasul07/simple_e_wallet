package models

type Wallet struct {
	ID      string `db:"id"`
	UserID  string `db:"user_id"`
	Balance int64 `db:"balance"`
}

type TopUpRequest struct {
	WalletID string `json:"wallet_id" binding:"required"`
	Amount   string `json:"amount" binding:"required"`
}

const (
	MaxBalanceUnidentified = 1000000
	MaxBalanceIdentified   = 10000000
)