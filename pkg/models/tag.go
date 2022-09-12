package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	TagName   string             `bson:"tag_name" json:"tagName"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

type TagById struct {
	ID string `bson:"_id" json:"id"`
}

type TagByName struct {
	TagName string `bson:"tag_name" json:"tagName"`
}
