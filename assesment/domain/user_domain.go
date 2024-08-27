package domain

import (
	"context"
	"time"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt"
	
)

// Task represents the task data model.
type Loan struct {
    LoanID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    BorrowerName  string             `bson:"borrower_name" json:"borrowerName"`
    Principal     float64            `bson:"principal" json:"principal"`
    InterestRate  float64            `bson:"interest_rate" json:"interestRate"`
    TermMonths    int                `bson:"term_months" json:"termMonths"`
    StartDate     time.Time          `bson:"start_date" json:"startDate"`
    MonthlyPayment float64           `bson:"monthly_payment" json:"monthlyPayment"`
}


type LoanRepository interface {
    CreateLoan(ctx context.Context, loan *Loan) (primitive.ObjectID, error)
    GetLoanByID(ctx context.Context, id primitive.ObjectID) (*Loan, error)
    UpdateLoan(ctx context.Context, loan *Loan) error
    DeleteLoan(ctx context.Context, id primitive.ObjectID) error
    ListLoans(ctx context.Context) ([]*Loan, error)
}

// // TaskUsecase defines the business logic methods for tasks.
// type TaskUsecase interface {
// 	// GetAllTasks retrieves all tasks for a user or all users if admin.
// 	GetAllTasks(c context.Context, isadmin bool, userid primitive.ObjectID) ([]Task, error)
// 	// GetTask retrieves a specific task by ID.
// 	GetTask(c context.Context, id string, isadmin bool, userid string) (Task, error)
// 	// AddTask creates a new task in the storage.
// 	AddTask(c context.Context, task Task) error
// 	// SetTask updates an existing task by ID.
// 	UpdateTask(c context.Context, id string, updatedTask Task, isadmin bool) error
// 	// DeleteTask removes a task by ID.
// 	DeleteTask(c context.Context, id string, userid string, isadmin bool) error
// }

/*
=========== The Models and Interfaces for User =============
*/

// User represents the user data model.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Role    string              `json:"role"`

	IsActive bool			`json:"isactivated"`
	Username string              `json:"username"`
	ActivationToken string         `bson:"activation_token,omitempty" json:"activation_token,omitempty"`
	TokenCreatedAt  time.Time      `bson:"token_created_at"`
}

// UserRepository defines the methods for interacting with the user data storage.
type UserRepository interface {
	// RegisterUserDb registers a new user in the database.
	Register(user User) (error)
	// LoginUserDb authenticates a user and returns a status code and token.
	LoginUser(user User) (int, string, error)
	// DeleteUser removes a user by ID.
	DeleteUser(id string) (int, error)
	GetUserByUsernameOrEmail(username, email string) (User, error)
}

// UserUsecase defines the business logic methods for users.
type UserUsecase interface {
	// RegisterUserDb registers a new user in the database.
	Register(user User) (error)
	// LoginUserDb authenticates a user and returns a status code and token.
	LoginUser(c context.Context, user User) (int, string, error)
	// DeleteUser removes a user by ID.
	DeleteUser(c context.Context, id string) (int, error)
}











type CustomError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func New(message string, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *CustomError) Error() string {
	return e.Message
}

