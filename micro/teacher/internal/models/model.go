package models

import "time"

type Teacher struct {
	UserID             uint       `gorm:"index" json:"user_id"`
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Bio                string     `json:"bio"`
	Name               string     `json:"name"`
	LanguageLevel      string     `gorm:"size:100" json:"language_level"`
	PricePerHour       float64    `json:"price_per_hour"`
	AvailableStartTime string     `gorm:"type:VARCHAR(8)" json:"available_start_time"` // format: "HH:mm"
	AvailableEndTime   string     `gorm:"type:VARCHAR(8)" json:"available_end_time"`
	ProfileImage       string     `json:"profile_image"`
	Schedules          []Schedule `gorm:"foreignKey:TeacherID" json:"schedules,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

type Schedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeacherID uint      `gorm:"index;type:int;not null" json:"teacher_id"`
	Teacher   *Teacher  `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"teacher"`
	Date      time.Time `gorm:"type:date" json:"date"`
	StartTime string    `gorm:"type:VARCHAR(8)" json:"start_time"`
	EndTime   string    `gorm:"type:VARCHAR(8)" json:"end_time"`
	Status    string    `gorm:"type:enum('available','booked','cancelled');default:'available'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TeacherResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Bio            string     `json:"bio"`
	LanguageLevel  string     `json:"language_level"`
	PricePerHour   float64    `json:"price_per_hour"`
	AvailableStart string     `json:"available_start_time"`
	AvailableEnd   string     `json:"available_end_time"`
	ProfileImage   string     `json:"profile_image"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
	Schedules      []Schedule `json:"schedules"`
}

type TeacherRequest struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Bio            string  `json:"bio"`
	LanguageLevel  string  `json:"language_level"`
	PricePerHour   float64 `json:"price_per_hour"`
	AvailableStart string  `json:"available_start_time"`
	AvailableEnd   string  `json:"available_end_time"`
	// AvailableDays  string  `json:"available_days"`
	ProfileImage string `json:"profile_image"`
}
