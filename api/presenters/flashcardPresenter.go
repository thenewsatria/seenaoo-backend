package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flashcard struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	FrontImagePath   string             `bson:"front_image_path" json:"frontImagePath"`
	BackImagePath    string             `bson:"back_image_path" json:"backImagePath"`
	FrontText        string             `bson:"front_text" json:"frontText"`
	BackText         string             `bson:"back_text" json:"backText"`
	Question         string             `bson:"question" json:"question"`
	FlashCardCoverId primitive.ObjectID `bson:"flashcard_cover_id" json:"flashcardCoverId"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

type FlashcardDetail struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	FrontImagePath   string             `bson:"front_image_path" json:"frontImagePath"`
	BackImagePath    string             `bson:"back_image_path" json:"backImagePath"`
	FrontText        string             `bson:"front_text" json:"frontText"`
	BackText         string             `bson:"back_text" json:"backText"`
	Question         string             `bson:"question" json:"question"`
	FlashCardCoverId primitive.ObjectID `bson:"flashcard_cover_id" json:"flashcardCoverId"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
	Hints            []FlashcardHint    `json:"hints" bson:"hints"`
}

func FlashcardSuccessResponse(f *models.Flashcard) *fiber.Map {
	flashcard := &Flashcard{
		ID:               f.ID,
		FrontImagePath:   f.FrontImagePath,
		BackImagePath:    f.BackImagePath,
		FrontText:        f.FrontText,
		BackText:         f.BackText,
		Question:         f.Question,
		FlashCardCoverId: f.FlashCardCoverId,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
	return &fiber.Map{
		"status": "success",
		"data":   flashcard,
	}
}

func FlashcardDetailSuccessResponse(f *models.Flashcard, hints *[]models.FlashcardHint) *fiber.Map {
	flashcard := &FlashcardDetail{
		ID:               f.ID,
		FrontImagePath:   f.FrontImagePath,
		BackImagePath:    f.BackImagePath,
		FrontText:        f.FrontText,
		BackText:         f.BackText,
		Question:         f.Question,
		FlashCardCoverId: f.FlashCardCoverId,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
		Hints:            []FlashcardHint{},
	}

	for _, v := range *hints {
		flashcardHint := &FlashcardHint{
			ID:       v.ID,
			HintText: v.HintText,
		}

		flashcard.Hints = append(flashcard.Hints, *flashcardHint)
	}

	return &fiber.Map{
		"status": "success",
		"data":   flashcard,
	}
}
