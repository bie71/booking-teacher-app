package repository

import (
    "auth/internal/models"
    "gorm.io/gorm"
)

// ActivityLogRepository provides methods to persist and retrieve
// ActivityLog records from the database. It encapsulates the GORM
// operations to make unit testing and future modifications easier.
type ActivityLogRepository struct {
    db *gorm.DB
}

// NewActivityLogRepository instantiates a new ActivityLogRepository
// with the provided GORM DB connection.
func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
    return &ActivityLogRepository{db: db}
}

// Create inserts a new ActivityLog record into the database. It
// returns an error if the insert fails.
func (r *ActivityLogRepository) Create(log *models.ActivityLog) error {
    return r.db.Create(log).Error
}

// GetRecentByUser fetches the most recent activity logs for a given
// user. The results are ordered by CreatedAt descending and limited
// by the provided limit. If the query fails, it returns an error.
func (r *ActivityLogRepository) GetRecentByUser(userID uint, limit int) ([]models.ActivityLog, error) {
    var logs []models.ActivityLog
    if err := r.db.Where("user_id = ?", userID).
        Order("created_at desc").
        Limit(limit).
        Find(&logs).Error; err != nil {
        return nil, err
    }
    return logs, nil
}