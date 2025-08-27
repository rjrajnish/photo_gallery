package models
 
import "time"

type Photo struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"userId" json:"userId"`
	FolderID  string    `bson:"folderId" json:"folderId"`
	Filename  string    `bson:"filename" json:"filename"`
	Size      int64     `bson:"size" json:"size"`
	MegaNode  string    `bson:"megaNode" json:"megaNode"` // MEGA file node handle
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
