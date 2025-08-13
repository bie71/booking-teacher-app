package models

import "time"

type TeacherDashboardResponse struct {
	TeacherProfile   TeacherResponse `json:"teacher_profile"`
	Stats            DashboardStats  `json:"stats"`
	UpcomingBookings []BookingInfo   `json:"upcoming_bookings"`
	RecentStudents   []StudentInfo   `json:"recent_students"`
	CompletedLessons []BookingInfo   `json:"completed_lessons"`
}

type DashboardStats struct {
	TotalStudents    int     `json:"total_students"`
	UpcomingBookings int     `json:"upcoming_bookings"`
	CompletedLessons int     `json:"completed_lessons"`
	TotalEarnings    float64 `json:"total_earnings"`
}

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

type StudentInfo struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	TotalLessons int     `json:"total_lessons"`
	TotalSpent   float64 `json:"total_spent"`
}

type Response struct {
	Message  string                   `json:"message"`
	Bookings []TeacherBookingResponse `json:"bookings"`
}

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
