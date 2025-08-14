package models

import "time"

// ActivityLog represents a record of a user action in the system.
// Each entry is associated with a specific user and describes
// what action occurred along with an optional description. Logs
// are ordered by their creation time, allowing recent activity
// to be queried efficiently.
//
// The `json` tags ensure that when this struct is serialized to
// JSON (for example, in API responses), the field names are
// presented in a camelCase format. The GORM tags define the
// database column types and constraints.
type ActivityLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;type:int unsigned;not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Action      string    `gorm:"size:255;not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type FavoriteTeacher struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	TeacherID uint      `gorm:"index;not null" json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
}
