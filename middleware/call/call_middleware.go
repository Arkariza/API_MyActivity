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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort() 
			return
		}
		c.Set("ValidateCall", input) 
		c.Next() 
	}
}
