package service

import (
	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/repository"
)

type CompanyService struct {
	repo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) List(group string) ([]model.Company, error) {
	return s.repo.List(group)
}

func (s *CompanyService) GetByID(id string) (*model.Company, error) {
	return s.repo.GetByID(id)
}

func (s *CompanyService) Create(company *model.Company) error {
	return s.repo.Create(company)
}

func (s *CompanyService) Update(company *model.Company) error {
	return s.repo.Update(company)
}

func (s *CompanyService) Delete(id string) error {
	return s.repo.Delete(id)
}
