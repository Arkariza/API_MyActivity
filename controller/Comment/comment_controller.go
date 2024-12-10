package CommandController

import (
	"context"
	"net/http"
	"time"

	"github.com/Arkariza/API_MyActivity/models/CallAndMeet"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandController struct {
	Collection *mongo.Collection
}

func NewCommandController(collection *mongo.Collection) *CommandController {
	return &CommandController{Collection: collection}
}

func (cc *CommandController) CreateCommand(c *gin.Context) {
	var command models.Command
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command.ID = primitive.NewObjectID()
	if err := command.BeforeCreate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set default values"})
		return
	}

	_, err := cc.Collection.InsertOne(context.Background(), command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert command"})
		return
	}

	c.JSON(http.StatusCreated, command)
}

func (cc *CommandController) GetAllCommands(c *gin.Context) {
	cursor, err := cc.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch commands"})
		return
	}
	defer cursor.Close(context.Background())

	var commands []models.Command
	if err := cursor.All(context.Background(), &commands); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode commands"})
		return
	}

	c.JSON(http.StatusOK, commands)
}

func (cc *CommandController) GetCommandByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid command ID"})
		return
	}

	var command models.Command
	err = cc.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&command)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "command not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch command"})
		return
	}

	c.JSON(http.StatusOK, command)
}

func (cc *CommandController) UpdateCommand(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid command ID"})
		return
	}

	var updatedCommand models.Command
	if err := c.ShouldBindJSON(&updatedCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCommand.Date = time.Now()
	update := bson.M{"$set": updatedCommand}
	_, err = cc.Collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update command"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "command updated successfully"})
}

func (cc *CommandController) DeleteCommand(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid command ID"})
		return
	}

	_, err = cc.Collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete command"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "command deleted successfully"})
}