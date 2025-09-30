package main

import (
	"log"
	"os"
	"stackwhiz-portfolio-backend/internal/api"
	"stackwhiz-portfolio-backend/internal/config"
	"stackwhiz-portfolio-backend/internal/database"
	"stackwhiz-portfolio-backend/internal/middleware"
	"stackwhiz-portfolio-backend/internal/repository"
	"stackwhiz-portfolio-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// @title Portfolio API
// @version 1.0
// @description Professional portfolio backend API for portfolio
// @termsOfService http://swagger.io/terms/

// @contact.name
// @contact.url https://github.com/StackWhiz
// @contact.email

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient := database.InitializeRedis(cfg.RedisURL)

	// Initialize repositories
	profileRepo := repository.NewProfileRepository(db)
	experienceRepo := repository.NewExperienceRepository(db)
	skillRepo := repository.NewSkillRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	contactRepo := repository.NewContactRepository(db)

	// Initialize services
	profileService := service.NewProfileService(profileRepo, redisClient)
	experienceService := service.NewExperienceService(experienceRepo, redisClient)
	skillService := service.NewSkillService(skillRepo, redisClient)
	projectService := service.NewProjectService(projectRepo, redisClient)
	contactService := service.NewContactService(contactRepo, redisClient)
	authService := service.NewAuthService(cfg.JWTSecret)

	// Initialize handlers
	handlers := api.NewHandlers(
		profileService,
		experienceService,
		skillService,
		projectService,
		contactService,
		authService,
	)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(handlers *api.Handlers, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())
	router.Use(middleware.SecurityHeaders())

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		{
			public.GET("/profile", handlers.GetProfile)
			public.GET("/experiences", handlers.GetExperiences)
			public.GET("/skills", handlers.GetSkills)
			public.GET("/projects", handlers.GetProjects)
			public.POST("/contact", handlers.CreateContact)
		}

		// Admin routes (protected)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			admin.PUT("/profile", handlers.UpdateProfile)
			admin.POST("/experiences", handlers.CreateExperience)
			admin.PUT("/experiences/:id", handlers.UpdateExperience)
			admin.DELETE("/experiences/:id", handlers.DeleteExperience)
			admin.POST("/skills", handlers.CreateSkill)
			admin.PUT("/skills/:id", handlers.UpdateSkill)
			admin.DELETE("/skills/:id", handlers.DeleteSkill)
			admin.POST("/projects", handlers.CreateProject)
			admin.PUT("/projects/:id", handlers.UpdateProject)
			admin.DELETE("/projects/:id", handlers.DeleteProject)
			admin.GET("/contacts", handlers.GetContacts)
			admin.PUT("/contacts/:id/status", handlers.UpdateContactStatus)
		}

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}
	}

	return router
}
