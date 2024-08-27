package Infrastructure

import (
	"log"
	"strings"

	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The middleware for Authentication
func AuthMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization") // extracting the authentication value from the header
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		c.JSON(401, gin.H{"error": "Invalid authorization header"})
		c.Abort()
		return
	}

	tokenString := authParts[1]             //the token string
	token, err := TokenClaimer(tokenString) //verifying the token

	if err != nil {
		log.Println("Token parsing error:", err.Error())
		c.JSON(401, gin.H{"error": "Invalid JWT"})
		c.Abort()
		return
	}

	//extracting the map claims from the token
	if role, err := token.Claims.(jwt.MapClaims); err && token.Valid {
		log.Println(role["_id"])
		c.Set("isadmin", role["isadmin"])
		c.Set("userid", role["_id"])
	} else {
		c.JSON(401, gin.H{"error": "Invalid JWT"})
		c.Abort()
		return
	}

	c.Next()
}

// The middleware for Authentication
func AdminMiddleware(c *gin.Context) {
	isAdmin, exists := c.Get("isadmin") //fetching the data from the context
	if !exists || !isAdmin.(bool) {
		c.JSON(403, gin.H{"error": "Forbidden: You don't have admin privileges"})
		c.Abort()
		return
	}

	c.Next()
}
func TokenGenerator(id primitive.ObjectID, email string, Role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SECRET_KEY := os.Getenv("JWT_SECRET")

	var jwtSecret = []byte(SECRET_KEY)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":     id.Hex(),
		"email":   email,
		"role" : Role,
	})

	return token.SignedString(jwtSecret)
}

// a token claimer for extracting the necessary datas
func TokenClaimer(tokenstr string) (*jwt.Token, error) {

	SECRET_KEY := DotEnvLoader("JWT_SECRET")

	log.Println("secretkey: ", SECRET_KEY)

	return jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SECRET_KEY), nil
	})
}
