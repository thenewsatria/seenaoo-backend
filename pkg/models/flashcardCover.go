package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlashcardCover struct {
	FlashCardCoverSlug string               `bson:"flashcard_cover_slug" json:"flashcardCoverSlug"`
	Title              string               `bson:"title" json:"title"`
	Description        string               `bson:"description" json:"description"`
	Image_path         string               `bson:"image_path" json:"imagePath"`
	Tags               []primitive.ObjectID `bson:"tags" json:"tags"`
	UserId             primitive.ObjectID   `bson:"user_id" json:"userId"`
}
