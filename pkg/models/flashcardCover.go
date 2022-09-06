package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardCover struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Slug        string               `bson:"slug" json:"slug"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Image_path  string               `bson:"image_path" json:"imagePath"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Author      primitive.ObjectID   `bson:"user_id" json:"userId"`
	CreatedAt   time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updatedAt"`
}

type FlashcardCoverBySlug struct {
	Slug string `bson:"slug" json:"slug"`
}

type FlashcardCoverById struct {
	ID string `bson:"_id" json:"id"`
}
