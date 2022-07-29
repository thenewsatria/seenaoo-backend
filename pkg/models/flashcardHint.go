package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlashcardHint struct {
	ID          primitive.ObjectID
	HintText    string
	FlashcardId primitive.ObjectID
}
