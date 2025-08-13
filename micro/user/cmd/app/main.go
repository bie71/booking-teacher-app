package main

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/infrastructure/statistic"
	"auth/internal/infrastructure/supabase"
	"auth/internal/middleware"
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/service"
	"auth/seeder"
	"fmt"
	"os"
	"time"

	// Import models for auto migration

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	zerolog "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

func main() {
	c := config.LoadConfig()
	db := config.InitDB(c)
	// Automatically migrate database schemas for user-related models. This
	// ensures that tables are created/updated based on the Go models when
	// the service starts. It is idempotent and safe to run on each
	// startup.  Add new models here as they are introduced.
	{
		// Import model types anonymously to avoid unused import errors
		// and ensure AutoMigrate has type information.
		type (
			User            = models.User
			ActivityLog     = models.ActivityLog
			FavoriteTeacher = models.FavoriteTeacher
		)
		if err := db.AutoMigrate(&User{}, &ActivityLog{}, &FavoriteTeacher{}); err != nil {
			zerolog.Info().Err(err).Msg("failed to auto migrate user service database")
		}
	}
	seeder.SeedUsers(db)
	gin.SetMode(c.GIN_MODE)

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userRepo := repository.NewUserRepository(db)

	imageRepo := repository.NewAdminRepository(db)

	authService := service.NewJWTConfig(c.JWT.SecretKey, time.Duration(c.JWT.TokenDuration)*time.Hour)

	emailService := service.NewEmailService(
		c.SMTP.Host,
		c.SMTP.Port,
		c.SMTP.TimeoutDuration,
		c.SMTP.Username,
		c.SMTP.Password,
		c.SMTP.FromEmail,
		c.SMTP.TemplatePath,
		c.SMTP.TemplateLogoURL,
		c.SMTP.InsecureSkipVerify,
		c.SMTP.UseTLS,
	)

	restyInit := resty.New()
	restyInit.SetDebug(cast.ToBool(os.Getenv("DEBUG")))

	statistic := statistic.NewStatistic(&c.ServiceBooking, &c.ServiceTeacher, restyInit)

	// Initialize repositories for activity logs and favorites
	activityRepo := repository.NewActivityLogRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	// Instantiate the user service with additional repositories to enable
	// logging of recent actions and managing favorite teachers. Passing the
	// extra dependency ensures CreateActivityLog, GetRecentActivity, and
	// favorite teacher methods work as expected.
	userService := service.NewUserService(userRepo, emailService, c.JWT.TokenDuration, authService, statistic, activityRepo, favoriteRepo)
	userHandler := handler.NewHandler(userService, userRepo)

	initSupabase := supabase.InitUploadClient(&c.Client, restyInit)

	uploadHandler := handler.NewUploadHandler(initSupabase, &c.Client)

	adminHandler := handler.NewAdminHandler(imageRepo)
	userAdminHandler := handler.NewUserAdminHandler(userService, userRepo)

	v1 := r.Group("/api/v1")
	v1.POST("/register", userHandler.RegisterUser)
	v1.POST("/login", userHandler.LoginUser)
	v1.POST("/forgot-password", userHandler.ForgotPassword)
	v1.GET("/verify-reset-token", userHandler.VerifyResetToken)
	v1.POST("/reset-password", userHandler.ResetPassword)

	v1.GET("/profile", userHandler.GetProfile)
	v1.GET("/hero-image", adminHandler.GetHeroImage)

	// Admin routes
	admin := r.Group("/api/admin")
	if !c.IsNFT {
		admin.Use(middleware.AuthMiddleware(authService))
	}
	admin.POST("/upload-hero-image", uploadHandler.UploadHandler)
	admin.POST("/hero-image", adminHandler.SaveHeroImage)
	admin.DELETE("/hero-image", adminHandler.DeleteHeroImage)
	admin.GET("/dashboard-stats", userHandler.GetDashboardStats)

	// User management routes
	admin.GET("/users", userAdminHandler.GetUsers)
	admin.POST("/users", userAdminHandler.CreateUser)
	admin.GET("/users/:id", userAdminHandler.GetUser)
	admin.PUT("/users/:id", userAdminHandler.UpdateUser)
	admin.DELETE("/users/:id", userAdminHandler.DeleteUser)

	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(authService))
	auth.GET("/me", userHandler.GetMe)
	auth.PUT("/profile", userHandler.UpdateProfile)
	auth.POST("/change-password", userHandler.ChangePassword)
	auth.POST("/upload-image", uploadHandler.UploadHandler)

	// Activity endpoints for logging and fetching recent user activities
	auth.POST("/activity", userHandler.LogActivity)
	auth.GET("/activity/recent", userHandler.GetRecentActivity)

	// Favorite teachers endpoints for adding/removing favorites and
	// retrieving the current user's favorites
	auth.POST("/favorites", userHandler.ToggleFavorite)
	auth.GET("/favorites", userHandler.GetFavorites)

	// Internal endpoint for service-to-service activity logging.  This
	// endpoint bypasses authentication and should only be used by other
	// services (e.g. payment service) to record activities on behalf of
	// users.  It accepts a JSON payload with user_id, action and
	// description.
	r.POST("/api/v1/internal/activity", userHandler.LogActivityInternal)

	zerolog.Info().Msg("Starting server on port " + fmt.Sprint(c.AppPort))

	r.Run(fmt.Sprint(":", c.AppPort)) // default port from .env handled inside gin or set manually with ":8001"
}
