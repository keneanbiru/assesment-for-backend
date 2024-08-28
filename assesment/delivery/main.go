package main

import (
	"assesment/delivery/controllers"
	"assesment/delivery/routes"
	infrastructure"assesment/infrastructure"
	repositories"assesment/repo"
	"assesment/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	client := infrastructure.MongoDBInit() // MongoDB initialization

	// Initialize the repositories
	userRepo := repositories.NewUserRepository(client)
	loanRepo := repositories.NewLoanRepository(client)

	// Set up the token generator, password service, and use cases
	secretKey := "abebe"
	tokenGen := infrastructure.NewTokenGeneratorImpl(secretKey, userRepo)
	passwordService := infrastructure.NewPasswordService()

	// Set up the controllers
	userCtrl := controllers.NewUserController(usecase.NewUserUsecase(userRepo, tokenGen, passwordService))
	loanCtrl := controllers.NewLoanController(usecase.NewLoanUsecase(loanRepo))

	// Set up the router
	router := gin.Default()
	routes.SetupRoutes(router, userCtrl, loanCtrl)
	router.Run(":8080")
}
