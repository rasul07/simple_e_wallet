package storage

import (
	"database/sql"
	"time"

	"github.com/rasul07/alif-task/internal/models"
)

type WalletStorage struct {
	db *sql.DB
}

func NewWalletStorage(db *sql.DB) *WalletStorage {
	return &WalletStorage{db: db}
}

func (s *WalletStorage) CheckWalletExists(userID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM wallets WHERE user_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (s *WalletStorage) GetWallet(userID string) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := s.db.QueryRow("SELECT id, user_id, balance FROM wallets WHERE user_id=$1", userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletStorage) UpdateWalletBalance(walletID int64, newBalance int64) error {
	_, err := s.db.Exec("UPDATE wallets SET balance=$1 WHERE id=$2", newBalance, walletID)
	return err
}

func (s *WalletStorage) AddTransaction(walletID int64, amount int64) error {
	_, err := s.db.Exec("INSERT INTO transactions (wallet_id, amount, timestamp) VALUES ($1, $2, $3)", walletID, amount, time.Now())
	return err
}

func (s *WalletStorage) GetTransactions(userID string) (int, int64, error) {
	var count int
	var total int64

	err := s.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN wallets w ON t.wallet_id = w.id
		WHERE w.user_id=$1 AND t.timestamp >= DATE_TRUNC('month', CURRENT_DATE)
	`, userID).Scan(&count, &total)

	return count, total, err
}

func (s *WalletStorage) GetBalance(userID string) (int64, error) {
	var balance int64
	err := s.db.QueryRow("SELECT balance FROM wallets WHERE user_id=$1", userID).Scan(&balance)
	return balance, err
}

func (s *WalletStorage) IsIdentified(userID string) (bool, error) {
	var identified bool
	err := s.db.QueryRow("SELECT is_identified FROM users WHERE id=$1", userID).Scan(&identified)
	return identified, err
}