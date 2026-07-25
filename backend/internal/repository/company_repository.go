package repository

import (
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) List(group string) ([]model.Company, error) {
	var companies []model.Company
	query := r.db.Where("deleted_at IS NULL")

	if group != "" && group != "all" {
		query = query.Where("group = ?", group)
	}

	err := query.Order("name").Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) GetByID(id string) (*model.Company, error) {
	var company model.Company
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&company).Error
	return &company, err
}

func (r *CompanyRepository) Create(company *model.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) Update(company *model.Company) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) Delete(id string) error {
	return r.db.Delete(&model.Company{}, "id = ?", id).Error
}
