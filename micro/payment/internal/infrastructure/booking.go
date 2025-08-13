package infrastructure

import (
	"fmt"
	"log"
	"payment/internal/config"
	"payment/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

type Booking struct {
	restyClient    *resty.Client
	serviceBooking config.Service
}

// InternalBookingResponse represents the structure returned by the
// booking service's internal endpoint for a single booking.  It contains
// only the fields needed by the payment service to determine the owner
// of a booking.
type InternalBookingResponse struct {
	Booking struct {
		ID         uint    `json:"id"`
		UserID     uint    `json:"user_id"`
		ScheduleID uint    `json:"schedule_id"`
		Status     string  `json:"status"`
		PaymentID  *uint   `json:"payment_id"`
		TotalPrice float64 `json:"total_price"`
	} `json:"booking"`
}

// GetBooking fetches a booking record from the booking service's internal
// endpoint.  It retrieves the booking's user ID and other fields needed
// for payment handling.  If the booking does not exist, it returns an error.
func (b *Booking) GetBooking(id uint) (*InternalBookingResponse, error) {
	url := fmt.Sprintf("%s:%s/api/v1/internal/bookings/%d", b.serviceBooking.Host, b.serviceBooking.Port, id)
	log.Println(url, "fetch internal booking")
	var result InternalBookingResponse
	resp, err := b.restyClient.R().
		SetResult(&result).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("booking not found")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch booking: %s", resp.Status())
	}
	return &result, nil
}

func NewBooking(restyClient *resty.Client, serviceBooking config.Service) *Booking {
	return &Booking{
		restyClient: restyClient,
		serviceBooking: config.Service{
			Host: serviceBooking.Host,
			Port: serviceBooking.Port,
		},
	}
}

func (b *Booking) UpdateBookingStatus(bookingID, paymentID string, status string) error {

	url := fmt.Sprintf("%s:%s/private/bookings/%s/status", b.serviceBooking.Host, b.serviceBooking.Port, bookingID)
	log.Println(url, "update booking status")

	resp, err := b.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("status", status).
		SetQueryParam("paymentId", paymentID).
		Put(url)

	if resp.IsError() {
		return fmt.Errorf("failed to update booking status: %s", resp.Status())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to update booking status: %s", resp.Status())
	}

	if err != nil {
		return err
	}

	log.Println("update booking status success")

	return nil
}

func (b *Booking) CheckBookingExist(c *gin.Context, bookingID string) (bool, error) {
	url := fmt.Sprintf("%s:%s/api/v1/booking/%s", b.serviceBooking.Host, b.serviceBooking.Port, bookingID)

	response := model.BookingResponse{}
	token, _ := c.Get("token")
	resp, err := b.restyClient.R().
		SetAuthToken(cast.ToString(token)).
		SetResult(&response).
		Get(url)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != 200 {
		return false, nil
	}
	if response.ID == 0 {
		return true, nil
	}
	return false, nil
}
