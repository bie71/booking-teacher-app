package main

import (
	"os"
	"payment/internal/config"
	"payment/internal/handler"
	"payment/internal/infrastructure"
	"payment/internal/middleware"
	model "payment/internal/model"
	"payment/internal/repository"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	zerolog "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/veritrans/go-midtrans"
)

func main() {
	c := config.LoadConfig()
	db := config.InitDB(c)
	// Auto migrate payment-related models. Ensures that the payments and
	// payment_methods tables exist when the service starts.
	{
		type (
			Payment       = model.Payment
			PaymentMethod = model.PaymentMethod
		)
		if err := db.AutoMigrate(&Payment{}, &PaymentMethod{}); err != nil {
			zerolog.Info().Err(err).Msg("failed to auto migrate payment service database")
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

	midtransClient := midtrans.NewClient()
	midtransClient.ServerKey = c.Midtrans.ServerKey
	midtransClient.ClientKey = c.Midtrans.ClientKey
	midtransClient.APIEnvType = midtrans.Sandbox

	restryClient := resty.New()
	restryClient.SetDebug(cast.ToBool(os.Getenv("DEBUG")))

	serviceInfra := infrastructure.NewPaymentService(&midtransClient, c.Midtrans)
	serviceBooking := infrastructure.NewBooking(restryClient, c.ServiceBooking)
	// Initialize user client for calling user service's internal activity
	serviceUser := infrastructure.NewUser(restryClient, c.ServiceUser)
	// Initialize the payment service with user client so activities can be logged
	service := service.NewService(serviceInfra, repo, serviceBooking, serviceUser)

	// Initialize handlers
	paymentHandler := handler.NewHandler(service)
	paymentMethodCRUDHandler := handler.NewPaymentMethodCRUDHandler(service)

	// Payment routes
	api := r.Group("/api/v1")
	if !c.IsNFT {
		api.Use(middleware.AuthMiddleware(&c.JWT))
	}
	{
		api.POST("/payments", paymentHandler.CreatePayment)
		api.GET("/payments", paymentHandler.GetPayments)
		api.GET("/payment/:id", paymentHandler.GetPaymentById)
	}

	// Payment callback
	callback := r.Group("/api/v1")
	callback.POST("/payments/callback", paymentHandler.HandleWebhook)

	// Payment method CRUD routes
	crudMethods := r.Group("/api/v1/admin/payment-methods")
	crudMethods.Use(middleware.AuthMiddleware(&c.JWT))
	{
		crudMethods.GET("/", paymentMethodCRUDHandler.GetPaymentMethods)
		crudMethods.POST("/", paymentMethodCRUDHandler.CreatePaymentMethod)
		crudMethods.GET("/:id", paymentMethodCRUDHandler.GetPaymentMethod)
		crudMethods.PUT("/:id", paymentMethodCRUDHandler.UpdatePaymentMethod)
		crudMethods.DELETE("/:id", paymentMethodCRUDHandler.DeletePaymentMethod)
	}

	// Basic payment method routes
	v1 := r.Group("/api/v1/payment-methods")
	{
		v1.GET("/", paymentHandler.GetAll)
		v1.POST("/status", paymentHandler.SetActive)
	}

	zerolog.Info().Msg("Server is running on port " + c.AppPort)
	r.Run(":" + c.AppPort)
}
