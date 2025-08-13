package payment

import (
	"booking/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

type DataPaymentResponse struct {
	Data PaymentResponse `json:"data"`
}

type PaymentResponse struct {
	Id                    uint       `json:"id"`
	MidtransTransactionID string     `json:"midtrans_transaction_id"`
	Amount                float64    `json:"amount"`
	Status                string     `json:"status"`
	PaymentMethod         string     `json:"payment_method"`
	BookingID             uint       `json:"booking_id"`
	PaidAt                *time.Time `json:"paid_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type Payment struct {
	restyClient *resty.Client
	service     config.Service
}

func NewPaymentHttp(service config.Service, restyClient *resty.Client) *Payment {
	return &Payment{
		restyClient: restyClient,
		service:     service,
	}
}

func (p *Payment) GetPaymentByID(paymentID uint) (map[uint]*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payment/%d", p.service.Host, paymentID)

	var respPayment PaymentResponse

	resp, err := p.restyClient.R().
		SetResult(&respPayment).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("payment not found")
	}

	result := make(map[uint]*PaymentResponse)
	result[paymentID] = &respPayment

	return result, nil

}

func (p *Payment) GetPaymentDetail(c *gin.Context, paymentID uint) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payment/%d", p.service.Host, paymentID)

	var respPayment DataPaymentResponse

	token, _ := c.Get("token")
	resp, err := p.restyClient.R().
		SetAuthToken(cast.ToString(token)).
		SetResult(&respPayment).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("payment not found")
	}

	return &respPayment.Data, nil

}
