package controllers

import (
	"context"
	//"encoding/json"
	"net/http"
	"assesment/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related HTTP requests.
type UserController struct {
	userUsecase domain.UserUsecase
}

// NewUserController creates a new instance of UserController.
func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

// RegisterUser handles user registration.
func (uc *UserController) RegisterUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userUsecase.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// LoginUser handles user login and returns a JWT token.
func (uc *UserController) LoginUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	status, token, err := uc.userUsecase.LoginUser(context.Background(), user)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// UpdateUser handles updating a user's information.
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userUsecase.UpdateUser(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateUserPassword handles updating a user's password.
func (uc *UserController) UpdateUserPassword(c *gin.Context) {
	var request struct {
		ID           primitive.ObjectID `json:"id"`
		NewPassword  string             `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userUsecase.UpdateUserPassword(context.Background(), request.ID, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetUserByID retrieves a user by ID.
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := uc.userUsecase.GetUserByID(context.Background(), oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByEmail retrieves a user by email.
func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	user, err := uc.userUsecase.GetUserByEmail(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ActivateUser handles user account activation.
func (uc *UserController) ActivateUser(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	err := uc.userUsecase.ActivateUser(context.Background(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// RefreshToken handles refreshing a user's JWT token.
func (uc *UserController) RefreshToken(c *gin.Context) {
	var request struct {
		OldToken string `json:"old_token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newToken, err := uc.userUsecase.RefreshToken(context.Background(), request.OldToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// SendPasswordResetLink handles sending a password reset link to the user's email.
func (uc *UserController) SendPasswordResetLink(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userUsecase.SendPasswordResetLink(context.Background(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ResetPassword handles resetting the user's password using a reset token.
func (uc *UserController) ResetPassword(c *gin.Context) {
	var request struct {
		ResetToken   string `json:"reset_token"`
		NewPassword  string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userUsecase.ResetPassword(context.Background(), request.ResetToken, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}


func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = uc.userUsecase.DeleteUser(context.Background(), oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateUserRole handles updating a user's role.
