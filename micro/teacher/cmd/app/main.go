package main

import (
	"fmt"
	"os"
	"teacher/internal/config"
	"teacher/internal/handler"
	"teacher/internal/infrastructure/booking"
	"teacher/internal/infrastructure/supabase"
	"teacher/internal/middleware"
	"teacher/internal/models"
	"teacher/internal/repository"
	"teacher/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

func main() {
	c := config.LoadConfig()
	db := config.InitDB(c)
	// Auto migrate teacher-related models. Creates or updates tables based on
	// the current model definitions. Ensures that tables such as teachers
	// and schedules exist when the service starts.
	{
		type (
			Teacher  = models.Teacher
			Schedule = models.Schedule
		)
		if err := db.AutoMigrate(&Teacher{}, &Schedule{}); err != nil {
			log.Info().Err(err).Msg("failed to auto migrate teacher service database")
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

	teacherRepo := repository.NewRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	restyInit := resty.New()
	restyInit.SetDebug(cast.ToBool(os.Getenv("DEBUG")))

	supabaseService := supabase.InitUploadClient(&c.Client, restyInit)

	bookingService := booking.NewBookingService(restyInit, &c.ServiceBooking)

	teacherService := service.NewService(teacherRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	dashboardService := service.NewDashboardService(teacherRepo, bookingService)

	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	handlers := handler.NewHandler(teacherService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	uploadHandler := handler.NewUploadHandler(supabaseService, &c.Client)

	api := r.Group("/api/v1")
	{
		api.GET("/teachers", handlers.GetTeachers)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(&c.JWT))
		auth.GET("/teachers/me", handlers.GetMe)
		api.POST("/teachers", handlers.CreateTeacher)
		api.GET("/teachers/:id", handlers.GetTeacher)
		api.PUT("/teachers/:id", handlers.UpdateTeacher)
		api.DELETE("/teachers/:id", handlers.DeleteTeacher)

		api.GET("/total-teachers", handlers.CountTeachers)

		// api.POST("/schedule", scheduleHandler.BookSchedule)
		api.GET("/schedule/teacher/:teacher_id", scheduleHandler.GetScheduleAvailable)
		api.PUT("/cancel-schedule/:id", scheduleHandler.CancelSchedule)
		api.POST("/schedule", scheduleHandler.CreateSchedule)
		api.PUT("/schedule-status", scheduleHandler.UpdateScheduleStatus)
		api.PUT("/schedule/:id", scheduleHandler.UpdateSchedule)
		api.POST("/schedule/batch-detail", scheduleHandler.GetBatchScheduleDetail)
		api.GET("/schedules", scheduleHandler.GetSchedules)
		api.GET("/schedule/:id", scheduleHandler.GetScheduleById)
		api.DELETE("/schedule/:id", scheduleHandler.DeleteSchedule)
		api.POST("schedule/filter-by-teacher", scheduleHandler.FilterByTeacher)

		api.GET("/teachers/dashboard/:teacher_id", dashboardHandler.GetTeacherDashboard)

		api.POST("/upload-image", uploadHandler.UploadHandler)
	}

	log.Info().Msgf("Starting server on port %s", c.AppPort)

	r.Run(fmt.Sprint(":", c.AppPort))
}
