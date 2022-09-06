package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardCover struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Slug        string               `bson:"slug" json:"slug"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Image_path  string               `bson:"image_path" json:"imagePath"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Author      primitive.ObjectID   `bson:"user_id" json:"userId"`
	CreatedAt   time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updatedAt"`
}

func FlashcardCoverSuccessResponse(fcCover *models.FlashcardCover) *fiber.Map {
	flashcardCvr := &FlashcardCover{
		ID:          fcCover.ID,
		Slug:        fcCover.Slug,
		Title:       fcCover.Title,
		Description: fcCover.Description,
		Image_path:  fcCover.Image_path,
		Tags:        fcCover.Tags,
		Author:      fcCover.Author,
		CreatedAt:   fcCover.CreatedAt,
		UpdatedAt:   fcCover.UpdatedAt,
	}
	return &fiber.Map{
		"success": true,
		"data":    flashcardCvr,
		"error":   nil,
	}
}
