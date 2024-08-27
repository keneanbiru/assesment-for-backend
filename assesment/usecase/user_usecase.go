package usecase

import (
	domain "assesment/domain"
	Infrastructure "assesment/Infrastructure"
	"context"
	"time"
	"errors"
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
		userRepository:    userRepo,
		TokenGen:    tokenGen,
		PasswordSvc: passwordSvc,
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

	//Check if a user with the same username or email already exists
	// _, err := u.userRepository.GetUserByUsernameOrEmail(user.Username, user.Email)
	// if err == nil {
	// 	return domain.ErrUserAlreadyExists
	// } else if err != domain.ErrUserNotFound {
	// 	return domain.ErrInternalServer
	// }

	user.Role = "user"

	// Hash the password
	// hashedPassword, err := u.PasswordSvc.HashPassword(user.Password)
	// if err != nil {
	// 	return domain.ErrInternalServer
	// }

	// Generate an activation token
	token, err := Infrastructure.GenerateActivationToken()
	if err != nil {
		return domain.ErrInternalServer
	}

	//user.Password = hashedPassword
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
func (useru *userUsecase) LoginUser(c context.Context, user domain.User) (int, string, error) {
	// Delegate the login to the UserRepository
	return useru.userRepository.LoginUser(user)
}

// DeleteUser handles the deletion of a user by delegating to the repository.
func (useru *userUsecase) DeleteUser(c context.Context, id string) (int, error) {
	// Delegate the deletion to the UserRepository
	return useru.userRepository.DeleteUser(id)
}
