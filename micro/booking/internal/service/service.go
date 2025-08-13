package service

import (
	"booking/internal/infrastructure/payment"
	"booking/internal/infrastructure/schedule"
	"booking/internal/infrastructure/user"
	"booking/internal/model"
	"booking/internal/pkg"
	"booking/internal/repository"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"sort"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Service struct {
	bookingRepository *repository.Repository
	serviceHttp       *schedule.ScheduleHttp
	serviceUser       *user.UserService
	servicePayment    *payment.Payment
}

func NewService(
	bookingRepository *repository.Repository,
	serviceHttp *schedule.ScheduleHttp,
	serviceUser *user.UserService,
	servicePayment *payment.Payment,
) *Service {
	return &Service{
		bookingRepository: bookingRepository,
		serviceHttp:       serviceHttp,
		serviceUser:       serviceUser,
		servicePayment:    servicePayment,
	}
}

func (s *Service) CreateBooking(c *gin.Context, req model.BookingRequest) (*model.Booking, error) {
	// Check if schedule exists and is available
	schedule, err := s.serviceHttp.CheckScheduleAvailability(req.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("schedule not available: %w", err)
	}

	if schedule.Status != "available" {
		return nil, fmt.Errorf("schedule is not available for booking")
	}

	// Check for existing booking by same user for this schedule
	if !cast.ToBool(os.Getenv("IS_NFT")) {
		existing, err := s.bookingRepository.GetBookingsByUserIDAndId(req.UserID, req.ScheduleID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("you already have a booking for this schedule")
		}
	}

	booking := model.Booking{
		UserID:     req.UserID,
		ScheduleID: schedule.ID,
		Note:       req.Note,
		Status:     "pending",
		TotalPrice: req.TotalPrice,
	}

	if err := s.bookingRepository.CreateBooking(&booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	statusSchedule := "booked"
	if cast.ToBool(os.Getenv("IS_NFT")) {
		statusSchedule = "available"
	}
	err = s.serviceHttp.UpdateScheduleStatus(c, schedule.ID, statusSchedule)
	if err != nil {
		go s.bookingRepository.DeleteBooking(booking.ID)
		return nil, fmt.Errorf("failed to update schedule status: %w", err)
	}

	return &booking, nil
}

func (s *Service) RescheduleBooking(c *gin.Context, bookingID, newScheduleID, userID uint) error {
	booking, err := s.bookingRepository.GetBookingsByUserIDAndId(userID, bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.Status != "paid" && booking.Status != "pending" {
		return errors.New("cannot reschedule this booking")
	}

	_, err = s.serviceHttp.CheckScheduleAvailability(newScheduleID)
	if err != nil {
		return errors.New("new schedule not available")
	}

	if err := s.serviceHttp.UpdateScheduleStatus(c, booking.ScheduleID, "available"); err != nil {
		return errors.New("failed to free previous schedule")
	}
	if err := s.serviceHttp.UpdateScheduleStatus(c, newScheduleID, "booked"); err != nil {
		return errors.New("failed to book new schedule")
	}

	booking.ScheduleID = newScheduleID
	booking.Status = "rescheduled"
	booking.RescheduleFrom = &booking.ID
	booking.UpdatedAt = time.Now()

	if err := s.bookingRepository.UpdateBooking(booking); err != nil {
		return errors.New("failed to update booking")
	}

	return nil
}

func (s *Service) CancelBookingByID(c *gin.Context, id uint) (*model.Booking, error) {
	booking, err := s.bookingRepository.GetBooking(id)
	if err != nil {
		return nil, fmt.Errorf("booking not found")
	}

	if booking.Status != "pending" && booking.Status != "paid" {
		return nil, fmt.Errorf("only pending or paid bookings can be cancelled")
	}

	booking.Status = "cancelled"
	booking.UpdatedAt = time.Now()

	if err := s.bookingRepository.UpdateBooking(booking); err != nil {
		return nil, err
	}

	err = s.serviceHttp.CancelSchedule(c, booking.ScheduleID)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) GetUserBookings(c *gin.Context, userID uint, isAdmin bool, pg pkg.Pagination) (model.PaginatedBookingsResponse, error) {

	page, _ := strconv.Atoi(pg.PageStr)
	limit, _ := strconv.Atoi(pg.LimitStr)

	if limit == 0 {
		limit = 10
	}

	if page == 0 {
		page = 1
	}

	var startDate, endDate *time.Time
	layout := "2006-01-02"
	if pg.StartDateStr != "" {
		t, err := time.Parse(layout, pg.StartDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if pg.EndDateStr != "" {
		t, err := time.Parse(layout, pg.EndDateStr)
		if err == nil {
			endDate = &t
		}
	}

	totalRecords, err := s.bookingRepository.CountBookings(userID, false, pg.Status, startDate, endDate)
	if err != nil {
		return model.PaginatedBookingsResponse{}, err
	}

	bookings, err := s.bookingRepository.GetBookingsByUserID(userID, isAdmin, pg.Status, startDate, endDate, page, limit)
	if err != nil {
		return model.PaginatedBookingsResponse{}, err
	}

	scheduleIds := make([]uint, 0, len(bookings))
	for _, v := range bookings {
		scheduleIds = append(scheduleIds, v.ScheduleID)
	}

	schedules, err := s.serviceHttp.FetchSchedulesByIDs(c, scheduleIds)
	if err != nil {
		return model.PaginatedBookingsResponse{}, err
	}

	var mappedBookings []model.BookingResponse
	if pg.TeacherID != "" {
		result, err := s.filterByTeacherID(
			s.mapBookingResponse(bookings, schedules),
			cast.ToUint(pg.TeacherID),
		)
		if err != nil {
			return model.PaginatedBookingsResponse{}, err
		}
		mappedBookings = result
	} else {
		mappedBookings = s.mapBookingResponse(bookings, schedules)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	return model.PaginatedBookingsResponse{
		Data: mappedBookings,
		Pagination: model.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      totalRecords,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) mapBookingResponse(bookings []model.Booking, schedulesMap map[uint]model.ScheduleResponse) []model.BookingResponse {
	var result []model.BookingResponse

	for _, booking := range bookings {

		user, err := s.serviceUser.GetUserById(booking.UserID)
		if err != nil {
			continue
		}

		resp := model.BookingResponse{
			ID:         booking.ID,
			Status:     booking.Status,
			Note:       booking.Note,
			CreatedAt:  booking.CreatedAt,
			TotalPrice: booking.TotalPrice,
			User:       &user,
		}

		if schedule, ok := schedulesMap[booking.ScheduleID]; ok {
			resp.Schedule = &schedule
		}

		result = append(result, resp)
	}

	return result
}

func (s *Service) filterByTeacherID(bookings []model.BookingResponse, teacherID uint) ([]model.BookingResponse, error) {
	scheduleIDs := make([]string, len(bookings))
	for i, b := range bookings {
		scheduleIDs[i] = cast.ToString(b.Schedule.ID)
	}

	scheduleResp, err := s.serviceHttp.CallScheduleServiceToGetByTeacher(cast.ToString(teacherID), scheduleIDs)
	if err != nil {
		return nil, err
	}

	validScheduleMap := make(map[string]bool)
	for _, id := range scheduleResp.ValidScheduleIDs {
		validScheduleMap[id] = true
	}

	filteredBookings := make([]model.BookingResponse, 0)
	for _, b := range bookings {
		if validScheduleMap[cast.ToString(b.Schedule.ID)] {
			filteredBookings = append(filteredBookings, b)
		}
	}

	return filteredBookings, nil
}

func (s *Service) UpdateBookingStatus(id uint, paymentID uint, status string) (*model.Booking, error) {
	booking, err := s.bookingRepository.GetBooking(id)
	if err != nil {
		return nil, fmt.Errorf("booking not found")
	}

	booking.Status = status
	booking.PaymentID = &paymentID
	booking.UpdatedAt = time.Now()

	if err := s.bookingRepository.UpdateBooking(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) GetBookings(pg pkg.Paginate) (pkg.ResponsePaginate, error) {
	return s.bookingRepository.GetBookings(pg)
}

// GetUpcomingLessons returns a list of upcoming lessons (bookings) for a given user.
// A lesson is considered upcoming if the booking status is "paid" and the
// schedule date/time is in the future relative to the current time.  This
// function will fetch up to the specified limit of upcoming bookings,
// sorted in ascending order by the schedule date/time.  If limit is 0 or
// negative, a default of 5 will be used.
func (s *Service) GetUpcomingLessons(c *gin.Context, userID uint, limit int) ([]model.BookingResponse, error) {
	if limit <= 0 {
		limit = 5
	}

	// Fetch paid bookings for the user.  We request a generous limit to
	// ensure we have enough bookings to filter for upcoming lessons.
	// Note: offset 0 and limit 100 should suffice for most use cases.  If
	// needed, this can be adjusted or made configurable.
	bookings, err := s.bookingRepository.GetBookingsByUserID(userID, false, "paid", nil, nil, 0, 100)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return []model.BookingResponse{}, nil
	}

	// Collect schedule IDs to fetch schedule details in batch
	scheduleIDs := make([]uint, 0, len(bookings))
	for _, b := range bookings {
		scheduleIDs = append(scheduleIDs, b.ScheduleID)
	}

	// Fetch schedule details from the schedule service.  Pass the
	// gin.Context so that any auth tokens on the request can be forwarded.
	schedules, err := s.serviceHttp.FetchSchedulesByIDs(c, scheduleIDs)
	if err != nil {
		return nil, err
	}

	// Map bookings to responses (includes user and schedule details)
	mapped := s.mapBookingResponse(bookings, schedules)

	// Filter out bookings whose schedule date/time is not in the future.
	upcoming := make([]model.BookingResponse, 0)
	now := time.Now()
	for _, booking := range mapped {
		if booking.Schedule == nil {
			continue
		}
		// Parse schedule date and start time.  The schedule service
		// returns date in YYYY-MM-DD and start_time in HH:MM format.  We
		// combine them into a single timestamp for comparison.
		// If parsing fails, skip the booking.
		dateStr := booking.Schedule.Date
		timeStr := booking.Schedule.StartTime
		dt, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", dateStr, timeStr))
		if err != nil {
			continue
		}
		if dt.After(now) && booking.Status == "paid" {
			upcoming = append(upcoming, booking)
		}
	}

	// Sort upcoming bookings by schedule date/time ascending
	sort.Slice(upcoming, func(i, j int) bool {
		a := upcoming[i]
		b := upcoming[j]
		if a.Schedule == nil || b.Schedule == nil {
			return false
		}
		at, err1 := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", a.Schedule.Date, a.Schedule.StartTime))
		bt, err2 := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", b.Schedule.Date, b.Schedule.StartTime))
		if err1 != nil || err2 != nil {
			return false
		}
		return at.Before(bt)
	})

	// Limit the number of upcoming lessons returned
	if len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}
	return upcoming, nil
}

func (s *Service) CheckBookingExists(id uint) (*model.BookingResponse, error) {
	booking, error := s.bookingRepository.GetBooking(id)

	if error != nil {
		return nil, errors.New("booking not found")
	}

	return &model.BookingResponse{
		ID:        booking.ID,
		Status:    booking.Status,
		Note:      booking.Note,
		CreatedAt: booking.CreatedAt,
	}, nil
}

func (s *Service) BookingDetail(c *gin.Context, id uint) (*model.BookingResponse, error) {
	booking, error := s.bookingRepository.GetBooking(id)
	if error != nil {
		return nil, errors.New("booking not found")
	}

	schedule, err := s.serviceHttp.GetScheduleByID(booking.ScheduleID)
	if err != nil {
		return nil, err
	}

	user, err := s.serviceUser.GetUserById(booking.UserID)

	if err != nil {
		return nil, err
	}

	payment := &payment.PaymentResponse{}
	if booking.PaymentID != nil {
		payment, err = s.servicePayment.GetPaymentDetail(c, *booking.PaymentID)
		if err != nil {
			return nil, err
		}

	}

	return &model.BookingResponse{
		ID:         booking.ID,
		Status:     booking.Status,
		Note:       booking.Note,
		CreatedAt:  booking.CreatedAt,
		TotalPrice: booking.TotalPrice,
		Schedule:   schedule,
		User:       &user,
		Payment:    payment,
	}, nil
}

// GetRawBooking retrieves the raw Booking model by its ID.  This method is
// intended for internal use by other services (e.g. payment) that need
// access to the underlying booking fields, such as the UserID associated
// with a booking.  If the booking does not exist, it returns an error
// indicating the booking was not found.
func (s *Service) GetRawBooking(id uint) (*model.Booking, error) {
	booking, err := s.bookingRepository.GetBooking(id)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}

// GetTeacherBookings mengambil semua bookings untuk teacher tertentu
func (s *Service) GetTeacherBookings(c *gin.Context, teacherID uint, pg pkg.Pagination) ([]model.TeacherBookingResponse, error) {
	page, _ := strconv.Atoi(pg.PageStr)
	limit, _ := strconv.Atoi(pg.LimitStr)
	offset := (page - 1) * limit

	var startDate, endDate *time.Time
	layout := "2006-01-02"
	if pg.StartDateStr != "" {
		t, err := time.Parse(layout, pg.StartDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if pg.EndDateStr != "" {
		t, err := time.Parse(layout, pg.EndDateStr)
		if err == nil {
			endDate = &t
		}
	}

	// Get all bookings first
	allBookings, err := s.bookingRepository.GetBookings(pkg.Paginate{
		Page:  1,
		Limit: 1000, // Get all bookings to filter by teacher
	})
	if err != nil {
		return nil, err
	}

	bookings := allBookings.Data.([]model.Booking)
	if len(bookings) == 0 {
		return []model.TeacherBookingResponse{}, nil
	}

	// Get schedule IDs
	scheduleIds := make([]uint, 0, len(bookings))
	for _, v := range bookings {
		scheduleIds = append(scheduleIds, v.ScheduleID)
	}

	// Get schedule details
	schedules, err := s.serviceHttp.FetchSchedulesByIDs(c, scheduleIds)
	if err != nil {
		return nil, err
	}

	// Filter by teacher ID and status/date
	var responses []model.TeacherBookingResponse
	for _, booking := range bookings {

		user, err := s.serviceUser.GetUserById(booking.UserID)
		if err != nil {
			return nil, err
		}

		schedule, ok := schedules[booking.ScheduleID]
		if !ok || schedule.Teacher.ID != teacherID {
			continue
		}

		// Filter by status if provided
		if pg.Status != "" && booking.Status != pg.Status {
			continue
		}

		// Filter by date range if provided
		if startDate != nil {
			scheduleDate, err := time.Parse("2006-01-02", schedule.Date)
			if err == nil && scheduleDate.Before(*startDate) {
				continue
			}
		}
		if endDate != nil {
			scheduleDate, err := time.Parse("2006-01-02", schedule.Date)
			if err == nil && scheduleDate.After(*endDate) {
				continue
			}
		}

		response := model.TeacherBookingResponse{
			BookingID:   booking.ID,
			StudentID:   booking.UserID,
			StudentName: user.Name,
			Date:        schedule.Date,
			StartTime:   schedule.StartTime,
			EndTime:     schedule.EndTime,
			Status:      booking.Status,
			Price:       schedule.TotalPrice,
			CreatedAt:   booking.CreatedAt,
		}

		responses = append(responses, response)
	}

	// Apply pagination
	if offset >= len(responses) {
		return []model.TeacherBookingResponse{}, nil
	}

	end := offset + limit
	if end > len(responses) {
		end = len(responses)
	}

	return responses[offset:end], nil
}

func (s *Service) DeleteBooking(c *gin.Context, id uint) error {

	booking, err := s.bookingRepository.GetBooking(id)
	if err != nil {
		return err
	}

	err = s.bookingRepository.DeleteBooking(id)
	if err != nil {
		return err
	}

	go s.serviceHttp.UpdateScheduleStatus(c, booking.ScheduleID, "available")

	return nil
}

func (s *Service) ChangeStatus(id uint, status string) error {
	return s.bookingRepository.ChangeStatus(id, status)
}

func (s *Service) TotalPricePaidBookings() (float64, int64, error) {

	revenue, err := s.bookingRepository.SumTotalPricePaidBookings()

	if err != nil {
		return 0, 0, err
	}

	count, err := s.bookingRepository.CountBookings(0, true, "", nil, nil)

	if err != nil {
		return 0, 0, err
	}

	return revenue, int64(count), err

}

func (s *Service) GetEnrichedBookings(c *gin.Context, pg pkg.Paginate, status string, q string, startDate, endDate string) (model.PaginatedBookingsResponse, error) {
	// If q provided, fetch larger set then filter in-memory for cross-service fields
	basePg := pg
	if q != "" {
		basePg.Page = 1
		if basePg.Limit < 200 {
			basePg.Limit = 200
		}
	}
	resp, err := s.bookingRepository.GetBookingsFiltered(basePg, status, startDate, endDate)
	if err != nil {
		return model.PaginatedBookingsResponse{}, err
	}
	bookings, ok := resp.Data.([]model.Booking)
	if !ok {
		return model.PaginatedBookingsResponse{}, errors.New("invalid data type from repository")
	}

	// Collect schedule IDs
	scheduleIDs := make([]uint, 0, len(bookings))
	for _, b := range bookings {
		scheduleIDs = append(scheduleIDs, b.ScheduleID)
	}
	schedulesMap, _ := s.serviceHttp.FetchSchedulesByIDs(c, scheduleIDs)
	items := s.mapBookingResponse(bookings, schedulesMap)

	// Enrich with payment details (per booking)
	for i, b := range bookings {
		if b.PaymentID != nil {
			pay, err := s.servicePayment.GetPaymentDetail(c, *b.PaymentID)
			if err == nil && i < len(items) {
				items[i].Payment = pay
			}
		}
	}
	/*PAYMENT_ENRICH*/
	if q != "" {
		qLower := strings.ToLower(q)
		filtered := make([]model.BookingResponse, 0, len(items))
		for _, it := range items {
			match := false
			if it.User != nil {
				if strings.Contains(strings.ToLower(it.User.Name), qLower) || strings.Contains(strings.ToLower(it.User.Email), qLower) {
					match = true
				}
			}
			if !match && it.Schedule.Teacher != nil {
				if strings.Contains(strings.ToLower(it.Schedule.Teacher.Name), qLower) {
					match = true
				}
			}
			if !match && it.Payment != nil {
				if strings.Contains(strings.ToLower(it.Payment.PaymentMethod), qLower) || strings.Contains(strings.ToLower(it.Payment.Status), qLower) {
					match = true
				}
			}
			if !match {
				if strings.Contains(strings.ToLower(fmt.Sprintf("%d", it.ID)), qLower) {
					match = true
				}
			}
			if !match && it.Schedule != nil {
				if strings.Contains(strings.ToLower(fmt.Sprintf("%d", it.Schedule.ID)), qLower) {
					match = true
				}
			}
			if match {
				filtered = append(filtered, it)
			}
		}
		total := len(filtered)
		if pg.Page < 1 {
			pg.Page = 1
		}
		totalPages := int(math.Ceil(float64(total) / float64(pg.Limit)))
		start := (pg.Page - 1) * pg.Limit
		end := start + pg.Limit
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		items = filtered[start:end]
		return model.PaginatedBookingsResponse{Data: items, Pagination: model.PaginationMeta{Page: pg.Page, Limit: pg.Limit, Total: total, TotalPages: totalPages}}, nil
	}

	return model.PaginatedBookingsResponse{Data: items, Pagination: model.PaginationMeta{Page: pg.Page, Limit: pg.Limit, Total: resp.Pagination.TotalData, TotalPages: resp.Pagination.TotalPage}}, nil
}
