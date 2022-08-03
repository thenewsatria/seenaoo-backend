package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flashcard struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	FrontImagePath   string             `bson:"front_image_path" json:"frontImagePath"`
	BackImagePath    string             `bson:"back_image_path" json:"backImagePath"`
	FrontText        string             `bson:"front_text" json:"frontText"`
	BackText         string             `bson:"back_text" json:"backText"`
	Question         string             `bson:"question" json:"question"`
	FlashCardCoverId primitive.ObjectID `bson:"flashcard_cover_id" json:"flashcardCoverId"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type DeleteFlashcardRequest struct {
	ID string `json:"id"`
}

type ReadFlashcardRequest struct {
	ID string `json:"id"`
}
