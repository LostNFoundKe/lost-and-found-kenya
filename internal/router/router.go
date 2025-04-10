package router

import (
	"github.com/gin-gonic/gin"
	"lostnfound-api/internal/config"
	"lostnfound-api/internal/handler"
	"lostnfound-api/internal/middleware"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter(
	cfg *config.Config,

	itemHandler *handler.ItemHandler,

) *gin.Engine {
	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Public routes
		//api.POST("/register", authHandler.Register)
		//api.POST("/login", authHandler.Login)

		// Item public routes

		/// TODO
		//api.GET("/items/public", itemHandler.ListPublic)
		//api.GET("/items/public/:id", itemHandler.GetPublicByID)
		//api.GET("/items/search", itemHandler.Search)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.JWT(cfg.JWTSecret))
		{
			// Item routes
			protected.POST("/items", itemHandler.Create)
			protected.GET("/items", itemHandler.List)
			protected.GET("/items/:id", itemHandler.GetByID)
			protected.PUT("/items/:id", itemHandler.Update)
			protected.DELETE("/items/:id", itemHandler.Delete)

			/// TODO

			//protected.POST("/items/:id/images", itemHandler.UploadImage)

			// User routes
			//protected.GET("/users/me", userHandler.GetProfile)
			//protected.PUT("/users/me", userHandler.UpdateProfile)

			/// TODO

			// Admin routes
			//admin := protected.Group("/admin")
			//admin.Use(middleware.AdminOnly())
			{
				//admin.GET("/users", userHandler.ListUsers)
				//admin.PUT("/users/:id", userHandler.UpdateUser)
				//admin.DELETE("/users/:id", userHandler.DeleteUser)
			}
		}
	}

	return router
}
