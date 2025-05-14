package routes

import (
	"hospital-management/controllers"
	"hospital-management/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", controllers.LoginUser)
	}

	// Protected routes (JWT required)
	protected := router.Group("/api")
	protected.Use(middlewares.JWTMiddleware())
	{
		// User management (admin-only in real apps)
		protected.POST("/users", controllers.RegisterUser)

		// Patients routes — shared between roles, role check happens in controller
		protected.POST("/patients", middlewares.RoleMiddleware("receptionist"), controllers.CreatePatient)
		protected.GET("/patients", controllers.GetAllPatients)    // Shared route
		protected.GET("/patients/:id", controllers.GetPatient)    // Shared route
		protected.PUT("/patients/:id", controllers.UpdatePatient) // Shared route
		protected.DELETE("/patients/:id", middlewares.RoleMiddleware("receptionist"), controllers.DeletePatient)
	}
}
