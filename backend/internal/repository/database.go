package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/yourusername/jobtracker/backend/internal/config"
	"github.com/yourusername/jobtracker/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *gorm.DB
	var err error

	// 重试连接数据库
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				// 连接池配置（D9）
				sqlDB.SetMaxOpenConns(25)                  // 最大打开连接数
				sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
				sqlDB.SetConnMaxLifetime(5 * time.Minute)  // 连接最大存活时间
				sqlDB.SetConnMaxIdleTime(1 * time.Minute)  // 空闲连接最大存活时间

				if err = sqlDB.Ping(); err == nil {
					slog.Info("Database connected successfully",
						"host", cfg.DBHost,
						"port", cfg.DBPort,
						"database", cfg.DBName,
						"max_open_conns", 25,
						"max_idle_conns", 10,
					)

					// 自动迁移
					if err := db.AutoMigrate(
						&model.User{},
						&model.Company{},
						&model.Application{},
						&model.Resume{},
					); err != nil {
						slog.Warn("Auto migration failed", "error", err)
					}

					// 加载种子数据
					if err := SeedCompanies(db); err != nil {
						slog.Warn("Seed data failed", "error", err)
					}

					return db, nil
				}
			}
		}
		slog.Warn("Failed to connect to database", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 30 attempts: %v", err)
}