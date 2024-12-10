package CallControllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CallController struct {
	collection *mongo.Collection
}

func NewCallController(collection *mongo.Collection) *CallController {
	return &CallController{
		collection: collection,  
	}
}

type AddCallRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	NumPhone   string `json:"numphone" binding:"required"`
	Note       string `json:"note"`
}

type UpdateCallRequest struct {
	ClientName string `json:"client_name"`
	NumPhone   string `json:"numphone"`
	Note       string `json:"note"`
}

func (cc *CallController) AddCall(c *gin.Context) {
	var input AddCallRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	call := bson.M{
		"_id":         primitive.NewObjectID(),
		"client_name": input.ClientName,
		"numphone":    input.NumPhone,
		"note":        input.Note,
		"date":        time.Now(),
		"deleted_at":  nil,
	}

	_, err := cc.collection.InsertOne(context.Background(), call)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add call",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Call added successfully",
		"data":    call,
	})
}


// GetCalls - Ambil daftar panggilan dengan pagination
func (cc *CallController) GetCalls(c *gin.Context) {
	limit := 10
	page := 1
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if limit > 100 {
		limit = 100 // Batasi maksimum 100 per halaman
	}

	skip := (page - 1) * limit
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := cc.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch calls"})
		return
	}
	defer cursor.Close(context.Background())

	var calls []bson.M
	if err := cursor.All(context.Background(), &calls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode calls"})
		return
	}

	totalCount, err := cc.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count calls"})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	response := gin.H{
		"data": calls,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  totalCount,
			"total_pages":  totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetCallByID - Ambil data panggilan berdasarkan ID
func (cc *CallController) GetCallByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	var call bson.M
	err = cc.collection.FindOne(context.Background(), filter).Decode(&call)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Call not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

// UpdateCall - Perbarui data panggilan
func (cc *CallController) UpdateCall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input UpdateCallRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	update := bson.M{"$set": bson.M{
		"client_name": input.ClientName,
		"numphone":    input.NumPhone,
		"note":        input.Note,
	}}
	_, err = cc.collection.UpdateOne(context.Background(), bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update call", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call updated successfully"})
}

// DeleteCall - Hapus data panggilan (soft delete)
func (cc *CallController) DeleteCall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	deletedAt := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": deletedAt}}
	_, err = cc.collection.UpdateOne(context.Background(), bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete call", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call deleted successfully"})
}