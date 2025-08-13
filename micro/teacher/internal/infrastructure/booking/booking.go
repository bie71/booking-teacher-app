package booking

import (
	"fmt"
	"teacher/internal/config"
	"teacher/internal/models"

	"github.com/go-resty/resty/v2"
)

type BookingService struct {
	Client *resty.Client
	Cfg    *config.Service
}

func NewBookingService(initRestyClient *resty.Client, cfg *config.Service) *BookingService {
	return &BookingService{
		Client: initRestyClient,
		Cfg:    cfg,
	}
}

func (s *BookingService) GetBookingByTeacherId(teacherId uint) ([]models.TeacherBookingResponse, error) {
    // Use the internal endpoint on the booking service that does not require
    // authentication.  This allows the teacher service to retrieve teacher
    // bookings without forwarding the user's JWT.  See booking service
    // cmd/app/main.go for route definition.
    url := fmt.Sprintf("%s/api/v1/internal/bookings/teacher/%d", s.Cfg.Host, teacherId)

	var bookings models.Response

	resp, err := s.Client.R().SetResult(&bookings).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("booking service returned status: %d", resp.StatusCode())
	}

	if len(bookings.Bookings) == 0 {
		return nil, nil
	}

	return bookings.Bookings, nil

}
