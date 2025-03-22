package LeadController

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Arkariza/API_MyActivity/models/ManageLead"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadController struct {
	collection *mongo.Collection
}

func NewLeadController(collection *mongo.Collection) *LeadController {
	return &LeadController{collection: collection}
}

type AddLeadRequest struct {
	UserID      primitive.ObjectID `json:"user_id"`
	ClientName  string             `json:"clientname" binding:"required"`
	NumPhone    string             `json:"numphone" binding:"required"`
	Priority    string             `json:"priority" binding:"required"`
	Information string             `json:"information"`
	NoPolicy    int64              `json:"no_policy"`
	Status      string             `json:"status"`
	TypeLead    string             `json:"type_lead"`
}

func ValidateLeadInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is missing"})
			c.Abort()
			return
		}

		var input AddLeadRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("lead_input", input)
		c.Next()
	}
}
 

func (lc *LeadController) GetAllLead(c *gin.Context) {
	cursor, err := lc.collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leads"})
		return
	}
	defer cursor.Close(context.Background())

	var leads []models.Lead
	if err := cursor.All(context.Background(), &leads); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode leads"})
		return
	}

	totalCount, err := lc.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count leads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leads": leads,
		"total": totalCount,
	})
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

func (lc *LeadController) getUsersByRole(role int) ([]string, error) {
    userCollection := lc.collection.Database().Collection("users")
    cursor, err := userCollection.Find(context.Background(), bson.M{"role": role})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var users []struct {
        UserName string `bson:"username"` // Asumsikan field username ada di koleksi user
    }
    if err := cursor.All(context.Background(), &users); err != nil {
        return nil, err
    }

    var userNames []string
    for _, user := range users {
        userNames = append(userNames, user.UserName)
    }
    return userNames, nil
}

func (cc *LeadController) AddLead(c *gin.Context, req AddLeadRequest) (*models.Lead, error) {
    userRole, roleExists := c.Get("Role")
    userID, idExists := c.Get("UserID")

    if !roleExists || !idExists {
        handleError(c, http.StatusForbidden, "User role or ID missing", nil)
        return nil, errors.New("user role or ID missing")
    }

    if c.Request.Body == nil {
        handleError(c, http.StatusBadRequest, "Empty request body", nil)
        return nil, errors.New("empty request body")
    }

    body, readErr := io.ReadAll(c.Request.Body)
    if readErr != nil {
        handleError(c, http.StatusBadRequest, "Cannot read request body", readErr)
        return nil, readErr
    }

    c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

    if err := json.Unmarshal(body, &req); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid JSON", err)
        return nil, err
    }

    _, err := validateToken(c)
    if err != nil {
        handleError(c, http.StatusUnauthorized, "Invalid authentication", err)
        return nil, err
    }

    const startRange int64 = 3200000000
    const rangeSize int64 = 100000000
    rand.Seed(time.Now().UnixNano())
    randomNoPolicy := startRange + rand.Int63n(rangeSize)

    lead := models.Lead{
        ID:          primitive.NewObjectID(),
        UserID:      req.UserID,
        NumPhone:    req.NumPhone,
        Priority:    req.Priority,
        Latitude:    0,
        Longitude:   0,
        CreateAt:    time.Now(),
        DateSubmit:  time.Time{},
        ClientName:  req.ClientName,
        Information: req.Information,
        NoPolicy:    randomNoPolicy,
    }

    switch userRole.(int) {
    case 1: // Self
        lead.Status = models.StatusOpen
        lead.TypeLead = models.TypeSelf

    case 2: 
        lead.Status = models.StatusPending
        lead.TypeLead = models.TypeReferral

        userNames, err := cc.getUsersByRole(1)
        if err != nil || len(userNames) == 0 {
            handleError(c, http.StatusInternalServerError, "No users with Role 1 available", err)
            return nil, errors.New("no users with Role 1 found")
        }
        lead.AsignTo = userNames[rand.Intn(len(userNames))]

    default:
        handleError(c, http.StatusForbidden, "Invalid user role for this operation", nil)
        return nil, errors.New("user role not allowed")
    }

    parsedID, parseErr := primitive.ObjectIDFromHex(userID.(string))
    if parseErr != nil {
        handleError(c, http.StatusBadRequest, "Invalid user ID format", parseErr)
        return nil, parseErr
    }
    lead.UserID = parsedID

    _, dbErr := cc.collection.InsertOne(c, lead)
    if dbErr != nil {
        handleError(c, http.StatusInternalServerError, "Failed to save lead", dbErr)
        return nil, dbErr
    }

    return &lead, nil
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		c.JSON(statusCode, gin.H{
			"error":   message,
			"details": err.Error(),
		})
	} else {
		c.JSON(statusCode, gin.H{
			"error": message,
		})
	}
}

func (lc *LeadController) GetLeadByID(c *gin.Context) {
	leadID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}
	var lead models.Lead
	err = lc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&lead)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lead"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"lead": lead,
	})
}

func (cc *LeadController) ChangeLeadStatus(c *gin.Context) {
    type ChangeStatusRequest struct {
        LeadID string `json:"lead_id"`
        Status string `json:"status"`
    }
    var req ChangeStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request body", err)
        return
    }
    leadID, err := primitive.ObjectIDFromHex(req.LeadID)
    if err != nil {
        handleError(c, http.StatusBadRequest, "Invalid lead ID format", err)
        return
    }
    req.Status = models.StatusOpen

    filter := bson.M{"_id": leadID}
    update := bson.M{"$set": bson.M{"status": req.Status, "updated_at": time.Now()}}

    result, err := cc.collection.UpdateOne(c, filter, update)
    if err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to update lead status", err)
        return
    }

    if result.MatchedCount == 0 {
        handleError(c, http.StatusNotFound, "Lead not found", nil)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Lead status updated successfully",
        "lead_id": req.LeadID,
        "status": req.Status,
    })
}