package CallControllers

import (
	"context"
	"errors"
	"fmt"
	"math"
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

type CallController struct {
	collection *mongo.Collection
}

func NewCallController(collection *mongo.Collection) *CallController {
	return &CallController{
		collection: collection,
	}
}

type AddCallRequest struct {
    ClientName      string `json:"client_name" binding:"required"`
    PhoneNum        string `json:"phonenum" binding:"required"`
    Note            string `json:"note,omitempty"`
    ProspectStatus  string `json:"prospect_status,omitempty"`
    CallResult      string `json:"call_result,omitempty"`
}

type UpdateCallRequest struct {
	ClientName     string `json:"client_name"`
	PhoneNum       string `json:"Phone_num"`
	Note           string `json:"note"`
	ProspectStatus string `json:"prospect_status"`
	CallResult     string `json:"call_result"`
}

func validateToken(c *gin.Context) (string, error) {
    authHeader := c.GetHeader("Authorization")
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return "", errors.New("invalid token format")
    }

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
        return "", errors.New("empty token")
    }

	return tokenString, nil
}

func (cc *CallController) AddCall(c *gin.Context, req AddCallRequest) (*models.Call, error) {
    _, err := validateToken(c)
    if err != nil {
        handleError(c, http.StatusUnauthorized, "Invalid authentication", err)
        return nil, err
    }

    call := models.Call{
        ID:              primitive.NewObjectID(),
        ClientName:      req.ClientName,
        PhoneNum:        req.PhoneNum,
        Note:            req.Note,
        CreatedAt:       time.Now(),
        Date:            time.Now(),
        ProspectStatus:  req.ProspectStatus,
        CallResult:      req.CallResult,
    }

    if call.ProspectStatus == "" {
        call.ProspectStatus = "new"
    }
    if call.CallResult == "" {
        call.CallResult = "Pending"
    }
    if call.Note == "" {
        call.Note = "No additional notes provided."
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = cc.collection.InsertOne(ctx, call)
    if err != nil {
        return nil, fmt.Errorf("failed to create call: %v", err)
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Call created successfully",
        "data":    call,
    })

    return &call, nil
}

func handleError(c *gin.Context, i int, s string, err error) {
	panic("unimplemented")
}

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
		limit = 100
	}

	skip := (page - 1) * limit

	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	if searchQuery := c.Query("search"); searchQuery != "" {
		filter["$or"] = []bson.M{
			{"client_name": bson.M{"$regex": primitive.Regex{Pattern: searchQuery, Options: "i"}}},
			{"phonenum": bson.M{"$regex": primitive.Regex{Pattern: searchQuery, Options: "i"}}},
		}
	}

	if status := c.Query("status"); status != "" {
		filter["prospect_status"] = status
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "date", Value: -1}})

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

func (cc *CallController) UpdateCall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	if req.ClientName != "" {
		req.ClientName = strings.TrimSpace(req.ClientName)
		if len(req.ClientName) < 2 || len(req.ClientName) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Client name must be between 2 and 100 characters",
			})
			return
		}
	}

	validStatuses := map[string]bool{
		"new":          true,
		"in_progress":  true,
		"contacted":    true,
		"qualified":    true,
		"unqualified":  true,
		"follow_up":    true,
	}
	if req.ProspectStatus != "" && !validStatuses[req.ProspectStatus] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid prospect status",
		})
		return
	}

	update := bson.M{"$set": bson.M{
		"client_name":     req.ClientName,
		"phone_num":       req.PhoneNum,
		"note":            req.Note,
		"prospect_status": req.ProspectStatus,
		"call_result":     req.CallResult,
	}}

	for k, v := range update["$set"].(bson.M) {
		if v == "" {
			delete(update["$set"].(bson.M), k)
		}
	}

	result, err := cc.collection.UpdateOne(
		context.Background(), 
		bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}, 
		update,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update call",
			"details": err.Error(),
		})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Call not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call updated successfully"})
}

func (cc *CallController) DeleteCall(c *gin.Context) {
    id, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    result, err := cc.collection.DeleteOne(
        context.Background(), 
        bson.M{"_id": id},
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to delete call",
            "details": err.Error(),
        })
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Call not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Call deleted successfully"})
}