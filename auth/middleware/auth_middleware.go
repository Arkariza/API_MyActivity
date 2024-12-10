package AuthMiddleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/Arkariza/API_MyActivity/auth"
)

func AuthMiddleware(authCommand *auth.AuthCommand) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Authorization header is required",
            })
            ctx.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Invalid authorization header format",
            })
            ctx.Abort()
            return
        }

        claims, err := authCommand.ValidateToken(parts[1])
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Invalid token",
                "error":   err.Error(),
            })
            ctx.Abort()
            return
        }

        user, err := authCommand.GetUserFromToken(claims)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Invalid user token",
                "error":   err.Error(),
            })
            ctx.Abort()
            return
        }

        ctx.Set("user", user)
        ctx.Set("userID", user.ID) 
        ctx.Set("userRole", user.Role)
        ctx.Next()
    }
}


func AuthenticateUser(authCommand *auth.AuthCommand) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Authorization header is required",
            })
            ctx.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Invalid authorization header format",
            })
            ctx.Abort()
            return
        }

        claims, err := authCommand.ValidateToken(parts[1])
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "status":  false,
                "message": "Invalid token",
                "error":   err.Error(),
            })
            ctx.Abort()
            return
        }

        ctx.Set("claims", claims) 
        ctx.Next()
    }
}