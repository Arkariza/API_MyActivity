package LeadMiddleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
    secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
    return &AuthMiddleware{
        secretKey: secretKey,
    }
}

func (m *AuthMiddleware) AuthenticateUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(m.secretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        if userID, exists := claims["user_id"].(float64); exists {
            c.Set("user_id", int(userID))
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func SetLeadStatus() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
            c.Abort()
            return
        }

        var status string
        switch userID.(int) {
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

func ValidateLeadInput() gin.HandlerFunc {
    return func(c *gin.Context) {
        var lead struct {
            NumPhone    *int        `json:"NumPhone"`
            ClientName  *string     `json:"ClientName"`
            Priority    *string     `json:"Priority"`
            Information *string     `json:"Information"`
            Latitude    *float64    `json:"latitude"`
            Longitude   *float64    `json:"longitude"`
            DateSubmit  *time.Time  `json:"DateSubmit"`
        }

        if err := c.ShouldBindJSON(&lead); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
            c.Abort()
            return
        }

        if lead.NumPhone == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "NumPhone is required"})
            c.Abort()
            return
        }

        if lead.ClientName == nil || *lead.ClientName == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ClientName is required"})
            c.Abort()
            return
        }

        if lead.Priority == nil || *lead.Priority == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Priority is required"})
            c.Abort()
            return
        }

        if lead.Latitude == nil || lead.Longitude == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Location coordinates are required"})
            c.Abort()
            return
        }

        c.Set("validated_lead", lead)
        c.Next()
    }
}