package service

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rasul07/alif-task/internal/models"
	"github.com/rasul07/alif-task/internal/storage"
)

type WalletService interface {
	CheckWalletExists(walletID, userID string) (bool, error)
	TopUpWallet(walletID, userID, amount string) error
	GetTransactions(walletID, userID string) (int, string, error)
	GetBalance(walletID, userID string) (string, error)
}

type walletService struct {
	storage *storage.WalletStorage
}

func NewWalletService(db *sql.DB) WalletService {
	return &walletService{storage: storage.NewWalletStorage(db)}
}

func (s *walletService) CheckWalletExists(walletID, userID string) (bool, error) {
	return s.storage.CheckWalletExists(walletID, userID)
}

func (s *walletService) TopUpWallet(walletID, userID, amount string) error {
	wallet, err := s.storage.GetWallet(walletID, userID)
	if err != nil {
		return errors.Wrap(err, "Error getting wallet")
	}

	// Check if user is identified
	isIdentified, err := s.storage.IsIdentified(wallet.UserID)
	if err != nil {
		return err
	}

	// Convert amount being added to numeric type
	amountInt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	newBalance := wallet.Balance + int64(amountInt)
	maxBalance := models.MaxBalanceUnidentified // if user is not identified max balance is 10.000
	if isIdentified {
		// if user is identified max balance is 100.000
		maxBalance = models.MaxBalanceIdentified
	}

	// Check possible balance overflow
	if newBalance > int64(maxBalance) {
		return fmt.Errorf("top-up would exceed maximum balance")
	}

	err = s.storage.UpdateWalletBalance(wallet.ID, userID, newBalance, int64(amountInt))
	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) GetTransactions(walletID, userID string) (int, string, error) {
	_, err := s.storage.GetWallet(walletID, userID)
	if err != nil {
		return 0, "", errors.Wrap(err, "couldn't get this wallet")
	}

	trCount, total, err := s.storage.GetTransactions(walletID)
	if err != nil {
		return 0, "", err
	}

	totalStr := fmt.Sprintf("%.2f", float64(total/100))

	return trCount, totalStr, err
}

func (s *walletService) GetBalance(walletID, userID string) (string, error) {
	balance, err := s.storage.GetBalance(walletID, userID)
	if err != nil {
		return "", err
	}

	balanceStr := fmt.Sprintf("%.2f", float64(balance/100))

	return balanceStr, err
}
