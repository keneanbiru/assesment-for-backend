package controllers

import (
	"net/http"
	domain "assesment/domain"

	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanController struct {
	
	userusecase domain.UserUsecase
}

// NewTaskController creates a new instance of TaskController with the provided task and user usecases.
func NewLoanController( usermgr domain.UserUsecase) *LoanController {
	return &LoanController{
		
		userusecase: usermgr,
	}
}



// RegisterUser handles the registration of a new user.
func (controller *LoanController) RegisterUser(c *gin.Context) {
	var user domain.User

	// Bind the JSON request body to the user object
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	// Register the user using the user usecase
	err = controller.userusecase.Register(user)
	if err != nil {
		switch err {
		case domain.ErrMissingRequiredFields:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		case domain.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case domain.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password format"})
		case domain.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		case domain.ErrFailedToSendEmail:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send activation email"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server errorrrrrrrr"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully registered. Check your email for activation."})
}


// LoginUser handles user login and returns a token.
func (controller *LoanController) LoginUser(c *gin.Context) {
	var user domain.User

	// Bind the JSON request body to the user object
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Authenticate the user and get a token
	code, token, err := controller.userusecase.LoginUser(c, user)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully logged in", "token": token})
}

// DeleteUser removes a user identified by their ID (admin operation).
func (controller *LoanController) DeleteUser(c *gin.Context) {
	id := c.Param("id") // Extract the user ID from the URL path

	// Delete the user using the user usecase
	code, erro := controller.userusecase.DeleteUser(c, id)
	if erro == nil {
		c.IndentedJSON(code, gin.H{"message": "User successfully deleted"})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
