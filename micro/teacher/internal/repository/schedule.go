package repository

import (
	"log"
	"math"
	"teacher/internal/models"
	"teacher/internal/pkg"
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *Schedule {
	return &Schedule{
		DB: db,
	}
}

func (s *Schedule) GetTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := s.DB.First(&teacher, id).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (s *Schedule) HasScheduleConflict(teacherID uint, date time.Time, start, end string) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Schedule{}).
		Where("teacher_id = ? AND date = ? AND start_time < ? AND end_time > ? AND status = ?", teacherID, date, end, start, "booked").
		Count(&count).Error
	return count > 0, err
}

func (s *Schedule) CreateSchedule(schedule *models.Schedule) error {
	return s.DB.Create(schedule).Error
}

func (s *Schedule) GetAvailableSchedules(teacherID uint, paginate *pkg.Paginate) (pkg.ResponsePaginate, error) {
	var schedules []models.Schedule

	if paginate.Page < 1 {
		paginate.Page = 1
	}

	if paginate.Limit <= 0 || paginate.Limit > 100 {
		paginate.Limit = 10
	}

	var count int64
	err := s.DB.Model(&models.Schedule{}).
		Where("teacher_id = ? AND status = ? AND date >= ?", teacherID, "available", time.Now().Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	offset := (paginate.Page - 1) * paginate.Limit
	err = s.DB.Offset(offset).Limit(paginate.Limit).
		Where("teacher_id = ? AND status = ? AND date >= ?", teacherID, "available", time.Now().Format("2006-01-02")).
		Find(&schedules).Error

	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: schedules,
		Pagination: pkg.PaginationPage{
			CurrentPage: paginate.Page,
			TotalPage:   int(math.Ceil(float64(count) / float64(paginate.Limit))),
			TotalData:   int(count),
			Limit:       paginate.Limit,
		},
	}, nil
}

func (s *Schedule) CancelSchedule(id uint) error {
	return s.DB.Model(&models.Schedule{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (s *Schedule) GetSchedulesById(id uint) (models.Schedule, error) {
	var schedule models.Schedule
	if err := s.DB.Preload("Teacher").First(&schedule, id).Error; err != nil {
		return models.Schedule{}, err
	}
	return schedule, nil
}

func (s *Schedule) DeleteSchedule(id uint) error {
	return s.DB.Delete(&models.Schedule{}, id).Error
}

func (s *Schedule) UpdateScheduleStatus(id uint, status string) error {
	log.Println(status, id)
	return s.DB.Model(&models.Schedule{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (s *Schedule) GetBatchScheduleDetail(ids []uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := s.DB.Preload("Teacher").
		Where("id IN ?", ids).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *Schedule) FindByTeacherAndIDs(teacherID uint, ids []uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := s.DB.
		Where("teacher_id = ? AND id IN ?", teacherID, ids).
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *Schedule) GetSchedules(paginate pkg.Paginate) (pkg.ResponsePaginate, error) {
	var schedules []models.Schedule
	var count int64

	if paginate.Page < 1 {
		paginate.Page = 1
	}

	if paginate.Limit <= 0 || paginate.Limit > 100 {
		paginate.Limit = 10
	}

	err := s.DB.Model(&models.Schedule{}).Count(&count).Error
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	offset := (paginate.Page - 1) * paginate.Limit
	err = s.DB.Preload("Teacher").Offset(offset).Limit(paginate.Limit).Find(&schedules).Error
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: schedules,
		Pagination: pkg.PaginationPage{
			CurrentPage: paginate.Page,
			TotalPage:   int(math.Ceil(float64(count) / float64(paginate.Limit))),
			TotalData:   int(count),
			Limit:       paginate.Limit,
		},
	}, nil
}

func (s *Schedule) UpdateSchedule(id uint, schedule *models.Schedule) error {
	return s.DB.Model(&models.Schedule{}).
		Where("id = ?", id).
		Updates(schedule).Error
}
