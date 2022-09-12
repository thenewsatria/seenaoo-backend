package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	TagName string             `bson:"tag_name" json:"tagName"`
}

func TagSuccessResponse(t *models.Tag) *fiber.Map {
	tag := &Tag{
		ID:      t.ID,
		TagName: t.TagName,
	}
	return &fiber.Map{
		"success": true,
		"data":    tag,
		"error":   nil,
	}
}
