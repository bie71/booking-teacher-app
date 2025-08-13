package handler

import (
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/service"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Handler struct {
	userService *service.UserService
	repoUser    *repository.UserRepository
	// adminHandler *handler.AdminHandler // Removed due to missing import
}

func NewHandler(userService *service.UserService, repoUser *repository.UserRepository) *Handler {
	// uploadDir := filepath.Join(".", "uploads", "hero")
	// adminHandler := handler.NewAdminHandler(uploadDir)

	return &Handler{
		userService: userService,
		repoUser:    repoUser,
		// adminHandler: adminHandler,
	}
}

// LogActivity handles POST requests to record a new activity for the
// authenticated user. The request body should contain an `action`
// describing the event and an optional `description` providing
// additional details. If validation fails, a 400 response is
// returned. If logging fails, a 500 response is returned.
func (h *Handler) LogActivity(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	var req struct {
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := cast.ToUint(userIDStr)
	if err := h.userService.CreateActivityLog(c, userID, req.Action, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "activity logged successfully"})
}

// LogActivityInternal handles service-to-service requests to record an activity log
// for a user.  This endpoint does not require authentication and expects
// a JSON payload with fields "user_id", "action", and "description".  It
// calls the user service's CreateActivityLog method and returns a success
// response when the activity is recorded.  In case of invalid input or
// errors from the service layer, it returns an appropriate HTTP error.
func (h *Handler) LogActivityInternal(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.UserID == 0 || req.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and action are required"})
		return
	}
	if err := h.userService.CreateActivityLog(c, req.UserID, req.Action, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Activity recorded successfully"})
}

// GetRecentActivity handles GET requests to fetch recent activity logs
// for the authenticated user. An optional `limit` query parameter
// controls how many records are returned. If the parameter is not
// provided or invalid, a default of 5 entries is used. The logs are
// returned in descending order by creation time.
func (h *Handler) GetRecentActivity(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	limit := 5
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	userID := cast.ToUint(userIDStr)
	logs, err := h.userService.GetRecentActivity(c, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// ToggleFavorite handles POST requests to add or remove a teacher from
// the authenticated user's list of favorites. The request body should
// include a `teacher_id` field (integer) and an optional `favorite`
// boolean indicating whether to add (true) or remove (false) the
// favorite. If `favorite` is omitted or true, the teacher will be
// added to the user's favorites. If `favorite` is false, the teacher
// will be removed. A 200 response is returned on success.
func (h *Handler) ToggleFavorite(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	var req struct {
		TeacherID int   `json:"teacher_id" binding:"required"`
		Favorite  *bool `json:"favorite"` // optional: true to add, false to remove
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := cast.ToInt(userIDStr)
	// Determine whether to add or remove favorite
	shouldAdd := true
	if req.Favorite != nil {
		shouldAdd = *req.Favorite
	}
	var err error
	if shouldAdd {
		err = h.userService.AddFavoriteTeacher(c, userID, req.TeacherID)
	} else {
		err = h.userService.RemoveFavoriteTeacher(c, userID, req.TeacherID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "favorite updated successfully"})
}

// GetFavorites handles GET requests to retrieve a list of teacher IDs
// that the authenticated user has marked as favorites. Returns an
// array of integers under the `data` key. If there are no
// favorites, an empty array is returned. A 200 response is
// always returned unless authentication fails.
func (h *Handler) GetFavorites(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	userID := cast.ToInt(userIDStr)
	teacherIDs, err := h.userService.GetFavoriteTeachers(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teacherIDs})
}
func (h *Handler) RegisterUser(c *gin.Context) {
	var user models.RegisterRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.userService.Register(c, user)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) LoginUser(c *gin.Context) {
	var user models.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.userService.Login(c, user)
	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "user not found" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) VerifyResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := h.repoUser.FindUserByResetToken(token)
	if err != nil || user.ResetExpiration.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token valid",
		"user_id": user.ID, // bisa disimpan sementara di frontend
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var input models.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ResetPassword(c, input.Token, input.NewPassword)

	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// ForgotPassword handles the password reset request. It accepts an email
// address in the request body and, if a user with that email exists,
// triggers the generation of a reset token and sends a password reset
// email. The response is intentionally generic to avoid revealing
// whether the email exists in the system.
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Call the service to process the forgot password request. Errors are
	// ignored in the response to prevent account enumeration.
	_ = h.userService.ProcessForgotPassword(c, req.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists in our system, a password reset link has been sent",
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	log.Println(userID)
	user, err := h.userService.GetUserById(c, cast.ToUint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c, cast.ToUint(userID), req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "email already taken" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ChangePassword(c, cast.ToUint(userID), req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "current password is incorrect" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.userService.GetDashboardStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.Query("user_id")
	log.Println(userID)
	user, err := h.userService.GetUserById(c, cast.ToUint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
