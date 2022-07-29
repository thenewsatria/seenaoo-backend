package presenters

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlashcardHint struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	HintText string             `json:"hintText" bson:"hint_text"`
}
