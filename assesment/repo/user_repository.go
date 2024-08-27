package repository

import (
	"context"
	"errors"
	"log"
	domain "assesment/domain"
	infrastructure "assesment/Infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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
		collection: mongoClient.Database("loan").Collection("users"),
	}
}

// GetUserByID retrieves a user by ID.
func (ur *UserRepository) GetUserByID(id primitive.ObjectID) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// UpdateUser updates a user's data.
func (ur *UserRepository) UpdateUser(user domain.User) error {
	_, err := ur.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// RegisterUserDb registers a new user in the database.
func (userepo *UserRepository) Register(user domain.User) error {
	collection := userepo.collection

	if collection == nil {
		return errors.New("database collection is not initialized")
	}

	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Err()
	if err == nil {
		return domain.ErrUserAlreadyExists
	}

	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	user.ID = primitive.NewObjectID()

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

	collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)

	log.Println(existingUser, user)

	if !infrastructure.PasswordComparator(existingUser.Password, user.Password) {
		return http.StatusUnauthorized, "", errors.New("invalid email or password")
	}

	jwtToken, err := infrastructure.TokenGenerator(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return http.StatusInternalServerError, "", errors.New("internal server error")
	}

	return http.StatusOK, jwtToken, nil
}

// DeleteUser removes a user by ID.
func (userepo *UserRepository) DeleteUser(id primitive.ObjectID) error {
	collection := userepo.collection

	
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return  nil
}

// GetUserByEmail retrieves a user by email.
func (ur *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// UpdateUserPassword updates the user's password in the database.
func (ur *UserRepository) UpdateUserPassword(user domain.User) error {
	_, err := ur.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"password": user.Password}},
	)
	return err
}
