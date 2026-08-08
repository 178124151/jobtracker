package service

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/repository"
	"github.com/yourusername/jobtracker/backend/pkg/utils"
)

type HealthCheckService struct {
	companyRepo *repository.CompanyRepository
	client      *utils.HTTPClient
}

func NewHealthCheckService(companyRepo *repository.CompanyRepository) *HealthCheckService {
	return &HealthCheckService{
		companyRepo: companyRepo,
		client:      utils.NewHTTPClient(20*time.Second, 2), // 20s超时，重试2次
	}
}

// CheckAllCompanies 检查所有公司网站状态
func (s *HealthCheckService) CheckAllCompanies() {
	companies, err := s.companyRepo.List("")
	if err != nil {
		slog.Error("Failed to list companies for health check", "error", err)
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 并发限制

	for _, company := range companies {
		wg.Add(1)
		go func(c model.Company) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			s.checkCompany(&c)
		}(company)
	}

	wg.Wait()
	slog.Info("Health check completed", "companies_count", len(companies))
}

func (s *HealthCheckService) checkCompany(company *model.Company) {
	now := time.Now()

	statusCode, _, err := s.client.Get(company.Website)

	if err != nil {
		slog.Warn("Health check failed",
			"company", company.Name,
			"website", company.Website,
			"error", err,
		)
		var timeoutErr *utils.TimeoutError
		if errors.As(err, &timeoutErr) {
			company.HealthStatus = "YELLOW"
		} else {
			company.HealthStatus = "RED"
		}
		company.LastChecked = &now
		s.companyRepo.Update(company)
		return
	}

	if statusCode >= 200 && statusCode < 400 {
		company.HealthStatus = "GREEN"
	} else {
		company.HealthStatus = "YELLOW"
	}

	company.LastChecked = &now
	s.companyRepo.Update(company)
}
