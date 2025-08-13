package models

type ScheduleResponse struct {
	ID         uint    `json:"id"`
	Status     string  `json:"status"`
	TeacherID  uint    `json:"teacher_id"`
	Date       string  `json:"date"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	TotalPrice float64 `json:"total_price"`
	Teacher    TeacherResponse
}

type ScheduleRequest struct {
	TeacherID uint   `json:"teacher_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

type ScheduleFilterRequest struct {
	TeacherID   string   `json:"teacher_id" binding:"required"`
	ScheduleIDs []string `json:"schedule_ids" binding:"required"`
}

type ScheduleFilterResponse struct {
	ValidScheduleIDs []string `json:"valid_schedule_ids"`
}
