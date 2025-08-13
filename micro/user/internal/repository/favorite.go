package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// FavoriteRepository manages favorite teachers data in the database.
// It allows users to add and remove teachers from their list of favorites
// and retrieve the list of favorite teacher IDs for a given user.
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository creates a new instance of FavoriteRepository.
func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

// AddFavorite associates a teacher with a user as a favorite. If the
// association already exists, this method returns nil without
// creating a duplicate entry.
func (r *FavoriteRepository) AddFavorite(ctx context.Context, userID int, teacherID int) error {
	if r.db == nil {
		return errors.New("database not initialized")
	}
	// Use INSERT ... ON CONFLICT DO NOTHING to avoid duplicate entries
	return r.db.WithContext(ctx).Exec(
		"INSERT IGNORE INTO favorite_teachers (user_id, teacher_id) VALUES (?, ?);",
		userID, teacherID,
	).Error
}

// RemoveFavorite removes a teacher from a user's list of favorites. If
// the association does not exist, no error is returned.
func (r *FavoriteRepository) RemoveFavorite(ctx context.Context, userID int, teacherID int) error {
	if r.db == nil {
		return errors.New("database not initialized")
	}
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM favorite_teachers WHERE user_id = ? AND teacher_id = ?",
		userID, teacherID,
	).Error
}

// GetFavoritesByUser retrieves all teacher IDs favorited by a given user.
func (r *FavoriteRepository) GetFavoritesByUser(ctx context.Context, userID int) ([]int, error) {
	if r.db == nil {
		return nil, errors.New("database not initialized")
	}
	var teacherIDs []int
	rows, err := r.db.WithContext(ctx).Raw(
		"SELECT teacher_id FROM favorite_teachers WHERE user_id = ?",
		userID,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		teacherIDs = append(teacherIDs, id)
	}
	return teacherIDs, nil
}
