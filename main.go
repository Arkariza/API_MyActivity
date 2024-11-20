package main

import (
	"log"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/controller/Lead"
	"github.com/Arkariza/API_MyActivity/controller/User"
	"github.com/Arkariza/API_MyActivity/middleware/Lead"
	"github.com/Arkariza/API_MyActivity/models"
	"github.com/gin-gonic/gin"
)

func main() {
    models.ConnectDatabase()
    defer models.DisconnectDatabase()

    r := gin.Default()

    authCommand := auth.NewAuthCommand(models.GetCollection("users"))
    userController := UserControllers.NewUserController(authCommand)
    leadController := LeadControllers.NewLeadController(models.GetCollection("leads"))
    authMiddleware := LeadMiddleware.NewAuthMiddleware(authCommand.GetSecretKey())
    
    api := r.Group("/api")
    {
        api.POST("/register", userController.Register)
        api.POST("/login", userController.Login)

        leads := api.Group("/leads")
        leads.Use(authMiddleware.AuthenticateUser())
        leads.Use(LeadMiddleware.SetLeadStatus())
        {
            leads.POST("/add", LeadMiddleware.ValidateLeadInput(), leadController.AddLead)
            leads.GET("/", leadController.GetLeads)
            leads.GET("/:id", leadController.GetLeadByID)      
            leads.POST("/referral", LeadMiddleware.ValidateLeadInput(), leadController.AddReferral)
        }
    }

    if err := r.Run(":8080"); err != nil {
        log.Fatal("Error starting server:", err)
    }
}