package http

import (
	"errors"
	"net/http"
	"strconv"

	"transaction-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.TransactionUsecase
}

func NewHandler(uc *usecase.TransactionUsecase) Handler {
	return Handler{uc: uc}
}

func (h Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", h.healthz)
	r.POST("/transactions/payments", h.processPayment)
	r.GET("/transactions", h.list)
	r.GET("/transactions/:id", h.detail)
}

type processPaymentRequest struct {
	ExternalID   string `json:"external_id"`
	Type         string `json:"type" binding:"required"`
	FromWalletID string `json:"from_wallet_id" binding:"required"`
	ToWalletID   string `json:"to_wallet_id" binding:"required"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
}

// processPayment godoc
// @Summary Process payment transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "Idempotency key"
// @Param payload body processPaymentRequest true "Payment payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/payments [post]
func (h Handler) processPayment(c *gin.Context) {
	var req processPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	idemKey := c.GetHeader("Idempotency-Key")
	out, replayed, err := h.uc.ProcessPayment(c.Request.Context(), usecase.ProcessPaymentInput{
		IdempotencyKey: idemKey,
		ExternalID:     req.ExternalID,
		Type:           req.Type,
		FromWalletID:   req.FromWalletID,
		ToWalletID:     req.ToWalletID,
		Amount:         req.Amount,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrIdempotencyRequired), errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrIdempotencyMismatch):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrIdempotencyBusy):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"replayed":    replayed,
		"transaction": out,
	})
}

// list godoc
// @Summary List transaction history
// @Tags transactions
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h Handler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	rows, err := h.uc.ListHistory(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": rows, "limit": limit, "offset": offset})
}

// detail godoc
// @Summary Get transaction detail
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/{id} [get]
func (h Handler) detail(c *gin.Context) {
	tx, err := h.uc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, usecase.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

// healthz godoc
// @Summary Health check
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h Handler) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"service": "transaction-service", "status": "ok"})
}

