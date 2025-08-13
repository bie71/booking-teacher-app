package repository

import (
	"errors"
	"math"
	"teacher/internal/models"
	"teacher/internal/pkg"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateTeacher(teacher *models.Teacher) error {
	return r.db.Create(teacher).Error
}

func (r *Repository) GetTeachers(paginate *pkg.Paginate) (pkg.ResponsePaginate, error) {
	var teachers []models.Teacher

	if paginate.Page < 1 {
		paginate.Page = 1
	}

	if paginate.Limit <= 0 || paginate.Limit > 100 {
		paginate.Limit = 10
	}

	var total int64
	r.db.Model(&models.Teacher{}).Count(&total)

	offset := (paginate.Page - 1) * paginate.Limit
	err := r.db.Preload("Schedules").Offset(offset).Limit(paginate.Limit).Find(&teachers).Error

	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: teachers,
		Pagination: pkg.PaginationPage{
			CurrentPage: paginate.Page,
			TotalPage:   int(math.Ceil(float64(total) / float64(paginate.Limit))),
			TotalData:   int(total),
			Limit:       paginate.Limit,
		},
	}, nil
}

func (r *Repository) UpdateTeacher(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *Repository) DeleteTeacher(teacher *models.Teacher) error {
	return r.db.Delete(teacher).Error
}
func (r *Repository) GetTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.First(&teacher, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *Repository) CountTeachers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Teacher{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}


func (r *Repository) GetTeacherByUserID(userID uint) (*models.Teacher, error) {
	var t models.Teacher
	err := r.db.Preload("Schedules").Where("user_id = ?", userID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
