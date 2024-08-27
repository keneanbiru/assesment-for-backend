package usecase

import (
	"context"
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	domain "assesment/domain"
	Infrastructure "assesment/Infrastructure"
)

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	userRepository domain.UserRepository
	TokenGen       domain.TokenGenerator
	PasswordSvc    domain.PasswordService
}

// NewUserUsecase creates a new instance of userUsecase.
func NewUserUsecase(userRepo domain.UserRepository, tokenGen domain.TokenGenerator, passwordSvc domain.PasswordService) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepo,
		TokenGen:       tokenGen,
		PasswordSvc:    passwordSvc,
	}
}

// RegisterUser handles user registration by delegating to the repository.
func (u *userUsecase) Register(user domain.User) error {
	// Validate input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return domain.ErrMissingRequiredFields
	}

	if !Infrastructure.IsValidEmail(user.Email) {
		return domain.ErrInvalidEmail
	}

	if !Infrastructure.IsValidPassword(user.Password) {
		return domain.ErrInvalidPassword
	}

	// Check if a user with the same email already exists
	_, err := u.userRepository.GetUserByEmail(user.Email)
	if err == nil {
		return domain.ErrUserAlreadyExists
	} else if err != domain.ErrUserNotFound {
		return domain.ErrInternalServer
	}

	user.Role = "user"

	// Hash the password
	hashedPassword, err := u.PasswordSvc.HashPassword(user.Password)
	if err != nil {
		return domain.ErrInternalServer
	}
	user.Password = hashedPassword

	// Generate an activation token
	token, err := Infrastructure.GenerateActivationToken()
	if err != nil {
		return domain.ErrInternalServer
	}

	user.ActivationToken = token
	user.TokenCreatedAt = time.Now()

	// Register the user in the repository
	err = u.userRepository.Register(user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.ErrUserAlreadyExists
		}
		return domain.ErrInternalServer
	}

	// Send activation email
	err = Infrastructure.SendActivationEmail(user.Email, token)
	if err != nil {
		return domain.ErrFailedToSendEmail
	}

	return nil
}

// LoginUser handles user login by delegating to the repository and returns a JWT token.
func (u *userUsecase) LoginUser(c context.Context, user domain.User) (int, string, error) {
	// Delegate the login to the UserRepository
	return u.userRepository.LoginUser(user)
}

// UpdateUser handles updating a user's data.
func (u *userUsecase) UpdateUser(c context.Context, user domain.User) error {
	// Validate input
	if user.ID.IsZero() {
		return domain.ErrInvalidUserID
	}

	// Update the user in the repository
	err := u.userRepository.UpdateUser(user)
	if err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// UpdateUserPassword handles updating a user's password.
func (u *userUsecase) UpdateUserPassword(c context.Context, id primitive.ObjectID, newPassword string) error {
	// Validate input
	if id.IsZero() || newPassword == "" {
		return domain.ErrMissingRequiredFields
	}

	// Hash the new password
	hashedPassword, err := u.PasswordSvc.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternalServer
	}

	// Create a user object with updated password
	user := domain.User{
		ID:       id,
		Password: hashedPassword,
	}

	// Update the password in the repository
	err = u.userRepository.UpdateUserPassword(user)
	if err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// GetUserByID retrieves a user by ID.
func (u *userUsecase) GetUserByID(c context.Context, id primitive.ObjectID) (domain.User, error) {
	// Retrieve the user from the repository
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, domain.ErrInternalServer
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email.
func (u *userUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	// Retrieve the user from the repository
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, domain.ErrInternalServer
	}
	return user, nil
}

// ActivateUser handles user account activation.
func (u *userUsecase) ActivateUser(c context.Context, token string) error {
	// Get the user by activation token
	user, err := u.userRepository.GetUserByEmail(token)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternalServer
	}

	// Check if the token is expired
	if time.Since(user.TokenCreatedAt) > 24*time.Hour {
		return domain.ErrTokenExpired
	}

	// Activate the user
	user.IsActive = true
	user.ActivationToken = ""
	err = u.userRepository.UpdateUser(user)
	if err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// RefreshToken refreshes the user's JWT token.
func (u *userUsecase) RefreshToken(c context.Context, oldToken string) (string, error) {
	// Verify and parse the old token
	user, err := u.TokenGen.RefreshToken(oldToken)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	// Generate a new token
	newToken, err := u.TokenGen.GenerateToken(user)
	if err != nil {
		return "", domain.ErrInternalServer
	}

	return newToken, nil
}

// SendPasswordResetLink handles sending a password reset link to the user's email.
func (u *userUsecase) SendPasswordResetLink(c context.Context, email string) error {
	// Validate the email
	if !Infrastructure.IsValidEmail(email) {
		return domain.ErrInvalidEmail
	}

	// Check if the user exists
	_, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternalServer
	}

	// Generate a reset token
	resetToken, err := Infrastructure.GeneratePasswordResetToken()
	if err != nil {
		return domain.ErrInternalServer
	}

	// Send the reset link via email
	err = Infrastructure.SendPasswordResetEmail(email, resetToken)
	if err != nil {
		return domain.ErrFailedToSendEmail
	}

	return nil
}

// ResetPassword handles resetting the user's password using a reset token.
func (u *userUsecase) ResetPassword(c context.Context, resetToken string, newPassword string) error {
	// Verify the reset token and get the user
	user, err := u.TokenGen.VerifyResetToken(resetToken)
	if err != nil {
		return domain.ErrInvalidToken
	}

	// Hash the new password
	hashedPassword, err := u.PasswordSvc.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternalServer
	}

	// Update the password in the repository
	user.Password = hashedPassword
	err = u.userRepository.UpdateUserPassword(*user)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	// You might want to add additional logic here, such as checking if the user exists
	// before attempting to delete, logging the deletion, etc.

	// Call the repository method to delete the user
	return u.userRepository.DeleteUser( id)
}
