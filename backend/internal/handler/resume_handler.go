package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/repository"
)

type ResumeHandler struct {
	repo *repository.ResumeRepository
}

func NewResumeHandler(repo *repository.ResumeRepository) *ResumeHandler {
	return &ResumeHandler{repo: repo}
}

func (h *ResumeHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	resumes, err := h.repo.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to fetch resumes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    resumes,
	})
}

func (h *ResumeHandler) Get(c *gin.Context) {
	id := c.Param("id")

	resume, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1201,
			"message": "Resume not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    resume,
	})
}

type SaveResumeRequest struct {
	ID       string `json:"id"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Template string `json:"template"`
}

func (h *ResumeHandler) Save(c *gin.Context) {
	var req SaveResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1101,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	if req.ID != "" {
		// Update existing
		resume, err := h.repo.GetByID(req.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    1201,
				"message": "Resume not found",
			})
			return
		}

		resume.Title = req.Title
		if req.Content != "" {
			resume.Content = req.Content
		}
		if req.Template != "" {
			resume.Template = req.Template
		}

		if err := h.repo.Update(resume); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "Failed to update resume",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data":    resume,
		})
	} else {
		// Create new
		resume := &model.Resume{
			UserID:    userID,
			Title:     req.Title,
			Content:   req.Content,
			Template:  req.Template,
			IsDefault: true,
		}

		if resume.Template == "" {
			resume.Template = "classic"
		}

		if err := h.repo.Create(resume); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "Failed to create resume",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data":    resume,
		})
	}
}

func (h *ResumeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to delete resume",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}