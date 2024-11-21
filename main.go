package main

import (
	"log"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/controller/Call"
	"github.com/Arkariza/API_MyActivity/controller/User"
	"github.com/Arkariza/API_MyActivity/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
    models.ConnectDatabase()
    defer models.DisconnectDatabase()

    r := gin.Default()

    authCommand := auth.NewAuthCommand(models.GetCollection("users"))
    userController := UserControllers.NewUserController(authCommand)
    callController := CallControllers.NewCallController(&gorm.DB{})


    api := r.Group("/api")  
    {
        api.POST("/register", userController.Register)
        api.POST("/login", userController.Login)
        api.POST("/addcall", callController.AddCall)
    }

    if err := r.Run(":8080"); err != nil {
        log.Fatal("Error starting server:", err)
    }
}