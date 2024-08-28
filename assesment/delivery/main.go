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
	secretkey := "abebe"
	tokenGen := infrastructure.NewTokenGeneratorImpl(secretkey,user_repo)
	passwordService := infrastructure.NewPasswordService()

	//set-up the controllers
	cont := controllers.NewUserController(usecase.NewUserUsecase(user_repo,tokenGen,passwordService))

	//the router gateway
	router := gin.Default()
	routers.SetupRoutes(router, cont)
	router.Run(":8080")
}
