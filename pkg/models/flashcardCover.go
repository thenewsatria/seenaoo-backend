package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardCover struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id" validate:"required"`
	Slug        string               `bson:"slug" json:"slug" validate:"required"`
	Title       string               `bson:"title" json:"title" validate:"required,min=5,max=128"`
	Description string               `bson:"description" json:"description" validate:"min=5,max=255"`
	Image_path  string               `bson:"image_path" json:"imagePath" validate:"omitempty,file"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Author      string               `bson:"author" json:"author" validate:"required,lowercase,alphanum,max=25,min=5"`
	CreatedAt   time.Time            `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type FlashcardCoverRequest struct {
	Title       string   `bson:"title" json:"title" validate:"required,min=5,max=128"`
	Description string   `bson:"description" json:"description" validate:"omitempty,min=5,max=255"`
	Image_path  string   `bson:"image_path" json:"imagePath" validate:"omitempty,file"`
	Tags        []string `bson:"tags" json:"tags" validate:"omitempty,dive,omitempty,lowercase,alphanum,max=32"`
}

type FlashcardCoverBySlug struct {
	Slug string `bson:"slug" json:"slug" validate:"required"`
}

type FlashcardCoverById struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}
