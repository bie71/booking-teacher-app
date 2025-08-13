package handler

import (
	"fmt"
	"net/http"
	"payment/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/spf13/cast"
)

type Handler struct {
	paymentService *service.Service
}

func NewHandler(paymentService *service.Service) *Handler {
	return &Handler{
		paymentService: paymentService,
	}
}

func (c *Handler) CreatePayment(ctx *gin.Context) {
	var req struct {
		BookingID     uint   `json:"booking_id"`
		Amount        int64  `json:"amount"`
		PaymentMethod string `json:"payment_method"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	orderID := fmt.Sprintf("BOOK-%d-%d", req.BookingID, time.Now().Unix())

	url, err := c.paymentService.CreatePayment(ctx, orderID, req.BookingID, req.Amount, req.PaymentMethod)
	if err != nil {
		if err.Error() == "booking not found" || err.Error() == "payment method not found or not active" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"snap_url": url})
}

func (c *Handler) HandleWebhook(ctx *gin.Context) {
	var notif snap.ResponseWithMap
	if err := ctx.ShouldBindJSON(&notif); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notif"})
		return
	}

	// Update status
	status := notif["transaction_status"].(string)
	orderID := notif["order_id"].(string)
	payment, err := c.paymentService.Handle(orderID, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to handle payment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"payment": payment})
}

func (h *Handler) GetAll(ctx *gin.Context) {
	methods, err := h.paymentService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": methods})
}

func (h *Handler) SetActive(ctx *gin.Context) {
	var req struct {
		Code   string `json:"code"`
		Active bool   `json:"active"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.paymentService.SetActive(ctx, req.Code, req.Active)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment method"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *Handler) GetPaymentById(ctx *gin.Context) {
	id := ctx.Param("id")
	payment, err := h.paymentService.GetPaymentById(cast.ToUint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": payment})
}

func (h *Handler) GetPayments(ctx *gin.Context) {
	payments, err := h.paymentService.GetPayments(cast.ToInt(ctx.Query("page")), cast.ToInt(ctx.Query("limit")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payments"})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}
