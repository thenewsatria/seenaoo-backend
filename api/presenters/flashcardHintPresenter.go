package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardHint struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	HintText string             `json:"hintText" bson:"hint_text"`
}

func FlashcardHintInsertSuccessResponse(fh *models.FlashcardHint) *fiber.Map {
	flashcardHint := FlashcardHint{
		ID:       fh.ID,
		HintText: fh.HintText,
	}
	return &fiber.Map{
		"status": true,
		"data":   flashcardHint,
		"error":  nil,
	}
}
