package pkg

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Paginate struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationPage struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	Limit       int `json:"limit"`
}

type ResponsePaginate struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationPage `json:"pagination"`
}

type TimeSchedule struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}

func ParseTimeSchedule(dateRaw, start, end string) (TimeSchedule, error) {
	date, err := time.Parse("2006-01-02", dateRaw)
	if err != nil {
		return TimeSchedule{}, errors.New("invalid date format, should be 2006-01-02")
	}
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return TimeSchedule{}, errors.New("invalid start time format, should be 15:04")
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return TimeSchedule{}, errors.New("invalid end time format, should be 15:04")
	}
	t := TimeSchedule{
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return t, nil
}

var DayMap = map[time.Weekday]string{
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
	time.Sunday:    "Minggu",
}

func IsDayAvailableForTeacher(availableDaysStr string, bookingDate time.Time) bool {
	availableDays := strings.Split(availableDaysStr, ",")
	bookingDayIndo := DayMap[bookingDate.Weekday()]
	for _, day := range availableDays {
		if strings.TrimSpace(day) == bookingDayIndo {
			return true
		}
	}
	return false
}

func ParseHHMM(input string) (time.Time, error) {
	return time.Parse("15:04", input)
}

func IsWithinRange(startStr, endStr, checkStr string) (bool, error) {
	start, err := ParseHHMM(startStr)
	if err != nil {
		return false, err
	}
	end, err := ParseHHMM(endStr)
	if err != nil {
		return false, err
	}
	check, err := ParseHHMM(checkStr)
	if err != nil {
		return false, err
	}
	return (check.Equal(start) || check.After(start)) && check.Before(end), nil
}

func ParseTime(input string) (time.Time, error) {
	t, err := time.Parse("15:04", input)
	if err != nil {
		return time.Time{}, errors.New("invalid time format, should be 15:04")
	}
	return t, nil
}

func NormalizeTime(input time.Time) string {
	t := input
	return t.Format("15:04")
}

func CalculateDuration(start string, end string) (durationHours float64, err error) {
	layout := "15:04"

	startTime, err := time.Parse(layout, start)
	if err != nil {
		return 0, fmt.Errorf("failed to parse start time: %w", err)
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return 0, fmt.Errorf("failed to parse end time: %w", err)
	}

	// Hitung durasi dalam jam
	duration := endTime.Sub(startTime).Hours()
	if duration <= 0 {
		return 0, fmt.Errorf("end time must be after start time")
	}

	return duration, nil
}
