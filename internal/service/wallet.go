package service

import (
	"database/sql"
	"fmt"
	"log"
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
	storage storage.WalletStorager
	logger  *log.Logger
}

func NewWalletService(db *sql.DB) WalletService {
	return &walletService{
		storage: storage.NewWalletStorage(db),
		logger:  log.New(log.Writer(), "WalletService: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (s *walletService) CheckWalletExists(walletID, userID string) (bool, error) {
	s.logger.Printf("Checking if wallet exists: walletID=%s, userID=%s", walletID, userID)
	exists, err := s.storage.CheckWalletExists(walletID, userID)
	if err != nil {
		s.logger.Printf("Error checking wallet existence: %v", err)
		return false, err
	}

	return exists, nil
}

func (s *walletService) TopUpWallet(walletID, userID, amount string) error {
	s.logger.Printf("Topping up wallet: walletID=%s, userID=%s, amount=%s", walletID, userID, amount)
	wallet, err := s.storage.GetWallet(walletID, userID)
	if err != nil {
		s.logger.Printf("Error getting wallet: %v", err)
		return errors.Wrap(err, "Error getting wallet")
	}

	// Check if user is identified
	isIdentified, err := s.storage.IsIdentified(wallet.UserID)
	if err != nil {
		s.logger.Printf("Error checking if user is identified: %v", err)
		return err
	}

	// Convert amount being added to numeric type
	amountInt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		s.logger.Printf("Error parsing amount: %v", err)
		return err
	}
	newAmount := int64(amountInt)*100

	newBalance := wallet.Balance + newAmount
	maxBalance := models.MaxBalanceUnidentified // if user is not identified max balance is 10.000
	if isIdentified {
		// if user is identified max balance is 100.000
		maxBalance = models.MaxBalanceIdentified
	}

	// Check possible balance overflow
	if newBalance > int64(maxBalance) {
		s.logger.Printf("Top-up would exceed maximum balance: current=%d, new=%d, max=%d", wallet.Balance, newBalance, maxBalance)
		return fmt.Errorf("top-up would exceed maximum balance")
	}

	err = s.storage.UpdateWalletBalance(wallet.ID, userID, newBalance, newAmount)
	if err != nil {
		s.logger.Printf("Error updating wallet balance: %v", err)
		return err
	}

	return nil
}

func (s *walletService) GetTransactions(walletID, userID string) (int, string, error) {
	s.logger.Printf("Getting transactions: walletID=%s, userID=%s", walletID, userID)
	_, err := s.storage.GetWallet(walletID, userID)
	if err != nil {
		s.logger.Printf("Error getting wallet: %v", err)
		return 0, "", errors.Wrap(err, "couldn't get this wallet")
	}

	trCount, total, err := s.storage.GetTransactions(walletID)
	if err != nil {
		s.logger.Printf("Error getting transactions: %v", err)
		return 0, "", err
	}

	totalStr := fmt.Sprintf("%.2f", float64(total/100))

	return trCount, totalStr, err
}

func (s *walletService) GetBalance(walletID, userID string) (string, error) {
	s.logger.Printf("Getting balance: walletID=%s, userID=%s", walletID, userID)
	balance, err := s.storage.GetBalance(walletID, userID)
	if err != nil {
		s.logger.Printf("Error getting balance: %v", err)
		return "", err
	}

	balanceStr := fmt.Sprintf("%.2f", float64(balance/100))
	s.logger.Printf("Balance retrieved: %s", balanceStr)

	return balanceStr, err
}
