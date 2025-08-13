package models

import "time"

type User struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"size:100;not null"`
	Email           string `gorm:"size:100;unique;not null"`
	PasswordHash    string `gorm:"type:text;not null"`
	Role            string `gorm:"type:enum('user','admin','teacher');not null"`
	ProfileImage    string
	ResetToken      string
	ResetExpiration *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserResponse struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	Email        string `json:"email" binding:"required,email"`
	ProfileImage string `json:"profile_image"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type DashboardStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	TotalTeachers int64 `json:"total_teachers"`
	TotalBookings int64 `json:"total_bookings"`
	TotalRevenue  int64 `json:"total_revenue"`
}

// Admin user management models
type CreateUserRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Role         string `json:"role" binding:"required,oneof=user admin teacher"`
	ProfileImage string `json:"profile_image"`
}

type UpdateUserAdminRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	Email        string `json:"email" binding:"required,email"`
	Role         string `json:"role" binding:"required,oneof=user admin teacher"`
	ProfileImage string `json:"profile_image"`
}

type UserListResponse struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
