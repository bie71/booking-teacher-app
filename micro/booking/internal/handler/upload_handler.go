package handler

import (
	"booking/internal/config"
	"booking/internal/infrastructure/supabase"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *supabase.Client
	cfg           *config.Client
}

func NewUploadHandler(uploadService *supabase.Client, cfg *config.Client) *UploadHandler {
	return &UploadHandler{uploadService: uploadService, cfg: cfg}
}
func (h *UploadHandler) UploadHandler(c *gin.Context) {
	// Ambil file dari form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
		return
	}

	safeFilename := sanitizeFilename(fileHeader.Filename)

	// Simpan ke temporary directory OS
	tempPath := filepath.Join(os.TempDir(), safeFilename) // ✅ Simpan file temporer di /tmp/

	if err := c.SaveUploadedFile(fileHeader, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer os.Remove(tempPath)

	objectName := "uploads/" + safeFilename

	err = h.uploadService.UploadToSupabase(objectName, tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	// Public URL jika bucket diatur "public"
	fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		h.cfg.Endpoint,
		h.cfg.BucketName,
		objectName,
	)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_url": fileURL})
}

func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Hapus karakter aneh, hanya izinkan huruf, angka, titik, dash, underscore
	re := regexp.MustCompile(`[^a-z0-9._-]+`)
	name = re.ReplaceAllString(name, "")

	// Batasi panjang nama file kalau perlu
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}
