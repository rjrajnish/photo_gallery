package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Rjrajnish99/photo_gallery/backend/internal/db"
	"github.com/Rjrajnish99/photo_gallery/backend/internal/models"
	"github.com/Rjrajnish99/photo_gallery/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderHandler struct{}

func NewFolderHandler() *FolderHandler { return &FolderHandler{} }

type createFolderReq struct {
	Name string `json:"name" binding:"required"`
}

func (h *FolderHandler) List(c *gin.Context) {
	userId := c.GetString("userId")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, _ := db.DB.Collection("folders").Find(ctx, bson.M{"userId": userId})
	var out []models.Folder
	_ = cur.All(ctx, &out)
	c.JSON(http.StatusOK, out)
}

func (h *FolderHandler) Create(c *gin.Context) {
	userId := c.GetString("userId")
	var body createFolderReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Find user to get root node
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = db.DB.Collection("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&user)

	rootNode, _ := services.Mega().FindNodeByHandle(user.RootNode)
	dir, err := services.Mega().CreateFolder(rootNode, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create folder on MEGA"})
		return
	}
	f := models.Folder{
		UserID:    userId,
		Name:      body.Name,
		MegaNode:  dir.GetHash(),
		CreatedAt: time.Now(),
	}
	res, err := db.DB.Collection("folders").InsertOne(ctx, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
		return
	}
	f.ID = res.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusOK, f)
}

type deleteFolderReq struct {
	FolderID string `json:"folderId" binding:"required"`
}

func (h *FolderHandler) Delete(c *gin.Context) {
	userId := c.GetString("userId")
	var body deleteFolderReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var folder models.Folder
	if err := db.DB.Collection("folders").FindOne(ctx, bson.M{"_id": body.FolderID, "userId": userId}).Decode(&folder); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	// Delete photos in folder
	db.DB.Collection("photos").DeleteMany(ctx, bson.M{"folderId": folder.ID, "userId": userId})

	// Delete MEGA node
	node, _ := services.Mega().FindNodeByHandle(folder.MegaNode)
	_ = services.Mega().DeleteNode(node)

	// Delete folder doc
	db.DB.Collection("folders").DeleteOne(ctx, bson.M{"_id": folder.ID})
	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}
