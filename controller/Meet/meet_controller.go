package MeetControllers

import (
	"context"
	"net/http"
	"strconv"
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

func (mc *MeetController) AddMeet(c *gin.Context) {
	meetData, exists := c.Get("meet_data")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid meet data",
		})
		return
	}

	meet := meetData.(models.Meet)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := mc.collection.InsertOne(ctx, meet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create meet",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Meet created successfully",
		"id":      result.InsertedID,
	})
}

func (mc *MeetController) ViewMeets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	status := c.Query("status")
	clientName := c.Query("client_name")

	filter := bson.M{
		"user_id": userID,
	}

	if status != "" {
		filter["prospect_status"] = status
	}
	if clientName != "" {
		filter["client_name"] = bson.M{"$regex": clientName, "$options": "i"}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := mc.collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve meets",
		})
		return
	}
	defer cursor.Close(ctx)

	var meets []models.Meet
	if err = cursor.All(ctx, &meets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse meets",
		})
		return
	}

	total, err := mc.collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count meets",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid meet ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	var updatedMeet models.Meet
	if err := c.ShouldBindJSON(&updatedMeet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := updatedMeet.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	filter := bson.M{
		"_id":      objectID,
		"user_id":  userID,
	}

	update := bson.M{
		"$set": bson.M{
			"client_name":      updatedMeet.ClientName,
			"address":          updatedMeet.Address,
			"prospect_status":  updatedMeet.ProspectStatus,
			"latitude":         updatedMeet.Latitude,
			"longitude":        updatedMeet.Longitude,
			"meet_result":      updatedMeet.MeetResult,
		},
	}

	result, err := mc.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update meet",
		})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Meet not found or you don't have permission to update",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meet updated successfully",
	})
}

func (mc *MeetController) DeleteMeet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meetID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(meetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid meet ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	result, err := mc.collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete meet",
		})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Meet not found or you don't have permission to delete",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meet deleted successfully",
	})
}

func (mc *MeetController) GetMeetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meetID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(meetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid meet ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	var meet models.Meet
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	err = mc.collection.FindOne(ctx, filter).Decode(&meet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Meet not found or you don't have permission to view",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve meet",
		})
		return
	}

	c.JSON(http.StatusOK, meet)
}