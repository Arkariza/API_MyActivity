package LeadMiddleware

import (
	"fmt"
	"net/http"
	"strings"
    // "strconv"
    
	"github.com/Arkariza/API_MyActivity/controller/Lead"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LeadMiddleware struct{
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


func ValidateLeadInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LeadController.LeadController
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			c.Abort()
			return
		}
		c.Set("leadInput", input)
		c.Next()
	}
}
