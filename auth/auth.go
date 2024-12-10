package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Arkariza/API_MyActivity/models/User"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecretEnvKey = "JWT_SECRET"

type AuthCommand struct {
	collection *mongo.Collection
}

func NewAuthCommand(collection *mongo.Collection) *AuthCommand {
	return &AuthCommand{collection: collection}
}

func (c *AuthCommand) GetSecretKey() string {
	return os.Getenv(jwtSecretEnvKey)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	PhoneNum string `json:"phone_num"`
	Role     int    `json:"role"`
}

type TokenResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	Username    string  	 `json:"username"`
	Role        int          `json:"role"`
	PhoneNum    string       `json:"phone_num"`
	Email       string       `json:"email"`
}

func (c *AuthCommand) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	var user models.User
	err := c.collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("invalid credentials")
	} else if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.LastLogin = time.Now()
	_, err = c.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"last_login": user.LastLogin}})
	if err != nil {
		return nil, err
	}

	token, err := c.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   24 * 60 * 60, 
		Username:    user.Username,
		Role:        user.Role,
		PhoneNum:    user.PhoneNum,
		Email:       user.Email,
	}, nil
}

func (c *AuthCommand) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	var existingUser models.User
	err := c.collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		PhoneNum:  req.PhoneNum,
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	_, err = c.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *AuthCommand) generateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"phone_num":user.PhoneNum,
		"email":	user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (c *AuthCommand) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.GetSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (c *AuthCommand) GetUserFromToken(claims jwt.MapClaims) (*models.User, error) {
	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := c.collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}