package domain

import "time"

// Loan represents the loan entity in the domain layer.
type Loan struct {
    ID        uint      `json:"id"`
    UserID    uint      `json:"user_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`   // possible values: "pending", "approved", "rejected"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// LoanRepository provides an interface for loan-related operations in the repository layer.
type LoanRepository interface {
    ApplyForLoan(loan Loan) error               // Method to apply for a loan
    GetLoanByID(id uint) (Loan, error)          // Method to retrieve a loan by its ID
    GetAllLoans(status, order string) ([]Loan, error) // Method to retrieve all loans with optional filtering by status and sorting
    UpdateLoanStatus(id uint, status string) error // Method to update the status of a loan (approve/reject)
    DeleteLoan(id uint) error                   // Method to delete a loan by its ID
}

// LoanUsecase provides an interface for loan-related business logic in the use case layer.
type LoanUsecase interface {
    ApplyForLoan(loan Loan) error               // Method to apply for a loan
    GetLoanByID(id uint) (Loan, error)          // Method to retrieve a loan by its ID
    GetAllLoans(status, order string) ([]Loan, error) // Method to retrieve all loans with optional filtering by status and sorting
    ApproveLoan(id uint) error                  // Method to approve a loan
    RejectLoan(id uint) error                   // Method to reject a loan
    DeleteLoan(id uint) error                   // Method to delete a loan by its ID
}
