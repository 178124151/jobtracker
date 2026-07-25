package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/config"
	"github.com/yourusername/jobtracker/backend/internal/handler"
	"github.com/yourusername/jobtracker/backend/internal/middleware"
	"github.com/yourusername/jobtracker/backend/internal/repository"
	"github.com/yourusername/jobtracker/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	companyRepo := repository.NewCompanyRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	userRepo := repository.NewUserRepository(db)
	resumeRepo := repository.NewResumeRepository(db)

	companySvc := service.NewCompanyService(companyRepo)
	appSvc := service.NewApplicationService(appRepo)
	userSvc := service.NewUserService(userRepo)
	healthCheckSvc := service.NewHealthCheckService(companyRepo)

	go func() {
		log.Println("Running initial health check...")
		healthCheckSvc.CheckAllCompanies()
		log.Println("Initial health check completed")

		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Running scheduled health check...")
			healthCheckSvc.CheckAllCompanies()
			log.Println("Health check completed")
		}
	}()

	companyHandler := handler.NewCompanyHandler(companySvc)
	appHandler := handler.NewApplicationHandler(appSvc)
	userHandler := handler.NewUserHandler(userSvc)
	resumeHandler := handler.NewResumeHandler(resumeRepo)

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthRequired(), userHandler.Logout)
			auth.GET("/me", middleware.AuthRequired(), userHandler.GetMe)
		}

		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.List)
			companies.GET("/:id", companyHandler.Get)
			companies.POST("", middleware.AuthRequired(), companyHandler.Create)
			companies.PUT("/:id", middleware.AuthRequired(), companyHandler.Update)
			companies.DELETE("/:id", middleware.AuthRequired(), companyHandler.Delete)
		}

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

		applications := api.Group("/applications")
		{
			applications.GET("", middleware.AuthRequired(), appHandler.List)
			applications.POST("", middleware.AuthRequired(), appHandler.Create)
			applications.PUT("/:id", middleware.AuthRequired(), appHandler.Update)
			applications.DELETE("/:id", middleware.AuthRequired(), appHandler.Delete)
		}

		// Resume routes
		resumes := api.Group("/resumes")
		{
			resumes.GET("", resumeHandler.List)
			resumes.GET("/:id", resumeHandler.Get)
			resumes.POST("", resumeHandler.Save)
			resumes.DELETE("/:id", resumeHandler.Delete)
		}

		sre := api.Group("/sre")
		{
			sre.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
			sre.GET("/costs", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "costs endpoint"})
			})
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}