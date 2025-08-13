package handler

import (
	"net/http"
	"strconv"
	"teacher/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetTeacherDashboard(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	if teacherIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacher_id is required"})
		return
	}

	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher_id format"})
		return
	}

	dashboardData, err := h.dashboardService.GetTeacherDashboard(uint(teacherID))
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboardData)
}