var (
	// General Errors
	ErrNotFound       = New("resource not found", http.StatusNotFound)
	ErrInternalServer = New("internal server error", http.StatusInternalServerError)
	ErrBadRequest     = New("bad request", http.StatusBadRequest)
	ErrUnauthorized   = New("unauthorized", http.StatusUnauthorized)
	ErrForbidden      = New("forbidden", http.StatusForbidden)

	// Auth-specific Errors
	ErrUserNotFound         = New("user not found", http.StatusNotFound)
	ErrRefreshTokenNotFound = New("refresh token not found", http.StatusNotFound)
	ErrInvalidRefreshToken  = New("invalid refresh token", http.StatusUnauthorized)
	ErrInvalidAccessToken   = New("invalid access token", http.StatusUnauthorized)
	ErrExpiredAccessToken   = New("expired access token", http.StatusUnauthorized)
	ErrExpiredRefreshToken  = New("expired refresh token", http.StatusUnauthorized)
	ErrInvalidToken         = New("invalid token", http.StatusUnauthorized)
	ErrInvalidRole          = New("invalid role", http.StatusUnauthorized)
	ErrInvalidUserID        = New("invalid user id", http.StatusBadRequest)
	ErrInvalidDeviceID      = New("invalid device id", http.StatusBadRequest)
	ErrDeviceNotFound       = New("device not found", http.StatusNotFound)
	ErrInvalidEmail         = New("invalid email", http.StatusBadRequest)
	ErrInvalidPassword      = New("invalid password", http.StatusBadRequest)

	ErrFailedToUpdateUser    = New("failed to update user", http.StatusInternalServerError)
	ErrMissingRequiredFields = New("missing required fields", http.StatusBadRequest)
	ErrInvalidUpdateRequest  = New("invalid update request", http.StatusBadRequest)
	ErrFailedToSendEmail     = New("failed to send email", http.StatusInternalServerError)
	ErrActivationFailed      = New("account activation failed", http.StatusInternalServerError)
	ErrFailedToDeleteUser    = New("failed to delete user", http.StatusInternalServerError)
	ErrFailedToDeleteAccount = New("failed to delete account", http.StatusInternalServerError)
	ErrFailedToUploadImage   = New("failed to upload image", http.StatusInternalServerError)
	ErrFailedToUpdateProfile = New("failed to update profile", http.StatusInternalServerError)

	// Oauth-specific Errors
	ErrFailedToFindOrCreateUser = New("failed to find or create user", http.StatusInternalServerError)
	ErrFailedToGenerateToken    = New("failed to generate token", http.StatusInternalServerError)

	// Blog-specific Errors
	ErrFailedToCreateBlog        = New("failed to create blog", http.StatusInternalServerError)
	ErrBlogNotFound              = New("blog not found", http.StatusNotFound)
	ErrPostNotFound              = New("blog post not found", http.StatusNotFound)
	ErrCommentNotFound           = New("comment not found", http.StatusNotFound)
	ErrFailedToDeleteBlog        = New("failed to delete blog", http.StatusInternalServerError)
	ErrFailedToUpdateBlog        = New("failed to update blog", http.StatusInternalServerError)
	ErrFailedToRetrieveBlogs     = New("failed to retrieve blogs", http.StatusInternalServerError)
	ErrFailedToRetrieveUserBlogs = New("failed to retrieve user blogs", http.StatusInternalServerError)
	ErrReplyNotFound             = New("reply not found", http.StatusNotFound)
	ErrLikeAlreadyExists         = New("like already exists", http.StatusConflict)
	ErrLikeNotFound              = New("like not found", http.StatusNotFound)
	ErrUserAlreadyExists         = New("user already exists", http.StatusConflict)
	ErrInvalidCredentials        = New("invalid credentials", http.StatusUnauthorized)
	ErrInsufficientRights        = New("insufficient rights", http.StatusForbidden)
	ErrAdminRoleRequired         = New("admin role required", http.StatusForbidden)
	ErrUserRoleRequired          = New("user role required", http.StatusForbidden)

	// Comment-specific Errors
	ErrFailedToCreateComment       = New("failed to create comment", http.StatusInternalServerError)
	ErrFailedToUpdateComment       = New("failed to update comment", http.StatusInternalServerError)
	ErrFailedToDeleteComment       = New("failed to delete comment", http.StatusInternalServerError)
	ErrFailedToGetComment          = New("failed to get comment", http.StatusInternalServerError)
	ErrInvalidCommentID            = New("invalid comment ID", http.StatusBadRequest)
	ErrInvalidPaginationParameters = New("invalid pagination parameters", http.StatusBadRequest)
	ErrFailedToCreateReply         = New("failed to create reply", http.StatusInternalServerError)
	ErrFailedToUpdateReply         = New("failed to update reply", http.StatusInternalServerError)
	ErrFailedToDeleteReply         = New("failed to delete reply", http.StatusInternalServerError)
	ErrFailedToGetReplies          = New("failed to get replies", http.StatusInternalServerError)
	ErrInvalidReplyID              = New("invalid reply ID", http.StatusBadRequest)
	ErrFailedToLikeComment         = New("failed to like comment", http.StatusInternalServerError)
	ErrFailedToUnlikeComment       = New("failed to unlike comment", http.StatusInternalServerError)
	ErrFailedToLikeReply           = New("failed to like reply", http.StatusInternalServerError)
	ErrFailedToUnlikeReply         = New("failed to unlike reply", http.StatusInternalServerError)
	ErrFailedToGetComments         = New("failed to get comments", http.StatusInternalServerError)

	// Like-specific Errors
	ErrFailedToLikePost   = New("failed to like post", http.StatusInternalServerError)
	ErrFailedToUnlikePost = New("failed to unlike post", http.StatusInternalServerError)
)






type RefreshToken struct {
	Token     string    `bson:"token" json:"token"`
	DeviceID  string    `bson:"device_id" json:"device_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type LogInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type TokenGenerator interface {
	GenerateToken(user User) (string, error)
	GenerateRefreshToken(user User) (string, error)
	RefreshToken(token string) (string, error)
}

type TokenVerifier interface {
	VerifyToken(token string) (*User, error)
	VerifyRefreshToken(token string) (*User, error)
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}



type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id" json:"userId"` // References the User ObjectID
	AccessToken  string             `bson:"access_token" json:"accessToken"`
	RefreshToken string             `bson:"refresh_token" json:"refreshToken"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expiresAt"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

type JwtCustomClaims struct {
	Authorized bool   `json:"authorized"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	Username   string `json:"username"`
	IsActivated bool `json:"is_activated"`

	jwt.StandardClaims
}