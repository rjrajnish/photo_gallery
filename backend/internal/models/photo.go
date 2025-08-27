package models

import "time"

type Photo struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"userId" json:"userId"`
	FolderID  string    `bson:"folderId" json:"folderId"`
	Filename  string    `bson:"filename" json:"filename"`
	URL       string    `bson:"-" json:"url"`      // computed field for frontend
	Size      int64     `bson:"size" json:"size"`
	MegaNode  string    `bson:"megaNode" json:"megaNode"` // Mega file node handle
	MegaKey   string    `bson:"megaKey" json:"megaKey"`   // Mega decryption key
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
