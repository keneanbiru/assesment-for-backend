package usecase

import (
    "assesment/domain"
    "errors"
    "time"
)

type loanUsecase struct {
    loanRepo domain.LoanRepository
}

// NewLoanUsecase creates a new instance of LoanUsecase.
func NewLoanUsecase(loanRepo domain.LoanRepository) domain.LoanUsecase {
    return &loanUsecase{
        loanRepo: loanRepo,
    }
}

// ApplyForLoan allows a user to submit a loan application.
func (uc *loanUsecase) ApplyForLoan(loan domain.Loan) error {
    // Add business logic here (e.g., validation, setting default values)
    if loan.Amount <= 0 {
        return errors.New("invalid loan amount")
    }

    loan.Status = "pending"
    loan.CreatedAt = time.Now()
    loan.UpdatedAt = time.Now()

    return uc.loanRepo.ApplyForLoan(loan)
}

// GetLoanByID retrieves the loan status by ID.
func (uc *loanUsecase) GetLoanByID(id uint) (domain.Loan, error) {
    loan, err := uc.loanRepo.GetLoanByID(id)
    if err != nil {
        return domain.Loan{}, err
    }

    return loan, nil
}

// GetAllLoans retrieves all loan applications, with optional filtering and sorting.
func (uc *loanUsecase) GetAllLoans(status, order string) ([]domain.Loan, error) {
    return uc.loanRepo.GetAllLoans(status, order)
}

// ApproveLoan allows an admin to approve a loan.
func (uc *loanUsecase) ApproveLoan(id uint) error {
    // Business logic: Validate the loan approval
    loan, err := uc.loanRepo.GetLoanByID(id)
    if err != nil {
        return err
    }

    if loan.Status != "pending" {
        return errors.New("only pending loans can be approved")
    }

    return uc.loanRepo.UpdateLoanStatus(id, "approved")
}

// RejectLoan allows an admin to reject a loan.
func (uc *loanUsecase) RejectLoan(id uint) error {
    // Business logic: Validate the loan rejection
    loan, err := uc.loanRepo.GetLoanByID(id)
    if err != nil {
        return err
    }

    if loan.Status != "pending" {
        return errors.New("only pending loans can be rejected")
    }

    return uc.loanRepo.UpdateLoanStatus(id, "rejected")
}

// DeleteLoan allows an admin to delete a loan by its ID.
func (uc *loanUsecase) DeleteLoan(id uint) error {
    return uc.loanRepo.DeleteLoan(id)
}
