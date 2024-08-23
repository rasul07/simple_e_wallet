package service

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/rasul07/alif-task/internal/models"
	"github.com/rasul07/alif-task/internal/storage"
)

type WalletService interface {
	CheckWalletExists(userID int) (bool, error)
	TopUpWallet(userID int, amount string) error
	GetTransactions(userID int) (int, string, error)
	GetBalance(userID int) (string, error)
}

type walletService struct {
	storage *storage.WalletStorage
}

func NewWalletService(db *sql.DB) WalletService {
	return &walletService{storage: storage.NewWalletStorage(db)}
}

func (s *walletService) CheckWalletExists(walletID int) (bool, error) {
	return s.storage.CheckWalletExists(walletID)
}

func (s *walletService) TopUpWallet(walletID int, amount string) error {
	wallet, err := s.storage.GetWallet(walletID)
	if err != nil {
		return err
	}

	isIdentified, err := s.storage.IsIdentified(wallet.UserID)
	if err != nil {
		return err
	}

	currentBalance, err := strconv.Atoi(wallet.Balance)
	if err != nil {
		return err
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return err
	}
	
	newBalance := int64(currentBalance) + int64(amountInt)
	maxBalance := models.MaxBalanceUnidentified
	if isIdentified {
		maxBalance = models.MaxBalanceIdentified
	}

	if newBalance > int64(maxBalance) {
		return fmt.Errorf("top-up would exceed maximum balance")
	}

	err = s.storage.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return err
	}

	err = s.storage.AddTransaction(wallet.ID, int64(amountInt))
	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) GetTransactions(walletID int) (int, string, error) {
	trCount, total, err := s.storage.GetTransactions(walletID)
	if err != nil {
		return 0, "", err
	}

	totalStr := strconv.Itoa(int(total))

	return trCount, totalStr, err
}

func (s *walletService) GetBalance(walletID int) (string, error) {
	balance, err := s.storage.GetBalance(walletID)
	if err != nil {
		return "", err
	}

	balanceStr := strconv.Itoa(int(balance))

	return balanceStr, err
}