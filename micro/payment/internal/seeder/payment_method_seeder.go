package seeder

import (
	"log"
	"payment/internal/model"

	"gorm.io/gorm"
)

func SeedPaymentMethods(db *gorm.DB) error {
	paymentMethods := []model.PaymentMethod{
		{
			Code:     "credit_card",
			Name:     "Credit Card",
			IsActive: true,
		},
		{
			Code:     "bank_transfer",
			Name:     "Bank Transfer",
			IsActive: true,
		},
		{
			Code:     "gopay",
			Name:     "GoPay",
			IsActive: true,
		},
		{
			Code:     "ovo",
			Name:     "OVO",
			IsActive: true,
		},
		{
			Code:     "dana",
			Name:     "DANA",
			IsActive: true,
		},
		{
			Code:     "shopeepay",
			Name:     "ShopeePay",
			IsActive: true,
		},
		{
			Code:     "linkaja",
			Name:     "LinkAja",
			IsActive: true,
		},
		{
			Code:     "bca_va",
			Name:     "BCA Virtual Account",
			IsActive: true,
		},
		{
			Code:     "bni_va",
			Name:     "BNI Virtual Account",
			IsActive: true,
		},
		{
			Code:     "bri_va",
			Name:     "BRI Virtual Account",
			IsActive: true,
		},
		{
			Code:     "mandiri_va",
			Name:     "Mandiri Virtual Account",
			IsActive: true,
		},
		{
			Code:     "permata_va",
			Name:     "Permata Virtual Account",
			IsActive: true,
		},
		{
			Code:     "indomaret",
			Name:     "Indomaret",
			IsActive: true,
		},
		{
			Code:     "alfamart",
			Name:     "Alfamart",
			IsActive: true,
		},
	}

	for _, method := range paymentMethods {
		var existingMethod model.PaymentMethod
		result := db.Where("code = ?", method.Code).First(&existingMethod)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new payment method
			if err := db.Create(&method).Error; err != nil {
				log.Printf("Error creating payment method %s: %v", method.Code, err)
				return err
			}
			log.Printf("Created payment method: %s", method.Name)
		} else if result.Error != nil {
			log.Printf("Error checking payment method %s: %v", method.Code, result.Error)
			return result.Error
		} else {
			// Update existing payment method
			existingMethod.Name = method.Name
			existingMethod.IsActive = method.IsActive
			if err := db.Save(&existingMethod).Error; err != nil {
				log.Printf("Error updating payment method %s: %v", method.Code, err)
				return err
			}
			log.Printf("Updated payment method: %s", method.Name)
		}
	}

	log.Println("Payment methods seeding completed successfully")
	return nil
}
