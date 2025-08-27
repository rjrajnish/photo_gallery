package models

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"` // bcrypt hash
	RootNode  string    `bson:"rootNode" json:"rootNode"` // MEGA node handle for the user's root folder
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
