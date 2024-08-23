package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/rasul07/alif-task/internal/models"
)

type WalletStorage struct {
	db *sql.DB
}

func NewWalletStorage(db *sql.DB) *WalletStorage {
	return &WalletStorage{db: db}
}

func (s *WalletStorage) CheckWalletExists(walletID, userID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM wallets WHERE id=$1 and user_id=$2)", walletID, userID).Scan(&exists)
	return exists, err
}

func (s *WalletStorage) GetWallet(walletID, userID string) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := s.db.QueryRow("SELECT id, user_id, balance FROM wallets WHERE id=$1 and user_id=$2", walletID, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletStorage) UpdateWalletBalance(walletID, userID string, newBalance, amount int64) error {
	tx, err := s.db.BeginTx(context.TODO(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return errors.Wrap(err, "unable to begin transaction to save card relation")
	}

	_, err = tx.Exec("UPDATE wallets SET balance=$1 WHERE id=$2 and user_id=$3", newBalance, walletID, userID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(rollbackErr, "unable to rollback transaction")
		}

		return errors.Wrap(err, "unable to execute transaction")
	}

	_, err = tx.Exec("INSERT INTO transactions (wallet_id, amount, created_at) VALUES ($1, $2, $3)", walletID, amount, time.Now())
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(rollbackErr, "unable to rollback transaction")
		}

		return errors.Wrap(err, "unable to execute transaction")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "unable to commit transaction")
	}

	return err
}

func (s *WalletStorage) GetTransactions(walletID string) (int, int64, error) {
	var count int
	var total int64

	err := s.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN wallets w ON t.wallet_id = w.id
		WHERE w.id=$1 AND t.created_at >= DATE_TRUNC('month', CURRENT_DATE)
	`, walletID).Scan(&count, &total)

	return count, total, err
}

func (s *WalletStorage) GetBalance(walletID, userID string) (int64, error) {
	var balance int64
	err := s.db.QueryRow("SELECT balance FROM wallets WHERE id=$1 and user_id=$2", walletID, userID).Scan(&balance)
	return balance, err
}

func (s *WalletStorage) IsIdentified(userID string) (bool, error) {
	var identified bool
	err := s.db.QueryRow("SELECT is_identified FROM users WHERE id=$1", userID).Scan(&identified)
	return identified, err
}
