package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardHint struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	HintText    string             `bson:"hint_text" json:"hintText"`
	FlashcardId primitive.ObjectID `bson:"flashcard_id" json:"flashcardId"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReadFlashcardHintRequest struct {
	ID string `json:"id"`
}
