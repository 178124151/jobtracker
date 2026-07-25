package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/service"
)

type CompanyHandler struct {
	svc *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

func (h *CompanyHandler) List(c *gin.Context) {
	group := c.Query("group")

	companies, err := h.svc.List(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to fetch companies",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    companies,
	})
}

func (h *CompanyHandler) Get(c *gin.Context) {
	id := c.Param("id")

	company, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1201,
			"message": "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    company,
	})
}

type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Website     string `json:"website" binding:"required"`
	Industry    string `json:"industry"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1101,
			"message": "Invalid request",
		})
		return
	}

	// TODO: Create company
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func (h *CompanyHandler) Update(c *gin.Context) {
	// TODO: Update company
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func (h *CompanyHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "Failed to delete company",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
