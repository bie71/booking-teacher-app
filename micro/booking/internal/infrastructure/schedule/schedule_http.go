package schedule

import (
	"booking/internal/config"
	"booking/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type ScheduleResponse struct {
	ID        uint   `json:"id"`
	Status    string `json:"status"`
	TeacherID uint   `json:"teacher_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ScheduleHttp struct {
	restyClient *resty.Client
	service     config.Service
}

func NewScheduleHttp(service config.Service, restyClient *resty.Client) *ScheduleHttp {
	return &ScheduleHttp{
		restyClient: restyClient,
		service:     service,
	}
}

func (s *ScheduleHttp) CheckScheduleAvailability(id uint) (*ScheduleResponse, error) {
	url := fmt.Sprintf("%s:%s/api/v1/schedule/%d", s.service.Host, s.service.Port, id)

	resp, err := s.restyClient.R().
		SetResult(&ScheduleResponse{}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error contacting teacher service: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("schedule not found or unavailable")
	}

	schedule := resp.Result().(*ScheduleResponse)

	return schedule, nil
}

func (s *ScheduleHttp) UpdateScheduleStatus(c *gin.Context, scheduleID uint, status string) error {
	url := fmt.Sprintf("%s:%s/api/v1/schedule-status", s.service.Host, s.service.Port)

	query := map[string]string{
		"id":     fmt.Sprintf("%d", scheduleID),
		"status": status,
	}

	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.GetHeader("Authorization")).
		SetQueryParams(query).
		Put(url)

	if err != nil {
		return fmt.Errorf("failed to update schedule status: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("teacher service returned non-200: %v", resp.Status())
	}

	return nil
}
func (s *ScheduleHttp) CancelSchedule(c *gin.Context, scheduleID uint) error {
	url := fmt.Sprintf("%s:%s/api/v1/cancel-schedule/%d", s.service.Host, s.service.Port, scheduleID)

	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.GetHeader("Authorization")).
		Put(url)

	if err != nil {
		return fmt.Errorf("failed to update schedule status: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("teacher service returned non-200: %v", resp.Status())
	}

	return nil
}

func (s *ScheduleHttp) FetchSchedulesByIDs(c *gin.Context, scheduleIDs []uint) (map[uint]model.ScheduleResponse, error) {
	url := fmt.Sprintf("%s:%s/api/v1/schedule/batch-detail", s.service.Host, s.service.Port)

	schedules := []model.ScheduleResponse{}

	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.GetHeader("Authorization")).
		SetBody(gin.H{"ids": scheduleIDs}).
		SetResult(&schedules).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedules by IDs: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("teacher service returned non-200: %v", resp.Status())
	}

	result := make(map[uint]model.ScheduleResponse)
	for _, v := range schedules {
		result[v.ID] = v
	}

	return result, nil

}

func (s *ScheduleHttp) CallScheduleServiceToGetByTeacher(teacherID string, scheduleIDs []string) (*ScheduleFilterResponse, error) {
	req := ScheduleFilterRequest{
		TeacherID:   teacherID,
		ScheduleIDs: scheduleIDs,
	}

	url := fmt.Sprintf("%s:%s/api/v1/schedules/filter-by-teacher", s.service.Host, s.service.Port)

	var resp ScheduleFilterResponse

	scheduleFilterResponse, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&resp).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedules by IDs: %v", err)
	}
	if scheduleFilterResponse.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("teacher service returned non-200: %v", scheduleFilterResponse.Status())
	}

	return &resp, nil

}

func (s *ScheduleHttp) GetScheduleByID(id uint) (*model.ScheduleResponse, error) {
	url := fmt.Sprintf("%s:%s/api/v1/schedule/%d", s.service.Host, s.service.Port, id)

	var schedule *model.ScheduleResponse
	resp, err := s.restyClient.R().
		SetResult(&schedule).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error contacting teacher service: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("schedule not found or unavailable")
	}

	return schedule, nil
}
