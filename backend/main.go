package main

import (
	"log"
	"task-manager/backend/internal/handlers"
	"task-manager/backend/internal/repositories"
	"task-manager/backend/internal/services"
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

	// TODO(1-8): Replace each nil below with its service (services.NewAuthService(), NewTaskService(),
	//            NewUserService(), NewRegisterService()) as you implement them; the impls must satisfy the interfaces.
	// TODO(2): Fail fast at startup if JWT_SECRET is empty.
	authHandler := handlers.NewAuthHandler(db, nil)

	taskHandler := handlers.NewTaskHandler(db, nil)

	refreshHandler := handlers.NewRefreshHandler(db, nil)

	userHandler := handlers.NewUserHandler(db, nil)

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
		// TODO(6): Everything below must be protected: add middleware.AuthMiddleware() with .Use() on these groups
		//          and middleware.RequirePermission("task", "read"|"write"|"delete") per route.
		taskRoutes := v1.Group("/tasks")
		{
			taskRoutes.POST("", taskHandler.CreateTask)
			taskRoutes.PUT("/:id", taskHandler.UpdateTask)
			taskRoutes.DELETE("/:id", taskHandler.DeleteTask)
			taskRoutes.GET("/:id", taskHandler.GetTaskByID)
			taskRoutes.GET("", taskHandler.GetTasks)
		}
		// TODO(6): /users management (list, delete) is admin-only; own profile/tasks are allowed for the owner.
		userRoutes := v1.Group("/users")
		{
			userRoutes.DELETE("/:user_id", userHandler.DeleteUser)
			userRoutes.GET("", userHandler.GetUsers)
			userRoutes.GET("/:user_id/tasks", taskHandler.GetTasksByUser)
			userRoutes.GET("/profile", userHandler.GetUserProfile)
			userRoutes.GET("/profile/:user_id", userHandler.GetUserProfileByUserId)
		}
	}
	r.Run(":8080")
}
