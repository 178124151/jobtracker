package service

import (
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/jobtracker/backend/internal/model"
	"github.com/yourusername/jobtracker/backend/internal/repository"
)

type HealthCheckService struct {
	companyRepo *repository.CompanyRepository
	client      *http.Client
}

func NewHealthCheckService(companyRepo *repository.CompanyRepository) *HealthCheckService {
	return &HealthCheckService{
		companyRepo: companyRepo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckAllCompanies 检查所有公司网站状态
func (s *HealthCheckService) CheckAllCompanies() {
	companies, err := s.companyRepo.List("")
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 限制并发数

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
}

func (s *HealthCheckService) checkCompany(company *model.Company) {
	now := time.Now()
	
	resp, err := s.client.Get(company.Website)
	if err != nil {
		company.HealthStatus = "RED"
		company.LastChecked = &now
		s.companyRepo.Update(company)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		company.HealthStatus = "GREEN"
	} else {
		company.HealthStatus = "YELLOW"
	}
	
	company.LastChecked = &now
	s.companyRepo.Update(company)
}