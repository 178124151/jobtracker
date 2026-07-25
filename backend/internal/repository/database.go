package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/jobtracker/backend/internal/config"
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *gorm.DB
	var err error

	// 重试连接数据库
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// 测试连接
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					log.Println("Database connected successfully")
					// 自动迁移
					if err := db.AutoMigrate(
						&model.User{},
						&model.Company{},
						&model.Application{},
					); err != nil {
						log.Printf("Auto migration failed: %v", err)
					}
					// 加载种子数据
					if err := SeedCompanies(db); err != nil {
						log.Printf("Seed data failed: %v", err)
					}
					return db, nil
				}
			}
		}
		log.Printf("Failed to connect to database (attempt %d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 30 attempts: %v", err)
}