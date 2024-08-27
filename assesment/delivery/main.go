package main

import (
	infrastructure "assesment/Infrastructure"
	"assesment/delivery/controllers"
	routers "assesment/delivery/routes"
	repositories "assesment/repo"
	usecase "assesment/usecase"

	//"time"

	"github.com/gin-gonic/gin"
)

func main() {

	client := infrastructure.MongoDBInit() //mongodb initialization

	//initialization of the repositories

	user_repo := repositories.NewUserRepository(client)
	tokenGen := infrastructure.NewTokenGenerator()
	passwordService := infrastructure.NewPasswordService()

	//set-up the controllers
	cont := controllers.NewLoanController(usecase.NewUserUsecase(user_repo, tokenGen, passwordService))

	//the router gateway
	router := gin.Default()
	routers.SetupRoutes(router, cont)
	router.Run(":8080")
}
