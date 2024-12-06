package MeetMiddleware

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/Arkariza/API_MyActivity/models/CallAndMeet"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
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

func ValidateMeetRequest() gin.HandlerFunc {
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

        c.Set("meet_data", meet)
        c.Next()
    }
}