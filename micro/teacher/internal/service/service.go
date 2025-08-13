package service

import (
	"errors"
	"teacher/internal/models"
	"teacher/internal/pkg"
	"teacher/internal/repository"
	"time"
)

type Service struct {
	teacherRepo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		teacherRepo: repo,
	}
}

func (s *Service) GetTeachers(paginate *pkg.Paginate) (pkg.ResponsePaginate, error) {
	response, err := s.teacherRepo.GetTeachers(paginate)
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	// Transform Teacher models to TeacherResponse format
	teachers, ok := response.Data.([]models.Teacher)
	if !ok {
		return response, nil
	}

	var teacherResponses []models.TeacherResponse
	for _, teacher := range teachers {
		teacherResponse := models.TeacherResponse{
			ID:             teacher.ID,
			Name:           teacher.Name,
			Bio:            teacher.Bio,
			LanguageLevel:  teacher.LanguageLevel,
			PricePerHour:   teacher.PricePerHour,
			AvailableStart: teacher.AvailableStartTime,
			AvailableEnd:   teacher.AvailableEndTime,
			ProfileImage:   teacher.ProfileImage,
			CreatedAt:      teacher.CreatedAt.String(),
			UpdatedAt:      teacher.UpdatedAt.String(),
			Schedules:      teacher.Schedules,
		}
		teacherResponses = append(teacherResponses, teacherResponse)
	}

	response.Data = teacherResponses
	return response, nil
}

func (s *Service) CreateTeacher(teacher models.TeacherRequest) error {

	availableStart, err := pkg.ParseTime(teacher.AvailableStart)
	if err != nil {
		return err
	}
	availableEnd, err := pkg.ParseTime(teacher.AvailableEnd)
	if err != nil {
		return err
	}

	teacherReq := models.Teacher{
		Bio:                teacher.Bio,
		Name:               teacher.Name,
		LanguageLevel:      teacher.LanguageLevel,
		PricePerHour:       teacher.PricePerHour,
		AvailableStartTime: pkg.NormalizeTime(availableStart),
		AvailableEndTime:   pkg.NormalizeTime(availableEnd),
		// AvailableDays:      teacher.AvailableDays,
		ProfileImage: teacher.ProfileImage,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.teacherRepo.CreateTeacher(&teacherReq)

}

func (s *Service) GetTeacherByID(id uint) (*models.TeacherResponse, error) {
	var teacher models.TeacherResponse
	data, err := s.teacherRepo.GetTeacherByID(id)
	if err != nil {
		return nil, err
	}

	teacher.ID = data.ID
	teacher.Bio = data.Bio
	teacher.Name = data.Name
	teacher.LanguageLevel = data.LanguageLevel
	teacher.PricePerHour = data.PricePerHour
	teacher.AvailableStart = data.AvailableStartTime
	teacher.AvailableEnd = data.AvailableEndTime
	// teacher.AvailableDays = data.AvailableDays
	teacher.ProfileImage = data.ProfileImage
	teacher.CreatedAt = data.CreatedAt.String()
	teacher.UpdatedAt = data.UpdatedAt.String()

	return &teacher, nil
}

func (s *Service) UpdateTeacher(teacher models.TeacherRequest) error {

	availableStart, err := pkg.ParseTime(teacher.AvailableStart)
	if err != nil {
		return err
	}
	availableEnd, err := pkg.ParseTime(teacher.AvailableEnd)
	if err != nil {
		return err
	}

	dataTeacher, err := s.teacherRepo.GetTeacherByID(teacher.ID)
	if err != nil {
		return err
	}

	if dataTeacher == nil {
		return errors.New("teacher not found")
	}

	return s.teacherRepo.UpdateTeacher(&models.Teacher{
		ID:                 teacher.ID,
		Name:               teacher.Name,
		Bio:                teacher.Bio,
		LanguageLevel:      teacher.LanguageLevel,
		PricePerHour:       teacher.PricePerHour,
		AvailableStartTime: pkg.NormalizeTime(availableStart),
		AvailableEndTime:   pkg.NormalizeTime(availableEnd),
		// AvailableDays:      teacher.AvailableDays,
		ProfileImage: teacher.ProfileImage,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
}

func (s *Service) DeleteTeacher(teacher models.TeacherRequest) error {
	return s.teacherRepo.DeleteTeacher(&models.Teacher{
		ID: teacher.ID,
	})
}

func (s *Service) CountTeachers() (int64, error) {
	return s.teacherRepo.CountTeachers()
}


func (s *Service) GetTeacherByUserID(userID uint) (*models.Teacher, error) {
	return s.teacherRepo.GetTeacherByUserID(userID)
}
