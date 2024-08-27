package routes

import (
	controllers "assesment/delivery/controllers"
	infrastructure "assesment/infrastructure"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes and configures the routes for the application.
func SetupRoutes(gino *gin.Engine, userCtrl *controllers.UserController, loanCtrl *controllers.LoanController) {

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

	// Public routes for loans
	// Route to apply for a loan
	gino.POST("/loans", loanCtrl.ApplyForLoan)
	// Route to get a loan by ID
	gino.GET("/loans/:id", loanCtrl.GetLoanByID)
	// Route to get all loans with optional filtering and sorting
	gino.GET("/loans", loanCtrl.GetAllLoans)

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
			
			// Admin-specific routes for loans
			// Route to approve a loan
			admin.POST("/admin/loans/:id/approve", loanCtrl.ApproveLoan)
			// Route to reject a loan
			admin.POST("/admin/loans/:id/reject", loanCtrl.RejectLoan)
			// Route to delete a loan
			admin.DELETE("/admin/loans/:id", loanCtrl.DeleteLoan)
		}
	}
}
