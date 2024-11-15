package UserControllers

import (
	"context"
	"net/http"

	"github.com/Arkariza/API_MyActivity/auth"
	"github.com/Arkariza/API_MyActivity/models/User"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
    authCommand *auth.AuthCommand
}

func NewUserController(authCommand *auth.AuthCommand) *UserController {
    return &UserController{
        authCommand: authCommand,
    }
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    PhoneNum string `json:"phone_num" binding:"required"`
    Role     int    `json:"role" binding:"required,oneof=1 2"`
}

func (c *UserController) Register(ctx *gin.Context) {
    var request RegisterRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    cmdRequest := auth.RegisterRequest{
        Username: request.Username,
        Email:    request.Email,
        Password: request.Password,
        PhoneNum: request.PhoneNum,
        Role:     request.Role,
    }
    ctxRequest := context.Background()
    user, err := c.authCommand.Register(ctxRequest, cmdRequest)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Registration failed",
            "error":   err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "status":  true,
        "message": "Registration successful",
        "data":    user,
    })
}

func (c *UserController) Login(ctx *gin.Context) {
    var request LoginRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    cmdRequest := auth.LoginRequest{
        Username: request.Username,
        Password: request.Password,
    }
    ctxRequest := context.Background()
    tokenResponse, err := c.authCommand.Login(ctxRequest, cmdRequest)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "status":  false,
            "message": "Login failed",
            "error":   err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Login successful",
        "data":    tokenResponse,
    })
}

func (c *UserController) GetProfile(ctx *gin.Context) {
    user, exists := ctx.Get("user")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "status":  false,
            "message": "Unauthorized",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Profile retrieved successfully",
        "data":    user,
    })
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
    var request struct {
        Email    string `json:"email" binding:"omitempty,email"`
        PhoneNum string `json:"phone_num"`
        Image    string `json:"image"`
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    currentUser, exists := ctx.Get("user")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "status":  false,
            "message": "Unauthorized",
        })
        return
    }

    user := currentUser.(*models.User)

    if request.Email != "" {
        user.Email = request.Email
    }
    if request.PhoneNum != "" {
        user.PhoneNum = request.PhoneNum
    }
    if request.Image != "" {
        user.Image = request.Image
    }

    if err := ctx.MustGet("db").(*gorm.DB).Save(user).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "Failed to update profile",
            "error":   err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Profile updated successfully",
        "data":    user,
    })
}