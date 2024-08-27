package Infrastructure

import (
	"assesment/config"
	"assesment/domain"
	"time"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

//tokenGenerator implementation
type TokenGeneratorImpl struct {
	secretKey string
	userRepo  domain.UserRepository
}

// NewTokenGeneratorImpl creates a new instance of TokenGeneratorImpl.
func NewTokenGeneratorImpl(secretKey string, userRepo domain.UserRepository) domain.TokenGenerator {
	return &TokenGeneratorImpl{
		secretKey: secretKey,
		userRepo:  userRepo,
	}
}

// GenerateToken generates an access token for the user
func (tg *TokenGeneratorImpl) GenerateToken(user domain.User) (string, error) {
	accessTokenSecret := []byte(config.EnvConfigs.JwtSecret)
	accessTokenExpiryHour := config.EnvConfigs.AccessTokenExpiryHour

	claims := domain.JwtCustomClaims{
		Authorized:  true,
		UserID:      user.ID.Hex(),
		Role:        user.Role,
		Username:    user.Username,
		IsActivated: user.IsActive,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(accessTokenExpiryHour)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(accessTokenSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

// GenerateRefreshToken generates a refresh token for the user
func (tg *TokenGeneratorImpl) GenerateRefreshToken(user domain.User) (string, error) {
	refreshTokenSecret := []byte(config.EnvConfigs.JwtRefreshSecret)
	refreshTokenExpiryHour := config.EnvConfigs.RefreshTokenExpiryHour

	claims := domain.JwtCustomClaims{
		Authorized:  true,
		UserID:      user.ID.Hex(),
		Role:        user.Role,
		Username:    user.Username,
		IsActivated: user.IsActive,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(refreshTokenExpiryHour)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(refreshTokenSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

// RefreshToken parses and verifies a refresh token and returns the user ID
func (t *TokenGeneratorImpl) RefreshToken(token string) (domain.User, error) {
	if token == "" {
		return domain.User{}, domain.ErrInvalidToken
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return domain.User{}, domain.ErrTokenExpired
			}
		}
		return domain.User{}, domain.ErrInvalidToken
	}

	// Extract claims and validate
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return domain.User{}, domain.ErrInvalidToken
	}

	// Extract user ID from claims
	userIDHex, ok := claims["user_id"].(string)
	if !ok {
		return domain.User{}, domain.ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return domain.User{}, domain.ErrInvalidToken
	}

	// Retrieve the user by ID (assumes a method exists in the repository)
	user, err := t.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, domain.ErrInternalServer
	}

	return user, nil
}


func (t *TokenGeneratorImpl) VerifyResetToken(token string) (*domain.User, error) {
	if token == "" {
		return nil, domain.ErrInvalidToken
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, domain.ErrTokenExpired
			}
		}
		return nil, domain.ErrInvalidToken
	}

	// Extract claims and validate
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Extract user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	// Convert userID to primitive.ObjectID if necessary
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Retrieve the user by ID
	user, err := t.userRepo.GetUserByID(oid)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrInternalServer
	}

	return &user, nil
}