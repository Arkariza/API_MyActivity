package LeadControllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Arkariza/API_MyActivity/models/ManageLead"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LeadController struct {
    collection *mongo.Collection
}

func NewLeadController(collection *mongo.Collection) *LeadController {
    return &LeadController{
        collection: collection,
    }
}

func SetLeadStatusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt("user_id")

        var status string
        switch userID {
        case 1:
            status = "Self"
        case 2:
            status = "Referral"
        default:
            status = "Unknown"
        }

        c.Set("lead_status", status)
        c.Next()
    }
}

func (lc *LeadController) AddLead(c *gin.Context) {
    var lead models.Lead

    if err := c.ShouldBindJSON(&lead); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    lead.ID = primitive.NewObjectID()
    lead.CreateAt = time.Now()

    status, exists := c.Get("lead_status")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Status not set"})
        return
    }
    lead.Status = status.(string)

    _, err := lc.collection.InsertOne(context.Background(), lead)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save lead"})
        return
    }

    c.JSON(http.StatusCreated, lead)
}

func (lc *LeadController) AddReferral(c *gin.Context) {
    var leads []models.Lead

    if err := c.ShouldBindJSON(&leads); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    status, exists := c.Get("lead_status")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Status not set"})
        return
    }

    documents := make([]interface{}, len(leads))
    for i := range leads {
        leads[i].ID = primitive.NewObjectID()
        leads[i].CreateAt = time.Now()
        leads[i].Status = status.(string)
        documents[i] = leads[i]
    }

    _, err := lc.collection.InsertMany(context.Background(), documents)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save referrals"})
        return
    }

    c.JSON(http.StatusCreated, leads)
}

func (lc *LeadController) GetLeads(c *gin.Context) {
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

    skip := (page - 1) * limit

    filter := bson.M{}

    if status := c.Query("status"); status != "" {
        filter["status"] = status
    }

    if startDate := c.Query("start_date"); startDate != "" {
        startTime, err := time.Parse("2006-01-02", startDate)
        if err == nil {
            if _, exists := filter["createAt"]; !exists {
                filter["createAt"] = bson.M{}
            }
            filter["createAt"].(bson.M)["$gte"] = startTime
        }
    }

    if endDate := c.Query("end_date"); endDate != "" {
        endTime, err := time.Parse("2006-01-02", endDate)
        if err == nil {
            if _, exists := filter["createAt"]; !exists {
                filter["createAt"] = bson.M{}
            }
            filter["createAt"].(bson.M)["$lte"] = endTime.Add(24 * time.Hour)
        }
    }

    findOptions := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(limit)).
        SetSort(bson.D{{Key: "createAt", Value: -1}})

    cursor, err := lc.collection.Find(context.Background(), filter, findOptions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
        return
    }
    defer cursor.Close(context.Background())

    totalCount, err := lc.collection.CountDocuments(context.Background(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count leads"})
        return
    }

    var leads []models.Lead
    if err := cursor.All(context.Background(), &leads); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode leads"})
        return
    }

    totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

    response := gin.H{
        "data": leads,
        "meta": gin.H{
            "current_page": page,
            "per_page":    limit,
            "total_items": totalCount,
            "total_pages": totalPages,
        },
    }

    c.JSON(http.StatusOK, response)
}

func (lc *LeadController) GetLeadByID(c *gin.Context) {
    id, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    filter := bson.M{"_id": id}
    var lead models.Lead
    err = lc.collection.FindOne(context.Background(), filter).Decode(&lead)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, lead)
}

// func (lc *LeadController) UpdateLead(c *gin.Context){

// }

// func (lc *LeadController) DeleteLead(c *gin.Context){

// }