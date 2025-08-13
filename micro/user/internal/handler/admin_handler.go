package handler

import (
	"auth/internal/models"
	"auth/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-specific endpoints
type AdminHandler struct {
	imageRepo *repository.AdminRepository
}

func NewAdminHandler(imageRepo *repository.AdminRepository) *AdminHandler {
	return &AdminHandler{
		imageRepo: imageRepo,
	}
}

func (h *AdminHandler) GetHeroImage(c *gin.Context) {
	key := c.Query("key")
	imageURL, err := h.imageRepo.GetImageUrl(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Hero image URL",
		"image_url": imageURL,
	})
}

func (h *AdminHandler) SaveHeroImage(c *gin.Context) {
	var req models.HeroImage

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.imageRepo.SaveImageUrl(req.KeyImage, req.ImageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hero image saved successfully"})
}

func (h *AdminHandler) DeleteHeroImage(c *gin.Context) {
	key := c.Query("key")
	err := h.imageRepo.DeleteImageUrl(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hero image deleted successfully"})
}
