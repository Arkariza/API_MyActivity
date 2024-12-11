package CommentController

import (
	"context"
	"errors"
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

type CommentController struct {
	Collection *mongo.Collection
}

func NewCommentController(collection *mongo.Collection) *CommentController {
	return &CommentController{Collection: collection}
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

func (cc *CommentController) CreateComment(c *gin.Context) {
    _, err := validateToken(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
        return
    }

    var comment models.Comment
    if err := c.ShouldBindJSON(&comment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    comment.ID = primitive.NewObjectID()

    if err := comment.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if comment.Date.IsZero() {
        comment.Date = time.Now()
    }
    result, err := cc.Collection.InsertOne(context.Background(), comment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "failed to insert comment",
            "details": err.Error(),
        })
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "message": "Comment created successfully",
        "comment": comment,
        "insertedID": result.InsertedID,
    })
}


func (cc *CommentController) GetAllComments(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	skip := (pageNum - 1) * limitNum

	totalCount, err := cc.Collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count comments"})
		return
	}

	cursor, err := cc.Collection.Find(context.Background(), bson.M{}, 
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limitNum)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"page":     pageNum,
		"limit":    limitNum,
		"total":    totalCount,
	})
}

func (cc *CommentController) GetCommentByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	var comment models.Comment
	err = cc.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (cc *CommentController) UpdateComment(c *gin.Context) {

	_, err := validateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user role not found"})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	var updatedComment models.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userRole.(int) != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only users with Role 1 can update comments"})
		return
	}

	if err := updatedComment.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedComment.Date = time.Now()

	update := bson.M{"$set": updatedComment}
	result, err := cc.Collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update comment"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"id":      objectID.Hex(),
	})
}

func (cc *CommentController) DeleteComment(c *gin.Context) {
	_, err := validateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}
	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user role not found"})
		return
	}
	if userRole.(int) != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only users with Role 1 can delete comments"})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	result, err := cc.Collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
		"id":      objectID.Hex(),
	})
}