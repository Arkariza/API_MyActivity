package LeadController

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arkariza/API_MyActivity/models/ManageLead"
	"github.com/gin-gonic/gin"
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

    if !roleExists || !idExists {
        handleError(c, http.StatusForbidden, "User role or ID missing", nil)
        return nil, errors.New("user role or ID missing")
    }
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
    }
    switch userRole.(int) {
    case 1:
        lead.Status = models.StatusPending
        lead.TypeLead = models.TypeSelf
    case 2:
        lead.Status = models.StatusPending
        lead.TypeLead = models.TypeReferral
    default:
        handleError(c, http.StatusForbidden, "Invalid user role for this operation", nil)
        return nil, errors.New("invalid user role")
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