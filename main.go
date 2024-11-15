package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/Arkariza/API_MyActivity/models"
    "github.com/Arkariza/API_MyActivity/controller/User"
    "github.com/Arkariza/API_MyActivity/auth"
)

func main() {
    models.ConnectDatabase()
    defer models.DisconnectDatabase()

    r := gin.Default()

    authCommand := auth.NewAuthCommand(models.GetCollection("users"))
    userController := UserControllers.NewUserController(authCommand)

    api := r.Group("/api")
    {
        api.POST("/register", userController.Register)
        api.POST("/login", userController.Login)
    }

    if err := r.Run(":8080"); err != nil {
        log.Fatal("Error starting server:", err)
    }
}