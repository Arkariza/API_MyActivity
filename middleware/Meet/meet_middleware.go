package MeetMiddleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/Arkariza/API_MyActivity/models/CallAndMeet"
)

type MeetMiddleware struct {
	secretKey string
}

func NewMeetMiddleware(secretKey string) *MeetMiddleware {
	return &MeetMiddleware{
		secretKey: secretKey,
	}
}

func (m *MeetMiddleware) AuthenticateMeet() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		if userID, exists := claims["user_id"].(float64); exists {
			c.Set("user_id", int(userID))
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *MeetMiddleware) ValidateMeetRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var meet models.Meet
		if err := c.ShouldBindJSON(&meet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if err := meet.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		userID, exists := c.Get("user_id")
		if exists {
			meet.UserID = userID.(int)
		}

		c.Set("meet_data", meet)
		c.Next()
	}
}