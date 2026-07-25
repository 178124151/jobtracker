package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/service"
)

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

func (h *ApplicationHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	apps, err := h.svc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to fetch applications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    apps,
	})
}

type CreateApplicationRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	Position    string `json:"position" binding:"required"`
	AppliedAt   string `json:"applied_at" binding:"required"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1101,
			"message": "Invalid request",
		})
		return
	}

	userID := c.GetString("user_id")

	app := &model.Application{
		UserID:      userID,
		CompanyName: req.CompanyName,
		Position:    req.Position,
		Status:      req.Status,
		Notes:       req.Notes,
	}

	if err := h.svc.Create(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to create application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    app,
	})
}

func (h *ApplicationHandler) Update(c *gin.Context) {
	// TODO: Update application
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func (h *ApplicationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to delete application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
