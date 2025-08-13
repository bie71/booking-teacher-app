package service

import (
	"fmt"
	"teacher/internal/infrastructure/booking"
	"teacher/internal/models"
	"teacher/internal/repository"
)

type DashboardService struct {
	teacherRepo    *repository.Repository
	serviceBooking *booking.BookingService
}

func NewDashboardService(teacherRepo *repository.Repository, serviceBooking *booking.BookingService) *DashboardService {
	return &DashboardService{
		teacherRepo:    teacherRepo,
		serviceBooking: serviceBooking,
	}
}

func (s *DashboardService) GetTeacherDashboard(teacherID uint) (*models.TeacherDashboardResponse, error) {
	// Get teacher profile
	teacher, err := s.teacherRepo.GetTeacherByID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	// Fetch real booking data from booking service
	bookings, err := s.fetchBookingsByTeacher(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	// Calculate dashboard stats from real data
	stats := s.calculateRealDashboardStats(bookings)

	// Get upcoming bookings (next 7 days)
	upcomingBookings := s.getRealUpcomingBookings(bookings)

	// Get completed lessons
	completedLessons := s.getRealCompletedLessons(bookings)

	// Get unique students
	recentStudents := s.getRealStudents(bookings)

	// Build response
	response := &models.TeacherDashboardResponse{
		TeacherProfile:   s.mapToTeacherResponse(teacher),
		Stats:            stats,
		UpcomingBookings: upcomingBookings,
		RecentStudents:   recentStudents,
		CompletedLessons: completedLessons,
	}

	return response, nil
}

func (s *DashboardService) fetchBookingsByTeacher(teacherID uint) ([]models.BookingInfo, error) {
	bookings, err := s.serviceBooking.GetBookingByTeacherId(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	bookingInfo := make([]models.BookingInfo, 0, len(bookings))
	for _, v := range bookings {
		bookingInfo = append(bookingInfo, models.BookingInfo{
			ID:          v.BookingID,
			StudentID:   v.StudentID,
			StudentName: v.StudentName,
			StartTime:   v.StartTime,
			EndTime:     v.EndTime,
			Status:      v.Status,
			Date:        v.Date,
			Price:       v.Price,
		})
	}

	return bookingInfo, nil
}

func (s *DashboardService) calculateRealDashboardStats(bookings []models.BookingInfo) models.DashboardStats {
	stats := models.DashboardStats{
		TotalStudents:    0,
		UpcomingBookings: 0,
		CompletedLessons: 0,
		TotalEarnings:    0,
	}

	// Track unique students
	studentSet := make(map[uint]bool)

	for _, booking := range bookings {
		studentSet[booking.StudentID] = true

		switch booking.Status {
		case "booked":
			stats.UpcomingBookings++
		case "completed":
			stats.CompletedLessons++
		}

		if booking.Status == "completed" {
			stats.TotalEarnings += booking.Price
		}
	}

	stats.TotalStudents = len(studentSet)
	return stats
}

func (s *DashboardService) getRealUpcomingBookings(bookings []models.BookingInfo) []models.BookingInfo {
	var upcoming []models.BookingInfo

	for _, booking := range bookings {
		if booking.Status == "booked" {
			upcoming = append(upcoming, booking)
		}
	}

	// Limit to next 5 bookings
	if len(upcoming) > 5 {
		return upcoming[:5]
	}
	return upcoming
}

func (s *DashboardService) getRealCompletedLessons(bookings []models.BookingInfo) []models.BookingInfo {
	var completed []models.BookingInfo

	for _, booking := range bookings {
		if booking.Status == "completed" {
			completed = append(completed, booking)
		}
	}

	// Limit to last 5 completed
	if len(completed) > 5 {
		return completed[len(completed)-5:]
	}
	return completed
}

func (s *DashboardService) getRealStudents(bookings []models.BookingInfo) []models.StudentInfo {
	studentMap := make(map[uint]*models.StudentInfo)

	for _, booking := range bookings {
		if _, exists := studentMap[booking.StudentID]; !exists {
			studentMap[booking.StudentID] = &models.StudentInfo{
				ID:           booking.StudentID,
				Name:         booking.StudentName,
				Email:        "", // Would need to fetch from user service
				TotalLessons: 0,
				TotalSpent:   0,
			}
		}

		student := studentMap[booking.StudentID]
		student.TotalLessons++
		if booking.Status == "completed" {
			student.TotalSpent += booking.Price
		}
	}

	var students []models.StudentInfo
	for _, student := range studentMap {
		students = append(students, *student)
	}

	// Limit to 5 recent students
	if len(students) > 5 {
		return students[:5]
	}
	return students
}

func (s *DashboardService) mapToTeacherResponse(teacher *models.Teacher) models.TeacherResponse {
	return models.TeacherResponse{
		ID:             teacher.ID,
		Name:           teacher.Name,
		Bio:            teacher.Bio,
		LanguageLevel:  teacher.LanguageLevel,
		PricePerHour:   teacher.PricePerHour,
		AvailableStart: teacher.AvailableStartTime,
		AvailableEnd:   teacher.AvailableEndTime,
		ProfileImage:   teacher.ProfileImage,
		CreatedAt:      teacher.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      teacher.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
