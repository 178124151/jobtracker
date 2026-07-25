package service

import (
	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/repository"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) ListByUser(userID string) ([]model.Application, error) {
	return s.repo.ListByUser(userID)
}

func (s *ApplicationService) Create(app *model.Application) error {
	return s.repo.Create(app)
}

func (s *ApplicationService) Update(app *model.Application) error {
	return s.repo.Update(app)
}

func (s *ApplicationService) Delete(id string) error {
	return s.repo.Delete(id)
}
