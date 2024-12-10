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
		AllowOrigins:     []string{"http://localhost:65516"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authCommand := auth.NewAuthCommand(models.GetCollection("users"))
	userController := UserControllers.NewUserController(authCommand)
	leadController := LeadController.NewLeadController(models.GetCollection("leads"))
	meetController := MeetControllers.NewMeetController(models.GetCollection("meet"))

	leadMiddleware := LeadMiddleware.NewLeadMiddleware(authCommand.GetSecretKey())
	meetMiddleware := MeetMiddleware.NewMeetMiddleware(authCommand.GetSecretKey())

	api := r.Group("/api")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		leads := api.Group("/leads")
		leads.Use(leadMiddleware.AuthenticateLead())
		{
			leads.POST("/auto-status", leadController.CreateLeadWithAutoStatus)
		}

		meets := api.Group("/meets")
		meets.Use(meetMiddleware.AuthenticateMeet())
		{
			meets.POST("/add", meetMiddleware.AuthenticateMeet(), func(c *gin.Context) {
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
