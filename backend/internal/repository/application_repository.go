package repository

import (
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) ListByUser(userID string) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("applied_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) GetByID(id string) (*model.Application, error) {
	var app model.Application
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&app).Error
	return &app, err
}

func (r *ApplicationRepository) Create(app *model.Application) error {
	return r.db.Create(app).Error
}

func (r *ApplicationRepository) Update(app *model.Application) error {
	return r.db.Save(app).Error
}

func (r *ApplicationRepository) Delete(id string) error {
	return r.db.Delete(&model.Application{}, "id = ?", id).Error
}
