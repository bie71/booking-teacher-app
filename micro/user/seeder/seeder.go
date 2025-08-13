package seeder

import (
	"auth/internal/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {

	user := models.User{}

	db.Find(&models.User{}).Where("email = ?", "admin@example.com").First(&user)
	if user.ID > 0 {
		log.Println("Admin user already exists")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:         "Admin",
		Email:        "admin@example.com",
		PasswordHash: string(password),
		Role:         "admin",
	}

	db.Save(&admin)
}
