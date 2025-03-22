package UserControllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/models/User"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	authCommand *auth.AuthCommand
	collection  *mongo.Collection
}

func NewUserController(authCommand *auth.AuthCommand, collection *mongo.Collection) *UserController {
	return &UserController{
		authCommand: authCommand,
		collection:  collection,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	PhoneNum string `json:"phone_num" binding:"required"`
	Role     int    `json:"role" binding:"required,oneof=1 2"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	cmdRequest := auth.RegisterRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		PhoneNum: request.PhoneNum,
		Role:     request.Role,
	}
	ctxRequest := context.Background()
	user, err := c.authCommand.Register(ctxRequest, cmdRequest)
	if err != nil {
		log.Printf("Registration failed: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Registration failed",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Registration successful",
		"data":    user,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	cmdRequest := auth.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	}
	ctxRequest := context.Background()
	tokenResponse, err := c.authCommand.Login(ctxRequest, cmdRequest)
	if err != nil {
		log.Printf("Login failed: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login successful",
		"data":    tokenResponse,
	})
}

func (uc *UserController) GetUser(c *gin.Context) {
	cursor, err := uc.collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Failed to fetch users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Printf("Failed to decode users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	totalCount, err := uc.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Failed to count users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": totalCount,
	})
}



func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %s\n", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid ID format",
		})
		return
	}

	var user models.User
	err = uc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found for ID: %s\n", id)
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "User not found",
			})
		} else {
			log.Printf("Failed to retrieve user: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to retrieve user",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func (uc *UserController) EditProfile(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %s\n", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid ID format",
		})
		return
	}
	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("Failed to parse request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	updateData.ID = primitive.NilObjectID
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"username":   updateData.Username,
			"email":      updateData.Email,
			"phone_num":  updateData.PhoneNum,
			"image":      updateData.Image,
			"updated_at": time.Now(),
		},
	}

	result, err := uc.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	if result.MatchedCount == 0 {
		log.Printf("User not found for ID: %s\n", id)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Profile updated successfully",
	})
}