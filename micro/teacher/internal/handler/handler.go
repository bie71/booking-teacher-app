package handler

import (
	"net/http"
	"teacher/internal/models"
	"teacher/internal/pkg"
	"teacher/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Handler struct {
	teacherService *service.Service
}

func NewHandler(teacherService *service.Service) *Handler {
	return &Handler{
		teacherService: teacherService,
	}
}

func (h *Handler) GetTeachers(c *gin.Context) {
	var paginate pkg.Paginate
	page := c.Query("page")
	limit := c.Query("limit")

	if page != "" {
		paginate.Page = cast.ToInt(page)
	}
	if limit != "" {
		paginate.Limit = cast.ToInt(limit)
	}

	response, err := h.teacherService.GetTeachers(&paginate)
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *Handler) CreateTeacher(c *gin.Context) {
	var teacher models.TeacherRequest
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.teacherService.CreateTeacher(teacher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

func (h *Handler) DeleteTeacher(c *gin.Context) {

	var id string
	id = c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return

	}

	var teacher models.TeacherRequest
	teacher.ID = cast.ToUint(id)

	err := h.teacherService.DeleteTeacher(teacher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) UpdateTeacher(c *gin.Context) {

	params := c.Params.ByName("id")
	if params == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var teacher models.TeacherRequest
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacher.ID = cast.ToUint(params)

	err := h.teacherService.UpdateTeacher(teacher)
	if err != nil {
		if err.Error() == "invalid time format, should be 2006-01-02 15:04" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) GetTeacher(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	teacher, err := h.teacherService.GetTeacherByID(cast.ToUint(id))
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) CountTeachers(c *gin.Context) {
	count, err := h.teacherService.CountTeachers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_teachers": count})
}


func (h *Handler) GetMe(c *gin.Context) {
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"})
		return
	}
	userID := uidVal.(uint)
	teacher, err := h.teacherService.GetTeacherByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"teacher not found for this user"})
		return
	}
	c.JSON(http.StatusOK, teacher)
}
