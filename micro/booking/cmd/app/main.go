package main

import (
	"booking/internal/config"
	"booking/internal/handler"
	"booking/internal/infrastructure/payment"
	"booking/internal/infrastructure/schedule"
	"booking/internal/infrastructure/supabase"
	"booking/internal/infrastructure/user"
	"booking/internal/middleware"
	model "booking/internal/model"
	"booking/internal/repository"
	"booking/internal/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	zerolog "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

func main() {
	c := config.LoadConfig()
	db := config.InitDB(c)
	// Auto migrate booking-related models to ensure the bookings table exists.
	{
		type Booking = model.Booking
		if err := db.AutoMigrate(&Booking{}); err != nil {
			zerolog.Info().Err(err).Msg("failed to auto migrate booking service database")
		}
	}
	gin.SetMode(c.GIN_MODE)

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	repo := repository.NewRepository(db)

	restyInit := resty.New()
	restyInit.SetDebug(cast.ToBool(os.Getenv("DEBUG")))

	supabaseService := supabase.InitUploadClient(&c.Client, restyInit)

	serviceSchedule := schedule.NewScheduleHttp(c.ServiceSchedule, restyInit)

	userService := user.NewUserService(restyInit, c.ServiceUser)

	paymentService := payment.NewPaymentHttp(c.ServicePayment, restyInit)

	service := service.NewService(repo, serviceSchedule, userService, paymentService)

	uploadHandler := handler.NewUploadHandler(supabaseService, &c.Client)

	handler := handler.NewHandler(service)

	api := r.Group("/api/v1")
	if !c.IsNFT {
		api.Use(middleware.AuthMiddleware(&c.JWT))
	}
	auth := api.Group("")
	auth.POST("/auth/logout", handler.Logout)
	auth.GET("/admin/bookings", handler.GetBookingsAdmin)
	{
		api.POST("/bookings", handler.CreateBooking)
		api.GET("/bookings", handler.GetBookings)
		api.GET("/booking/:id", handler.GetBooking)
		api.DELETE("/booking/:id", handler.DeleteBooking)
		api.GET("/booking-detail/:id", handler.GetBookingDetail)
		api.POST("/bookings/:id/reschedule", handler.RescheduleBooking)
		api.POST("/bookings/:id/cancel", handler.CancelBooking)
		api.GET("/bookings/user/:user_id", handler.GetBookingsByUserID)
		// Endpoint to fetch upcoming lessons for a user.  Returns the list of
		// upcoming bookings (paid status with future schedule) for the given
		// user ID.  Optional query parameter "limit" can specify how many
		// upcoming lessons to retrieve (default 5).
		api.GET("/bookings/user/:user_id/upcoming-lessons", handler.GetUpcomingLessons)
		api.GET("/bookings/teacher/:teacher_id", handler.GetTeacherBookings)
		api.PUT("/bookings/:id/status", handler.UpdateBookingStatus)
		api.PUT("/booking-change/:id/status/:status", handler.ChangeStatus)

		api.POST("/upload-image", uploadHandler.UploadHandler)
		api.GET("/total-bookings", handler.TotalPricePaidBookings)
	}

	// Internal routes that do not require authentication. These endpoints
	// are intended for service-to-service calls (e.g. teacher service or payment
	// service fetching bookings) and should not be exposed publicly. They are
	// deliberately placed outside the authentication middleware group.
	// Endpoint to fetch teacher bookings for dashboard and other services
	r.GET("/api/v1/internal/bookings/teacher/:teacher_id", handler.GetTeacherBookingsInternal)

	// Endpoint to fetch a single booking by ID (including user ID) for internal
	// services.  This internal endpoint exposes the raw booking record with
	// associated user ID so that other services (e.g. payment) can determine
	// the owner of a booking without requiring authentication.  It should not
	// be exposed to clients directly.
	r.GET("/api/v1/internal/bookings/:id", handler.GetBookingInternal)

	r.PUT("/private/bookings/:id/status", handler.UpdateBookingStatus)

	zerolog.Info().Msg("Server running on port " + fmt.Sprint(":", c.AppPort))

	r.Run(fmt.Sprint(":", c.AppPort))

}
