package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "assesment/domain"
    //"assesment/usecase"
)

// LoanController handles HTTP requests related to loans.
type LoanController struct {
    loanUsecase domain.LoanUsecase
}

// NewLoanController creates a new instance of LoanController.
func NewLoanController(loanUsecase domain.LoanUsecase) *LoanController {
    return &LoanController{
        loanUsecase: loanUsecase,
    }
}

// ApplyForLoan handles the request to apply for a loan.
func (lc *LoanController) ApplyForLoan(c *gin.Context) {
    var loan domain.Loan
    if err := c.ShouldBindJSON(&loan); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := lc.loanUsecase.ApplyForLoan(loan); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan application submitted successfully"})
}

// GetLoanByID handles the request to retrieve a loan by ID.
func (lc *LoanController) GetLoanByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    loan, err := lc.loanUsecase.GetLoanByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, loan)
}

// GetAllLoans handles the request to retrieve all loans.
func (lc *LoanController) GetAllLoans(c *gin.Context) {
    status := c.Query("status")
    order := c.Query("order")

    loans, err := lc.loanUsecase.GetAllLoans(status, order)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, loans)
}

// ApproveLoan handles the request to approve a loan.
func (lc *LoanController) ApproveLoan(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := lc.loanUsecase.ApproveLoan(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan approved successfully"})
}

// RejectLoan handles the request to reject a loan.
func (lc *LoanController) RejectLoan(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := lc.loanUsecase.RejectLoan(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan rejected successfully"})
}

// DeleteLoan handles the request to delete a loan.
func (lc *LoanController) DeleteLoan(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := lc.loanUsecase.DeleteLoan(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
