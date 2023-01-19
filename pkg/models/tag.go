package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	TagName   string             `bson:"tag_name" json:"tagName" validate:"required,lowercase,alphanum,max=64"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type TagById struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}

type TagByName struct {
	TagName string `bson:"tag_name" json:"tagName" validate:"required,lowercase,alphanum,max=64"`
}
