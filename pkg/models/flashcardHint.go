package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardHint struct {
	ID               primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	HintText         string             `bson:"hint_text" json:"hintText" validate:"required,min=5,max=32"`
	FlashcardCoverId primitive.ObjectID `bson:"flashcard_cover_id" json:"flashcardCoverId" validate:"required"`
	FlashcardId      primitive.ObjectID `bson:"flashcard_id" json:"flashcardId" validate:"required"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type FlashcardHintByIdRequest struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}
