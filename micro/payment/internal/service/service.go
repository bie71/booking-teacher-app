package service

import (
	"context"
	"errors"
	"log"
	"payment/internal/infrastructure"
	"payment/internal/model"
	"payment/internal/pkg"
	"payment/internal/repository"
	"strings"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Service struct {
	infra          *infrastructure.PaymentService
	serviceBooking *infrastructure.Booking
	repository     *repository.Repository
	serviceUser    *infrastructure.User
}

func NewService(
	infra *infrastructure.PaymentService,
	repository *repository.Repository,
	serviceBooking *infrastructure.Booking,
	serviceUser *infrastructure.User,
) *Service {
	return &Service{
		infra:          infra,
		repository:     repository,
		serviceBooking: serviceBooking,
		serviceUser:    serviceUser,
	}
}

func (s *Service) CreatePayment(c *gin.Context, orderID string, bookingID uint, amount int64, paymentMethod string) (string, error) {

	exist, err := s.serviceBooking.CheckBookingExist(c, cast.ToString(bookingID))
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("booking not found")
	}

	method, err := s.repository.GetPaymentMethodByName(context.Background(), strings.ToLower(paymentMethod))

	if err != nil || !method.IsActive {
		log.Println(err, method)
		return "", errors.New("payment method not found or not active")
	}

	redirectUrl, err := s.infra.CreatePayment(orderID, bookingID, amount)
	if err != nil {
		log.Println(err)
		return "", err
	}

	payment := &model.Payment{
		MidtransTransactionID: orderID,
		Amount:                float64(amount),
		Status:                "pending",
		PaymentMethod:         method.Name,
		BookingID:             bookingID,
	}

	err = s.repository.CreatePayment(payment)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// After successfully creating the payment record, record a recent activity
	// for the user associated with this booking.  We fetch the booking to
	// determine the user ID from the booking service's internal endpoint.  If
	// the booking cannot be retrieved or the user service call fails, we log
	// the error but do not block the payment creation flow.
	if s.serviceUser != nil && s.serviceBooking != nil {
		// Retrieve booking details to get the user ID
		if bookingDetail, err := s.serviceBooking.GetBooking(bookingID); err == nil {
			userID := bookingDetail.Booking.UserID
			// Compose a description in Indonesian indicating that the user created a payment
			go func() {
				description := fmt.Sprintf("Pengguna membuat pembayaran untuk pemesanan #%d", bookingID)
				if err := s.serviceUser.CreateActivityLog(userID, "create_payment", description); err != nil {
					log.Printf("failed to log activity for payment creation: %v", err)
				}
			}()
		} else {
			log.Printf("failed to fetch booking for activity log: %v", err)
		}
	}

	return redirectUrl, nil
}

func (s *Service) Handle(trxId string, status string) (*model.Payment, error) {

	payment, err := s.repository.GetPaymentByTrxID(trxId)
	if err != nil {
		return nil, err
	}

	payment.Status = status

	switch status {
	case "settlement":
		now := time.Now()
		payment.PaidAt = &now
	case "expire":
		payment.Status = "failed"
	}

	err = s.repository.UpdatePayment(payment)
	if err != nil {
		return nil, err
	}

	if status == "settlement" {
		err = s.serviceBooking.UpdateBookingStatus(cast.ToString(payment.BookingID), cast.ToString(payment.ID), "paid")
		if err != nil {
			return nil, err
		}
		// Record an activity log when payment is settled (paid) indicating the user
		// has completed the payment and the booking is paid.  Fetch the user ID
		// from the booking service and create the activity.  Errors from
		// logging are recorded but do not interrupt the callback flow.
		if s.serviceUser != nil && s.serviceBooking != nil {
			if bookingDetail, err2 := s.serviceBooking.GetBooking(payment.BookingID); err2 == nil {
				userID := bookingDetail.Booking.UserID
				go func() {
					description := fmt.Sprintf("Pembayaran untuk pemesanan #%d berhasil dan booking telah dibayar", payment.BookingID)
					if err3 := s.serviceUser.CreateActivityLog(userID, "payment_success", description); err3 != nil {
						log.Printf("failed to log activity for payment success: %v", err3)
					}
				}()
			} else {
				log.Printf("failed to fetch booking for payment success activity: %v", err2)
			}
		}
	}
	if status == "expire" || status == "cancel" {
		err = s.serviceBooking.UpdateBookingStatus(cast.ToString(payment.BookingID), cast.ToString(payment.ID), "cancelled")
		if err != nil {
			return nil, err
		}
		// Optionally, you could log cancellation or expiry activities here if desired
	}

	return payment, nil
}

func (s *Service) GetAll(ctx context.Context) ([]model.PaymentMethod, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) SetActive(ctx context.Context, code string, active bool) error {
	return s.repository.UpdateStatus(ctx, code, active)
}

// CRUD operations for payment methods
func (s *Service) CreatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error {
	return s.repository.CreatePaymentMethod(ctx, method)
}

func (s *Service) GetPaymentMethodByID(ctx context.Context, id uint) (*model.PaymentMethod, error) {
	return s.repository.GetPaymentMethodByID(ctx, id)
}

func (s *Service) GetPaymentMethods(ctx context.Context, page int, limit int) ([]model.PaymentMethod, int64, error) {
	return s.repository.GetPaymentMethods(ctx, page, limit)
}

func (s *Service) UpdatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error {
	return s.repository.UpdatePaymentMethod(ctx, method)
}

func (s *Service) DeletePaymentMethod(ctx context.Context, id uint) error {
	return s.repository.DeletePaymentMethod(ctx, id)
}

func (s *Service) UpdatePaymentMethodStatus(ctx context.Context, code string, active bool) error {
	return s.repository.UpdateStatus(ctx, code, active)
}

func (s *Service) GetPaymentById(id uint) (*model.PaymentResponse, error) {
	payment, err := s.repository.GetPaymentById(id)
	if err != nil {
		return nil, err
	}
	return &model.PaymentResponse{
		Id:                    payment.ID,
		MidtransTransactionID: payment.MidtransTransactionID,
		Amount:                payment.Amount,
		Status:                payment.Status,
		PaymentMethod:         payment.PaymentMethod,
		BookingID:             payment.BookingID,
		PaidAt:                payment.PaidAt,
		CreatedAt:             payment.CreatedAt,
		UpdatedAt:             payment.UpdatedAt,
	}, nil
}

func (s *Service) GetPayments(pg int, limit int) (pkg.ResponsePaginate, error) {
	return s.repository.GetPayments(pg, limit)
}
