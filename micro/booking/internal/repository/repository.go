package repository

import (
	"booking/internal/model"
	"booking/internal/pkg"
	"math"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{Db: db} }

func (r *Repository) CreateBooking(b *model.Booking) error {
	return r.Db.Create(b).Error
}

func (r *Repository) UpdateBooking(b *model.Booking) error {
	return r.Db.Save(b).Error
}

func (r *Repository) DeleteBooking(id uint) error {
	return r.Db.Delete(&model.Booking{}, id).Error
}

func (r *Repository) GetBooking(id uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.Db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

// Basic paginated list
func (r *Repository) GetBookings(pagination pkg.Paginate) (pkg.ResponsePaginate, error) {
	var bookings []model.Booking
	var total int64
	q := r.Db.Model(&model.Booking{})

	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 10
	}
	offset := (pagination.Page - 1) * pagination.Limit

	if err := q.Count(&total).Error; err != nil {
		return pkg.ResponsePaginate{}, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(pagination.Limit).Find(&bookings).Error; err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: bookings,
		Pagination: pkg.PaginationPage{
			CurrentPage: pagination.Page,
			TotalPage:   int(math.Ceil(float64(total) / float64(pagination.Limit))),
			TotalData:   int(total),
			Limit:       pagination.Limit,
		},
	}, nil
}

// DB-level filter by status and date range
func (r *Repository) GetBookingsFiltered(pagination pkg.Paginate, status string, startDateStr, endDateStr string) (pkg.ResponsePaginate, error) {
	var bookings []model.Booking
	var total int64

	q := r.Db.Model(&model.Booking{})

	if status != "" {
		q = q.Where("status = ?", status)
	}
	if startDateStr != "" {
		q = q.Where("created_at >= ?", startDateStr)
	}
	if endDateStr != "" {
		q = q.Where("created_at <= ?", endDateStr)
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 10
	}
	offset := (pagination.Page - 1) * pagination.Limit

	if err := q.Count(&total).Error; err != nil {
		return pkg.ResponsePaginate{}, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(pagination.Limit).Find(&bookings).Error; err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: bookings,
		Pagination: pkg.PaginationPage{
			CurrentPage: pagination.Page,
			TotalPage:   int(math.Ceil(float64(total) / float64(pagination.Limit))),
			TotalData:   int(total),
			Limit:       pagination.Limit,
		},
	}, nil
}

// ChangeStatus updates booking status by id
func (r *Repository) ChangeStatus(id uint, status string) error {
	return r.Db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}

// Get a booking by user and id (ownership check)
func (r *Repository) GetBookingsByUserIDAndId(userID, id uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.Db.Where("user_id = ? AND id = ?", userID, id).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

// Get user bookings with optional admin override, filters, and offset pagination
func (r *Repository) GetBookingsByUserID(userID uint, isAdmin bool, status string, startDate, endDate *time.Time, page, limit int) ([]model.Booking, error) {
	var bookings []model.Booking
	q := r.Db.Model(&model.Booking{})

	if !isAdmin {
		q = q.Where("user_id = ?", userID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if startDate != nil {
		q = q.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		q = q.Where("created_at <= ?", *endDate)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *Repository) CountBookings(userID uint, isAdmin bool, status string, startDate, endDate *time.Time) (int, error) {
	var count int64
	q := r.Db.Model(&model.Booking{})

	if !isAdmin {
		q = q.Where("user_id = ?", userID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if startDate != nil {
		q = q.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		q = q.Where("created_at <= ?", *endDate)
	}

	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// SumTotalPricePaidBookings sums TotalPrice where status='paid'
func (r *Repository) SumTotalPricePaidBookings() (float64, error) {
	var sum float64
	if err := r.Db.Model(&model.Booking{}).Where("status in (?, ?)", "paid", "completed").Select("COALESCE(SUM(total_price),0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
