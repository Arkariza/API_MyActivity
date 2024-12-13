package CallMiddleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Arkariza/API_MyActivity/models/CallAndMeet"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Call struct {
	ID             string    `json:"id,omitempty"`         
	ClientName     string    `json:"client_name" binding:"required"`
	PhoneNum       string    `json:"phonenum" binding:"required"`
	ProspectStatus string    `json:"prospect_status,omitempty"`
	Date           time.Time `json:"date" binding:"required"`
	Note           string    `json:"note,omitempty"`
}

type CallMiddleware struct {
	secretKey string
}

func NewCallMiddleware(secretKey string) *CallMiddleware {
	return &CallMiddleware{
		secretKey: secretKey,
	}
}

func (m *CallMiddleware) AuthenticateCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(m.secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func ValidateCallRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var call models.Call
		if err := c.ShouldBindJSON(&call); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		if err := call.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("call_data", call)
		c.Next()
	}
}