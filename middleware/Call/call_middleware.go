package CallMiddleware

import (
	"net/http"
	"time" 

	"github.com/gin-gonic/gin"                          
)

type Call struct {
	ID             string    `json:"id,omitempty"`         
	ClientName     string    `json:"client_name" binding:"required"`
	Numphone       int       `json:"numphone" binding:"required"`
	ProspectStatus string    `json:"prospect_status,omitempty"`
	Date           time.Time `json:"date" binding:"required"`
	Note           string    `json:"note,omitempty"`
}

func ValidateAddCall() gin.HandlerFunc {
	return func(c *gin.Context) {
 
		var input Call
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
			c.Abort() 
			return
		}
		c.Set("validate_call", input) 

		c.Next() 
	}
}

func setCallTimestampMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentTime := time.Now()
		c.Set("call_timestamp", currentTime)
		c.Next()
	}
}

func SetCallTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientName := c.GetInt("client_name")
		var callType string

		switch ClientName {
		case 1:
			callType = "Personal"
		case 2:
			callType = "Business"
		case 3:
			callType = "Unknown"
		}
		c.Set("call_type", callType)
		c.Next()
	}
}

func ExampleMiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Processed-Time", time.Now().Format(time.RFC3339))
		c.Next()
	}
} 

