package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/t3rm1n4l/go-mega"

	"github.com/Rjrajnish99/photo-gallery/backend/internal/db"
	"github.com/Rjrajnish99/photo-gallery/backend/internal/models"
	"github.com/Rjrajnish99/photo-gallery/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler { return &AuthHandler{} }

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	_ = godotenv.Load()
	var body registerReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users := db.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ensure unique email
	count, _ := users.CountDocuments(ctx, bson.M{"email": body.Email})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	// Create user root folder in MEGA
	root, _ := services.Mega().Root()
	uuidRoot := "user-" + uuid.NewString()
	dir, err := services.Mega().CreateFolder(root, uuidRoot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user root on MEGA"})
		return
	}

	hash, _ := services.HashPassword(body.Password)
	user := models.User{
		Email:     body.Email,
		Password:  hash,
		RootNode:  dir.GetHash(), // MEGA handle
		CreatedAt: time.Now(),
	}
	_, err = users.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body loginReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users := db.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := users.FindOne(ctx, bson.M{"email": body.Email}).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !services.CheckPassword(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := services.CreateJWT(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "email": user.Email})
}
