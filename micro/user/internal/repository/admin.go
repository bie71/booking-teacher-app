package repository

import (
	"auth/internal/models"
	"log"

	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		DB: db,
	}
}

func (r *AdminRepository) SaveImageUrl(key string, imageUrl string) error {
	return r.DB.Save(&models.HeroImage{KeyImage: key, ImageUrl: imageUrl}).Error
}
func (r *AdminRepository) GetImageUrl(key string) (string, error) {
	log.Println(key)
	var image models.HeroImage
	if err := r.DB.Where("key_image = ?", key).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return image.ImageUrl, nil
}

func (r *AdminRepository) DeleteImageUrl(key string) error {
	return r.DB.Where("key_image = ?", key).Delete(&models.HeroImage{}).Error
}
