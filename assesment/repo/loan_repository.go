package repository

import (
    "assesment/domain"
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoanRepository struct {
    database   *mongo.Database
    collection *mongo.Collection
}

// NewLoanRepository creates a new instance of LoanRepository.
func NewLoanRepository(mongoClient *mongo.Client) domain.LoanRepository {
    return &LoanRepository{
        database:   mongoClient.Database("loan"),
        collection: mongoClient.Database("loan").Collection("loans"),
    }
}

// ApplyForLoan inserts a new loan application into the MongoDB collection.
func (r *LoanRepository) ApplyForLoan(loan domain.Loan) error {
    loan.CreatedAt = time.Now()
    loan.UpdatedAt = time.Now()
    _, err := r.collection.InsertOne(context.Background(), loan)
    return err
}

// GetLoanByID retrieves a loan by its ID from the MongoDB collection.
func (r *LoanRepository) GetLoanByID(id uint) (domain.Loan, error) {
    var loan domain.Loan
    objectID := primitive.NewObjectIDFromTimestamp(time.Unix(int64(id), 0)) // assuming the id is a Unix timestamp

    err := r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&loan)
    return loan, err
}

// GetAllLoans retrieves all loans from the MongoDB collection, optionally filtering by status and sorting.
// GetAllLoans retrieves all loans from the MongoDB collection, optionally filtering by status and sorting.
func (r *LoanRepository) GetAllLoans(status, order string) ([]domain.Loan, error) {
    var loans []domain.Loan
    filter := bson.M{}

    // Apply status filter if provided
    if status != "" && status != "all" {
        filter["status"] = status
    }

    // Set the sorting order
    sortOrder := 1
    if order == "desc" {
        sortOrder = -1
    }

    // Define options with sort order
    findOptions := options.Find().SetSort(bson.M{"created_at": sortOrder})

    // Query the collection with filter and options
    cursor, err := r.collection.Find(context.Background(), filter, findOptions)
    if err != nil {
        return nil, err
    }

    // Decode all results into the loans slice
    if err := cursor.All(context.Background(), &loans); err != nil {
        return nil, err
    }

    return loans, nil
}
// UpdateLoanStatus updates the status of a loan in the MongoDB collection.
func (r *LoanRepository) UpdateLoanStatus(id uint, status string) error {
    objectID := primitive.NewObjectIDFromTimestamp(time.Unix(int64(id), 0)) // assuming the id is a Unix timestamp

    _, err := r.collection.UpdateOne(
        context.Background(),
        bson.M{"_id": objectID},
        bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
    )
    return err
}

// DeleteLoan deletes a loan by its ID from the MongoDB collection.
func (r *LoanRepository) DeleteLoan(id uint) error {
    objectID := primitive.NewObjectIDFromTimestamp(time.Unix(int64(id), 0)) // assuming the id is a Unix timestamp

    _, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
    return err
}
