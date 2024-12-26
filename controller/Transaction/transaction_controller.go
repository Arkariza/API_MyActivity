package TransactionController

import (
	"context"
	"net/http"

	"github.com/Arkariza/API_MyActivity/models/ManageLead"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
