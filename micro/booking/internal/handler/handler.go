package handler

import (
	"booking/internal/model"
	"booking/internal/pkg"
	"booking/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateBooking(c *gin.Context) {
	var req model.BookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.service.CreateBooking(c, req)
	if err != nil {
		if err.Error() == "schedule is not available" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully", "booking": resp})
}

func (h *Handler) RescheduleBooking(c *gin.Context) {
	var req model.BookingRescheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.RescheduleBooking(c, req.BookingID, req.NewScheduleID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reschedule booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking rescheduled successfully"})
}

func (h *Handler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	booking, err := h.service.CancelBookingByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
		"data":    booking,
	})
}

func (h *Handler) GetBookingsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pg := pkg.Pagination{
		PageStr:      c.Query("page"),
		LimitStr:     c.Query("limit"),
		StartDateStr: c.Query("start_date"),
		EndDateStr:   c.Query("end_date"),
		Status:       c.Query("status"),
		TeacherID:    c.Query("teacher_id"),
	}

	isAdmin := false
	value, _ := c.Get("role")
	if cast.ToString(value) == "admin" {
		isAdmin = true
	}

	bookings, err := h.service.GetUserBookings(c, uint(userID), isAdmin, pg)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	status := c.Query("status")
	paymentid := c.Query("paymentId")
	booking, err := h.service.UpdateBookingStatus(uint(id), cast.ToUint(paymentid), status)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated successfully",
		"data":    booking,
	})
}

func (h *Handler) GetBookings(c *gin.Context) {
	pg := pkg.Pagination{
		PageStr:  c.Query("page"),
		LimitStr: c.Query("limit"),
	}

	paginate := pkg.Paginate{
		Page:  cast.ToInt(pg.PageStr),
		Limit: cast.ToInt(pg.LimitStr),
	}

	bookings, err := h.service.GetBookings(paginate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// GetTeacherBookingsInternal handles requests to retrieve bookings for a teacher without
// requiring authentication. This endpoint is intended for internal service-to-service
// communication (e.g. teacher dashboard in the teacher microservice) and should not
// be exposed to public clients. It accepts the same query parameters as
// GetTeacherBookings: page, limit, start_date, end_date and status.
func (h *Handler) GetTeacherBookingsInternal(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID"})
		return
	}

	pg := pkg.Pagination{
		PageStr:      c.Query("page"),
		LimitStr:     c.Query("limit"),
		StartDateStr: c.Query("start_date"),
		EndDateStr:   c.Query("end_date"),
		Status:       c.Query("status"),
	}

	// Directly call the service to retrieve teacher bookings without checking role.
	bookings, err := h.service.GetTeacherBookings(c, uint(teacherID), pg)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// Wrap the slice in a response structure to match the shape expected by
	// downstream services (e.g. teacher service).  Returning an object with
	// a "bookings" key allows the caller to unmarshal into its own
	// Response struct.
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"bookings": bookings,
	})
}

// GetBookingInternal handles internal requests to retrieve a booking record
// by its ID.  This endpoint is intended for service-to-service calls
// (e.g. payment service) that need access to the raw booking model including
// the associated user ID.  It does not require authentication and returns
// a minimal JSON payload containing the booking fields used by other
// services.  If the booking is not found, it returns HTTP 404.
func (h *Handler) GetBookingInternal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	booking, err := h.service.GetRawBooking(uint(id))
	if err != nil {
		// Return 404 if booking not found; propagate other errors as 500
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Respond with a sanitized version of the booking containing the fields
	// required by internal consumers.  We avoid exposing internal fields
	// unnecessarily.  The struct literal here is anonymous; the JSON keys
	// correspond to the field names in the Booking model.
	c.JSON(http.StatusOK, gin.H{"booking": gin.H{
		"id":          booking.ID,
		"user_id":     booking.UserID,
		"schedule_id": booking.ScheduleID,
		"status":      booking.Status,
		"payment_id":  booking.PaymentID,
		"total_price": booking.TotalPrice,
	}})
}

func (h *Handler) GetBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	booking, err := h.service.CheckBookingExists(uint(id))
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

// GetTeacherBookings handler untuk endpoint teacher bookings
func (h *Handler) GetTeacherBookings(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID"})
		return
	}

	pg := pkg.Pagination{
		PageStr:      c.Query("page"),
		LimitStr:     c.Query("limit"),
		StartDateStr: c.Query("start_date"),
		EndDateStr:   c.Query("end_date"),
		Status:       c.Query("status"),
	}

	bookings, err := h.service.GetTeacherBookings(c, uint(teacherID), pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"message":  "Teacher bookings retrieved successfully",
	})
}

// GetUpcomingLessons handles requests to retrieve upcoming lessons for a user.
// It expects a user ID as a URL parameter and an optional "limit" query
// parameter specifying how many upcoming lessons to return (default 5 if
// unspecified or invalid).  The handler calls the booking service to
// fetch upcoming lessons and returns them in a JSON response.
func (h *Handler) GetUpcomingLessons(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse limit from query parameter.  Use default 5 on error or if
	// non-positive value provided.
	limit := 5
	if q := c.Query("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 {
			limit = v
		}
	}

	upcoming, err := h.service.GetUpcomingLessons(c, uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming lessons"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": upcoming})
}

func (h *Handler) GetBookingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	booking, err := h.service.BookingDetail(c, uint(id))
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

func (h *Handler) DeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteBooking(c, uint(id))
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	status := c.Param("status")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.ChangeStatus(uint(id), strings.ToLower(status))
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}

func (h *Handler) TotalPricePaidBookings(c *gin.Context) {
	totalRevenue, totalBooking, err := h.service.TotalPricePaidBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_bookings": totalBooking, "total_revenue": totalRevenue})
}


func (h *Handler) GetBookingsAdmin(c *gin.Context) {
	pg := pkg.Pagination{ PageStr: c.Query("page"), LimitStr: c.Query("limit") }
	paginate := pkg.Paginate{ Page: cast.ToInt(pg.PageStr), Limit: cast.ToInt(pg.LimitStr) }
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	q := c.Query("q")
	res, err := h.service.GetEnrichedBookings(c, paginate, status, q, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookings": res.Data, "pagination": res.Pagination})
}
