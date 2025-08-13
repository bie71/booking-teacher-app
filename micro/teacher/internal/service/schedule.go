package service

import (
	"errors"
	"log"
	"teacher/internal/models"
	"teacher/internal/pkg"
	"teacher/internal/repository"

	"github.com/spf13/cast"
)

type ScheduleService struct {
	scheduleRepo *repository.Schedule
}

func NewScheduleService(scheduleRepo *repository.Schedule) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *ScheduleService) BookScheduleService(input models.Schedule) (*models.Schedule, error) {
	// 1. Validasi teacher exist
	teacher, err := s.scheduleRepo.GetTeacherByID(input.TeacherID)
	if err != nil {
		return nil, errors.New("teacher not found")
	}

	isValidStart, err := pkg.IsWithinRange(teacher.AvailableStartTime, teacher.AvailableEndTime, input.StartTime)
	if err != nil || !isValidStart {
		return nil, errors.New("start time outside teacher availability")
	}

	isValidEnd, err := pkg.IsWithinRange(teacher.AvailableStartTime, teacher.AvailableEndTime, input.EndTime)
	if err != nil || !isValidEnd {
		return nil, errors.New("end time outside teacher availability")
	}

	// 3. Validasi bentrok jadwal
	conflict, err := s.scheduleRepo.HasScheduleConflict(input.TeacherID, input.Date, input.StartTime, input.EndTime)
	if err != nil {
		return nil, errors.New("failed to check schedule")
	}
	if conflict {
		return nil, errors.New("schedule conflict, please choose another time")
	}

	input.Status = "booked"
	if err := s.scheduleRepo.CreateSchedule(&input); err != nil {
		return nil, errors.New("failed to create schedule")
	}

	return &input, nil
}

func (s *ScheduleService) CancelScheduleService(id uint) error {

	_, err := s.scheduleRepo.GetSchedulesById(id)
	if err != nil {
		return errors.New("schedule not found")
	}

	if err := s.scheduleRepo.CancelSchedule(id); err != nil {
		return errors.New("failed to cancel schedule")
	}
	return nil
}

func (s *ScheduleService) GetAvailableScheduleService(teacherID uint, paginate *pkg.Paginate) (pkg.ResponsePaginate, error) {
	schedules, err := s.scheduleRepo.GetAvailableSchedules(teacherID, paginate)
	if err != nil {
		return pkg.ResponsePaginate{}, errors.New("failed to get available schedules")
	}
	return schedules, nil
}

func (s *ScheduleService) GetScheduleService(id uint) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.GetSchedulesById(id)
	if err != nil {
		return nil, errors.New("failed to get schedule")
	}
	return &schedule, nil
}

func (s *ScheduleService) DeleteScheduleService(id uint) error {
	if err := s.scheduleRepo.DeleteSchedule(id); err != nil {
		return errors.New("failed to delete schedule")
	}
	return nil
}

func (s *ScheduleService) CreateScheduleService(schedule *models.Schedule) error {

	teacher, err := s.scheduleRepo.GetTeacherByID(schedule.TeacherID)

	if err != nil {
		return errors.New("teacher not found")
	}

	isValidStart, err := pkg.IsWithinRange(teacher.AvailableStartTime, teacher.AvailableEndTime, schedule.StartTime)
	if err != nil || !isValidStart {
		return errors.New("start time outside teacher availability")
	}

	isValidEnd, err := pkg.IsWithinRange(teacher.AvailableStartTime, teacher.AvailableEndTime, schedule.EndTime)
	if err != nil || !isValidEnd {
		return errors.New("end time outside teacher availability")
	}

	conflict, err := s.scheduleRepo.HasScheduleConflict(schedule.TeacherID, schedule.Date, schedule.StartTime, schedule.EndTime)
	if err != nil {
		return errors.New("failed to check schedule")
	}
	if conflict {
		return errors.New("schedule conflict, please choose another time")
	}

	if err := s.scheduleRepo.CreateSchedule(schedule); err != nil {
		return errors.New("failed to create schedule")
	}
	return nil
}

func (s *ScheduleService) UpdateScheduleService(id uint, status string) error {

	_, err := s.scheduleRepo.GetSchedulesById(id)
	if err != nil {
		return errors.New("schedule not found")
	}

	if err := s.scheduleRepo.UpdateScheduleStatus(id, status); err != nil {
		return errors.New("failed to update schedule")
	}
	return nil
}

func (s *ScheduleService) GetBatchScheduleDetailService(ids []uint) ([]models.ScheduleResponse, error) {
	schedules, err := s.scheduleRepo.GetBatchScheduleDetail(ids)
	if err != nil {
		return nil, errors.New("failed to get batch schedule detail")
	}

	var schedulesResponse []models.ScheduleResponse
	for _, schedule := range schedules {

		totalDuration, err := pkg.CalculateDuration(schedule.StartTime, schedule.EndTime)
		if err != nil {
			return nil, errors.New("failed to calculate duration")
		}

		schedulesResponse = append(schedulesResponse, models.ScheduleResponse{
			ID:         schedule.ID,
			Status:     schedule.Status,
			TeacherID:  schedule.TeacherID,
			Date:       schedule.Date.Format("2006-01-02"),
			StartTime:  schedule.StartTime,
			EndTime:    schedule.EndTime,
			TotalPrice: totalDuration * float64(schedule.Teacher.PricePerHour),
			Teacher: models.TeacherResponse{
				ID:           schedule.Teacher.ID,
				Name:         schedule.Teacher.Name,
				Bio:          schedule.Teacher.Bio,
				PricePerHour: schedule.Teacher.PricePerHour,
				ProfileImage: schedule.Teacher.ProfileImage,
			},
		})
	}

	return schedulesResponse, nil
}

func (s *ScheduleService) FetchTeacherIdAndIds(teacherID string, ids []string) ([]string, error) {

	idsInt := make([]uint, len(ids))
	for i, id := range ids {
		idsInt[i] = cast.ToUint(id)
	}

	schedules, err := s.scheduleRepo.FindByTeacherAndIDs(cast.ToUint(teacherID), idsInt)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to fetch teacher id and ids")
	}

	validIDs := make([]string, len(schedules))
	for i, s := range schedules {
		validIDs[i] = cast.ToString(s.ID)
	}

	return validIDs, nil
}

func (s *ScheduleService) Schedules(paginate *pkg.Paginate) (pkg.ResponsePaginate, error) {
	schedules, err := s.scheduleRepo.GetSchedules(*paginate)
	if err != nil {
		return pkg.ResponsePaginate{}, errors.New("failed to get schedules")
	}
	return schedules, nil
}

func (s *ScheduleService) Update(schedule *models.Schedule) error {

	_, err := s.scheduleRepo.GetSchedulesById(schedule.ID)
	if err != nil {
		return errors.New("schedule not found")
	}

	if err := s.scheduleRepo.UpdateSchedule(schedule.ID, schedule); err != nil {
		return errors.New("failed to update schedule")
	}
	return nil
}
