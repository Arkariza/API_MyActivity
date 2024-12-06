package MeetControllers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Arkariza/API_MyActivity/models/CallAndMeet"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MeetController struct {
	collection *mongo.Collection
}

func NewMeetController(collection *mongo.Collection) *MeetController {
	return &MeetController{
		collection: collection,
	}
}

type AddMeetRequest struct {
	ClientName string  `json:"client_name" binding:"required,min=2,max=100"`
	PhoneNum   string  `json:"phone_num" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	Note       string  `json:"note"`
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return page, limit
}

func handleError(c *gin.Context, statusCode int, message string, details error) {
	c.JSON(statusCode, gin.H{
		"error":   message,
		"details": details.Error(),
	})
}

func validateToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid token format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	return tokenString, nil
}

func (mc *MeetController) AddMeet(c *gin.Context, req AddMeetRequest) (*models.Meet, error) {
	token, err := validateToken(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, "Invalid token", err)
		return nil, err
	}

	fmt.Println("Validated Token:", token) // Debugging

	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found")
	}

	currentUser := user.(*models.User)

	meet := models.Meet{
		UserID:         currentUser.ID,
		ClientName:     req.ClientName,
		PhoneNum:       req.PhoneNum,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Note:           req.Note,
		Address:        req.Address,
		CreatedAt:      time.Now(),
		ProspectStatus: "potential",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := mc.collection.InsertOne(ctx, meet)
	if err != nil {
		return nil, fmt.Errorf("failed to create meet: %v", err)
	}

	meet.ID = result.InsertedID.(primitive.ObjectID)

	return &meet, nil
}

// ViewMeets function
func (mc *MeetController) ViewMeets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	status := c.Query("status")
	clientName := c.Query("client_name")

	filter := bson.M{"user_id": userID}
	if status != "" {
		filter["prospect_status"] = status
	}
	if clientName != "" {
		filter["client_name"] = bson.M{"$regex": clientName, "$options": "i"}
	}

	page, limit := parsePagination(c)
	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := mc.collection.Find(ctx, filter, findOptions)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve meets", err)
		return
	}
	defer cursor.Close(ctx)

	var meets []models.Meet
	if err = cursor.All(ctx, &meets); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to parse meets", err)
		return
	}

	total, err := mc.collection.CountDocuments(ctx, filter)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to count meets", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meets": meets,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": limit,
		},
	})
}

func (mc *MeetController) UpdateMeet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meetID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(meetID)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid meet ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var updatedMeet models.Meet
	if err := c.ShouldBindJSON(&updatedMeet); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	filter := bson.M{"_id": objectID, "user_id": userID}
	update := bson.M{"$set": updatedMeet}

	result, err := mc.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update meet", err)
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meet not found or no permission to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meet updated successfully"})
}

func (mc *MeetController) DeleteMeet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meetID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(meetID)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid meet ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	filter := bson.M{"_id": objectID, "user_id": userID}

	result, err := mc.collection.DeleteOne(ctx, filter)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete meet", err)
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meet not found or no permission to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meet deleted successfully"})
}

func (mc *MeetController) GetMeetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meetID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(meetID)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid meet ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var meet models.Meet
	filter := bson.M{"_id": objectID, "user_id": userID}

	err = mc.collection.FindOne(ctx, filter).Decode(&meet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meet not found or no permission to view"})
			return
		}
		handleError(c, http.StatusInternalServerError, "Failed to retrieve meet", err)
		return
	}

	c.JSON(http.StatusOK, meet)
}
