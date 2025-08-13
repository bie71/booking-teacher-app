package schedule

import "time"

type ScheduleFilterRequest struct {
	TeacherID   string   `json:"teacher_id"`
	ScheduleIDs []string `json:"schedule_ids"`
}

type ScheduleFilterResponse struct {
	ValidScheduleIDs []string `json:"valid_schedule_ids"`
}

type Teacher struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	Bio                string      `json:"bio"`
	Name               string      `json:"name"`
	LanguageLevel      string      `gorm:"size:10" json:"language_level"`
	PricePerHour       float64     `json:"price_per_hour"`
	AvailableStartTime string      `gorm:"type:VARCHAR(8)" json:"available_start_time"` // format: "HH:mm"
	AvailableEndTime   string      `gorm:"type:VARCHAR(8)" json:"available_end_time"`
	ProfileImage       string      `json:"profile_image"`
	Schedules          *[]Schedule `gorm:"foreignKey:TeacherID" json:"schedules,omitempty"`
}

type Schedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeacherID uint      `gorm:"index,not null" json:"teacher_id"`
	Teacher   *Teacher  `gorm:"foreignKey:TeacherID" json:"teacher"`
	Date      time.Time `gorm:"type:date" json:"date"`
	StartTime string    `gorm:"type:VARCHAR(8)" json:"start_time"`
	EndTime   string    `gorm:"type:VARCHAR(8)" json:"end_time"`
	Status    string    `gorm:"type:enum('available','booked','cancelled');default:'available'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
