package routers

import (
	"assesment/delivery/controllers"
	infrastructure "assesment/Infrastructure"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes and configures the routes for the application.
func SetupRoutes(gino *gin.Engine, taskmgr *controllers.LoanController) {

	// Public routes
	// Route for user registration
	gino.POST("/register", taskmgr.RegisterUser)
	// Route for user login
	gino.POST("/login", taskmgr.LoginUser)

	// Protected routes group
	auth := gino.Group("/")
	// Apply authentication middleware to protected routes
	auth.Use(infrastructure.AuthMiddleware)
	{
		
		

		// Admin-specific endpoint group
		admin := auth.Group("/")
		// Apply admin middleware to admin-specific routes
		admin.Use(infrastructure.AdminMiddleware)
		{
			// Route to delete a user by ID (admin operation)
			admin.DELETE("/users/:id", taskmgr.DeleteUser)
		}
	}
}
