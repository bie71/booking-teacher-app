package repository

import (
	"auth/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateResetToken(user *models.User, token string, expiration time.Time) error {
	user.ResetToken = token
	user.ResetExpiration = &expiration

	return r.db.Save(user).Error
}

func (r *UserRepository) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) DeleteUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Unscoped().Delete(user).Error
}

func (r *UserRepository) FindUserByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserPasswordAndClearToken(userID uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash":    hashedPassword,
		"reset_token":      nil,
		"reset_expiration": nil,
	}).Error
}

func (r *UserRepository) GetUsersWithPagination(ctx context.Context, page, limit int, search, role string) ([]models.UserListResponse, int64, error) {
	var users []models.UserListResponse
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	// Apply search filter
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply role filter
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.
		Select("id, name, email, role, profile_image, created_at, updated_at").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, userID).Error
}

func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role != ?", "admin").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
