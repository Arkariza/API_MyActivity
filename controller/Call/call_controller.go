package CallControllers

import (
	"net/http"
	"time"

	"github.com/Arkariza/API_MyActivity/models/CallAndMeet"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type CallController struct {
	DB *gorm.DB 
}

func NewCallController(db *gorm.DB) *CallController {
	return &CallController{DB: db} 
}

type AddCall struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientName     string             `json:"client_name" binding:"required"`
	Numphone       int                `json:"numphone" binding:"required"`
	ProspectStatus string             `json:"prospect_status"`
	Date           time.Time          `json:"date" binding:"required"`
	Note           string             `json:"note"`
}

func (cc *CallController) CallList(c *gin.Context) {
	var calls []models.Call 

	if err := cc.DB.Find(&calls).Error; err != nil { // Gorm query untuk mengambil data
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving calls", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calls": calls})
}

func (cc *CallController) AddCall(c *gin.Context) {
	var input AddCall

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	call := models.Call{
		ClientName:     input.ClientName,
		Numphone:       input.Numphone,
		ProspectStatus: input.ProspectStatus,
		Date:           input.Date,
		Note:           input.Note,
	}

	if err := cc.DB.Create(&call).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add call", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call added successfully", "call": call})
}
