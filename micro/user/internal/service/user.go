package service

import (
	"auth/internal/infrastructure/statistic"
	"auth/internal/models"
	"auth/internal/repository"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	failedToRegisterUser = "failed to register user"
	userNotFound         = "user not found"
)

type UserService struct {
	repoUser         *repository.UserRepository
	emailService     *EmailService
	expiration       int
	jwtConfig        *JWTConfig
	serviceStatistic *statistic.Client
	activityRepo     *repository.ActivityLogRepository
	favoriteRepo     *repository.FavoriteRepository
}

func NewUserService(
	repoUser *repository.UserRepository,
	emailService *EmailService,
	expiration int,
	jwtConfig *JWTConfig,
	serviceStatistic *statistic.Client,
	activityRepo *repository.ActivityLogRepository,
	favoriteRepo *repository.FavoriteRepository,
) *UserService {
	return &UserService{
		repoUser:         repoUser,
		emailService:     emailService,
		expiration:       expiration,
		jwtConfig:        jwtConfig,
		serviceStatistic: serviceStatistic,
		activityRepo:     activityRepo,
		favoriteRepo:     favoriteRepo,
	}
}

func (s *UserService) Register(ctx context.Context, req models.RegisterRequest) (*models.UserResponse, error) {

	existingUser, _ := s.repoUser.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New(failedToRegisterUser)
	}

	role := "user"
	if req.Role != "" {
		switch strings.ToLower(req.Role) {
		case "admin":
			role = "admin"
		case "teacher":
			role = "teacher"
		case "user":
			role = "user"
		default:
			return nil, errors.New("invalid role")
		}
	}

	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.repoUser.CreateUser(ctx, newUser); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, errors.New(failedToRegisterUser)
	}

	return &models.UserResponse{
		Id:           newUser.ID,
		Name:         newUser.Name,
		Email:        newUser.Email,
		Role:         newUser.Role,
		ProfileImage: newUser.ProfileImage,
		CreatedAt:    newUser.CreatedAt,
		UpdatedAt:    newUser.UpdatedAt,
	}, nil

}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {

	user, err := s.repoUser.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by email")
		return nil, errors.New(userNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Error().Err(err).Msg("Failed to compare password")
		return nil, errors.New("invalid password")
	}

	token, err := s.jwtConfig.GenerateToken(strconv.Itoa(int(user.ID)), user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		Token:  token,
		UserId: strconv.Itoa(int(user.ID)),
		User: models.UserResponse{
			Id:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Role:         user.Role,
			ProfileImage: user.ProfileImage,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.repoUser.GetUserById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by id")
		return nil, errors.New(userNotFound)
	}

	return &models.UserResponse{
		Id:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *UserService) ProcessForgotPassword(ctx context.Context, email string) error {
	user, err := s.repoUser.GetByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by email")
		return errors.New(userNotFound)
	}

	token := uuid.NewString()
	// The expiration uses the JWT token duration configured for the service.
	expiration := time.Now().Add(time.Duration(s.expiration) * time.Hour)

	err = s.repoUser.UpdateResetToken(user, token, expiration)
	if err != nil {
		return err
	}

	// Build a reset link based on the FRONTEND_URL environment variable. If
	// FRONTEND_URL is not set (e.g. during local development), default to
	// http://localhost:3000. The reset link will be used in the email
	// template and should route to the recovery page in the frontend.
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	// The path can be adjusted to match the actual route in the Vue app.
	resetLink := fmt.Sprintf("%s/recover-password?token=%s", strings.TrimRight(frontendURL, "/"), token)

	// Send the password reset email asynchronously. The email template
	// expects a full URL for ResetLink, not just a token.
	go s.emailService.SendPasswordResetEmail(email, resetLink)

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, token string, newPassword string) error {

	user, err := s.repoUser.FindUserByResetToken(token)
	if err != nil || user.ResetExpiration.Before(time.Now()) {
		return errors.New("invalid or expired token")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	err = s.repoUser.UpdateUserPasswordAndClearToken(user.ID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	// Check if user exists
	user, err := s.repoUser.GetUserById(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by id")
		return nil, errors.New(userNotFound)
	}

	// Check if email is already taken by another user
	if req.Email != user.Email {
		existingUser, _ := s.repoUser.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already taken")
		}
	}

	// Update user fields
	user.Name = req.Name
	user.Email = req.Email
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	if err := s.repoUser.UpdateUser(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return nil, errors.New("failed to update profile")
	}

	return &models.UserResponse{
		Id:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, req models.ChangePasswordRequest) error {
	// Get user
	user, err := s.repoUser.GetUserById(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by id")
		return errors.New(userNotFound)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		log.Error().Err(err).Msg("Current password is incorrect")
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash new password")
		return errors.New("failed to change password")
	}

	// Update password
	if err := s.repoUser.UpdateUserPasswordAndClearToken(user.ID, string(hashedPassword)); err != nil {
		log.Error().Err(err).Msg("Failed to update password")
		return errors.New("failed to change password")
	}

	return nil
}

func (s *UserService) GetDashboardStats(ctx *gin.Context) (*models.DashboardStatsResponse, error) {
	usersTotal, err := s.repoUser.CountUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard stats")
		return nil, errors.New("failed to get dashboard stats")
	}

	totalBooking, err := s.serviceStatistic.GetTotalBookings(ctx)
	log.Info().Interface("totalBooking", totalBooking)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard stats")
		return nil, errors.New("failed to get dashboard stats")
	}

	totalTeacher, err := s.serviceStatistic.GetTotalTeachers(ctx)
	log.Info().Interface("totalTeacher", totalTeacher)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard stats")
		return nil, errors.New("failed to get dashboard stats")
	}

	stats := &models.DashboardStatsResponse{
		TotalUsers:    usersTotal,
		TotalBookings: totalBooking.TotalBookings,
		TotalRevenue:  int64(totalBooking.TotalRevenue),
		TotalTeachers: totalTeacher.TotalTeachers,
	}

	return stats, nil
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	existingUser, _ := s.repoUser.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New("failed to create user")
	}

	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	if err := s.repoUser.CreateUser(ctx, newUser); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, errors.New("failed to create user")
	}

	return &models.UserResponse{
		Id:           newUser.ID,
		Name:         newUser.Name,
		Email:        newUser.Email,
		Role:         newUser.Role,
		ProfileImage: newUser.ProfileImage,
		CreatedAt:    newUser.CreatedAt,
		UpdatedAt:    newUser.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uint, req models.UpdateUserAdminRequest) (*models.UserResponse, error) {
	user, err := s.repoUser.GetUserById(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by id")
		return nil, errors.New("user not found")
	}

	// Check if email is already taken by another user
	if req.Email != user.Email {
		existingUser, _ := s.repoUser.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already taken")
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	if err := s.repoUser.UpdateUser(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return nil, errors.New("failed to update user")
	}

	return &models.UserResponse{
		Id:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.repoUser.GetUserById(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by id")
		return errors.New("user not found")
	}

	if err := s.repoUser.DeleteUser(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return errors.New("failed to delete user")
	}

	return nil
}

// CreateActivityLog records a new activity entry for a user. The action
// should be a short string describing what the user did (e.g.,
// "booking_created", "profile_updated"), and the description can
// provide additional context. This method delegates to the
// ActivityLogRepository to persist the log. If no activityRepo is
// configured, it returns nil without error to avoid breaking
// functionality when logging is optional.
func (s *UserService) CreateActivityLog(ctx context.Context, userID uint, action, description string) error {
	if s.activityRepo == nil {
		return nil
	}
	logEntry := &models.ActivityLog{
		UserID:      userID,
		Action:      action,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.activityRepo.Create(logEntry)
}

// GetRecentActivity retrieves the most recent activity logs for a
// specific user. The limit parameter controls how many entries
// are returned. If the ActivityLogRepository is not configured,
// it returns an empty slice without error.
func (s *UserService) GetRecentActivity(ctx context.Context, userID uint, limit int) ([]models.ActivityLog, error) {
	if s.activityRepo == nil {
		return []models.ActivityLog{}, nil
	}
	if limit <= 0 {
		limit = 5
	}
	return s.activityRepo.GetRecentByUser(userID, limit)
}

// AddFavoriteTeacher records a teacher as a favorite for a given user.
// If the favorite repository is not configured or the operation fails,
// an error is returned. Duplicate favorites are ignored by the
// repository (no error returned).
func (s *UserService) AddFavoriteTeacher(ctx context.Context, userID int, teacherID int) error {
	if s.favoriteRepo == nil {
		return errors.New("favorite repository not configured")
	}
	err := s.favoriteRepo.AddFavorite(ctx, userID, teacherID)
	if err != nil {
		return err
	}

	go s.CreateActivityLog(ctx, uint(userID), "favorite_teacher_added", fmt.Sprintf("Added teacher %d as a favorite", teacherID))

	return nil

}

// RemoveFavoriteTeacher removes a teacher from a user's list of favorites.
// If the favorite repository is not configured, an error is returned.
func (s *UserService) RemoveFavoriteTeacher(ctx context.Context, userID int, teacherID int) error {
	if s.favoriteRepo == nil {
		return errors.New("favorite repository not configured")
	}
	err := s.favoriteRepo.RemoveFavorite(ctx, userID, teacherID)
	if err != nil {
		return err
	}

	go s.CreateActivityLog(ctx, uint(userID), "favorite_teacher_removed", fmt.Sprintf("Removed teacher %d from favorites", teacherID))

	return nil

}

// GetFavoriteTeachers returns a slice of teacher IDs that the user has
// marked as favorites. If the favorite repository is not configured,
// it returns an empty slice without error.
func (s *UserService) GetFavoriteTeachers(ctx context.Context, userID int) ([]int, error) {
	if s.favoriteRepo == nil {
		return []int{}, nil
	}
	return s.favoriteRepo.GetFavoritesByUser(ctx, userID)
}
