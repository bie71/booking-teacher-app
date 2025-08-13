package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentMethodHandler struct {
	paymentService interface {
		GetPaymentMethods(ctx context.Context, page int, limit int) (interface{}, error)
		UpdatePaymentMethodStatus(ctx context.Context, code string, active bool) error
	}
}

func NewPaymentMethodHandler(paymentService interface {
	GetPaymentMethods(ctx context.Context, page int, limit int) (interface{}, error)
	UpdatePaymentMethodStatus(ctx context.Context, code string, active bool) error
}) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentMethodHandler) GetPaymentMethods(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	methods, err := h.paymentService.GetPaymentMethods(context.Background(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, methods)
}

func (h *PaymentMethodHandler) UpdatePaymentMethodStatus(c *gin.Context) {
	var req struct {
		Code   string `json:"code" binding:"required"`
		Active bool   `json:"active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.paymentService.UpdatePaymentMethodStatus(context.Background(), req.Code, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method status updated successfully"})
}
