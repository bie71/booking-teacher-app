package handler

import (
	"errors"
	"net/http"
	"teacher/internal/models"
	"teacher/internal/pkg"
	"teacher/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (s *ScheduleHandler) BookSchedule(c *gin.Context) {
	var req struct {
		TeacherID uint   `json:"teacher_id"`
		Date      string `json:"date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	parseTime, err := pkg.ParseTimeSchedule(req.Date, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := parseTime.Date
	startTime := parseTime.StartTime
	endTime := parseTime.EndTime

	schedule := models.Schedule{
		TeacherID: req.TeacherID,
		Date:      date,
		StartTime: pkg.NormalizeTime(startTime),
		EndTime:   pkg.NormalizeTime(endTime),
	}

	result, err := s.scheduleService.BookScheduleService(schedule)
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booked", "schedule": result})
}

func (s *ScheduleHandler) GetScheduleAvailable(c *gin.Context) {
	var paginate pkg.Paginate
	if err := c.ShouldBindQuery(&paginate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.scheduleService.GetAvailableScheduleService(cast.ToUint(c.Params.ByName("teacher_id")), &paginate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (s *ScheduleHandler) CancelSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	err := s.scheduleService.CancelScheduleService(cast.ToUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule canceled"})

}

func (s *ScheduleHandler) GetScheduleById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	schedule, err := s.scheduleService.GetScheduleService(cast.ToUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (s *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	err := s.scheduleService.DeleteScheduleService(cast.ToUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

func (s *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parseTime, err := pkg.ParseTimeSchedule(req.Date, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := parseTime.Date
	startTime := parseTime.StartTime
	endTime := parseTime.EndTime

	schedule := models.Schedule{
		TeacherID: req.TeacherID,
		Date:      date,
		StartTime: pkg.NormalizeTime(startTime),
		EndTime:   pkg.NormalizeTime(endTime),
	}

	err = s.scheduleService.CreateScheduleService(&schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule created"})
}

func (s *ScheduleHandler) UpdateScheduleStatus(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	status := c.Query("status")

	err := s.scheduleService.UpdateScheduleService(cast.ToUint(id), status)
	if err != nil {
		if err.Error() == "schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule status updated"})

}

func (s *ScheduleHandler) GetBatchScheduleDetail(c *gin.Context) {
	var ids struct {
		Ids []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schedules, err := s.scheduleService.GetBatchScheduleDetailService(ids.Ids)
	if err != nil {
		if errors.Is(err, errors.New("schedule not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (s *ScheduleHandler) FilterByTeacher(c *gin.Context) {
	var req models.ScheduleFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validIDs, err := s.scheduleService.FetchTeacherIdAndIds(req.TeacherID, req.ScheduleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get schedules"})
		return
	}

	c.JSON(http.StatusOK, models.ScheduleFilterResponse{
		ValidScheduleIDs: validIDs,
	})
}

func (s *ScheduleHandler) GetSchedules(c *gin.Context) {
	var paginate pkg.Paginate
	if err := c.ShouldBindQuery(&paginate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.scheduleService.Schedules(&paginate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	parseTime, err := pkg.ParseTimeSchedule(req.Date, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := parseTime.Date
	startTime := parseTime.StartTime
	endTime := parseTime.EndTime

	schedule := models.Schedule{
		ID:        cast.ToUint(id),
		Date:      date,
		StartTime: pkg.NormalizeTime(startTime),
		EndTime:   pkg.NormalizeTime(endTime),
		Status:    req.Status,
	}

	err = s.scheduleService.Update(&schedule)
	if err != nil {
		if err.Error() == "schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated"})
}
