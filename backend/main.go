package main

import (
	"log"
	"task-manager/backend/internal/handlers"
	"task-manager/backend/internal/middleware"
	"task-manager/backend/internal/repositories"
	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	dbCfg := repositories.NewDatabaseConfig()
	db, err := dbCfg.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}
	defer sqlDB.Close()

	if utils.GetEnv("JWT_SECRET", "") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// TODO(7-8): Replace the remaining nil services (NewTaskService()) as you implement them.
	authHandler := handlers.NewAuthHandler(db, services.NewAuthService())

	taskHandler := handlers.NewTaskHandler(db, nil)

	refreshHandler := handlers.NewRefreshHandler(db, services.NewAuthService())

	userHandler := handlers.NewUserHandler(db, services.NewUserService())

	registrationHandler := handlers.NewRegisterHandler(db, services.NewRegisterService())

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://host.docker.internal"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// TODO(9, optional): rate limiting middleware on /auth/login to slow brute-force attempts.
	v1 := r.Group("/api/v1")
	{
		// Public: no token needed.
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", registrationHandler.Registration)
			authRoutes.POST("/login", authHandler.Token)
			authRoutes.POST("/refresh", refreshHandler.Refresh)
		}
		taskRoutes := v1.Group("/tasks")
		taskRoutes.Use(middleware.AuthMiddleware())
		{
			taskRoutes.POST("", middleware.RequirePermission(utils.ResourceTask, utils.ActionWrite), taskHandler.CreateTask)
			taskRoutes.PUT("/:id", middleware.RequirePermission(utils.ResourceTask, utils.ActionWrite), taskHandler.UpdateTask)
			taskRoutes.DELETE("/:id", middleware.RequirePermission(utils.ResourceTask, utils.ActionDelete), taskHandler.DeleteTask)
			taskRoutes.GET("/:id", middleware.RequirePermission(utils.ResourceTask, utils.ActionRead), taskHandler.GetTaskByID)
			taskRoutes.GET("", middleware.RequirePermission(utils.ResourceTask, utils.ActionRead), taskHandler.GetTasks)
		}
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			userRoutes.DELETE("/:user_id", middleware.RequirePermission(utils.ResourceUser, utils.ActionDelete), userHandler.DeleteUser)
			userRoutes.GET("", middleware.RequirePermission(utils.ResourceUser, utils.ActionRead), userHandler.GetUsers)
			userRoutes.GET("/:user_id/tasks", middleware.RequirePermission(utils.ResourceTask, utils.ActionRead), middleware.RequireOwnerOrAdmin("user_id"), taskHandler.GetTasksByUser)
			userRoutes.GET("/profile", middleware.RequirePermission(utils.ResourceProfile, utils.ActionRead), userHandler.GetUserProfile)
			userRoutes.GET("/profile/:user_id", middleware.RequirePermission(utils.ResourceProfile, utils.ActionRead), middleware.RequireOwnerOrAdmin("user_id"), userHandler.GetUserProfileByUserId)
		}
	}
	r.Run(":8080")
}
