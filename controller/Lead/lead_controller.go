package LeadController

import (
	"context"
	"net/http"
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

func (lc *LeadController) AddLeadToDB(lead models.Lead) error {
	_, err := lc.collection.InsertOne(context.Background(), lead)
	return err
}

func (lc *LeadController) CreateLeadWithAutoStatus(c *gin.Context) {
	var input models.LeadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole, exists := c.Get("Role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var status string
	switch userRole.(int) {
	case 1:
		status = "Referral"
	case 2:
		status = "Self"
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user ID for this operation"})
		return
	}

	lead := models.Lead{
		ID:          primitive.NewObjectID(),
		NumPhone:    input.NumPhone,
		Priority:    input.Priority,
		Latitude:    0,
		Longitude:   0,
		CreateAt:    time.Now(),
		DateSubmit:  time.Time{},
		ClientName:  input.ClientName,
		IdBFA:       0,
		IdRefeal:    0,
		TypeLead:    "",
		NoPolicy:    input.NoPolicy,
		Information: input.Information,
		Status:      status,
	}

	if err := lc.AddLeadToDB(lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lead"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lead created successfully with auto-status",
		"lead":    lead,
	})
}
