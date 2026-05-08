package http

import (
	"errors"
	"net/http"

	"wallet-service/internal/domain"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	uc *usecase.WalletUsecase
}

func NewHandler(uc *usecase.WalletUsecase) Handler {
	return Handler{uc: uc}
}

func (h Handler) RegisterRoutes(r *gin.Engine, secret string) {
	r.GET("/healthz", h.healthz)

	protected := r.Group("/")
	protected.Use(AuthMiddleware(secret))
	protected.POST("/wallets", h.createWallet)
	protected.POST("/wallet/topups", h.topup)
	protected.POST("/wallet/withdrawals", h.withdraw)
	protected.POST("/wallet/transfers", h.transfer)
	protected.GET("/wallets/:userID/balance", h.balance)
}

type balanceResponse struct {
	Balance domain.Wallet `json:"balance"`
}


type createWalletRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Currency string `json:"currency"`
}

// createWallet godoc
// @Summary Create wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param payload body createWalletRequest true "Create wallet payload"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallets [post]
func (h Handler) createWallet(c *gin.Context) {
	var req createWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	w, err := h.uc.CreateWallet(c.Request.Context(), req.UserID, req.Currency)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrWalletExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, w)
}

type amountRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

// topup godoc
// @Summary Topup wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param payload body amountRequest true "Topup payload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/topups [post]
func (h Handler) topup(c *gin.Context) {
	var req amountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	w, err := h.uc.Topup(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

// withdraw godoc
// @Summary Withdraw from wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param payload body amountRequest true "Withdraw payload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/withdrawals [post]
func (h Handler) withdraw(c *gin.Context) {
	var req amountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	w, err := h.uc.Withdraw(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

type transferRequest struct {
	FromUserID string `json:"from_user_id" binding:"required"`
	ToUserID   string `json:"to_user_id" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
}

// transfer godoc
// @Summary Transfer between wallets
// @Tags wallet
// @Accept json
// @Produce json
// @Param payload body transferRequest true "Transfer payload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/transfers [post]
func (h Handler) transfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	fromWallet, toWallet, err := h.uc.Transfer(c.Request.Context(), req.FromUserID, req.ToUserID, req.Amount)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"from_wallet": fromWallet,
		"to_wallet":   toWallet,
	})
}

// balance godoc
// @Summary Get wallet balance
// @Tags wallet
// @Produce json
// @Param userID path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} balanceResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallets/{userID}/balance [get]
func (h Handler) balance(c *gin.Context) {
	userID := c.Param("userID")
	balance, err := h.uc.BalanceInquiry(c.Request.Context(), userID)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, balanceResponse{
		Balance: balance,
	})
}


func (h Handler) handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrWalletNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrInsufficientBalance):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrConcurrentUpdate):
		c.JSON(http.StatusConflict, gin.H{"error": "concurrent update, retry request"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// healthz godoc
// @Summary Health check
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h Handler) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"service": "wallet-service", "status": "ok"})
}

