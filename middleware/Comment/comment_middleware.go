package CommadnMiddleware

import (
	"fmt"
    "net/http"
    "strings"
     
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
)

type CommandMiddleware struct {
	secretKey string
}

func NewCommandMiddleware(secretKey string) *CommandMiddleware {
	return &CommandMiddleware{
		secretKey: secretKey,
	}
}

func (m *CommandMiddleware) AuthenticateCommand() gin.HandlerFunc {
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
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token must start with 'Bearer '",
			})
			c.Abort()
			return
		}

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
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing user ID in token",
			})
			c.Abort()
			return
		}

		role, ok := claims["role"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing role in token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", int(role))

 		fmt.Printf("Authenticated User ID: %s, Role: %d\n", userID, int(role))

		c.Next()
	}
}