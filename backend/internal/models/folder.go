package models

import "time"

type Folder struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"userId" json:"userId"`
	Name      string    `bson:"name" json:"name"`
	MegaNode  string    `bson:"megaNode" json:"megaNode"` // MEGA folder node handle
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
