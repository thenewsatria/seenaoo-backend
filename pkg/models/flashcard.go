package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flashcard struct {
	ID               primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	FrontImagePath   string             `bson:"front_image_path" json:"frontImagePath" validate:"omitempty,file"`
	BackImagePath    string             `bson:"back_image_path" json:"backImagePath" validate:"omitempty,file"`
	FrontText        string             `bson:"front_text" json:"frontText" validate:"required,min=5,max=128"`
	BackText         string             `bson:"back_text" json:"backText" validate:"required,min=5,max=128"`
	Question         string             `bson:"question" json:"question" validate:"min=5,max=32"`
	FlashCardCoverId primitive.ObjectID `bson:"flashcard_cover_id" json:"flashcardCoverId" validate:"required"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at" validate:"required"`
}

type FlashcardByIdRequest struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}
