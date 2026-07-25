package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/gorm"
)

type CompanySeedData struct {
	Version     string `json:"version"`
	LastUpdated string `json:"last_updated"`
	Companies   []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		NameEN      string `json:"name_en"`
		Website     string `json:"website"`
		Industry    string `json:"industry"`
		Group       string `json:"group"`
		Description string `json:"description"`
		CareersPage string `json:"careers_page"`
	} `json:"companies"`
}

func SeedCompanies(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&model.Company{}).Count(&count)
	if count > 0 {
		log.Printf("Companies table already has %d records, skipping seed", count)
		return nil
	}

	// 读取JSON文件
	jsonFile := "data/companies.json"
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		log.Printf("Seed file %s not found, skipping seed", jsonFile)
		return nil
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	var seedData CompanySeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return err
	}

	// 插入数据
	for _, c := range seedData.Companies {
		company := model.Company{
			ID:           uuid.New().String(),
			Name:         c.Name,
			NameEN:       c.NameEN,
			Website:      c.Website,
			Industry:     c.Industry,
			Group:        c.Group,
			Description:  c.Description,
			IsPreset:     true,
			HealthStatus: "UNKNOWN",
		}
		if err := db.Create(&company).Error; err != nil {
			log.Printf("Failed to seed company %s: %v", c.Name, err)
		}
	}

	log.Printf("Seeded %d companies from %s", len(seedData.Companies), jsonFile)
	return nil
}