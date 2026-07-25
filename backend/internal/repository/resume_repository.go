package repository

import (
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/gorm"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) ListByUser(userID string) ([]model.Resume, error) {
	var resumes []model.Resume
	err := r.db.Where("user_id = ? OR user_id IS NULL OR user_id = ''", userID).
		Order("updated_at DESC").
		Find(&resumes).Error
	return resumes, err
}

func (r *ResumeRepository) GetByID(id string) (*model.Resume, error) {
	var resume model.Resume
	err := r.db.Where("id = ?", id).First(&resume).Error
	return &resume, err
}

func (r *ResumeRepository) GetDefault() (*model.Resume, error) {
	var resume model.Resume
	err := r.db.Where("is_default = true").First(&resume).Error
	if err != nil {
		// Return empty resume if no default found
		return &model.Resume{
			Title:    "Untitled Resume",
			Content:  "{}",
			Template: "classic",
		}, nil
	}
	return &resume, err
}

func (r *ResumeRepository) Create(resume *model.Resume) error {
	return r.db.Create(resume).Error
}

func (r *ResumeRepository) Update(resume *model.Resume) error {
	return r.db.Save(resume).Error
}

func (r *ResumeRepository) Delete(id string) error {
	return r.db.Delete(&model.Resume{}, "id = ?", id).Error
}