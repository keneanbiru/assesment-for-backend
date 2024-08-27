package Infrastructure


import (
	//"errors"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
)








// Function to validate email format
func IsValidEmail(email string) bool {
	// Simple regex for validating email format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Function to validate password strength
func IsValidPassword(password string) bool {
	// Example: password must be at least 8 characters long
	return len(password) >= 8
}

// Generate a new activation token
func GenerateActivationToken() (string, error) {
	// Example: Generate a random token (for demonstration purposes)
	// In a real application, use a secure method to generate tokens
	return "random-activation-token", nil
}

// Send activation email
func SendActivationEmail(email, token string) error {
	// Placeholder for sending email
	// Use an actual email service to send the activation email
	return nil
}

// Generate a new password reset token
func GeneratePasswordResetToken() (string, error) {
	// Example: Generate a random token (for demonstration purposes)
	// In a real application, use a secure method to generate tokens
	return "random-reset-token", nil
}

// Send password reset email
func SendPasswordResetEmail(email, token string) error {
	// Placeholder for sending email
	// Use an actual email service to send the reset email
	return nil
}

// Token claims structure
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Generate JWT token
func GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

// Refresh JWT token
func RefreshToken(oldToken string) (string, error) {
	// Example: parse and validate the old token, then generate a new one
	// For simplicity, we'll just generate a new token without validating the old one
	return GenerateToken("user_id_from_old_token")
}

// Verify and parse JWT token
func VerifyResetToken(resetToken string) (*Claims, error) {
	// Example: parse and verify the reset token
	// For simplicity, we'll return a dummy user ID
	return &Claims{UserID: "user_id_from_reset_token"}, nil
}
