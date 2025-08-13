package handler

import (
	"context"
	"net/http"
	"payment/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentMethodCRUDHandler struct {
	paymentService interface {
		CreatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error
		GetPaymentMethodByID(ctx context.Context, id uint) (*model.PaymentMethod, error)
		GetPaymentMethods(ctx context.Context, page int, limit int) ([]model.PaymentMethod, int64, error)
		UpdatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error
		DeletePaymentMethod(ctx context.Context, id uint) error
		UpdatePaymentMethodStatus(ctx context.Context, code string, active bool) error
	}
}

func NewPaymentMethodCRUDHandler(paymentService interface {
	CreatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error
	GetPaymentMethodByID(ctx context.Context, id uint) (*model.PaymentMethod, error)
	GetPaymentMethods(ctx context.Context, page int, limit int) ([]model.PaymentMethod, int64, error)
	UpdatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id uint) error
	UpdatePaymentMethodStatus(ctx context.Context, code string, active bool) error
}) *PaymentMethodCRUDHandler {
	return &PaymentMethodCRUDHandler{
		paymentService: paymentService,
	}
}

// CreatePaymentMethod creates a new payment method
func (h *PaymentMethodCRUDHandler) CreatePaymentMethod(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Name     string `json:"name" binding:"required"`
		IsActive bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	method := &model.PaymentMethod{
		Code:     req.Code,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := h.paymentService.CreatePaymentMethod(context.Background(), method); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment method created successfully", "data": method})
}

// GetPaymentMethods returns paginated payment methods
func (h *PaymentMethodCRUDHandler) GetPaymentMethods(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	methods, total, err := h.paymentService.GetPaymentMethods(context.Background(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": methods,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
			"total_items":  total,
			"limit":        limit,
		},
	})
}

// GetPaymentMethod returns a single payment method by ID
func (h *PaymentMethodCRUDHandler) GetPaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	method, err := h.paymentService.GetPaymentMethodByID(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": method})
}

// UpdatePaymentMethod updates an existing payment method
func (h *PaymentMethodCRUDHandler) UpdatePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		IsActive *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	method, err := h.paymentService.GetPaymentMethodByID(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	method.Name = req.Name
	if req.IsActive != nil {
		method.IsActive = *req.IsActive
	}

	if err := h.paymentService.UpdatePaymentMethod(context.Background(), method); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method updated successfully", "data": method})
}

// DeletePaymentMethod deletes a payment method
func (h *PaymentMethodCRUDHandler) DeletePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.paymentService.DeletePaymentMethod(context.Background(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}
