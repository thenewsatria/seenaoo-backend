package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlashcardCover struct {
	FlashCardCoverSlug string
	Title              string
	Description        string
	Image_path         string
	Tags               []string
	UserId             primitive.ObjectID
}
