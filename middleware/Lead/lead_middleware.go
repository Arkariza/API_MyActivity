package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Arkariza/API_MyActivity/controller/Lead"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LeadMiddleware struct {
	secretKey string
}

func NewLeadMiddleware(secretKey string) *LeadMiddleware {
	return &LeadMiddleware{
		secretKey: secretKey,
	}
}
func (m *LeadMiddleware) AuthenticateLead() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }
            return []byte(m.secretKey), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        userID, userOk := claims["user_id"].(string)
        role, roleOk := claims["role"].(float64)

        if !userOk || !roleOk {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
            c.Abort()
            return
        }
        c.Set("UserID", userID)
        c.Set("Role", int(role))
        
        c.Next()
    }
}

func ValidateLeadInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LeadController.AddLeadRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("lead_input", input)
		c.Next()
	}
}

func BodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		fmt.Println("Request Body:", string(body))

		c.Next()
	}
}
