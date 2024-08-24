package service

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rasul07/alif-task/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of WalletStorage
type MockWalletStorage struct {
	mock.Mock
}

func (m *MockWalletStorage) CheckWalletExists(walletID, userID string) (bool, error) {
	args := m.Called(walletID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockWalletStorage) GetWallet(walletID, userID string) (*models.Wallet, error) {
	args := m.Called(walletID, userID)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletStorage) UpdateWalletBalance(walletID, userID string, newBalance, amount int64) error {
	args := m.Called(walletID, userID, newBalance, amount)
	return args.Error(0)
}

func (m *MockWalletStorage) GetTransactions(walletID string) (int, int64, error) {
	args := m.Called(walletID)
	return args.Int(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockWalletStorage) GetBalance(walletID, userID string) (int64, error) {
	args := m.Called(walletID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockWalletStorage) IsIdentified(userID string) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func TestCheckWalletExists(t *testing.T) {
	mockStorage := new(MockWalletStorage)
	service := &walletService{storage: mockStorage, logger: log.Default()}

	walletID1 := uuid.New().String()
	userID1 := uuid.New().String()
	walletID2 := uuid.New().String()
	userID2 := uuid.New().String()

	t.Run("Wallet exists", func(t *testing.T) {
		mockStorage.On("CheckWalletExists", walletID1, userID1).Return(true, nil).Once()

		exists, err := service.CheckWalletExists(walletID1, userID1)

		assert.True(t, exists)
		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Wallet does not exist", func(t *testing.T) {
		mockStorage.On("CheckWalletExists", walletID2, userID2).Return(false, nil).Once()

		exists, err := service.CheckWalletExists(walletID2, userID2)

		assert.False(t, exists)
		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestTopUpWallet(t *testing.T) {
	mockStorage := new(MockWalletStorage)
	service := &walletService{storage: mockStorage, logger: log.Default()}

	walletID1 := uuid.New().String()
	userID1 := uuid.New().String()
	walletID2 := uuid.New().String()
	userID2 := uuid.New().String()

	t.Run("Successful top-up", func(t *testing.T) {
		wallet := &models.Wallet{ID: walletID1, UserID: userID1, Balance: 5000}
		mockStorage.On("GetWallet", walletID1, userID1).Return(wallet, nil).Once()
		mockStorage.On("IsIdentified", userID1).Return(true, nil).Once()
		mockStorage.On("UpdateWalletBalance", walletID1, userID1, int64(15000), int64(10000)).Return(nil).Once()

		err := service.TopUpWallet(walletID1, userID1, "100.00")

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Top-up exceeds maximum balance", func(t *testing.T) {
		wallet := &models.Wallet{ID: walletID2, UserID: userID2, Balance: 9000000}
		mockStorage.On("GetWallet", walletID2, userID2).Return(wallet, nil).Once()
		mockStorage.On("IsIdentified", userID2).Return(false, nil).Once()

		err := service.TopUpWallet(walletID2, userID2, "20000.00")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "top-up would exceed maximum balance")
		mockStorage.AssertExpectations(t)
	})
}

func TestGetTransactions(t *testing.T) {
	mockStorage := new(MockWalletStorage)
	service := &walletService{storage: mockStorage, logger: log.Default()}

	walletID1 := uuid.New().String()
	userID1 := uuid.New().String()
	walletID2 := uuid.New().String()
	userID2 := uuid.New().String()

	t.Run("Successful get transactions", func(t *testing.T) {
		wallet := &models.Wallet{ID: walletID1, UserID: userID1, Balance: 10000}
		mockStorage.On("GetWallet", walletID1, userID1).Return(wallet, nil).Once()
		mockStorage.On("GetTransactions", walletID1).Return(5, int64(50000), nil).Once()

		count, total, err := service.GetTransactions(walletID1, userID1)

		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		assert.Equal(t, "500.00", total)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Wallet not found", func(t *testing.T) {
		mockStorage.On("GetWallet", walletID2, userID2).Return((*models.Wallet)(nil), errors.New("wallet not found")).Once()

		_, _, err := service.GetTransactions(walletID2, userID2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "couldn't get this wallet")
		mockStorage.AssertExpectations(t)
	})
}

func TestGetBalance(t *testing.T) {
	mockStorage := new(MockWalletStorage)
	service := &walletService{storage: mockStorage, logger: log.Default()}

	walletID1 := uuid.New().String()
	userID1 := uuid.New().String()
	walletID2 := uuid.New().String()
	userID2 := uuid.New().String()

	t.Run("Successful get balance", func(t *testing.T) {
		mockStorage.On("GetBalance", walletID1, userID1).Return(int64(10000), nil).Once()

		balance, err := service.GetBalance(walletID1, userID1)

		assert.NoError(t, err)
		assert.Equal(t, "100.00", balance)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Error getting balance", func(t *testing.T) {
		mockStorage.On("GetBalance", walletID2, userID2).Return(int64(0), errors.New("database error")).Once()

		_, err := service.GetBalance(walletID2, userID2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockStorage.AssertExpectations(t)
	})
}