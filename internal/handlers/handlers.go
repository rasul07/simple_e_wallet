package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rasul07/alif-task/internal/models"
	"github.com/rasul07/alif-task/internal/service"
)

type Handler struct {
	walletService service.WalletService
}

func NewHandler(walletService service.WalletService) *Handler {
	return &Handler{walletService: walletService}
}

// CheckWalletExists godoc
// @Summary Check if a wallet exists
// @Description Check if a wallet exists for the given wallet ID
// @Tags wallet
// @Accept json
// @Produce json
// @Param X-UserId header int true "User ID"
// @Param X-Digest header string true "Digest"
// @Param wallet_id body int true "Wallet ID"
// @Success 200 {object} map[string]bool
// @Router /v1/wallet/check [post]
func (h *Handler) CheckWalletExists(c *gin.Context) {
	var walletID int

	if err := c.ShouldBindJSON(&walletID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	exists, err := h.walletService.CheckWalletExists(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

// TopUpWallet godoc
// @Summary Top up a wallet
// @Description Top up a wallet with the given amount
// @Tags wallet
// @Accept json
// @Produce json
// @Param X-UserId header int true "User ID"
// @Param X-Digest header string true "Digest"
// @Param request body models.TopUpRequest true "Top up request"
// @Success 200 {object} map[string]string
// @Router /v1/wallet/topup [post]
func (h *Handler) TopUpWallet(c *gin.Context) {
	var request models.TopUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.walletService.TopUpWallet(request.WalletID, request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet topped up successfully"})
}

// GetTransactions godoc
// @Summary Get transactions for the current month
// @Description Get the total number and amount of transactions for the current month
// @Tags wallet
// @Accept json
// @Produce json
// @Param X-UserId header int true "User ID"
// @Param X-Digest header string true "Digest"
// @Param wallet_id body int true "Wallet ID"
// @Success 200 {object} map[string]interface{}
// @Router /v1/wallet/transactions [post]
func (h *Handler) GetTransactions(c *gin.Context) {
	var walletID int

	if err := c.ShouldBindJSON(&walletID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	count, total, err := h.walletService.GetTransactions(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"total": total,
	})
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get the current balance of a wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param X-UserId header int true "User ID"
// @Param X-Digest header string true "Digest"
// @Param wallet_id body int true "Wallet ID"
// @Success 200 {object} map[string]float64
// @Router /v1/wallet/balance [post]
func (h *Handler) GetBalance(c *gin.Context) {
	var walletID int

	if err := c.ShouldBindJSON(&walletID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	balance, err := h.walletService.GetBalance(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
