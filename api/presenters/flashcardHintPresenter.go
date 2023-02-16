package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardHint struct {
	ID       primitive.ObjectID `json:"id"`
	HintText string             `json:"hintText"`
}

func FlashcardHintSuccessResponse(fh *models.FlashcardHint) *fiber.Map {
	flashcardHint := FlashcardHint{
		ID:       fh.ID,
		HintText: fh.HintText,
	}
	return &fiber.Map{
		"status": "success",
		"data":   flashcardHint,
	}
}
