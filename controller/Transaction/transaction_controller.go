package TransactionController

import (
	"context"
	"net/http"

	"github.com/Arkariza/API_MyActivity/models/ManageLead"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionController struct {
	leadCollection *mongo.Collection
}

func NewTransactionController(leadCollection *mongo.Collection) *TransactionController {
	return &TransactionController{leadCollection: leadCollection}
}

func (tc *TransactionController) GetAllTransactions(c *gin.Context) {
	filter := bson.M{"status": "Pending"}

	cursor, err := tc.leadCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads", "details": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var pendingLeads []models.Lead
	if err := cursor.All(context.Background(), &pendingLeads); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode leads", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pending_leads": pendingLeads})
}

func (tc *TransactionController) UpdateTransactionStatus(c *gin.Context) {
	id := c.Param("_id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}
	filter := bson.M{"_id": objID, "status": "Pending"}
	update := bson.M{"$set": bson.M{"status": "Open"}}
	result, err := tc.leadCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status", "details": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No transaction found with the given ID and status 'Pending'"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated to 'Open'", "id": id})
}