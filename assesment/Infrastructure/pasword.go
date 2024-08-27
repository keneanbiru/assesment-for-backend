package Infrastructure

import (
	"golang.org/x/crypto/bcrypt"
    "errors"
	"log"
)

// PasswordServiceImpl is a struct that implements the PasswordService interface.
type PasswordServiceImpl struct{}

// NewPasswordService creates a new instance of PasswordServiceImpl.
func NewPasswordService() *PasswordServiceImpl {
    return &PasswordServiceImpl{}
}

// HashPassword hashes the given password using bcrypt.
// func (p *PasswordServiceImpl) PasswordHasher(password string) (string, error) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//     return string(bytes), err
// }





func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// PasswordComparator compares a plain text password with a hashed password.
func PasswordComparator(hashedPassword, plainPassword string) bool {
	// Compare the hashed password with the plain password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
// CheckPasswordHash compares the given password with the stored hash.
func (p *PasswordServiceImpl) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}


func (ps *PasswordServiceImpl) HashPassword(password string) (string, error) {
	// Generate a hash from the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if the provided password matches the hashed password.


// PasswordHasher is similar to HashPassword; included here as per the interface.
func (ps *PasswordServiceImpl) PasswordHasher(password string) (string, error) {
	return ps.HashPassword(password)
}