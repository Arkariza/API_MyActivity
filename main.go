package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/controller/Lead"
	"github.com/Arkariza/API_MyActivity/controller/Meet"
	"github.com/Arkariza/API_MyActivity/controller/User"
	"github.com/Arkariza/API_MyActivity/middleware/Meet"
	"github.com/Arkariza/API_MyActivity/middleware/Lead"
	"github.com/Arkariza/API_MyActivity/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()
	defer models.DisconnectDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:57038"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	authCommand := auth.NewAuthCommand(models.GetCollection("users"))
	userController := UserControllers.NewUserController(authCommand)
	leadController := LeadControllers.NewLeadController(models.GetCollection("leads"))
	meetController := MeetControllers.NewMeetController(models.GetCollection("meet"))

	authMiddleware := LeadMiddleware.NewAuthMiddleware(authCommand.GetSecretKey())
	leadMiddleware := LeadMiddleware.SetLeadStatus()
	meetMiddleware := MeetMiddleware.NewMeetMiddleware(authCommand.GetSecretKey())


	api := r.Group("/api")
	{

		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		leads := api.Group("/leads")
		leads.Use(authMiddleware.AuthenticateUser())
		leads.Use(leadMiddleware)
		{
			leads.POST("/add", LeadMiddleware.ValidateLeadInput(), leadController.AddLead)
			leads.GET("/", leadController.GetLeads)
			leads.GET("/:id", leadController.GetLeadByID)
			leads.POST("/referral", LeadMiddleware.ValidateLeadInput(), leadController.AddReferral)
		}

		meets := api.Group("/meets")
		meets.Use(meetMiddleware.AuthenticateMeet())
		{meets.POST("/add", MeetMiddleware.ValidateMeetRequest(), func(c *gin.Context) {
			var req MeetControllers.AddMeetRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		
			meet, err := meetController.AddMeet(c, req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		
			c.JSON(http.StatusCreated, gin.H{
				"message": "Meet created successfully",
				"data":    meet,
			})
		})
			meets.GET("/", meetController.ViewMeets)
			meets.GET("/:id", meetController.GetMeetByID)
			meets.DELETE("/:id", meetController.DeleteMeet)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}