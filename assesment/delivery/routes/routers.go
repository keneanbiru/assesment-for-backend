package routes
import (
	controllers "assesment/delivery/controllers"
	infrastructure "assesment/Infrastructure"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes and configures the routes for the application.
func SetupRoutes(gino *gin.Engine, userCtrl *controllers.UserController) {

	// Public routes
	// Route for user registration
	gino.POST("/auth/register", userCtrl.RegisterUser)
	// Route for user login
	gino.POST("/auth/login", userCtrl.LoginUser)
	// Route to send password reset link
	gino.POST("/auth/reset-password", userCtrl.SendPasswordResetLink)
	// Route to reset password using a token
	gino.POST("/auth/reset-password/:token", userCtrl.ResetPassword)
	// Route to activate a user's account
	gino.GET("/auth/activate/:token", userCtrl.ActivateUser)
	// Route to refresh a user's JWT token
	gino.POST("/auth/refresh-token", userCtrl.RefreshToken)

	// Protected routes group
	auth := gino.Group("/")
	// Apply authentication middleware to protected routes
	auth.Use(infrastructure.AuthMiddleware)
	{
		// Route to get the current user's profile
		auth.GET("/user/profile", userCtrl.GetUserByID)
		// Route to update the current user's profile
		auth.PUT("/user/update", userCtrl.UpdateUser)
		// Route to update the current user's password
		auth.POST("/user/update-password", userCtrl.UpdateUserPassword)

		// Admin-specific endpoint group
		admin := auth.Group("/")
		// Apply admin middleware to admin-specific routes
		admin.Use(infrastructure.AdminMiddleware)
		{
			// Route to get all users (admin operation)
			admin.GET("/admin/users", userCtrl.GetUserByID)
			// Route to get a user by ID (admin operation)
			admin.GET("/admin/users/:id", userCtrl.GetUserByID)
			// Route to delete a user by ID (admin operation)
			admin.DELETE("/admin/users/:id", userCtrl.DeleteUser)
			
		}
	}
}
