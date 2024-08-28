package domain

import (
	"context"
	//"net/http"
	"time"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Loan represents the loan data model.
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	Role           string             `json:"role"`
	IsActive       bool               `json:"isActivated"`
	Username       string             `json:"username"`
	ActivationToken string            `bson:"activation_token,omitempty" json:"activationToken,omitempty"`
	TokenCreatedAt time.Time          `bson:"token_created_at,omitempty" json:"tokenCreatedAt,omitempty"`
}

// UserRepository defines the methods for interacting with the user data storage in the domain layer.
type UserRepository interface {
	GetUserByID(id primitive.ObjectID) (User, error)
	GetUserByEmail(email string) (User, error)
	Register(user User) error
	LoginUser(user User) (int, string, error)
	UpdateUser(user User) error
	UpdateUserPassword(user User) error
	DeleteUser(id primitive.ObjectID) error
	
}


type UserUsecase interface {
	Register(user User) error
	LoginUser(c context.Context, user User) (int, string, error)
	UpdateUser(c context.Context, user User) error
	UpdateUserPassword(c context.Context, id primitive.ObjectID, newPassword string) error
	GetUserByID(c context.Context, id primitive.ObjectID) (User, error)
	GetUserByEmail(c context.Context, email string) (User, error)
	ActivateUser(c context.Context, token string) error
	RefreshToken(c context.Context, oldToken string) (string, error)
	SendPasswordResetLink(c context.Context, email string) error
	ResetPassword(c context.Context, resetToken string, newPassword string) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error }

// CustomError represents a custom error with a message and status code.
type CustomError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// New creates a new custom error.
func New(message string, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return e.Message
}

// Common errors
var (
	ErrUserAlreadyExists = New("user already exists", 409)
)
// Token represents the JWT token model.
type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id" json:"userId"`
	AccessToken  string             `bson:"access_token" json:"accessToken"`
	RefreshToken string             `bson:"refresh_token" json:"refreshToken"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expiresAt"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// JwtCustomClaims represents the custom claims for JWT tokens.
type JwtCustomClaims struct {
	Authorized  bool   `json:"authorized"`
	UserID      string `json:"user_id"`
	Role        string `json:"role"`
	Username    string `json:"username"`
	IsActivated bool   `json:"is_activated"`
	jwt.StandardClaims
}

// TokenGenerator defines the methods for generating tokens.
type TokenGenerator interface {
	GenerateToken(user User) (string, error)
	GenerateRefreshToken(user User) (string, error)
	RefreshToken(token string) (User, error)
	VerifyResetToken(token string) (*User, error)
}

// TokenVerifier defines the methods for verifying tokens.
type TokenVerifier interface {
	VerifyToken(token string) (*User, error)
	VerifyRefreshToken(token string) (*User, error)
}

// PasswordService defines the methods for password hashing and verification.
type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	PasswordHasher(password string) (string, error) 
}

// CustomError represents a custom error with a message and status code.
// type CustomError struct {
// 	Message    string `json:"message"`
// 	StatusCode int    `json:"status_code"`
// }

// New creates a new custom error.
// func New(message string, statusCode int) *CustomError {
// 	return &CustomError{
// 		Message:    message,
// 		StatusCode: statusCode,
// 	}
// }

// Error implements the error interface for CustomError.
// func (e *CustomError) Error() string {
// 	return e.Message
// }

// Common errors
// var (
// 	ErrNotFound                 = New("resource not found", http.StatusNotFound)
// 	ErrInternalServer           = New("internal server error", http.StatusInternalServerError)
// 	ErrBadRequest               = New("bad request", http.StatusBadRequest)
// 	ErrUnauthorized             = New("unauthorized", http.StatusUnauthorized)
// 	ErrForbidden                = New("forbidden", http.StatusForbidden)
// 	ErrUserNotFound             = New("user not found", http.StatusNotFound)
// 	ErrInvalidToken             = New("invalid token", http.StatusUnauthorized)
// 	ErrInvalidCredentials       = New("invalid credentials", http.StatusUnauthorized)
// 	//ErrUserAlreadyExists        = New("user already exists", http.StatusConflict)
// 	ErrFailedToCreateBlog       = New("failed to create blog", http.StatusInternalServerError)
// 	ErrBlogNotFound             = New("blog not found", http.StatusNotFound)
// 	ErrFailedToUpdateBlog       = New("failed to update blog", http.StatusInternalServerError)
// 	ErrFailedToRetrieveBlogs    = New("failed to retrieve blogs", http.StatusInternalServerError)
// 	ErrFailedToDeleteBlog       = New("failed to delete blog", http.StatusInternalServerError)
// 	ErrFailedToCreateComment    = New("failed to create comment", http.StatusInternalServerError)
// 	ErrCommentNotFound          = New("comment not found", http.StatusNotFound)
// 	ErrFailedToUpdateComment    = New("failed to update comment", http.StatusInternalServerError)
// 	ErrFailedToDeleteComment    = New("failed to delete comment", http.StatusInternalServerError)
// 	ErrFailedToLikeComment      = New("failed to like comment", http.StatusInternalServerError)
// 	ErrFailedToUnlikeComment    = New("failed to unlike comment", http.StatusInternalServerError)
// 	ErrFailedToCreateReply      = New("failed to create reply", http.StatusInternalServerError)
// 	ErrFailedToUpdateReply      = New("failed to update reply", http.StatusInternalServerError)
// 	ErrFailedToDeleteReply      = New("failed to delete reply", http.StatusInternalServerError)
// 	ErrFailedToGetComments      = New("failed to get comments", http.StatusInternalServerError)
// 	ErrInvalidCommentID         = New("invalid comment ID", http.StatusBadRequest)
// 	ErrInvalidPaginationParameters = New("invalid pagination parameters", http.StatusBadRequest)
// )


var (
	ErrMissingRequiredFields   = errors.New("missing required fields")
	ErrInvalidEmail            = errors.New("invalid email format")
	ErrInvalidPassword         = errors.New("invalid password format")
	//ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInternalServer          = errors.New("internal server error")
	ErrFailedToSendEmail       = errors.New("failed to send email")
	ErrInvalidUserID           = errors.New("invalid user ID")
	ErrUserNotFound            = errors.New("user not found")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("invalid token")
)