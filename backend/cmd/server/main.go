package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/config"
	"github.com/yourusername/jobtracker/backend/internal/handler"
	"github.com/yourusername/jobtracker/backend/internal/middleware"
	"github.com/yourusername/jobtracker/backend/internal/repository"
	"github.com/yourusername/jobtracker/backend/internal/service"
)

func main() {
	// 初始化结构化日志
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	companyRepo := repository.NewCompanyRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	userRepo := repository.NewUserRepository(db)
	resumeRepo := repository.NewResumeRepository(db)

	companySvc := service.NewCompanyService(companyRepo)
	appSvc := service.NewApplicationService(appRepo)
	userSvc := service.NewUserService(userRepo)
	healthCheckSvc := service.NewHealthCheckService(companyRepo)

	// 后台健康检查任务
	go func() {
		slog.Info("Running initial health check...")
		healthCheckSvc.CheckAllCompanies()
		slog.Info("Initial health check completed")

		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			slog.Info("Running scheduled health check...")
			healthCheckSvc.CheckAllCompanies()
			slog.Info("Health check completed")
		}
	}()

	companyHandler := handler.NewCompanyHandler(companySvc)
	appHandler := handler.NewApplicationHandler(appSvc)
	userHandler := handler.NewUserHandler(userSvc)
	resumeHandler := handler.NewResumeHandler(resumeRepo)
	healthHandler := handler.NewHealthHandler(db)

	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.StructuredLogger())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.Recovery())

	// 健康检查端点
	r.GET("/healthz", healthHandler.Liveness)
	r.GET("/readyz", healthHandler.Readiness)

	// API 路由
	api := r.Group("/api/v1")
	{
		// SRE 端点
		sre := api.Group("/sre")
		{
			sre.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
			sre.GET("/metrics", handler.GetMetrics)
			sre.GET("/prometheus", handler.GetPrometheusMetrics)
			sre.GET("/costs", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "costs endpoint"})
			})
		}

		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthRequired(), userHandler.Logout)
			auth.GET("/me", middleware.AuthRequired(), userHandler.GetMe)
		}

		// 公司路由
		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.List)
			companies.GET("/:id", companyHandler.Get)
			companies.POST("", middleware.AuthRequired(), companyHandler.Create)
			companies.PUT("/:id", middleware.AuthRequired(), companyHandler.Update)
			companies.DELETE("/:id", middleware.AuthRequired(), companyHandler.Delete)
		}

		// SME公司
		api.GET("/sme-companies", func(c *gin.Context) {
			jsonFile := "data/sme_companies.json"
			data, err := os.ReadFile(jsonFile)
			if err != nil {
				c.JSON(500, gin.H{"code": 5000, "message": "Failed to read SME companies data"})
				return
			}
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				c.JSON(500, gin.H{"code": 5000, "message": "Failed to parse SME companies data"})
				return
			}
			c.JSON(200, gin.H{"code": 0, "data": result["companies"], "message": "ok"})
		})

		// 投递记录
		applications := api.Group("/applications")
		{
			applications.GET("", middleware.AuthRequired(), appHandler.List)
			applications.POST("", middleware.AuthRequired(), appHandler.Create)
			applications.PUT("/:id", middleware.AuthRequired(), appHandler.Update)
			applications.DELETE("/:id", middleware.AuthRequired(), appHandler.Delete)
		}

		// 简历
		resumes := api.Group("/resumes")
		{
			resumes.GET("", resumeHandler.List)
			resumes.GET("/:id", resumeHandler.Get)
			resumes.POST("", resumeHandler.Save)
			resumes.DELETE("/:id", resumeHandler.Delete)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// 优雅关闭
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited gracefully")
}