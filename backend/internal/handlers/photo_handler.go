package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rjrajnish/photo_gallery/backend/internal/db"
	"github.com/rjrajnish/photo_gallery/backend/internal/models"
	"github.com/rjrajnish/photo_gallery/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhotoHandler struct{}

func NewPhotoHandler() *PhotoHandler { return &PhotoHandler{} }

func (h *PhotoHandler) List(c *gin.Context) {
	userId := c.GetString("userId")
	folderId := c.Query("folderId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId}
	if folderId != "" {
		filter["folderId"] = folderId
	}
	cur, _ := db.DB.Collection("photos").Find(ctx, filter)
	var out []models.Photo
	_ = cur.All(ctx, &out)
	c.JSON(http.StatusOK, out)
}

func (h *PhotoHandler) Upload(c *gin.Context) {

	userId := c.GetString("userId")
	folderId := c.PostForm("folderId")

	if folderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folderId required"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fid, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}
	// find folder mega node
	var folder models.Folder
	fmt.Println("	folderId", fid, userId)
	if err := db.DB.Collection("folders").FindOne(ctx, bson.M{"_id": fid, "userId": userId}).Decode(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder not found"})
		return
	}
	parent, _ := services.Mega().FindNodeByHandle(folder.MegaNode)

	form, _ := c.MultipartForm()
	files := form.File["files"]
	var saved []models.Photo

	for _, file := range files {
		f, _ := file.Open()
		defer f.Close()
		buf := new(bytes.Buffer)
		_, _ = io.Copy(buf, f)

		node, err := services.Mega().UploadBytes(parent, file.Filename, buf.Bytes())
		if err != nil {
			continue
		}
		photo := models.Photo{
			UserID:    userId,
			FolderID:  folderId,
			Filename:  file.Filename,
			Size:      file.Size,
			MegaNode:  node.GetHash(),
			CreatedAt: time.Now(),
		}
		res, _ := db.DB.Collection("photos").InsertOne(ctx, photo)
		photo.ID = res.InsertedID.(primitive.ObjectID).Hex()
		saved = append(saved, photo)
	}
	c.JSON(http.StatusOK, saved)
}

type deleteReq struct {
	PhotoIDs []string `json:"photoIds"`
}

func (h *PhotoHandler) Delete(c *gin.Context) {
	userId := c.GetString("userId")
	var body deleteReq
	if err := c.ShouldBindJSON(&body); err != nil || len(body.PhotoIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photoIds required"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	photosCol := db.DB.Collection("photos")
	for _, id := range body.PhotoIDs {
		var p models.Photo
		if err := photosCol.FindOne(ctx, bson.M{"_id": id, "userId": userId}).Decode(&p); err == nil {
			if node, err := services.Mega().FindNodeByHandle(p.MegaNode); err == nil {
				_ = services.Mega().DeleteNode(node)
			}
			photosCol.DeleteOne(ctx, bson.M{"_id": id})
		}
	}
	c.JSON(http.StatusOK, gin.H{"deleted": len(body.PhotoIDs)})
}
func (h *PhotoHandler) Stream(c *gin.Context) {
    userId := c.GetString("userId")
    photoId := c.Param("id")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pid, err := primitive.ObjectIDFromHex(photoId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
        return
    }

    var p models.Photo
    if err := db.DB.Collection("photos").FindOne(ctx, bson.M{
        "_id":    pid,
        "userId": userId,
    }).Decode(&p); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
        return
    }

    node, err := services.Mega().FindNodeByHandle(p.MegaNode)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "mega node not found"})
        return
    }

    c.Header("Content-Type", "image/jpeg")
    c.Status(http.StatusOK)
    _ = services.Mega().Download(node, c.Writer)
}

