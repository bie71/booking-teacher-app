package model

import "time"

type Booking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	ScheduleID     uint      `gorm:"index" json:"schedule_id"`
	Status         string    `gorm:"type:enum('pending','paid','cancelled','rescheduled');default:'pending'" json:"status"`
	PaymentID      *uint     `json:"payment_id"`
	RescheduleFrom *uint     `json:"reschedule_from"`
	Note           string    `json:"note"`
	TotalPrice     float64   `json:"total_price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BookingInfo struct untuk response endpoint teacher bookings
type BookingInfo struct {
	ID          uint    `json:"id"`
	StudentID   uint    `json:"student_id"`
	StudentName string  `json:"student_name"`
	Date        string  `json:"date"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Status      string  `json:"status"`
	Price       float64 `json:"price"`
}

// TeacherBookingResponse untuk response endpoint teacher bookings
type TeacherBookingResponse struct {
	BookingID   uint      `json:"booking_id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name"`
	Date        string    `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookingRequest struct {
	ScheduleID uint    `json:"schedule_id"`
	UserID     uint    `json:"user_id"`
	Note       string  `json:"note"`
	TotalPrice float64 `json:"total_price"`
}

type BookingRescheduleRequest struct {
	NewScheduleID uint `json:"schedule_id"`
	UserID        uint `json:"user_id"`
	BookingID     uint `json:"booking_id"`
}
