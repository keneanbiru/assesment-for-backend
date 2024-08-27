package repositories

import (
	"context"
	"errors"
	"log"
	"net/http"
	domain "assesment/domain"
	infrastructure "assesment/Infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository implements the UserRepository interface for MongoDB.
type UserRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(mongoClient *mongo.Client) domain.UserRepository {
	return &UserRepository{
		database:   mongoClient.Database("loan"),
		collection: mongoClient.Database("loan").Collection("users"), // Changed collection name to "users"
	}
}


func (ur *UserRepository) GetUserByUsernameOrEmail(username, email string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(),  bson.M{"username": username, "email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// RegisterUserDb registers a new user in the database.
func (userepo *UserRepository) Register(user domain.User) error {
	collection := userepo.collection

	// Ensure the collection is not nil
	if collection == nil {
		return errors.New("database collection is not initialized")
	}

	// Check if a user with the same email already exists
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Err()
	if err == nil {
		return domain.ErrUserAlreadyExists // use a domain-specific error
	}
	// if err != mongo.ErrNoDocuments { // Check for other errors
	// 	return err
	// }

	// Hash the password before storing it
	hashedPassword, err := infrastructure.PasswordHasher(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate a new ObjectID for the user
	user.ID = primitive.NewObjectID()

	// Insert the new user into the collection
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}


// LoginUserDb authenticates a user and returns a JWT token.
func (userepo *UserRepository) LoginUser(user domain.User) (int, string, error) {
	collection := userepo.collection

	var existingUser domain.User

	// Retrieve the user with the provided email
	collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)

	log.Println(existingUser, user)

	// Check if the provided password matches the stored hashed password
	if !infrastructure.PasswordComparator(existingUser.Password, user.Password) {
		return http.StatusUnauthorized, "", errors.New("invalid email or password")
	}

	// Generate a JWT token for the authenticated user
	jwtToken, err := infrastructure.TokenGenerator(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return http.StatusInternalServerError, "", errors.New("internal server error")
	}

	return http.StatusOK, jwtToken, nil
}

// DeleteUser removes a user by ID.
func (userepo *UserRepository) DeleteUser(id string) (int, error) {
	collection := userepo.collection

	ido, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ido}

	// Delete the user from the collection
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil || result.DeletedCount == 0 {
		return http.StatusNotFound, errors.New("user not found")
	}

	return http.StatusOK, nil
}
