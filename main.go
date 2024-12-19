package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/controller/Call"
	"github.com/Arkariza/API_MyActivity/controller/Comment"
	"github.com/Arkariza/API_MyActivity/controller/Lead"
	"github.com/Arkariza/API_MyActivity/controller/Meet"
	"github.com/Arkariza/API_MyActivity/controller/User"
	"github.com/Arkariza/API_MyActivity/middleware/Call"
	"github.com/Arkariza/API_MyActivity/middleware/Comment"
	"github.com/Arkariza/API_MyActivity/middleware/Lead"
	"github.com/Arkariza/API_MyActivity/middleware/Meet"
	"github.com/Arkariza/API_MyActivity/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()
	defer models.DisconnectDatabase()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:50574"},
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
	callController := CallControllers.NewCallController(models.GetCollection("call"))
	commentController := CommentController.NewCommentController(models.GetCollection("comments"))

	leadMiddleware := middleware.NewLeadMiddleware(authCommand.GetSecretKey())
	meetMiddleware := MeetMiddleware.NewMeetMiddleware(authCommand.GetSecretKey())
	callMiddleware := CallMiddleware.NewCallMiddleware(authCommand.GetSecretKey())
	commentMiddleware := CommentMiddleware.NewCommentMiddleware(authCommand.GetSecretKey())

	api := r.Group("/api")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		leads := api.Group("/leads")
		leads.Use(leadMiddleware.AuthenticateLead())
		{
			leads.POST("/add", leadMiddleware.AuthenticateLead(), func(c *gin.Context) {
				var req LeadController.AddLeadRequest
			
				lead, err := leadController.AddLead(c, req)
				if err != nil {
					return 
				}
			
				c.JSON(http.StatusCreated, gin.H{
					"message": "Lead has been created",
					"data":    lead,
				})
			})
			leads.GET("/", leadController.GetAllLead)
		}

		meets := api.Group("/meets")
		meets.Use(meetMiddleware.AuthenticateMeet())
		{
			meets.POST("/add", func(c *gin.Context) {
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
		
		calls := api.Group("/calls")
		calls.Use(callMiddleware.AuthenticateCall())
		{
			calls.POST("/add", func(c *gin.Context) {
				var req CallControllers.AddCallRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Invalid request",
						"details": err.Error(),
					})
					return
				}

				call, err := callController.AddCall(c, req)
				if err != nil {
					log.Printf("Error adding call: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to create call",
						"details": err.Error(),
					})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"message": "Call has been created",
					"data":    call,
				})
			})
			calls.GET("/", callController.GetCalls)
			calls.GET("/:id", callController.GetCallByID)
		}

		comments := api.Group("/comments")
		comments.Use(commentMiddleware.AuthenticateComment())
		{
			comments.POST("/add", commentController.CreateComment)
			comments.GET("/", commentController.GetAllComments)
			comments.GET("/:id", commentController.GetCommentByID)
			comments.PUT("/:id", commentController.UpdateComment)
			comments.DELETE("/:id", commentController.DeleteComment)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
