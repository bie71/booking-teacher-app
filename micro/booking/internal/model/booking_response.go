package model

import (
	"booking/internal/infrastructure/payment"
	"booking/internal/infrastructure/user"
	"time"
)

type TeacherResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Bio          string  `json:"bio"`
	Price        float64 `json:"price_per_hour"`
	ProfileImage string  `json:"profile_image"`
}

type ScheduleResponse struct {
	ID         uint             `json:"id"`
	Date       string           `json:"date"`
	StartTime  string           `json:"start_time"`
	EndTime    string           `json:"end_time"`
	Status     string           `json:"status"`
	TotalPrice float64          `json:"total_price"`
	Teacher    *TeacherResponse `json:"teacher,omitempty"`
}

type BookingResponse struct {
	ID         uint                     `json:"id"`
	Status     string                   `json:"status"`
	Note       string                   `json:"note"`
	CreatedAt  time.Time                `json:"created_at"`
	TotalPrice float64                  `json:"total_price"`
	Schedule   *ScheduleResponse        `json:"schedule,omitempty"`
	User       *user.UserResponse       `json:"user,omitempty"`
	Payment    *payment.PaymentResponse `json:"payment,omitempty"`
}

type PaginatedBookingsResponse struct {
	Data       []BookingResponse `json:"bookings"`
	Pagination PaginationMeta    `json:"pagination"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
